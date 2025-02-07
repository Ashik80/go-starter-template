const template = document.createElement("template");
template.innerHTML = `
  <svg id="open-eye" class="hidden absolute right-2 top-2 text-gray-900 dark:text-white cursor-pointer" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
    <path d="M1 12C1 12 5 3 12 3s11 9 11 9-4 9-11 9S1 12 1 12z"/>
    <circle cx="12" cy="12" r="3"/>
  </svg>
  <svg id="closed-eye" class="absolute right-2 top-2 text-gray-900 dark:text-white cursor-pointer" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
    <path d="M1 12C1 12 5 21 12 21s11-9 11-9-4-9-11-9S1 12 1 12z"/>
    <line x1="1" y1="12" x2="23" y2="12"/>
  </svg>
`;

class PasswordInputControl extends HTMLElement {
  constructor() {
    super();
    this.appendChild(template.content.cloneNode(true));
    this.initializeInput = this.initializeInput.bind(this);
    this.showPassword = this.showPassword.bind(this);
    this.hidePassword = this.hidePassword.bind(this);
    this.setupEventListeners = this.setupEventListeners.bind(this);
    this.cleanupEventListeners = this.cleanupEventListeners.bind(this);
  }

  initializeInput() {
    this.passwordInput = this.querySelector("#password");
    this.openEye = this.querySelector("#open-eye");
    this.closedEye = this.querySelector("#closed-eye");
    this.setupEventListeners();
  }

  showPassword() {
    this.passwordInput.type = "text";
    this.closedEye.classList.add("hidden");
    this.openEye.classList.remove("hidden");
  }

  hidePassword() {
    this.passwordInput.type = "password";
    this.closedEye.classList.remove("hidden");
    this.openEye.classList.add("hidden");
  }

  setupEventListeners() {
    this.closedEye.addEventListener("click", this.showPassword);
    this.openEye.addEventListener("click", this.hidePassword);
  }

  cleanupEventListeners() {
    this.closedEye.removeEventListener("click", this.showPassword);
    this.openEye.removeEventListener("click", this.hidePassword);
  }

  connectedCallback() {
    this.initializeInput();
    this.addEventListener("htmx:afterSwap", this.initializeInput);
  }

  disconnectedCallback() {
    this.cleanupEventListeners();
    this.removeEventListener("htmx:afterSwap", this.initializeInput);
  }
}

customElements.define("password-input-control", PasswordInputControl);
