<template>
    <div v-if="filters?.length" class="row g-4">
        <div v-for="filter in filters" :key="filter.id" class="col-12 col-md-6 col-lg-4">
            <NuxtLink
                :to="`/panel/wanted_filters/${filter.id}`"
                class="card shadow-sm h-100 text-decoration-none filter-card"
                :class="filter.deletedAt ? 'bg-dark-subtle' : ''"
            >
                <div class="card-body filter-details">
                    <h5 class="card-title mb-3">{{ filter.filterData.name }}</h5>
                    <div v-if="filter.filterData.siteNames?.length">
                        <strong>{{ $t('panel.user_wanted_filters.sites') }}:</strong>
                        {{ filter.filterData.siteNames.join(', ') }}
                    </div>
                    <div v-if="filter.filterData.keywords?.length">
                        <strong>{{ $t('panel.user_wanted_filters.keywords') }}:</strong>
                        {{ filter.filterData.keywords.join(', ') }}
                    </div>
                    <div v-if="filter.filterData.titleContains">
                        <strong>{{ $t('panel.user_wanted_filters.title_contains') }}:</strong>
                        {{ filter.filterData.titleContains }}
                    </div>
                </div>
            </NuxtLink>
        </div>
    </div>
    <div v-else class="alert alert-info text-center mt-3">{{ $t('panel.user_wanted_filters.no_filters') }}</div>
</template>

<script setup lang="ts">
import type { UserWantedPagesFilter } from '~/types/user_wanted_pages_filter'

defineProps<{
    filters: UserWantedPagesFilter[]
}>()
</script>

<style scoped lang="scss">
.filter-card {
    border-radius: 1rem;
    transition: box-shadow 0.2s, transform 0.2s;
    color: inherit;
    border: 1px solid #e5e7eb;
    cursor: pointer;
    &:hover {
        box-shadow: 0 4px 24px rgba(0, 0, 0, 0.14);
        transform: translateY(-2px) scale(1.02);
        border-color: #cbd5e1;
        background: linear-gradient(135deg, #f1f5f9 60%, #e2e8f0 100%);
        text-decoration: none;
    }
}
.filter-details {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
}
.filter-details strong {
    color: #2563eb;
    font-weight: 600;
    margin-right: 0.3em;
}
.filter-details div {
    font-size: 1rem;
    line-height: 1.5;
    color: #334155;
    background: rgba(226, 232, 240, 0.18);
    border-radius: 0.4em;
    padding: 0.2em 0.5em;
}
</style>
