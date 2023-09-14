<template>
  <a-spin :spinning="loading">
    <a-row>
      <a-col :span="12">
        <a-card>
          <a-statistic title="Number of Users" :value="resp?.MetaData.Count" />
        </a-card>
      </a-col>
    </a-row>
    <br />
    <a-table
      :columns="columns"
      :pagination="{ total: resp?.MetaData.Count, current: offset, pageSize: limit }"
      :data-source="resp?.Result"
    >
      <template #headerCell="{ column }">
        <template v-if="column.key === 'id'">
          <NumberOutlined />
        </template>
        <template v-if="column.key === 'admin'"> Role </template>
      </template>

      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'admin'">
          <a-tag :color="record.Admin ? 'volcano' : 'green'">
            {{ (record.Admin && 'Admin') || 'Normal' }}
          </a-tag>
        </template>
        <template v-else-if="column.key === 'actions'">
          <a-button type="primary" :danger="true" @click="showChangePasswordModal(record)">
            <template #icon>
              <KeyOutlined />
            </template>
            Change Password
          </a-button>
        </template>
      </template>
    </a-table>
  </a-spin>
  <a-modal
    warning
    v-model:open="changePasswordModalVisible"
    :confirm-loading="changePasswordModalLoading"
    @ok="confirmChangePassword"
  >
    <template #title>
      <ExclamationCircleOutlined style="color: #d89614" /> Changing
      {{ currentSelectedUser?.Username }}'s Password
    </template>
    <a-form :model="changePasswordFormState">
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
    </a-form>
  </a-modal>
</template>

<script setup lang="ts">
import type { listUsersResponse, errorResponse, userInfo } from '@/lib/api';
import { listUsers, changeUserPassword } from '@/lib/api';
import { message } from 'ant-design-vue';
import { ref, reactive, computed } from 'vue';
import { NumberOutlined, KeyOutlined, ExclamationCircleOutlined } from '@ant-design/icons-vue';

const limit = ref(5);
const offset = ref(1);
const loading = ref<boolean>(true);
const resp = ref<listUsersResponse>({
  MetaData: {
    Count: NaN,
  },
  Result: [],
});
const changePasswordModalVisible = ref<boolean>(false);
const changePasswordModalLoading = ref<boolean>(false);
const currentSelectedUser = ref<userInfo | null>(null);

const showChangePasswordModal = (user: userInfo) => {
  currentSelectedUser.value = user;
  changePasswordModalVisible.value = true;
};

interface ChangePasswordFormState {
  password: string;
  repeatPassword: string;
}

const changePasswordFormState = reactive<ChangePasswordFormState>({
  password: '',
  repeatPassword: '',
});

const passwordsAreEqual = computed(() => {
  return changePasswordFormState.password === changePasswordFormState.repeatPassword;
});

const confirmChangePassword = () => {
  changePasswordModalLoading.value = true;
  if (passwordsAreEqual.value) {
    changeUserPassword((currentSelectedUser.value as userInfo).Id, {
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
        changePasswordModalLoading.value = false;
        changePasswordModalVisible.value = false;
      });
  } else {
    message.error('Passwords should be equal');
    changePasswordModalLoading.value = false;
  }
};

listUsers(limit.value, offset.value - 1)
  .then((data) => {
    resp.value = data as listUsersResponse;
  })
  .catch((data) => {
    message.error((data as errorResponse).message);
  })
  .finally(() => {
    loading.value = false;
  });

const columns = [
  {
    name: 'Id',
    dataIndex: 'Id',
    key: 'id',
  },
  {
    title: 'Username',
    dataIndex: 'Username',
    key: 'username',
  },
  {
    title: 'Admin',
    dataIndex: 'Admin',
    key: 'admin',
  },
  {
    title: 'CreatedAt',
    dataIndex: 'CreatedAt',
    key: 'createdAt',
  },
  {
    title: 'UpdatedAt',
    dataIndex: 'UpdatedAt',
    key: 'updatedAt',
  },
  {
    title: 'Version',
    key: 'version',
    dataIndex: 'Version',
  },
  {
    title: 'Actions',
    key: 'actions',
  },
];
</script>
