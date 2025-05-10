<script lang="ts" setup>
import { IconType, type Icon } from '@/entities/note';


defineProps<{
  icon: Icon
}>()

</script>

<template>
  <img v-if="icon.type == IconType.IMG" class="icon" :src="icon.data" alt="">
  <div v-if="icon.type == IconType.SVG" v-html="icon.data" class="icon"></div>
</template>


<style>

.icon{
  @apply inline text-2xl;
}
.icon>svg {
  @apply text-invert-bg;
  width: 1em;
  height: 1em;
  display: inline;
}
</style>
