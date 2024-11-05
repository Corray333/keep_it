import { ref } from 'vue'
import { defineStore } from 'pinia'
import type { User } from '@/entities/user'

export const useAccountStore = defineStore('account', () => {
  const authorization = ref<string>()
  const user = ref<User>()

  return { authorization, user }
})
