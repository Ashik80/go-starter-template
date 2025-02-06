(() => {
  let passwordInput = document.querySelector("#password");
  let openEye = document.querySelector("#open-eye");
  let closedEye = document.querySelector("#closed-eye");

  function setupEventListeners() {
    reEvent(closedEye, "click", () => {
      passwordInput.type = "text";
      closedEye.classList.add("hidden");
      openEye.classList.remove("hidden");
    });

    reEvent(openEye, "click", () => {
      passwordInput.type = "password";
      closedEye.classList.remove("hidden");
      openEye.classList.add("hidden");
    });
  }

  reEvent(document, "htmx:afterSwap", () => {
    passwordInput = document.querySelector("#password");
    openEye = document.querySelector("#open-eye");
    closedEye = document.querySelector("#closed-eye");
    setupEventListeners();
  });

  setupEventListeners();
})();
