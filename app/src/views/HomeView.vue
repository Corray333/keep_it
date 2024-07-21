<script setup>
import { ref, onBeforeMount } from 'vue'
import NoteCard from '../components/NoteCard.vue'
import axios from 'axios'

const notes = ref([{}])



// onBeforeMount(() => {

//     notes.value[0] = JSON.parse(`{
//             "id":"88004ea3-1c51-4802-b54d-62678db7563b",
//             "source":"tg",
//             "icon":{
//                 "type": "emoji",
//                 "emoji":"🎉"
//             },
//             "from":{
//                 "text":"Jun so Soon",
//                 "link":"https://t.me/jun_so_soon"
//             },
//             "tags":[
//                 {
//                 "text":"books",
//                 "color":"8BFF89"
//                 },
//                 {
//                 "text":"review",
//                 "color":"FFC75A"
//                 }
//             ],
//             "font":"Mono",
//             "created_at":1718712074,
//             "copied_at":1718718074,
//             "timestamps":[
//                 {
//                 "timestamp":1718728074,
//                 "notification":true
//                 }
//             ],
//             "cover":"https://images.unsplash.com/photo-1699695726278-f3f57f9356c4?q=80&w=1932&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D",
//             "content":[
//                 {
//                     "type":"h1",  
//                     "rich_text":{
//                         "plain_text":"Test text",
//                         "meta": [
//                             {
//                                 "offset":0,
//                                 "length":4,
//                                 "link":"",
//                                 "color":"#FFC75A",
//                                 "width":"bold",
//                                 "italic":false,
//                                 "underline":false,
//                                 "cross_out":true
//                             },
//                             {
//                                 "offset":4,
//                                 "length":5,
//                                 "link":"",
//                                 "italic":false,
//                                 "underline":false,
//                                 "cross_out":true
//                             }
//                         ]
//                         }
//                     },
//                     {
//                     "type":"checkbox",
//                     "checked":false,
//                     "rich_text":{
//                         "plain_text":"Test",
//                         "meta": [
//                             {
//                             "offset":0,
//                             "length":2,
//                             "link":"",
//                             "width":"bold",
//                             "italic":false,
//                             "underline":false,
//                             "cross_out":true
//                             }
//                         ]
//                         }
//                     },
//                 {
//                     "type":"img",
//                     "src":"https://images.unsplash.com/photo-1718439111428-f6ef86aae18d",
//                     "caption":"Sample text"
//                 }
//             ]
//         }
//     `)
//     notes.value[1] = JSON.parse(`{
//             "id":"88004ea3-1c51-4802-b54d-62678db7563b",
//             "source":"tg",
//             "icon":{
//                 "type": "emoji",
//                 "emoji":"🎉"
//             },
//             "from":{
//                 "text":"Jun so Soon",
//                 "link":"https://t.me/jun_so_soon"
//             },
//             "tags":[
//                 {
//                 "text":"books",
//                 "color":"8BFF89"
//                 },
//                 {
//                 "text":"review",
//                 "color":"FFC75A"
//                 }
//             ],
//             "font":"Mono",
//             "created_at":1718712074,
//             "copied_at":1718718074,
//             "timestamps":[
//                 {
//                 "timestamp":1718728074,
//                 "notification":true
//                 }
//             ],
//             "cover":"https://images.unsplash.com/photo-1699695726278-f3f57f9356c4?q=80&w=1932&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D",
//             "content":[
//                 {
//                     "type":"h1",  
//                     "rich_text":{
//                         "plain_text":"Test text",
//                         "meta": [
//                             {
//                                 "offset":0,
//                                 "length":4,
//                                 "link":"",
//                                 "color":"#FFC75A",
//                                 "width":"bold",
//                                 "italic":false,
//                                 "underline":false,
//                                 "cross_out":true
//                             },
//                             {
//                                 "offset":4,
//                                 "length":5,
//                                 "link":"",
//                                 "italic":false,
//                                 "underline":false,
//                                 "cross_out":true
//                             }
//                         ]
//                         }
//                     },
//                     {
//                     "type":"checkbox",
//                     "checked":false,
//                     "rich_text":{
//                         "plain_text":"Test",
//                         "meta": [
//                             {
//                             "offset":0,
//                             "length":2,
//                             "link":"",
//                             "width":"bold",
//                             "italic":false,
//                             "underline":false,
//                             "cross_out":true
//                             }
//                         ]
//                         }
//                     },
//                 {
//                     "type":"img",
//                     "src":"https://images.unsplash.com/photo-1718439111428-f6ef86aae18d",
//                     "caption":"Sample text"
//                 }
//             ]
//         }
//     `)
// })

const loadNotes = async () =>{
    try {
        const {data} = await axios.get(`${import.meta.env.VITE_API_URL}/notes`,{
            headers: {
                'Authorization': store.state.AccessToken,
            }
        })

        notes.value = data.notes
    } catch (error) {
        console.log(error)
    }
}

onBeforeMount(()=>{
    loadNotes()
})

</script>

<template>
    <section class="p-5">
        <section class="top grid md:grid-cols-2 lg:grid-cols-3 2xl:grid-cols-4 gap-2">
            <NoteCard v-for="(note, i) of notes" :key="i" :note="note" />
        </section>
    </section>
</template>


<style scoped></style>