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

    <a-form-item
      label="Repeat Password"
      name="repeatPassword"
      :validate-status="!passwordsAreEqual && 'error'"
      help="Passwords should be equal"
      :rules="[{ required: true, message: 'Repeat your password!' }]"
    >
      <a-input-password v-model:value="formState.repeatPassword">
        <template #prefix>
          <LockOutlined />
        </template>
      </a-input-password>
    </a-form-item>

    <a-form-item>
      <a-button :disabled="disabled" type="primary" html-type="submit"> Register </a-button>
      Or
      <RouterLink to="/user/login">login!</RouterLink>
    </a-form-item>
  </a-form>
</template>

<script lang="ts" setup>
import router from '@/router';
import { register } from '@/lib/api';
import { message } from 'ant-design-vue';
import { computed, reactive } from 'vue';
import type { errorResponse, loginInfo } from '@/lib/api';
import { LockOutlined, UserOutlined } from '@ant-design/icons-vue';

interface FormState {
  username: string;
  password: string;
  repeatPassword: string;
}

const formState = reactive<FormState>({
  username: '',
  password: '',
  repeatPassword: '',
});

async function onFinish(values: any) {
  const finish = message.loading('Signing up');

  const info = <loginInfo>{
    Password: values.password,
    Username: values.username,
  };

  const resp = await register(info);
  if ((resp as errorResponse).message) {
    finish();
    message.error((resp as errorResponse).message);
  } else {
    setTimeout(finish, 1000);
    message.success(`Now Login!`);
    router.push('/user/login');
  }
}

function onFinishFailed(errorInfo: any) {
  message.error(errorInfo);
}

const disabled = computed(() => {
  return !(formState.username && formState.password);
});

const passwordsAreEqual = computed(() => {
  return formState.password === formState.repeatPassword;
});
</script>
