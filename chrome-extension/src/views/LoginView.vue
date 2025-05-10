<script lang="ts" setup>
import { useAccountStore } from '@/stores/account';


const accStore = useAccountStore()

</script>

<template>
    <main>
        <h1>Login first</h1>
        <p>Token: {{ accStore.authorization }}</p>
    </main>
</template>

<style scoped>

</style>
