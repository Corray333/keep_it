<script lang="ts" setup>
import type { Checkbox, Image, Note, Text } from '@/entities/note';
import NoteIcon from '../NoteIcon.vue';
import ContentCheckbox from './content/ContentCheckbox.vue';
import ContentH1 from './content/ContentH1.vue';
import ContentImage from './content/ContentImage.vue';
import TagsBlock from './TagsBlock.vue';
import ContentP from './content/ContentP.vue';

import {  ref } from 'vue';
import { Accordion, AccordionContent, AccordionHeader, AccordionPanel, Image as Img } from 'primevue';
import { TimeFromUnix } from '@/helpers/time';
import CheckBox from '../CheckBox.vue';
import { adjustHexColor } from '@/helpers/color';


const emit = defineEmits(['select-toggle'])

const props = defineProps<{
    note: Note,
    selected: boolean,
    selectable: boolean
}>()

const showModal = ref(false)

const openNote = ()=>{
    if (props.selectable) {
        emit('select-toggle')
    } else {
        showModal.value = true
    }
}



</script>

<template>
    <Transition name="delay">
        <section @click.self="showModal = false" @contextmenu.stop v-show="showModal" class="note-page-wrapper">
            <Transition name="scale">
                <div v-if="showModal" class="note-page-container">
                    <div class="note-page">

                        <Img v-show="note.cover" alt="Image" width="250" preview>
                            <template #image>
                                <img :src="note.cover" class="note-page-cover" alt="image" />
                            </template>
                            <template #preview="slotProps: any">
                                <img :src="note.cover" alt="preview" :style="slotProps.style" @click="slotProps.onClick" />
                            </template>
                        </Img>
                        
                        <div class="note-page-header">
                            <Accordion unstyled class="w-full" expandIcon="pi" collapseIcon="pi">
                                <AccordionPanel unstyled value="0" class="w-full">
                                    <AccordionHeader unstyled class="w-full">
                                        <div class="note-page-header-tab">
                                            <div class="note-page-header-label">
                                                <NoteIcon :icon="note.icon.data" />
                                                <a :href="note.original" target="_blank"><p>{{ note.title }}</p></a>
                                            </div>
                
                                            <TagsBlock :tags="note.tags" :noteID="note.id" /> 
                                        </div>
                                    </AccordionHeader>
                                    <AccordionContent unstyled>
                                        <div class="note-page-header-more">
                                            <div class="tags-picker">
                                                <span class="select-tag" v-for="(tag, i) of note.tags" :key="i" 
                                                    :style="{ backgroundColor: tag.color, color: adjustHexColor(tag.color, 30, 40) }">
                                                    <p>{{ tag.text }}</p>
                                                    <CloseIcon />
                                                </span>
                                            </div>
                                            <p>Copied at: {{ TimeFromUnix(note.copiedAt) }}</p>
                                            <p>Created at: {{ TimeFromUnix(note.createdAt) }}</p>
                                            <p>Source: {{ note.source }}</p>
                                        </div>
                                    </AccordionContent>
                                </AccordionPanel>
                            </Accordion>
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
                </div>
            </Transition>
        </section>
    </Transition>

        <div @click="openNote" class="note-card">
            <div class="note-card-header">
                <div class="note-card-header-label">
                    <CheckBox v-if="selectable" :model-value="selected" :disabled="false" />
                    <NoteIcon :icon="note.icon.data" />
                    <a :href="note.original" target="_blank" @click.stop><p>{{ note.title }}</p></a>
                </div>
    
                <TagsBlock :tags="note.tags" :noteID="note.id" />
            </div>
    
            <div class="note-card-body">
                <img v-show="note.cover" class="note-card-body-cover" :src="note.cover">
    
                <div class="note-card-body-content w-full">
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

.tags-picker {
    @apply w-full flex flex-wrap gap-2 text-black
}

.note-card{
    @apply relative z-0 bg-primary-bg rounded-2xl shadow-lg border-2 border-invert-bg-50 w-full cursor-pointer;
    @apply duration-300 hover:shadow-lg hover:shadow-invert-bg-opacity;
    word-break: break-word;
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
    @apply flex gap-2 font-medium w-full
}

.note-card-header-label p, .note-card-header-label a{
    @apply overflow-hidden whitespace-nowrap text-ellipsis text-nowrap
}

.note-page-wrapper{
    @apply fixed w-screen h-screen p-8 top-0 left-0 bg-black bg-opacity-50 z-50 flex justify-center items-center;
}
.note-page-container{
    @apply max-h-full h-fit w-full lg:w-2/3 bg-primary-bg 2xl:w-1/2  overflow-y-scroll border-2 border-invert-bg rounded-2xl;
}
.note-page{
    @apply w-full flex flex-col rounded-2xl bg-primary-bg;
}
.note-page-cover{
    @apply w-full h-96 object-cover rounded-t-2xl;
}
.note-page-header{
    @apply p-4 flex items-center gap-2 border-b-2 w-full justify-between;
}
.note-page-header-label{
    @apply flex items-center gap-2 text-xl
}
.note-page-header-tab{
    @apply flex justify-between items-center
}
.note-page-content{
    @apply p-4 flex flex-col gap-2;
}


</style>