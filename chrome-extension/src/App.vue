<script setup lang="ts">
import { onMounted } from 'vue';
import { RouterLink, RouterView, useRoute } from 'vue-router'
import { UserTransport } from './services/userService';

const route = useRoute()
onMounted(()=>
  UserTransport.renewTokens()
)
</script>

<template>
  <RouterView />
</template>

<style scoped>

</style>
