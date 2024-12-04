<script lang="ts" setup>

import CheckBox from '@/components/CheckBox.vue';
import CloseIcon from '@/components/icons/close-icon.vue';
import NoteCard from '@/components/note/NoteCard.vue';
import NoteCardSkeleton from '@/components/note/NoteCardSkeleton.vue';
import type { Category } from '@/entities/category';
import { NoteTemplate, type Note, type Tag } from '@/entities/note';
import { adjustHexColor } from '@/helpers/color';
import { CategoryService } from '@/service/category';
import { NoteService } from '@/service/note';
import { useNotesStore } from '@/stores/notes';
import { ContextMenu, Dialog, Toolbar, Tree } from 'primevue';
import { computed, onBeforeMount, ref, watch } from 'vue';

const noteService = new NoteService()

const notesUndefined = [new NoteTemplate(), new NoteTemplate(), new NoteTemplate(), new NoteTemplate(), new NoteTemplate(), new NoteTemplate(), new NoteTemplate(), new NoteTemplate(), new NoteTemplate(),]

const notes = ref<Note[]>([])

const undefineNotes = () => {
    notes.value = notesUndefined
}

const pickedTags = ref<Tag[]>([])

const filteredTags = computed(() => {
    if (!pickedTags.value) return notesStore.tags
    return notesStore.tags.filter(tag => !pickedTags.value.find(t => t.text === tag.text))
})

const pickTag = (tag: Tag) => {
    selectedNotes.value = []
    pickedTags.value.push(tag)
}

const unpickTag = (tag: string) => {
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
    undefineNotes()
    notes.value = await noteService.getNotes(offset.value, pickedTags.value.map(tag => tag.text), pickedCategory.value)
}

onBeforeMount(async () => {
    undefineNotes()
    await fetchCategories()
    fetchNotes()
})

const notesStore = useNotesStore()

const selectedNotes = ref<string[]>([])

// Ссылка на контекстное меню
const contextMenu = ref<{show: (event: MouseEvent)=>void}>();

// Текущая заметка
const currentNote = ref<Note>();

const menuItems = ref<{label: string, command: () => void, icon: string}[]>([]);

const deselectNote = (note: Note) => {
    selectedNotes.value = selectedNotes.value.filter(id => id !== note.id);
};

// Показ контекстного меню
const showContextMenu = (event: MouseEvent, note: Note) => {
    currentNote.value = note; // Устанавливаем текущую заметку
    if (selectedNotes.value.find(id => id === note.id)){
        menuItems.value = [
            {
                label: "Отменить выбор",
                icon: 'pi pi-check-circle',
                command: () => deselectNote(note),
            },
        ];
    } else {
        menuItems.value = [
        {
            label: "Выбрать заметку",
            icon: 'pi pi-circle',
            command: () => selectNote(note),
        },
    ];
    }
    

    // Показываем меню
    if (contextMenu.value) contextMenu.value.show(event)
};

// Добавление заметки в список выбранных
const selectNote = (note: Note) => {
    if (!selectedNotes.value.includes(note.id)) {
        selectedNotes.value.push(note.id);
    }
};

const selectNoteToggle = (noteID: string) =>{
    if (selectedNotes.value.includes(noteID)){
        selectedNotes.value = selectedNotes.value.filter(id => id !== noteID)
    } else {
        selectedNotes.value.push(noteID)
    }
}

const deselectAll = ()=>{
    selectedNotes.value = []
}

const selectAllToggle = ()=>{
    if (allSelected.value) selectedNotes.value = []
    else selectedNotes.value = notes.value.map(note => note.id)
}

const deleteSelectedNotes = async()=>{
    undefineNotes()
    if (!await noteService.deleteNotes(selectedNotes.value)) return
    deselectAll()
    fetchNotes()
}

const allSelected = computed(()=>{
    return selectedNotes.value.length == notes.value.length
})


const fetchCategories = async () => {
    notesStore.categories = await CategoryService.getCategories()
}

const newCategoryName = ref<string>('')
const newCategoryParent = ref()

const pickedCategory = ref('')

watch(newCategoryParent, (newVal)=>{
    if (Object.keys(newVal).length>0) {
        pickedCategory.value = Object.keys(newVal)[0]
        fetchNotes()
    } else {
        pickedCategory.value = ''
        fetchNotes()
    }
    
})

const createCategory = async () => {
    if (!newCategoryName.value) return
    let parent = ''
    if (newCategoryParent.value && Object.keys(newCategoryParent.value).length>0) parent = Object.keys(newCategoryParent.value)[0]
    await CategoryService.createCategory(newCategoryName.value, parent)
    fetchCategories()
    closeNewCategoryDialog()
}

const showNewCategoryDialog = ref<boolean>(false)

const openNewCategoryDialog = () => {
    showNewCategoryDialog.value = true
}

const closeNewCategoryDialog = () => {
    showNewCategoryDialog.value = false
}



</script>

<template>
    <Dialog v-model:visible="showNewCategoryDialog" modal header="New category" >
        <input v-model="newCategoryName" type="text" placeholder="Category name">
        <div class="flex gap-2">
            <button class="button" @click="closeNewCategoryDialog">Cancel</button>
            <button class="button" @click="createCategory">Create</button>
        </div>
    </Dialog>

    <section class="notes-page">
        <h1>Notes</h1>

        <div class="category">
            <div class="flex gap-2 items-center w-fit">
                <h2>Categories</h2>
                <button @click="openNewCategoryDialog"><i class="pi pi-plus"></i></button>
            </div>

            <Tree :value="notesStore.categories" v-model:selection-keys="newCategoryParent" selection-mode="single" class="w-fit min-w-64"/>
        </div>

        <div class="tags">
            <h2>Tags</h2>
            <div class="tags-picker">
                <span class="select-tag" v-for="(tag, i) of pickedTags" :key="i" @click="unpickTag(tag.text)"
                    :style="{ backgroundColor: tag.color, color: adjustHexColor(tag.color, 30, 40) }">
                    <p>{{ tag.text }}</p>
                    <CloseIcon />
                </span>
                <span class="select-tag" v-for="(tag, i) of filteredTags" :key="i" @click="pickTag(tag)"
                    :style="{ backgroundColor: tag.color, color: adjustHexColor(tag.color, 30, 40) }">
                    <p>{{ tag.text }}</p>
                </span>
            </div>
        </div>

        <div class="home-body">
            <Toolbar>
                <template #start>
                    <button @click="selectAllToggle" class=" flex text-invert-bg items-center gap-2"><CheckBox :model-value="allSelected" disabled/>Select all</button>
                </template>

                <template #end>
                    <button @click="deleteSelectedNotes" class=" button"><i class="pi pi-trash"></i></button>
                </template>
            </Toolbar>
            <div class="notes w-full">
                <div v-for="note of notes" :key="note.id" @contextmenu="showContextMenu($event, note)" :class="{shake: selectedNotes.length > 0}">
                    <NoteCard :note="note" :selectable="selectedNotes.length > 0" :selected="selectedNotes.includes(note.id)" @select-toggle="selectNoteToggle(note.id)" v-if="note.id != ''" />
                    <NoteCardSkeleton v-else/>
                </div>
                <ContextMenu 
                :model="menuItems"
                ref="contextMenu" 
                />
            </div>
        </div>


    </section>
</template>

<style>

.shake {
    animation: shake 0.25s infinite;
}

@keyframes shake {
    0%, 100% { transform: rotate(0.5); }
    50% { transform: rotate(-0.5deg); }
}

.p-toolbar{
    background: var(--primary-bg) !important;
}

</style>


<style scoped>


.notes-page {
    @apply flex flex-col gap-4 p-4 xl:p-8 xl:px-16
}

.home-body {
    @apply flex flex-col gap-4
}

.notes {
    @apply grid grid-cols-1 lg:grid-cols-2 xl:grid-cols-3 3xl:grid-cols-4 4xl:grid-cols-5 gap-4 w-full
}

.tags {
    @apply flex flex-col gap-2
}

.tags-picker {
    @apply w-full flex flex-wrap gap-2 text-black
}
</style>