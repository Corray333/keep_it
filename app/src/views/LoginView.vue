<script lang="ts" setup>

import { onBeforeMount, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { loadLocaleMessages } from '@/i18n'
import GoogleIcon from '@/components/icons/google-icon.vue';
import EmailIcon from '@/components/icons/email-icon.vue';
import TelegramIcon from '@/components/icons/telegram-icon.vue';
import { UserService } from '@/service/user';
import { useComponentsStore } from '@/stores/components';
import { useAccountStore } from '@/stores/account';
import type { User } from '@/entities/user';
import { useRouter } from 'vue-router';
import QRCodeStyling from 'qr-code-styling';
import { baseURL } from '@/helpers/data';


const router = useRouter()

enum Page {
    USERNAME,
    REGISTER,
    LOGIN,
    GET_CODE,
    METHOD_TG,
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

const passwordRequirements = ref([
    {
        regex: /^(?=.*[a-z]).+$/, // At least one lowercase letter
        translationKey: 'passwordReqs.lowercase',
        passedPassword: false,
        passedRepeatPassword: false
    },
    {
        regex: /^(?=.*[A-Z]).+$/, // At least one uppercase letter
        translationKey: 'passwordReqs.uppercase',
        passedPassword: false,
        passedRepeatPassword: false
    },
    {
        regex: /^(?=.*\d).+$/, // At least one digit
        translationKey: 'passwordReqs.digit',
        passedPassword: false,
        passedRepeatPassword: false
    },
    {
        regex: /^.{8,}$/, // Minimum 8 characters
        translationKey: 'passwordReqs.minLength',
        passedPassword: false,
        passedRepeatPassword: false
    },
    {
        regex: /^.{0,32}$/, // Maximum 32 characters
        translationKey: 'passwordReqs.maxLength',
        passedPassword: false,
        passedRepeatPassword: false
    }
])

const checkPasswords = () => {
    for (let i = 0; i < passwordRequirements.value.length; i++) {
        passwordRequirements.value[i].passedPassword = password.value.match(passwordRequirements.value[i].regex) != null
        passwordRequirements.value[i].passedRepeatPassword = repeatPassword.value.match(passwordRequirements.value[i].regex) != null
    }
}

watch(password, checkPasswords)
watch(repeatPassword, checkPasswords)

const tgLink = ref<string>('')
const tgQrCanvas = ref<HTMLCanvasElement>()

const generateTgCode = () => {
    const request = {
        username: username.value,
        type: userExists.value ? CodeType.LOGIN : CodeType.SIGNUP,
        syn: Math.floor(Math.random() * 1000000)
    }
    const jsonString = JSON.stringify(request)

    const base64Encoded = btoa(jsonString)

    tgLink.value = `https://t.me/keep_it_authorization_bot?start=${base64Encoded}`

    try {
        const qrCode = new QRCodeStyling({
            width: 300,
            height: 300,
            type: "svg",
            data: `https://t.me/keep_it_authorization_bot?start=${base64Encoded}`,
            image: `${baseURL}/static/images/identity/keep-it-logo-color.png`,
            dotsOptions: {
                color: "var(--secondary-bg)",
                type: "rounded"
            },
            backgroundOptions: {
                color: "var(--invert-bg)",
            },
            imageOptions: {
                crossOrigin: "anonymous",
                margin: 5
            }
        })

        qrCode.append(tgQrCanvas.value);
    } catch (error) {
        console.log(error)
        componentsStore.newError("Couldn't generate QR code")
    }

    const interval = setInterval(async ()=>{
        const exists = await userService.codeExists(username.value, request.syn)
        if (exists) {
            clearInterval(interval)
            currentPage.value = Page.ENTER_CODE
        }
    }, 1000)

    currentPage.value = Page.METHOD_TG
}

const checkCode = async (event: Event, i: number) => {
    code.value[i] = code.value[i].toUpperCase()
    if (code.value[i].length != 1 || !code.value[i].match(/[A-Z0-9]/)) {
        return
    }

    if (i < code.value.length - 1) {
        ((event?.target as HTMLElement)?.nextElementSibling as HTMLElement)?.focus()
    } else {
        login()
        return
    }

    for (let j = 0; j < code.value.length; j++) {
        if (code.value[j] == '') {
            return
        }
    }
}

const pasteCode = (event: ClipboardEvent) => {
    event.preventDefault()

    const pastedData = event.clipboardData?.getData('text') || ''

    code.value = pastedData.substring(0, 4).split('')

    login()
}

const login = async () => {
    let authData: {
        authorization: string
        user: User
    } | null
    if (userExists.value) {
        authData = await userService.logIn(username.value, password.value, code.value.join(''))
    } else {
        authData = await userService.signUp(username.value, password.value, code.value.join(''))
    }

    if (!authData) {
        return
    }

    const accountStore = useAccountStore()
    accountStore.setAuthorization(authData.authorization)
    accountStore.setUser(authData.user)

    router.push('/home')
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

            <div class="login-page signup-password">
                <div class="content">
                    <h3>{{ t('registerPage.title') }}</h3>
                    <div class="input-with-reqs group">
                        <input v-model="password" type="password" :placeholder="t('registerPage.password')">
                        <div class="requirements">
                            <p class="requirement" v-for="(req, i) of passwordRequirements" :key="i" :class="{'correct':req.passedPassword}">{{ t(req.translationKey) }}</p>
                        </div>
                    </div>
                    <div class="input-with-reqs group">
                        <input v-model="repeatPassword" type="password" :placeholder="t('registerPage.repeatPassword')">
                        <div class="requirements">
                            <p class="requirement" v-for="(req, i) of passwordRequirements" :key="i" :class="{'correct':req.passedRepeatPassword}">{{ t(req.translationKey) }}</p>
                        </div>
                    </div>
                    <button @click="currentPage = Page.GET_CODE" class="btn-type-1">{{ t('continueBtn') }}</button>
                </div>
            </div>

            <div class="login-page login-password">
                <div class="content">
                    <h3 class=" whitespace-pre-line">{{ t('loginPage.title') }}</h3>
                    <div class="input-with-reqs group">
                        <input v-model="password" type="password" :placeholder="t('registerPage.password')">
                        <div class="requirements">
                            <p class="requirement" v-for="(req, i) of passwordRequirements" :key="i" :class="{'correct':req.passedPassword}">{{ t(req.translationKey) }}</p>
                        </div>
                    </div>
                    
                    <button @click="currentPage = Page.GET_CODE" class="btn-type-1">{{ t('continueBtn') }}</button>
                </div>
            </div>

            <div class="login-page choose-verification-method">
                <div class="content">
                    <h3 class="whitespace-pre-line">{{ t('verify.title') }}</h3>

                    <div class="verify-row">
                        <div class="verify-email text-active">
                            <EmailIcon class="text-4xl" />
                            <p>use email</p>
                        </div>

                        <div @click="generateTgCode" class="verify-telegram text-active">
                            <TelegramIcon class="text-4xl" />
                            <p>use telegram</p>
                        </div>
                    </div>
                </div>
            </div>

            <div class="login-page telegram">
                <div class="content">
                    <p class="font-mono text-xl font-bold">{{ t('verifyTelegram.scanQR') }}</p>
                    <div class="rounded-2xl overflow-hidden p-2 bg-white" ref="tgQrCanvas"></div>
                    <a :href="tgLink" target="_blank">{{ t('verifyTelegram.useLink') }}</a>
                </div>
            </div>

            <div class="login-page">
                <div class="content">
                    <h3 class=" whitespace-pre-line">{{ t('verifyCode.title') }}</h3>

                    <div class="code">
                        <input v-for="(char, i) of code" :key="i" v-model="code[i]" type="text" maxlength="1"
                            @input="checkCode($event, i)" @paste="pasteCode">
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

.requirement{
    text-align: left;
}
.requirement.correct{
    color: var(--correct)
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

.code>input:focus{
    animation: input-blink 1s infinite;
    caret-color: transparent;
}

@keyframes input-blink {
    0% {
        background-color: var(--primary-bg);
    }
    50% {
        background-color: var(--secondary-bg);
    }
    100% {
        background-color: var(--primary-bg);
    }
    
}

.input-with-reqs {
    @apply relative flex w-full;
}

.input-with-reqs>input {
    @apply w-full
}

.requirements {
    @apply absolute top-0 ml-4 scale-0 w-full left-full bg-secondary-bg shadow-xl rounded-2xl p-4 group-focus-within:scale-100 duration-300;
}
</style>