<template>
    <div class="container-fluid py-5">
        <div class="row justify-content-center">
            <div class="col-12">
                <div class="mb-3 d-flex align-items-center gap-3">
                    <ElementsButtonBack/>
                    <h2 class="mb-0">{{ $t('panel.user_wanted_filters.edit_filter') }}</h2>
                </div>
                <div v-if="error" class="alert alert-danger">{{ error }}</div>
                <div v-else-if="filter" class="card shadow-sm">
                    <div class="card-body">
                        <PanelUserWantedFiltersForm :filter="filter" />
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import type { UserWantedPagesFilter } from '~/types/user_wanted_pages_filter'

definePageMeta({
    layout: 'panel',
    middleware: ['auth'],
})

const route = useRoute()
const userWantedFilters = useUserWantedFilters()
const filter = ref<UserWantedPagesFilter | null>(null)
const error = ref<string | null>(null)

const { filter: fetchedFilter, error: fetchError } = await userWantedFilters.fetchFilter(Number(route.params.id))
filter.value = fetchedFilter
error.value = fetchError
</script>
