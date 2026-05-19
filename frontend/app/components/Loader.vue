<template>
  <div class="loader-overlay">
    <div class="spinner-border text-success" role="status" style="width: 3rem; height: 3rem;">
      <span class="visually-hidden">{{ $t('loading') }}...</span>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.loader-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  background: rgba(255,255,255,0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
}
.spinner-border {
  width: 3rem;
  height: 3rem;
}
</style>