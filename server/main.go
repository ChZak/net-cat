package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

type Client struct {
	name string
	ch   chan<- string // Canal pour envoyer les messages au client
}

var (
	entering       = make(chan Client)
	leaving        = make(chan Client)
	messages       = make(chan string) // Tous les messages des clients
	maxClients     = 10
	clientCount    = 0
	messageHistory []string
)

func main() {
	port := "8989"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}

	listener, err := net.Listen("tcp", "localhost:"+port)
	if err != nil {
		log.Fatal(err)
	}

	go broadcaster()
	log.Println("Server is running on port ", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}

		go handleConn(conn)
	}
}

func broadcaster() {
	clients := make(map[Client]bool) // Tous les clients connectés

	for {
		select {
		case msg := <-messages:

			for cli := range clients {
				cli.ch <- msg
			}

		case cli := <-entering:
			clients[cli] = true

		case cli := <-leaving:
			delete(clients, cli)
			close(cli.ch)
		}
	}
}

func handleConn(conn net.Conn) {
	if clientCount >= maxClients {
		fmt.Fprintln(conn, "Server is full, try again later")
		conn.Close()
		return
	}

	clientCount++

	ch := make(chan string) // Canal pour les messages du client
	go clientWriter(conn, ch)

	who := getClientName(conn)
	if who == "" {
		conn.Close()
		return
	}

	ch <- "Welcome to the chat, " + who
	for _, msg := range messageHistory {
		ch <- msg
	}

	messages <- fmt.Sprintf("[%s][%s] has arrived", time.Now().Format("2023-12-01 15:05:04"), who)
	entering <- Client{who, ch}

	input := bufio.NewScanner(conn)
	for input.Scan() {
		text := input.Text()
		if text != "" {
			msg := fmt.Sprintf("[%s][%s]: %s", time.Now().Format("2023-12-01 15:05:04"), who, text)
			messages <- msg
			messageHistory = append(messageHistory, msg)
		}
	}

	leaving <- Client{who, ch}
	messages <- fmt.Sprintf("[%s][%s] has left", time.Now().Format("2023-12-01 15:05:04"), who)
	conn.Close()
	clientCount--
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg) // Ignorer les erreurs potentielles
	}
}

func getClientName(conn net.Conn) string {
	fmt.Fprintln(conn, "Enter your name :")
	nameScanner := bufio.NewScanner(conn)

	if nameScanner.Scan() {
		return nameScanner.Text()
	}

	return ""
}
