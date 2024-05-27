import define from '../../utils/define.js';
import escapeString from '../../utils/escapeString.js';
import { globalBus } from '../../utils/events.js';

const template = ({ data }) => /*html*/`
  <a href='${data.url}'>
    <p class='link'>${data.url}</p>
    <p class='title'>${data.title}</p>
    <p class='extract'>${data.extract}</p>
  </a>
`;

export default define('result', class extends HTMLLIElement {
  constructor() {
    super();
    this.classList.add('result');
    this.__setup();
  }

  __setup() {
    this.innerHTML = template({ data: {
      url: this.dataset.url,
      title: this.__handleBold(JSON.parse(this.dataset.title)),
      extract: this.__handleBold(JSON.parse(this.dataset.extract))
     }});
    this.__events();
  }

  __events() {
    this.addEventListener('keydown', (e) => {
      if (this.firstElementChild === document.activeElement) {
        if (e.key === 'ArrowDown') {
          e.preventDefault();
          this?.nextElementSibling?.firstElementChild.focus();
        }
        if (e.key === 'ArrowUp') {
          e.preventDefault();
          if (this.previousElementSibling)
            this.previousElementSibling?.firstElementChild.focus();
          else {
            const focusSearchEvent = new CustomEvent('focus-search');
            globalBus.dispatch(focusSearchEvent);
          }
        }
      }
    })
  }

  __handleBold(input) {
    let text = '';
    for (const part of input) {
      if (part.is_bold) text += `<strong>${escapeString(part.value)}</strong>`;
      else text += escapeString(part.value);
    }
    return text;
  }
}, { extends: 'li' });