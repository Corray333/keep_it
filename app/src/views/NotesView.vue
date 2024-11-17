<script lang="ts" setup>

import CloseIcon from '@/components/icons/close-icon.vue';
import NoteCard from '@/components/note/NoteCard.vue';
import type { Note, Tag } from '@/entities/note';
import { NoteService } from '@/service/note';
import { useNotesStore } from '@/stores/notes';
import { computed, onBeforeMount, ref, watch } from 'vue';

const noteService = new NoteService()

const notes = ref<Note[]>([])

const pickedTags = ref<Tag[]>([])

const filteredTags = computed(()=>{
    if (!pickedTags.value) return notesStore.tags
    return notesStore.tags.filter(tag => !pickedTags.value.find(t => t.text === tag.text))
})

const unpickTag = (tag: string)=>{
    pickedTags.value = pickedTags.value.filter(t => t.text !== tag)
}

watch(
    () => pickedTags.value,
    () => {
        fetchNotes()
    },
    { deep: true }
)

const offset = ref(0)


const fetchNotes = async () => {
    notes.value = await noteService.getNotes(offset.value, pickedTags.value.map(tag => tag.text))
}

onBeforeMount(()=>{
    fetchNotes()
})

const notesStore = useNotesStore()


</script>

<template>
    <section class="notes-page">
        <h1>Notes</h1>

        <div class="tags">
            <h2>Tags</h2>
            <div class="tags-picker">
                <span class="select-tag" v-for="(tag, i) of pickedTags" :key="i" @click="unpickTag(tag.text)" :style="{backgroundColor: tag.color}"><p>{{ tag.text }}</p><CloseIcon /></span>
                <span class="select-tag" v-for="(tag, i) of filteredTags" :key="i" @click="pickedTags.push(tag)" :style="{backgroundColor: tag.color}"><p>{{ tag.text }}</p></span>
            </div>
        </div>
        
        <div class="home-body">
            <div class="notes w-full">
                <NoteCard v-for="note of notes" :key="note.id" :note="note"/>
            </div>
        </div>

    </section>
</template>


<style scoped>

.notes-page{
    @apply flex flex-col gap-4 p-4
}

.home-body{
    
}

.notes{
    @apply grid grid-cols-1 lg:grid-cols-2 xl:grid-cols-3 gap-4 w-full
}

.tags{
    @apply flex flex-col gap-2
}

.tags-picker{
    @apply w-full flex flex-wrap gap-2 text-black
}

</style>