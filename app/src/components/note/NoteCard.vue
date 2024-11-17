<script lang="ts" setup>
import type { Checkbox, Image, Note, Text } from '@/entities/note';
import NoteIcon from '../NoteIcon.vue';
import ContentCheckbox from './content/ContentCheckbox.vue';
import ContentH1 from './content/ContentH1.vue';
import ContentImage from './content/ContentImage.vue';
import TagsBlock from './TagsBlock.vue';
import ContentP from './content/ContentP.vue';
import { ref } from 'vue';


defineProps<{
    note: Note
}>()

const showModal = ref(false)



</script>

<template>
    <Transition name="delay">
        <section @click.self="showModal = false" v-show="showModal" class="note-page-wrapper">
            <Transition name="scale">
                <div v-if="showModal" class="note-page">
                    <img class="note-page-cover" :src="note.cover">

                    <div class="note-page-label">
                        <NoteIcon :icon="note.icon.data" />
                        <a :href="note.original" target="_blank"><p>{{ note.title }}</p></a>
                    </div>
        
                    <div class="note-page-content">
                        <p v-for="(content_el, i) of note.content" :key="`el${i}`">
                            <ContentH1 v-if="content_el.type == 'h1'" :element="(content_el as Text)" />
                            <ContentP v-if="content_el.type == 'p'" :element="(content_el as Text)" />
                            <ContentCheckbox v-if="content_el.type == 'checkbox'" :element="(content_el as Checkbox)" />
                            <ContentImage v-if="content_el.type == 'img'" :element="(content_el as Image)"/>
                        </p>
                    </div>
                </div>
            </Transition>
        </section>
    </Transition>

    <div @click="showModal = true" class="note-card">
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
                    <ContentH1 v-if="content_el.type == 'h1'" :element="(content_el as Text)" />
                    <ContentP v-if="content_el.type == 'p'" :element="(content_el as Text)" />
                    <ContentCheckbox v-if="content_el.type == 'checkbox'" :element="(content_el as Checkbox)" />
                    <ContentImage v-if="content_el.type == 'img'" :element="(content_el as Image)"/>
                </p>
            </div>
        </div>
    </div>
</template>


<style scoped>

.delay-enter-active,
.delay-leave-active {
    transition: opacity 0.5s ease;
}

.delay-enter-from,
.delay-leave-to {
    opacity: 0;
}

.scale-enter-active,
.scale-leave-active {
    transition: transform 0.5s ease;
}

.scale-enter-from,
.scale-leave-to {
    transform: scale(0);
}

.note-card{
    @apply relative z-0 bg-primary-bg rounded-2xl shadow-lg border-2 border-invert-bg-50 w-full cursor-pointer;
    @apply duration-300 hover:scale-105
}

.note-card-header{
    @apply p-4 border-b-2 border-invert-bg-50 flex items-center justify-between;
}

.note-card-body{
    @apply p-4 flex h-36 overflow-hidden gap-2 text-sm;
}

.note-card-body-cover{
    @apply w-28 h-28 object-cover aspect-square rounded-lg;
}

.note-card-header-label{
    @apply flex gap-2 font-medium
}

.note-page-wrapper{
    @apply fixed w-screen h-screen p-8 top-0 left-0 bg-black bg-opacity-50 z-50 flex justify-center items-center;
}
.note-page{
    @apply w-full lg:w-2/3 2xl:w-1/2 flex flex-col border-2 border-invert-bg rounded-2xl;
}
.note-page-cover{
    @apply w-full h-96 object-cover rounded-t-2xl;
}
.note-page-label{
    @apply p-4 flex items-center gap-2 text-xl border-b-2;
}
.note-page-content{
    @apply p-4 flex flex-col gap-2;
}

</style>