<script setup lang="ts">
import type { VNode } from 'vue';
import { theme } from 'ant-design-vue';
import { getBrowserTheme } from '@/lib/utils';
import type { VueCookies } from 'vue-cookies';
import MoonIcon from '@/components/MoonIcon.vue';
import { RouterLink, RouterView } from 'vue-router';
import { h, inject, reactive, ref, watch } from 'vue';
import { type loginResponse, loginStateCookie, setToken } from '@/lib/api';
import {
  GithubOutlined,
  GroupOutlined,
  HeartFilled,
  HomeOutlined,
  InfoOutlined,
  KeyOutlined,
  LinkOutlined,
  LoginOutlined,
  LogoutOutlined,
  UserOutlined,
  VerifiedOutlined,
} from '@ant-design/icons-vue';

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
  openKeys: [],
  preOpenKeys: [],
  collapsed: false,
  selectedKeys: [currentRouteName || '/'],
});

const items = reactive<MenuItem[]>([
  {
    key: '/',
    title: 'Home',
    disabled: !loggedIn,
    icon: () => h(HomeOutlined),
    label: (loggedIn && h(RouterLink, { to: '/' }, 'Home')) || 'Home',
  },
  {
    key: '/url',
    title: 'URL',
    disabled: !loggedIn,
    icon: () => h(LinkOutlined),
    label: (loggedIn && h(RouterLink, { to: '/url' }, 'URL')) || 'URL',
  },
  {
    key: '/entity',
    title: 'Entity',
    disabled: !Admin,
    icon: () => h(GroupOutlined),
    label: (Admin && h(RouterLink, { to: '/entity' }, 'Entity')) || 'Entity',
  },
  {
    key: 'user',
    label: 'User',
    title: 'User',
    icon: () => h(UserOutlined),
    children: [
      {
        disabled: !loggedIn,
        title: 'Change Password',
        icon: () => h(KeyOutlined),
        key: '/user/change-password',
        label:
          (loggedIn && h(RouterLink, { to: '/user/change-password' }, 'Change Password')) ||
          'Change Password',
      },
      {
        disabled: !Admin,
        key: '/user/manage',
        title: 'Manage Users',
        icon: () => h(VerifiedOutlined),
        label: (Admin && h(RouterLink, { to: '/user/manage' }, 'Manage Users')) || 'Manage Users',
      },
      {
        title: 'Login',
        key: '/user/login',
        icon: () => h(LoginOutlined),
        label: h(RouterLink, { to: '/user/login' }, 'Login'),
      },
      {
        danger: true,
        title: 'Logout',
        disabled: !loggedIn,
        key: '/user/logout',
        icon: () => h(LogoutOutlined),
        label: (loggedIn && h(RouterLink, { to: '/user/logout' }, 'Logout')) || 'Logout',
      },
    ],
  },
  {
    key: '/about',
    title: 'About',
    icon: () => h(InfoOutlined),
    label: h(RouterLink, { to: '/about' }, 'About'),
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
