/**
 * This file is made for tweaking parameters on the front-end
 * without having to dive in the source code.
 * 
 */

export default {
  componentPrefix: 'turdsearch',
  publicApiURL: 'http://localhost:8080/',
  searchQueryParam: 'q',
  commands: {
    'go: ': 'https://',
    'search: google.com ': 'https://www.google.com/search?q=',
  }
}
