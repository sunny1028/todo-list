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
const selectedAll = ref(false)
const searchInput = ref<HTMLInputElement | null>(null)

const filterStatus = computed(() => (route.query.filter as string) || '')

function filteredTodos() {
  let list = [...store.todos]

  if (filterStatus.value === 'active') list = list.filter((t) => !t.completed)
  else if (filterStatus.value === 'completed') list = list.filter((t) => t.completed)

  if (priorityFilter.value) list = list.filter((t) => t.priority === priorityFilter.value)
  if (tagFilter.value) list = list.filter((t) => t.tags.includes(tagFilter.value))
  if (search.value.trim()) {
    const q = search.value.trim().toLowerCase()
    list = list.filter((t) => t.title.toLowerCase().includes(q))
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
}

function toggleSelectAll() {
  selectedAll.value = !selectedAll.value
  store.todos.forEach((t) => {
    if (t.completed !== selectedAll.value) {
      store.toggle(t.id)
    }
  })
}

// Drag & drop
const dragIndex = ref<number | null>(null)

function onDragStart(idx: number) {
  dragIndex.value = idx
}

function onDragOver(e: DragEvent) {
  e.preventDefault()
}

async function onDrop(idx: number) {
  if (dragIndex.value === null || dragIndex.value === idx) return
  const list = filteredTodos()
  const ids = list.map((t) => t.id)
  const [moved] = ids.splice(dragIndex.value, 1)
  ids.splice(idx, 0, moved)
  // Optimistic UI update
  store.todos = ids.map((id) => store.todos.find((t) => t.id === id)!).filter(Boolean)
  await store.reorder(ids)
  dragIndex.value = null
}

// Keyboard shortcuts
function onKeydown(e: KeyboardEvent) {
  if ((e.ctrlKey || e.metaKey) && e.key === 'k') {
    e.preventDefault()
    searchInput.value?.focus()
  }
  if ((e.ctrlKey || e.metaKey) && e.key === 'n') {
    e.preventDefault()
    // Focus the first input in TodoForm
    const input = document.querySelector<HTMLInputElement>('input[placeholder="添加待办事项…"]')
    input?.focus()
  }
  if (e.key === 'Escape') {
    search.value = ''
    priorityFilter.value = ''
    tagFilter.value = ''
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

    <!-- Filters -->
    <div class="flex gap-2 mb-3 flex-wrap">
      <input
        ref="searchInput"
        v-model="search"
        type="text"
        placeholder="搜索… (Ctrl+K)"
        class="flex-1 min-w-[120px] px-3 py-2 border border-gray-200 dark:border-gray-800 rounded-lg text-sm outline-none focus:border-indigo-400 bg-white dark:bg-gray-900 dark:text-gray-100 transition-colors"
      />
      <select
        v-model="priorityFilter"
        class="px-3 py-2 border border-gray-200 dark:border-gray-800 rounded-lg text-sm bg-white dark:bg-gray-900 dark:text-gray-100 outline-none"
      >
        <option value="">全部优先级</option>
        <option value="high">高</option>
        <option value="medium">中</option>
        <option value="low">低</option>
      </select>
      <select
        v-model="tagFilter"
        class="px-3 py-2 border border-gray-200 dark:border-gray-800 rounded-lg text-sm bg-white dark:bg-gray-900 dark:text-gray-100 outline-none"
      >
        <option value="">全部标签</option>
        <option v-for="tag in allTags" :key="tag" :value="tag">{{ tag }}</option>
      </select>
      <select
        v-model="sortBy"
        class="px-3 py-2 border border-gray-200 dark:border-gray-800 rounded-lg text-sm bg-white dark:bg-gray-900 dark:text-gray-100 outline-none"
      >
        <option value="created">最新</option>
        <option value="priority">优先级</option>
        <option value="due">截止日期</option>
      </select>
    </div>

    <!-- Batch actions -->
    <div v-if="store.todos.length > 0" class="flex items-center justify-between mb-3 text-xs text-gray-400 dark:text-gray-500">
      <span>{{ store.todos.filter(t => !t.completed).length }} 项未完成 / 共 {{ store.todos.length }} 项</span>
      <div class="flex gap-2">
        <button
          @click="toggleSelectAll"
          class="px-3 py-1 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors"
        >
          {{ selectedAll ? '取消全选' : '全选' }}
        </button>
        <button
          @click="store.downloadExport('json')"
          class="px-3 py-1 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors"
        >
          导出 JSON
        </button>
        <button
          @click="store.downloadExport('csv')"
          class="px-3 py-1 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors"
        >
          导出 CSV
        </button>
        <button
          v-if="completedCount > 0"
          @click="showClearConfirm = true"
          class="px-3 py-1 rounded-lg text-red-500 hover:bg-red-50 dark:hover:bg-red-900/30 transition-colors"
        >
          清除已完成 ({{ completedCount }})
        </button>
      </div>
    </div>

    <div v-if="store.loading">
      <Skeleton :count="5" />
    </div>
    <div v-else-if="store.error" class="text-center py-10 text-red-400 text-sm">{{ store.error }}</div>
    <div v-else-if="filteredTodos().length === 0" class="text-center py-16 text-gray-300 dark:text-gray-700 text-sm italic">
      暂无待办，添加第一个吧！
    </div>

    <TransitionGroup v-else name="list" tag="div">
      <div
        v-for="(todo, idx) in filteredTodos()"
        :key="todo.id"
        draggable="true"
        @dragstart="onDragStart(idx)"
        @dragover="onDragOver"
        @drop="onDrop(idx)"
      >
        <TodoItem :todo="todo" />
      </div>
    </TransitionGroup>

    <!-- Keyboard hints -->
    <div class="mt-8 text-center text-[11px] text-gray-300 dark:text-gray-700 hidden md:block">
      <kbd class="px-1 py-0.5 rounded bg-gray-100 dark:bg-gray-800 text-gray-400">Ctrl+N</kbd> 新建
      &nbsp;
      <kbd class="px-1 py-0.5 rounded bg-gray-100 dark:bg-gray-800 text-gray-400">Ctrl+K</kbd> 搜索
      &nbsp;
      <kbd class="px-1 py-0.5 rounded bg-gray-100 dark:bg-gray-800 text-gray-400">Esc</kbd> 清除筛选
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
.list-enter-active,
.list-leave-active {
  transition: all 0.3s ease;
}
.list-enter-from {
  opacity: 0;
  transform: translateY(-12px);
}
.list-leave-to {
  opacity: 0;
  transform: translateX(20px);
}
</style>
