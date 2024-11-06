<script lang="ts" setup>
import { Message } from '@/entities/message'


defineProps({
    message: {
        type: Message,
        required: true
    }
})

</script>

<template>
    <div class="message" :class="message.type">
        <p>{{ message.text }}</p>
    </div>
</template>


<style scoped>

.message {
    @apply py-4 px-8 rounded-2xl border-2 bg-primary-bg shadow-md w-full flex justify-center;
}

.message.error {
    @apply border-red-500 text-red-500;
}

.message.success {
    @apply border-green-500 text-green-500;
}

.message.warning {
    @apply border-yellow-500 text-yellow-500;
}

.message.info{
    @apply border-white text-white;
}

</style>