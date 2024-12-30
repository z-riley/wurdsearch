import define from "../utils/define.js";

const template = () => /*html*/ `
  <header class="search-menu">
    <div class="branding">
      <img class="brand-icon" src="/images/logo.png" width="40" height="60" alt="logo">
      <span class="brand-title">Wurdsearch</span>
    </div>
    <wurdsearch-search-bar></wurdsearch-search-bar>
  </header>
    <main>
      <wurdsearch-results></wurdsearch-results>
    </main>
`;

export default define(
  "app",
  class extends HTMLElement {
    constructor() {
      super();
      this.__setup();
    }

    __setup() {
      this.innerHTML = template();
    }
  }
);
