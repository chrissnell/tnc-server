tnc-server
==========

tnc-server is a multiplexing network server for KISS-enabled Amateur Radio packet terminal node controllers (TNCs).   It provides a way to share a TNC amongst multiple read/write, read-only, and write-only clients.   tnc-server attaches to a serial port and sends all received KISS messages to all connected network clients.   The clients talk to tnc-server over TCP and can run locally (on the same machine that's attached to the TNC) or remote (across the Internet).  
