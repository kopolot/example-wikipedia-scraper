<template>
    <div class="container-fluid py-5">
        <div class="row justify-content-center">
            <div class="col-12">
                <div class="card shadow-sm">
                    <div class="card-body">
                        <h1 class="text-center mb-4">{{ $t('panel.page_records.title') }}</h1>

                        <div class="d-flex flex-wrap justify-content-between align-items-center mb-3 gap-2">
                            <div class="d-flex align-items-center gap-2 flex-wrap">
                                <label class="form-label mb-0 small">{{ $t('panel.page_records.sort') }}:</label>
                                <select v-model="sort" class="form-select form-select-sm" style="width:auto">
                                    <option value="created_at">{{ $t('panel.page_records.sort_created_at') }}</option>
                                    <option value="updated_at">{{ $t('panel.page_records.sort_updated_at') }}</option>
                                    <option value="title">{{ $t('panel.page_records.sort_title') }}</option>
                                </select>
                                <select v-model="order" class="form-select form-select-sm" style="width:auto">
                                    <option value="desc">{{ $t('panel.page_records.order_desc') }}</option>
                                    <option value="asc">{{ $t('panel.page_records.order_asc') }}</option>
                                </select>
                            </div>
                            <div class="d-flex align-items-center gap-2 flex-wrap">
                                <ElementsSectionsPagination v-model:page="page" v-model:limit="limit" :max-page="maxPage" />
                                <button class="btn btn-outline-secondary btn-sm" @click="resetFilters">
                                    {{ $t('panel.page_records.reset_filters') }}
                                </button>
                                <button class="btn btn-primary btn-sm" :disabled="loading" @click="loadPageRecords">
                                    <BootstrapIcon name="arrow-clockwise"/>
                                    {{ $t('panel.page_records.refresh') }}
                                </button>
                            </div>
                        </div>

                        <div v-if="loading" class="text-center text-primary my-4">{{ $t('panel.page_records.loading') }}</div>
                        <div v-if="error" class="alert alert-danger text-center mx-auto" style="max-width:500px;">{{ error }}</div>
                        <PanelPageRecordsList v-if="!loading && !error" :page-records="pageRecords" />
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import type { PageType } from '~/types/page';

definePageMeta({
    layout: 'panel',
    middleware: ['auth']
})

const { fetchPageRecords } = usePages();

const pageRecords = ref<PageType[]>([]);
const loading = ref(false);
const error = ref<string | null>(null);
const limit = ref(20);
const page = ref(1);
const maxPage = ref(1);

const sort = ref('created_at');
const order = ref('desc');



async function loadPageRecords() {
    loading.value = true;
    error.value = null;
    const result = await fetchPageRecords(limit.value, page.value, sort.value, order.value);
    pageRecords.value = result.pages;
    maxPage.value = result.maxPage;
    error.value = result.error;
    loading.value = false;
}

function resetFilters() {
    page.value = 1;
}

watch(limit, (newLimit, oldLimit) => {
    const firstRecord = oldLimit * (page.value - 1) + 1;
    page.value = Math.ceil(firstRecord / newLimit);
});

watch([limit, page, sort, order], () => {
    loadPageRecords();
});

onMounted(async () => {
    loadPageRecords();
});
</script>

<style scoped>
</style>
