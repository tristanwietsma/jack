import sys
import socket

HOST = '0.0.0.0'
PORT = 2000

s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
s.connect((HOST, PORT))

for i in range(200):
	s.send("SET key 123")
	s.recv(1024)

s.close()
