<script setup>
import { Icon } from '@iconify/vue'
import { useStore } from 'vuex'
import axios from 'axios'

const store = useStore()

const getUser = async ()=>{
    try {
        let { data } = await axios.get(`${import.meta.env.VITE_API_URL}/users/0`, {
            headers: {
                'Authorization': store.state.AccessToken,
            }
        })
        console.log(data)
    } catch (error) {
        console.log(error)
        alert(error)
    }
}

</script>

<template>
    <header class="flex items-center w-full p-5 pb-0 gap-5">
        <div class="search w-full relative flex items-center justify-end">
            <input type="text" placeholder="search" class="p-2 rounded-full bg-[#333] w-full">
            <button
                class=" absolute text-xl h-10 aspect-square flex justify-center items-center border-l-2 border-[#8f8f8f]">
                <Icon icon="mdi:filter-outline" />
            </button>
        </div>
        <button @click="getUser" class=" flex h-10 items-center justify-center rounded-full bg-[#333] aspect-square text-xl">
            <Icon icon="iconamoon:profile" />
        </button>
    </header>
</template>


<style></style>