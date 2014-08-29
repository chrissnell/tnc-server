// tnc-server
// A serial/TCP bridge for connecting multiple read/write clients to an AX.25/KISS TNC device.
// Written in the Go programming language
// (c) 2014, Christopher Snell

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/tv42/topic"
	"io"
	"log"
	"net"
)

var (
	listen *string
	debug  *bool
)

// The smallest KISS packet that we can ever expect to see.
const reasonableSize = 15

// This function sets up our TCP listener and our serial port ReadWriteCloser
// and handles incoming connections.
func newSerialListener(serialport io.ReadWriteCloser) {

	// We're going to use Topic to handle the one -> many distribution
	// of our serial->net traffic
	top := topic.New()
	defer close(top.Broadcast)

	// We use a channel as a FIFO buffer to handle the many -> one
	// writes from our net->serial traffic
	msg := make(chan []byte, 15)

	l, err := net.Listen("tcp", *listen)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	log.Printf("Listening for connections on %v\n", *listen)

	// Launch a broadcaster in a goroutine that reads off the serial port and sends to Topic
	go serialReaderBroadcaster(top, serialport)

	// Launch a writer in a goroutine to receive our incoming writes and write them to serial
	go serialWriter(serialport, msg)

	for {
		// Wait for a connection.
		conn, err := l.Accept()
		log.Printf("Answered incoming client connection from %v\n", conn.RemoteAddr())
		if err != nil {
			log.Fatal(err)
		}

		// Set up a consumer channel for this new connection.  All reads off the
		// serial port will be sent over this channel
		consumer := make(chan interface{}, 1)
		top.Register(consumer)

		// Start a consumer goroutine to take the serial data off the consumer
		// channel and send it to this network connection
		go serialReaderConsumer(consumer, conn, top)

		// Start a writer in a goroutine to read off the network and send all messages
		// to the message buffer
		go serialWriterConnection(conn, msg)
	}
}

// This function reads off the serial port and sends what it gets to Topic
func serialReaderBroadcaster(top *topic.Topic, serialport io.ReadWriteCloser) {

	// Wrap the goserial's ReadWriteCloser with a bufio.Reader so we can do fancy stuffs.
	sr := bufio.NewReader(serialport)

	for {

		frame := []byte{}

		framestarted := false

		bcount := 0

		//for len(frame) <= reasonableSize {
		for bcount <= 400 {

			c, err := sr.ReadByte()
			if err != nil {
				log.Printf("Error reading bytes from serial: %v\n", err)
			}

		readloop:
			switch c {
			case 0xE:
				// The first character read from an I2C TNC will be a 0xE.
				// It should be discarded, as it is not part of a frame.
			case 0x2C:
				// I2C Readers are handled a little differently than standard serial port Readers.
				// When reading an I2C filehandle, it will return a steady stream of 0x2C bytes until
				// some actual data is sent across.  For this reason, we test to see the AX.25 frame
				// has begun (i.e. we've received a 0xC0 FEND) and we only append the 0xE if it's part
				// of an actual frame
				if framestarted {
					// We've received a 0xE and a frame has already begun so it will be added
					// to our frame
					frame = append(frame, c)
					bcount++
					rs := c >> 1
					log.Printf("%4d \t%#x \t%v \t%v\t%08b\n", bcount, c, string(c), string(rs), c)
				}
			case 0xC0:

				log.Printf("FEND received.  Size of frame so far: %v", len(frame))

				if framestarted {

					log.Println("End FEND received")

					// If we receive a 0xC0 (FEND) and the frame has already started, this is probably
					// the end of our frame.  The outer loop will continue to run if our frame has not
					// reached reasonableSize.
					frame = append(frame, c)
					bcount++
					rs := c >> 1
					log.Printf("%4d \t%#x \t%v \t%v\t%08b\n", bcount, c, string(c), string(rs), c)
					break readloop
				} else {

					log.Println("Beginning FEND received")

					// We've received a 0xC0 (FEND) and the frame has not yet started so we know that
					// this marks the beginning of a new frame.  Append the byte to frame and set
					// framestarted to true.
					frame = append(frame, c)
					bcount++
					framestarted = true
					rs := c >> 1
					log.Printf("%4d \t%#x \t%v \t%v\t%08b\n", bcount, c, string(c), string(rs), c)
				}
			default:
				frame = append(frame, c)
				bcount++
				rs := c >> 1
				log.Printf("%4d \t%#x \t%v \t%v\t%08b\n", bcount, c, string(c), string(rs), c)
			}

		}

		// We've received a frame that's too long or too short.  Probably junk.  Discared and continue
		if bcount >= 400 || bcount < 15 {
			log.Printf("Breaking.  bcount: %v\n", bcount)
			continue
		}

		log.Println("Sending frame")

		if *debug {
			fmt.Println("Byte#\tHexVal\tChar\tChar>>1\tBinary")
			fmt.Println("-----\t------\t----\t-------\t------")
			for k, v := range frame {
				rs := v >> 1
				fmt.Printf("%4d \t%#x \t%v \t%v\t%08b\n", k, v, string(v), string(rs), v)
			}
		}

		// Send our received frame to Topic for distribution to the consumer(s)
		top.Broadcast <- frame

	}

}

// This function reads from the Topic consumer and writes what it gets to the network
func serialReaderConsumer(consumer chan interface{}, conn net.Conn, top *topic.Topic) {
	defer conn.Close()
	defer top.Unregister(consumer)

	for {
		select {
		// A new message was received from this Topic consumer
		case msg, ok := <-consumer:
			if ok {
				i, err := conn.Write(msg.([]byte))
				if err != nil {
					log.Printf("Error writing %v bytes to %v: %v\n", i, conn.RemoteAddr(), err)
					log.Println("Client hung up.  Closing connection.")
					conn.Close()
					return
				}
			} else {
				log.Printf("Channel closed (%v)", conn.RemoteAddr())
				break
			}
		}
	}
}

// This function reads off the network and sends everything it gets to the msg channel
// for consumption by serialWriter()
func serialWriterConnection(conn net.Conn, msg chan []byte) {
	var err error

	for {

		frame := []byte{}
		var first_byte byte

		var frame_buffer bytes.Buffer

		// Wrap a bufio.Reader around our net.Conn
		r := bufio.NewReader(conn)

		// Read our first byte, a 0xc0 and add it to the frame
		first_byte, err = r.ReadByte()
		if err != nil {
			log.Printf("Error reading bytes from %v: %v", conn.RemoteAddr(), err)
			log.Println("Client hung up.  Closing connection.")
			conn.Close()
			return
		}

		frame_buffer.WriteByte(first_byte)

		for len(frame) <= reasonableSize {

			// Read until we see a 0xc0, and store this in the frame (including that 0xc0 byte)
			frame, err = r.ReadBytes(byte(0xc0))

			if *debug {
				fmt.Println("Byte#\tHexVal\tChar\tChar>>1\tBinary")
				fmt.Println("-----\t------\t----\t-------\t------")
				for k, v := range frame {
					rs := v >> 1
					fmt.Printf("%4d \t%#x \t%v \t%v\t%08b\n", k, v, string(v), string(rs), v)
				}
			}

			if err != nil {
				log.Printf("Error reading bytes from %v: %v", conn.RemoteAddr(), err)
				log.Println("Client hung up.  Closing connection.")
				conn.Close()
				return
			}
		}

		frame_buffer.Write(frame)

		frame_out := frame_buffer.Bytes()

		// Send the frame we just read off the network to the message buffer for eventual
		// write to serial
		msg <- frame_out

	}

}

// This function reads frames off the buffered msg channel and writes them to serial
func serialWriter(s io.ReadWriteCloser, msg chan []byte) {
	for {
		select {
		case msg, ok := <-msg:
			if ok {
				_, err := s.Write(msg)
				if err != nil {
					log.Printf("Unable to write to serial port: %v\n", err)
				}
			} else {
				log.Println("serialWriter message channel closed.")
				break
			}
		}
	}
}
