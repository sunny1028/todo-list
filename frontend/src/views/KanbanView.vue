<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useTodos } from '../stores/todo'
import { useLists } from '../stores/lists'
import type { Todo } from '../types/todo'
import * as api from '../api/todos'

const store = useTodos()
const listStore = useLists()

const columns = ref<{ id: number; name: string; color: string; todos: Todo[] }[]>([])
const dragTodo = ref<Todo | null>(null)

async function load() {
  await listStore.fetchLists()
  // Build columns: uncategorized + each list
  const raw = await api.listTodos().then(r => r.data)
  columns.value = [
    { id: 0, name: '未分类', color: '#9ca3af', todos: raw.filter(t => t.list_id === 0 && !t.archived) },
    ...listStore.lists.map(l => ({
      id: l.id,
      name: l.name,
      color: l.color,
      todos: raw.filter(t => t.list_id === l.id && !t.archived),
    })),
  ]
}

function onDragStart(todo: Todo) { dragTodo.value = todo }
function onDragOver(e: DragEvent) { e.preventDefault() }

async function onDrop(col: typeof columns.value[0]) {
  if (!dragTodo.value || dragTodo.value.list_id === col.id) return
  await store.editTodo(dragTodo.value.id, { list_id: col.id })
  dragTodo.value.list_id = col.id
  dragTodo.value = null
  load()
}

onMounted(load)
</script>

<template>
  <router-link to="/" class="inline-flex items-center gap-1 text-xs text-gray-400 hover:text-indigo-500 transition-colors mb-3">
    <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="15 18 9 12 15 6"/></svg>
    返回列表
  </router-link>
  <div class="flex gap-4 overflow-x-auto pb-4 -mx-4 px-4" style="min-height: calc(100vh - 200px)">
    <div v-for="col in columns" :key="col.id"
      class="flex-shrink-0 w-64 bg-gray-50 dark:bg-gray-900/50 rounded-xl border border-gray-200 dark:border-gray-800 p-3"
      @dragover="onDragOver" @drop="onDrop(col)">
      <div class="flex items-center gap-2 mb-3 px-1">
        <span class="w-3 h-3 rounded-full shrink-0" :style="{ background: col.color }" />
        <span class="text-sm font-semibold text-gray-700 dark:text-gray-300">{{ col.name }}</span>
        <span class="text-xs text-gray-400 ml-auto">{{ col.todos.length }}</span>
      </div>
      <div class="space-y-2">
        <div v-for="todo in col.todos" :key="todo.id"
          draggable="true" @dragstart="onDragStart(todo)"
          class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-2.5 cursor-grab active:cursor-grabbing shadow-sm"
          :class="{ 'opacity-50': todo.completed }">
          <div class="text-[13px] font-medium truncate text-gray-800 dark:text-gray-200"
            :class="{ 'line-through': todo.completed }">{{ todo.title }}</div>
          <div class="flex items-center gap-2 mt-1.5">
            <span class="text-[10px] px-1.5 py-0.5 rounded-full"
              :class="{
                'bg-blue-100 text-blue-700': todo.priority === 'low',
                'bg-yellow-100 text-yellow-700': todo.priority === 'medium',
                'bg-red-100 text-red-700': todo.priority === 'high',
              }">{{ { low: '低', medium: '中', high: '高' }[todo.priority] }}</span>
            <span v-if="todo.due_date" class="text-[10px] text-gray-400">{{ new Date(todo.due_date).toLocaleDateString('zh-CN') }}</span>
          </div>
        </div>
        <div v-if="col.todos.length === 0" class="text-center py-8 text-xs text-gray-300 dark:text-gray-700">
          拖拽待办到这里
        </div>
      </div>
    </div>
  </div>
</template>
