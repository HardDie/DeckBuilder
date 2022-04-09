<template>
  <v-app>
    <v-main>
      <div class="main-wrapper">
        <router-view/>
      </div>
    </v-main>
  </v-app>
</template>

<script>

export default {
  name: 'App'
}
</script>

<style lang="scss">
.main-wrapper {
  width: 100%;
  height: 100%;
  padding: 20px;
}
</style>
