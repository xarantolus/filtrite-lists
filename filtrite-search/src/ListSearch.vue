<script lang="ts">
import { defineComponent, type PropType } from 'vue'
import type { FilterListData } from './model/FilterList';
import FuzzySearch from 'fuzzy-search';
import FilterList from './components/FilterList.vue';

export default defineComponent({
    props: {
        filter_lists: {
            type: Array as PropType<Array<FilterListData>>,
            required: true,
        },
    },
    data() {
        return {
            query: "",
            searcher: new FuzzySearch(this.filter_lists ?? [], ["display_name", "urls.title", "repo_owner", "repo_name"], {
                sort: true,
            }),
        };
    },
    computed: {
        searchResults(): Array<FilterListData> {
            return this.searcher.search(this.query);
        }
    },
    components: { FilterList }
})
</script>

<template>
    <div class="search-field">
        <h5 class="title is-5">Bromite filter search</h5>
        <div class="field">
            <div class="control">
                <input class="input" autofocus type="text" v-model="query" />

                <ul>
                    <FilterList
                        v-for="item in searchResults"
                        v-bind:key="item.filter_file_url"
                        :list="item"
                    ></FilterList>
                </ul>
            </div>
        </div>
    </div>
</template>
