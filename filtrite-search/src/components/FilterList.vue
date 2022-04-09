<script lang="ts">
import { defineComponent, type PropType } from 'vue'
import type { FilterListData } from '../model/FilterList';

import copy from 'copy-to-clipboard';

export default defineComponent({
    props: {
        list: {
            type: Object as PropType<FilterListData>,
            required: true
        }
    },
    data() {
        return {
            copied: false,
            error: "",
        }
    },
    methods: {
        copyURL() {
            try {
                copy(this.list.filter_file_url);
            } catch (e) {
                this.error = String(e);
            }

            this.copied = true;
            setTimeout(() => {
                this.copied = false;
                this.error = "";
            }, 750);
        }
    }
})
</script>

<template>
    <li>
        <div class="card filter-box">
            <div class="card-content">
                <h4 class="title is-4">{{ list.display_name }}</h4>
                <h5 class="subtitle is-5">
                    <a
                        target="_blank"
                        :href="'https://github.com/' + list.repo_owner"
                    >@{{ list.repo_owner }}</a>
                </h5>
                <details class="content has-text-left">
                    <summary>{{ list.urls.length }} included list{{ list.urls.length == 1 ? '' : 's' }}</summary>
                    <ul>
                        <li v-bind:key="item.url" v-for="item in list.urls">
                            <a :href="item.url">{{ item.title }}</a>
                        </li>
                    </ul>
                </details>
            </div>
            <footer class="card-footer">
                <a
                    @click.prevent="copyURL"
                    class="card-footer-item copy"
                    :href="list.filter_file_url"
                >{{ error ? 'Error!' : (copied ? 'Copied!' : 'Copy filter URL') }}</a>
                <a
                    target="_blank"
                    :href="'https://github.com/' + list.repo_owner + '/' + list.repo_name"
                    class="card-footer-item github"
                >View on GitHub</a>
            </footer>
        </div>
    </li>
</template>

<style scoped>
.filter-box {
 background:   #222;
    border: 2px solid grey;
    margin-bottom: 2%;
}

summary {
    width: 100%;
    cursor: pointer;
    padding: 2.5%;
}

.copy {
    background-color: green;
    color: #fff;
}
.github {
    background-color: rgb(63, 63, 63);
    color: #ddd;
}
</style>
