<template>
    <div class="container-fluid py-5">
        <div class="row justify-content-center">
            <div class="col-12">
                <div class="mb-3 d-flex align-items-center gap-3">
                    <ElementsButtonBack/>
                    <h2 class="mb-0">{{ $t('panel.user_wanted_filters.add_new_filter') }}</h2>
                </div>
                <div class="card shadow-sm">
                    <div class="card-body">
                        <PanelUserWantedFiltersForm :filter="null" />
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
definePageMeta({
    layout: 'panel',
    middleware: ['auth'],
})
</script>
