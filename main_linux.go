// +build linux

// tnc-server
// A serial/TCP bridge for connecting multiple read/write clients to an AX.25/KISS TNC device.
// Written in the Go programming language
// (c) 2014, Christopher Snell
//
// main_linux.go - Initialization functions for Linux (including I2C support)

package main

import (
	"flag"
	"github.com/chrissnell/i2c"
	"github.com/tarm/goserial"
	"io"
	"log"
	"os"
	"os/signal"
)

func main() {

	var s io.ReadWriteCloser
	var err error

	port := flag.String("port", "", "Serial port device (e.g. /dev/ttyUSB0)")
	baud := flag.Int("baud", 4800, "Baud rate for serial device (default: 4800")
	i2cbus := flag.Int("i2cbus", 0, "I2C bus number (0, 1, 2, etc.)")
	i2caddr := flag.Uint("i2caddr", 0, "I2C device address (e.g. 0x27)")
	listen = flag.String("listen", "0.0.0.0:6700", "Address/port to listen on (defaults to 0.0.0.0:6700)")
	debug = flag.Bool("debug", false, "Enable debugging information (default: false)")
	flag.Parse()

	// Spin off a goroutine to watch for a SIGINT and die if we get one
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	if len(*port) > 0 {

		sc := &serial.Config{Name: *port, Baud: *baud}

		s, err = serial.OpenPort(sc)
		if err != nil {
			log.Fatal(err)
		}

	} else if *i2caddr != 0 {
		s, err = i2c.New(uint8(*i2caddr), *i2cbus)
		if err != nil {
			log.Fatalf("Error opening I2C address %v on bus %v: %v\n", *i2caddr, *i2cbus, err)
		}
	} else {
		log.Fatalln("Must pass -port argument (for serial) or -i2cbus and -i2caddr arguments (for I2C).   Use -h flag for help.")
	}

	defer s.Close()

	go newSerialListener(s)

	<-sig
	log.Println("SIGINT received.  Shutting down...")
	os.Exit(1)
}
