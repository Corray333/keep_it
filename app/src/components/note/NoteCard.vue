<script lang="ts" setup>
import type { Checkbox, H1, Image, Note } from '@/entities/note';
import NoteIcon from '../NoteIcon.vue';
import ContentCheckbox from './content/ContentCheckbox.vue';
import ContentH1 from './content/ContentH1.vue';
import ContentImage from './content/ContentImage.vue';
import TagsBlock from './TagsBlock.vue';


defineProps<{
    note: Note
}>()

</script>

<template>
    <div class="note-card">
        <div class="note-card-header">
            <div class="note-card-header-label">
                <NoteIcon :icon="note.icon.data" />
                <a :href="note.original" target="_blank"><p>{{ note.title }}</p></a>
            </div>

            <TagsBlock :tags="note.tags" :noteID="note.id" />
        </div>

        <div class="note-card-body">
            <img class="note-card-body-cover" :src="note.cover">

            <div class="note-card-body-content">
                <p v-for="(content_el, i) of note.content" :key="`el${i}`">
                    <ContentH1 v-if="content_el.type == 'h1'" :element="(content_el as H1)" />
                    <ContentCheckbox v-if="content_el.type == 'checkbox'" :element="(content_el as Checkbox)" />
                    <ContentImage v-if="content_el.type == 'img'" :element="(content_el as Image)"/>
                </p>
            </div>
        </div>
    </div>
</template>


<style scoped>

.note-card{
    @apply bg-primary-bg rounded-2xl shadow-lg border-2 border-invert-bg-50 w-full;
}

.note-card-header{
    @apply p-2 border-b-2 border-invert-bg-50 flex items-center justify-between;
}

.note-card-body{
    @apply p-4 flex h-36 overflow-hidden gap-2 text-sm;
}

.note-card-body-cover{
    @apply w-28 h-28 object-cover rounded-lg;
}

.note-card-header-label{
    @apply flex gap-2
}



</style>