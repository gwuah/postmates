const WebSocket = require("ws");

const origin = {
  latitude: 5.6796946725653745,
  longitude: -0.2447180449962616,
};

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
          origin,
          destination: {
            longitude: 2.4345545,
            latitude: 4.054594095,
          },
        })
      );
    }, 1000);
  });

  ws.on("message", function (data) {
    let parsed = JSON.parse(data);
    console.log(JSON.stringify(parsed, null, 4));
    // console.log(data);
  });

  ws.on("error", function (data) {
    console.log("Error connecting");
  });
}

function main() {
  connect(process.argv[2]);
}

main();
