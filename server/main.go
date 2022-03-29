// Chat Server
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type client chan string

var chMessage = make(chan string)
var chEnteringClient = make(chan client)
var chLeavingClient = make(chan client)

func handleConnection(conn net.Conn) {
	fmt.Println(len(chEnteringClient))
	clientName := conn.RemoteAddr()
	chOutGoingMsg := make(chan string)
	chEnteringClient <- chOutGoingMsg
	fmt.Println("Client connected :", clientName)

	go clientWriter(conn, chOutGoingMsg)

	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		msg := scanner.Text()
		chMessage <- fmt.Sprintf("%s => %s\n", clientName, msg)
	}
	chLeavingClient <- chOutGoingMsg
	fmt.Printf("%s disconnected.\n", clientName)
}

func boardcaster() {
	clients := make(map[client]bool)

	for {
		select {
		case m := <-chMessage:
			for k := range clients {
				k <- m
			}
		case cli := <-chEnteringClient:
			clients[cli] = true
		case cli := <-chLeavingClient:
			clients[cli] = false
		}
	}
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		// io.Copy(conn, msg)
		fmt.Fprintf(conn, msg)
	}
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	go boardcaster()
	for {
		conn, err := listener.Accept() // waiting...
		if err != nil {
			// ignore
			log.Println(err)
		}
		go handleConnection(conn)
	}
}
