chrome.runtime.onMessage.addListener((message, sender, sendResponse) => {
    if (message.type === "FETCH_CONTENT") {
      sendResponse({ data: document.body.innerHTML });
    }
  });
  