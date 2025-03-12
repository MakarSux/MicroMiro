// src/components/Canvas.vue
<script setup>
import { ref, computed, onMounted } from 'vue';
import { useCanvasStore } from '@/store/canvas';
import Shape from './elements/Shape.vue';
import Toolbar from './Toolbar.vue';
import PropertiesPanel from './PropertiesPanel.vue';

const store = useCanvasStore();
const svgCanvas = ref(null);

// Базовые размеры области просмотра
const baseWidth = 1000;
const baseHeight = 1000;

// Реактивные переменные
const zoom = ref(1);
const panX = ref(0);
const panY = ref(0);
const isPanning = ref(false);

const shapes = computed(() => store.shapes);

const handleZoom = (event) => {
  event.preventDefault();
  const scale = event.deltaY > 0 ? 0.9 : 1.1; // Уменьшение или увеличение масштаба
  const rect = svgCanvas.value.getBoundingClientRect();
  const mouseX = (event.clientX - rect.left) / zoom.value - panX.value;
  const mouseY = (event.clientY - rect.top) / zoom.value - panY.value;

  zoom.value *= scale;

  // Корректируем смещение, чтобы зум происходил относительно мыши
  panX.value = (event.clientX - rect.left) / zoom.value - mouseX;
  panY.value = (event.clientY - rect.top) / zoom.value - mouseY;
};

const startPan = (event) => {
  if (event.button === 0) { // Только левая кнопка мыши
    isPanning.value = true;
    window.addEventListener('mousemove', pan);
    window.addEventListener('mouseup', stopPan);
  }
};

const pan = (event) => {
  if (isPanning.value) {
    const rect = svgCanvas.value.getBoundingClientRect();
    panX.value += event.movementX / zoom.value;
    panY.value += event.movementY / zoom.value;
  }
};

const stopPan = () => {
  isPanning.value = false;
  window.removeEventListener('mousemove', pan);
  window.removeEventListener('mouseup', stopPan);
};

// Вычисляемый viewBox
const viewBox = computed(() => `0 0 ${baseWidth} ${baseHeight}`);

onMounted(() => {
  console.log('Canvas initialized');
});
</script>

<template>
  <div class="canvas-container" @wheel="handleZoom" @mousedown="startPan">
    <Toolbar />
    <PropertiesPanel />
    <svg ref="svgCanvas" :viewBox="viewBox">
      <defs>
        <pattern id="grid" width="20" height="20" patternUnits="userSpaceOnUse">
          <path d="M 20 0 L 0 0 0 20" fill="none" stroke="gray" stroke-width="0.5" />
        </pattern>
      </defs>
      <g :transform="`translate(${panX}, ${panY}) scale(${zoom})`">
        <rect x="0" y="0" :width="baseWidth" :height="baseHeight" fill="url(#grid)" />
        <Shape v-for="shape in shapes" :key="shape.id" :shape="shape" />
      </g>
    </svg>
  </div>
</template>

<style>
.canvas-container {
  width: 100%;
  height: 100vh;
  overflow: hidden;
}
svg {
  width: 100%;
  height: 100%;
}
</style>