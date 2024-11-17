<script lang="ts" setup>
import type { Text } from '@/entities/note'


defineProps<{
    element: Text
}>()

</script>

<template>
    <h1>
        <span v-if="element.rich_text.meta">
            <span 
            v-for="(piece, i) of element.rich_text.meta" :key="i">
                <a v-if="piece.link" :href="piece.link" target="_blank">
                    <span
                        :style="{ color: piece.color, fontWeight: piece.weight }"
                        >
                        {{ element.rich_text.plain_text.substring(piece.offset, piece.offset + piece.length) }}
                    </span>
                </a>
                <span v-else
                    :style="{ color: piece.color, fontWeight: piece.weight }"
                    >
                    {{ element.rich_text.plain_text.substring(piece.offset, piece.offset + piece.length) }}
                </span>
                <span v-if="i < element.rich_text.meta.length-1">
                    {{ element.rich_text.plain_text.substring(piece.offset + piece.length, element.rich_text.meta[i+1].offset) }}
                </span>
                <span v-else>
                    {{ element.rich_text.plain_text.substring(piece.offset + piece.length, element.rich_text.plain_text.length) }}
                </span>
            </span>
        </span>
        <span v-else>{{ element.rich_text.plain_text}}</span>
    </h1>
</template>


<style scoped>

h1{
    white-space: pre-line;
}

</style>