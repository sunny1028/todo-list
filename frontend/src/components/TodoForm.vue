<script setup lang="ts">
import { computed, ref } from 'vue'
import { useTodos } from '../stores/todo'
import { parseTodoInput } from '../utils/parseTodoInput'
import DatePicker from './ui/DatePicker.vue'
import Select from './ui/Select.vue'

const store = useTodos()
const title = ref('')
const description = ref('')
const priority = ref('medium')
const effort = ref('')
const recurrence = ref('')
const dueDate = ref('')
const tags = ref('')
const expanded = ref(false)

const parsed = computed(() => {
  const t = title.value.trim()
  if (!t) return null
  return parseTodoInput(t)
})

const previewChips = computed(() => {
  const p = parsed.value
  if (!p) return []
  const chips: { type: string; label: string }[] = []

  if (p.priority) {
    const labels: Record<string, string> = { low: '低', medium: '中', high: '高' }
    chips.push({ type: 'priority', label: labels[p.priority] })
  }
  if (p.tags) {
    p.tags.split(',').forEach(tag => chips.push({ type: 'tag', label: `#${tag}` }))
  }
  if (p.effort) {
    const labels: Record<string, string> = { easy: '简单', medium: '中等', hard: '困难' }
    chips.push({ type: 'effort', label: labels[p.effort] })
  }
  if (p.recurrence) {
    const labels: Record<string, string> = { daily: '每天', weekly: '每周', monthly: '每月' }
    chips.push({ type: 'recurrence', label: labels[p.recurrence] })
  }
  if (p.dueDate) {
    const today = new Date().toISOString().slice(0, 10)
    const tomorrow = new Date(Date.now() + 86400000).toISOString().slice(0, 10)
    if (p.dueDate === today) chips.push({ type: 'date', label: '今天' })
    else if (p.dueDate === tomorrow) chips.push({ type: 'date', label: '明天' })
    else chips.push({ type: 'date', label: p.dueDate })
  }

  return chips
})

const emit = defineEmits<{ created: [] }>()

async function handleSubmit() {
  const p = parsed.value
  const finalTitle = p?.title || title.value.trim()
  if (!finalTitle) return
  await store.addTodo({
    title: finalTitle,
    description: description.value.trim() || undefined,
    priority: priority.value !== 'medium' ? priority.value : (p?.priority || 'medium'),
    effort: effort.value || p?.effort || undefined,
    recurrence: recurrence.value || p?.recurrence || undefined,
    tags: tags.value.trim() || p?.tags || undefined,
    due_date: dueDate.value || p?.dueDate || null,
  })
  title.value = ''
  description.value = ''
  priority.value = 'medium'
  effort.value = ''
  recurrence.value = ''
  dueDate.value = ''
  tags.value = ''
  expanded.value = false
  emit('created')
}

function handleKeydown(e: KeyboardEvent) {
  if (e.key === 'Enter' && !e.shiftKey) {
    e.preventDefault()
    handleSubmit()
  }
  if (e.key === 'Escape') {
    expanded.value = false
  }
}
</script>

<template>
  <div class="bg-white dark:bg-gray-900 rounded-xl border border-gray-200 dark:border-gray-800 shadow-sm mb-6 transition-colors">
    <div class="flex items-center gap-2 p-3">
      <!-- Checkbox placeholder -->
      <span class="w-[18px] h-[18px] rounded-full border-2 border-gray-300 dark:border-gray-600 shrink-0" />

      <!-- Title input -->
      <input
        id="new-todo-title"
        v-model="title"
        type="text"
        placeholder="添加待办事项…"
        @focus="expanded = true"
        @keydown="handleKeydown"
        class="flex-1 min-w-0 text-sm bg-transparent border-none outline-none text-gray-800 dark:text-gray-100 placeholder-gray-400 dark:placeholder-gray-500"
      />

      <!-- Preview chips from natural language parsing -->
      <div v-if="previewChips.length" class="flex flex-wrap gap-1 shrink-0">
        <span
          v-for="chip in previewChips"
          :key="chip.type + chip.label"
          class="inline-flex items-center px-1.5 py-0.5 rounded text-[10px] font-medium"
          :class="{
            'bg-orange-50 text-orange-600 dark:bg-orange-900/30 dark:text-orange-400': chip.type === 'priority',
            'bg-blue-50 text-blue-600 dark:bg-blue-900/30 dark:text-blue-400': chip.type === 'date',
            'bg-purple-50 text-purple-600 dark:bg-purple-900/30 dark:text-purple-400': chip.type === 'tag',
            'bg-green-50 text-green-600 dark:bg-green-900/30 dark:text-green-400': chip.type === 'effort',
            'bg-cyan-50 text-cyan-600 dark:bg-cyan-900/30 dark:text-cyan-400': chip.type === 'recurrence',
          }">
          {{ chip.label }}
        </span>
      </div>

      <!-- Action icons -->
      <div class="flex items-center gap-0.5 shrink-0">
        <button type="button" @click="expanded = !expanded"
          class="w-7 h-7 flex items-center justify-center rounded-lg text-gray-400 hover:text-indigo-500 hover:bg-indigo-50 dark:hover:bg-indigo-900/30 transition-colors"
          :class="{ 'text-indigo-500 bg-indigo-50 dark:bg-indigo-900/30': expanded }"
          title="更多选项">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="5" r="1"/><circle cx="12" cy="12" r="1"/><circle cx="12" cy="19" r="1"/></svg>
        </button>
        <button type="submit" @click="handleSubmit" :disabled="!title.trim()"
          class="px-3 py-1.5 bg-indigo-500 text-white rounded-lg text-xs font-semibold hover:bg-indigo-600 disabled:opacity-40 disabled:cursor-default transition-colors">
          添加
        </button>
      </div>
    </div>

    <!-- Expanded row -->
    <div v-if="expanded" class="px-3 pb-3 pt-0 border-t border-gray-100 dark:border-gray-800 flex flex-wrap items-center gap-2">
      <!-- Description -->
      <input
        v-model="description"
        type="text"
        placeholder="备注…"
        class="flex-1 min-w-[120px] px-2.5 py-1.5 border border-gray-150 dark:border-gray-800 rounded-lg text-xs outline-none focus:border-indigo-400 bg-gray-50 dark:bg-gray-800 dark:text-gray-100 transition-colors"
      />
      <!-- Date -->
      <DatePicker v-model="dueDate" />
      <!-- Recurrence -->
      <Select v-model="recurrence" :options="[
        { label: '不重复', value: '' },
        { label: '每天', value: 'daily' },
        { label: '每周', value: 'weekly' },
        { label: '每月', value: 'monthly' },
      ]" />
      <!-- Priority -->
      <Select v-model="priority" :options="[
        { label: '低优先级', value: 'low' },
        { label: '中优先级', value: 'medium' },
        { label: '高优先级', value: 'high' },
      ]" />
      <!-- Effort -->
      <Select v-model="effort" :options="[
        { label: '工作量', value: '' },
        { label: '简单', value: 'easy' },
        { label: '中等', value: 'medium' },
        { label: '困难', value: 'hard' },
      ]" />
      <!-- Tags -->
      <input
        v-model="tags"
        type="text"
        placeholder="标签…"
        class="w-[110px] px-2.5 py-1.5 border border-gray-150 dark:border-gray-800 rounded-lg text-xs outline-none focus:border-indigo-400 bg-gray-50 dark:bg-gray-800 dark:text-gray-100 transition-colors"
      />
    </div>
  </div>
</template>
