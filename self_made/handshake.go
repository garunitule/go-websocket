package main

import (
	"bufio"
	"crypto/sha1"
	b64 "encoding/base64"
	"fmt"
	"net/http"
)

func validateHandShakeReq(r *http.Request, w http.ResponseWriter) error {
	if r.Method != http.MethodGet {
		err := fmt.Errorf("Not a GET request")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return err
	}

	if r.Header.Get("Upgrade") != "websocket" {
		err := fmt.Errorf("Not a websocket request")
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	if r.Header.Get("Host") != "" && r.Header.Get("Connection") != "" && r.Header.Get("Sec-WebSocket-Key") != "" && r.Header.Get("Sec-WebSocket-Version") != "" {
		err := fmt.Errorf("Missing required headers")
		w.WriteHeader(http.StatusBadRequest)
		return err
	}
	return nil
}

func openHandShake(r *http.Request, buf *bufio.ReadWriter) *bufio.ReadWriter {
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
	return buf
}
