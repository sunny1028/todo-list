<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { useTodos } from '../stores/todo'
import { useLists } from '../stores/lists'
import { useAuth } from '../stores/auth'
import type { Todo } from '../types/todo'
import TodoForm from './TodoForm.vue'
import TodoItem from './TodoItem.vue'
import Skeleton from './ui/Skeleton.vue'
import Select from './ui/Select.vue'
import ConfirmDialog from './ui/ConfirmDialog.vue'
import { useToast } from '../stores/toast'
import * as api from '../api/todos'

const store = useTodos()
const route = useRoute()
const listStore = useLists()
const authStore = useAuth()

const isReadonly = computed(() => {
  const list = listStore.lists.find(l => l.id === store.currentListId)
  if (!list) return false
  return list.user_id !== authStore.userId && list.permission === 'view'
})

const search = ref('')
const priorityFilter = ref('')
const tagFilter = ref('')
const sortBy = ref('default')
const showClearConfirm = ref(false)
const showMore = ref(false)
const searchInput = ref<HTMLInputElement | null>(null)
const fileInput = ref<HTMLInputElement | null>(null)

async function handleImport(e: Event) {
  const input = e.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) return
  const text = await file.text()
  const format = file.name.endsWith('.csv') ? 'csv' : 'json'
  try {
    await api.importTodos(format, text, store.currentListId || undefined)
    store.fetchTodos()
    useToast().show(`导入成功`, 'success')
  } catch {
    useToast().show('导入失败', 'error')
  }
  input.value = ''
  showMore.value = false
}

function triggerImport() {
  fileInput.value?.click()
}

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
    case 'default':
      // Keep backend sort_order
      break
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

interface GroupedItem {
  type: 'group'
  label: string
  count: number
}
type ListItem = Todo | GroupedItem
function isGroup(item: ListItem): item is GroupedItem {
  return (item as GroupedItem).type === 'group'
}
function asTodo(item: ListItem): Todo {
  return item as Todo
}

const groupedTodos = computed(() => {
  const list = filteredTodos()
  if (list.length === 0) return [] as ListItem[]

  const today = new Date().toISOString().slice(0, 10)
  const tomorrow = new Date(Date.now() + 86400000).toISOString().slice(0, 10)
  const afterTomorrow = new Date(Date.now() + 172800000).toISOString().slice(0, 10)

  const groups: { label: string; todos: typeof list }[] = [
    { label: '已过期', todos: [] },
    { label: '今天', todos: [] },
    { label: '明天', todos: [] },
    { label: '后天', todos: [] },
    { label: '更晚', todos: [] },
  ]

  for (const t of list) {
    if (!t.due_date || t.due_date === today) { groups[1].todos.push(t); continue }
    if (t.due_date < today) { groups[0].todos.push(t); continue }
    if (t.due_date === tomorrow) { groups[2].todos.push(t); continue }
    if (t.due_date === afterTomorrow) { groups[3].todos.push(t); continue }
    groups[4].todos.push(t)
  }

  const result: ListItem[] = []
  for (const g of groups) {
    if (g.todos.length === 0) continue
    result.push({ type: 'group', label: g.label, count: g.todos.length })
    result.push(...g.todos)
  }
  return result
})

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

const dragId = ref<number | null>(null)
const dragOverId = ref<number | null>(null)

function onDragStart(todo: { id: number }) {
  dragId.value = todo.id
}

function onDragOver(e: DragEvent, todo: { id: number }) {
  e.preventDefault()
  if (dragId.value && dragId.value !== todo.id) {
    dragOverId.value = todo.id
  }
}

function onDragLeave() {
  dragOverId.value = null
}

async function onDrop(todo: { id: number }) {
  dragOverId.value = null
  if (!dragId.value || dragId.value === todo.id) {
    dragId.value = null
    return
  }
  const srcId = dragId.value
  dragId.value = null

  const arr = store.todos
  const srcIdx = arr.findIndex(t => t.id === srcId)
  const dstIdx = arr.findIndex(t => t.id === todo.id)
  if (srcIdx === -1 || dstIdx === -1) return

  const [item] = arr.splice(srcIdx, 1)
  arr.splice(dstIdx, 0, item)

  const ids = arr.map(t => t.id)
  await store.reorder(ids)
}

function onDragEnd() {
  dragId.value = null
  dragOverId.value = null
}

function onKeydown(e: KeyboardEvent) {
  if ((e.ctrlKey || e.metaKey) && e.key === 'k') {
    e.preventDefault()
    searchInput.value?.focus()
  }
  if (e.altKey && e.key === 'n') {
    e.preventDefault()
    document.getElementById('new-todo-title')?.focus()
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
    <TodoForm v-if="!isReadonly" @created="() => {}" />

    <div class="flex gap-2 mb-3 flex-wrap">
      <input
        ref="searchInput"
        v-model="search"
        type="text"
        placeholder="搜索… (Ctrl+K)"
        class="flex-1 min-w-[120px] px-3 py-2 border border-gray-200 dark:border-gray-800 rounded-lg text-sm outline-none focus:border-indigo-400 bg-white dark:bg-gray-900 dark:text-gray-100 transition-colors"
      />
      <Select v-model="priorityFilter" :options="[
        { label: '全部优先级', value: '' },
        { label: '高', value: 'high' },
        { label: '中', value: 'medium' },
        { label: '低', value: 'low' },
      ]" />
      <Select v-model="tagFilter" :options="[
        { label: '全部标签', value: '' },
        ...allTags.map(t => ({ label: t, value: t })),
      ]" />
      <Select v-model="sortBy" :options="[
        { label: '默认', value: 'default' },
        { label: '最新', value: 'created' },
        { label: '优先级', value: 'priority' },
        { label: '截止日期', value: 'due' },
      ]" />

      <!-- More menu -->
      <div class="relative">
        <button @click="showMore = !showMore" class="px-3 py-2 border border-gray-200 dark:border-gray-800 rounded-lg text-sm bg-white dark:bg-gray-900 text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300 transition-colors">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="5" r="1"/><circle cx="12" cy="12" r="1"/><circle cx="12" cy="19" r="1"/></svg>
        </button>
        <div v-if="showMore" class="absolute right-0 top-full mt-1 bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700 shadow-lg py-1 z-10 min-w-[130px]">
          <button @click="store.downloadExport('json'); showMore = false" class="w-full text-left px-4 py-2 text-sm text-gray-600 dark:text-gray-400 hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors">导出 JSON</button>
          <button @click="store.downloadExport('csv'); showMore = false" class="w-full text-left px-4 py-2 text-sm text-gray-600 dark:text-gray-400 hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors">导出 CSV</button>
          <div class="border-t border-gray-100 dark:border-gray-700 my-1" />
          <button @click="triggerImport" class="w-full text-left px-4 py-2 text-sm text-gray-600 dark:text-gray-400 hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors">导入 (JSON/CSV)</button>
          <div class="border-t border-gray-100 dark:border-gray-700 my-1" />
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
    <div v-else-if="groupedTodos.length === 0" class="text-center py-16 text-gray-300 dark:text-gray-700 text-sm italic">
      暂无待办
    </div>

    <TransitionGroup v-else name="list" tag="div">
      <template v-for="item in groupedTodos" :key="isGroup(item) ? item.label : asTodo(item).id">
        <div v-if="isGroup(item)" class="text-[11px] font-semibold text-gray-400 dark:text-gray-500 tracking-wider px-1 py-2 first:pt-0">
          {{ item.label }}
          <span class="font-normal text-gray-300 dark:text-gray-600 ml-1">{{ item.count }}</span>
        </div>
        <div v-else
          draggable="true"
          @dragstart="onDragStart(asTodo(item))"
          @dragover="onDragOver($event, asTodo(item))"
          @dragleave="onDragLeave"
          @drop="onDrop(asTodo(item))"
          @dragend="onDragEnd"
          :class="{ 'opacity-40': dragId === asTodo(item).id, 'border-t-2 border-indigo-400': dragOverId === asTodo(item).id }"
        >
          <TodoItem :todo="asTodo(item)" :readonly="isReadonly" />
        </div>
      </template>
    </TransitionGroup>

    <div class="mt-8 text-center text-[11px] text-gray-300 dark:text-gray-700 hidden md:block">
      <kbd class="px-1 py-0.5 rounded bg-gray-100 dark:bg-gray-800 text-gray-400">Alt+N</kbd> 新建
      <kbd class="px-1 py-0.5 rounded bg-gray-100 dark:bg-gray-800 text-gray-400 ml-2">Ctrl+K</kbd> 搜索
      <kbd class="px-1 py-0.5 rounded bg-gray-100 dark:bg-gray-800 text-gray-400 ml-2">?</kbd> 快捷键
    </div>
  </div>

  <input ref="fileInput" type="file" accept=".json,.csv" class="hidden" @change="handleImport" />

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
