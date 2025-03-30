chrome.runtime.onMessage.addListener((message, sender, sendResponse) => {
  if (message.type === "GET_PAGE_CONTENT") {
    chrome.tabs.query({ active: true, currentWindow: true }, (tabs) => {
      if (tabs.length === 0 || !tabs[0].id) {
        sendResponse({ data: "No active tab" });
        return;
      }

      chrome.scripting.executeScript(
        {
          target: { tabId: tabs[0].id },
          files: ["content-script.js"]
        },
        () => {
          chrome.tabs.sendMessage(tabs[0].id, { type: "FETCH_CONTENT" }, sendResponse);
        }
      );
    });

    return true; // Указывает, что sendResponse будет вызван асинхронно
  }
});
