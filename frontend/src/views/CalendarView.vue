<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import type { Todo } from '../types/todo'
import * as api from '../api/todos'

const today = new Date()
const currentYear = ref(today.getFullYear())
const currentMonth = ref(today.getMonth())
const todos = ref<Todo[]>([])
const selectedDate = ref<string | null>(null)

const monthLabel = computed(() => `${currentYear.value}年${currentMonth.value + 1}月`)

const days = computed(() => {
  const first = new Date(currentYear.value, currentMonth.value, 1)
  const last = new Date(currentYear.value, currentMonth.value + 1, 0)
  const startDay = first.getDay()
  const days: (number | null)[] = []

  for (let i = 0; i < startDay; i++) days.push(null)
  for (let d = 1; d <= last.getDate(); d++) days.push(d)

  return days
})

async function load() {
  const res = await api.listTodos()
  todos.value = res.data.filter(t => !t.archived && t.due_date)
}

function dateStr(d: number) {
  const m = String(currentMonth.value + 1).padStart(2, '0')
  const day = String(d).padStart(2, '0')
  return `${currentYear.value}-${m}-${day}`
}

function todosForDay(d: number) {
  const ds = dateStr(d)
  return todos.value.filter(t => t.due_date === ds)
}

function prevMonth() {
  if (currentMonth.value === 0) { currentMonth.value = 11; currentYear.value-- }
  else currentMonth.value--
}
function nextMonth() {
  if (currentMonth.value === 11) { currentMonth.value = 0; currentYear.value++ }
  else currentMonth.value++
}

function todayDateStr() { return today.toISOString().slice(0, 10) }
function isToday(d: number) { return dateStr(d) === todayDateStr() }

onMounted(load)
</script>

<template>
  <div class="max-w-lg mx-auto">
    <div class="flex items-center justify-between mb-4">
      <button @click="prevMonth" class="text-gray-400 hover:text-gray-600 dark:hover:text-gray-300 px-2">&larr;</button>
      <h2 class="text-lg font-bold text-gray-800 dark:text-gray-100">{{ monthLabel }}</h2>
      <button @click="nextMonth" class="text-gray-400 hover:text-gray-600 dark:hover:text-gray-300 px-2">&rarr;</button>
    </div>

    <!-- Day headers -->
    <div class="grid grid-cols-7 mb-1">
      <div v-for="d in ['日','一','二','三','四','五','六']" :key="d" class="text-center text-xs text-gray-400 dark:text-gray-600 py-1">{{ d }}</div>
    </div>

    <!-- Calendar grid -->
    <div class="grid grid-cols-7 border-t border-l border-gray-200 dark:border-gray-800 rounded-xl overflow-hidden">
      <div v-for="(d, i) in days" :key="i"
        @click="d && (selectedDate = dateStr(d))"
        class="border-r border-b border-gray-200 dark:border-gray-800 p-1.5 min-h-[60px] cursor-pointer hover:bg-gray-50 dark:hover:bg-gray-800/50 transition-colors"
        :class="{ 'bg-white dark:bg-gray-900': !isToday(d || 0), 'bg-indigo-50 dark:bg-indigo-900/20': isToday(d || 0) }">
        <template v-if="d">
          <div class="text-xs font-semibold mb-0.5" :class="{ 'text-indigo-600 dark:text-indigo-400': isToday(d), 'text-gray-600 dark:text-gray-400': !isToday(d) }">{{ d }}</div>
          <div v-for="todo in todosForDay(d).slice(0, 2)" :key="todo.id" class="text-[9px] truncate px-1 py-0.5 rounded mb-0.5"
            :class="{
              'bg-red-100 dark:bg-red-900/30 text-red-700 dark:text-red-400': todo.priority === 'high',
              'bg-yellow-100 dark:bg-yellow-900/30 text-yellow-700 dark:text-yellow-400': todo.priority === 'medium',
              'bg-blue-100 dark:bg-blue-900/30 text-blue-700 dark:text-blue-400': todo.priority === 'low',
            }">{{ todo.title }}</div>
          <div v-if="todosForDay(d).length > 2" class="text-[9px] text-gray-400 px-1">+{{ todosForDay(d).length - 2 }}</div>
        </template>
      </div>
    </div>

    <!-- Selected day todos -->
    <div v-if="selectedDate" class="mt-4 bg-white dark:bg-gray-900 rounded-xl border border-gray-200 dark:border-gray-800 p-4">
      <div class="flex items-center justify-between mb-3">
        <h3 class="text-sm font-semibold text-gray-700 dark:text-gray-300">{{ selectedDate }}</h3>
        <button @click="selectedDate = null" class="text-xs text-gray-400 hover:text-gray-600">&times;</button>
      </div>
      <div v-if="todos.filter(t => t.due_date === selectedDate).length === 0" class="text-xs text-gray-300 dark:text-gray-700 py-4 text-center">暂无待办</div>
      <div v-for="todo in todos.filter(t => t.due_date === selectedDate)" :key="todo.id"
        class="flex items-center gap-2 py-2 border-b border-gray-50 dark:border-gray-800 last:border-0">
        <span class="text-[13px] text-gray-800 dark:text-gray-200 flex-1">{{ todo.title }}</span>
        <span class="text-[10px] px-1.5 py-0.5 rounded-full"
          :class="{
            'bg-blue-100 text-blue-700': todo.priority === 'low',
            'bg-yellow-100 text-yellow-700': todo.priority === 'medium',
            'bg-red-100 text-red-700': todo.priority === 'high',
          }">{{ { low: '低', medium: '中', high: '高' }[todo.priority] }}</span>
      </div>
    </div>
  </div>
</template>
