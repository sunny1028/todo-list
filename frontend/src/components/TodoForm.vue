<script setup lang="ts">
import { ref } from 'vue'
import { useTodos } from '../stores/todo'
import DatePicker from './ui/DatePicker.vue'

const store = useTodos()
const title = ref('')
const description = ref('')
const priority = ref('medium')
const dueDate = ref('')
const tags = ref('')

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
  emit('created')
}
</script>

<template>
  <form @submit.prevent="handleSubmit" class="flex flex-wrap gap-2 mb-6">
    <input
      id="new-todo-title"
      v-model="title"
      type="text"
      placeholder="添加待办事项… (Alt+N)"
      class="flex-1 min-w-[180px] px-3 py-2.5 border border-gray-200 dark:border-gray-800 rounded-lg text-sm outline-none focus:border-indigo-400 bg-white dark:bg-gray-900 dark:text-gray-100 transition-colors"
    />
    <input
      v-model="description"
      type="text"
      placeholder="备注（可选）"
      class="flex-1 min-w-[150px] px-3 py-2.5 border border-gray-200 dark:border-gray-800 rounded-lg text-sm outline-none focus:border-indigo-400 bg-white dark:bg-gray-900 dark:text-gray-100 transition-colors"
    />
    <DatePicker v-model="dueDate" />
    <select
      v-model="priority"
      class="px-3 py-2.5 border border-gray-200 dark:border-gray-800 rounded-lg text-sm bg-white dark:bg-gray-900 dark:text-gray-100 outline-none"
    >
      <option value="low">低优先级</option>
      <option value="medium">中优先级</option>
      <option value="high">高优先级</option>
    </select>
    <input
      v-model="tags"
      type="text"
      placeholder="标签（逗号分隔）"
      class="w-[160px] px-3 py-2.5 border border-gray-200 dark:border-gray-800 rounded-lg text-sm outline-none focus:border-indigo-400 bg-white dark:bg-gray-900 dark:text-gray-100 transition-colors"
    />
    <button
      type="submit"
      :disabled="!title.trim()"
      class="px-5 py-2.5 bg-indigo-500 text-white rounded-lg text-sm font-semibold hover:bg-indigo-600 disabled:opacity-40 disabled:cursor-default transition-colors"
    >
      添加
    </button>
  </form>
</template>
