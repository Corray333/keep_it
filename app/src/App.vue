<script setup>
import TopBar from './components/TopBar.vue'
import NavBar from './components/NavBar.vue'


</script>

<template>
  <section class="flex flex-col gap-5">
    <TopBar v-if="!$route.meta.hideNav" />
    <main class="flex">
      <NavBar v-if="!$route.meta.hideNav" />
      <router-view class=" min-h-screen w-full"></router-view>  
    </main>
  </section>
</template>

<style scoped></style>
