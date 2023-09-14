import { createRouter, createWebHistory } from 'vue-router';
import HomeView from '../views/HomeView.vue';

const router = createRouter({
  history: createWebHistory('/ui'),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
    },
    {
      path: '/about',
      name: 'about',
      component: () => import('../views/AboutView.vue'),
    },
    {
      path: '/entity',
      name: 'entity',
      component: () => import('../views/EntityView.vue'),
    },
    {
      path: '/url',
      name: 'url',
      component: () => import('../views/URLView.vue'),
    },
    {
      path: '/user/change-password',
      name: 'userChangePassword',
      component: () => import('../views/user/ChangePasswordView.vue'),
    },
    {
      path: '/user/manage',
      name: 'userManage',
      component: () => import('../views/user/ManageView.vue'),
    },
    {
      path: '/user/logout',
      name: 'userLogout',
      component: () => import('../views/user/LogoutView.vue'),
    },
    {
      path: '/user/login',
      name: 'userLogin',
      component: () => import('../views/user/LoginView.vue'),
    },
    {
      path: '/user/register',
      name: 'userRegister',
      component: () => import('../views/user/RegisterView.vue'),
    },
    {
      path: '/error/404',
      name: '404Error',
      component: () => import('../views/errorPages/404View.vue'),
    },
    {
      path: '/error/500',
      name: '500Error',
      component: () => import('../views/errorPages/500View.vue'),
    },
    {
      path: '/error/403',
      name: '403Error',
      component: () => import('../views/errorPages/403View.vue'),
    },
    {
      path: '/:pathMatch(.*)*',
      name: '404Error',
      component: () => import('../views/errorPages/404View.vue'),
    },
  ],
});

export default router;
