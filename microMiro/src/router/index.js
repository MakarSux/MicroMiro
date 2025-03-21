import { createRouter, createWebHistory } from 'vue-router'
import Canvas from '@/components/Canvas.vue'


const routes = [
    {
        path: '/',
        name: 'Canvas',
        component: Canvas,
    },
]


const router = createRouter({
    history: createWebHistory(),
    routes,
})

export default router