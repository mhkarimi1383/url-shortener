<script setup lang="ts">
import dayjs from 'dayjs';
import { ref, reactive } from 'vue';
import { message } from 'ant-design-vue';
import { pageToOffset } from '@/lib/utils';
import type { FormProps } from 'ant-design-vue';
import relativeTime from 'dayjs/plugin/relativeTime';
import { NumberOutlined } from '@ant-design/icons-vue';
import {
  listEntities,
  createEntity,
  deleteEntity,
  type errorResponse,
  type entityCreateRequest,
  type listEntitiesResponse,
} from '@/lib/api';

dayjs.extend(relativeTime);

const limit = ref(5);
const page = ref(1);
const loading = ref<boolean>(true);
const resp = ref<listEntitiesResponse>({
  Result: [],
  MetaData: {
    Count: NaN,
    TotalVisit: NaN,
  },
});

const loadEntities = function () {
  listEntities(limit.value, pageToOffset(page.value, limit.value))
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

const handleFinish: FormProps['onFinish'] = () => {
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
    key: 'visitCount',
    title: 'VisitCount',
    dataIndex: 'VisitCount',
  },
  {
    key: 'lastVisitedAt',
    title: 'LastVisitedAt',
    dataIndex: 'LastVisitedAt',
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
      <a-col :span="6">
        <a-card title="Statistics" style="height: 100%">
          <a-card-grid style="width: 50%"
            ><a-statistic title="Number of Entities" :value="resp?.MetaData.Count"
          /></a-card-grid>
          <a-card-grid style="width: 50%"
            ><a-statistic title="Total Visit Count" :value="resp?.MetaData.TotalVisit"
          /></a-card-grid>
        </a-card>
      </a-col>
      <a-col :span="18">
        <a-card title="Create Entity">
          <a-form
            layout="inline"
            :model="formState"
            @finish="handleFinish"
            @finish-failed="handleFinishFailed"
          >
            <a-card-grid :bordered="false" style="width: 45%">
              <a-form-item>
                <a-input v-model:value="formState.Name" placeholder="Name"> </a-input>
              </a-form-item>
            </a-card-grid>
            <a-card-grid :bordered="false" style="width: 45%">
              <a-form-item>
                <a-input v-model:value="formState.Description" placeholder="Description"> </a-input>
              </a-form-item>
            </a-card-grid>
            <a-card-grid :bordered="false" style="width: 10%">
              <a-form-item>
                <a-button
                  type="primary"
                  html-type="submit"
                  :disabled="formState.Name === '' || formState.Description === ''"
                >
                  Create
                </a-button>
              </a-form-item>
            </a-card-grid>
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
      <template #bodyCell="{ column, record, text }">
        <template v-if="['updatedAt', 'createdAt', 'lastVisitedAt'].includes(column.key)">
          {{ (text && dayjs(text).toNow(false)) || 'NaN' }}
        </template>
        <template v-else-if="column.key === 'actions'">
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
      v-model:current="page"
      v-model:page-size="limit"
      style="float: right"
      :total="resp?.MetaData.Count"
      @change="loadEntities"
    ></a-pagination>
  </a-spin>
</template>
