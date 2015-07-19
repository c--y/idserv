# coding=utf-8

import socket
import struct
import time

packet = 'me'

N = 100000

HOST = 'localhost'
PORT = 8000

def bench():
    t0 = time.time()
    socket.setdefaulttimeout(0.1)
    sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    sock.connect((HOST, PORT))
    for i in xrange(N):
        try:
            #sock.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
            sock.send(packet)
            ret = sock.recv(8)
            #print struct.unpack('<Q', ret)[0] 
        except:
            continue

    sock.close()
    t1 = time.time()
    print '%d times generate: %d ms' % (N, t1 * 1000 - t0 * 1000)

if __name__ == '__main__':
    bench()
