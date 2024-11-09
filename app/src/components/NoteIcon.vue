<script lang="ts" setup>

defineProps<{
    icon: string
}>()

</script>

<template>
    <img class="icon" :src="icon" alt="">
</template>


<style scoped>

.icon{
    width: 1.5em;
    height: 1.5em;
    display: inline;
}

</style>