<script setup lang="ts">
import dayjs from 'dayjs';
import { ref, reactive } from 'vue';
import { message } from 'ant-design-vue';
import { pageToOffset } from '@/lib/utils';
import type { FormProps } from 'ant-design-vue';
import relativeTime from 'dayjs/plugin/relativeTime';
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

dayjs.extend(relativeTime);

const limit = ref(5);
const page = ref(1);
const loading = ref<boolean>(true);
const resp = ref<listUrlsResponse>({
  Result: [],
  MetaData: {
    Count: NaN,
    TotalVisit: NaN,
  },
});
const createdModalView = ref(false);
const createResult = ref<urlCreateResponse>();

const closeResultModal = function () {
  createdModalView.value = false;
};
const loadListUrls = function () {
  listUrls(limit.value, pageToOffset(page.value, limit.value))
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

const handleFinish: FormProps['onFinish'] = () => {
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
      <a-col :span="6">
        <a-card title="Statistics" style="height: 100%">
          <a-card-grid style="width: 50%"
            ><a-statistic title="Number of URL(s)" :value="resp?.MetaData.Count"
          /></a-card-grid>
          <a-card-grid style="width: 50%"
            ><a-statistic title="Total Visit Count" :value="resp?.MetaData.TotalVisit"
          /></a-card-grid>
        </a-card>
      </a-col>
      <a-col :span="18">
        <a-card title="Create URL">
          <a-form
            layout="inline"
            :model="formState"
            @finish="handleFinish"
            @finish-failed="handleFinishFailed"
          >
            <a-card-grid :bordered="false" style="width: 25%">
              <a-form-item>
                <a-select ref="select" v-model:value="formState.Entity">
                  <a-select-option
                    v-for="entity in entities?.Result"
                    :key="entity.Id"
                    :value="entity.Id"
                    >{{ entity.Id }} - {{ entity.Name }}</a-select-option
                  >
                </a-select>
              </a-form-item>
            </a-card-grid>
            <a-card-grid :bordered="false" style="width: 40%">
              <a-form-item :rules="[{ required: true, message: 'Please enter full URL!' }]">
                <a-input v-model:value="formState.FullUrl" placeholder="Long URL"> </a-input>
              </a-form-item>
            </a-card-grid>
            <a-card-grid :bordered="false" style="width: 25%">
              <a-form-item>
                <a-input
                  v-model:value="formState.ShortCode"
                  placeholder="Custom Short Code (optional)"
                >
                </a-input>
              </a-form-item>
            </a-card-grid>
            <a-card-grid style="width: 10%" :bordered="false">
              <a-form-item>
                <a-button type="primary" html-type="submit" :disabled="formState.FullUrl === ''">
                  Create
                </a-button>
              </a-form-item>
            </a-card-grid>
          </a-form>
          <a-modal v-model:open="createdModalView" title="Create Result">
            <a-result status="success" title="Successfully Created URL" @ok="closeResultModal">
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
      <template #bodyCell="{ column, record, text }">
        <template v-if="['updatedAt', 'createdAt', 'lastVisitedAt'].includes(column.key)">
          {{ (text && dayjs(text).toNow(false)) || 'NaN' }}
        </template>
        <template v-else-if="column.key === 'actions'">
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
      v-model:current="page"
      v-model:page-size="limit"
      style="float: right"
      :total="resp?.MetaData.Count"
      @change="loadListUrls"
    ></a-pagination>
  </a-spin>
</template>
