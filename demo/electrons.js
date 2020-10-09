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

class Courier {
  constructor(id, coord) {
    this.appState = {
      id,
      coord,
      state: "awaiting_dispatch",
      courier: null,
      current_delivery: null,
    };
  }

  _initialization() {
    this.ws = new WebSocket(
      `ws://localhost:8080/v1/courier/realtime/${this.appState.id}`
    );

    this.ws.on("message",  (msg) => {
      this.handleMessage(msg)
    })

    this.ws.on("error", (data) =>{
      console.log("Error connecting", data);
    });
  }

  _handleNewDelivery(parsed) {
    console.log(`NewDeliveryRequest Recieved ${this.appState.id} `);
    if (this.appState.id == "2") {
      console.log(`ID(${this.appState.id}) >>> `, JSON.stringify(parsed, null, 4));
      this.appState.current_delivery = parsed.delivery;
      this.ws.send(
        JSON.stringify({
          meta: {
            type: "AcceptDelivery",
          },
          deliveryId: parsed.delivery.id,
        })
      );
    }
  }

  _sendLocationUpdate() {
    const deliveryId = this.appState.current_delivery
      ? this.appState.current_delivery.id
      : null;

    setInterval(() => {
      this.ws.send(
        JSON.stringify({
          meta: {
            type: "LocationUpdate",
          },
          id: id,
          latitude: this.appState.coord.latitude,
          longitude: this.appState.coord.longitude,
          state: this.appState.state,
          deliveryId,
        })
      );
    }, 3000);
  }

  handleMessage(message) {
    let parsed = JSON.parse(message);
    switch (parsed.meta.type) {
      case "NewDelivery":
        this._handleNewDelivery(parsed);
        break;
    }
  }
}


function main() {
  var c1 = new Courier("1", defaultCabPositions[0])
  c1._initialization()
  var c2 = new Courier("2", defaultCabPositions[1])
  c2._initialization()
  var c3 = new Courier("3", defaultCabPositions[2])
  c3._initialization()
}

main();
