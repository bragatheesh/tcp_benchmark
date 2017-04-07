#include <stdio.h>
#include <string.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <stdlib.h>
#include <unistd.h>
#include <arpa/inet.h>
#include <netdb.h>
#include <fcntl.h>
#include <sys/time.h>

#define BUFFSIZE 1024
//#define NUM_TESTS 10

int
main(int argc, char **argv)
{
    if(argc < 3){
        printf("Usage ./client hostname port [number_of_tests default: 10]");
        return -1;
    }
    int port_number;
    int server_sock;
    int sockaddr_size;
    int ret = 0;
    int loops = 0;
    int size;
    int NUM_TESTS = 10;
    unsigned int secs = 0;
    unsigned int usecs = 0;
    unsigned int sent_size = 0;
    unsigned int recvd_size = 0;
    char* message = "Hello from client";
    char recv_buff[BUFFSIZE];
    struct sockaddr_in server;
    struct hostent *host;
    struct timeval timeout;
    struct timeval start;
    struct timeval end;

    port_number = atoi(argv[2]);

    if(argc == 4){
        NUM_TESTS = atoi(argv[3]);
    }

    host = gethostbyname(argv[1]);
    if(host == NULL){
        printf("Host %s could not be found through gethostbyname\n", argv[1]);
        return -1;
    }

    //CreateSocket(PktIOSocket) (uint32, error)
    server_sock = socket(AF_INET, SOCK_STREAM, 0);
    if(server_sock < 0){
        printf("Error creating socket\n");
        return -1;
    }

    timeout.tv_sec = 5;
    timeout.tv_usec = 0;
    //SetSockoptTCPMD5(uint32, string, string)
    if(setsockopt (server_sock, SOL_SOCKET, SO_RCVTIMEO, (char *)&timeout, 
                                                     sizeof(timeout)) < 0){
        printf("setsockopt(SO_RCVTIMEO) failed\n");
        //DeleteSocket(uint32) error
        close(server_sock);
        return -1;
    }

    sockaddr_size = sizeof(server);
    bzero(&server, sockaddr_size);
    server.sin_family = AF_INET;
    bcopy((char *)host->h_addr, (char *)&server.sin_addr.s_addr, host->h_length);
    server.sin_port = htons(port_number);

    //Connect(uint32, PktIOSocket, time.Duration) error
    if(connect(server_sock, (struct sockaddr *)&server, sockaddr_size) < 0){
        printf("Failed to connect to %s\n", argv[1]);
        //DeleteSocket(uint32) error
        close(server_sock);
        return -1;
    }

    while(loops < NUM_TESTS){

        memset(&recv_buff, '\0', BUFFSIZE);
        size = strlen(message);
        gettimeofday(&start, NULL);

        //Write([]byte) (int, error)
        if(send(server_sock, message, size, 0) < 0){
            printf("Could not send message\n");
            //DeleteSocket(uint32) error
            close(server_sock);
            return -1;
        }

        //Read([]byte) (int, error)
        ret = recv(server_sock, recv_buff, BUFFSIZE, 0);
        gettimeofday(&end, NULL);
        if(ret < 0){
            printf("Error timed out on recv\n");
            //DeleteSocket(uint32) error
            close(server_sock);
            return -1;
        }
        //printf("Message from server: %s\n", recv_buff);
        secs += end.tv_sec - start.tv_sec;
        usecs += end.tv_usec - start.tv_usec;
        sent_size += size;
        recvd_size += (unsigned int)strlen(recv_buff);
        loops++;
        
    }
    //DeleteSocket(uint32) error
    close(server_sock);
    printf("Sent total of %u bytes and recvd total of %u bytes\nAvg RTT: %d seconds %d usec\n", sent_size, recvd_size, secs/NUM_TESTS, usecs/NUM_TESTS);
    return 1;
}
