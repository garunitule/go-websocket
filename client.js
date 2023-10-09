let ws = new WebSocket("ws://localhost:8080");
ws.onopen = function() {
  console.log("Connected!");
  ws.send("Hello, Server!");
};
ws.onmessage = function(event) {
  console.log("Received:", event.data);
};