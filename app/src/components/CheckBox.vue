<script lang="ts" setup>
import { ref, defineProps, defineEmits, watch } from 'vue'

const props = defineProps<{ modelValue: boolean, disabled: boolean }>()
const emit = defineEmits(['update:modelValue'])

const done = ref(props.modelValue)

watch(() => props.modelValue, (newValue) => {
    done.value = newValue
})

const toggleDone = () => {
    if (props.disabled) return
    done.value = !done.value
    emit('update:modelValue', done.value)
}

</script>

<template>
    <div class="checkbox-wrapper">
        <div @click="toggleDone"
            class=" checkbox">
            <svg xmlns="http://www.w3.org/2000/svg" width="12" height="10" viewBox="0 0 12 10" fill="none">
                <path :style="done ? 'stroke-dashoffset: -40;' : 'stroke-dashoffset: -20;'"
                    d="M1 5.5L4.18849 8.28992C4.62644 8.67313 5.29783 8.6029 5.64699 8.13735L11 1" stroke="white"
                    stroke-width="2" stroke-linecap="round" />
            </svg>
        </div>
    </div>
</template>

<style scoped>

.checkbox{
    @apply  cursor-pointer flex justify-center items-center border-2 border-white rounded-md;
    width: 1.125em;
    height: 1.125em;
}

.checkbox-wrapper{
    @apply inline-block p-1;
}

path {
    stroke-dasharray: 20;
    transition: all 0.3s;
}

</style>