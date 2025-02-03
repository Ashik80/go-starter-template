(() => {
  if (customElements.get("todo-detail")) {
    return;
  }

  class TodoDetail extends HTMLElement {
    constructor() {
      super();
      this.details = this.querySelector("#details");
      this.editForm = this.querySelector("#edit-form");
      this.editButton = this.querySelector("#edit-button");
      this.cancelButton = this.querySelector("#cancel-button");
      this.deleteButton = this.querySelector("#delete-button");

      this.hideDetailsAndShowForm = this.hideDetailsAndShowForm.bind(this);
      this.hideFormAndShowDetails = this.hideFormAndShowDetails.bind(this);
      this.setupEventListeners = this.setupEventListeners.bind(this);
      this.removeEventListeners = this.removeEventListeners.bind(this);
      this.reinitializeForm = this.reinitializeForm.bind(this);
      this.shouldSwapAfterDelete = this.shouldSwapAfterDelete.bind(this);
    }

    hideDetailsAndShowForm() {
      this.details.style.display = "none";
      this.editForm.style.display = "block";
    }

    hideFormAndShowDetails() {
      this.details.style.display = "block";
      this.editForm.style.display = "none";
    }

    connectedCallback() {
      this.hideFormAndShowDetails();
      this.setupEventListeners();
      this.editForm.addEventListener("htmx:afterSwap", this.reinitializeForm);
    }

    disconnectedCallback() {
      this.removeEventListeners();
    }

    reinitializeForm() {
      this.editForm = this.querySelector("#edit-form");
      this.cancelButton = this.querySelector("#cancel-button");
      this.setupEventListeners();
    }

    setupEventListeners() {
      this.editForm.addEventListener("submit", this.hideFormAndShowDetails);
      this.editButton.addEventListener("click", this.hideDetailsAndShowForm);
      this.cancelButton.addEventListener("click", this.hideFormAndShowDetails);
      document.addEventListener("htmx:beforeSwap", this.shouldSwapAfterDelete);
    }

    shouldSwapAfterDelete(e) {
      if (e.detail.xhr.status === 404) {
        e.detail.shouldSwap = true;
      }
    }

    removeEventListeners() {
      this.editForm.removeEventListener("submit", this.hideFormAndShowDetails);
      this.editButton.removeEventListener("click", this.hideDetailsAndShowForm);
      this.cancelButton.removeEventListener(
        "click",
        this.hideFormAndShowDetails,
      );
      this.editForm.removeEventListener(
        "htmx:afterSwap",
        this.reinitializeForm,
      );
      document.removeEventListener(
        "htmx:beforeSwap",
        this.shouldSwapAfterDelete,
      );
    }
  }

  customElements.define("todo-detail", TodoDetail);
})();
