<template>
    <form @submit.prevent="onSubmit">
        <div v-if="error" class="alert alert-danger text-center mb-3">{{ error }}</div>
        <div class="row g-3">
            <div class="col-12">
                <div class="mb-3">
                    <label for="filter_name" class="form-label">{{ $t('panel.user_wanted_filters.name') }}</label>
                    <input
                        id="filter_name"
                        v-model="localCriteria.name"
                        type="text"
                        name="name"
                        class="form-control"
                        required
                        :placeholder="$t('panel.user_wanted_filters.name_placeholder')"
                    >
                </div>
            </div>
            <div class="col-md-6">
                <div class="mb-3">
                    <label for="sites" class="form-label">{{ $t('panel.user_wanted_filters.sites') }}</label>
                    <select id="sites" v-model="localCriteria.siteNames" name="sites[]" multiple class="form-select">
                        <option v-for="site in availableSites" :key="site" :value="site">{{ site }}</option>
                    </select>
                    <div class="form-text">{{ $t('panel.user_wanted_filters.sites_hint') }}</div>
                </div>
            </div>
            <div class="col-md-6">
                <div class="mb-3">
                    <label for="keywords" class="form-label">{{ $t('panel.user_wanted_filters.keywords') }}</label>
                    <input
                        id="keywords"
                        v-model="keywordsInput"
                        type="text"
                        name="keywords"
                        class="form-control"
                        :placeholder="$t('panel.user_wanted_filters.keywords_placeholder')"
                    >
                    <div class="form-text">{{ $t('panel.user_wanted_filters.keywords_hint') }}</div>
                </div>
            </div>
            <div class="col-12">
                <div class="mb-3">
                    <label for="title_contains" class="form-label">{{ $t('panel.user_wanted_filters.title_contains') }}</label>
                    <input
                        id="title_contains"
                        v-model="localCriteria.titleContains"
                        type="text"
                        name="title_contains"
                        class="form-control"
                        :placeholder="$t('panel.user_wanted_filters.title_contains_placeholder')"
                    >
                </div>
            </div>
        </div>
        <button type="submit" class="btn btn-primary w-100 mt-3">{{ $t('panel.user_wanted_filters.save') }}</button>
    </form>
</template>

<script setup lang="ts">
import type { UserWantedPageCriteria, UserWantedPagesFilter } from '~/types/user_wanted_pages_filter'

const props = defineProps<{ filter: UserWantedPagesFilter | null }>()

const availableSites = ['wikipedia.pl', 'example']

const localCriteria = ref<UserWantedPageCriteria>(props.filter?.filterData ? { ...props.filter.filterData } : {
    name: '',
    siteNames: [],
    keywords: [],
    titleContains: '',
})

const keywordsInput = ref((props.filter?.filterData?.keywords ?? []).join(', '))
const error = ref<string | null>(null)
const userWantedFilters = useUserWantedFilters()
const router = useRouter()

const onSubmit = async () => {
    const filterCriterias: UserWantedPageCriteria = {
        name: localCriteria.value.name.trim(),
        siteNames: localCriteria.value.siteNames?.length ? [...localCriteria.value.siteNames] : [],
        keywords: keywordsInput.value
            .split(',')
            .map(k => k.trim())
            .filter(Boolean),
        titleContains: localCriteria.value.titleContains?.trim() ?? '',
    }

    if (!filterCriterias.name) {
        error.value = 'panel.user_wanted_filters.name_required'
        return
    }

    let upsertError: string | null = null
    if (props.filter?.id) {
        const { error: updateError } = await userWantedFilters.updateFilter(props.filter.id, filterCriterias)
        upsertError = updateError
    } else {
        const { error: createError } = await userWantedFilters.createFilter(filterCriterias)
        upsertError = createError
    }

    error.value = upsertError
    if (!upsertError) {
        await router.push('/panel/wanted_filters')
    }
}
</script>
