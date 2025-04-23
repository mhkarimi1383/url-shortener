<script setup lang="ts">
import { message } from 'ant-design-vue';
import { computed, reactive, ref } from 'vue';
import { changeUserPassword, type errorResponse } from '@/lib/api';

interface ChangePasswordFormState {
  password: string;
  repeatPassword: string;
}

const changePasswordLoading = ref<boolean>(false);

const changePasswordFormState = reactive<ChangePasswordFormState>({
  password: '',
  repeatPassword: '',
});

const passwordsAreEqual = computed(() => {
  return changePasswordFormState.password === changePasswordFormState.repeatPassword;
});

const confirmChangePassword = () => {
  changePasswordLoading.value = true;
  if (passwordsAreEqual.value) {
    changeUserPassword({
      Password: changePasswordFormState.password,
    })
      .then((data) => {
        if ((data as errorResponse).message) {
          message.error((data as errorResponse).message);
        } else {
          message.success('Password changed successfully');
        }
      })
      .catch((data) => {
        message.error((data as errorResponse).message);
      })
      .finally(() => {
        changePasswordLoading.value = false;
      });
  } else {
    message.error('Passwords should be equal');
    changePasswordLoading.value = false;
  }
};
</script>

<template>
  <a-form :model="changePasswordFormState" @finish="confirmChangePassword">
    <a-form-item
      label="Password"
      name="password"
      :rules="[{ required: true, message: 'Please input your password!' }]"
    >
      <a-input-password v-model:value="changePasswordFormState.password">
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
      <a-input-password v-model:value="changePasswordFormState.repeatPassword">
        <template #prefix>
          <LockOutlined />
        </template>
      </a-input-password>
    </a-form-item>
    <a-form-item>
      <a-button type="primary" html-type="submit">Submit</a-button>
    </a-form-item>
  </a-form>
</template>
