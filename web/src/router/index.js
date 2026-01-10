import { createRouter, createWebHistory } from 'vue-router'
import DashboardView from '../views/DashboardView.vue'
import FeedView from '../views/FeedView.vue'

const routes = [
  { path: '/', name: 'dashboard', component: DashboardView },
  { path: '/feed', name: 'feed', component: FeedView }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
