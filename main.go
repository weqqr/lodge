package main

import (
	"context"
	"log"
	"net"
)

type Server struct {
}

func (s *Server) Serve(ctx context.Context, address string) error {
	const maxPacketSize = 1536

	config := &net.ListenConfig{}

	conn, err := config.ListenPacket(ctx, "udp", address)
	if err != nil {
		return err
	}

	defer conn.Close()

	quit := make(chan struct{})

	go func() {
		for {
			buf := make([]byte, maxPacketSize)

			n, addr, err := conn.ReadFrom(buf)
			if err != nil {
				log.Printf("network read error: %v", err)

				quit <- struct{}{}
			}

			packetData := buf[:n]

			s.handlePacket(packetData, addr)
		}
	}()

	select {
	case <-ctx.Done():
	case <-quit:
	}

	return nil
}

func (s *Server) handlePacket(data []byte, addr net.Addr) {
	log.Printf("data from %v: %v", addr, data)
}

func main() {
	ctx := context.Background()

	address := "0.0.0.0:31000"
	server := Server{}

	log.Printf("address: %v", address)

	server.Serve(ctx, address)
}
