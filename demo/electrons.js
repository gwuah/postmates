const WebSocket = require("ws");


const defaultCabPositions = [
  { latitude: 5.68435947053963, longitude: -0.25092702358961105, color: "blue", name: "GGG"},
  { latitude: 5.682181546392756, longitude: -0.25024037808179855, color: "blue", name: "HHH" },
  { latitude: 5.682565886546736, longitude: -0.24792294949293137, color: "blue", name: "III" },
  { latitude: 5.685000034897074, longitude: -0.244661383330822, color: "blue", name: "BBB" },
  { latitude: 5.685597894320817, longitude: -0.24564843624830246, color: "blue", name: "CCC" },
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

    this.ws.on("message", (msg) => {
      console.log("new message");
      this.handleMessage(msg);
    });

    this.ws.on("error", (data) => {
      console.log("Error connecting", data);
    });

    console.log(`Courier ${this.appState.id} has been instantiated & is sending location updates every 3 seconds`);

    this._sendLocationUpdate();
  }

  _handleNewDelivery(parsed) {
    console.log(`From Courier (${this.appState.id}) -> NewDeliveryRequest Recieved!`);
    // if (this.appState.id == "III") {
    //   console.log(
    //     `ID(${this.appState.id}) >>> `,
    //     JSON.stringify(parsed, null, 4)
    //   );
    //   this.appState.current_delivery = parsed.delivery;
    //   this.ws.send(
    //     JSON.stringify({
    //       meta: {
    //         type: "AcceptDelivery",
    //       },
    //       deliveryId: parsed.delivery.id,
    //     })
    //   );
    // }
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
          id: this.appState.id,
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
  var c1 = new Courier("GGG", defaultCabPositions[0]);
  c1._initialization();
  var c2 = new Courier("HHH", defaultCabPositions[1]);
  c2._initialization();
  var c3 = new Courier("III", defaultCabPositions[2]);
  c3._initialization();
  var c4 = new Courier("BBB", defaultCabPositions[3]);
  c4._initialization();
  var c5 = new Courier("CCC", defaultCabPositions[4]);
  c5._initialization();
}

main();
