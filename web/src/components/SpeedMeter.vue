<template>
  <div>
    <svg viewBox="0 0 1000 250">
      <g>
        <polygon :points="speedoPoints" fill="rgba(33, 150, 243,0.5)" />
      </g>
      <g>
        <text x="100" y="200" class="speed">
          {{ lastSpeedPoint }} {{ units }}
        </text>
      </g>
    </svg>
  </div>
</template>

<script>
export default {
  name: "SpeedMeter",

  props: {
    testData: undefined,
    units: undefined
  },
  data: () => ({}),

  computed: {
    speedoPoints: function() {
      const numSamples = this.testData.length;
      const maxElement = Math.max(...this.testData);
      let calculatedPoints = this.testData.map((x, i) => {
        return `${(1000 / numSamples) * i},${250 - (x / maxElement) * 250}`;
      });

      return `-0,250 ${calculatedPoints.join(" ")} 1000,250`;
    },
    lastSpeedPoint: function() {
      return this.testData[this.testData.length - 1];
    }
  }
};
</script>

<style scoped>
.speed {
  font-family: "Noto Sans JP", sans-serif;
  font-size: 80px;
}
</style>
