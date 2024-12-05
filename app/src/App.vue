<script lang="ts" setup>
import { onBeforeMount } from 'vue';
import MessageStack from './components/MessageStack.vue'
import NavBar from './components/NavBar.vue';
import { useI18n } from 'vue-i18n';

const i18n = useI18n()

onBeforeMount(()=>{
  const lang = localStorage.getItem('lang') || 'en'
  console.log("lang", lang)
  i18n.locale.value = lang
})

</script>

<template>
    <section class="main">
        <MessageStack class="fixed top-0 right-0 z-50"/>
        <NavBar/>
        <RouterView/>
    </section>
</template>


<style scoped>

.main {
  @apply bg-primary-bg text-invert-bg min-h-screen md:flex p-2
}

</style>