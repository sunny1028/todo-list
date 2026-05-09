<script setup lang="ts">
import { ref } from 'vue'
import * as shareApi from '../../api/shares'
import { useLists } from '../../stores/lists'
import { useToast } from '../../stores/toast'

const emit = defineEmits<{ close: [] }>()

const code = ref('')
const loading = ref(false)
const listStore = useLists()

async function handleJoin() {
  if (!code.value.trim()) return
  loading.value = true
  try {
    await shareApi.joinList(code.value.trim().toUpperCase())
    await listStore.fetchLists()
    useToast().show('已加入共享列表', 'success')
    emit('close')
  } catch {
    useToast().show('加入失败，请检查邀请码', 'error')
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <Teleport to="body">
  <div class="fixed inset-0 z-50 flex items-center justify-center bg-black/30" @click.self="emit('close')">
    <div class="bg-white dark:bg-gray-800 rounded-2xl shadow-xl p-6 w-[360px] max-w-[90vw]">
      <h3 class="text-base font-bold text-gray-800 dark:text-gray-100 mb-4">加入共享列表</h3>
      <input v-model="code" @keydown.enter="handleJoin"
        placeholder="输入邀请码"
        class="w-full px-3 py-2.5 border border-gray-200 dark:border-gray-700 rounded-lg text-sm bg-white dark:bg-gray-800 dark:text-gray-100 outline-none focus:border-indigo-400 mb-3 uppercase"
        maxlength="8" />
      <button @click="handleJoin" :disabled="loading || !code.trim()"
        class="w-full py-2 bg-indigo-500 text-white rounded-lg text-sm font-semibold hover:bg-indigo-600 disabled:opacity-40 transition-colors mb-2">
        {{ loading ? '加入中…' : '加入' }}
      </button>
      <button @click="emit('close')"
        class="w-full py-2 text-gray-400 rounded-lg text-sm hover:text-gray-600 dark:hover:text-gray-300 transition-colors">
        取消
      </button>
    </div>
  </div>
  </Teleport>
</template>
