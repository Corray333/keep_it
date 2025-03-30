<script setup lang="ts">
import { useAccountStore } from '@/stores/account';


const accStore = useAccountStore()

</script>

<template>
  <main>
    <h1>You are logged in</h1>
    <p>Token: {{ accStore.authorization }}</p>
  </main>
</template>
