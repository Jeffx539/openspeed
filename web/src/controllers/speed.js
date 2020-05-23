class SpeedController {
  constructor(endpoint) {
    this.endpoint = endpoint;
  }

  executeDownload(chunkSize, progressCallback, completeCallback) {
    this.xhr = new XMLHttpRequest();
    this.chunkSize = chunkSize;
    let that = this;

    this.xhr.onload = function() {
      if (that.xhr.status == 200) {
        completeCallback({
          chunkSize: that.chunkSize
        });
      } else {
        console.log("Request error: " + that.xhr.statusText);
      }
    };

    this.xhr.open(
      "GET",
      `${this.endpoint}/chunk?size=${chunkSize}&cachebust=${performance.now()}`,
      true
    );
    this.xhr.onprogress = event => {
      let elapsed = (performance.now() - that.startTime) / 1000;
      progressCallback({
        speed: (event.loaded / elapsed) * 8,
        progress: event.loaded / event.total
      });
    };
    this.xhr.responseType = "arraybuffer";
    this.startTime = performance.now();
    this.xhr.send();
  }

  executeUpload(chunkSize, progressCallback, completeCallback) {
    this.xhr = new XMLHttpRequest();
    var blob = new ArrayBuffer(chunkSize);
    var uInt8Array = new Uint8Array(blob);
    this.chunkSize = chunkSize;
    let that = this;

    this.xhr.onload = function() {
      if (that.xhr.status == 200) {
        completeCallback({
          chunkSize: that.chunkSize
        });
      } else {
        console.log("Request error: " + that.xhr.statusText);
      }
    };

    this.xhr.open(
      "POST",
      `${
        this.endpoint
      }/upload?size=${chunkSize}&cachebust=${performance.now()}`,
      true
    );
    this.xhr.upload.onprogress = event => {
      let elapsed = (performance.now() - that.startTime) / 1000;
      progressCallback({
        speed: (event.loaded / elapsed) * 8,
        progress: event.loaded / event.total
      });
    };
    this.startTime = performance.now();
    this.xhr.send(uInt8Array);
  }
}

export default SpeedController;
