<template>
  <v-content>
    <v-container>
      <v-row wrap>
        <v-col cols="6">
          <v-card>
            <v-card-title
              ><v-icon>mdi-download </v-icon>
              <span class="d-none d-sm-block"
                >Download Speed</span
              ></v-card-title
            >
            <v-card-text>
              <SpeedMeter :testData="downloadTestData" units="Mbps" />
            </v-card-text>
          </v-card>
        </v-col>
        <v-col cols="6">
          <v-card>
            <v-card-title
              ><v-icon>mdi-upload </v-icon>
              <span class="d-none d-sm-block">Upload Speed</span></v-card-title
            >
            <v-card-text>
              <SpeedMeter :testData="uploadTestData" units="Mbps" />
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>
      <v-row wrap>
        <v-col cols="6" class="d-none d-sm-block">
          <v-card height="100%">
            <v-card-title
              ><v-icon>mdi-information-outline </v-icon> Details</v-card-title
            >
            <v-card-text v-if="connectionInfo">
              <v-list-item two-line>
                <v-list-item-content>
                  <v-list-item-title>ISP</v-list-item-title>
                  <v-list-item-subtitle
                    >{{ connectionInfo.autonomousSystemOrganisation }} (AS
                    {{
                      connectionInfo.autonomousSystemNumber
                    }})</v-list-item-subtitle
                  >
                </v-list-item-content>
              </v-list-item>
              <v-list-item two-line>
                <v-list-item-content>
                  <v-list-item-title>IP Address</v-list-item-title>
                  <v-list-item-subtitle>{{
                    connectionInfo.remoteAddress
                  }}</v-list-item-subtitle>
                </v-list-item-content>
              </v-list-item>
            </v-card-text>

            <v-card-text v-else>
              Loading Connection Details
            </v-card-text>
          </v-card>
        </v-col>
        <v-col cols="12" sm="6">
          <v-card height="100%">
            <v-card-title
              ><v-icon>mdi-format-text-wrapping-wrap</v-icon
              >Latency</v-card-title
            >
            <v-card-text>
              <SpeedMeter :testData="pingData" units="ms" />
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>

      <v-row>
        <v-col cols="12">
          <v-select
            :disabled="testing"
            :items="serverList"
            label="Test Server"
            v-model="speedServer"
            outlined
          ></v-select>
        </v-col>
      </v-row>
      <v-row>
        <v-col cols="12">
          <v-btn
            color="primary"
            depressed
            :disabled="testing"
            x-large
            fill-width
            @click="executeSpeedTest"
            width="100%"
            height="5vh"
            >{{ testingText }}</v-btn
          >
        </v-col>
      </v-row>
    </v-container>
  </v-content>
</template>

<script>
import SpeedMeter from "@/components/SpeedMeter.vue";
import InfoController from "@/controllers/info.js";
import SpeedController from "@/controllers/speed.js";
import PingController from "@/controllers/ping.js";

export default {
  name: "Home",
  components: {
    SpeedMeter
  },

  computed: {
    testingText: function() {
      return this.testing ? "Please Wait" : "Begin Test";
    }
  },

  data: () => ({
    downloadTestData: [0],
    uploadTestData: [0],
    pingData: [0],
    testing: true,
    connectionInfo: false,
    speedServer: {
      name: "Standalone Server",
      httpEndpoint: `${window.location.protocol}//${window.location.host}`,
      websocketEndpoint: `${window.location.protocol == "http:" ? "ws:" : "wss:"}//${window.location.host}`
    },

    serverList: [
      {
        text: "Standalone Server",
        value: {
          name: "Standalone Server",
          httpEndpoint: `${window.location.protocol}//${window.location.host}`,
          websocketEndpoint: `${window.location.protocol == "http:" ? "ws:" : "wss:"}//${window.location.host}`
        }
      }
    ]
  }),

  methods: {
    executeSpeedTest: function() {
      let that = this;
      that.testing = true;
      this.uploadTestData = [0];
      this.downloadTestData = [0];
      this.pingData = [0];
      this.PingController = new PingController(
        this.speedServer.websocketEndpoint,
        duration => {
          that.pingData.push(duration);
        }
      );
      this.SpeedController = new SpeedController(this.speedServer.httpEndpoint);
      this.SpeedController.executeDownload(
        64777216,
        event => {
          this.downloadTestData.push((event.speed / 1000000).toFixed(2));
        },
        () => {
          this.executeUpload();
        }
      );
      this.pingInterval = setInterval(() => {
        that.PingController.ping();
      }, 200);
    },

    executeUpload: function() {
      let that = this;
      this.SpeedController.executeUpload(
        64777216,
        event => {
          this.uploadTestData.push((event.speed / 1000000).toFixed(2));
        },
        () => {
          clearInterval(this.pingInterval);
          that.testing = false;
        }
      );
    }
  },
  async mounted() {
    let connectionQuery = await InfoController.queryInfo(
      this.speedServer.httpEndpoint
    );
    this.connectionInfo = connectionQuery;
    this.testing = false;
  }
};
</script>
