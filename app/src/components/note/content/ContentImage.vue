<script lang="ts" setup>
import type { Image } from '@/entities/note';


const props = defineProps<{
    element: Image
}>()

</script>

<template>
    <div class="image" :style="{justifyContent: element.align}">
        <img :src="element.src" alt="" :style="{width: element.width+'%'}">
    </div>
</template>


<style scoped>

.image{
    @apply w-full flex
}
img{
    @apply rounded-lg
}

</style>