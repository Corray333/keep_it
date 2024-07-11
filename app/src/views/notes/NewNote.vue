<script setup>
import { ref, computed } from 'vue'
import { Icon, addIcon } from '@iconify/vue'
import axios from 'axios'

// TODO: change that
addIcon('mdi:account', {
    body: '<path fill-rule="evenodd" clip-rule="evenodd" d="M6.66667 5.18519C5.84847 5.18519 5.18519 5.84847 5.18519 6.66667V17.3333C5.18519 18.1515 5.84847 18.8148 6.66667 18.8148H17.3333C18.1515 18.8148 18.8148 18.1515 18.8148 17.3333V6.66667C18.8148 5.84847 18.1515 5.18519 17.3333 5.18519H15.7037C15.3764 5.18519 15.1111 4.91987 15.1111 4.59259C15.1111 4.26531 15.3764 4 15.7037 4H17.3333C18.8061 4 20 5.19391 20 6.66667V17.3333C20 18.8061 18.8061 20 17.3333 20H6.66667C5.19391 20 4 18.8061 4 17.3333V6.66667C4 5.19391 5.19391 4 6.66667 4H8.2963C8.62358 4 8.88889 4.26531 8.88889 4.59259C8.88889 4.91987 8.62358 5.18519 8.2963 5.18519H6.66667Z" fill="currentColor"/><path fill-rule="evenodd" clip-rule="evenodd" d="M15.2593 7.85185C14.7683 7.85185 14.3704 8.24982 14.3704 8.74074C14.3704 9.23166 14.7683 9.62963 15.2593 9.62963C15.7502 9.62963 16.1481 9.23166 16.1481 8.74074C16.1481 8.24982 15.7502 7.85185 15.2593 7.85185ZM13.1852 8.74074C13.1852 7.59526 14.1138 6.66667 15.2593 6.66667C16.4047 6.66667 17.3333 7.59526 17.3333 8.74074C17.3333 9.88622 16.4047 10.8148 15.2593 10.8148C14.6947 10.8148 14.1828 10.5892 13.8088 10.2233L10.8012 11.7605C10.8102 11.8391 10.8148 11.919 10.8148 12C10.8148 12.1761 10.7929 12.3471 10.7516 12.5104L13.8088 14.073C14.1828 13.7071 14.6947 13.4815 15.2593 13.4815C16.4047 13.4815 17.3333 14.4101 17.3333 15.5556C17.3333 16.701 16.4047 17.6296 15.2593 17.6296C14.1138 17.6296 13.1852 16.701 13.1852 15.5556C13.1852 15.4025 13.2018 15.2534 13.2332 15.1098L10.1415 13.5296C9.77249 13.8677 9.28072 14.0741 8.74076 14.0741C7.59528 14.0741 6.66669 13.1455 6.66669 12C6.66669 10.8545 7.59528 9.92593 8.74076 9.92593C9.37981 9.92593 9.95136 10.2149 10.3318 10.6694L13.2332 9.18648C13.2018 9.0429 13.1852 8.89376 13.1852 8.74074ZM7.85187 12C7.85187 11.5091 8.24984 11.1111 8.74076 11.1111C9.23168 11.1111 9.62965 11.5091 9.62965 12C9.62965 12.4909 9.23168 12.8889 8.74076 12.8889C8.24984 12.8889 7.85187 12.4909 7.85187 12ZM14.3704 15.5556C14.3704 15.0646 14.7683 14.6667 15.2593 14.6667C15.7502 14.6667 16.1481 15.0646 16.1481 15.5556C16.1481 16.0465 15.7502 16.4444 15.2593 16.4444C14.7683 16.4444 14.3704 16.0465 14.3704 15.5556Z" fill="currentColor"/>',
    width: 24,
    height: 24,
});


const file = ref(null)


const newNote = ref({
    source: "keep_it",
    icon: {
        type: "icon",
        icon: "mdi:account",
        color: "#FFFFFF"
    },
    from: {
        "text": "Keep it"
    },
    tags: [

    ],
    title: "",
    cover: "",
    content: [
        {
            type: "p",
            rich_text: {
                plain_text: "Test text",
                meta: [
                    {
                        offset: 0,
                        length: 4,
                        link: "",
                        color: "#FFC75A",
                        width: "bold",
                        italic: false,
                        underline: false,
                        cross_out: true
                    }
                ]
            }
        }
    ]
})

const content = computed(() => {
    const result = []
    for (let note of newNote.value.content) {
        if (note.rich_text) {
            let temp = {
                type: note.type,
                rich_text: []
            }
            let i = 0
            for (let block of note.rich_text.meta) {
                if (i < block.offset) {
                    temp.rich_text.push({
                        plain_text: note.rich_text.plain_text.substr(i, block.offset),
                        text_style: {
                        }
                    })
                }
                temp.rich_text.push({
                    plain_text: note.rich_text.plain_text.substr(block.offset, block.offset + block.length),
                    text_style: {
                        color: block.color,
                        width: block.width,
                        italic: block.italic,
                        underline: block.underline,
                        cross_out: block.cross_out
                    }
                })
                i += block.length
            }
            if (i < note.rich_text.plain_text.length) {
                temp.rich_text.push({
                    plain_text: note.rich_text.plain_text.substr(i, note.rich_text.plain_text.length),
                    text_style: {
                    }
                })
            }
            result.push(temp)
        } else {
            result.push(note)
        }
    }
    return result
})

const handleCoverUpload = async (event) => {
    if (event.target.files[0].size > 5 * 1024 * 1024) {
        return
    }
    file.value = event.target.files[0]
    const reader = new FileReader()

    reader.onload = (e) => {
        posts.value[0].cover = e.target.result
    }
    reader.readAsDataURL(event.target.files[0])
    let url = await loadFile()
    for (let i = 0; i < posts.value.length; i++) {
        posts.value[i].cover = url
    }
}

const loadFile = async () => {
    try {
        let formData = new FormData()
        formData.append('file', file.value)

        let { data } = await axios.post(`${import.meta.env.VITE_API_URL}/upload/image`, formData, {
            headers: {
                "Authorization": localStorage.getItem('Authorization')
            }
        })
        return data.url
    } catch (err) {
        console.log(err)
    }
}

const handleImageUpload = async (event) => {
    if (event.target.files[0].size > 5 * 1024 * 1024) {
        return
    }
    file.value = event.target.files[0]
    const reader = new FileReader()

    reader.readAsDataURL(event.target.files[0])
    let file_name = await loadFile()
    console.log(`![image](${file_name})`)
    navigator.clipboard.writeText(`![image](${file_name})`)
}

const handleInput = (text, event) => {
    console.log(text, inputText)
    const inputText = event.target.innerText
    text.plain_text = inputText

}

const test = ()=>{
    console.log(window.getSelection())
}

</script>

<template>
    <section class=" text-white">
        <div class="tools fixed bottom-0 w-full p-2 bottom-bar-shadow">
            <ul class="flex w-full justify-around">
                <li>
                    <Icon icon="mdi:image-plus-outline" class="text-3xl"></Icon>
                </li>
                <li>
                    <Icon icon="mdi:checkbox-marked-outline" class="text-3xl"></Icon>
                </li>
                <li>
                    <Icon icon="mdi:text-box-edit-outline" class="text-3xl"></Icon>
                </li>
                <li>
                    <Icon icon="mdi:microphone" class="text-3xl"></Icon>
                </li>
                <li>
                    <Icon icon="mdi:dots-vertical" class="text-3xl"></Icon>
                </li>
            </ul>
        </div>
        <div class="profile_photo relative">
            <input @input="changed = true" type="file" id="coverInput" class="hidden" @change="handleCoverUpload" />
            <label for="coverInput"
                class="text-center absolute mx-auto bg-gray-900 bg-opacity-80 h-full w-full flex items-center justify-center opacity-0 duration-300 cursor-pointer  hover:opacity-100">
                <Icon icon="mdi:user" />
            </label>
            <img :src="file ? posts[0].cover : 'https://image-tc.galaxy.tf/wijpeg-1vskcxut1q5gyupy5uj05agr1/placeholder_standard.jpg?crop=38%2C0%2C1365%2C1024'"
                alt="" class="w-full h-32 object-cover border-white">
        </div>
        <div class="main p-5">
            <h1 class="flex">
                <Icon :icon="newNote.icon.icon" class="text-4xl" />
                <input v-model="newNote.title" class=" bg-transparent font-bold" type="text" placeholder="Title">
            </h1>
            <div class="content border-2 border-[#8f8f8f] min-h-10 text-sm p-2 rounded-xl" contenteditable>
                <!-- <span role="textbox" contenteditable class="relative flex items-center w-fit"
                    @keydown.enter="newTask(index)"
                    spellcheck="false"
                        @input="handleInput(task, $event)"
                        @focus="handleFocus"
                        @blur="handleBlur"
                        v-text="task.task"></span> -->

                <div v-for="(block, i) of content" :key="i" contenteditable @click="test()">
                    <span v-if="block.type == 'h1'" class=" text-sm font-bold">
                        <span role="textbox" v-for="(text, j) of block.rich_text" :key="j" :style="(text.text_style.color ? `color:#${text.text_style.color};` : '') +
                            (text.text_style.width ? `font-weight:${text.text_style.width};` : '') +
                            (text.text_style.italic ? `font-style:italic;` : '') +
                            (text.text_style.underline ? `text-decoration:underline;` : '')// Add cross-out
                            " @input="handleInput(task, $event)">
                                {{ text.plain_text }}
                        </span>
                    </span>
                    <span role="textbox" v-if="block.type == 'p'" class=" w-fit" v-for="(text, j) of block.rich_text" :key="j" @click="test()"
                        
                        :style="(text.text_style.color ? `color:#${text.text_style.color};` : '') +
                                (text.text_style.width ? `font-weight:${text.text_style.width};` : '') +
                                (text.text_style.italic ? `font-style:italic;` : '') +
                                (text.text_style.underline ? `text-decoration:underline;` : '')// Add cross-out
                                ">
                        {{ text.plain_text }}
                    </span>
                    <span v-if="block.type == 'checkbox'" class="flex gap-1 items-center">
                        <input type="checkbox" :checked="block.checked" disabled>
                        <p v-for="(text, j) of block.rich_text" :key="j" :style="(text.text_style.color ? `color:#${text.text_style.color};` : '') +
                (text.text_style.width ? `font-weight:${text.text_style.width};` : '') +
                (text.text_style.italic ? `font-style:italic;` : '') +
                (text.text_style.underline ? `text-decoration:underline;` : '')// Add cross-out
                ">
                            {{ text.plain_text }}
                        </p>
                    </span>
                    <img v-if="block.type == 'img'" :src="block.src" alt="" class=" w-full h-min object-cover">
                </div>
            </div>
        </div>
    </section>
</template>


<style></style>