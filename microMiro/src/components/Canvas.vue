// src/components/Canvas.vue
<template>
  <div class="canvas-container" ref="canvasContainer">
    <div 
      class="canvas" 
      ref="canvas" 
      @mousedown="startDrag" 
      @mousemove="onDrag" 
      @mouseup="endDrag"
      @mouseleave="endDrag"
      @wheel="handleZoom"
      :style="{ transform: `translate(${position.x}px, ${position.y}px) scale(${scale})` }"
    >
      <div 
        v-for="shape in shapes" 
        :key="shape.id" 
        class="shape"
        :class="{ 'selected': selectedShape && selectedShape.id === shape.id }"
        :style="getShapeStyle(shape)"
        @mousedown.stop="selectShape(shape, $event)"
      >
        <div v-if="shape.type === 'text'" class="text-content">
          <textarea 
            v-if="editingText && selectedShape && selectedShape.id === shape.id"
            v-model="editingTextContent"
            @blur="saveTextContent"
            @keydown.enter.prevent="saveTextContent"
            ref="textEditor"
            :style="{ width: '100%', height: '100%' }"
          ></textarea>
          <div v-else class="text-display" @dblclick="startEditingText(shape)">
            {{ shape.content }}
          </div>
        </div>
      </div>
    </div>

    <div class="toolbar">
      <button @click="addRectangle" title="Добавить прямоугольник">
        <span class="icon">□</span>
      </button>
      <button @click="addCircle" title="Добавить круг">
        <span class="icon">○</span>
      </button>
      <button @click="addText" title="Добавить текст">
        <span class="icon">T</span>
      </button>
      <button @click="deleteSelected" title="Удалить выбранный элемент" :disabled="!selectedShape">
        <span class="icon">🗑️</span>
      </button>
      <div class="zoom-controls">
        <button @click="zoomOut" title="Уменьшить">-</button>
        <span>{{ Math.round(scale * 100) }}%</span>
        <button @click="zoomIn" title="Увеличить">+</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, nextTick, defineProps, defineEmits } from 'vue';
import { useCanvasStore } from '@/store/canvas';
import { v4 as uuidv4 } from 'uuid';

const props = defineProps({
  initialElements: {
    type: Array,
    default: () => []
  }
});

const emit = defineEmits(['element-added', 'element-updated', 'element-deleted']);

const canvasStore = useCanvasStore();
const shapes = computed(() => canvasStore.shapes);

const canvas = ref(null);
const canvasContainer = ref(null);
const textEditor = ref(null);

// Состояние перетаскивания
const isDragging = ref(false);
const dragStartPos = ref({ x: 0, y: 0 });
const selectedShape = ref(null);
const isMovingShape = ref(false);

// Состояние масштабирования и позиции
const scale = ref(1);
const position = ref({ x: 0, y: 0 });

// Состояние редактирования текста
const editingText = ref(false);
const editingTextContent = ref('');

// Инициализация канваса
onMounted(() => {
  // Загружаем начальные элементы, если они переданы
  if (props.initialElements && props.initialElements.length > 0) {
    canvasStore.setShapes(props.initialElements);
  }
  
  // Центрируем канвас при инициализации
  if (canvasContainer.value) {
    const containerRect = canvasContainer.value.getBoundingClientRect();
    position.value = {
      x: containerRect.width / 2,
      y: containerRect.height / 2
    };
  }
});

// Функции для работы с фигурами
const addRectangle = () => {
  const newShape = {
    id: uuidv4(),
    type: 'rectangle',
    x: 100,
    y: 100,
    width: 150,
    height: 100,
    color: getRandomColor(),
    zIndex: shapes.value.length + 1
  };
  
  canvasStore.addShape(newShape);
  emit('element-added', newShape);
};

const addCircle = () => {
  const newShape = {
    id: uuidv4(),
    type: 'circle',
    x: 100,
    y: 100,
    radius: 50,
    color: getRandomColor(),
    zIndex: shapes.value.length + 1
  };
  
  canvasStore.addShape(newShape);
  emit('element-added', newShape);
};

const addText = () => {
  const newShape = {
    id: uuidv4(),
    type: 'text',
    x: 100,
    y: 100,
    width: 200,
    height: 100,
    content: 'Двойной клик для редактирования',
    color: '#333333',
    zIndex: shapes.value.length + 1
  };
  
  canvasStore.addShape(newShape);
  emit('element-added', newShape);
};

const deleteSelected = () => {
  if (selectedShape.value) {
    const shapeId = selectedShape.value.id;
    canvasStore.removeShape(shapeId);
    emit('element-deleted', shapeId);
    selectedShape.value = null;
  }
};

const getShapeStyle = (shape) => {
  const baseStyle = {
    position: 'absolute',
    left: `${shape.x}px`,
    top: `${shape.y}px`,
    zIndex: shape.zIndex,
    backgroundColor: shape.color
  };

  if (shape.type === 'rectangle' || shape.type === 'text') {
    return {
      ...baseStyle,
      width: `${shape.width}px`,
      height: `${shape.height}px`,
      borderRadius: '4px'
    };
  } else if (shape.type === 'circle') {
    return {
      ...baseStyle,
      width: `${shape.radius * 2}px`,
      height: `${shape.radius * 2}px`,
      borderRadius: '50%'
    };
  }

  return baseStyle;
};

// Функции для перетаскивания
const startDrag = (event) => {
  if (!selectedShape.value) {
    isDragging.value = true;
    dragStartPos.value = {
      x: event.clientX,
      y: event.clientY
    };
    event.preventDefault();
  }
};

const onDrag = (event) => {
  if (isDragging.value) {
    const dx = event.clientX - dragStartPos.value.x;
    const dy = event.clientY - dragStartPos.value.y;
    
    position.value = {
      x: position.value.x + dx,
      y: position.value.y + dy
    };
    
    dragStartPos.value = {
      x: event.clientX,
      y: event.clientY
    };
  } else if (isMovingShape.value && selectedShape.value) {
    const dx = (event.clientX - dragStartPos.value.x) / scale.value;
    const dy = (event.clientY - dragStartPos.value.y) / scale.value;
    
    const updatedShape = { ...selectedShape.value };
    updatedShape.x += dx;
    updatedShape.y += dy;
    
    canvasStore.updateShape(updatedShape);
    selectedShape.value = updatedShape;
    
    dragStartPos.value = {
      x: event.clientX,
      y: event.clientY
    };
  }
};

const endDrag = () => {
  if (isMovingShape.value && selectedShape.value) {
    emit('element-updated', selectedShape.value);
  }
  
  isDragging.value = false;
  isMovingShape.value = false;
};

const selectShape = (shape, event) => {
  event.stopPropagation();
  selectedShape.value = shape;
  isMovingShape.value = true;
  dragStartPos.value = {
    x: event.clientX,
    y: event.clientY
  };
};

// Функции для масштабирования
const handleZoom = (event) => {
  event.preventDefault();
  const delta = event.deltaY > 0 ? -0.1 : 0.1;
  const newScale = Math.max(0.1, Math.min(3, scale.value + delta));
  
  // Вычисляем позицию курсора относительно канваса
  const rect = canvasContainer.value.getBoundingClientRect();
  const mouseX = event.clientX - rect.left;
  const mouseY = event.clientY - rect.top;
  
  // Вычисляем точку на холсте, над которой находится курсор
  const pointX = (mouseX - position.value.x) / scale.value;
  const pointY = (mouseY - position.value.y) / scale.value;
  
  // Вычисляем новую позицию, чтобы масштабирование происходило относительно курсора
  const newX = mouseX - pointX * newScale;
  const newY = mouseY - pointY * newScale;
  
  scale.value = newScale;
  position.value = { x: newX, y: newY };
};

const zoomIn = () => {
  const newScale = Math.min(3, scale.value + 0.1);
  
  // Масштабирование относительно центра видимой области
  if (canvasContainer.value) {
    const rect = canvasContainer.value.getBoundingClientRect();
    const centerX = rect.width / 2;
    const centerY = rect.height / 2;
    
    const pointX = (centerX - position.value.x) / scale.value;
    const pointY = (centerY - position.value.y) / scale.value;
    
    const newX = centerX - pointX * newScale;
    const newY = centerY - pointY * newScale;
    
    scale.value = newScale;
    position.value = { x: newX, y: newY };
  } else {
    scale.value = newScale;
  }
};

const zoomOut = () => {
  const newScale = Math.max(0.1, scale.value - 0.1);
  
  // Масштабирование относительно центра видимой области
  if (canvasContainer.value) {
    const rect = canvasContainer.value.getBoundingClientRect();
    const centerX = rect.width / 2;
    const centerY = rect.height / 2;
    
    const pointX = (centerX - position.value.x) / scale.value;
    const pointY = (centerY - position.value.y) / scale.value;
    
    const newX = centerX - pointX * newScale;
    const newY = centerY - pointY * newScale;
    
    scale.value = newScale;
    position.value = { x: newX, y: newY };
  } else {
    scale.value = newScale;
  }
};

// Функции для редактирования текста
const startEditingText = (shape) => {
  if (shape.type === 'text') {
    selectedShape.value = shape;
    editingTextContent.value = shape.content;
    editingText.value = true;
    
    nextTick(() => {
      if (textEditor.value) {
        textEditor.value.focus();
      }
    });
  }
};

const saveTextContent = () => {
  if (selectedShape.value && editingText.value) {
    const updatedShape = { ...selectedShape.value, content: editingTextContent.value };
    canvasStore.updateShape(updatedShape);
    selectedShape.value = updatedShape;
    editingText.value = false;
    
    emit('element-updated', updatedShape);
  }
};

// Вспомогательные функции
const getRandomColor = () => {
  const colors = [
    '#FFD1DC', // Светло-розовый
    '#FFECB3', // Светло-желтый
    '#B3E5FC', // Светло-голубой
    '#C8E6C9', // Светло-зеленый
    '#D1C4E9', // Светло-фиолетовый
    '#F5F5F5', // Светло-серый
    '#BBDEFB', // Светло-синий
    '#DCEDC8'  // Светло-лаймовый
  ];
  return colors[Math.floor(Math.random() * colors.length)];
};
</script>

<style scoped>
.canvas-container {
  position: relative;
  width: 100%;
  height: 100%;
  overflow: hidden;
  background-color: #f5f5f5;
  background-image: 
    linear-gradient(rgba(0, 0, 0, 0.1) 1px, transparent 1px),
    linear-gradient(90deg, rgba(0, 0, 0, 0.1) 1px, transparent 1px);
  background-size: 20px 20px;
  background-position: center;
}

.canvas {
  position: absolute;
  transform-origin: 0 0;
  cursor: grab;
  width: 100%;
  height: 100%;
}

.canvas:active {
  cursor: grabbing;
}

.shape {
  position: absolute;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
  cursor: move;
  user-select: none;
  display: flex;
  align-items: center;
  justify-content: center;
}

.shape.selected {
  outline: 2px solid #4a6cf7;
  z-index: 1000 !important;
}

.text-content {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 10px;
  overflow: hidden;
}

.text-display {
  width: 100%;
  height: 100%;
  overflow: auto;
  cursor: text;
}

textarea {
  border: none;
  background: transparent;
  resize: none;
  font-family: inherit;
  font-size: inherit;
  outline: none;
  padding: 0;
}

.toolbar {
  position: absolute;
  top: 10px;
  left: 10px;
  background-color: white;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  padding: 8px;
  display: flex;
  gap: 5px;
  z-index: 1000;
}

.toolbar button {
  width: 36px;
  height: 36px;
  border: none;
  background-color: #f5f5f5;
  border-radius: 4px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
}

.toolbar button:hover {
  background-color: #e0e0e0;
}

.toolbar button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.icon {
  font-size: 18px;
}

.zoom-controls {
  display: flex;
  align-items: center;
  margin-left: 10px;
  gap: 5px;
}

.zoom-controls span {
  font-size: 14px;
  width: 50px;
  text-align: center;
}
</style>