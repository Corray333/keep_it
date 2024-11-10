import { ref } from 'vue'
import { defineStore } from 'pinia'
import type { Tag } from '@/entities/note'

export const useNotesStore = defineStore('notes', () => {
  const tags = ref<Tag[]>([])

  return { tags }
})
