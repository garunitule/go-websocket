package main

import (
	"crypto/sha1"
	b64 "encoding/base64"
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

	// implement the rest of the handshake
	key := r.Header.Get("Sec-WebSocket-Key")
	h := sha1.New()
	// TODO: 定数で良い理由を調査
	h.Write([]byte(key + "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"))
	accept := b64.StdEncoding.EncodeToString(h.Sum(nil))
	fmt.Fprintf(buf, "HTTP/1.1 101 Switching Protocols\r\n")
	fmt.Fprintf(buf, "Upgrade: websocket\r\n")
	fmt.Fprintf(buf, "Connection: Upgrade\r\n")
	fmt.Fprintf(buf, "Sec-WebSocket-Accept: %s\r\n", accept)
	if r.Header.Get("Sec-WebSocket-Protocol") != "" {
		fmt.Fprintf(buf, "Sec-WebSocket-Protocol: %s\r\n", r.Header.Get("Sec-WebSocket-Protocol"))
	}
	fmt.Fprintf(buf, "\r\n")

	buf.Flush()
	go handleConnection(conn)
}

func main() {
	http.HandleFunc("/", wsHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
