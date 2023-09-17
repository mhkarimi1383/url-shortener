<script setup lang="ts">
import {
  KeyOutlined,
  HeartFilled,
  UserOutlined,
  HomeOutlined,
  InfoOutlined,
  LinkOutlined,
  GroupOutlined,
  LoginOutlined,
  LogoutOutlined,
  GithubOutlined,
  VerifiedOutlined,
} from '@ant-design/icons-vue';
import type { VNode } from 'vue';
import { theme } from 'ant-design-vue';
import type { VueCookies } from 'vue-cookies';
import { getBrowserTheme } from '@/lib/utils';
import MoonIcon from '@/components/MoonIcon.vue';
import { RouterLink, RouterView } from 'vue-router';
import { reactive, watch, h, inject, ref } from 'vue';
import { setToken, loginStateCookie, type loginResponse } from '@/lib/api';

const currentTheme = ref<string>(localStorage.getItem('theme') || getBrowserTheme() || 'light');

const $cookies = inject<VueCookies>('$cookies');

const changeTheme = () => {
  if (currentTheme.value === 'dark') {
    currentTheme.value = 'light';
  } else {
    currentTheme.value = 'dark';
  }
  localStorage.setItem('theme', currentTheme.value);
};

let Admin = false;
let loggedIn = false;

if ($cookies?.get(loginStateCookie)) {
  const info = $cookies.get(loginStateCookie) as loginResponse;
  setToken(info.Token);
  loggedIn = true;
  Admin = info.Info.Admin;
}

interface MenuItem {
  icon: () => VNode;
  title: string;
  key: string;
  danger?: boolean;
  disabled?: boolean;
  label: string | VNode;
  children?: MenuItem[];
}

const currentRouteName = window.location.pathname.replace('/ui', '');
const state = reactive({
  collapsed: false,
  selectedKeys: [currentRouteName || '/'],
  openKeys: [],
  preOpenKeys: [],
});

const items = reactive<MenuItem[]>([
  {
    key: '/',
    icon: () => h(HomeOutlined),
    disabled: !loggedIn,
    label: (loggedIn && h(RouterLink, { to: '/' }, 'Home')) || 'Home',
    title: 'Home',
  },
  {
    key: '/url',
    icon: () => h(LinkOutlined),
    disabled: !loggedIn,
    label: (loggedIn && h(RouterLink, { to: '/url' }, 'URL')) || 'URL',
    title: 'URL',
  },
  {
    key: '/entity',
    disabled: !Admin,
    icon: () => h(GroupOutlined),
    label: (Admin && h(RouterLink, { to: '/entity' }, 'Entity')) || 'Entity',
    title: 'Entity',
  },
  {
    key: 'user',
    icon: () => h(UserOutlined),
    label: 'User',
    title: 'User',
    children: [
      {
        icon: () => h(KeyOutlined),
        key: '/user/change-password',
        disabled: !loggedIn,
        label:
          (loggedIn && h(RouterLink, { to: '/user/change-password' }, 'Change Password')) ||
          'Change Password',
        title: 'Change Password',
      },
      {
        icon: () => h(VerifiedOutlined),
        disabled: !Admin,
        key: '/user/manage',
        label: (Admin && h(RouterLink, { to: '/user/manage' }, 'Manage Users')) || 'Manage Users',
        title: 'Manage Users',
      },
      {
        icon: () => h(LoginOutlined),
        key: '/user/login',
        label: h(RouterLink, { to: '/user/login' }, 'Login'),
        title: 'Login',
      },
      {
        icon: () => h(LogoutOutlined),
        danger: true,
        disabled: !loggedIn,
        key: '/user/logout',
        label: (loggedIn && h(RouterLink, { to: '/user/logout' }, 'Logout')) || 'Logout',
        title: 'Logout',
      },
    ],
  },
  {
    key: '/about',
    icon: () => h(InfoOutlined),
    label: h(RouterLink, { to: '/about' }, 'About'),
    title: 'About',
  },
]);

watch(
  () => state.openKeys,
  (_val, oldVal) => {
    state.preOpenKeys = oldVal;
  },
);
</script>

<template>
  <a-config-provider
    :theme="{
      algorithm: (currentTheme === 'dark' && theme.darkAlgorithm) || theme.defaultAlgorithm,
    }"
  >
    <a-float-button
      @click="changeTheme"
      :type="(currentTheme === 'dark' && 'primary') || 'default'"
    >
      <template #icon>
        <MoonIcon />
      </template>
    </a-float-button>
    <a-page-header
      :ghost="false"
      :style="`
        border: 1px solid ${(currentTheme === 'dark' && '#424242FF') || '#d9d9d9ff'};
        border-radius: 5px;
        height: max-content;
        margin-top: 0.25%;
        margin-bottom: 0.25%;
        margin-left: 0.25%;
        margin-right: 0.25%;`"
      title="URL Shortener"
      :avatar="{
        src: '/ui/logo.svg',
        shape: 'square',
        style: currentTheme === 'dark' && 'filter: invert(100%);',
      }"
      sub-title="Simple and minimalism URL Shortener"
    >
      <template #extra>
        <a-button href="https://github.com/mhkarimi1383/url-shortener">
          <template #icon>
            <GithubOutlined />
          </template>
        </a-button>
      </template>
    </a-page-header>
    <a-layout>
      <a-layout-sider
        :theme="currentTheme"
        breakpoint="sm"
        v-model:collapsed="state.collapsed"
        collapsible
      >
        <a-menu
          v-model:openKeys="state.openKeys"
          v-model:selectedKeys="state.selectedKeys"
          mode="vertical"
          :inline-collapsed="state.collapsed"
          :items="items"
        ></a-menu>
      </a-layout-sider>
      <a-layout style="margin-left: 3%; margin-top: 2%; margin-bottom: 2%; margin-right: 3%">
        <a-layout-content>
          <RouterView />
        </a-layout-content>
      </a-layout>
    </a-layout>
    <a-layout-footer style="text-align: center">
      Made with
      <HeartFilled style="color: red" />
      by <a href="https://github.com/mhkarimi1383">mhkarimi1383</a>
    </a-layout-footer>
  </a-config-provider>
</template>
