<script lang="ts" setup>
import type { H1 } from '@/entities/note'


defineProps<{
    element: H1
}>()

</script>

<template>
    <h1>
        <span 
            v-for="(piece, i) of element.rich_text.meta" :key="i">
            <span
                :style="{ color: piece.color, fontWeight: piece.weight }"
                >
                {{ element.rich_text.plain_text.substring(piece.offset, piece.offset + piece.length) }}
            </span>
            <span v-if="i<element.rich_text.meta.length-1">
                {{ element.rich_text.plain_text.substring(piece.offset + piece.length, element.rich_text.meta[i+1].offset) }}
            </span>
            <span v-else>
                {{ element.rich_text.plain_text.substring(piece.offset + piece.length, element.rich_text.plain_text.length) }}
            </span>
        </span>
    </h1>
</template>


<style scoped>

h1{
    white-space: pre-line;
}

</style>