<script setup>
import { ref } from 'vue'
import axios from 'axios'
import { Icon } from '@iconify/vue'
import {useRouter} from 'vue-router'
import {useStore} from 'vuex'
import { refreshTokens } from '../utils/helpers'

const router = useRouter()
const store = useStore()

const slide = ref(0)

const username = ref("")
const password = ref("")
const repeatPassword = ref("")
const code = ref(["", "", "", ""])

const checkUsername = async () => {
    try {
        let { data } = await axios.get(`${import.meta.env.VITE_API_URL}/users/check-username?username=${username.value}`)
        if (data.found) {
            slide.value = 3
        } else {
            slide.value = 1
        }
    } catch (error) {
        alert(error)
    }
}

const sendCodeInTelegram = () => {
    // TODO: validate username
    let request = {
        username: username.value,
        type_id: 1
    }
    const jsonString = JSON.stringify(request)

    const base64Encoded = btoa(jsonString)

    window.open(`https://t.me/keep_it_authorization_bot?start=${base64Encoded}`, "_blank").focus()
}

const signUp = async () => {
    try {
        let { data } = await axios.post(`${import.meta.env.VITE_API_URL}/users/signup`, {
            username: username.value,
            password: password.value,
            code: code.value.join("")
        })
        store.commit("setAccess", data.authorization)
        console.log(data.authorization)
        router.push("/home")
    } catch (error) {
        console.log(error)
        alert(error)
    }
}

const logIn = async () => {
    try {
        let { data } = await axios.post(`${import.meta.env.VITE_API_URL}/users/login`, {
            username: username.value,
            password: password.value,
        }, {withCredentials:true})
        store.commit("setAccess", data.authorization)
        console.log(data.authorization)
        router.push("/home")
    } catch (error) {
        console.log(error)
        alert(error)
    }
}

const verifyCode = async () => {
    try {
        let { data } = await axios.post(`${import.meta.env.VITE_API_URL}/users/check-code`, {
            username: username.value,
            code: code.value.join("")
        })
        if (data.valid) {
            slide.value = 2
        } else {
            alert("Wrong code!")

        }
    } catch (error) {
        console.log(error)
        alert(error)
    }
}

const checkCode = async (event, i) => {
    code.value[i] = code.value[i].toUpperCase()
    if (code.value[i].length != 1 || !code.value[i].match(/[A-Z0-9]/)) {
        return
    }
    if (i < code.value.length - 1) {
        event.target.nextElementSibling.focus()
    }
    for (let j = 0; j < code.value.length; j++) {
        if (code.value[j] == '') {
            return
        }
    }

}


</script>

<template>
    <section class="flex justify-center items-center text-center">
        <article v-if="slide == 0" class="flex  w-72 flex-col items-center gap-5">
            <h1 class="text-xl font-bold">Hi ✌️ Who are you?</h1>
            <input v-model="username" placeholder="username" type="text"
                class=" w-full p-2 rounded-full bg-[var(--semi-dark)]">
            <button @click="checkUsername" class=" w-full p-2 rounded-full bg-white text-black">Continue</button>
        </article>
        <article v-if="slide == 3" class="flex  w-72 flex-col items-center gap-5">
            <h1 class="text-xl font-bold">Welcome back 👋</h1>
            <input v-model="password" placeholder="password" type="password"
                class=" w-full p-2 rounded-full bg-[var(--semi-dark)]">
            <button @click="logIn" class=" w-full p-2 rounded-full bg-white text-black">Continue</button>
        </article>
        <article v-if="slide == 2" class="flex  w-72 flex-col items-center gap-5">
            <h1 class="text-xl font-bold">Let's protect your account with password🕵</h1>
            <input v-model="password" placeholder="password" type="password"
                class=" w-full p-2 rounded-full bg-[var(--semi-dark)]">
            <input v-model="repeatPassword" placeholder="repeat password" type="password"
                class=" w-full p-2 rounded-full bg-[var(--semi-dark)]">
            <button @click="signUp()" class=" w-full p-2 rounded-full bg-white text-black">Keep it⚡</button>
        </article>
        <article v-if="slide == 1" class="flex  w-72 flex-col items-center gap-5">
            <h1 class="text-xl font-bold">We need to verify <br /> your account 🔍</h1>
            <button @click="sendCodeInTelegram(); slide = 4" class=" flex flex-col items-center gap-2">
                <Icon icon="mingcute:telegram-fill" class=" text-[6rem]" />
                <p>use telegram</p>
            </button>
        </article>
        <article v-if="slide == 4" class="flex  w-72 flex-col items-center gap-5">
            <h1 class="text-xl font-bold">Let's check your code✅</h1>
            <div class="flex gap-2">
                <input v-for="(cell, i) of code" :key="i" v-model="code[i]" @input="checkCode($event, i)" type="text"
                    class=" text-[3rem] p-2 rounded-xl w-full border-2 border-white bg-[var(--dark)] text-center">
            </div>
            <button @click="verifyCode" class=" w-full p-2 rounded-full bg-white text-black">Continue</button>
        </article>
    </section>
</template>


<style></style>