import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'

const router = createRouter({
  history: createWebHistory("/ui"),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView
    },
    {
      path: '/about',
      name: 'about',
      component: () => import('../views/AboutView.vue')
    },
    {
      path: '/entity',
      name: 'entity',
      component: () => import('../views/EntityView.vue')
    },
    {
      path: '/url',
      name: 'url',
      component: () => import('../views/URLView.vue')
    },
    {
      path: '/user/change-password',
      name: 'userChangePassword',
      component: () => import('../views/user/ChangePasswordView.vue')
    },
    {
      path: '/user/manage',
      name: 'userManage',
      component: () => import('../views/user/ManageView.vue')
    },
    {
      path: '/user/logout',
      name: 'userLogout',
      component: () => import('../views/user/LogoutView.vue')
    },
    {
      path: '/user/login',
      name: 'userLogin',
      component: () => import('../views/user/LoginView.vue')
    },
  ]
})

export default router
