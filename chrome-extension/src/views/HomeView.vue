<script setup lang="ts">
import { UserTransport } from '@/services/userService';
import { useAccountStore } from '@/stores/account';
import { ref } from 'vue';

const pageContent = ref('');
const accStore = useAccountStore();

function grabContent() {
  if (!chrome?.tabs || !chrome?.runtime) {
    console.warn('Browser does not support chrome extension APIs.');
    pageContent.value = 'Extension APIs not available.';
    return;
  }

  chrome.tabs.query({ active: true, currentWindow: true }, (tabs) => {
    const currentTab = tabs[0];
    const currentUrl = currentTab?.url || 'Unknown URL';

    chrome.runtime.sendMessage(
      { type: 'GET_PAGE_CONTENT' },
      (response: { data?: string } = {}) => {
        pageContent.value = response?.data || 'No content found';
        if (response?.data) {
          UserTransport.processPage(currentUrl, pageContent.value);
        }
      }
    );
  });
}

const hasSelection = ref(false);

function grabSelection() {
  console.log('grabSelection');
  chrome.tabs.query({ active: true, currentWindow: true }, (tabs) => {
    const currentTab = tabs[0];
    const currentUrl = currentTab?.url || 'Unknown URL';

    chrome.runtime.sendMessage(
      { type: 'GET_SELECTED_CONTENT' },
      (response: { data?: string } = {}) => {
        pageContent.value = response?.data || 'No selection found';
        console.log(response?.data);
        if (response?.data) {
          UserTransport.processPage(currentUrl, pageContent.value);
        }
      }
    );
  });
}

// Слушаем события обновления выделения
window.addEventListener('focus', checkSelection);
window.addEventListener('mouseup', checkSelection);
window.addEventListener('keyup', checkSelection);

function checkSelection() {
  chrome.tabs.query({ active: true, currentWindow: true }, (tabs) => {
    const tabId = tabs[0]?.id;
    if (!tabId) return;

    chrome.scripting?.executeScript?.(
      {
        target: { tabId },
        func: () => window.getSelection()?.toString().trim().length > 0,
      },
      (results) => {
        hasSelection.value = results?.[0]?.result || false;
      }
    );
  });
}


</script>

<template>
  <main>
    <h2>Keep it</h2>
    <button class="button-1" @click="grabContent">Grab content</button>
    <button class="button-1" @click="grabSelection">Grab selection</button>
  </main>
</template>



<style scoped>

main{
  @apply p-4 flex flex-col items-center gap-4 bg-primary-bg text-invert-bg;
}
main>h2{
  @apply text-center text-2xl font-bold;
}

</style>
