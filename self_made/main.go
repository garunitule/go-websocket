package main

import (
	"fmt"
	wb "github.com/garunitule/go-websocket/self_made/websocket"
	"log"
	"net"
	"net/http"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	for {
		message, err := wb.ReadMessage(conn)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(message)

		err = wb.WriteMessage(conn, "received!")
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := wb.Upgrade(r, w)
	if err != nil {
		fmt.Println(err)
		return
	}
	go handleConnection(conn)
}

func main() {
	http.HandleFunc("/", wsHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
