<script setup lang="ts">
import { computed, ref } from 'vue'
import type { Todo, Subtask } from '../types/todo'
import { useTodos } from '../stores/todo'
import { useRouter } from 'vue-router'
import { useResponsive } from '../composables/useResponsive'
import * as api from '../api/todos'
import ConfirmDialog from './ui/ConfirmDialog.vue'

const props = defineProps<{ todo: Todo }>()
const store = useTodos()
const router = useRouter()
const { isMobile } = useResponsive()

const editing = ref(false)
const editTitle = ref('')
const editDescription = ref('')
const editPriority = ref<'low' | 'medium' | 'high'>('medium')
const editDueDate = ref('')
const editTags = ref('')
const showConfirm = ref(false)
const showSubtasks = ref(false)
const subtasks = ref<Subtask[]>([])
const subInput = ref('')

const isOverdue = computed(() => {
  if (!props.todo.due_date || props.todo.completed) return false
  return new Date(props.todo.due_date) < new Date(new Date().toDateString())
})

const isDueToday = computed(() => {
  if (!props.todo.due_date || props.todo.completed) return false
  return props.todo.due_date === new Date().toISOString().slice(0, 10)
})

async function loadSubtasks() {
  if (!showSubtasks.value) return
  subtasks.value = await store.fetchSubtasks(props.todo.id)
}

async function addSubtask() {
  if (!subInput.value.trim()) return
  try {
    await api.createSubtask(props.todo.id, subInput.value.trim())
    subInput.value = ''
    loadSubtasks()
  } catch { /* silent */ }
}

async function toggleSubtask(id: number) {
  await api.toggleSubtask(id)
  loadSubtasks()
}

async function deleteSubtask(id: number) {
  await api.deleteSubtask(id)
  loadSubtasks()
}

const subProgress = computed(() => {
  if (subtasks.value.length === 0) return null
  return subtasks.value.filter(s => s.completed).length
})

function toggleSubtasks() {
  showSubtasks.value = !showSubtasks.value
  if (showSubtasks.value) loadSubtasks()
}

function startEdit() {
  editTitle.value = props.todo.title
  editDescription.value = props.todo.description
  editPriority.value = props.todo.priority
  editDueDate.value = props.todo.due_date?.slice(0, 10) || ''
  editTags.value = props.todo.tags
  editing.value = true
}

function cancelEdit() { editing.value = false }

async function saveEdit() {
  if (!editTitle.value.trim()) return
  await store.editTodo(props.todo.id, {
    title: editTitle.value.trim(),
    description: editDescription.value.trim() || undefined,
    priority: editPriority.value,
    tags: editTags.value.trim() || undefined,
    due_date: editDueDate.value || null,
  })
  editing.value = false
}

async function handleToggle() { await store.toggle(props.todo.id) }
async function handleDelete() { await store.removeTodo(props.todo.id) }

function openDetail() {
  if (isMobile.value) router.push(`/todo/${props.todo.id}`)
}

function formatDate(d: string) { return new Date(d).toLocaleDateString('zh-CN') }
function priorityLabel(p: string) { return { low: '低', medium: '中', high: '高' }[p] || p }
function parseTags(s: string) { return s ? s.split(',').filter(Boolean).map(t => t.trim()) : [] }
</script>

<template>
  <div
    class="bg-white dark:bg-gray-900 rounded-xl border mb-2 shadow-sm transition-colors"
    :class="{
      'border-gray-100 dark:border-gray-800': !isOverdue && !isDueToday,
      'border-red-300 dark:border-red-800 bg-red-50/30 dark:bg-red-950/20': isOverdue,
      'border-amber-300 dark:border-amber-800 bg-amber-50/30 dark:bg-amber-950/20': isDueToday,
      'opacity-60': todo.completed,
    }"
  >
    <!-- Main row -->
    <div class="flex items-center justify-between gap-3 p-3.5" @click="openDetail">
      <template v-if="!editing">
        <div class="flex items-center gap-3 min-w-0 flex-1">
          <input type="checkbox" :checked="todo.completed" @change="handleToggle" @click.stop
            class="w-[18px] h-[18px] accent-indigo-500 cursor-pointer shrink-0" />
          <div class="flex flex-col min-w-0">
            <span class="text-[15px] font-semibold truncate"
              :class="{ 'line-through text-gray-400 dark:text-gray-600': todo.completed, 'text-gray-800 dark:text-gray-100': !todo.completed }">
              {{ todo.title }}
            </span>
            <div class="flex items-center gap-1.5 mt-0.5 flex-wrap">
              <span v-for="tag in parseTags(todo.tags)" :key="tag"
                class="text-[10px] px-1.5 py-0.5 rounded-md bg-indigo-50 dark:bg-indigo-900/30 text-indigo-600 dark:text-indigo-400">
                {{ tag }}
              </span>
              <span v-if="subProgress !== null" class="text-[10px] text-gray-400">
                {{ subProgress }}/{{ subtasks.length }}
              </span>
            </div>
          </div>
        </div>
        <div class="flex items-center gap-1.5 shrink-0">
          <span class="text-[11px] font-semibold px-2 py-0.5 rounded-full"
            :class="{
              'bg-blue-100 dark:bg-blue-900/40 text-blue-700 dark:text-blue-400': todo.priority === 'low',
              'bg-yellow-100 dark:bg-yellow-900/40 text-yellow-700 dark:text-yellow-400': todo.priority === 'medium',
              'bg-red-100 dark:bg-red-900/40 text-red-700 dark:text-red-400': todo.priority === 'high',
            }">{{ priorityLabel(todo.priority) }}</span>
          <span v-if="todo.due_date" class="text-[11px]"
            :class="{
              'text-red-500 font-semibold': isOverdue,
              'text-amber-500 font-semibold': isDueToday,
              'text-gray-400 dark:text-gray-500': !isOverdue && !isDueToday && !todo.completed,
            }">{{ isOverdue ? '已过期' : isDueToday ? '今天' : formatDate(todo.due_date) }}</span>

          <!-- Subtask toggle -->
          <button @click.stop="toggleSubtasks" class="w-7 h-7 flex items-center justify-center rounded-lg text-gray-300 dark:text-gray-600 hover:text-indigo-500 transition-colors"
            :class="{ 'text-indigo-500': showSubtasks }">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="8" y1="6" x2="21" y2="6"/><line x1="8" y1="12" x2="21" y2="12"/><line x1="8" y1="18" x2="21" y2="18"/>
              <line x1="3" y1="6" x2="3.01" y2="6"/><line x1="3" y1="12" x2="3.01" y2="12"/><line x1="3" y1="18" x2="3.01" y2="18"/>
            </svg>
          </button>

          <button @click.stop="startEdit"
            class="w-7 h-7 flex items-center justify-center rounded-lg text-gray-300 dark:text-gray-600 hover:text-gray-500 transition-colors">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/></svg>
          </button>

          <!-- Archive or Delete -->
          <button v-if="todo.completed && !todo.archived" @click.stop="store.archiveTodo(todo.id)"
            class="w-7 h-7 flex items-center justify-center rounded-lg text-gray-300 dark:text-gray-600 hover:text-amber-500 transition-colors">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="21 8 21 21 3 21 3 8"/><rect x="1" y="3" width="22" height="5"/><line x1="10" y1="12" x2="14" y2="12"/></svg>
          </button>
          <button v-else @click.stop="showConfirm = true"
            class="w-7 h-7 flex items-center justify-center rounded-lg text-gray-300 dark:text-gray-600 hover:text-red-400 transition-colors">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="3 6 5 6 21 6"/><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/></svg>
          </button>
        </div>
      </template>
      <template v-else>
        <div @click.stop class="flex flex-wrap gap-2 w-full">
          <input v-model="editTitle" class="flex-1 min-w-[100px] px-2.5 py-2 border border-gray-200 dark:border-gray-700 rounded-lg text-sm outline-none focus:border-indigo-400 bg-white dark:bg-gray-800 dark:text-gray-100" placeholder="标题" />
          <input v-model="editDescription" class="flex-1 min-w-[100px] px-2.5 py-2 border border-gray-200 dark:border-gray-700 rounded-lg text-sm outline-none focus:border-indigo-400 bg-white dark:bg-gray-800 dark:text-gray-100" placeholder="备注" />
          <input v-model="editDueDate" type="date" class="px-2.5 py-2 border border-gray-200 dark:border-gray-700 rounded-lg text-sm bg-white dark:bg-gray-800 dark:text-gray-100 outline-none" />
          <select v-model="editPriority" class="px-2.5 py-2 border border-gray-200 dark:border-gray-700 rounded-lg text-sm bg-white dark:bg-gray-800 dark:text-gray-100 outline-none">
            <option value="low">低</option><option value="medium">中</option><option value="high">高</option>
          </select>
          <input v-model="editTags" class="w-[120px] px-2.5 py-2 border border-gray-200 dark:border-gray-700 rounded-lg text-sm bg-white dark:bg-gray-800 dark:text-gray-100 outline-none" placeholder="标签" />
          <button @click="saveEdit" class="px-3 py-2 bg-indigo-500 text-white rounded-lg text-sm font-semibold">保存</button>
          <button @click="cancelEdit" class="px-3 py-2 bg-gray-200 dark:bg-gray-700 text-gray-600 dark:text-gray-300 rounded-lg text-sm font-semibold">取消</button>
        </div>
      </template>
    </div>

    <!-- Subtasks panel -->
    <div v-if="showSubtasks" class="px-4 pb-3 border-t border-gray-50 dark:border-gray-800 pt-2" @click.stop>
      <div v-for="st in subtasks" :key="st.id" class="flex items-center gap-2 py-1.5">
        <input type="checkbox" :checked="st.completed" @change="toggleSubtask(st.id)"
          class="w-[15px] h-[15px] accent-indigo-500 cursor-pointer shrink-0" />
        <span class="text-[13px] flex-1" :class="{ 'line-through text-gray-300 dark:text-gray-600': st.completed, 'text-gray-600 dark:text-gray-400': !st.completed }">{{ st.title }}</span>
        <button @click="deleteSubtask(st.id)" class="text-gray-300 dark:text-gray-700 hover:text-red-400 shrink-0">
          <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
        </button>
      </div>
      <form @submit.prevent="addSubtask" class="flex gap-1.5 mt-1">
        <input v-model="subInput" class="flex-1 px-2.5 py-1.5 border border-gray-150 dark:border-gray-800 rounded-lg text-xs bg-gray-50 dark:bg-gray-800 dark:text-gray-200 outline-none" placeholder="添加子任务…" />
        <button type="submit" class="px-3 py-1.5 bg-indigo-500 text-white rounded-lg text-xs font-semibold">添加</button>
      </form>
    </div>
  </div>

  <ConfirmDialog
    :open="showConfirm"
    title="删除待办"
    :message="`确定要删除「${todo.title}」吗？（10秒内可撤销）`"
    @confirm="handleDelete(); showConfirm = false"
    @cancel="showConfirm = false"
  />
</template>
