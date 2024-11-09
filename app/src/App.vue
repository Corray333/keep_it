<script lang="ts" setup>
import MessageStack from './components/MessageStack.vue'


</script>

<template>
    <section class="main">
        <MessageStack class="fixed top-0 right-0 z-50"/>
        <RouterView/>
    </section>
</template>


<style scoped>

.main {
  @apply bg-primary-bg text-invert-bg min-h-screen
}

</style>