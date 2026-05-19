<template>
    <div class="container-fluid py-5">
        <div class="row justify-content-center">
            <div class="col-12">
                <div class="card shadow-sm">
                    <div class="card-body">
                        <h1 class="text-center mb-4">{{ $t('panel.title') }}</h1>
                        <!-- tutaj możesz dodać kolejne elementy panelu -->
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">

definePageMeta({
    layout: 'panel',
    middleware: ['auth']
})

</script>

<style scoped lang="scss">
/* Możesz dodać własne style, jeśli chcesz nadpisać domyślne style Bootstrap */
</style>