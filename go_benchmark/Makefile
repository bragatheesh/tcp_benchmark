all:
	#go tool cgo benchmark.go
	go build benchmark_server.go
	go build benchmark_client.go

server:
	./benchmark_server 9090

client:
	./benchmark_client localhost 9090 50

check:
	@- ./benchmark_server 9090 & ./benchmark_client localhost 9090 5||:

clean:
	rm -rf benchmark_server benchmark_client
	
