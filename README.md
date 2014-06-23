# tnc-server

**tnc-server** is a multiplexing network server for KISS-enabled Amateur Radio packet terminal node controllers (TNCs).   It provides a way to share a TNC amongst multiple read/write, read-only, and write-only clients.   **tnc-server** attaches to a serial port and sends all received KISS messages to all connected network clients.   The clients talk to **tnc-server** over TCP and can run locally (on the same machine that's attached to the TNC) or remotely (across the Internet).  

**tnc-server** is designed to allow many simultaneous client connections.  An example use case might be an Emergency Operations Center (EOC) that has multiple APRS users but only one TNC and radio transmitter.   In this scenario, they would run **tnc-server** on a server or workstation attached to the TNC (could be as simple and inexpensive as a [Raspberry Pi](http://raspberrypi.org)) and connect that machine on a LAN with the other workstations.   **tnc-server** can handle multiple types of workstations simultaneously.   You could connect workstations running APRSISCE/32 running on Windows, Xastir on Linux, and a software digipeater running aprx on a Mac simultaneously.

tnc-server is written in the [Go Programming Language](http://golang.org/)

## Using tnc-server

You will need a computer with a serial port that's attached to a TNC that is able to speak the KISS protocol.   tnc-server does not currently support the "TNC2" protocol.  

### Linux and Mac OS X
Download the appropriate **tnc-server** package for your architecture from the one of the links below. 
```
Usage:
./tnc-server [-port=/path/to/serialdevice] [-baud=BAUDRATE] [-listen=IPADDRESS:PORT]

-port - the serial device where the KISS TNC is attached.  Default: /dev/ttyUSB0

-baud - the baudrate to talk to the TNC  Default: 4800

-listen - the IPADDRES:PORT to listen for incoming client connections.  Default: 0.0.0.0:6700  (all IPs on port 6700)

```

### Windows
Download the appropriate **tnc-server** package for your architecture from the one of the links below.   See below for virtual COM port emulation, if you plan on running a Windows-based APRS client.
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
**tnc-server** works very nicely with [aprx](http://wiki.ham.fi/Aprx.en) using aprx's KISS-over-TNC feature.   To use it, simply include a stanza like this in your aprx.conf, substituting your own callsign and optional SSID, and the IP address of your tnc-server:

```
<interface>
  tcp-device 127.0.0.1 6700 KISS
  callsign YOURCALL-SSID
  tx-ok true
</interface>
```

If you're running aprx on the same machine as **tnc-server**, using 127.0.0.1 as the IP address.   Otherwise, use your machine's IP address here.

## Using with Xastir
To use **tnc-server** with Xastir, you will need to download and install [remserial](http://lpccomp.bc.ca/remserial/).   You'll run remserial and give it the address of your **tnc-server**, as well as the local pseudo-tty (Linux version of virtual serial ports) that Xastir will attach to.

Example:

```
% sudo ./remserial -r 10.50.0.25 -p 6700 -s "4800" -l /dev/remserial1 /dev/ptmx
% sudo chmod 666 /dev/remserial1
```

In this example, we're connecting to a TNC server at IP 10.50.0.25 (port 6700) at 4800 baud and mapping that back to /dev/remserial1.   Then we're running chmod to make that virtual serial port read/write accessible to non-root users (you).

Next, fire up Xastir and navigate to the Interface Control menu.  Create a new interface (type: **Serial KISS TNC**) with /dev/remserial1 as the **TNC Port**.  Set your port baud rate to **4800** and choose the iGating options that you want.  Check "Allow Transmitting" if you want Xastir to transmit.  Choose a reasonable APRS digipeater path for your area.   Leave the KISS parameters in their default settings and click **Ok**.   Go back to Interface Control, select your new interface and click the Start button.  It should start hearing stations off the air at this point.

## Windows Virtual COM port
You don't need to install a virtual COM port to run **tnc-server** on Windows.   However, if you want to use Windows-based APRS software that expects a COMn port (like COM1, etc), you'll need to use com2tcp from the [com0com project](http://com0com.sourceforge.net/).

To get this working, download com0com [here](http://sourceforge.net/projects/com0com/files/com0com/3.0.0.0/com0com-3.0.0.0-i386-and-x64-unsigned.zip/download).  Windows 7 users, download the signed version of com0com [here](https://code.google.com/p/powersdr-iq/downloads/detail?name=setup_com0com_W7_x64_signed.exe&can=2&q=).  

Once you have this package installed, you'll run com2tcp like this:

```
    com2tcp \\.\CNCB0 127.0.0.1 6700
```

You'll want to substitute the IP address of your **tnc-server**.  CNCB0 refers to COM2 in com0com parlance.   For more info on what to put here, check out the [README file for com0com](http://com0com.cvs.sourceforge.net/viewvc/com0com/com0com/ReadMe.txt?revision=RELEASED) and the [README file for com2tcp](http://com0com.cvs.sourceforge.net/*checkout*/com0com/com2tcp/ReadMe.txt?revision=RELEASED).

## TNCs known to work with tnc-server
[Argent Data Tracker2](http://www.argentdata.com/products/tracker2.html)

If you've tested **tnc-server** with another TNC, let me know and I will add it to this list.


## Building your own binaries
If you want to modify tnc-server and build your own binaries, you'll need a working installation of the [Go Programming Language](http://golang.org).  Once you have that...

```
% go get github.com/tarm/goserial
% go get github.com/tv42/topic
% go build
```
