<script setup>
import TopBar from './components/TopBar.vue'
import NavBar from './components/NavBar.vue'
import {refreshTokens} from './utils/helpers'
import { onBeforeMount } from 'vue'
import { useStore } from 'vuex'

const store = useStore()

onBeforeMount(()=>{
  refreshTokens(store)
})

</script>

<template>
  <TopBar v-if="!$route.meta.hideNav" />
  <NavBar v-if="!$route.meta.hideNav" />
  <router-view class=" min-h-screen w-full"></router-view>
</template>

<style scoped>

</style>
