<template>
    <div class="container-fluid py-5">
        <div class="row justify-content-center">
            <div class="col-12">
                <div class="filter-header d-flex justify-content-between align-items-center flex-wrap mb-4 gap-2">
                    <ElementsButtonBack/>
                    <h2 class="mb-0">{{ filter?.filterData?.name || $t('panel.user_wanted_filters.details') }}</h2>
                    <div class="filter-actions d-flex gap-2">
                        <button class="btn btn-primary" @click="onEdit">
                            <i class="bi bi-pencil me-1"/>{{ $t('common.edit') }}
                        </button>
                        <button class="btn btn-danger" @click="onDelete">
                            <i class="bi bi-trash me-1"/>{{ $t('common.delete') }}
                        </button>
                    </div>
                    <h4 v-if="filter?.deletedAt" class="w-100 text-danger">{{ $t('panel.user_wanted_filters.disabled') }}</h4>
                </div>
                <div class="card shadow-sm filter-details-card mb-4">
                    <div class="card-body p-0">
                        <table class="table mb-0">
                            <tbody>
                                <tr>
                                    <th>{{ $t('panel.user_wanted_filters.name') }}</th>
                                    <td>{{ criteria?.name }}</td>
                                </tr>
                                <tr v-if="criteria?.siteNames?.length">
                                    <th>{{ $t('panel.user_wanted_filters.sites') }}</th>
                                    <td>{{ criteria.siteNames.join(', ') }}</td>
                                </tr>
                                <tr v-if="criteria?.keywords?.length">
                                    <th>{{ $t('panel.user_wanted_filters.keywords') }}</th>
                                    <td>{{ criteria.keywords.join(', ') }}</td>
                                </tr>
                                <tr v-if="criteria?.titleContains">
                                    <th>{{ $t('panel.user_wanted_filters.title_contains') }}</th>
                                    <td>{{ criteria.titleContains }}</td>
                                </tr>
                            </tbody>
                        </table>
                    </div>
                </div>
                <h3 class="mt-5">{{ $t('panel.user_wanted_filters.matching_pages') }}</h3>
                <PanelPageRecordsList :page-records="pages" :error="error" />
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import type { UserWantedPagesFilter, UserWantedPageCriteria } from '~/types/user_wanted_pages_filter'
import type { PageType } from '~/types/page'

definePageMeta({
    layout: 'panel',
    middleware: ['auth'],
})

const route = useRoute()
const router = useRouter()
const userWantedFilters = useUserWantedFilters()

const filter = ref<UserWantedPagesFilter | null>(null)
const criteria = ref<UserWantedPageCriteria | null>(null)
const pages = ref<PageType[]>([])
const error = ref<string | null>(null)

const fetchFilterAndPages = async () => {
    const id = Number(route.params.id)
    const { filter: fetchedFilter, error: fetchError } = await userWantedFilters.fetchFilter(id)
    error.value = fetchError
    filter.value = fetchedFilter || null
    criteria.value = fetchedFilter ? fetchedFilter.filterData : null

    if (fetchedFilter?.id) {
        const { pages: fetchedPages, error: pagesError } = await userWantedFilters.fetchPagesByFilters([fetchedFilter.id])
        if (pagesError) {
            error.value = pagesError
        } else {
            pages.value = fetchedPages
        }
    }
}

const onEdit = () => {
    router.push(`/panel/wanted_filters/edit/${route.params.id}`)
}

const onDelete = async () => {
    if (!confirm('Czy na pewno chcesz usunąć ten filtr?')) {
        return
    }
    const { success, error: deleteError } = await userWantedFilters.deleteFilter(Number(route.params.id))
    if (!success) {
        alert('Nie udało się usunąć filtra: ' + deleteError)
        return
    }
    await router.push('/panel/wanted_filters')
}

onMounted(fetchFilterAndPages)
</script>
