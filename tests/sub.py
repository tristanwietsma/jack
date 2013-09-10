import sys
import socket

HOST = '0.0.0.0'
PORT = 2000

s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
s.connect((HOST, PORT))
s.send("SUB testKey")
while True:
	inc = s.recv(1024)
	print "subscribe...", inc
s.close()
