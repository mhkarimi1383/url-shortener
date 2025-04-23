<template>
  <a-form :model="formState" @finish="onFinish" @finish-failed="onFinishFailed">
    <a-form-item
      label="Username"
      name="username"
      :rules="[{ required: true, message: 'Please input your username!' }]"
    >
      <a-input v-model:value="formState.username">
        <template #prefix>
          <UserOutlined />
        </template>
      </a-input>
    </a-form-item>

    <a-form-item
      label="Password"
      name="password"
      :rules="[{ required: true, message: 'Please input your password!' }]"
    >
      <a-input-password v-model:value="formState.password">
        <template #prefix>
          <LockOutlined />
        </template>
      </a-input-password>
    </a-form-item>

    <a-form-item>
      <a-form-item name="remember" no-style>
        <a-checkbox v-model:checked="formState.remember"
          >Remember me <small>(Username and Password will be saved as cookie)</small></a-checkbox
        >
      </a-form-item>
    </a-form-item>

    <a-form-item>
      <a-button :disabled="disabled" type="primary" html-type="submit"> Log in </a-button>
      Or
      <RouterLink to="/user/register">register now!</RouterLink>
    </a-form-item>
  </a-form>
</template>

<script lang="ts" setup>
import router from '@/router';
import { message } from 'ant-design-vue';
import type { VueCookies } from 'vue-cookies';
import { computed, inject, reactive } from 'vue';
import { LockOutlined, UserOutlined } from '@ant-design/icons-vue';
import type { errorResponse, loginInfo, loginResponse } from '@/lib/api';
import { login, loginInfoCookie, loginStateCookie, setToken } from '@/lib/api';

const $cookies = inject<VueCookies>('$cookies');

interface FormState {
  username: string;
  password: string;
  remember: boolean;
}

const formState = reactive<FormState>({
  username: '',
  password: '',
  remember: false,
});

if ($cookies?.get(loginInfoCookie)) {
  const info = $cookies.get(loginInfoCookie) as loginInfo;
  formState.remember = true;
  formState.password = info.Password;
  formState.username = info.Username;
}

async function onFinish(values: any) {
  const finish = message.loading('Logging in');
  const info = <loginInfo>{
    Password: values.password,
    Username: values.username,
  };
  const resp = await login(info);
  if ((resp as loginResponse).Token) {
    if (values.remember === true) {
      $cookies?.set(loginInfoCookie, info);
    }
    $cookies?.set(loginStateCookie, resp);
    setTimeout(finish, 1000);
    message.success(`Welcome, ${(resp as loginResponse).Info.Username}`);
    setToken((resp as loginResponse).Token);
    router.push('/').finally(() => location.reload());
  } else {
    finish();
    message.error((resp as errorResponse).message);
  }
}

function onFinishFailed(errorInfo: any) {
  message.error(errorInfo);
}

const disabled = computed(() => {
  return !(formState.username && formState.password);
});
</script>
<style scoped></style>
