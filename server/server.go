package main

import (
	"bufio"
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

// Server with clients
type Server struct {
	network string
	host    string
	port    string
	clients []*Client
}

func newServer(network string, host string, port string) *Server {
	log.Printf("Creating new %v server on %v:%v ...\n", connType, connHost, connPort)

	server := Server{
		network: network,
		host:    host,
		port:    port,
		clients: []*Client{},
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

	log.Printf("Client %v connected\n", client.conn.RemoteAddr().String())

	s.enqueueClient(&client)
	go s.handleClient(&client)
}

func (s *Server) enqueueClient(client *Client) {
	// Find available client slot
	for i, c := range s.clients {
		if c == nil || c.conn == nil {
			s.clients[i] = client
			return
		}
	}

	// No existing clients or all are active
	s.clients = append(s.clients, client)
}

func (s *Server) handleClient(c *Client) {
	for {
		buffer, err := bufio.NewReader(c.conn).ReadBytes('\n')

		if err != nil {
			log.Printf("Client %v disconnected\n", c.conn.RemoteAddr().String())
			c.conn.Close()
			c.conn = nil
			return
		}

		s.receiveMessage(buffer, c)
	}
}

func (s *Server) receiveMessage(msg []byte, c *Client) {
	// TODO: Validate payload
	log.Printf("Received message from client %v: %v\n", c.conn.RemoteAddr().String(), string(msg[:len(msg)-1]))
	s.broadcast(msg, c)
}

func (s *Server) broadcast(msg []byte, client *Client) {
	log.Printf("Broadcasting message to %v clients: %v", len(s.clients), string(msg[:len(msg)-1]))
	for _, c := range s.clients {
		if c == client || c.conn == nil {
			continue
		}

		go s.sendToClient(msg, c)
	}
}

func (s *Server) sendToClient(msg []byte, c *Client) {
	_, err := c.conn.Write(msg)
	if err != nil {
		log.Printf("Error sending message to client %v: %v", c.conn.RemoteAddr().String(), err.Error())
	}
}

func main() {
	s := newServer(connType, connHost, connPort)
	s.Listen()
}
