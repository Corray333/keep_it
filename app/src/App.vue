<script setup>
import TopBar from './components/TopBar.vue'

</script>

<template>
  <TopBar v-if="$route.path!=='/login'" />
  <router-view class=" min-h-screen w-full p-5"></router-view>
</template>

<style scoped>

</style>
