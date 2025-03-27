<template>
  <div class="login-form">
    <h2>Вход в систему</h2>
    <form @submit.prevent="login">
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
          placeholder="Введите пароль"
        />
      </div>
      <div v-if="error" class="error-message">
        {{ error }}
      </div>
      <button type="submit" :disabled="loading">
        {{ loading ? 'Вход...' : 'Войти' }}
      </button>
      <div class="form-footer">
        <p>Нет аккаунта? <a href="#" @click.prevent="$emit('switch-to-register')">Зарегистрироваться</a></p>
      </div>
    </form>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import axios from 'axios';

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1';

const emit = defineEmits(['login-success', 'switch-to-register']);
const router = useRouter();

const email = ref('');
const password = ref('');
const error = ref('');
const loading = ref(false);

const login = async () => {
  try {
    loading.value = true;
    error.value = '';
    
    const response = await axios.post(`${API_URL}/login`, {
      email: email.value,
      password: password.value
    });
    
    const token = response.data.token;
    localStorage.setItem('token', token);
    
    // Получаем информацию о пользователе
    const userResponse = await axios.get(`${API_URL}/protected/profile`, {
      headers: {
        Authorization: `Bearer ${token}`
      }
    });
    
    const userData = {
      id: userResponse.data.user_id,
      email: userResponse.data.email,
      token: token
    };
    
    localStorage.setItem('user', JSON.stringify(userData));
    
    emit('login-success', userData);
    router.push('/boards');
  } catch (err) {
    console.error('Ошибка входа:', err);
    if (err.response) {
      error.value = err.response.data.error || 'Ошибка входа. Проверьте email и пароль.';
    } else {
      error.value = 'Ошибка соединения с сервером. Попробуйте позже.';
    }
  } finally {
    loading.value = false;
  }
};
</script>

<style scoped>
.login-form {
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
  margin-bottom: 15px;
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
