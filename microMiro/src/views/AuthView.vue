<template>
  <div class="auth-container">
    <div class="auth-logo">
      <h1>MicroMiro</h1>
      <p>Создавайте и делитесь интерактивными досками</p>
    </div>
    <div class="auth-form">
      <transition name="fade" mode="out-in">
        <LoginForm v-if="showLogin" @login-success="handleLoginSuccess" @switch-to-register="showLogin = false" />
        <RegisterForm v-else @register-success="handleRegisterSuccess" @switch-to-login="showLogin = true" />
      </transition>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import LoginForm from '@/components/auth/LoginForm.vue';
import RegisterForm from '@/components/auth/RegisterForm.vue';

const router = useRouter();
const showLogin = ref(true);

const handleLoginSuccess = (userData) => {
  // Сохраняем данные пользователя в localStorage
  localStorage.setItem('user', JSON.stringify(userData));
  router.push('/boards');
};

const handleRegisterSuccess = ({ message }) => {
  showLogin.value = true;
  // Можно добавить уведомление об успешной регистрации
  alert(message || 'Регистрация прошла успешно! Теперь вы можете войти в систему.');
};
</script>

<style scoped>
.auth-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  background-color: #f5f7fb;
  padding: 20px;
}

.auth-logo {
  text-align: center;
  margin-bottom: 30px;
}

.auth-logo h1 {
  font-size: 2.5rem;
  color: #4a6cf7;
  margin-bottom: 10px;
}

.auth-logo p {
  color: #666;
  font-size: 1.1rem;
}

.auth-form {
  width: 100%;
  max-width: 450px;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
