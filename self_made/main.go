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
		strPayload, err := readRequestPayload(conn)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(strPayload)

		return
	}
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%+v\n", r)
	err := validateHandShakeReq(r, w)
	if err != nil {
		fmt.Println(err)
		return
	}

	hijacker, ok := w.(http.Hijacker)
	if !ok {
		fmt.Println("Not a hijacker")
		http.Error(w, "WebSocket upgrade failed", http.StatusInternalServerError)
		return
	}

	conn, buf, err := hijacker.Hijack()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "WebSocket upgrade failed", http.StatusInternalServerError)
		return
	}

	buf = openHandShake(r, buf)
	buf.Flush()
	go handleConnection(conn)
}

func main() {
	http.HandleFunc("/", wsHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
