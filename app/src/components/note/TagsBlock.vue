<script lang="ts" setup>
import type { Tag } from '@/entities/note'
import PlusIcon from '../icons/plus-icon.vue'
import { useNotesStore } from '@/stores/notes';
import { NoteService } from '@/service/note';
import { computed, onBeforeMount, ref } from 'vue';
import CloseIcon from '../icons/close-icon.vue';

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
    const colors = ['#FFC0CB', '#FFB6C1', '#FF69B4', '#FF1493', '#DB7093', '#C71585', '#FFA07A', '#FA8072', '#E9967A', '#F08080', '#CD5C5C', '#DC143C', '#B22222', '#8B0000', '#FF0000', '#FF6347', '#FF4500', '#FF8C00', '#FFA500', '#FFD700', '#FFFF00', '#FFFFE0', '#FFFACD', '#FAFAD2', '#FFEFD5', '#FFE4B5', '#FFDAB9', '#EEE8AA', '#F0E68C', '#BDB76B', '#ADFF2F', '#7FFF00', '#7CFC00', '#00FF00', '#32CD32', '#98FB98', '#90EE90', '#00FA9A', '#00FF7F', '#3CB371', '#2E8B57', '#228B22', '#008000', '#006400', '#9ACD32', '#6B8E23', '#808000', '#556B2F', '#66CDAA', '#8FBC8F', '#20B2AA', '#008B8B', '#008080', '#00CED1', '#48D1CC', '#40E0D0', '#7FFFD4', '#66CDAA', '#00FA9A', '#00FF7F', '#3CB371', '#2E8B57', '#228B22', '#008000', '#006400', '#9ACD32', '#6B8E23', '#808000', '#556B2F', '#66CDAA', '#8FBC8F', '#20B2AA', '#008B8B', '#008080', '#00CED1', '#48D1CC', '#40E0D0', '#7FFFD4', '#66CDAA', '#00FA9A', '#00FF7F', '#3CB371', '#2E8B57', '#228B22', '#008000', '#006400', '#9ACD32', '#6B8E23', '#808000', '#556B2F']
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
        <span class="tag new-tag bg-secondary-bg" @click.stop="showAddTagMenu = true" v-click-outside="closeAddTagMenu">
            <PlusIcon/>

            <div class="add-tag-menu" :class="{hidden: !showAddTagMenu}">
                <input type="text" placeholder="Tag" v-model="newTag" @keyup.enter="createNewTag">
                <div class="add-tag-menu-list">
                    <span class="select-tag" v-for="(tag, i) of tags" :key="i" @click="removeTagFromNote(tag)" :style="{backgroundColor: tag.color}"><p>{{ tag.text }}</p><CloseIcon /></span>
                    <span class="select-tag" v-for="(tag, i) of filteredTags" :key="i" @click="addTagToNote(tag)" :style="{backgroundColor: tag.color}"><p>{{ tag.text }}</p></span>
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
    @apply w-4 h-4 rounded-full relative flex justify-center cursor-pointer duration-300 hover:scale-125;
}
.new-tag{
    @apply hover:scale-100;
}
.tag-label{
    @apply absolute bottom-full mb-2 px-2 rounded-full bg-secondary-bg shadow-lg scale-0 group-hover:scale-75 origin-bottom duration-300;
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