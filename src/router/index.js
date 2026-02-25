import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import ContractorLoginView from "@/views/ContractorLoginView.vue";
import ContractorRegisterView from "@/views/ContractorRegisterView.vue";
import ContractorDashboardView from "@/views/ContractorDashboardView.vue";

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView
    },
    {
      path: '/about',
      name: 'about',
      // route level code-splitting
      // this generates a separate chunk (About.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: () => import('../views/AboutView.vue')
    },
    {
      path: '/admin',
      name: 'admin',
      component: () => import('../views/AdminViewClaude.vue')
    },
    {
      path: '/contractors/register',
      name: 'contractor-register',
      component: ContractorRegisterView
    },
    {
      path: '/contractors/login',
      name: 'contractor-login',
      component: ContractorLoginView
    },
    {
      path: '/contractors/dashboard',
      name: 'contractor-dashboard',
      component: ContractorDashboardView
    }
  ]
})

export default router
