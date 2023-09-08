<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router'
import { reactive, watch, h, inject } from 'vue';
import type { VNode } from 'vue';
import { setToken, loginStateCookie, type loginResponse } from '@/lib/api';
import {
  UserOutlined,
  HomeOutlined,
  InfoOutlined,
  LogoutOutlined,
  KeyOutlined,
  LinkOutlined,
  VerifiedOutlined,
  GroupOutlined,
  GithubOutlined,
} from '@ant-design/icons-vue';
import type { VueCookies } from 'vue-cookies';

const $cookies = inject<VueCookies>("$cookies");

if ($cookies?.get(loginStateCookie)) {
  const info = $cookies.get(loginStateCookie) as loginResponse;
  setToken(info.Token);
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

const Admin = false;

const currentRouteName = window.location.pathname
const state = reactive({
  collapsed: false,
  selectedKeys: [currentRouteName || "/"],
  openKeys: [],
  preOpenKeys: [],
});

const items = reactive<MenuItem[]>([
  {
    key: '/',
    icon: () => h(HomeOutlined),
    label: h(RouterLink, { to: '/' }, 'Home'),
    title: 'Home',
  },
  {
    key: '/about',
    icon: () => h(InfoOutlined),
    label: h(RouterLink, { to: '/about' }, 'About'),
    title: 'About',
  },
  {
    key: '/url',
    icon: () => h(LinkOutlined),
    label: h(RouterLink, { to: '/url' }, 'URL'),
    title: 'URL',
  },
  {
    key: '/entity',
    disabled: !Admin,
    icon: () => h(GroupOutlined),
    label: Admin && h(RouterLink, { to: '/entity' }, 'Entity') || 'Entity',
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
        label: h(RouterLink, { to: '/user/change-password' }, 'Change Password'),
        title: 'Change Password',
      },
      {
        icon: () => h(VerifiedOutlined),
        disabled: !Admin,
        key: '/user/manage',
        label: Admin && h(RouterLink, { to: '/user/manage' }, 'Manage Users') || 'Manage Users',
        title: 'Manage Users',
      },
      {
        icon: () => h(LogoutOutlined),
        danger: true,
        key: '/user/logout',
        label: h(RouterLink, { to: '/user/logout' }, 'Logout'),
        title: 'Logout',
      },
    ],
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
  <a-page-header style="border: 1px solid rgb(235, 237, 240);
        height: max-content;
        margin-top: 0.25%;
        margin-bottom: 0.25%;
        margin-left: 0.25%;
        margin-right: 0.25%;" title="URL Shortener" :avatar="{ src: 'https://www.antdv.com/assets/logo.1ef800a8.svg' }"
    sub-title="Simple and minimalism URL Shortener">
    <template #extra>
      <a-button href="https://github.com/mhkarimi1383/url-shortener">
        <template #icon>
          <GithubOutlined />
        </template>
      </a-button>
    </template>
  </a-page-header>
  <a-layout>
    <a-layout-sider breakpoint="sm" theme="light" v-model:collapsed="state.collapsed" collapsible>
      <a-menu v-model:openKeys="state.openKeys" v-model:selectedKeys="state.selectedKeys" mode="vertical"
        :inline-collapsed="state.collapsed" :items="items"></a-menu>
    </a-layout-sider>
    <a-layout style="margin-left: 3%;
        margin-top: 2%;
        margin-bottom: 2%;
        margin-right: 3%">
      <a-layout-content>
        <RouterView />
      </a-layout-content>
    </a-layout>
  </a-layout>
</template>

<style scoped>
</style>
