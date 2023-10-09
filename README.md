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

- client.jsを実行し、メッセージの送受信を確認

## TODO
- [] Sec-WebSocket-Extensionsヘッダの解析とhandlerへの反映