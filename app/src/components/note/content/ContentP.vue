<script lang="ts" setup>
import type { Text } from '@/entities/note'


defineProps<{
    element: Text
}>()

</script>

<template>
    <p>
        <span v-if="element.rich_text.meta && element.rich_text.meta.length">
            <span>{{ element.rich_text.plain_text.substring(0, element.rich_text.meta[0].offset) }}</span>
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
    </p>
</template>


<style scoped>

*{
    @apply whitespace-pre-line break-words
}

</style>
