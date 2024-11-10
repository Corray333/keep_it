<script lang="ts" setup>

import NoteCard from '@/components/note/NoteCard.vue';
import type { Note } from '@/entities/note';
import { NoteService } from '@/service/note';
import { onBeforeMount, ref } from 'vue';

const noteService = new NoteService()

const notes = ref<Note[]>([])


const fetchNotes = async () => {
    notes.value = await noteService.getNotes(notes.value.length)
}

onBeforeMount(()=>{
    fetchNotes()
})



</script>

<template>
    <section class="home">
        <h1 class="text-4xl font-bold text-secondary-bg">Home</h1>
        
        <div class="home-body">
            <div class="notes w-full">
                <NoteCard v-for="note of notes" :key="note.id" :note="note"/>
            </div>
        </div>

    </section>
</template>


<style scoped>

.home{
    @apply flex flex-col min-h-screen w-full p-4
}

.home-body{
    @apply grid grid-cols-4 w-full
}

.notes{
    @apply relative flex flex-col gap-2
}

</style>