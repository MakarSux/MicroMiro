// src/store/canvas.js
import { defineStore } from 'pinia';
import { ref } from 'vue';

export const useCanvasStore = defineStore('canvas', () => {
    const shapes = ref([]);
  const selectedElementId = ref(null);

  const addShape = (shape) => {
    shapes.value.push(shape);
  };

  const updateShapePosition = (id, x, y) => {
    const shape = shapes.value.find((s) => s.id === id);
    if (shape) {
      shape.x = x;
      shape.y = y;
    }
  };

  const selectElement = (id) => {
    selectedElementId.value = id;
  };

  return { shapes, selectedElementId, addShape, updateShapePosition, selectElement };
});