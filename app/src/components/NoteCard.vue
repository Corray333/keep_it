<script setup>
import { Icon } from '@iconify/vue'
import { formToJSON } from 'axios';
import { ref, computed } from 'vue'

const props = defineProps(['note'])

const content = computed(()=>{
    const result = []
    for (let note of props.note.content){
        if (note.rich_text){
            let temp = {
                type: note.type,
                rich_text:[]
            }
            let i = 0
            for (let block of note.rich_text.meta){
                if (i<block.offset){
                    temp.rich_text.push({
                        plain_text: note.rich_text.plain_text.substr(i, block.offset),
                        text_style:{
                        }
                    })
                }
                temp.rich_text.push({
                    plain_text: note.rich_text.plain_text.substr(block.offset, block.offset + block.length),
                    text_style:{
                        color: block.color,
                        width:block.width,
                        italic:block.italic,
                        underline:block.underline,
                        cross_out:block.cross_out
                    }
                })
                i += block.length
            }
            if (i < note.rich_text.plain_text.length){
                temp.rich_text.push({
                        plain_text: note.rich_text.plain_text.substr(i, note.rich_text.plain_text.length),
                        text_style:{
                        }
                    })
            }
            result.push(temp)
        } else {
            result.push(note)
        }
    }
    console.log(result)
    return result
})


</script>

<template>
    <article class=" border-2 border-[#8f8f8f] rounded-xl">
       <div class="flex items-center p-2 border-b-2 border-[#8f8f8f] justify-between">
            <span class="header flex gap-2">
                <Icon v-if="note.icon.type=='icon'" class=" text-xl" :class="`text-[${note.icon.color}]`" :icon="note.icon.icon" />
                <p v-else>{{ note.icon.emoji }}</p>
                <a :href="note.from.link"><i>{{ note.from.text }}</i></a>
            </span>
            <div class="tags grid grid-rows-2 grid-cols-2 gap-1">
                <span 
                    v-for="i in 4"
                    class=" w-2 h-2 rounded-full"
                    :style="i-1 < note.tags.length?`background-color:#${note.tags[i-1].color};`:'background-color:#8f8f8f'"
                >
                </span>
            </div>
       </div>
       <div class="flex gap-2 p-2 text-xs">
           <img :src="note.cover" alt="Note cover" class=" aspect-square w-28 h-28 object-cover rounded-md">
           <div class="content h-28 w-full overflow-hidden">
     
                <div v-for="(block, i) of content" :key="i">
                    <p v-if="block.type=='h1'" class=" text-sm font-bold">
                        <span v-for="(text, j) of block.rich_text" :key="j" 
                        :style="
                            (text.text_style.color ? `color:${text.text_style.color};`:'')+
                            (text.text_style.width ? `font-weight:${text.text_style.width};`:'')+
                            (text.text_style.italic ? `font-style:italic;`:'')+
                            (text.text_style.underline ? `text-decoration:underline;`:'')// Add cross-out
                        ">
                            {{ text.plain_text }}
                        </span>
                    </p>
                    <p v-if="block.type=='paragraph'" class="">
                        <span v-for="(text, j) of block.rich_text" :key="j" 
                        :style="
                            (text.text_style.color ? `color:${text.text_style.color};`:'')+
                            (text.text_style.width ? `font-weight:${text.text_style.width};`:'')+
                            (text.text_style.italic ? `font-style:italic;`:'')+
                            (text.text_style.underline ? `text-decoration:underline;`:'')// Add cross-out
                        ">
                            {{ text.plain_text }}
                        </span>
                    </p>
                    <p v-if="block.type=='checkbox'" class="flex items-center">
                        <input type="checkbox" :checked="block.checked" disabled class=" mr-1">
                        <span v-for="(text, j) of block.rich_text" :key="j" 
                        :style="
                            (text.text_style.color ? `color:${text.text_style.color};`:'')+
                            (text.text_style.width ? `font-weight:${text.text_style.width};`:'')+
                            (text.text_style.italic ? `font-style:italic;`:'')+
                            (text.text_style.underline ? `text-decoration:underline;`:'')// Add cross-out
                        ">
                            {{ text.plain_text }}
                        </span>
                    </p>
                    <img v-if="block.type=='img'" :src="block.src" alt="" class=" w-full h-min object-cover">
                </div>
           </div>
       </div>
    </article>
</template>


<style></style>