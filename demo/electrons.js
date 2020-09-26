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
          latitude: coord.latitude,
          longitude: coord.longitude,
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
          latitude: coord.latitude,
          longitude: coord.longitude,
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
