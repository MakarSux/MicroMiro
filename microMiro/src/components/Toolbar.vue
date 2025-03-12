// src/components/Toolbar.vue
<template>
  <div class="toolbar">
    <button @click="addRectangle">Добавить прямоугольник</button>
  </div>
</template>

<script setup>
import { useCanvasStore } from '@/store/canvas';

const store = useCanvasStore();

const addRectangle = () => {
  const newShape = {
    id: Date.now(), // Уникальный ID на основе времени
    type: 'rectangle',
    x: 100, // Начальная позиция X
    y: 100, // Начальная позиция Y
    width: 100,
    height: 50,
    color: '#ff0000', // Цвет по умолчанию — красный
  };
  store.addShape(newShape);
};
</script>

<style>
.toolbar {
  position: absolute;
  top: 10px;
  left: 10px;
  background: #fff;
  padding: 10px;
  border-radius: 5px;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
}
</style>