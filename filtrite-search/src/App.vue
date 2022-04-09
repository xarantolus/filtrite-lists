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
  --green: #00ff55;
  --yellow: #f7ae40;
  --blue: #6fcaff;
  --card-color: #fff;
  --button-color: #bbb;
}
@media (prefers-color-scheme: dark) {
  :root {
    --font-color: #eed7c0;
    --border-color: #333;
    --green: #158c11;
    --yellow: #df8b1d;
    --blue: #004770;
    --card-color: #222;
    --button-color: #333;
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
