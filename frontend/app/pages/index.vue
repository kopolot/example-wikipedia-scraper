<template>
  <LoginForm/>
</template>

<script setup lang="ts">
if (useAuth().isAuthenticated()) {
  navigateTo('/panel');
}
</script>