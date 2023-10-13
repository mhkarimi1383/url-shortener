<script setup lang="ts">
import { ref, reactive } from 'vue';
import { message } from 'ant-design-vue';
import type { FormProps } from 'ant-design-vue';
import { NumberOutlined } from '@ant-design/icons-vue';
import {
  listUrls,
  createUrl,
  deleteUrl,
  listEntities,
  type errorResponse,
  type urlCreateRequest,
  type listUrlsResponse,
  type urlCreateResponse,
  type listEntitiesResponse,
} from '@/lib/api';

const limit = ref(5);
const offset = ref(1);
const loading = ref<boolean>(true);
const resp = ref<listUrlsResponse>({
  Result: [],
  MetaData: {
    Count: NaN,
  },
});
const createdModalView = ref(false);
const createResult = ref<urlCreateResponse>();

const closeResultModal = function () {
  createdModalView.value = false;
};
const loadListUrls = function () {
  listUrls(limit.value, offset.value - 1)
    .then((data) => {
      if ((data as errorResponse).message) {
        message.error((data as errorResponse).message);
      }
      resp.value = data as listUrlsResponse;
    })
    .catch((data) => {
      message.error((data as errorResponse).message);
    })
    .finally(() => {
      loading.value = false;
    });
};

async function deleteUrlTable(id: number) {
  await deleteUrl(id);
  loadListUrls();
}

const formState = reactive<urlCreateRequest>({
  Entity: 0,
  FullUrl: '',
});

const handleFinish: FormProps['onFinish'] = (_) => {
  loading.value = true;
  createUrl(formState)
    .then((data) => {
      if ((data as errorResponse).message) {
        message.error((data as errorResponse).message);
      } else {
        createResult.value = data as urlCreateResponse;
        createdModalView.value = true;
      }
    })
    .catch((data) => {
      message.error((data as errorResponse).message);
    })
    .finally(() => {
      loadListUrls();
      formState.Entity = 0;
      loading.value = false;
      formState.FullUrl = '';
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
    key: 'fullUrl',
    title: 'Full Url',
    dataIndex: 'FullUrl',
  },
  {
    key: 'shortCode',
    title: 'ShortCode',
    dataIndex: 'ShortCode',
  },
  {
    key: 'shortUrl',
    title: 'ShortUrl',
    dataIndex: 'ShortUrl',
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
    key: 'entity',
    title: 'Entity',
    dataIndex: ['Entity', 'Name'],
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
const entities = ref<listEntitiesResponse>();
listEntities(15, 0)
  .then((data) => {
    if ((data as errorResponse).message) {
      message.error((data as errorResponse).message);
    }
    entities.value = data as listEntitiesResponse;
  })
  .catch((data) => {
    message.error((data as errorResponse).message);
  })
  .finally(() => {
    loading.value = false;
  });
loadListUrls();
</script>

<template>
  <a-spin :spinning="loading">
    <a-row>
      <a-col :span="12">
        <a-card>
          <a-statistic title="Number of URL(s)" :value="resp?.MetaData.Count" />
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
              <a-select ref="select" v-model:value="formState.Entity" style="width: 20vh">
                <a-select-option
                  v-for="entity in entities?.Result"
                  v-bind:key="entity.Id"
                  :value="entity.Id"
                  >{{ entity.Id }} - {{ entity.Name }}</a-select-option
                >
              </a-select>
            </a-form-item>
            <a-form-item>
              <a-input v-model:value="formState.FullUrl" placeholder="Long URL"> </a-input>
            </a-form-item>
            <a-form-item>
              <a-button type="primary" html-type="submit" :disabled="formState.FullUrl === ''">
                Create
              </a-button>
            </a-form-item>
          </a-form>
          <a-modal v-model:open="createdModalView" title="Create Result">
            <a-result @ok="closeResultModal" status="success" title="Successfully Created URL">
              <template #subTitle>
                ShortURL:
                <a :href="createResult?.ShortUrl" target="_blank">{{ createResult?.ShortUrl }}</a>
                <br />
                ShortCode: {{ createResult?.ShortCode }}
              </template>
            </a-result>
            <template #footer>
              <a-button key="back" @click="closeResultModal">Return</a-button>
            </template>
          </a-modal>
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
            title="Are you sure delete this URL?"
            ok-text="Yes"
            cancel-text="No"
            @confirm="deleteUrlTable(record.Id)"
          >
            <a-button type="primary" :danger="true">
              <template #icon>
                <KeyOutlined />
              </template>
              Delete
            </a-button>
          </a-popconfirm>
        </template>
        <template v-else-if="column.key === 'shortUrl'">
          <a :href="record.ShortUrl" target="_blank">{{ record.ShortUrl }}</a>
        </template>
        <template v-else-if="column.key === 'fullUrl'">
          <a :href="record.FullUrl" target="_blank">{{ record.FullUrl }}</a>
        </template>
      </template>
    </a-table>
    <br />
    <a-pagination
      style="float: right"
      :total="resp?.MetaData.Count"
      v-model:current="offset"
      v-model:pageSize="limit"
      @change="loadListUrls"
    ></a-pagination>
  </a-spin>
</template>
