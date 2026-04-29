<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { useTodos } from '../stores/todo'
import TodoForm from './TodoForm.vue'
import TodoItem from './TodoItem.vue'
import Skeleton from './ui/Skeleton.vue'
import ConfirmDialog from './ui/ConfirmDialog.vue'

const store = useTodos()
const route = useRoute()

const search = ref('')
const priorityFilter = ref('')
const tagFilter = ref('')
const sortBy = ref('created')
const showClearConfirm = ref(false)
const showMore = ref(false)
const searchInput = ref<HTMLInputElement | null>(null)

const filterStatus = computed(() => (route.query.filter as string) || '')

function filteredTodos() {
  let list = [...store.todos]

  if (filterStatus.value === 'active') list = list.filter((t) => !t.completed && !t.archived)
  else if (filterStatus.value === 'completed') list = list.filter((t) => t.completed)
  else if (filterStatus.value === 'archived') list = list.filter((t) => t.archived)
  else list = list.filter((t) => !t.archived)

  if (priorityFilter.value) list = list.filter((t) => t.priority === priorityFilter.value)
  if (tagFilter.value) list = list.filter((t) => t.tags.includes(tagFilter.value))
  if (search.value.trim()) {
    const q = search.value.trim().toLowerCase()
    list = list.filter((t) => t.title.toLowerCase().includes(q) || t.description.toLowerCase().includes(q))
  }

  switch (sortBy.value) {
    case 'priority': {
      const order = { high: 0, medium: 1, low: 2 }
      list.sort((a, b) => (order[a.priority] ?? 1) - (order[b.priority] ?? 1))
      break
    }
    case 'created':
      list.sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime())
      break
    case 'due':
      list.sort((a, b) => {
        if (!a.due_date) return 1
        if (!b.due_date) return -1
        return new Date(a.due_date).getTime() - new Date(b.due_date).getTime()
      })
      break
  }

  return list
}

const completedCount = computed(() => store.todos.filter((t) => t.completed).length)
const archivedCount = computed(() => store.todos.filter((t) => t.archived).length)

const allTags = computed(() => {
  const set = new Set<string>()
  store.todos.forEach((t) => {
    if (t.tags) t.tags.split(',').filter(Boolean).forEach((tag) => set.add(tag.trim()))
  })
  return [...set].sort()
})

async function handleClearCompleted() {
  await store.clearCompleted()
  showClearConfirm.value = false
  showMore.value = false
}

function onKeydown(e: KeyboardEvent) {
  if ((e.ctrlKey || e.metaKey) && e.key === 'k') {
    e.preventDefault()
    searchInput.value?.focus()
  }
  if (e.key === 'Escape') {
    search.value = ''
    priorityFilter.value = ''
    tagFilter.value = ''
    showMore.value = false
    document.activeElement instanceof HTMLElement && document.activeElement.blur()
  }
}

onMounted(() => {
  store.fetchTodos()
  window.addEventListener('keydown', onKeydown)
})

onUnmounted(() => {
  window.removeEventListener('keydown', onKeydown)
})
</script>

<template>
  <div>
    <TodoForm @created="() => {}" />

    <div class="flex gap-2 mb-3 flex-wrap">
      <input
        ref="searchInput"
        v-model="search"
        type="text"
        placeholder="搜索… (Ctrl+K)"
        class="flex-1 min-w-[120px] px-3 py-2 border border-gray-200 dark:border-gray-800 rounded-lg text-sm outline-none focus:border-indigo-400 bg-white dark:bg-gray-900 dark:text-gray-100 transition-colors"
      />
      <select v-model="priorityFilter" class="px-3 py-2 border border-gray-200 dark:border-gray-800 rounded-lg text-sm bg-white dark:bg-gray-900 dark:text-gray-100 outline-none">
        <option value="">全部优先级</option>
        <option value="high">高</option>
        <option value="medium">中</option>
        <option value="low">低</option>
      </select>
      <select v-model="tagFilter" class="px-3 py-2 border border-gray-200 dark:border-gray-800 rounded-lg text-sm bg-white dark:bg-gray-900 dark:text-gray-100 outline-none">
        <option value="">全部标签</option>
        <option v-for="tag in allTags" :key="tag" :value="tag">{{ tag }}</option>
      </select>
      <select v-model="sortBy" class="px-3 py-2 border border-gray-200 dark:border-gray-800 rounded-lg text-sm bg-white dark:bg-gray-900 dark:text-gray-100 outline-none">
        <option value="created">最新</option>
        <option value="priority">优先级</option>
        <option value="due">截止日期</option>
      </select>

      <!-- More menu -->
      <div class="relative">
        <button @click="showMore = !showMore" class="px-3 py-2 border border-gray-200 dark:border-gray-800 rounded-lg text-sm bg-white dark:bg-gray-900 text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300 transition-colors">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="5" r="1"/><circle cx="12" cy="12" r="1"/><circle cx="12" cy="19" r="1"/></svg>
        </button>
        <div v-if="showMore" class="absolute right-0 top-full mt-1 bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700 shadow-lg py-1 z-10 min-w-[130px]">
          <button @click="store.downloadExport('json'); showMore = false" class="w-full text-left px-4 py-2 text-sm text-gray-600 dark:text-gray-400 hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors">导出 JSON</button>
          <button @click="store.downloadExport('csv'); showMore = false" class="w-full text-left px-4 py-2 text-sm text-gray-600 dark:text-gray-400 hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors">导出 CSV</button>
          <button v-if="completedCount > 0" @click="showClearConfirm = true" class="w-full text-left px-4 py-2 text-sm text-red-500 hover:bg-red-50 dark:hover:bg-red-900/30 transition-colors">清除已完成</button>
        </div>
      </div>
    </div>

    <div v-if="store.todos.length > 0" class="flex items-center justify-between mb-3 text-xs text-gray-400 dark:text-gray-500">
      <span>
        {{ store.todos.filter(t => !t.completed && !t.archived).length }} 项未完成 / 共 {{ store.todos.filter(t => !t.archived).length }} 项
        <span v-if="archivedCount > 0" class="text-gray-300 dark:text-gray-600"> · {{ archivedCount }} 归档</span>
      </span>
    </div>

    <div v-if="store.loading"><Skeleton :count="5" /></div>
    <div v-else-if="store.error" class="text-center py-10 text-red-400 text-sm">{{ store.error }}</div>
    <div v-else-if="filteredTodos().length === 0" class="text-center py-16 text-gray-300 dark:text-gray-700 text-sm italic">
      暂无待办
    </div>

    <TransitionGroup v-else name="list" tag="div">
      <div v-for="todo in filteredTodos()" :key="todo.id">
        <TodoItem :todo="todo" />
      </div>
    </TransitionGroup>

    <div class="mt-8 text-center text-[11px] text-gray-300 dark:text-gray-700 hidden md:block">
      <kbd class="px-1 py-0.5 rounded bg-gray-100 dark:bg-gray-800 text-gray-400">Ctrl+N</kbd> 新建
      <kbd class="px-1 py-0.5 rounded bg-gray-100 dark:bg-gray-800 text-gray-400 ml-2">Ctrl+K</kbd> 搜索
      <kbd class="px-1 py-0.5 rounded bg-gray-100 dark:bg-gray-800 text-gray-400 ml-2">?</kbd> 快捷键
    </div>
  </div>

  <ConfirmDialog
    :open="showClearConfirm"
    title="清除已完成"
    :message="`确定要删除全部 ${completedCount} 个已完成项吗？`"
    @confirm="handleClearCompleted"
    @cancel="showClearConfirm = false"
  />
</template>

<style scoped>
.list-enter-active, .list-leave-active { transition: all 0.3s ease; }
.list-enter-from { opacity: 0; transform: translateY(-12px); }
.list-leave-to { opacity: 0; transform: translateX(20px); }
</style>
