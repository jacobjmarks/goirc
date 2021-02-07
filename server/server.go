package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
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

// RemoteAddr returns the remote address of the client
func (c *Client) RemoteAddr() string {
	return c.conn.RemoteAddr().String()
}

// Server with clients
type Server struct {
	network string
	host    string
	port    string
	clients map[*Client]*Client
}

func newServer(network string, host string, port string) *Server {
	log.Printf("Creating new %v server on %v:%v ...\n", connType, connHost, connPort)

	server := Server{
		network: network,
		host:    host,
		port:    port,
		clients: make(map[*Client]*Client),
	}

	return &server
}

// Listen opens the server and awaits incoming client connections
func (s *Server) Listen() {
	listener, err := net.Listen(s.network, s.host+":"+s.port)
	if err != nil {
		log.Fatal("Error listening:", err.Error())
	}
	defer listener.Close()

	log.Printf("Listening for incoming connections ...\n")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("Error accepting connection:", err.Error())
		}

		s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	// TODO: Validate incoming connection

	client := Client{
		conn: conn,
	}

	log.Printf("Client %v connected\n", client.RemoteAddr())

	s.enqueueClient(&client)
	go s.handleClient(&client)
}

func (s *Server) enqueueClient(c *Client) {
	s.clients[c] = c
	go s.broadcast([]byte(fmt.Sprintf("A user has connected to the server. There are now %v others here.\n", len(s.clients)-1)), c)
}

func (s *Server) handleClient(c *Client) {
	c.conn.Write([]byte(fmt.Sprintf("Welcome. There are %v others currently here.\n", len(s.clients)-1)))

	for {
		buffer, err := bufio.NewReader(c.conn).ReadBytes('\n')

		if err != nil {
			log.Printf("Client %v disconnected\n", c.RemoteAddr())
			c.conn.Close()
			delete(s.clients, c)
			go s.broadcast([]byte(fmt.Sprintf("A user has disconnected from the server. There are now %v others here.\n", len(s.clients)-1)), nil)
			return
		}

		go s.receiveMessage(buffer, c)
	}
}

func (s *Server) receiveMessage(msg []byte, c *Client) {
	// TODO: Validate payload
	log.Printf("Received message from client %v: %v\n", c.RemoteAddr(), string(msg[:len(msg)-1]))
	go s.broadcast(msg, c)
}

func (s *Server) broadcast(msg []byte, client *Client) {
	log.Printf("Broadcasting message to %v clients: %v", len(s.clients), string(msg[:len(msg)-1]))
	for _, c := range s.clients {
		if c == client {
			continue
		}

		go s.sendToClient(msg, c)
	}
}

func (s *Server) sendToClient(msg []byte, c *Client) {
	_, err := c.conn.Write(msg)
	if err != nil {
		log.Printf("Error sending message to client %v: %v", c.RemoteAddr(), err.Error())
	}
}

func main() {
	s := newServer(connType, connHost, connPort)
	s.Listen()
}
