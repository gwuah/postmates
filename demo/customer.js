const WebSocket = require("ws");

function connect(id) {
  console.log(`Customer ${id} initating a connection ... `);
  let ws = new WebSocket(`ws://localhost:8080/v1/customer/realtime/${id}`);

  ws.on("open", (e) => {
    console.log("connection successful");
  });

  ws.on("message", function (data) {
    console.log(data);
  });

  ws.on("error", function (data) {
    console.log("Error connecting");
  });
}

function parseMessage(message) {}

function main() {
  connect(process.argv[2]);
}

main();
