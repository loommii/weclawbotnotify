<template>
  <div>
    <div
      v-if="isError"
      class="flex flex-col items-center justify-center py-12"
    >
      <i class="pi pi-exclamation-circle mb-3 text-4xl text-muted-color"></i>
      <p class="mb-4 text-sm text-muted-color">
        加载失败，点击重试
      </p>
      <Button
        label="重试"
        icon="pi pi-refresh"
        @click="$emit('retry')"
      />
    </div>

    <template v-else>
      <DataTable
        :value="displayData"
        :loading="isRefreshing"
        striped-rows
        responsive-layout="scroll"
        class="text-sm"
      >
        <template #empty>
          <div class="flex flex-col items-center justify-center py-12">
            <i class="pi pi-key mb-3 text-4xl text-muted-color"></i>
            <p class="mb-1 text-base font-semibold text-color">
              暂无 Application
            </p>
            <p class="mb-4 text-sm text-muted-color">
              创建第一个推送凭证开始使用
            </p>
            <Button
              label="创建 Application"
              icon="pi pi-plus"
              @click="$emit('create')"
            />
          </div>
        </template>

        <Column field="name" class="min-w-[160px]">
          <template #header>
            <span class="text-[13px] font-medium uppercase text-muted-color">
              应用名称
            </span>
          </template>
          <template #body="{ data }">
            <Skeleton v-if="data.skeleton" class="my-2" width="60%" />
            <span v-else class="font-medium text-color">{{ data.name }}</span>
          </template>
        </Column>

        <Column field="description" class="min-w-[200px]">
          <template #header>
            <span class="text-[13px] font-medium uppercase text-muted-color">
              描述
            </span>
          </template>
          <template #body="{ data }">
            <Skeleton v-if="data.skeleton" class="my-2" width="80%" />
            <span v-else class="text-muted-color">
              {{ data.description || '—' }}
            </span>
          </template>
        </Column>

        <Column field="createdAt" class="min-w-[160px]">
          <template #header>
            <span class="text-[13px] font-medium uppercase text-muted-color">
              创建时间
            </span>
          </template>
          <template #body="{ data }">
            <Skeleton v-if="data.skeleton" class="my-2" width="50%" />
            <span v-else class="text-xs text-muted-color">
              {{ formatDate(data.createdAt) }}
            </span>
          </template>
        </Column>

        <Column field="lastUsedAt" class="min-w-[160px]">
          <template #header>
            <span class="text-[13px] font-medium uppercase text-muted-color">
              最后使用
            </span>
          </template>
          <template #body="{ data }">
            <Skeleton v-if="data.skeleton" class="my-2" width="50%" />
            <span v-else class="text-xs text-muted-color">
              {{ data.lastUsedAt ? formatDate(data.lastUsedAt) : '从未使用' }}
            </span>
          </template>
        </Column>

        <Column class="min-w-[80px]" style="text-align: right">
          <template #header>
            <span class="text-[13px] font-medium uppercase text-muted-color">
              操作
            </span>
          </template>
          <template #body="{ data }">
            <div v-if="!data.skeleton" class="flex justify-end">
              <Button
                label="删除"
                text
                severity="danger"
                size="small"
                @click="$emit('delete', data)"
              />
            </div>
          </template>
        </Column>
      </DataTable>

      <Paginator
        v-if="totalRecords > 0"
        :first="first"
        :rows="rows"
        :total-records="totalRecords"
        :rows-per-page-options="[10, 20, 50]"
        class="mt-4"
        @page="$emit('page-change', $event)"
      />
    </template>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import Skeleton from 'primevue/skeleton'
import Paginator from 'primevue/paginator'
import type { ApplicationInfo } from '@/types/application'

const SKELETON_COUNT = 5

interface SkeletonRow extends Partial<ApplicationInfo> {
  id: number
  skeleton: true
}

const props = defineProps<{
  applications: ApplicationInfo[]
  isFirstLoad: boolean
  isRefreshing: boolean
  isError: boolean
  first: number
  rows: number
  totalRecords: number
}>()

defineEmits<{
  create: []
  delete: [app: ApplicationInfo]
  retry: []
  'page-change': [event: { page: number; rows: number }]
}>()

const skeletonRows: SkeletonRow[] = Array.from({ length: SKELETON_COUNT }, (_, i) => ({
  id: -(i + 1),
  skeleton: true,
  name: '',
  description: '',
  createdAt: '',
  lastUsedAt: null,
}))

const displayData = computed(() => {
  if (props.isFirstLoad) return skeletonRows
  return props.applications
})

function formatDate(dateStr: string): string {
  const ts = parseInt(dateStr, 10)
  if (isNaN(ts) || ts === 0) return '—'
  const date = new Date(ts * 1000)
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  return `${year}-${month}-${day} ${hours}:${minutes}`
}
</script>
