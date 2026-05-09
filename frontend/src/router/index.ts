import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', component: HomeView },
    { path: '/join', component: HomeView },
    { path: '/todo/:id', component: () => import('../views/TodoDetailView.vue') },
    { path: '/stats', component: () => import('../views/StatsView.vue') },
    { path: '/focus', component: () => import('../views/FocusView.vue') },
    { path: '/board', component: () => import('../views/KanbanView.vue') },
    { path: '/calendar', component: () => import('../views/CalendarView.vue') },
    { path: '/login', component: () => import('../views/LoginView.vue') },
  ],
})

export default router
