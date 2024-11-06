<script lang="ts" setup>
import { useComponentsStore } from '@/stores/components';
import MessageComponent from './MessageComponent.vue';


const componentsStore = useComponentsStore()

</script>

<template>
    <Transition name="delay">
        <section v-show="componentsStore.messages.length > 0" class="msg-stack">
            <TransitionGroup name="slide-down">
                <MessageComponent v-for="(msg, i) of componentsStore.messages" :key="i" :message="msg" />
            </TransitionGroup>
        </section>
    </Transition>
</template>


<style scoped>

.delay-enter-active,
.delay-leave-active {
    transition: opacity 0.25s ease;
}

.delay-enter-from,
.delay-leave-to {
    opacity: 0;
}

.slide-down-enter-active,
.slide-down-leave-active {
    transition: transform 0.25s ease;
}

.slide-down-enter-from,
.slide-down-leave-to {
    transform: translateY(-100%);
}

.msg-stack{
    @apply flex flex-col gap-4 mr-4 mt-4 w-full max-w-96;
}
.msg-stack>*{
    @apply absolute duration-300 top-0 right-0;
    transform: translateY(30px) scale(0.9);
}

.msg-stack>:nth-last-child(2) {
    transform: translateY(15px) scale(0.95);
}

.msg-stack>:nth-last-child(1) {
    transform: translateY(0px) scale(1);
}

</style>