package main

import (
	"bufio"
	"log"
	"net"
	"os"
)

const (
	connHost = "localhost"
	connPort = "8080"
	connType = "tcp"
)

// Client of the server
type Client struct {
	conn net.Conn
}

func main() {
	log.Printf("Connecting to %v server %v:%v ...\n", connType, connHost, connPort)

	conn, err := net.Dial(connType, connHost+":"+connPort)
	if err != nil {
		log.Fatal("Error connecting to server: ", err.Error())
	}
	defer conn.Close()

	client := Client{
		conn: conn,
	}

	go client.handleServerMessage()
	client.handleInput()
}

func (c *Client) handleInput() {
	reader := bufio.NewReader(os.Stdin)

	for {
		input, _ := reader.ReadString('\n')

		_, err := c.conn.Write([]byte(input))
		if err != nil {
			log.Printf("Error sending message to server: ", err.Error())
		}
	}
}

func (c *Client) handleServerMessage() {
	for {
		message, _ := bufio.NewReader(c.conn).ReadString('\n')

		log.Print("Server sent: ", message)
	}
}
