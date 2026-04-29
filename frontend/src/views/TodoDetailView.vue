<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useTodos } from '../stores/todo'
import type { Todo, Attachment } from '../types/todo'
import * as api from '../api/todos'
import ConfirmDialog from '../components/ui/ConfirmDialog.vue'
import { useToast } from '../stores/toast'

const route = useRoute()
const router = useRouter()
const store = useTodos()
const toast = useToast()

const todo = ref<Todo | null>(null)
const editing = ref(false)
const editTitle = ref('')
const editDescription = ref('')
const editPriority = ref<'low' | 'medium' | 'high'>('medium')
const editDueDate = ref('')
const editTags = ref('')
const showConfirm = ref(false)

const attachments = ref<Attachment[]>([])
const uploading = ref(false)

onMounted(async () => {
  const id = Number(route.params.id)
  await store.fetchTodos()
  todo.value = store.todos.find((t) => t.id === id) ?? null
  if (todo.value) {
    editTitle.value = todo.value.title
    editDescription.value = todo.value.description
    editPriority.value = todo.value.priority
    editDueDate.value = todo.value.due_date?.slice(0, 10) || ''
    editTags.value = todo.value.tags
    loadAttachments()
  }
})

async function loadAttachments() {
  if (!todo.value) return
  try {
    const res = await api.listAttachments(todo.value.id)
    attachments.value = res.data
  } catch { /* silent */ }
}

async function uploadFile(e: Event) {
  const input = e.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file || !todo.value) return
  uploading.value = true
  try {
    await api.uploadAttachment(todo.value.id, file)
    toast.show('上传成功', 'success')
    loadAttachments()
  } catch {
    toast.show('上传失败', 'error')
  } finally {
    uploading.value = false
    input.value = ''
  }
}

async function deleteAttachment(id: number) {
  try {
    await api.deleteAttachment(id)
    attachments.value = attachments.value.filter((a) => a.id !== id)
    toast.show('附件已删除', 'success')
  } catch {
    toast.show('删除失败', 'error')
  }
}

function isImage(mime: string) {
  return mime.startsWith('image/')
}

function formatSize(bytes: number) {
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
  return (bytes / (1024 * 1024)).toFixed(1) + ' MB'
}

async function save() {
  if (!todo.value || !editTitle.value.trim()) return
  await store.editTodo(todo.value.id, {
    title: editTitle.value.trim(),
    description: editDescription.value.trim() || undefined,
    priority: editPriority.value,
    due_date: editDueDate.value || null,
    tags: editTags.value.trim() || undefined,
  })
  todo.value = store.todos.find((t) => t.id === todo.value!.id) ?? null
  editing.value = false
}

async function remove() {
  if (!todo.value) return
  await store.removeTodo(todo.value.id)
  router.replace('/')
}

function formatDate(d: string) {
  return new Date(d).toLocaleDateString('zh-CN')
}

function priorityLabel(p: string) {
  return { low: '低', medium: '中', high: '高' }[p] || p
}

function parseTags(s: string) {
  return s ? s.split(',').filter(Boolean).map((t) => t.trim()) : []
}
</script>

<template>
  <div v-if="todo" class="max-w-lg mx-auto">
    <button
      @click="router.back()"
      class="text-indigo-600 dark:text-indigo-400 text-sm mb-4 flex items-center gap-1"
    >
      &larr; 返回
    </button>

    <template v-if="!editing">
      <div class="bg-white dark:bg-gray-900 rounded-xl border border-gray-200 dark:border-gray-800 p-5 shadow-sm">
        <div class="flex items-center justify-between mb-3">
          <h2 class="text-xl font-bold text-gray-800 dark:text-gray-100" :class="{ 'line-through opacity-50': todo.completed }">
            {{ todo.title }}
          </h2>
          <span
            class="text-xs font-semibold px-2 py-0.5 rounded-full"
            :class="{
              'bg-blue-100 dark:bg-blue-900/40 text-blue-700 dark:text-blue-400': todo.priority === 'low',
              'bg-yellow-100 dark:bg-yellow-900/40 text-yellow-700 dark:text-yellow-400': todo.priority === 'medium',
              'bg-red-100 dark:bg-red-900/40 text-red-700 dark:text-red-400': todo.priority === 'high',
            }"
          >
            {{ priorityLabel(todo.priority) }}
          </span>
        </div>
        <p v-if="todo.description" class="text-gray-500 dark:text-gray-400 text-sm mb-3">{{ todo.description }}</p>
        <div v-if="parseTags(todo.tags).length > 0" class="flex gap-1.5 mb-3 flex-wrap">
          <span v-for="tag in parseTags(todo.tags)" :key="tag"
            class="text-[11px] px-2 py-0.5 rounded-md bg-indigo-50 dark:bg-indigo-900/30 text-indigo-600 dark:text-indigo-400"
          >{{ tag }}</span>
        </div>
        <p class="text-xs text-gray-400 dark:text-gray-500">创建于 {{ formatDate(todo.created_at) }}</p>
        <p v-if="todo.due_date" class="text-xs text-gray-400 dark:text-gray-500">截止 {{ formatDate(todo.due_date) }}</p>

        <!-- Attachments -->
        <div class="mt-4 pt-4 border-t border-gray-100 dark:border-gray-800">
          <div class="flex items-center justify-between mb-2">
            <span class="text-xs font-semibold text-gray-400">附件</span>
            <label class="cursor-pointer text-xs text-indigo-500 hover:text-indigo-600 transition-colors">
              {{ uploading ? '上传中…' : '+ 添加' }}
              <input type="file" class="hidden" @change="uploadFile" :disabled="uploading" />
            </label>
          </div>
          <div v-if="attachments.length === 0" class="text-xs text-gray-300 dark:text-gray-700">暂无附件</div>
          <div v-for="att in attachments" :key="att.id" class="flex items-center justify-between py-1.5">
            <a :href="api.attachmentUrl(att.id)" target="_blank" class="flex items-center gap-2 text-sm text-gray-600 dark:text-gray-400 hover:text-indigo-500 truncate">
              <span v-if="isImage(att.mime_type)" class="w-8 h-8 rounded-lg bg-gray-100 dark:bg-gray-800 overflow-hidden shrink-0">
                <img :src="api.attachmentUrl(att.id)" class="w-full h-full object-cover" />
              </span>
              <span v-else class="w-8 h-8 rounded-lg bg-gray-100 dark:bg-gray-800 flex items-center justify-center text-xs shrink-0">📄</span>
              <span class="truncate">{{ att.filename }}</span>
              <span class="text-xs text-gray-400 shrink-0">{{ formatSize(att.size) }}</span>
            </a>
            <button @click="deleteAttachment(att.id)" class="text-gray-400 hover:text-red-500 ml-2 shrink-0">
              <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
            </button>
          </div>
        </div>

        <div class="flex gap-2 mt-4">
          <button @click="editing = true" class="px-4 py-2 text-sm bg-indigo-500 text-white rounded-lg hover:bg-indigo-600 transition-colors">编辑</button>
          <button @click="showConfirm = true" class="px-4 py-2 text-sm bg-red-500 text-white rounded-lg hover:bg-red-600 transition-colors">删除</button>
        </div>
      </div>
    </template>

    <template v-else>
      <div class="bg-white dark:bg-gray-900 rounded-xl border border-gray-200 dark:border-gray-800 p-5 shadow-sm flex flex-col gap-3">
        <input v-model="editTitle" class="border border-gray-200 dark:border-gray-700 rounded-lg px-3 py-2 text-sm bg-white dark:bg-gray-800 dark:text-gray-100 outline-none" placeholder="标题" />
        <input v-model="editDescription" class="border border-gray-200 dark:border-gray-700 rounded-lg px-3 py-2 text-sm bg-white dark:bg-gray-800 dark:text-gray-100 outline-none" placeholder="备注" />
        <input v-model="editDueDate" type="date" class="border border-gray-200 dark:border-gray-700 rounded-lg px-3 py-2 text-sm bg-white dark:bg-gray-800 dark:text-gray-100 outline-none" />
        <select v-model="editPriority" class="border border-gray-200 dark:border-gray-700 rounded-lg px-3 py-2 text-sm bg-white dark:bg-gray-800 dark:text-gray-100 outline-none">
          <option value="low">低</option>
          <option value="medium">中</option>
          <option value="high">高</option>
        </select>
        <input v-model="editTags" class="border border-gray-200 dark:border-gray-700 rounded-lg px-3 py-2 text-sm bg-white dark:bg-gray-800 dark:text-gray-100 outline-none" placeholder="标签（逗号分隔）" />
        <div class="flex gap-2">
          <button @click="save" class="px-4 py-2 text-sm bg-indigo-500 text-white rounded-lg">保存</button>
          <button @click="editing = false" class="px-4 py-2 text-sm bg-gray-200 dark:bg-gray-700 text-gray-600 dark:text-gray-300 rounded-lg">取消</button>
        </div>
      </div>
    </template>
  </div>

  <div v-else class="text-center py-20 text-gray-400 dark:text-gray-600">未找到该待办。</div>

  <ConfirmDialog
    :open="showConfirm"
    title="删除待办"
    :message="`确定要删除「${todo?.title}」吗？此操作无法撤销。`"
    @confirm="remove(); showConfirm = false"
    @cancel="showConfirm = false"
  />
</template>
