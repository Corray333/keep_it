<script setup>
import { Icon } from '@iconify/vue'
import { ref } from 'vue'

const props = defineProps(['note'])


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
           <img :src="note.cover" alt="Note cover" class=" w-28 h-28 object-cover rounded-md">
           <div class="content h-28 overflow-hidden">
                <div v-for="(block, i) of note.content" :key="i">
                    <p v-if="block.type=='h1'" class="flex text-sm font-bold">
                        <p v-for="(text, j) of block.rich_text" :key="j" 
                        :style="
                            (text.text_style.color ? `color:#${text.text_style.color};`:'')+
                            (text.text_style.width ? `font-weight:${text.text_style.width};`:'')+
                            (text.text_style.italic ? `font-style:italic;`:'')+
                            (text.text_style.underline ? `text-decoration:underline;`:'')// Add cross-out
                        ">
                            {{ text.plain_text }}
                        </p>
                    </p>
                    <span v-if="block.type=='checkbox'" class="flex gap-1 items-center">
                        <input type="checkbox" :checked="block.checked" disabled>
                        <p v-for="(text, j) of block.rich_text" :key="j" 
                        :style="
                            (text.text_style.color ? `color:#${text.text_style.color};`:'')+
                            (text.text_style.width ? `font-weight:${text.text_style.width};`:'')+
                            (text.text_style.italic ? `font-style:italic;`:'')+
                            (text.text_style.underline ? `text-decoration:underline;`:'')// Add cross-out
                        ">
                            {{ text.plain_text }}
                        </p>
                    </span>
                    <img v-if="block.type=='img'" :src="block.src" alt="" class=" w-full h-min object-cover">
                </div>
           </div>
       </div>
    </article>
</template>


<style></style>