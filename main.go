package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	CtrlCHandler()

	var host string
	var port int
	var timeout int
	flag.IntVar(&timeout, "timeout", 100, "timeout in sec (int)")
	flag.StringVar(&host, "host", "localhost", "host")
	flag.IntVar(&port, "port", 80, "port")
	flag.Parse()

	log.Printf("Host: %v, Port: %v, Timeout: %v", host, port, timeout)

	dial(time.Duration(timeout)*time.Second, host, port)

}

func dial(timeout time.Duration, host string, port int) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%v:%v", host, port))
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Connected to %v\n", conn.RemoteAddr())

	if err = conn.SetDeadline(time.Now().Add(timeout)); err != nil {
		log.Fatal(err)
	}

	connectionClosed := false

	go func(conn net.Conn) {
		reader := bufio.NewReader(conn)
		for {
			str, err := reader.ReadString('\n')

			log.Print(str)

			if err == io.EOF {
				log.Println("EOF!")
				connectionClosed = true
				conn.Close()
				break
			}

			if err != nil {
				log.Fatal(err)
			}

		}
	}(conn)
	reader := bufio.NewReader(os.Stdin)
	for {
		if connectionClosed {
			return
		}
		str, err := reader.ReadString('\n')

		if connectionClosed {
			return
		}

		if err != nil {
			log.Fatal(err)
		}

		_, err = fmt.Fprint(conn, str)

		if err != nil {
			log.Fatal(err)
		}

	}
}

func CtrlCHandler() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\rCtrl+C pressed")
		os.Exit(0)
	}()
}
