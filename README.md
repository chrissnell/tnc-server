# tnc-server
### [Download tnc-server](https://github.com/chrissnell/tnc-server/releases)
**tnc-server** is a multiplexing network server for KISS-enabled Amateur Radio packet terminal node controllers (TNCs).   It provides a way to share a TNC amongst multiple read/write, read-only, and write-only clients.   **tnc-server** connects to a TNC via serial port or TCP and sends all received KISS messages to all connected network clients.   The clients talk to **tnc-server** over TCP and can run locally (on the same machine that's attached to the TNC) or remotely (across the Internet).

The key difference between **tnc-server** and other remote serial software is that **tnc-server** understands AX.25 and is designed to allow many simultaneous client connections.  Packets sent to **tnc-server** are written to the TNC in a first-in first-out manner and will not clobber each other.  Likewise, incoming RF traffic through the TNC will be distributed to all connected clients.

tnc-server is written in the [Go Programming Language](http://golang.org/)

## Using tnc-server

You will need a KISS-capable TNC connected via serial port or accessible over TCP (e.g. Direwolf, UZ7HO soundmodem).  tnc-server does not currently support the "TNC2" protocol.

### Linux and Mac OS X
Download the appropriate **tnc-server** package for your architecture from the [releases page](https://github.com/chrissnell/tnc-server/releases).
```
Usage:
./tnc-server [-port=/path/to/serialdevice] [-baud=BAUDRATE] [-listen=IPADDRESS:PORT]

-port - serial device or tcp:host:port for network KISS.  Default: /dev/ttyUSB0
        Examples:
          -port=/dev/ttyUSB0           (serial)
          -port=tcp:192.168.1.100:8001 (TCP, e.g. Direwolf)

-baud - the baud rate for the serial device (ignored for TCP).  Default: 4800

-listen - the IPADDRESS:PORT to listen for incoming client connections.  Default: 0.0.0.0:6700  (all IPs on port 6700)

-debug - enable debug output with hex dumps of received frames
```

On Linux you can have `systemd` automatically start `tnc-server` when the machines boots. Put this
in `/etc/systemd/system/tnc-server.service`:
```
# /etc/systemd/system/tnc-server.service
[Unit]
Description=TNC Server - Multiplexing network server for KISS TNCs
Documentation=https://github.com/chrissnell/tnc-server
After=network.target

[Service]
Type=simple
EnvironmentFile=/etc/default/tnc-server
ExecStart=/usr/local/bin/tnc-server -port=${TNC_DEVICE} -baud=${TNC_BAUD} -listen=${TNC_LISTEN}
Restart=on-failure
RestartSec=5

# Hardening
NoNewPrivileges=true
ProtectSystem=strict
ProtectHome=true
PrivateTmp=true
ProtectKernelTunables=true
ProtectControlGroups=true
# Allow access to serial devices
SupplementaryGroups=dialout

[Install]
WantedBy=multi-user.target
```

Then put this configuration file in `/etc/default/tnc-server`:
```
# /etc/default/tnc-server
# Configuration for the tnc-server systemd service

# Serial device path, or tcp:host:port for network KISS (e.g. Direwolf)
# Examples:
#   /dev/ttyACM0              (serial)
#   tcp:192.168.1.100:8001    (TCP KISS, e.g. Direwolf)
TNC_DEVICE=/dev/ttyACM0

# Baud rate for serial device (ignored for TCP connections)
TNC_BAUD=1200

# IP address and port to listen for incoming client connections
TNC_LISTEN=0.0.0.0:6700
```

Adjust the settings in the config file as necessary. Then:
```
$ sudo systemctl daemon-reload
$ sudo systemctl enable --now tnc-server
```

### Windows
Download the appropriate **tnc-server** package for your architecture from the [releases page](https://github.com/chrissnell/tnc-server/releases).   See below for virtual COM port emulation, if you plan on running a Windows-based APRS client.
```
Usage:

Open a command-prompt in the directory where you have the tnc-server.exe binary and run it like this:

tnc-server.exe [-port=COMnn] [-baud=BAUDRATE] [-listen=IPADDRESS:PORT]

-port - serial device or tcp:host:port for network KISS.  Default: COM1
        Examples:
          -port=COM3                    (serial)
          -port=tcp:192.168.1.100:8001  (TCP, e.g. Direwolf)

-baud - the baud rate for the serial device (ignored for TCP).  Default: 4800

-listen - the IPADDRESS:PORT to listen for incoming client connections.  Default: 0.0.0.0:6700  (all IPs on port 6700)
```

## TCP KISS (Direwolf, soundmodem, etc.)

If your TNC exposes a KISS interface over TCP rather than a serial port, you can connect directly without needing socat or other bridging tools:

```
./tnc-server -port=tcp:192.168.1.100:8001 -listen=:6700
```

If the TCP connection to the TNC is lost (e.g. Direwolf restarts), tnc-server will automatically reconnect every 5 seconds until the TNC is available again. Connected clients remain attached through the reconnection.

## Using tnc-server with aprx
**tnc-server** works very nicely with [aprx](http://wiki.ham.fi/Aprx.en) using aprx's KISS-over-TNC feature.   To use it, simply include a stanza like this in your aprx.conf, substituting your own callsign and optional SSID, and the IP address of your tnc-server:

```
<interface>
  tcp-device 127.0.0.1 6700 KISS
  callsign YOURCALL-SSID
  tx-ok true
</interface>
```

If you're running aprx on the same machine as **tnc-server**, using 127.0.0.1 as the IP address.   Otherwise, use your machine's IP address here.

## Using tnc-server   with APRSISCE/32
**tnc-server** plays nicely with APRSISCE/32.  Start APRSISCE/32, navigate to Configure -> Ports -> New Port... and choose Simply KISS.   Choose TCP/IP as the port type and fill in the IP and port of your **tnc-server** instance.

## Using with Xastir
To use **tnc-server** with Xastir, you will need to download and install [remserial](http://lpccomp.bc.ca/remserial/).   You'll run remserial and give it the address of your **tnc-server**, as well as the local pseudo-tty (Linux version of virtual serial ports) that Xastir will attach to.

Example:

```
% sudo ./remserial -r 10.50.0.25 -p 6700 -s "4800" -l /dev/remserial1 /dev/ptmx
% sudo chmod 666 /dev/remserial1
```

In this example, we're connecting to a TNC server at IP 10.50.0.25 (port 6700) at 4800 baud and mapping that back to /dev/remserial1.   Then we're running chmod to make that virtual serial port read/write accessible to non-root users (you).

Next, fire up Xastir and navigate to the Interface Control menu.  Create a new interface (type: **Serial KISS TNC**) with /dev/remserial1 as the **TNC Port**.  Set your port baud rate to **4800** and choose the iGating options that you want.  Check "Allow Transmitting" if you want Xastir to transmit.  Choose a reasonable APRS digipeater path for your area.   Leave the KISS parameters in their default settings and click **Ok**.   Go back to Interface Control, select your new interface and click the Start button.  It should start hearing stations off the air at this point.

## Using tnc-server to debug/develop APRS clients
It is possible to use **tnc-server** and the **socat** utility to create a virtual null-modem connection that allows you to debug or develop an APRS/AX.25 client without the need for multiple physical RS-232 ports and without a null-modem cable.  For more details, see [this blog post](http://output.chrissnell.com/post/94364500380/debugging-aprs-clients-with-a-virtual-null-modem-cable).

## Windows Virtual COM port
You don't need to install a virtual COM port to run **tnc-server** on Windows.   However, if you want to use Windows-based APRS software that expects a COMn port (like COM1, etc), you'll need to use com2tcp from the [com0com project](http://com0com.sourceforge.net/).

To get this working, download com0com [here](http://sourceforge.net/projects/com0com/files/com0com/3.0.0.0/com0com-3.0.0.0-i386-and-x64-unsigned.zip/download).  Windows 7 users, download the signed version of com0com [here](https://code.google.com/p/powersdr-iq/downloads/detail?name=setup_com0com_W7_x64_signed.exe&can=2&q=).

Once you have this package installed, you'll run com2tcp like this:

```
    com2tcp \\.\CNCB0 127.0.0.1 6700
```

You'll want to substitute the IP address of your **tnc-server**.  CNCB0 refers to COM2 in com0com parlance.   For more info on what to put here, check out the [README file for com0com](http://com0com.cvs.sourceforge.net/viewvc/com0com/com0com/ReadMe.txt?revision=RELEASED) and the [README file for com2tcp](http://com0com.cvs.sourceforge.net/*checkout*/com0com/com2tcp/ReadMe.txt?revision=RELEASED).

## TNCs known to work with tnc-server
- [Argent Data Tracker2](http://www.argentdata.com/products/tracker2.html)
- Kenwood TH-D75 handheld radio

If you've tested **tnc-server** with another TNC, let me know and I will add it to this list.


## Building from source

```
go build
```

## License
```
            DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE
                    Version 2, December 2004

 Copyright (C) 2017 Chris Snell <chris@chrissnell.com>

 Everyone is permitted to copy and distribute verbatim or modified
 copies of this license document, and changing it is allowed as long
 as the name is changed.

            DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE
   TERMS AND CONDITIONS FOR COPYING, DISTRIBUTION AND MODIFICATION

  0. You just DO WHAT THE FUCK YOU WANT TO.
```
