<script setup lang="ts">
import { ref } from 'vue';


const pageContent = ref('');

function grabContent() {
  chrome.runtime.sendMessage(
    { type: 'GET_PAGE_CONTENT' } as { type: string; data?: string },
    (response: { data?: string }) => {
      pageContent.value = response?.data || 'No content found';
    }
  );
}

</script>

<template>
  <main>
    <h1>You are logged in</h1>
    <button @click="grabContent">Grab content</button>
    <p>{{ pageContent }}</p>
  </main>
</template>