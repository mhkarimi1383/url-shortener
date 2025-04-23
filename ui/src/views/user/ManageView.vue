<template>
  <a-spin :spinning="loading">
    <a-row>
      <a-col :span="12">
        <a-card>
          <a-statistic title="Number of Users" :value="resp?.MetaData.Count" />
        </a-card>
      </a-col>
      <a-col :span="12">
        <a-card style="width: 100%; height: 100%; display: flex" body-style="align-self: center;">
          <a-form
            layout="inline"
            :model="formState"
            @finish="handleFinish"
            @finish-failed="handleFinishFailed"
          >
            <a-form-item>
              <a-input v-model:value="formState.Username" placeholder="Username"> </a-input>
            </a-form-item>
            <a-form-item>
              <a-input-password v-model:value="formState.Password" placeholder="Password">
              </a-input-password>
            </a-form-item>
            <a-form-item>
              <a-checkbox v-model:checked="formState.Admin">Admin</a-checkbox>
            </a-form-item>
            <a-form-item>
              <a-button
                type="primary"
                html-type="submit"
                :disabled="formState.Username === '' || formState.Password === ''"
              >
                Create
              </a-button>
            </a-form-item>
          </a-form>
        </a-card>
      </a-col>
    </a-row>
    <br />
    <a-table :columns="columns" :pagination="false" :data-source="resp?.Result">
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
    <br />
    <a-pagination
      v-model:current="offset"
      v-model:page-size="limit"
      style="float: right"
      :total="resp?.MetaData.Count"
      @change="loadUserList"
    ></a-pagination>
  </a-spin>
  <a-modal
    v-model:open="changePasswordModalVisible"
    warning
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
import { message } from 'ant-design-vue';
import { computed, reactive, ref } from 'vue';
import type { FormProps } from 'ant-design-vue';
import { adminChangeUserPassword, listUsers, adminCreateUser } from '@/lib/api';
import type { errorResponse, listUsersResponse, userInfo, createUserRequest } from '@/lib/api';
import { ExclamationCircleOutlined, KeyOutlined, NumberOutlined } from '@ant-design/icons-vue';

const limit = ref(5);
const offset = ref(1);
const loading = ref<boolean>(true);
const resp = ref<listUsersResponse>({
  Result: [],
  MetaData: {
    Count: NaN,
  },
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
    adminChangeUserPassword((currentSelectedUser.value as userInfo).Id, {
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

const loadUserList = function () {
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
};

const columns = [
  {
    key: 'id',
    name: 'Id',
    dataIndex: 'Id',
  },
  {
    key: 'username',
    title: 'Username',
    dataIndex: 'Username',
  },
  {
    key: 'admin',
    title: 'Admin',
    dataIndex: 'Admin',
  },
  {
    key: 'createdAt',
    title: 'CreatedAt',
    dataIndex: 'CreatedAt',
  },
  {
    key: 'updatedAt',
    title: 'UpdatedAt',
    dataIndex: 'UpdatedAt',
  },
  {
    key: 'version',
    title: 'Version',
    dataIndex: 'Version',
  },
  {
    key: 'actions',
    title: 'Actions',
  },
];

const formState = reactive<createUserRequest>({
  Username: '',
  Password: '',
  Admin: false,
});

const handleFinish: FormProps['onFinish'] = () => {
  loading.value = true;
  adminCreateUser(formState)
    .then((data) => {
      if ((data as errorResponse).message) {
        message.error((data as errorResponse).message);
      }
    })
    .catch((data) => {
      message.error((data as errorResponse).message);
    })
    .finally(() => {
      loadUserList();
      loading.value = false;
      formState.Username = '';
      formState.Password = '';
      formState.Admin = false;
    });
};
const handleFinishFailed: FormProps['onFinishFailed'] = (errors) => {
  console.log(errors);
};

loadUserList();
</script>
