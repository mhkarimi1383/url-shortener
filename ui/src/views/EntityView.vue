<script setup lang="ts">
import { ref } from 'vue';
import { listEntities } from '@/lib/api';
import { message } from 'ant-design-vue';
import type { listEntitiesResponse, errorResponse } from '@/lib/api';

const limit = ref(5);
const offset = ref(1);
const loading = ref<boolean>(true);
const resp = ref<listEntitiesResponse>({
  Result: [],
  MetaData: {
    Count: NaN,
  },
});

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
</script>

<template>
  <a-spin :spinning="loading">
    <a-row>
      <a-col :span="12">
        <a-card>
          <a-statistic title="Number of Entities" :value="resp?.MetaData.Count" />
        </a-card>
      </a-col>
    </a-row>
    <br />
  </a-spin>
</template>
