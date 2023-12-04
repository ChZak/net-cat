package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	file, err := os.ReadFile("welcome.txt")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(file))
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <host> <port>\n", os.Args[0])
		os.Exit(1)
	}

	host := os.Args[1]
	port := os.Args[2]
	address := host + ":" + port

	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}

	defer conn.Close()

	go handleIncomingMessages(conn)
	handleOutgoingMessages(conn)
}

func handleIncomingMessages(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if scanner.Err() != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", scanner.Err().Error())
	}
}

func handleOutgoingMessages(conn net.Conn) {
	consoleReader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		msg, _ := consoleReader.ReadString('\n')
		msg = strings.TrimSpace(msg)

		if msg == "" {
			fmt.Println("Vous ne pouvez pas envoyer un message vide")
			continue
		}

		if msg == "exit" {
			return
		}

		_, err := conn.Write([]byte(msg + "\n"))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
			return
		}
	}
}
