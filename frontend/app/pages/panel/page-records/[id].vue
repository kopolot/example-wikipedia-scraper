<template>
  <div class="container-fluid py-5">
    <div class="row justify-content-center">
      <div class="col-12">
        <div class="mb-3">
          <ElementsButtonBack />
        </div>
        <div v-if="loading" class="text-center text-primary my-4">
          <div class="spinner-border me-2" role="status"><span class="visually-hidden">Ładowanie...</span></div>
          Ładowanie...
        </div>
        <div v-else-if="error" class="alert alert-danger text-center mx-auto" style="max-width:500px;">{{ error }}</div>
        <div v-else-if="pageRecord" class="card shadow-sm page-record-card">
          <div class="card-body">
            <h2 class="mb-4 text-primary">{{ pageRecord.title }}</h2>
            <div class="row g-3 page-record-details">
                <!-- To tworzy krytyczna podatność na xss -->
                <div class="col-sm-12 col-md-12"><strong>{{ $t('page_record.content') }}:</strong><div v-html="pageRecord.content"/></div>
                <div class="col-12"><strong>{{ $t('page_record.textField1') }}:</strong> <div v-html="pageRecord.textField1"/></div>
                <div class="col-12"><strong>{{ $t('page_record.textField2') }}:</strong> <div v-html="pageRecord.textField2"/></div>
                <div class="col-12"><strong>{{ $t('page_record.textField3') }}:</strong> <div v-html="pageRecord.textField3"/> <img :src="pageRecord.textField3"></div>
                <!--  -->
                <div class="col-12"><strong>{{ $t('page_record.pageUrl') }}:</strong> <a :href="pageRecord.url" target="_blank">Link</a></div>
                <div class="col-sm-6 col-md-4"><strong>{{ $t('page_record.createdAt') }}:</strong> {{ pageRecord.createdAt.toLocaleString() }}</div>
                <div class="col-sm-6 col-md-4"><strong>{{ $t('page_record.updatedAt') }}:</strong> {{ pageRecord.updatedAt.toLocaleString() }}</div>
            </div>
            <!-- <div v-if="pageRecord.imagesUrls && pageRecord.imagesUrls.length" class="mt-4">
              <strong>{{ $t('page_record.images') }}:</strong>
              <div class="d-flex flex-wrap gap-2 mt-2">
                <img v-for="(img, idx) in pageRecord.imagesUrls" :key="idx" :src="img" alt="Zdjęcie strony" class="page-img">
              </div>
            </div> -->
          </div>
        </div>
        <div v-else class="alert alert-warning text-center mt-4">Nie znaleziono strony.</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { PageType } from '@/types/page';

definePageMeta({
    layout: 'panel',
    middleware: ['auth']
})

const route = useRoute();
const id = String(route.params.id ?? '');
const pageRecord = ref<PageType|null>(null);
const loading = ref(true);
const error = ref<string|null>(null);

const { fetchPageById } = usePages();

async function loadPageRecord() {
  loading.value = true;
  error.value = null;
  const result = await fetchPageById(id);
  pageRecord.value = result.pageRecord;
  error.value = result.error;
  loading.value = false;
}

onMounted(loadPageRecord);

</script>

<style lang="scss" scoped>
.page-record-card {
  border-left: 5px solid #1976d2;
}
.page-record-img {
  max-width: 120px;
  max-height: 90px;
  border-radius: 4px;
  box-shadow: 0 1px 4px rgba(0,0,0,0.08);
}
</style>
