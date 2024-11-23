<script lang="ts" setup>
import type { Tag } from '@/entities/note'
import PlusIcon from '../icons/plus-icon.vue'
import { useNotesStore } from '@/stores/notes';
import { NoteService } from '@/service/note';
import { computed, onBeforeMount, ref } from 'vue';
import CloseIcon from '../icons/close-icon.vue';
import { adjustHexColor } from '@/helpers/color';

const notesStore = useNotesStore()

const noteService = new NoteService()

onBeforeMount(async ()=>{
    const newTags = await noteService.getTags()
    if (newTags) notesStore.tags = newTags
    generateRandomColor()
})


const props = defineProps<{
    tags: Tag[],
    noteID: string
}>()

const filteredTags = computed(()=>{
    if (!props.tags) return notesStore.tags
    return notesStore.tags.filter(tag => !props.tags.find(t => t.text === tag.text))
})

const generateRandomColor = ()=>{
    const colors = [
  "#F5A6B2", // мягкий темно-розовый
  "#D68C99", // насыщенный розовый
  "#C6A69B", // теплый бежево-коричневый
  "#B5A28D", // глубокий светлый бежевый
  "#E2C6A1", // золотисто-бежевый
  "#D39C72", // теплый коричнево-желтый
  "#CDAF88", // мягкий карамельный
  "#A7C7C4", // темный мятный
  "#A8B9D4", // глубокий голубой
  "#B79FC8", // пастельно-фиолетовый с насыщенностью
  "#7EC8B2", // насыщенный мятный
  "#9E7BB3", // глубокий лавандовый
  "#C3B7A3", // теплый светло-коричневый
  "#6C7C6E", // темный оливково-зеленый
  "#9F8B98", // насыщенный сиреневый
  "#7E9AA2", // темный серо-голубой
  "#B4C2A8", // глубокий оливковый
  "#9F9F88", // теплый серо-желтый
  "#A3D0B7", // насыщенный зеленый
  "#8360A0", // темный фиолетовый
  "#7D9F9A", // глубокий бирюзово-зеленый
  "#7E5C61", // темный розово-коричневый
  "#A0D1C0", // насыщенный морской зеленый
  "#A9838B", // розово-коричневый
  "#6C7B6F", // глубокий серо-зеленый
  "#857F7C", // темный серый с зелеными оттенками
  "#A86D79", // теплый терракотовый
  "#7C7D9E"  // темный голубовато-серый
]

    return colors[Math.floor(Math.random() * colors.length)]
}

const newTag = ref<string>('')

const createNewTag = async ()=>{
    if (newTag.value){
        const newColor = generateRandomColor()
        const done = await noteService.addTagToNote(props.noteID, newTag.value, newColor, !notesStore.tags.find(tag => tag.text === newTag.value))
        if (done){
            const tag = {text: newTag.value, color: newColor}
            notesStore.tags.push(tag)
            if (!props.tags) props.tags = []
            props.tags.push(tag)
            newTag.value = ''
        }
    }
}

const removeTagFromNote = async (tag: Tag)=>{
    const done = await noteService.removeTagFromNote(props.noteID, tag.text)
    if (done){
        props.tags.splice(props.tags.indexOf(tag), 1)
    }
}

const addTagToNote = async (tag: Tag)=>{
    const done = await noteService.addTagToNote(props.noteID, tag.text, tag.color, !notesStore.tags.find(t => t.text === tag.text))
    if (done){
        props.tags.push(tag)
    }
}

const showAddTagMenu = ref(false)

const closeAddTagMenu = ()=>{
    showAddTagMenu.value = false
}

</script>

<template>
    <div class="tags">
        <div class="added-tags" v-if="tags">
            <span v-for="(tag, i) of tags.slice(0, 3)" :key="`tag${i}`" class="tag group" :style="{backgroundColor:tag.color}">
                <div class="tag-label" :style="{backgroundColor:tag.color}"><p>{{tag.text}}</p></div>
            </span>
        </div>
        <span class="tag new-tag bg-secondary-bg" @click.stop="showAddTagMenu = !showAddTagMenu" v-click-outside="closeAddTagMenu">
            <PlusIcon/>

            <div class="add-tag-menu" :class="{hidden: !showAddTagMenu}">
                <input type="text" placeholder="Tag" v-model="newTag" @keyup.enter="createNewTag">
                <div class="add-tag-menu-list">
                    <span class="select-tag" v-for="(tag, i) of tags" :key="i" @click="removeTagFromNote(tag)" :style="{backgroundColor: tag.color, color: adjustHexColor(tag.color, 30, 40)}"><p>{{ tag.text }}</p><CloseIcon /></span>
                    <span class="select-tag" v-for="(tag, i) of filteredTags" :key="i" @click="addTagToNote(tag)" :style="{backgroundColor: tag.color, color: adjustHexColor(tag.color, 30, 40)}"><p>{{ tag.text }}</p></span>
                </div>
            </div>
        </span>
    </div>
</template>


<style scoped>

.tags, .added-tags{
    @apply flex gap-1;
}

.tag{
    @apply w-4 h-4 rounded-full relative flex justify-center cursor-pointer duration-300;
}

.tag-label{
    @apply absolute bottom-full right-0 whitespace-nowrap mb-1 px-2 rounded-full bg-secondary-bg shadow-lg scale-0 group-hover:scale-100 origin-bottom-right duration-300;
}
.tag-label>p{
    @apply invert
}

.add-tag-menu{
    @apply absolute z-20 flex flex-col gap-2 p-2 bg-primary-bg rounded-2xl shadow-lg w-64 right-full ml-2 duration-300 will-change-transform origin-top-right;
}

.add-tag-menu>*{
    @apply w-full;
}

.add-tag-menu-list{
    @apply flex gap-2 flex-wrap
}



.add-tag-menu.hidden{
    @apply scale-0
}

</style>