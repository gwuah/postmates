const WebSocket = require("ws");

function connect(id) {
  console.log(`Customer ${id} initating a connection ... `);
  let ws = new WebSocket(`ws://localhost:8080/v1/customer/realtime/${id}`);

  ws.on("open", (e) => {
    console.log("connection successful");

    setTimeout(() => {
      ws.send(
        JSON.stringify({
          meta: {
            type: "DeliveryRequest",
          },
          productId: 1,
          customerID: 1,
          notes: "Handle it carefully",
          origin: {
            latitude: 5.683704379865978, 
            longitude: -0.24667987691595838,
          },
          destination: {
            latitude: 5.5952893396562615, 
            longitude: -0.23337027814121614
          },
        })
      );
    }, 1000);
  });

  ws.on("message", function (data) {
    let parsed = JSON.parse(data);
    console.log(JSON.stringify(parsed, null, 4));
  });

  ws.on("error", function (data) {
    console.log("Error connecting");
  });
}

function main() {
  connect("PostMaster");
}

main();
