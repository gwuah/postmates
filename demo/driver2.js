const WebSocket = require("ws");

function connect(id) {
  console.log(`Electron ${id} initating a connection ... `);
  let ws = new WebSocket(`ws://localhost:8080/v1/electron/realtime/${id}`);

  ws.on("open", (e) => {
    console.log("connection successful");

    setInterval(() => {
      ws.send(
        JSON.stringify({
          meta: {
            type: "IndexElectronLocation",
          },
          id: id,
          latitude: 5.680814496321185,
          longitude: -0.23598240546815585,
        })
      );
    }, 2000);
  });

  ws.on("message", function (data) {
    console.log("Location Indexed", data);
  });

  ws.on("error", function (data) {
    console.log("Error connecting");
  });
}

function main() {
  connect(process.argv[2]);
}

main();
