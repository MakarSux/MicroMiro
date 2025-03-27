<template>
  <AuthGuard>
    <div class="boards-container">
      <header class="boards-header">
        <h1>Мои доски</h1>
        <button class="create-board-btn" @click="showCreateBoardModal = true">
          <span class="plus-icon">+</span> Создать доску
        </button>
      </header>

      <div v-if="loading" class="loading">
        <div class="spinner"></div>
        <p>Загрузка досок...</p>
      </div>

      <div v-else-if="error" class="error-message">
        <p>{{ error }}</p>
        <button @click="fetchBoards">Попробовать снова</button>
      </div>

      <div v-else-if="boards.length === 0" class="no-boards">
        <div class="empty-state">
          <img src="@/assets/img/empty-boards.svg" alt="Нет досок" class="empty-icon">
          <h3>У вас пока нет досок</h3>
          <p>Создайте свою первую доску, чтобы начать работу</p>
          <button class="create-board-btn" @click="showCreateBoardModal = true">
            <span class="plus-icon">+</span> Создать доску
          </button>
        </div>
      </div>

      <div v-else class="boards-grid">
        <div v-for="board in boards" :key="board.id" class="board-card" @click="openBoard(board.id)">
          <div class="board-card-content">
            <h3>{{ board.title }}</h3>
            <p v-if="board.description">{{ board.description }}</p>
            <p v-else class="no-description">Без описания</p>
          </div>
          <div class="board-card-footer">
            <span class="created-date">Создано: {{ formatDate(board.created_at) }}</span>
            <span v-if="board.is_public" class="public-badge">Публичная</span>
          </div>
        </div>
      </div>

      <!-- Модальное окно создания доски -->
      <div v-if="showCreateBoardModal" class="modal-overlay" @click.self="showCreateBoardModal = false">
        <div class="modal-content">
          <h2>Создать новую доску</h2>
          <form @submit.prevent="createBoard">
            <div class="form-group">
              <label for="boardTitle">Название доски:</label>
              <input 
                type="text" 
                id="boardTitle" 
                v-model="newBoard.title" 
                required 
                placeholder="Введите название доски"
              />
            </div>
            <div class="form-group">
              <label for="boardDescription">Описание (необязательно):</label>
              <textarea 
                id="boardDescription" 
                v-model="newBoard.description" 
                placeholder="Введите описание доски"
                rows="3"
              ></textarea>
            </div>
            <div class="form-group checkbox-group">
              <input 
                type="checkbox" 
                id="isPublic" 
                v-model="newBoard.is_public"
              />
              <label for="isPublic">Публичная доска</label>
            </div>
            <div v-if="createError" class="error-message">
              {{ createError }}
            </div>
            <div class="modal-actions">
              <button type="button" class="cancel-btn" @click="showCreateBoardModal = false">Отмена</button>
              <button type="submit" class="create-btn" :disabled="createLoading">
                {{ createLoading ? 'Создание...' : 'Создать' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </AuthGuard>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import axios from 'axios';
import AuthGuard from '@/components/auth/AuthGuard.vue';

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1';
const router = useRouter();

// Состояние для списка досок
const boards = ref([]);
const loading = ref(true);
const error = ref('');

// Состояние для создания новой доски
const showCreateBoardModal = ref(false);
const newBoard = ref({
  title: '',
  description: '',
  is_public: false
});
const createLoading = ref(false);
const createError = ref('');

// Получение списка досок
const fetchBoards = async () => {
  loading.value = true;
  error.value = '';
  
  try {
    const token = localStorage.getItem('token');
    if (!token) {
      router.push('/auth');
      return;
    }
    
    const response = await axios.get(`${API_URL}/protected/boards`, {
      headers: {
        Authorization: `Bearer ${token}`
      }
    });
    
    boards.value = response.data;
  } catch (err) {
    console.error('Ошибка при получении досок:', err);
    if (err.response && err.response.status === 401) {
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      router.push('/auth');
    } else {
      error.value = 'Не удалось загрузить доски. Попробуйте позже.';
    }
  } finally {
    loading.value = false;
  }
};

// Создание новой доски
const createBoard = async () => {
  createLoading.value = true;
  createError.value = '';
  
  try {
    const token = localStorage.getItem('token');
    if (!token) {
      router.push('/auth');
      return;
    }
    
    await axios.post(`${API_URL}/protected/boards`, newBoard.value, {
      headers: {
        Authorization: `Bearer ${token}`
      }
    });
    
    // Сбрасываем форму и закрываем модальное окно
    newBoard.value = {
      title: '',
      description: '',
      is_public: false
    };
    showCreateBoardModal.value = false;
    
    // Обновляем список досок
    await fetchBoards();
  } catch (err) {
    console.error('Ошибка при создании доски:', err);
    if (err.response && err.response.status === 401) {
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      router.push('/auth');
    } else {
      createError.value = err.response?.data?.error || 'Не удалось создать доску. Попробуйте позже.';
    }
  } finally {
    createLoading.value = false;
  }
};

// Открытие доски для редактирования
const openBoard = (boardId) => {
  router.push(`/board/${boardId}`);
};

// Форматирование даты
const formatDate = (dateString) => {
  const date = new Date(dateString);
  return date.toLocaleDateString('ru-RU');
};

// Загружаем доски при монтировании компонента
onMounted(() => {
  fetchBoards();
});
</script>

<style scoped>
.boards-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.boards-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
}

.boards-header h1 {
  font-size: 24px;
  color: #333;
}

.create-board-btn {
  display: flex;
  align-items: center;
  background-color: #4a6cf7;
  color: white;
  border: none;
  border-radius: 4px;
  padding: 10px 16px;
  font-size: 14px;
  cursor: pointer;
  transition: background-color 0.3s;
}

.create-board-btn:hover {
  background-color: #3a5ce5;
}

.plus-icon {
  font-size: 18px;
  margin-right: 6px;
}

.loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 50px 0;
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

.error-message {
  text-align: center;
  color: #e74c3c;
  padding: 20px;
  background-color: #fdf0f0;
  border-radius: 4px;
  margin-bottom: 20px;
}

.error-message button {
  background-color: #e74c3c;
  color: white;
  border: none;
  border-radius: 4px;
  padding: 8px 16px;
  margin-top: 10px;
  cursor: pointer;
}

.error-message button:hover {
  background-color: #c0392b;
}

.no-boards {
  display: flex;
  justify-content: center;
  padding: 50px 0;
}

.empty-state {
  text-align: center;
  max-width: 400px;
}

.empty-icon {
  width: 80px;
  height: 80px;
  margin-bottom: 20px;
  opacity: 0.7;
}

.empty-state h3 {
  font-size: 20px;
  margin-bottom: 10px;
  color: #333;
}

.empty-state p {
  color: #666;
  margin-bottom: 20px;
}

.boards-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 20px;
}

.board-card {
  background-color: white;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  overflow: hidden;
  cursor: pointer;
  transition: transform 0.2s, box-shadow 0.2s;
  display: flex;
  flex-direction: column;
  height: 180px;
}

.board-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.15);
}

.board-card-content {
  padding: 20px;
  flex-grow: 1;
}

.board-card h3 {
  font-size: 18px;
  margin-bottom: 10px;
  color: #333;
}

.board-card p {
  font-size: 14px;
  color: #666;
  overflow: hidden;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
}

.no-description {
  font-style: italic;
  color: #999;
}

.board-card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 20px;
  background-color: #f9f9f9;
  font-size: 12px;
  color: #666;
}

.public-badge {
  background-color: #4CAF50;
  color: white;
  padding: 3px 8px;
  border-radius: 10px;
  font-size: 11px;
}

/* Модальное окно */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.modal-content {
  background-color: white;
  border-radius: 8px;
  padding: 25px;
  width: 90%;
  max-width: 500px;
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.3);
}

.modal-content h2 {
  margin-bottom: 20px;
  color: #333;
}

.form-group {
  margin-bottom: 15px;
}

.form-group label {
  display: block;
  margin-bottom: 5px;
  font-weight: 500;
  color: #555;
}

.form-group input[type="text"],
.form-group textarea {
  width: 100%;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 16px;
}

.checkbox-group {
  display: flex;
  align-items: center;
}

.checkbox-group input[type="checkbox"] {
  margin-right: 10px;
}

.checkbox-group label {
  margin-bottom: 0;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 20px;
}

.cancel-btn {
  background-color: #f1f1f1;
  color: #333;
  border: none;
  border-radius: 4px;
  padding: 10px 16px;
  cursor: pointer;
}

.cancel-btn:hover {
  background-color: #e1e1e1;
}

.create-btn {
  background-color: #4a6cf7;
  color: white;
  border: none;
  border-radius: 4px;
  padding: 10px 16px;
  cursor: pointer;
}

.create-btn:hover {
  background-color: #3a5ce5;
}

.create-btn:disabled {
  background-color: #a0a0a0;
  cursor: not-allowed;
}
</style>
