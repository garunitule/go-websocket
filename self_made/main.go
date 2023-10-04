package main

import (
	"crypto/sha1"
	b64 "encoding/base64"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
)

func handleConnection(conn net.Conn) {
	fmt.Println("called handleConnection")
	defer conn.Close()

	for {
		header := make([]byte, 2)
		_, err := io.ReadFull(conn, header)
		if err != nil {
			fmt.Println(err)
			// TODO: エラーハンドリング
			return
		}
		fmt.Println(header)

		fin := header[0]&0x80 != 0
		opcode := header[0] & 0x0F
		masked := header[1]&0x80 != 0
		length := int(header[1] & 0x7F)
		fmt.Printf("fin: %t, opcode: %d, masked: %t, length: %d\n", fin, opcode, masked, length)

		if length == 126 || length == 127 {
			length, err = extraPayloadLength(conn, length)
			if err != nil {
				fmt.Println(err)
				// TODO: エラーハンドリング
				return
			}
		} else if length > 127 {
			fmt.Println("Invalid payload length. Please check the RFC.")
			return
		}

		maskKey := make([]byte, 4)
		_, err = io.ReadFull(conn, maskKey)
		if err != nil {
			fmt.Println(err)
			// TODO: エラーハンドリング
			return
		}

		payload := make([]byte, length)
		_, err = io.ReadFull(conn, payload)
		if err != nil {
			fmt.Println(err)
			// TODO: エラーハンドリング
			return
		}

		// unmask payload using maskKey
		// TODO: 調査
		for i := 0; i < length; i++ {
			payload[i] ^= maskKey[i%4]
		}

		strPayload := string(payload)
		fmt.Println(strPayload)

		return
	}
}

// 延長ペイロード長を取得する
func extraPayloadLength(conn net.Conn, payloadLength int) (int, error) {
	var n int
	if payloadLength == 126 {
		n = 2
	} else if payloadLength == 127 {
		n = 8
	} else {
		return 0, fmt.Errorf("Invalid payload length")
	}

	extraPayloadLengthBuffer := make([]byte, n)
	_, err := io.ReadFull(conn, extraPayloadLengthBuffer)
	if err != nil {
		fmt.Println(err)
		// TODO: エラーハンドリング
		return 0, err
	}

	length := 0
	for i := 0; i < n; i++ {
		t := 8 * (n - i - 1)
		fmt.Println(extraPayloadLengthBuffer[i], t, int(extraPayloadLengthBuffer[i])<<t)
		length |= int(extraPayloadLengthBuffer[i]) << t
	}
	return length, nil
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
