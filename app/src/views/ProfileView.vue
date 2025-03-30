<script lang="ts" setup>

import VKLinkButton from '@/components/login/VKLinkButton.vue';
import { UserTransport } from '@/transport/user';
import { Select } from 'primevue';
import { onBeforeMount, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';

const selectedLang = ref()

const langs = ref([
    {name: 'English', code: 'en'},
    {name: 'Russian', code: 'ru'},
])

onBeforeMount(()=>{
    const lang = langs.value.find(lang => lang.code === localStorage.getItem('lang'))
    selectedLang.value = lang
})

const i18n = useI18n()

watch(selectedLang, (newVal) => {
    localStorage.setItem('lang', newVal.code)
    i18n.locale.value = newVal.code
})

const linkVk = async (vkIDToken: string)=>{
    try {
        await UserTransport.linkVK(vkIDToken)
    } catch (error) {
        console.log(error)
    }
}

</script>

<template>
    <section class="home">
        <h1 class="text-4xl font-bold text-secondary-bg">Profile</h1>
        
        <section>
            <Select :options="langs" option-label="name" v-model="selectedLang" />
        </section>

        <section class="linking">
            <VKLinkButton @auth="linkVk"/>
        </section>

    </section>
</template>


<style>

.p-select{
    width: 100% !important;
}

</style>


<style scoped>

.home{
    @apply flex flex-col min-h-screen w-full p-4
}

.home section{
    max-width: 20rem;
}

.home-body{
    @apply grid grid-cols-4 w-full
}

.notes{
    @apply relative flex flex-col gap-2
}

</style>