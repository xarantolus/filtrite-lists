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
        const processedFilterlists = this.filter_lists.map(list => {
            return {
                ...list,
                search_data: [list.display_name, list.repo_name, list.repo_owner, ...(list.urls.map(u => [u.title, u.url]).flat())].join(" ").toLowerCase()
            };
        });

        return {
            query: "",
            searcher: new FuzzySearch(processedFilterlists, ["display_name", "repo_name", "repo_owner", "search_data"], {
                sort: true,
                caseSensitive: false,
            }),
        };
    },
    computed: {
        searchResults(): Array<FilterListData> {
            let qsplit = this.query.toLowerCase().split(/\s+/);
            return this.searcher.search(this.query).filter(i => qsplit.every(w => i.search_data.includes(w)));
        }
    },
    components: { FilterList },
})
</script>

<template>
    <div class="search-field">
        <h5 class="title is-5">Bromite filter search</h5>
        <div class="field">
            <div class="control">
                <input
                    class="input"
                    placeholder="Search filter lists..."
                    autofocus
                    type="search"
                    :value="query"
                    @input="event => query = (event?.target as HTMLInputElement).value"
                />

                <ul class="spacing">
                    <FilterList
                        v-for="item in searchResults"
                        v-bind:key="item.filter_file_url"
                        :list="item"
                        :search="query"
                    ></FilterList>
                </ul>
            </div>
        </div>
    </div>
</template>



<style scoped>
.spacing {
    margin-top: 2%;
}
</style>
