<script setup>
import { Icon } from '@iconify/vue'

</script>

<template>
    <nav class="flex w-full absolute bottom-0 p-3 text-2xl justify-around bottom-bar-shadow rounded-t-xl">
        <router-link to="/home"><Icon icon="iconamoon:home" /></router-link>
        <router-link to="/notes"><Icon icon="hugeicons:profile" /></router-link>
        <router-link to="/notes/new-note"><Icon icon="tabler:plus" /></router-link>
        <router-link to="/gallery"><Icon icon="solar:gallery-bold" /></router-link>
        <router-link to="/stats"><Icon icon="material-symbols:timeline" /></router-link>
    </nav>
</template>


<style>



</style>