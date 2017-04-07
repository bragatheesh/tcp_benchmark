all:
	gcc -g client.c -o client
	gcc -g server.c -o server

check: all
	@- ./server 9090 20 & ./client localhost 9090 20 ||:

clean:
	rm -rf client server
