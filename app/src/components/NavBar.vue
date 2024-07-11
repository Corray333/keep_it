<script setup>
import { Icon } from '@iconify/vue'

</script>

<template>
    <nav class="flex w-full absolute bottom-0 md:flex-col md:w-fit md:p-5 md:top-0 md:relative md:rounded-xl h-fit bg-black gap-5 p-3 text-2xl justify-around bottom-bar-shadow rounded-t-xl">
        <router-link to="/home" class=" hover:text-accent duration-300"><Icon icon="iconamoon:home" /></router-link>
        <router-link to="/notes" class=" hover:text-accent duration-300"><Icon icon="hugeicons:profile" /></router-link>
        <router-link to="/notes/new-note" class=" hover:text-accent duration-300"><Icon icon="tabler:plus" /></router-link>
        <router-link to="/gallery" class=" hover:text-accent duration-300"><Icon icon="solar:gallery-bold" /></router-link>
        <router-link to="/stats" class=" hover:text-accent duration-300"><Icon icon="material-symbols:timeline" /></router-link>
    </nav>
</template>


<style>



</style>