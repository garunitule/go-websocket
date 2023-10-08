package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
)

func handleConnection(conn net.Conn) {
	fmt.Println("called handleConnection")
	defer conn.Close()

	for {
		strPayload, err := ReadMessage(conn)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(strPayload)

		return
	}
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := Upgrade(r, w)
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
