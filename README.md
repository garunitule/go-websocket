# go-websocket
GoでWebSocketを理解するためのリポジトリ

ディレクトリ構成
```txt
.
├── gorilla/          *gorilla/websocketを利用したWebSocketのサンプル実装
└── self_made/        *自作でWebSocketを実装
```

## 確認方法
- いずれかのサーバーを起動（ここではgorilla/websocketによるWebSocketサーバーを起動）
```
go run gorilla/main.go
```

- 下記のcurlコマンドでコネクション確立
```
curl -vvv -i -N -H "Connection: keep-alive, Upgrade" -H "Upgrade: websocket" -H "Sec-WebSocket-Version: 13" -H "Sec-WebSocket-Extensions: deflate-stream" -H "Sec-WebSocket-Key: WIY4slX50bnnSF1GaedKhg==" -H "Host: localhost:8080" -H "Origin:http://localhost:8080" http://localhost:8080
```

- TODO: WebSocketでメッセージやりとり

## TODO
- [] Sec-WebSocket-Extensionsヘッダの解析とhandlerへの反映