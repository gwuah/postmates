const WebSocket = require("ws");

const outsideScope = {
  latitude: 5.698188535023582,
  longitude: -0.239341780857103,
};

const defaultCabPositions = [
  {
    longitude: -0.2475990969444747,
    latitude: 5.684136332305188,
    color: "blue",
  },
  {
    longitude: -0.2397266058667604,
    latitude: 5.683835847589247,
    color: "blue",
  },
  {
    longitude: -0.24460022375167725,
    latitude: 5.677474538991623,
    color: "blue",
  },
];

// function electron(id) {
//   let ws = new WebSocket(`ws://localhost:8080/v1/electron/realtime/${id}`);
//   ws.on("message", function (data) {
//     console.log(`ID(${id}) >>> `, data);
//   });

//   ws.on("error", function (data) {
//     console.log("Error connecting");
//   });

//   return (coord) => {
//     setInterval(() => {
//       ws.send(
//         JSON.stringify({
//           meta: {
//             type: "IndexElectronLocation",
//           },
//           id: id,
//           latitude: coord.latitude,
//           longitude: coord.longitude,
//         })
//       );
//     }, 2000);
//   };
// }

function electron(id) {
  let ws = new WebSocket(`ws://localhost:8080/v1/electron/realtime/${id}`);
  ws.on("message", function (data) {
    parsed = JSON.parse(data);
    // console.log(JSON.stringify(parsed, null, 4));
    // console.log(`ID(${id}) >>> `, JSON.stringify(parsed, null, 4));
    console.log(`ID(${id}) >>> `, parsed.meta.type);

    if (parsed.meta.type == "NewDelivery" && id == "2") {
      console.log(`ID(${id}) >>> `, JSON.stringify(parsed, null, 4));

      setTimeout(() => {
        ws.send(
          JSON.stringify({
            meta: {
              type: "AcceptDelivery",
            },
            deliveryId: parsed.delivery.id,
          })
        );
      }, 1000);
    }
  });

  ws.on("error", function (data) {
    console.log("Error connecting", data);
  });

  return (coord) => {
    setInterval(() => {
      ws.send(
        JSON.stringify({
          meta: {
            type: "IndexElectronLocation",
          },
          id: id,
          latitude: coord.latitude,
          longitude: coord.longitude,
        })
      );
    }, 2000);
  };
}

function main() {
  setTimeout(() => {
    electron("1")(defaultCabPositions[0]);
  }, 1000);
  setTimeout(() => {
    electron("2")(defaultCabPositions[1]);
  }, 2000);
  setTimeout(() => {
    electron("3")(defaultCabPositions[2]);
  }, 3000);
}

main();
