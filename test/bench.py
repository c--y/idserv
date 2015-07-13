# coding=utf-8

import socket
import struct

packet = 'me'

N = 10000000

HOST = 'localhost'
PORT = 8000

def bench():
    socket.setdefaulttimeout(0.1)
    for i in xrange(N):
        try:
            sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
            sock.connect((HOST, PORT))
            #sock.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)

            sock.send(packet)
            ret = sock.recv(8)
            print struct.unpack('<Q', ret)[0] 
            sock.close()
        except:
            continue

if __name__ == '__main__':
    bench()
