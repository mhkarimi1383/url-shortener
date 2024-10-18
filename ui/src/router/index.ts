import { getToken } from '@/lib/api';
import HomeView from '../views/HomeView.vue';
import { createRouter, createWebHistory } from 'vue-router';

const router = createRouter({
  history: createWebHistory('/BASE_URI/ui'),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
    },
    {
      name: 'about',
      path: '/about',
      component: () => import('../views/AboutView.vue'),
    },
    {
      name: 'entity',
      path: '/entity',
      component: () => import('../views/EntityView.vue'),
    },
    {
      name: 'url',
      path: '/url',
      component: () => import('../views/URLView.vue'),
    },
    {
      name: 'userChangePassword',
      path: '/user/change-password',
      component: () => import('../views/user/ChangePasswordView.vue'),
    },
    {
      name: 'userManage',
      path: '/user/manage',
      component: () => import('../views/user/ManageView.vue'),
    },
    {
      name: 'userLogout',
      path: '/user/logout',
      component: () => import('../views/user/LogoutView.vue'),
    },
    {
      name: 'userLogin',
      path: '/user/login',
      component: () => import('../views/user/LoginView.vue'),
    },
    {
      name: 'userRegister',
      path: '/user/register',
      component: () => import('../views/user/RegisterView.vue'),
    },
    {
      name: '404Error',
      path: '/error/404',
      component: () => import('../views/errorPages/404View.vue'),
    },
    {
      name: '500Error',
      path: '/error/500',
      component: () => import('../views/errorPages/500View.vue'),
    },
    {
      name: '403Error',
      path: '/error/403',
      component: () => import('../views/errorPages/403View.vue'),
    },
    {
      name: '404Error',
      path: '/:pathMatch(.*)*',
      component: () => import('../views/errorPages/404View.vue'),
    },
  ],
});

const publicRoutes = [
  'home',
  'about',
  '404Error',
  '500Error',
  '403Error',
  'userLogin',
  'userRegister',
];

router.beforeEach(async (to, _) => {
  if (!getToken() && !(to.name && publicRoutes.includes(to.name.toString()))) {
    return { name: 'userLogin' };
  }
});

export default router;
