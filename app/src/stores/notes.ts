import { ref } from 'vue'
import { defineStore } from 'pinia'
import type { Tag } from '@/entities/note'
import type { Category } from '@/entities/category'

export const useNotesStore = defineStore('notes', () => {
  const tags = ref<Tag[]>([])
  const categories = ref<Category[]>([])

  return { tags, categories }
})
