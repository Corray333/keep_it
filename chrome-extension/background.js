chrome.runtime.onMessage.addListener((message, sender, sendResponse) => {
  if (message.type === "GET_PAGE_CONTENT") {
    chrome.tabs.query({ active: true, currentWindow: true }, (tabs) => {
      if (!tabs.length || !tabs[0].id) {
        sendResponse({ data: "No active tab" });
        return;
      }

      const tabId = tabs[0].id;

      // Вызов после инъекции скрипта
      const fetchContent = () => {
        chrome.tabs.sendMessage(tabId, { type: "FETCH_CONTENT" }, sendResponse);
      };

      // Fallback для Firefox (tabs.executeScript) и Chrome (scripting.executeScript)
      if (chrome.scripting && chrome.scripting.executeScript) {
        // Chrome
        chrome.scripting.executeScript(
          {
            target: { tabId },
            files: ["content-script.js"]
          },
          () => fetchContent()
        );
      } else if (chrome.tabs && chrome.tabs.executeScript) {
        // Firefox
        chrome.tabs.executeScript(tabId, { file: "content-script.js" }, () => fetchContent());
      } else {
        console.error("Script injection not supported in this browser.");
        sendResponse({ data: "Script injection not supported." });
      }
    });

    return true; // allow async sendResponse
  }
  if (message.type === "GET_SELECTED_CONTENT") {
    chrome.tabs.query({ active: true, currentWindow: true }, (tabs) => {
      if (!tabs.length || !tabs[0].id) {
        sendResponse({ data: "No active tab" });
        return;
      }

      const tabId = tabs[0].id;

      const fetchSelectedContent = () => {
        chrome.tabs.sendMessage(tabId, { type: "FETCH_SELECTION" }, sendResponse);
      };

      if (chrome.scripting && chrome.scripting.executeScript) {
        chrome.scripting.executeScript(
          {
            target: { tabId },
            files: ["content-script.js"]
          },
          () => fetchSelectedContent()
        );
      } else if (chrome.tabs && chrome.tabs.executeScript) {
        chrome.tabs.executeScript(tabId, { file: "content-script.js" }, () => fetchSelectedContent());
      } else {
        console.error("Script injection not supported in this browser.");
        sendResponse({ data: "Script injection not supported." });
      }
    });

    return true;
  }

});
