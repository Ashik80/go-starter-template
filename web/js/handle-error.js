document.addEventListener("htmx:beforeSwap", (e) => {
  if (e.detail.xhr.status === 400) {
    e.detail.isError = false;
    e.detail.shouldSwap = true;
  }
});
