/**
 * This file is mainly used as an entry point
 * to import components or define globals.
 * 
 */

// Waiting for top-level await to be better supported.
(async () => {
    import('./components/app.js');
    import("./components/organisms/search-bar.js");
    import("./components/organisms/results.js");
  }
)();
