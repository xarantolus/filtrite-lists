<script lang="ts">
import { defineComponent, type PropType } from 'vue'
import type { FilterListData } from '@/model/FilterList';

export default defineComponent({
    props: {
        list1: {
            type: Object as PropType<FilterListData>,
            required: true
        },
        list2: {
            type: Object as PropType<FilterListData>,
            required: true
        }
    },
    components: {
    },
    data() {
        let listsIntersect = this.list1.urls.filter(u1 => this.list2.urls.some(u2 => u1.url === u2.url));

        return {
            same: listsIntersect,
            list1only: this.list1.urls.filter(u => !listsIntersect.some(ui => u.url === ui.url)),
            list2only: this.list2.urls.filter(u => !listsIntersect.some(ui => u.url === ui.url))
        }
    },
    computed: {
        sameName(): boolean {
            return this.list1.display_name == this.list2.display_name;
        }
    },
    methods: {
    }
})
</script>

<template>
    <div class="card filter-box">
        <div class="card-content">
            <details v-if="same.length > 0" class="content has-text-left">
                <summary>{{ same.length }} {{ same.length == 1 ? 'list is' : 'lists are' }} in both filters</summary>
                <ul>
                    <li v-bind:key="item.url" v-for="item in same">
                        <a target="_blank" :href="item.url">{{ item.title }}</a>
                    </li>
                </ul>
            </details>
            <p v-else class="has-text-left">There is no overlap between the two lists.</p>
            <details v-if="list1only.length > 0" class="content has-text-left">
                <summary>{{ list1only.length }} {{ list1only.length == 1 ? 'list is' : 'lists are' }} only in {{ list1.display_name }}{{sameName ? ' (' + list1.repo_owner + ')' : '' }}</summary>
                <ul>
                    <li v-bind:key="item.url" v-for="item in list1only">
                        <a target="_blank" :href="item.url">{{ item.title }}</a>
                    </li>
                </ul>
            </details>
            <details v-if="list2only.length > 0" class="content has-text-left">
                <summary>{{ list2only.length }} {{ list2only.length == 1 ? 'list is' : 'lists are' }} only in {{ list2.display_name }}{{sameName ? ' (' + list2.repo_owner + ')' : '' }}</summary>
                <ul>
                    <li v-bind:key="item.url" v-for="item in list2only">
                        <a target="_blank" :href="item.url">{{ item.title }}</a>
                    </li>
                </ul>
            </details>
        </div>
    </div>
</template>

<style>
</style>
