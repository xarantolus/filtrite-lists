<script lang="ts">
import { defineComponent, type PropType } from 'vue'
import type { FilterListData } from '../model/FilterList';

export default defineComponent({
    props: {
        list: {
            type: Object as PropType<FilterListData>,
            required: true
        }
    },
    data() {
        return {}
    },
})
</script>

<template>
    <p>{{ list.display_name }} - {{ list.stars }}</p>
</template>
