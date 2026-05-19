<template>
  <header class="header d-flex justify-content-between align-items-center shadow-sm bg-white w-100 position-fixed px-5" style="z-index:10000;">
    <div class="left fw-bold fs-5">{{ $t('nav.title') }}</div>
    <div class="right d-flex align-items-center gap-3">
      <div class="user-menu d-flex align-items-center gap-2 px-3">
        <div ref="menuRef" class="dropdown" @click="toggleMenu" >
          <span class="dropdown-toggle dropdown-icon" role="button"/>
          <ul v-if="showMenu" class="dropdown-menu show dropdown-menu-end mt-2">
            <li @click="showMenu = false"><NuxtLink class="dropdown-item" to="/panel/user">{{ $t('nav.profile') }}</NuxtLink></li>
            <!-- <li @click="showMenu = false"><NuxtLink class="dropdown-item" to="/settings">{{ $t('nav.settings') }}</NuxtLink></li> -->
            <li><a class="dropdown-item text-danger" href="#" @click.prevent="logout">{{ $t('nav.logout') }}</a></li>
          </ul>
        </div>
        <span class="user-avatar bg-dark text-white d-flex align-items-center justify-content-center rounded-circle">{{ user?.username.charAt(0).toUpperCase() || 'U' }}</span>
        <span class="user-name text-dark fw-medium">{{ user?.username || user?.email }}</span>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
const { user, logout } = useAuth();
const showMenu = ref(false);
const menuRef = ref(null);

function toggleMenu() {
  showMenu.value = !showMenu.value;
}

onClickOutside(menuRef, () => showMenu.value = false)
</script>

<style lang="scss" scoped>
.header {
  height: $header-height;
}
.user-avatar {
  width: 32px;
  height: 32px;
  font-size: 1.1rem;
}
</style>
