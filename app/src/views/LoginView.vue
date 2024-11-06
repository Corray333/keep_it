<script lang="ts" setup>

import { onBeforeMount, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { loadLocaleMessages } from '@/i18n'
import GoogleIcon from '@/components/icons/google-icon.vue';
import EmailIcon from '@/components/icons/email-icon.vue';
import TelegramIcon from '@/components/icons/telegram-icon.vue';
import { UserService } from '@/service/user';
import { useComponentsStore } from '@/stores/components';

enum Page {
    USERNAME,
    REGISTER,
    LOGIN,
    GET_CODE,
    ENTER_CODE
}

enum CodeType {
    SIGNUP = 1,
    LOGIN,
}

const currentPage = ref<Page>(Page.USERNAME)

const username = ref<string>('')
const email = ref<string>('')
const password = ref<string>('')
const repeatPassword = ref<string>('')
const code = ref<string[]>(['', '', '', ''])

const userExists = ref<boolean>(false)

const { t, locale } = useI18n()

onBeforeMount(() => {
    loadLocaleMessages(locale.value, 'login')
})

const userService = new UserService()

const checkUsername = async () => {
    const user = await userService.findUserLogin(username.value)
    if (user != null) {
        username.value = user.username
        email.value = user.email
        currentPage.value = Page.LOGIN
        userExists.value = true
    } else {
        currentPage.value = Page.REGISTER
        userExists.value = false
    }
}

const sendCodeInTelegram = () => {
    const request = {
        username: username.value,
        type: userExists.value ? CodeType.LOGIN : CodeType.SIGNUP
    }
    const jsonString = JSON.stringify(request)

    const base64Encoded = btoa(jsonString)

    window.open(`https://t.me/keep_it_authorization_bot?start=${base64Encoded}`, "_blank")?.focus()
    currentPage.value = Page.ENTER_CODE
}

const checkCode = async (event: Event, i: number) => {
    code.value[i] = code.value[i].toUpperCase()
    if (code.value[i].length != 1 || !code.value[i].match(/[A-Z0-9]/)) {
        return
    }

    if (i < code.value.length - 1) {
        ((event?.target as HTMLElement)?.nextElementSibling as HTMLElement)?.focus()
    } else {
        if (userExists.value){
            userService.logIn(username.value, password.value, code.value.join(''))
        } else {
            userService.signUp(username.value, password.value, code.value.join(''))
        }
        return
    }

    for (let j = 0; j < code.value.length; j++) {
        if (code.value[j] == '') {
            return
        }
    }


}

const componentsStore = useComponentsStore()


</script>

<template>
    <section class="login">
        <div class="pages-carousel duration-300" :style="`transform: translateX(${-100 * currentPage}vw)`">

            <div class="login-page">
                <div class="content">
                    <h3>{{ t('usernamePage.title') }}</h3>
                    <input v-model="username" type="text" :placeholder="t('usernamePage.usernameOrEmail')">
                    <button @click="checkUsername" class="btn-type-1">{{ t('continueBtn') }}</button>
                    <span class="flex items-center gap-2">
                        <p>{{ t('usernamePage.logInUsing') }}</p>
                        <GoogleIcon @click="componentsStore.newError()" class=" text-xl text-active" />
                    </span>
                </div>
            </div>

            <div class="login-page">
                <div class="content">
                    <h3 v-html="t('registerPage.title')"></h3>
                    <input v-model="password" type="password" :placeholder="t('registerPage.password')">
                    <input v-model="repeatPassword" type="password" :placeholder="t('registerPage.repeatPassword')">
                    <button class="btn-type-1">{{ t('continueBtn') }}</button>
                </div>
            </div>

            <div class="login-page">
                <div class="content">
                    <h3 v-html="t('loginPage.title')"></h3>
                    <input v-model="password" type="password" :placeholder="t('loginPage.password')">
                    <button @click="currentPage = Page.GET_CODE" class="btn-type-1">{{ t('continueBtn') }}</button>
                </div>
            </div>

            <div class="login-page">
                <div class="content">
                    <h3 v-html="t('verify.title')"></h3>

                    <div class="verify-row">
                        <div class="verify-email text-active">
                            <EmailIcon class="text-4xl" />
                            <p>use email</p>
                        </div>

                        <div @click="sendCodeInTelegram" class="verify-telegram text-active">
                            <TelegramIcon class="text-4xl" />
                            <p>use telegram</p>
                        </div>
                    </div>
                </div>
            </div>

            <div class="login-page">
                <div class="content">
                    <h3 v-html="t('verifyCode.title')"></h3>

                    <div class="code">
                        <input v-for="(char, i) of code" :key="i" v-model="code[i]" type="text" maxlength="1" @input="checkCode($event, i)">
                    </div>
                </div>
            </div>

        </div>
    </section>
</template>


<style scoped>
.slide-enter-active,
.slide-leave-active {
    transition: transform 0.5s ease;
}

.slide-enter,
.slide-leave-to {
    transform: translateX(100%);
}

.slide-leave-active {
    position: absolute;
}

.login {
    @apply h-screen w-screen relative overflow-hidden text-center;
}

.pages-carousel {
    @apply w-fit h-full absolute flex bg-primary-bg top-0
}

.login-page {
    @apply w-screen h-full flex items-center justify-center
}

.content {
    @apply w-fit min-w-64 flex flex-col gap-2 items-center
}

.content>input {
    @apply w-full
}

h3 {
    @apply text-2xl text-accent;
}

.verify-row {
    @apply flex w-full items-center justify-between mt-8;
}

.verify-row>div {
    @apply flex flex-col gap-2 items-center
}

.code {
    @apply flex gap-2;
}

.code>input {
    @apply w-16 text-4xl font-mono p-4 rounded-md flex items-center justify-center text-center font-bold;
}
</style>