const WebSocket = require("ws");

const outsideScope = { lat: 5.698188535023582, lng: -0.239341780857103 };

const defaultCabPositions = [
  { lng: -0.2475990969444747, lat: 5.684136332305188, color: "blue" },
  { lng: -0.2397266058667604, lat: 5.683835847589247, color: "blue" },
  { lng: -0.24460022375167725, lat: 5.677474538991623, color: "blue" },
];

function electron(id) {
  let ws = new WebSocket(`ws://localhost:8080/v1/electron/realtime/${id}`);
  ws.on("message", function (data) {
    console.log(`ID(${id}) >>> `, data);
  });

  ws.on("error", function (data) {
    console.log("Error connecting");
  });

  return (coord) => {
    setInterval(() => {
      ws.send(
        JSON.stringify({
          meta: {
            type: "IndexElectronLocation",
          },
          id: id,
          lat: coord.lat,
          lng: coord.lng,
        })
      );
    }, 2000);
  };
}

function electron(id) {
  let ws = new WebSocket(`ws://localhost:8080/v1/electron/realtime/${id}`);
  ws.on("message", function (data) {
    console.log(`ID(${id}) >>> `, data);
  });

  ws.on("error", function (data) {
    console.log("Error connecting");
  });

  return (coord) => {
    setInterval(() => {
      ws.send(
        JSON.stringify({
          meta: {
            type: "IndexElectronLocation",
          },
          id: id,
          lat: coord.lat,
          lng: coord.lng,
        })
      );
    }, 2000);
  };
}

function main() {
  electron("1")(defaultCabPositions[0]);
  electron("2")(defaultCabPositions[1]);
  electron("3")(defaultCabPositions[2]);
}

main();
