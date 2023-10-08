package main

import (
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
