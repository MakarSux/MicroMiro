// src/components/PropertiesPanel.vue
<template>
  <div v-if="selectedShape" class="properties-panel">
    <h3>Свойства</h3>
    <label>
      Цвет:
      <input type="color" v-model="selectedShape.color" />
    </label>
  </div>
</template>

<script setup>
import { computed } from 'vue';
import { useCanvasStore } from '@/store/canvas';

const store = useCanvasStore();
const selectedShape = computed(() => {
  return store.shapes.find((s) => s.id === store.selectedElementId);
});
</script>

<style>
.properties-panel {
  position: absolute;
  top: 10px;
  right: 10px;
  background: #fff;
  padding: 10px;
  border-radius: 5px;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
}
</style>