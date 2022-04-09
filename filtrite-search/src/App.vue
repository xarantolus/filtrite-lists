<script setup lang="ts">
import ListSearch from "./ListSearch.vue";

import filter_lists from './filterlists.json';
</script>

<template>
  <header></header>

  <main>
    <ListSearch :filter_lists="filter_lists"></ListSearch>
  </main>
</template>

<style>
@import "../node_modules/bulma/css/bulma.min.css";
@import "../node_modules/bulma-prefers-dark/css/bulma-prefers-dark.css";

:root {
  --font-color: #1f2933;
  --border-color: #ccc;
  --table-hover-color: #f0e8ff;
  --green: #00ff55;
  --red: #ff94a1;
  --blue: #6fcaff;
  --hidden-text: #d1d1d1;
}
@media (prefers-color-scheme: dark) {
  :root {
    --font-color: #eed7c0;
    --border-color: #333;
    --table-hover-color: #1c0f35;
    --green: #158c11;
    --red: #860000;
    --blue: #004770;
    --hidden-text: #616161;
  }
  /* Some additional color fixes */
  .help {
    color: #aaa;
  }
}
html,
body {
  height: 100vh;
}
#app {
  max-width: 1000px;
  margin: 0 auto;
  padding: 1%;
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: var(--font-color);
  margin-top: 5%;
}
</style>
