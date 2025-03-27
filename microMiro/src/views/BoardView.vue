<template>
  <AuthGuard>
    <div class="board-view">
      <header class="board-header">
        <div class="board-info">
          <h1 v-if="!editingTitle" @click="startEditingTitle" class="board-title">
            {{ board.title || 'Без названия' }}
            <span class="edit-icon">✏️</span>
          </h1>
          <div v-else class="title-edit">
            <input 
              type="text" 
              v-model="editableTitle" 
              @blur="saveTitle" 
              @keyup.enter="saveTitle"
              ref="titleInput"
              placeholder="Введите название доски"
            />
          </div>
          <div class="board-actions">
            <button @click="goBack" class="back-btn">
              Назад к доскам
            </button>
            <span v-if="saving" class="saving-status">Сохранение...</span>
            <span v-else-if="lastSaved" class="saved-status">Сохранено {{ formatTime(lastSaved) }}</span>
          </div>
        </div>
      </header>

      <div v-if="loading" class="loading-container">
        <div class="spinner"></div>
        <p>Загрузка доски...</p>
      </div>

      <div v-else-if="error" class="error-container">
        <p>{{ error }}</p>
        <button @click="fetchBoard">Попробовать снова</button>
      </div>

      <div v-else class="canvas-container">
        <Canvas 
          :initialElements="boardElements" 
          @element-added="handleElementAdded"
          @element-updated="handleElementUpdated"
          @element-deleted="handleElementDeleted"
        />
      </div>
    </div>
  </AuthGuard>
</template>

<script setup>
import { ref, onMounted, nextTick } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import axios from 'axios';
import Canvas from '@/components/Canvas.vue';
import { useCanvasStore } from '@/store/canvas';
import AuthGuard from '@/components/auth/AuthGuard.vue';

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1';
const route = useRoute();
const router = useRouter();
const canvasStore = useCanvasStore();

// Состояние доски
const board = ref({
  id: null,
  title: '',
  description: '',
  is_public: false,
  creator_id: null,
  created_at: null,
  updated_at: null
});
const boardElements = ref([]);
const loading = ref(true);
const error = ref('');
const saving = ref(false);
const lastSaved = ref(null);

// Редактирование заголовка
const editingTitle = ref(false);
const editableTitle = ref('');
const titleInput = ref(null);

// Получение данных доски
const fetchBoard = async () => {
  loading.value = true;
  error.value = '';
  
  try {
    const token = localStorage.getItem('token');
    if (!token) {
      router.push('/auth');
      return;
    }
    
    const boardId = route.params.id;
    const response = await axios.get(`${API_URL}/protected/boards/${boardId}`, {
      headers: {
        Authorization: `Bearer ${token}`
      }
    });
    
    board.value = response.data.board;
    boardElements.value = response.data.elements || [];
    
    // Инициализируем store элементами доски
    canvasStore.setShapes(boardElements.value);
  } catch (err) {
    console.error('Ошибка при получении доски:', err);
    if (err.response && err.response.status === 401) {
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      router.push('/auth');
    } else if (err.response && err.response.status === 404) {
      error.value = 'Доска не найдена';
    } else {
      error.value = 'Не удалось загрузить доску. Попробуйте позже.';
    }
  } finally {
    loading.value = false;
  }
};

// Обработчики событий элементов доски
const handleElementAdded = async (element) => {
  try {
    saving.value = true;
    const token = localStorage.getItem('token');
    if (!token) {
      router.push('/auth');
      return;
    }
    
    const boardId = route.params.id;
    await axios.post(`${API_URL}/protected/boards/${boardId}/elements`, element, {
      headers: {
        Authorization: `Bearer ${token}`
      }
    });
    
    lastSaved.value = new Date();
  } catch (err) {
    console.error('Ошибка при добавлении элемента:', err);
    // Можно добавить уведомление об ошибке
  } finally {
    saving.value = false;
  }
};

const handleElementUpdated = async (element) => {
  try {
    saving.value = true;
    const token = localStorage.getItem('token');
    if (!token) {
      router.push('/auth');
      return;
    }
    
    const boardId = route.params.id;
    await axios.put(`${API_URL}/protected/boards/${boardId}/elements/${element.id}`, element, {
      headers: {
        Authorization: `Bearer ${token}`
      }
    });
    
    lastSaved.value = new Date();
  } catch (err) {
    console.error('Ошибка при обновлении элемента:', err);
    // Можно добавить уведомление об ошибке
  } finally {
    saving.value = false;
  }
};

const handleElementDeleted = async (elementId) => {
  try {
    saving.value = true;
    const token = localStorage.getItem('token');
    if (!token) {
      router.push('/auth');
      return;
    }
    
    const boardId = route.params.id;
    await axios.delete(`${API_URL}/protected/boards/${boardId}/elements/${elementId}`, {
      headers: {
        Authorization: `Bearer ${token}`
      }
    });
    
    lastSaved.value = new Date();
  } catch (err) {
    console.error('Ошибка при удалении элемента:', err);
    // Можно добавить уведомление об ошибке
  } finally {
    saving.value = false;
  }
};

// Редактирование заголовка доски
const startEditingTitle = () => {
  editableTitle.value = board.value.title;
  editingTitle.value = true;
  nextTick(() => {
    titleInput.value.focus();
  });
};

const saveTitle = async () => {
  if (editableTitle.value !== board.value.title) {
    try {
      saving.value = true;
      const token = localStorage.getItem('token');
      if (!token) {
        router.push('/auth');
        return;
      }
      
      const boardId = route.params.id;
      await axios.put(`${API_URL}/protected/boards/${boardId}`, {
        title: editableTitle.value,
        description: board.value.description,
        is_public: board.value.is_public
      }, {
        headers: {
          Authorization: `Bearer ${token}`
        }
      });
      
      board.value.title = editableTitle.value;
      lastSaved.value = new Date();
    } catch (err) {
      console.error('Ошибка при обновлении заголовка:', err);
      // Можно добавить уведомление об ошибке
    } finally {
      saving.value = false;
    }
  }
  
  editingTitle.value = false;
};

// Форматирование времени
const formatTime = (date) => {
  return date.toLocaleTimeString('ru-RU', { hour: '2-digit', minute: '2-digit' });
};

// Возврат к списку досок
const goBack = () => {
  router.push('/boards');
};

// Загружаем доску при монтировании компонента
onMounted(() => {
  fetchBoard();
});
</script>

<style scoped>
.board-view {
  display: flex;
  flex-direction: column;
  height: 100vh;
}

.board-header {
  background-color: #fff;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  padding: 10px 20px;
  z-index: 10;
}

.board-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.board-title {
  font-size: 20px;
  margin: 0;
  cursor: pointer;
  display: flex;
  align-items: center;
}

.edit-icon {
  font-size: 16px;
  margin-left: 8px;
  opacity: 0;
  transition: opacity 0.2s;
}

.board-title:hover .edit-icon {
  opacity: 1;
}

.title-edit input {
  font-size: 20px;
  padding: 5px 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
  width: 300px;
}

.board-actions {
  display: flex;
  align-items: center;
  gap: 15px;
}

.back-btn {
  background-color: #f1f1f1;
  color: #333;
  border: none;
  border-radius: 4px;
  padding: 8px 12px;
  cursor: pointer;
  font-size: 14px;
}

.back-btn:hover {
  background-color: #e1e1e1;
}

.saving-status {
  font-size: 14px;
  color: #f39c12;
}

.saved-status {
  font-size: 14px;
  color: #27ae60;
}

.loading-container, .error-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: calc(100vh - 60px);
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

.error-container p {
  color: #e74c3c;
  margin-bottom: 15px;
}

.error-container button {
  background-color: #e74c3c;
  color: white;
  border: none;
  border-radius: 4px;
  padding: 8px 16px;
  cursor: pointer;
}

.error-container button:hover {
  background-color: #c0392b;
}

.canvas-container {
  flex: 1;
  overflow: hidden;
  position: relative;
}
</style>
