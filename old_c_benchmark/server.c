#include <stdio.h>
#include <string.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <arpa/inet.h>
#include <stdlib.h>
#include <unistd.h>
#include <netdb.h>
#include <signal.h>
#include <fcntl.h>

//#define BUFFSIZE 1024
//#define NUM_TESTS 10

struct Endpoint{
    char* Vrf;
    char* IPAdress;
    uint16_t L4Port;
    char* Zone;
};

struct PktIOSocket{
    uint8_t Type;
    uint16_t Vrf;
    uint16_t Port;
    uint16_t Vlan;
    uint16_t EthType;
    uint16_t IPAddress;
    uint16_t RemoteIP;
    uint16_t L3Proto;
    uint16_t RemoteL4Port;
};


int server_sock;
int client_sock;

void
term(int signum)
{
    printf("\nCaught SIGINT\n");
    close(server_sock);
    close(client_sock);
    exit(0);
}


int main(int argc, char *argv[])
{
    int port_number;
    int BUFFSIZE = 1024;
    int NUM_TESTS = 10;
    int loops = 0;
    pid_t pid;
    int opt_val = 1;
    int sockaddr_size;
    char client_message[BUFFSIZE];
    struct sockaddr_in server, client;
    struct sigaction action;

    if(argc < 2){   //check for correct number of args
        printf("Usage ./server port [number_of_tests default: 10]\n");
        return -1;
    }

    port_number = atoi(argv[1]);    //convert our argument to an integer port

    if(argc == 3){
        NUM_TESTS = atoi(argv[2]);
    }

    memset(&action, 0, sizeof(struct sigaction));
    action.sa_handler = term;
    sigaction(SIGINT, &action, NULL);
    
    //CreateSocket(PktIOSocket) (uint32, error)
    server_sock = socket(AF_INET, SOCK_STREAM, 0);   //create our tcp socket
    if(server_sock < 0){
        printf("Error creating server socket\n");
        return -1;
    }
    
    //SetSockoptTCPMD5(uint32, string, string)
    //enable reusable addresses option on our socket
    if (setsockopt(server_sock, SOL_SOCKET, SO_REUSEADDR, &opt_val, sizeof(int)) < 0){
        printf("setsockopt(SO_REUSEADDR) failed\n");
    }

    //populate the server sockaddr_in struct and store the port
    server.sin_family = AF_INET;
    server.sin_addr.s_addr = INADDR_ANY;
    server.sin_port = htons(port_number);

    //bind our program to the socket and port
    if (bind(server_sock, (struct sockaddr*)&server, sizeof(server)) < 0) {
		printf("bind failled\n");
        //DeleteSocket(uint32) error
        close(server_sock);
		return -1;
	}
    //Listen(uint32) error
    listen(server_sock, 10);

    //main server loop
    //Accept(uint32) (uint32, error)
    client_sock = accept(server_sock, (struct sockaddr *)&client, &sockaddr_size);
    while(loops < NUM_TESTS)
    {
        sockaddr_size = sizeof(client);
        bzero(client_message, BUFFSIZE);
    
        //Read([]byte) (int, error)
        if(recv(client_sock, client_message, BUFFSIZE, 0) < 0){
            printf("Recvfrom failed\n");
            //DeleteSocket(uint32) error
            close(server_sock);
            close(client_sock);
            return -1;
        }

        
        //Write([]byte) (int, error)
        if(send(client_sock, client_message, strlen(client_message), 0) < 0){
            printf("Sending message failed\n");
            //DeleteSocket(uint32) error
            close(server_sock);
            close(client_sock);
            return -1;
        }

        loops++;

        //printf("Message from client: %s\n", client_message);
    }
    //DeleteSocket(uint32) error
    close(client_sock);
    close(server_sock);
}
