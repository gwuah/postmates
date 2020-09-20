const WebSocket = require("ws");

const origin = {
  lat: 5.6796946725653745,
  lng: -0.2447180449962616,
};

function connect(id) {
  console.log(`Customer ${id} initating a connection ... `);
  let ws = new WebSocket(`ws://localhost:8080/v1/customer/realtime/${id}`);

  ws.on("open", (e) => {
    console.log("connection successful");

    // setTimeout(() => {
    //   console.log("sending message");
    //   ws.send(
    //     JSON.stringify({
    //       meta: {
    //         type: "DeliveryRequest",
    //       },
    //       productId: 3,
    //       customerID: 1,
    //       notes: "Hello",
    //       origin: {
    //         longitude: 2.4345545,
    //         latitude: 4.054594095,
    //       },
    //       destination: {
    //         longitude: 2.4345545,
    //         latitude: 4.054594095,
    //       },
    //     })
    //   );
    // }, 1000);

    // setTimeout(() => {
    //   console.log("sending message");
    //   ws.send(
    //     JSON.stringify({
    //       meta: {
    //         type: "CancelDelivery",
    //       },
    //       tripId: 1,
    //     })
    //   );
    // }, 1000);

    setInterval(() => {
      ws.send(
        JSON.stringify({
          meta: {
            type: "GetClosestElectrons",
          },
          id: id,
          origin,
        })
      );
    }, 2000);
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
