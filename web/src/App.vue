<template>
  <v-app>
    <v-app-bar app color="primary" dark>
      <div>
        <v-img class="mr-3" src="@/assets/openspeed-logo.svg" width="30%">
        </v-img>
      </div>
      <!-- <v-toolbar-title>OpenSpeed</v-toolbar-title> -->
      <v-spacer></v-spacer>
    </v-app-bar>
    <router-view />
  </v-app>
</template>
<script>
export default {
  name: "App",
  components: {},

  data: () => ({})
};
</script>
