import define from '../../utils/define.js';

const template = () => /*html*/`
  <h1>
    Search something to get started
  </h1>
`;

export default define('home', class extends HTMLLIElement {
  constructor() {
    super();
    this.classList.add('home');
    this.__setup();
  }

  __setup() {
    this.innerHTML = template();
  }
}, { extends: 'li' });