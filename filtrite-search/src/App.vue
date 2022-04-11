<script lang="ts">
import { defineComponent } from "vue";
import ListSearch from "./ListSearch.vue";
import type { FilterListData } from "./model/FilterList";
import { jsonp } from 'vue-jsonp'

const filtriteListURL = "https://github.com/xarantolus/filtrite-lists/releases/latest/download/filterlists_jsonp.js";

class JSONPResponse {
  public date : Date;
  public lists: Array<FilterListData>;
}

export default defineComponent({
  data() {
    return {
      filter_lists: [] as Array<FilterListData>,
      last_update_date: new Date(),
      loading: true,
      error: ""
    }
  },
  async mounted() {
    try {
      let resp = await jsonp<JSONPResponse>(filtriteListURL, {
        callbackQuery: 'cb',
        callbackName: 'jsonp',
      }, 15000);

      this.filter_lists = resp.lists;
      this.last_update_date = resp.date;
    } catch (e: any) {
      this.error = JSON.stringify(e);
    } finally {
      this.loading = false;
    }
  },
  components: { ListSearch }
});
</script>

<template>
  <main>
    <p v-if="loading">Loading data...</p>
    <p v-else-if="error">Error loading data:<br>{{ error }}</p>
    <ListSearch v-else :filter_lists="filter_lists" :update_date="last_update_date" ></ListSearch>
  </main>
  <footer class="ghlink">
    This page is <a target="_blank" href="https://github.com/xarantolus/filtrite-lists">open-source</a>, please feel free to report any issues.</footer>
</template>

<style>@import "../node_modules/bulma/css/bulma.min.css";
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
    --border-color: #999;
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

  .invert-dm {
    filter: invert()
  }
}

html,
body {
  height: 100vh;
}

footer.ghlink {
  bottom: 0;
  padding-top: 2.5%;
  text-align: center;
}

#app {
  max-width: 1000px;
  margin: 0 auto;
  padding: 2.5%;
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: var(--font-color);
  margin-top: 5%;
}
</style>
