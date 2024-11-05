import { ref } from 'vue'
import { defineStore } from 'pinia'
import type { Message } from '@/entities/message'
import { MessageType } from '@/entities/message'

const showMessageTime = 3000

export const useComponentsStore = defineStore('components', () => {
  const messages = ref<Message[]>([])

  const newMessage = (message: Message) => {
    messages.value.push(message)
    setTimeout(()=>{
      messages.value.shift()
    }, showMessageTime)
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

  return { newMessage, newError }
})
