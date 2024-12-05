<script lang="ts" setup>
import LogoIcon from './icons/logo-icon.vue';


</script>

<template>
    <aside>
        <LogoIcon class="text-4xl text-invert-bg"/>

        <RouterLink to="/home">
            <i class="pi pi-home"></i>
        </RouterLink>

        <RouterLink to="/notes">
            <i class="pi pi-pen-to-square"></i>
        </RouterLink>

        <RouterLink to="/profile">
            <i class="pi pi-user"></i>
        </RouterLink>
    </aside>
</template>


<style scoped>

aside{
    @apply flex md:flex-col items-center gap-2 p-4 shadow-lg w-fit h-fit bg-primary-bg rounded-full border-invert-bg-50 border-2
}

i{
    @apply text-2xl duration-300 hover:scale-125
}
.router-link-active{
    @apply text-primary
}

</style>