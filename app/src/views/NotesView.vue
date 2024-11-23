<script lang="ts" setup>

import CloseIcon from '@/components/icons/close-icon.vue';
import NoteCard from '@/components/note/NoteCard.vue';
import NoteCardSkeleton from '@/components/note/NoteCardSkeleton.vue';
import { NoteTemplate, type Note, type Tag } from '@/entities/note';
import { adjustHexColor } from '@/helpers/color';
import { NoteService } from '@/service/note';
import { useNotesStore } from '@/stores/notes';
import { computed, onBeforeMount, ref, watch } from 'vue';

const noteService = new NoteService()

const notesUndefined = [new NoteTemplate(), new NoteTemplate(), new NoteTemplate(), new NoteTemplate(), new NoteTemplate(), new NoteTemplate(), new NoteTemplate(), new NoteTemplate(), new NoteTemplate(),]

const notes = ref<Note[]>()

const undefineNotes = ()=>{
    // notes.value = notesUndefined
}

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
        undefineNotes()
        fetchNotes()
    },
    { deep: true }
)

const offset = ref(0)


const fetchNotes = async () => {
    notes.value = await noteService.getNotes(offset.value, pickedTags.value.map(tag => tag.text))
}

onBeforeMount(()=>{
    undefineNotes()
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
                <span class="select-tag" v-for="(tag, i) of pickedTags" :key="i" @click="unpickTag(tag.text)" :style="{backgroundColor: tag.color, color: adjustHexColor(tag.color, 30, 40)}"><p>{{ tag.text }}</p><CloseIcon /></span>
                <span class="select-tag" v-for="(tag, i) of filteredTags" :key="i" @click="pickedTags.push(tag)" :style="{backgroundColor: tag.color, color: adjustHexColor(tag.color, 30, 40)}"><p>{{ tag.text }}</p></span>
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
    @apply flex flex-col gap-4 p-4 xl:p-8 xl:px-16
}

.home-body{
    
}

.notes{
    @apply grid grid-cols-1 lg:grid-cols-2 xl:grid-cols-3 3xl:grid-cols-4 4xl:grid-cols-5 gap-4 w-full
}

.tags{
    @apply flex flex-col gap-2
}

.tags-picker{
    @apply w-full flex flex-wrap gap-2 text-black
}

</style>