<template>
    <div v-if="error" class="alert alert-danger text-center mt-4">{{ error }}</div>
    <div v-else-if="pageRecords.length" class="row g-4 page-records-list">
        <div v-for="pageRecord in pageRecords" :key="pageRecord.id" class="col-12 col-md-6 col-lg-4">
            <NuxtLink :to="`/panel/page-records/${pageRecord.id}`" class="card shadow-sm h-100 text-decoration-none page-record-card page-record-link">
                <div class="card-body">
                    <div class="page-record-header d-flex justify-content-between align-items-center mb-2">
                        <strong>{{ pageRecord.title }}</strong>
                        <span class="page-record-date">{{ pageRecord.createdAt.toLocaleString() }}</span>
                    </div>
                    <div class="page-record-details d-flex flex-column gap-1">
                        <span>{{ $t('page_record.title') }}: {{ pageRecord.title }}</span>
                        <span>{{ $t('page_record.textField1') }}: {{ pageRecord.textField1?.substring(0, 100) }}</span>
                        <span>{{ $t('page_record.textField2') }}: {{ pageRecord.textField2?.substring(0, 100) }}</span>
                        <span>{{ $t('page_record.textField3') }}: {{ pageRecord.textField3?.substring(0, 100) }}</span>
                    </div>
                </div>
            </NuxtLink>
        </div>
    </div>
    <div v-else class="alert alert-info text-center mt-4">{{ $t('page_records.no_records') }}</div>
</template>

<script setup lang="ts">
import type { PageType } from '~/types/page';

defineProps<{
    pageRecords: PageType[];
    error?: string | null;
}>();
</script>

<style scoped>
.page-record-card {
  border-left: 5px solid #1976d2;
  border-radius: 8px;
  transition: box-shadow 0.2s, transform 0.2s;
  color: inherit;
  cursor: pointer;
  &:hover {
    box-shadow: 0 4px 16px rgba(25,118,210,0.15);
    transform: translateY(-2px) scale(1.01);
    background: #f5faff;
    text-decoration: none;
  }
}
</style>