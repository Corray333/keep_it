import { ref } from 'vue'
import { defineStore } from 'pinia'

export const useAccountStore = defineStore('account', () => {
  const authorization = ref<string>()

  const setAuthorization = (token: string) => {
    authorization.value = token
  }

  return { authorization, setAuthorization }
})
