<script setup>
import { RouterView, useRoute } from 'vue-router';
import { computed } from 'vue';
import AppHeader from '@/components/layout/AppHeader.vue';

const route = useRoute();

// Не показываем заголовок на странице авторизации
const showHeader = computed(() => {
  return route.path !== '/auth';
});
</script>

<template>
  <div class="app">
    <AppHeader v-if="showHeader" />
    <main class="main-content">
      <RouterView />
    </main>
  </div>
</template>

<style>
* {
  box-sizing: border-box;
  margin: 0;
  padding: 0;
}

body {
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
  color: #333;
  background-color: #f5f7fb;
  line-height: 1.6;
}

.app {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}

.main-content {
  flex: 1;
  display: flex;
  flex-direction: column;
}

a {
  color: #4a6cf7;
  text-decoration: none;
}

a:hover {
  text-decoration: underline;
}

button {
  cursor: pointer;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
}
</style>
