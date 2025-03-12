// src/components/elements/Shape.vue
<template>
  <rect
    :x="shape.x"
    :y="shape.y"
    :width="shape.width"
    :height="shape.height"
    :fill="shape.color"
    @mousedown="startDrag"
    @click="selectShape"
  />
</template>

<script setup>
import { ref } from 'vue';
import { useCanvasStore } from '@/store/canvas';

const props = defineProps({
  shape: { type: Object, required: true },
});

const store = useCanvasStore();
const isDragging = ref(false);
const offsetX = ref(0);
const offsetY = ref(0);

const startDrag = (event) => {
  isDragging.value = true;
  offsetX.value = event.clientX - props.shape.x;
  offsetY.value = event.clientY - props.shape.y;
  window.addEventListener('mousemove', drag);
  window.addEventListener('mouseup', stopDrag);
};

const drag = (event) => {
  if (isDragging.value) {
    const newX = event.clientX - offsetX.value;
    const newY = event.clientY - offsetY.value;
    store.updateShapePosition(props.shape.id, newX, newY);
  }
};

const stopDrag = () => {
  isDragging.value = false;
  window.removeEventListener('mousemove', drag);
  window.removeEventListener('mouseup', stopDrag);
};

const selectShape = () => {
  store.selectElement(props.shape.id);
};
</script>