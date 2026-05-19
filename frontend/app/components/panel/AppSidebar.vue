<template>
  <div :class="['sidebar-nav bg-dark text-white position-relative', { collapsed }]" :style="{ width: navWidth + 'px' }">
    <nav :style="{ width: navWidth + 'px' }" class="d-flex flex-column align-items-center pt-3">
      <button :aria-label="collapsed ? $t('menu.show') : $t('menu.hide')" class="toggle-btn btn btn-dark d-flex align-items-center border light-subtle shadow justify-content-center shadow-sm" @click="toggleSidebar" >
        <span v-if="collapsed"><BootstrapIcon name="list"/></span>
        <span v-else><BootstrapIcon name="chevron-left"/></span>
      </button>
      <div v-if="!collapsed" class="logo fw-bold fs-4 text-center my-4">LOGO</div>
      <ul v-if="!collapsed" class="nav flex-column w-100 px-2">
        <li class="nav-item mb-2"><NuxtLink class="nav-link text-white" to="/panel">{{$t('menu.dashboard')}}</NuxtLink></li>
        <li class="nav-item mb-2"><NuxtLink class="nav-link text-white" to="/panel/page-records">{{$t('menu.page_records')}}</NuxtLink></li>
        <li class="nav-item mb-2"><NuxtLink class="nav-link text-white" to="/panel/wanted_filters">{{$t('menu.wanted')}}</NuxtLink></li>
      </ul>
    </nav>
  </div>
</template>

<script lang="ts" setup>
const collapsed = useState('sidebar.collapsed', () => false);

onMounted(() => {
  const saved = localStorage.getItem('sidebar.collapsed');
  if (saved !== null) {
    collapsed.value = saved === 'true';
  }
});

const navWidth = computed(() => (collapsed.value ? 60 : 220));

function toggleSidebar() {
  collapsed.value = !collapsed.value;
}

watch(collapsed, (val) => {
  if (import.meta.client) {
    localStorage.setItem('sidebar.collapsed', val ? 'true' : 'false');
  }
});
</script>

<style lang="scss" scoped>
$trasition-duration: 0.2s;
$transition-ease: cubic-bezier(.4,0,.2,1);
.sidebar-nav {
  height: 100vh;
  min-width: 60px;
  transition: width $trasition-duration $transition-ease;
}
.toggle-btn {
  position: absolute;
  top: 18px;
  left: 100%;
  transform: translateX(-18px);
  width: 36px;
  height: 36px;
  font-size: 1.2rem;
  z-index: 10;
  transition:
    background $trasition-duration $transition-ease,
    transform $trasition-duration $transition-ease;
}
.logo {
  letter-spacing: 2px;
}
</style>
