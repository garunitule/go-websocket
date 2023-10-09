package websocket

import (
	"fmt"
	"io"
	"net"
)

func ReadMessage(conn net.Conn) (string, error) {
	header := make([]byte, 2)
	_, err := io.ReadFull(conn, header)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	// TODO: fin, opcodeに応じた処理を実装する
	// fin := header[0]&0x80 != 0
	// opcode := header[0] & 0x0F
	masked := header[1]&0x80 != 0
	length := int(header[1] & 0x7F)

	if length == 126 || length == 127 {
		length, err = extraPayloadLength(conn, length)
		if err != nil {
			fmt.Println(err)
			return "", err
		}
	} else if length > 127 {
		err = fmt.Errorf("Invalid payload length. Please check the RFC.")
		fmt.Println(err)
		return "", err
	}

	maskKey := make([]byte, 4)
	if masked {
		_, err = io.ReadFull(conn, maskKey)
		if err != nil {
			fmt.Println(err)
			return "", err
		}
	}

	payload := make([]byte, length)
	_, err = io.ReadFull(conn, payload)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	if masked {
		for i := 0; i < length; i++ {
			payload[i] ^= maskKey[i%4]
		}
	}

	strPayload := string(payload)
	return strPayload, nil
}

func WriteMessage(conn net.Conn, message string) error {
	var response []byte
	response = append(response, 0x81)
	length := len(message)
	if length < 126 {
		response = append(response, byte(length))
	} else if length < 65536 {
		response = append(response, 126)
	} else {
		response = append(response, 127)
	}

	if 126 <= length {
		var n int
		if length < 65536 {
			n = 2
		} else {
			n = 8
		}
		extendedPayloadLength := make([]byte, n)
		for i := range extendedPayloadLength {
			extendedPayloadLength[i] = byte(length >> 8 * (n - i - 1))
		}
		response = append(response, extendedPayloadLength...)
	}
	response = append(response, []byte(message)...)
	_, err := conn.Write(response)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
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
		return 0, err
	}

	length := 0
	for i := 0; i < n; i++ {
		t := 8 * (n - i - 1)
		length |= int(extraPayloadLengthBuffer[i]) << t
	}
	return length, nil
}
