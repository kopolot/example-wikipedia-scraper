<template>
  <div class="panel-layout d-flex flex-wrap min-vh-100 bg-light">
    <Loader v-if="loading" />
    <div class="topbar w-100 position-relative">
      <PanelAppHeader/>
    </div>
    <div class="main-content w-100 flex-grow-1 d-flex flex-nowrap min-width-0">
      <div class="sidebar bg-dark text-white position-relative shadow-sm">
        <PanelAppSidebar/>
      </div>
      <div class="content flex-grow-1 p-4 overflow-auto">
        <slot />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useNuxtApp } from '#app'
import Loader from '../components/Loader.vue'

const loading = ref(true)
const nuxtApp = useNuxtApp()

onMounted(() => {
  nextTick(() => {
    loading.value = false
  })
  nuxtApp.hook('page:start', () => {
    loading.value = true
  })
  nuxtApp.hook('page:finish', () => {
    loading.value = false
  })
})
</script>

<style lang="scss" scoped>
.sidebar {
  min-height: 100vh;
  flex-basis: fit-content;
  background: #23272f;
  color: #fff;
}
.topbar {
  height: $header-height;
}
</style>