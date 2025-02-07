const faviconLink = document.querySelector('link[rel="icon"]');

window
  .matchMedia("(prefers-color-scheme: dark)")
  .addEventListener("change", (event) => {
    if (event.matches) {
      // replace with light logo
      faviconLink.href = "https://flowbite.com/docs/images/logo.svg";
    } else {
      // replace with dark logo
      faviconLink.href = "https://flowbite.com/docs/images/logo.svg";
    }
  });

if (
  window.matchMedia &&
  window.matchMedia("(prefers-color-scheme: dark)").matches
) {
  // replace with light logo
  faviconLink.href = "https://flowbite.com/docs/images/logo.svg";
} else {
  // replace with dark logo
  faviconLink.href = "https://flowbite.com/docs/images/logo.svg";
}
