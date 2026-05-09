<script setup lang="ts">
import { onMounted, ref } from 'vue'
import * as shareApi from '../../api/shares'
import type { ListShare, ListShareMember } from '../../types/todo'
import { useToast } from '../../stores/toast'

const props = defineProps<{ listId: number; listName: string }>()
const emit = defineEmits<{ close: [] }>()

const share = ref<ListShare | null>(null)
const members = ref<ListShareMember[]>([])
const permission = ref<'view' | 'edit'>('edit')
const loading = ref(false)

onMounted(async () => {
  try {
    const res = await shareApi.getShare(props.listId)
    share.value = res.data.share
    members.value = res.data.members
  } catch {
    // not shared yet
  }
})

async function handleCreate() {
  loading.value = true
  try {
    const res = await shareApi.createShare(props.listId, permission.value)
    share.value = res.data
    useToast().show('共享链接已生成', 'success')
  } catch {
    useToast().show('创建共享失败', 'error')
  } finally {
    loading.value = false
  }
}

async function handleRevoke() {
  try {
    await shareApi.deleteShare(props.listId)
    share.value = null
    members.value = []
    useToast().show('共享已撤销', 'success')
  } catch {
    useToast().show('撤销失败', 'error')
  }
}

function copyLink() {
  if (!share.value) return
  const link = `${window.location.origin}/join?code=${share.value.code}`
  navigator.clipboard.writeText(link)
  useToast().show('链接已复制', 'success')
}
</script>

<template>
  <Teleport to="body">
  <div class="fixed inset-0 z-50 flex items-center justify-center bg-black/30" @click.self="emit('close')">
    <div class="bg-white dark:bg-gray-800 rounded-2xl shadow-xl p-6 w-[380px] max-w-[90vw]">
      <h3 class="text-base font-bold text-gray-800 dark:text-gray-100 mb-4">共享列表「{{ listName }}」</h3>

      <!-- Not shared yet: create -->
      <div v-if="!share" class="space-y-3">
        <div class="flex items-center gap-2">
          <label class="text-xs text-gray-500">权限</label>
          <select v-model="permission" class="px-2.5 py-1.5 border border-gray-200 dark:border-gray-700 rounded-lg text-sm bg-white dark:bg-gray-800 dark:text-gray-100 outline-none">
            <option value="view">仅查看</option>
            <option value="edit">可编辑</option>
          </select>
        </div>
        <button @click="handleCreate" :disabled="loading"
          class="w-full py-2 bg-indigo-500 text-white rounded-lg text-sm font-semibold hover:bg-indigo-600 disabled:opacity-40 transition-colors">
          {{ loading ? '创建中…' : '生成共享链接' }}
        </button>
      </div>

      <!-- Shared: show code & members -->
      <div v-else class="space-y-4">
        <div class="flex items-center gap-2">
          <span class="text-xs text-gray-400">邀请码</span>
          <code class="px-3 py-1.5 bg-gray-100 dark:bg-gray-700 rounded-lg text-lg font-mono font-bold text-indigo-600 dark:text-indigo-400 tracking-wider select-all">{{ share.code }}</code>
        </div>
        <button @click="copyLink"
          class="w-full py-2 border border-indigo-200 dark:border-indigo-800 text-indigo-600 dark:text-indigo-400 rounded-lg text-sm font-medium hover:bg-indigo-50 dark:hover:bg-indigo-900/30 transition-colors">
          复制分享链接
        </button>
        <div>
          <div class="text-xs text-gray-400 mb-2">成员 ({{ members.length }})</div>
          <div v-if="members.length === 0" class="text-xs text-gray-300 dark:text-gray-600">暂无成员</div>
          <div v-for="m in members" :key="m.id" class="text-sm text-gray-600 dark:text-gray-400 py-1">
            ID: {{ m.user_id }}
          </div>
        </div>
        <button @click="handleRevoke"
          class="w-full py-2 text-red-500 rounded-lg text-sm font-medium hover:bg-red-50 dark:hover:bg-red-900/20 transition-colors">
          撤销共享
        </button>
      </div>

      <button @click="emit('close')"
        class="w-full py-2 mt-3 text-gray-400 rounded-lg text-sm hover:text-gray-600 dark:hover:text-gray-300 transition-colors">
        关闭
      </button>
    </div>
  </div>
  </Teleport>
</template>
