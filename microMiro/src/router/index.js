import { createRouter, createWebHistory } from 'vue-router'
import Canvas from '../components/Canvas.vue'
import AuthView from '../views/AuthView.vue'
import BoardsView from '../views/BoardsView.vue'
import BoardView from '../views/BoardView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      redirect: '/boards'
    },
    {
      path: '/auth',
      name: 'auth',
      component: AuthView,
      meta: { requiresAuth: false }
    },
    {
      path: '/boards',
      name: 'boards',
      component: BoardsView,
      meta: { requiresAuth: true }
    },
    {
      path: '/board/:id',
      name: 'board',
      component: BoardView,
      meta: { requiresAuth: true }
    },
    {
      path: '/canvas',
      name: 'canvas',
      component: Canvas
    }
  ]
})

// Защита маршрутов
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth)

  // Если маршрут требует авторизации и токен отсутствует
  if (requiresAuth && !token) {
    next('/auth')
  } 
  // Если пользователь уже авторизован и пытается перейти на страницу авторизации
  else if (to.path === '/auth' && token) {
    next('/boards')
  } 
  else {
    next()
  }
})

export default router