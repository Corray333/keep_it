<script lang="ts" setup>
import type { Note } from '@/entities/note'
import { baseURL } from '@/helpers/data'
import NoteIcon from '../NoteIcon.vue';


const props = defineProps<{
    note: Note
}>()

</script>

<template>
    <div class="note-card">
        <div class="note-card-header">
            <NoteIcon :icon="note.icon.data" /><p>{{ note.title }}</p>
        </div>

        <div class="note-card-body">
            <div class="note-card-body-cover"></div>

            <div class="note-card-body-content">
                <p v-for="(content_el, i) of note.content" :key="`el${i}`">{{ content_el }}</p>
            </div>
        </div>
    </div>
</template>


<style scoped>

.note-card{
    @apply bg-primary-bg rounded-2xl shadow-lg border-2 border-invert-bg-50 w-full;
}

.note-card-header{
    @apply p-2 border-b-2 border-invert-bg-50 flex gap-2 items-center;
}

.note-card-body{
    @apply p-2 flex;
}

</style>