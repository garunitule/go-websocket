package main

import (
	"crypto/sha1"
	b64 "encoding/base64"
	"fmt"
	"log"
	"net"
	"net/http"
)

// TODO: implement the rest of the communication
func handleConnection(conn net.Conn) {

}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%+v\n", r)
	if r.Method != http.MethodGet {
		fmt.Println("Not a GET request")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("Upgrade") != "websocket" {
		fmt.Println("Not a websocket request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if r.Header.Get("Host") != "" && r.Header.Get("Connection") != "" && r.Header.Get("Sec-WebSocket-Key") != "" && r.Header.Get("Sec-WebSocket-Version") != "" {
		fmt.Println("Missing required headers")
		w.WriteHeader(http.StatusBadRequest)
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
	defer conn.Close()

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
	if r.Header.Get("Sec-WebSocket-Extensions") != "" {
		fmt.Fprintf(buf, "Sec-WebSocket-Extensions: %s\r\n", r.Header.Get("Sec-WebSocket-Extensions"))
	}
	buf.Flush()

	go handleConnection(conn)
}

func main() {
	http.HandleFunc("/", wsHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
