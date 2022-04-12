<script lang="ts">
import { defineComponent, type PropType } from 'vue'
import type { FilterListData } from '../model/FilterList';

import Highlighter from 'vue-highlight-words';

import copy from 'copy-to-clipboard';

export default defineComponent({
    props: {
        list: {
            type: Object as PropType<FilterListData>,
            required: true
        },
        search: {
            type: String,
        }
    },
    components: {
        Highlighter,
    },
    data() {
        return {
            copied: false,
            error: "",
        }
    },
    computed: {
        keywords(): Array<string> {
            return this.search?.toLowerCase()?.split(/\s+/).filter(v => v.trim().length > 0) ?? [];
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
        },
        urlMatches(url: string): boolean {
            if (this.keywords.length == 0) {
                return false;
            }
            return this.keywords.some(kw => url.toLowerCase().includes(kw))
        }
    }
})
</script>

<template>
    <li>
        <div class="card filter-box">
            <div class="card-content">
                <span v-if="list?.stars ?? 0 > 0" class="card-header-icon stars">⭐{{ list.stars }}</span>
                <h4>
                    <Highlighter highlightClassName="highlight" class="title is-4" :searchWords="keywords" :autoEscape="true" :textToHighlight="list.display_name" />
                </h4>
                <details class="content has-text-left">
                    <summary>{{ list.urls.length }} included list{{ list.urls.length == 1 ? '' : 's' }}</summary>
                    <ul>
                        <li v-bind:key="item.url" v-for="item in list.urls">
                            <a target="_blank" :href="item.url">
                                <Highlighter highlightClassName="highlight" :searchWords="keywords" :autoEscape="true" :textToHighlight="item.title" />
                                <span class="tag urlmatch" v-if="urlMatches(item.url)">URL match</span>
                            </a>
                        </li>
                    </ul>
                </details>
            </div>
            <footer class="card-footer">
                <a target="_blank" :href="'https://github.com/' + list.repo_owner + '/' + list.repo_name" class="card-footer-item github">
                    <span class="icon" style="padding-right: 2.5%;">
                        <img class="invert-dm" src="@/assets/GitHub-dark.png">
                    </span>
                    <Highlighter highlightClassName="highlight" :searchWords="keywords" :autoEscape="true" :textToHighlight="list.repo_owner + '/' + list.repo_name" />
                </a>
                <a @click.prevent="copyURL" class="card-footer-item copy" target="_blank" :href="list.filter_file_url">{{ error ? 'Error!' : (copied ? 'Copied!' : 'Copy filter URL') }}</a>
            </footer>
        </div>
    </li>
</template>

<style>
.filter-box {
    background: var(--card-color);
    border: 3px solid var(--border-color);
    margin-bottom: 2%;
}

.stars {
    position: absolute;
    top: 0;
    right: 0;
    user-select: none;
}

summary {
    margin-top: 0.5rem !important;
    width: 100%;
    cursor: pointer;
    padding: 0.5%;
}

.copy {
    background-color: var(--green);
    color: var(--font-color);
}

.github {
    background-color: var(--button-color);
    color: var(--font-color);
    border: none;
}

.highlight {
    background-color: var(--yellow);
}

.tag.urlmatch {
    margin-left: 1%;
    color: var(--font-color-on-yellow);
    background-color: var(--yellow) !important;
}
</style>
