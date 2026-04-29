import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', component: HomeView },
    { path: '/todo/:id', component: () => import('../views/TodoDetailView.vue') },
    { path: '/stats', component: () => import('../views/StatsView.vue') },
  ],
})

export default router
