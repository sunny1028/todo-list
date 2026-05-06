<script setup lang="ts">
import { ref } from 'vue'
import { useTodos } from '../stores/todo'
import DatePicker from './ui/DatePicker.vue'
import Select from './ui/Select.vue'

const store = useTodos()
const title = ref('')
const description = ref('')
const priority = ref('medium')
const dueDate = ref('')
const tags = ref('')
const expanded = ref(false)

const emit = defineEmits<{ created: [] }>()

async function handleSubmit() {
  if (!title.value.trim()) return
  await store.addTodo({
    title: title.value.trim(),
    description: description.value.trim() || undefined,
    priority: priority.value,
    tags: tags.value.trim() || undefined,
    due_date: dueDate.value || null,
  })
  title.value = ''
  description.value = ''
  priority.value = 'medium'
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
      <!-- Priority -->
      <Select v-model="priority" :options="[
        { label: '低优先级', value: 'low' },
        { label: '中优先级', value: 'medium' },
        { label: '高优先级', value: 'high' },
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
