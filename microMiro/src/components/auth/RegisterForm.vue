<template>
  <div class="register-form">
    <h2>Регистрация</h2>
    <form @submit.prevent="register">
      <div class="form-group">
        <label for="username">Имя пользователя:</label>
        <input 
          type="text" 
          id="username" 
          v-model="username" 
          required 
          placeholder="Введите имя пользователя"
        />
      </div>
      <div class="form-group">
        <label for="email">Email:</label>
        <input 
          type="email" 
          id="email" 
          v-model="email" 
          required 
          placeholder="Введите email"
        />
      </div>
      <div class="form-group">
        <label for="password">Пароль:</label>
        <input 
          type="password" 
          id="password" 
          v-model="password" 
          required 
          placeholder="Введите пароль (минимум 6 символов)"
          minlength="6"
        />
      </div>
      <div class="form-group">
        <label for="confirmPassword">Подтверждение пароля:</label>
        <input 
          type="password" 
          id="confirmPassword" 
          v-model="confirmPassword" 
          required 
          placeholder="Повторите пароль"
        />
        <div v-if="passwordMismatch" class="error-message">
          Пароли не совпадают
        </div>
      </div>
      <div v-if="error" class="error-message">
        {{ error }}
      </div>
      <button type="submit" :disabled="loading || passwordMismatch">
        {{ loading ? 'Регистрация...' : 'Зарегистрироваться' }}
      </button>
      <div class="form-footer">
        <p>Уже есть аккаунт? <a href="#" @click.prevent="$emit('switch-to-login')">Войти</a></p>
      </div>
    </form>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue';
import axios from 'axios';

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1';

const emit = defineEmits(['register-success', 'switch-to-login']);

const username = ref('');
const email = ref('');
const password = ref('');
const confirmPassword = ref('');
const error = ref('');
const loading = ref(false);

const passwordMismatch = computed(() => {
  return password.value && confirmPassword.value && password.value !== confirmPassword.value;
});

const register = async () => {
  if (passwordMismatch.value) {
    return;
  }
  
  try {
    loading.value = true;
    error.value = '';
    
    const response = await axios.post(`${API_URL}/register`, {
      username: username.value,
      email: email.value,
      password: password.value
    });
    
    emit('register-success', { message: response.data.message });
    emit('switch-to-login');
  } catch (err) {
    console.error('Ошибка регистрации:', err);
    if (err.response) {
      error.value = err.response.data.error || 'Ошибка регистрации. Попробуйте другие данные.';
    } else {
      error.value = 'Ошибка соединения с сервером. Попробуйте позже.';
    }
  } finally {
    loading.value = false;
  }
};
</script>

<style scoped>
.register-form {
  max-width: 400px;
  margin: 0 auto;
  padding: 20px;
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
}

h2 {
  text-align: center;
  margin-bottom: 20px;
  color: #333;
}

.form-group {
  margin-bottom: 15px;
}

label {
  display: block;
  margin-bottom: 5px;
  font-weight: 500;
  color: #555;
}

input {
  width: 100%;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 16px;
}

button {
  width: 100%;
  padding: 12px;
  background-color: #4a6cf7;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 16px;
  cursor: pointer;
  transition: background-color 0.3s;
}

button:hover {
  background-color: #3a5ce5;
}

button:disabled {
  background-color: #a0a0a0;
  cursor: not-allowed;
}

.error-message {
  color: #e74c3c;
  margin-top: 5px;
  font-size: 14px;
}

.form-footer {
  margin-top: 20px;
  text-align: center;
  font-size: 14px;
}

.form-footer a {
  color: #4a6cf7;
  text-decoration: none;
}

.form-footer a:hover {
  text-decoration: underline;
}
</style>
