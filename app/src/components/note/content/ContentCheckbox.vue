<script lang="ts" setup>
import CheckBox from '@/components/CheckBox.vue'
import type { Checkbox } from '@/entities/note'

const props = defineProps<{
    element: Checkbox
}>()

</script>

<template>
    <div class="checkbox">
        <CheckBox v-model="element.checked" disabled="true" class="inline" /> 
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
    </div>
</template>


<style scoped>

.checkbox>span{
    @apply inline
}

</style>