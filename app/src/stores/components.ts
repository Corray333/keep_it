import { ref } from 'vue'
import { defineStore } from 'pinia'
import type { MessageI } from '@/entities/message'
import { MessageType } from '@/entities/message'

const showMessageTime = 5000

export const useComponentsStore = defineStore('components', () => {
  const messages = ref<MessageI[]>([])

  const newMessage = (message: MessageI) => {
    messages.value.push(message)
    setTimeout(()=>{
      messages.value.shift()
    }, showMessageTime)
  }

  const lang = ref<string>('en')
  const setLang = (newLang: string) => {
    lang.value = newLang
  }

  const newError = (error?: string) =>{
    if (error == undefined){
      error = "Unknown error"
    }
    newMessage({
      type: MessageType.ERROR,
      text: error
    })
  }

  return { messages, newMessage, newError, lang, setLang }
})
