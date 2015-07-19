package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"
)

const (
	PacketSize = 2
)

var (
	packet = []byte{109, 101}
)

type Server struct {
	Workers []*Worker
	Config  *Config
	Out     chan uint64
}

func NewServer(config *Config) (*Server, error) {
	if config == nil {
		return nil, IllegalArgumentError
	}

	s := &Server{Workers: nil, Config: config, Out: make(chan uint64, 2000)}
	workers := make([]*Worker, config.NumWorkers)
	for i := 0; i < config.NumWorkers; i++ {
		workers[i] = &Worker{WorkerId: uint32(i), DataCenterId: config.DataCenterId, Out: s.Out}
	}
	s.Workers = workers
	return s, nil
}

func (s *Server) HandleRequest(conn net.Conn) {
	fmt.Println("Handler Request:", conn.RemoteAddr())
	defer conn.Close()
	buf := make([]byte, PacketSize)

	for {
		_, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading:", err.Error())
			return
		}

		if !bytes.Equal(packet, buf) {
			fmt.Println("Error packet.")
			binary.Write(conn, binary.LittleEndian, 0)
			return
		}

		id := <-s.Out
		binary.Write(conn, binary.LittleEndian, id)
	}
}

func (s *Server) Serv() {
	serverSock, err := net.Listen("tcp", s.Config.Host+":"+s.Config.Port)
	//serverSock, err := net.Listen("tcp", "0.0.0.0:8000")

	if err != nil {
		fmt.Println("Error Listening:", err.Error())
		os.Exit(1)
	}
	fmt.Println("Start listening:", s.Config.Host, ":", s.Config.Port)

	defer serverSock.Close()

	// Start workers.
	for _, worker := range s.Workers {
		go func() {
			for {
				worker.YieldId()
			}
		}()
	}

	for {
		conn, err := serverSock.Accept()
		fmt.Println("accept a conn")
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}

		go s.HandleRequest(conn)
	}
}

func main() {
	config := &Config{Host: "0.0.0.0", Port: "8000", DataCenterId: 10, NumWorkers: 5}
	server, _ := NewServer(config)
	server.Serv()
}
