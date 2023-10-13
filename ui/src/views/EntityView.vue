<script setup lang="ts">
import { ref, reactive } from 'vue';
import { message } from 'ant-design-vue';
import type { FormProps } from 'ant-design-vue';
import { NumberOutlined } from '@ant-design/icons-vue';
import {
  listEntities,
  createEntity,
  deleteEntity,
  type errorResponse,
  type entityCreateRequest,
  type listEntitiesResponse,
} from '@/lib/api';

const limit = ref(5);
const offset = ref(1);
const loading = ref<boolean>(true);
const resp = ref<listEntitiesResponse>({
  Result: [],
  MetaData: {
    Count: NaN,
  },
});

const loadEntities = function () {
  listEntities(limit.value, offset.value - 1)
    .then((data) => {
      if ((data as errorResponse).message) {
        message.error((data as errorResponse).message);
      }
      resp.value = data as listEntitiesResponse;
    })
    .catch((data) => {
      message.error((data as errorResponse).message);
    })
    .finally(() => {
      loading.value = false;
    });
};

async function deleteEntityTable(id: number) {
  await deleteEntity(id);
  loadEntities();
}

const formState = reactive<entityCreateRequest>({
  Name: '',
  Description: '',
});

const handleFinish: FormProps['onFinish'] = (_) => {
  loading.value = true;
  createEntity(formState)
    .then((data) => {
      if ((data as errorResponse).message) {
        message.error((data as errorResponse).message);
      }
    })
    .catch((data) => {
      message.error((data as errorResponse).message);
    })
    .finally(() => {
      loadEntities();
      loading.value = false;
      formState.Description = '';
      formState.Name = '';
    });
};
const handleFinishFailed: FormProps['onFinishFailed'] = (errors) => {
  console.log(errors);
};

const columns = [
  {
    key: 'id',
    name: 'Id',
    dataIndex: 'Id',
  },
  {
    key: 'name',
    title: 'Name',
    dataIndex: 'Name',
  },
  {
    key: 'description',
    title: 'Description',
    dataIndex: 'Description',
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
    key: 'creator',
    title: 'Creator',
    dataIndex: ['Creator', 'Username'],
  },
  {
    key: 'actions',
    title: 'Actions',
  },
];

loadEntities();
</script>

<template>
  <a-spin :spinning="loading">
    <a-row>
      <a-col :span="12">
        <a-card>
          <a-statistic title="Number of Entities" :value="resp?.MetaData.Count" />
        </a-card>
      </a-col>
      <a-col :span="12">
        <a-card style="width: 100%; height: 100%; display: flex" bodyStyle="align-self: center;">
          <a-form
            layout="inline"
            :model="formState"
            @finish="handleFinish"
            @finishFailed="handleFinishFailed"
          >
            <a-form-item>
              <a-input v-model:value="formState.Name" placeholder="Name"> </a-input>
            </a-form-item>
            <a-form-item>
              <a-input v-model:value="formState.Description" placeholder="Description"> </a-input>
            </a-form-item>
            <a-form-item>
              <a-button
                type="primary"
                html-type="submit"
                :disabled="formState.Name === '' || formState.Description === ''"
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
      </template>
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'actions'">
          <a-popconfirm
            title="Are you sure delete this entity (It will delete every link that refers to this entity)?"
            ok-text="Yes"
            cancel-text="No"
            @confirm="deleteEntityTable(record.Id)"
          >
            <a-button type="primary" :danger="true">
              <template #icon>
                <KeyOutlined />
              </template>
              Delete
            </a-button>
          </a-popconfirm>
        </template>
      </template>
    </a-table>
    <br />
    <a-pagination
      style="float: right"
      :total="resp?.MetaData.Count"
      v-model:current="offset"
      v-model:pageSize="limit"
      @change="loadEntities"
    ></a-pagination>
  </a-spin>
</template>
