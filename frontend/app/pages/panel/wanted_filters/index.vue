<template>
    <div class="container-fluid py-5">
        <div class="row justify-content-center">
            <section class="watched-filters mb-5 col-12">
                <div class="filters-header d-flex align-items-center justify-content-between mb-3">
                    <h2 class="mb-0">{{ $t('panel.user_wanted_filters.title') }}</h2>
                    <NuxtLink class="btn btn-primary" to="/panel/wanted_filters/add">
                        <i class="bi bi-plus-lg me-1"/>{{ $t('panel.user_wanted_filters.add') }}
                    </NuxtLink>
                </div>
                <div v-if="filtersLoading" class="text-center text-primary my-3">
                    <div class="spinner-border me-2" role="status">
                        <span class="visually-hidden">{{ $t('panel.loading') }}</span>
                    </div>
                    {{ $t('panel.loading') }}
                </div>
                <div v-if="filtersError" class="alert alert-danger text-center mx-auto" style="max-width:500px;">
                    {{ filtersError }}
                </div>
                <PanelUserWantedFiltersList :filters="filters"/>
            </section>

            <section class="pages-list col-12">
                <h2>{{ $t('panel.user_wanted_filters_pages.title') }}</h2>
                <div v-if="pagesLoading" class="text-center text-primary my-3">
                    <div class="spinner-border me-2" role="status">
                        <span class="visually-hidden">{{ $t('panel.loading') }}</span>
                    </div>
                    {{ $t('panel.loading') }}
                </div>
                <PanelPageRecordsList :page-records="pages" :error="pagesError" />
            </section>
        </div>
    </div>
</template>

<script setup lang="ts">
import type { UserWantedPagesFilter } from '~/types/user_wanted_pages_filter'
import type { PageType } from '~/types/page'

definePageMeta({
    layout: 'panel',
    middleware: ['auth'],
})

const userWantedFilters = useUserWantedFilters()

const filters = ref<UserWantedPagesFilter[]>([])
const filtersLoading = ref(false)
const filtersError = ref<string | null>(null)

const pages = ref<PageType[]>([])
const pagesLoading = ref(false)
const pagesError = ref<string | null>(null)

const loadFilters = async () => {
    filtersLoading.value = true
    filtersError.value = null
    const { filters: currentFilters, error } = await userWantedFilters.fetchFilters()
    filters.value = currentFilters ?? []
    filtersError.value = error
    filtersLoading.value = false
}

const loadPages = async () => {
    pagesLoading.value = true
    pagesError.value = null
    const { pages: currentPages, error } = await userWantedFilters.fetchFilteredPages()
    pages.value = currentPages ?? []
    pagesError.value = error
    pagesLoading.value = false
}

onMounted(() => {
    loadFilters()
    loadPages()
})
</script>
