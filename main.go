// +build !linux

// tnc-server
// A serial/TCP bridge for connecting multiple read/write clients to an AX.25/KISS TNC device.
// Written in the Go programming language
// (c) 2014, Christopher Snell
//
// main.go - Initialization functions for every OS besides Linux

package main

import (
	"flag"
	"github.com/tarm/goserial"
	"io"
	"log"
	"os"
	"os/signal"
)

func main() {

	var s io.ReadWriteCloser
	var err error

	port := flag.String("port", "", "Serial port device (e.g. /dev/ttyUSB0, COM1, etc.)")
	baud := flag.Int("baud", 4800, "Baud rate for serial device (default: 4800")
	listen = flag.String("listen", "0.0.0.0:6700", "Address/port to listen on (defaults to 0.0.0.0:6700)")
	debug = flag.Bool("debug", false, "Enable debugging information (default: false)")
	flag.Parse()

	// Spin off a goroutine to watch for a SIGINT and die if we get one
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	if len(*port) == 0 {
		log.Fatalln("Must specify a serial port with -port flag.  Use -h flag for help")
	}

	sc := &serial.Config{Name: *port, Baud: *baud}

	s, err = serial.OpenPort(sc)
	if err != nil {
		log.Fatal(err)
	}

	defer s.Close()

	go newSerialListener(s)

	<-sig
	log.Println("SIGINT received.  Shutting down...")
	os.Exit(1)
}
