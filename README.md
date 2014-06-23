# tnc-server

tnc-server is a multiplexing network server for KISS-enabled Amateur Radio packet terminal node controllers (TNCs).   It provides a way to share a TNC amongst multiple read/write, read-only, and write-only clients.   tnc-server attaches to a serial port and sends all received KISS messages to all connected network clients.   The clients talk to tnc-server over TCP and can run locally (on the same machine that's attached to the TNC) or remote (across the Internet).  

tnc-server is written in the [Go Programming Language](http://golang.org/)

## Using tnc-server

You will need a computer with a serial port that's attached to a TNC that is able to speak the KISS protocol.   tnc-server does not currently support the "TNC2" protocol.  

### Linux and Mac OS X
Download the appropriate tnc-server package for your architecture from the one of the links below. 
```
Usage:
./tnc-server [-port=/path/to/serialdevice] [-baud=BAUDRATE] [-listen=IPADDRESS:PORT]

-port - the serial device where the KISS TNC is attached.  Default: /dev/ttyUSB0

-baud - the baudrate to talk to the TNC  Default: 4800

-listen - the IPADDRES:PORT to listen for incoming client connections.  Default: 0.0.0.0:6700  (all IPs on port 6700)

```

### Windows
Download the appropriate tnc-server package for your architecture from the one of the links below. 
```
Usage:

Open a command-prompt in the directory where you have the tnc-server.exe binary and run it like this:

tnc-server.exe [-port=/path/to/serialdevice] [-baud=BAUDRATE] [-listen=IPADDRESS:PORT]

-port - the serial device where the KISS TNC is attached.  Default: COM1

-baud - the baudrate to talk to the TNC  Default: 4800

-listen - the IPADDRES:PORT to listen for incoming client connections.  Default: 0.0.0.0:6700  (all IPs on port 6700)

```

## Download

Linux AMD/Intel 64-bit:  http://island.nu/tnc-server/tnc-server-linux-amd64.tar.gz

Linux ARMv7 (BeagleBone/BeagleBoard):  http://island.nu/tnc-server/tnc-server-linux-arm7.tar.gz

Linux ARMv6 (Raspberry Pi, etc.):  http://island.nu/tnc-server/tnc-server-linux-armv6.tar.gz

Mac OS X:  http://island.nu/tnc-server/tnc-server-darwin-amd64.tar.gz

Windows 32-bit: http://island.nu/tnc-server/tnc-server-win32.zip

Windows 64-bit: http://island.nu/tnc-server/tnc-server-win64.zip

## Using tnc-server with aprx
tnc-server works very nicely with [aprx](http://wiki.ham.fi/Aprx.en) using aprx's KISS-over-TNC feature.   To use it, simply include a stanza like this in your aprx.conf, substituting your own callsign and optional SSID, and the IP address of your tnc-server:

```
<interface>
  tcp-device 127.0.0.1 6700 KISS
  callsign YOURCALL-SSID
  tx-ok true
</interface>
```

If you're running aprx on the same machine as tnc-server, using 127.0.0.1 as the IP address.   Otherwise, use your machine's IP address here.

## TNCs known to work with tnc-server
[Argent Data Tracker2](http://www.argentdata.com/products/tracker2.html)

If you've tested tnc-server with another TNC, let me know and I will add it to this list.


## Building your own binaries
If you want to modify tnc-server and build your own binaries, you'll need a working installation of the [Go Programming Language](http://golang.org).  Once you have that...

```
% go get github.com/tarm/goserial
% go get github.com/tv42/topic
% go build
```
