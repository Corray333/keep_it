<script setup>
import { ref, onBeforeMount } from 'vue'
import NoteCard from '../components/NoteCard.vue'

const notes = ref([{}])

onBeforeMount(() => {
    notes.value[0] = JSON.parse(`{
        "id":"88004ea3-1c51-4802-b54d-62678db7563b",
        "source":"tg",
        "icon":{
            "type": "icon",
            "icon":"mingcute:telegram-fill",
            "color":"#3DA7F4"
        },
        "from":{
            "text":"Jus so Soon",
            "link":"https://t.me/jun_so_soon"
        },
        "tags":[
            {
            "text":"books",
            "color":"8BFF89"
            },
            {
            "text":"review",
            "color":"FFC75A"
            }
        ],
        "font":"Mono",
        "created_at":1718712074,
        "copied_at":1718718074,
        "timestamps":[
            {
            "timestamp":1718728074,
            "notification":true
            }
        ],
        "cover":"https://images.unsplash.com/photo-1699695726278-f3f57f9356c4?q=80&w=1932&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D",
        "content":[
            {
            "type":"h1",
            "rich_text":[
                {
                "plain_text":"Test text",
                "link":"",
                "text_style":{
                    "color":"FFC75A",
                    "italic":false,
                    "underline":false,
                    "cross_out":true
                }
                },
                {
                "plain_text":", continue of text",
                "text_style":{
                    "width":"normal",
                    "italic":true,
                    "underline":false,
                    "cross_out":false
                }
                }
            ]
            },
            {
                "type":"checkbox",
                "checked":true,
                "rich_text":[
                        {
                        "plain_text":"Test text",
                        "text_style":{
                            "width":"bold",
                            "italic":false,
                            "underline":false,
                            "cross_out":true
                        }
                        }
                    ]
                },
            {
            "type":"img",
            "src":"https://images.unsplash.com/photo-1718439111428-f6ef86aae18d",
            "caption":"Sample text"
            }
        ]
        }
`)
})
</script>

<template>
    <section>
        <article>
            <h1>Home</h1>
            <NoteCard v-for="(note, i) of notes" :key="i" :note="note" />
        </article>
    </section>
</template>


<style></style>