<script lang="ts" setup>
import type { Image } from '@/entities/note';
import { Image as Img } from 'primevue';


const props = defineProps<{
    element: Image
}>()

</script>

<template>
    <div class="image" :style="{justifyContent: element.align}">
        <Img :src="element.src" preview :style="{width: element.width+'%'}"/>
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