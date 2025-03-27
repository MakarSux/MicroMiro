// src/store/canvas.js
import { defineStore } from 'pinia';

export const useCanvasStore = defineStore('canvas', {
  state: () => ({
    shapes: [],
    selectedShapeId: null
  }),
  
  getters: {
    selectedShape: (state) => {
      return state.shapes.find(shape => shape.id === state.selectedShapeId);
    }
  },
  
  actions: {
    addShape(shape) {
      this.shapes.push(shape);
    },
    
    updateShape(updatedShape) {
      const index = this.shapes.findIndex(shape => shape.id === updatedShape.id);
      if (index !== -1) {
        this.shapes[index] = updatedShape;
      }
    },
    
    removeShape(shapeId) {
      this.shapes = this.shapes.filter(shape => shape.id !== shapeId);
      if (this.selectedShapeId === shapeId) {
        this.selectedShapeId = null;
      }
    },
    
    selectShape(shapeId) {
      this.selectedShapeId = shapeId;
    },
    
    clearSelection() {
      this.selectedShapeId = null;
    },
    
    setShapes(shapes) {
      this.shapes = shapes;
    },
    
    clearShapes() {
      this.shapes = [];
      this.selectedShapeId = null;
    }
  }
});