<template>
    <label>
        {{ $t('panel.record.items_per_page') }}:
        <select v-model="limit" class="form-select d-inline-block w-auto ms-2">
            <option v-for="n in [10, 20, 50, 100]" :key="n" :value="n">{{ n }}</option>
        </select>
    </label>
    <nav aria-label="Paginacja stron" class="ms-2">
        <ul class="pagination mb-0">
            <li class="page-item" :class="{ disabled: page === 1 }">
                <button :disabled="page === 1" class="page-link" @click="setPage(page - 1)" >&laquo;</button>
            </li>
            <li v-for="p in pagesToShow" :key="p" class="page-item" :class="{ active: page === p }">
                <button class="page-link" @click="setPage(p)">{{ p }}</button>
            </li>
            <li class="page-item" :class="{ disabled: page === maxPage }">
                <button :disabled="page === maxPage" class="page-link" @click="setPage(page + 1)" >&raquo;</button>
            </li>
        </ul>
    </nav>
</template>

<script setup lang="ts">
const props = defineProps({
    limit: { type: Number, default: 20 },
    page: { type: Number, default: 1 },
    maxPage: { type: Number, default: 1 },
});
const emit = defineEmits(['update:page', 'update:limit']);

const maxPage = computed(() => props.maxPage);
const page = computed({
    get: () => props.page,
    set: (val: number) => emit('update:page', val)
});
const limit = computed({
    get: () => props.limit,
    set: (val: number) => emit('update:limit', val)
});

function setPage(p: number) {
    if (p >= 1 && p <= maxPage.value) {
        page.value = p;
    }
}

const pagesToShow = computed(() => {
    const total = maxPage.value;
    const current = page.value;
    const delta = 2;
    const pages: number[] = [];
    for (let i = Math.max(1, current - delta); i <= Math.min(total, current + delta); i++) {
        pages.push(i);
    }
    return pages;
});
</script>