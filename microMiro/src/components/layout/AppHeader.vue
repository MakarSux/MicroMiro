<template>
  <header class="app-header">
    <div class="logo">
      <router-link to="/">MicroMiro</router-link>
    </div>
    <nav v-if="isAuthenticated">
      <ul class="nav-links">
        <li><router-link to="/boards">Мои доски</router-link></li>
      </ul>
    </nav>
    <div class="user-menu" v-if="isAuthenticated">
      <div class="user-info" @click="toggleDropdown">
        <span class="user-email">{{ userEmail }}</span>
        <span class="dropdown-icon">▼</span>
      </div>
      <div class="dropdown-menu" v-if="showDropdown">
        <ul>
          <li><router-link to="/boards">Мои доски</router-link></li>
          <li><a href="#" @click.prevent="logout">Выйти</a></li>
        </ul>
      </div>
    </div>
  </header>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue';
import { useRouter } from 'vue-router';

const router = useRouter();
const showDropdown = ref(false);

// Получаем данные пользователя из localStorage
const isAuthenticated = computed(() => {
  return !!localStorage.getItem('token');
});

const userEmail = computed(() => {
  const user = localStorage.getItem('user');
  if (user) {
    try {
      const userData = JSON.parse(user);
      return userData.email || 'Пользователь';
    } catch (e) {
      return 'Пользователь';
    }
  }
  return 'Пользователь';
});

const toggleDropdown = (event) => {
  event.stopPropagation();
  showDropdown.value = !showDropdown.value;
};

const closeDropdown = (e) => {
  // Don't close if clicking inside the dropdown
  const dropdownEl = document.querySelector('.dropdown-menu');
  const userInfoEl = document.querySelector('.user-info');
  
  if (dropdownEl && dropdownEl.contains(e.target)) {
    return;
  }
  
  if (userInfoEl && userInfoEl.contains(e.target)) {
    return;
  }
  
  if (showDropdown.value) {
    showDropdown.value = false;
  }
};

const logout = () => {
  localStorage.removeItem('token');
  localStorage.removeItem('user');
  router.push('/auth');
};

// Закрываем выпадающее меню при клике вне его
onMounted(() => {
  document.addEventListener('click', closeDropdown);
});

onUnmounted(() => {
  document.removeEventListener('click', closeDropdown);
});
</script>

<style scoped>
.app-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 20px;
  height: 60px;
  background-color: white;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  position: relative;
  z-index: 100;
}

.logo a {
  font-size: 22px;
  font-weight: bold;
  color: #4a6cf7;
  text-decoration: none;
}

.nav-links {
  display: flex;
  list-style: none;
  margin: 0;
  padding: 0;
}

.nav-links li {
  margin-left: 20px;
}

.nav-links a {
  color: #333;
  text-decoration: none;
  font-size: 16px;
  transition: color 0.3s;
}

.nav-links a:hover,
.nav-links a.router-link-active {
  color: #4a6cf7;
}

.user-menu {
  position: relative;
}

.user-info {
  display: flex;
  align-items: center;
  cursor: pointer;
  padding: 8px 12px;
  border-radius: 4px;
}

.user-info:hover {
  background-color: #f5f5f5;
}

.user-email {
  margin-right: 8px;
  font-size: 14px;
}

.dropdown-icon {
  font-size: 10px;
  color: #666;
}

.dropdown-menu {
  position: absolute;
  top: 100%;
  right: 0;
  background-color: white;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  border-radius: 4px;
  width: 200px;
  margin-top: 5px;
}

.dropdown-menu ul {
  list-style: none;
  padding: 0;
  margin: 0;
}

.dropdown-menu li {
  padding: 0;
}

.dropdown-menu a {
  display: block;
  padding: 12px 16px;
  color: #333;
  text-decoration: none;
  transition: background-color 0.3s;
}

.dropdown-menu a:hover {
  background-color: #f5f5f5;
}
</style>
