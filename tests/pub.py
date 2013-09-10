import time
import socket

HOST = '0.0.0.0'
PORT = 2000

s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
s.connect((HOST, PORT))
s.send("PUB testKey")

response = s.recv(1024)
print response

while True:
	msg = str(time.time())
	print "publish...", msg
	s.send(msg)
	time.sleep(1)

s.close()
