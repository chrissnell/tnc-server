# tnc-server

tnc-server is a multiplexing network server for KISS-enabled Amateur Radio packet terminal node controllers (TNCs).   It provides a way to share a TNC amongst multiple read/write, read-only, and write-only clients.   tnc-server attaches to a serial port and sends all received KISS messages to all connected network clients.   The clients talk to tnc-server over TCP and can run locally (on the same machine that's attached to the TNC) or remote (across the Internet).  

## Using tnc-server

You will need a computer with a serial port that's attached to a TNC that is able to speak the KISS protocol.   tnc-server does not currently support the "TNC2" protocol.  

### Linux
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
./tnc-server [-port=/path/to/serialdevice] [-baud=BAUDRATE] [-listen=IPADDRESS:PORT]

-port - the serial device where the KISS TNC is attached.  Default: COM1

-baud - the baudrate to talk to the TNC  Default: 4800

-listen - the IPADDRES:PORT to listen for incoming client connections.  Default: 0.0.0.0:6700  (all IPs on port 6700)

```
