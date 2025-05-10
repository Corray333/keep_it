chrome.runtime.onMessage.addListener((message, sender, sendResponse) => {
  if (message.type === "FETCH_CONTENT") {
    sendResponse({ data: document.body.innerHTML });
  }

  if (message.type === "FETCH_SELECTION") {
    const selection = window.getSelection();
    const container = document.createElement("div");
    for (let i = 0; i < selection.rangeCount; i++) {
      console.log(selection.getRangeAt(i).cloneContents());
      container.appendChild(selection.getRangeAt(i).cloneContents());
    }
    sendResponse({ data: container.innerHTML || "No selection" });
  }
});
