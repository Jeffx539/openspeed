class PingController {
  constructor(endpoint, pongCB) {
    this.socket = new WebSocket(`${endpoint}/ping`);
    this.ready = false;
    this.pongCB = pongCB;
    let that = this;

    this.socket.addEventListener("open", function() {
      that.ready = true;
      console.log("Ping WS Opened and Ready");
    });

    this.socket.addEventListener("message", function() {
      console.log(new Date() - this.startInterval);
      let duration = new Date() - that.startInterval;
      that.pongCB(duration);
    });
  }

  ping() {
    if (!this.ready) {
      return;
    }
    this.startInterval = new Date();
    this.socket.send("ping");
  }
}

export default PingController;
