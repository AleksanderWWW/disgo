package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/AleksanderWWW/disgo/pdu"
)

func receiver() {
	// Set up UDP address
	addr, err := net.ResolveUDPAddr("udp", "localhost:8080")
	if err != nil {
		log.Fatal("Could not resolve address:", err)
	}

	// Start listening
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatal("Listen failed:", err)
	}
	defer conn.Close()

	// Buffer for incoming data
	buffer := make([]byte, 1024)
	n, _, err := conn.ReadFromUDP(buffer)
	if err != nil {
		log.Fatalf("Read error: %v", err)
	}
	stream := bytes.NewReader(buffer[:n])

	fire, err := pdu.ParseFirePDU(stream)
	if err != nil {
		log.Fatalf("Could not parse PDU: %s", err)
	}

	fmt.Println("Received PDU")
	fmt.Println(fire)
}

func sender() {
	msg := pdu.FirePDU{
		Header: pdu.PDUHeader{
			PDUType:    2,
			ExerciseID: 2,
		},
		FiringID: pdu.EntityId{
			Address: pdu.SimulationAddress{
				Site:        1,
				Application: 2,
			},
			Number: 5,
		},
		TargetID: pdu.EntityId{
			Address: pdu.SimulationAddress{
				Site:        1,
				Application: 2,
			},
			Number: 6,
		},
	}
	var buf bytes.Buffer

	err := msg.Serialize(&buf)
	if err != nil {
		log.Fatal("Could not serialize PDU:", err)
	}

	// Connect to server
	addr, err := net.ResolveUDPAddr("udp", "localhost:8080")
	if err != nil {
		log.Fatal("Could not resolve address:", err)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatal("Connection failed:", err)
	}
	defer conn.Close()

	// Send a message
	message := buf.Bytes()
	_, err = conn.Write(message)
	if err != nil {
		log.Printf("Send failed: %v", err)
		return
	}
}

func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		receiver()
		wg.Done()
	}()

	go func() {
		time.Sleep(time.Second)
		sender()
		wg.Done()
	}()

	wg.Wait()
}
