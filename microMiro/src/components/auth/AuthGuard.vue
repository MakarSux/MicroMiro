<template>
  <div v-if="isLoading" class="auth-loading">
    <div class="spinner"></div>
    <p>Проверка авторизации...</p>
  </div>
  <slot v-else></slot>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import axios from 'axios';

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1';
const router = useRouter();
const isLoading = ref(true);

onMounted(async () => {
  const token = localStorage.getItem('token');
  
  if (!token) {
    router.push('/auth');
    return;
  }
  
  try {
    // Проверяем валидность токена
    await axios.get(`${API_URL}/protected/profile`, {
      headers: {
        Authorization: `Bearer ${token}`
      }
    });
    
    // Если запрос успешен, токен валиден
    isLoading.value = false;
  } catch (err) {
    console.error('Ошибка проверки авторизации:', err);
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    router.push('/auth');
  }
});
</script>

<style scoped>
.auth-loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100vh;
  background-color: #f5f7fb;
}

.spinner {
  border: 4px solid rgba(0, 0, 0, 0.1);
  border-radius: 50%;
  border-top: 4px solid #4a6cf7;
  width: 40px;
  height: 40px;
  animation: spin 1s linear infinite;
  margin-bottom: 15px;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}
</style>
