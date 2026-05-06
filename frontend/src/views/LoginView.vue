<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuth } from '../stores/auth'
import { useTodos } from '../stores/todo'
import ConfirmDialog from '../components/ui/ConfirmDialog.vue'

const auth = useAuth()
const todoStore = useTodos()
const router = useRouter()

const username = ref('')
const password = ref('')
const error = ref('')
const loading = ref(false)
const mode = ref<'bind' | 'login'>(auth.hasPassword ? 'login' : 'bind')
const showMergeConfirm = ref(false)

// Check if anonymous user has data
let anonymousCount = 0
if (mode.value === 'login') {
  todoStore.fetchTodos().then(() => {
    anonymousCount = todoStore.todos.length
  })
}

async function submit() {
  if (!username.value.trim() || !password.value.trim()) {
    error.value = '请填写用户名和密码'
    return
  }

  if (mode.value === 'login' && anonymousCount > 0) {
    showMergeConfirm.value = true
    return
  }
  await doSubmit(false)
}

async function doSubmit(merge: boolean) {
  loading.value = true
  error.value = ''
  try {
    if (mode.value === 'bind') {
      await auth.bind(username.value.trim(), password.value)
    } else if (merge) {
      await auth.mergeLogin(username.value.trim(), password.value)
    } else {
      await auth.login(username.value.trim(), password.value)
    }
    router.push('/')
  } catch (e: any) {
    error.value = e.response?.data?.error || '操作失败'
  } finally {
    loading.value = false
    showMergeConfirm.value = false
  }
}

function handleKeydown(e: KeyboardEvent) {
  if (e.key === 'Enter') submit()
}
</script>

<template>
  <div class="min-h-[80vh] flex items-center justify-center">
    <div class="w-full max-w-sm bg-white dark:bg-gray-900 rounded-2xl border border-gray-200 dark:border-gray-800 shadow-sm p-8">
      <h2 class="text-xl font-bold text-gray-800 dark:text-gray-100 mb-2 text-center">
        {{ mode === 'bind' ? '设置账号' : '登录' }}
      </h2>
      <p class="text-sm text-gray-400 dark:text-gray-500 text-center mb-6">
        {{ mode === 'bind' ? '设置用户名和密码，永久保存你的数据' : '用已有的账号登录，恢复你的数据' }}
      </p>

      <div class="space-y-3">
        <div>
          <input
            v-model="username"
            type="text"
            placeholder="用户名"
            @keydown="handleKeydown"
            class="w-full px-3 py-2.5 border border-gray-200 dark:border-gray-800 rounded-xl text-sm bg-white dark:bg-gray-900 dark:text-gray-100 outline-none focus:border-indigo-400 transition-colors"
          />
        </div>
        <div>
          <input
            v-model="password"
            type="password"
            placeholder="密码"
            @keydown="handleKeydown"
            class="w-full px-3 py-2.5 border border-gray-200 dark:border-gray-800 rounded-xl text-sm bg-white dark:bg-gray-900 dark:text-gray-100 outline-none focus:border-indigo-400 transition-colors"
          />
        </div>
        <div v-if="error" class="text-xs text-red-500">{{ error }}</div>
        <button
          @click="submit"
          :disabled="loading || !username.trim() || !password.trim()"
          class="w-full py-2.5 bg-indigo-500 text-white rounded-xl text-sm font-semibold hover:bg-indigo-600 disabled:opacity-40 disabled:cursor-default transition-colors"
        >
          {{ loading ? '处理中…' : (mode === 'bind' ? '设置账号' : '登录') }}
        </button>
      </div>

      <div class="mt-4 text-center">
        <button @click="mode = mode === 'bind' ? 'login' : 'bind'" class="text-xs text-gray-400 hover:text-indigo-500 transition-colors">
          {{ mode === 'bind' ? '已有账号？去登录' : '没有账号？设置一个' }}
        </button>
      </div>
    </div>
  </div>

  <ConfirmDialog
    :open="showMergeConfirm"
    title="合并数据"
    :message="`当前有 ${anonymousCount} 条匿名待办数据。是否将这些数据合并到账号「${username}」中？`"
    confirm-text="合并"
    cancel-text="不合并"
    @confirm="doSubmit(true)"
    @cancel="doSubmit(false)"
  />
</template>
