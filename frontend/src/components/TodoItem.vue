<script setup lang="ts">
import { computed, ref } from 'vue'
import type { Todo, Subtask } from '../types/todo'
import { useTodos } from '../stores/todo'
import { useRouter } from 'vue-router'
import { useResponsive } from '../composables/useResponsive'
import * as api from '../api/todos'
import ConfirmDialog from './ui/ConfirmDialog.vue'
import Select from './ui/Select.vue'

const props = defineProps<{ todo: Todo }>()
const store = useTodos()
const router = useRouter()
const { isMobile } = useResponsive()

const editing = ref(false)
const editTitle = ref('')
const editDescription = ref('')
const editPriority = ref<'low' | 'medium' | 'high'>('medium')
const editEffort = ref('')
const editRecurrence = ref('')
const editDueDate = ref('')
const editTags = ref('')
const showConfirm = ref(false)
const showSubtasks = ref(false)
const subtasks = ref<Subtask[]>([])
const subInput = ref('')
const editingSubtask = ref(0)
const editSubtitle = ref('')

const isOverdue = computed(() => {
  if (!props.todo.due_date || props.todo.completed) return false
  return new Date(props.todo.due_date) < new Date(new Date().toDateString())
})

const isDueToday = computed(() => {
  if (!props.todo.due_date || props.todo.completed) return false
  return props.todo.due_date === new Date().toISOString().slice(0, 10)
})

async function loadSubtasks() {
  subtasks.value = await store.fetchSubtasks(props.todo.id)
}

function bumpParentCounts() {
  props.todo.subtask_count = subtasks.value.length
  props.todo.subtask_completed = subtasks.value.filter(s => s.completed).length
}

async function addSubtask() {
  if (!subInput.value.trim()) return
  try {
    const res = await api.createSubtask(props.todo.id, subInput.value.trim())
    subInput.value = ''
    subtasks.value.push(res.data)
    props.todo.subtask_count++
  } catch { /* silent */ }
}

async function toggleSubtask(id: number) {
  const st = subtasks.value.find(s => s.id === id)
  if (!st) return
  st.completed = !st.completed
  props.todo.subtask_completed += st.completed ? 1 : -1
  // Auto-complete parent when all subtasks are done
  if (subtasks.value.every(s => s.completed)) {
    props.todo.completed = true
  } else {
    props.todo.completed = false
  }
  try {
    await api.toggleSubtask(id)
  } catch {
    st.completed = !st.completed
    props.todo.subtask_completed += st.completed ? 1 : -1
  }
}

async function deleteSubtask(id: number) {
  const idx = subtasks.value.findIndex(s => s.id === id)
  if (idx === -1) return
  const wasCompleted = subtasks.value[idx].completed
  subtasks.value.splice(idx, 1)
  props.todo.subtask_count--
  if (wasCompleted) props.todo.subtask_completed--
  try {
    await api.deleteSubtask(id)
  } catch {
    loadSubtasks().then(() => bumpParentCounts())
  }
}

function startEditSubtask(st: Subtask) {
  editingSubtask.value = st.id
  editSubtitle.value = st.title
}

async function saveEditSubtask(id: number) {
  if (!editSubtitle.value.trim()) { editingSubtask.value = 0; return }
  try {
    await api.updateSubtask(id, { title: editSubtitle.value.trim() })
    loadSubtasks().then(() => bumpParentCounts())
  } catch { /* silent */ }
  editingSubtask.value = 0
}

function cancelEditSubtask() {
  editingSubtask.value = 0
}

let dragIndex: number | null = null

function onSubtaskDragstart(e: DragEvent, idx: number) {
  dragIndex = idx
  const el = e.target as HTMLElement
  el.classList.add('opacity-50')
  e.dataTransfer!.effectAllowed = 'move'
}

function onSubtaskDragover(e: DragEvent, _idx: number) {
  e.preventDefault()
  e.dataTransfer!.dropEffect = 'move'
}

function onSubtaskDrop(_e: DragEvent, idx: number) {
  if (dragIndex === null || dragIndex === idx) return
  const items = [...subtasks.value]
  const [moved] = items.splice(dragIndex, 1)
  items.splice(idx, 0, moved)
  subtasks.value = items
  dragIndex = null
  api.reorderSubtasks(items.map(s => s.id))
}

function onSubtaskDragend() {
  dragIndex = null
}

function toggleSubtasks() {
  showSubtasks.value = !showSubtasks.value
  if (showSubtasks.value) loadSubtasks()
}

function startEdit() {
  editTitle.value = props.todo.title
  editDescription.value = props.todo.description
  editPriority.value = props.todo.priority
  editEffort.value = props.todo.effort || ''
  editRecurrence.value = props.todo.recurrence || ''
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
    effort: editEffort.value || undefined,
    recurrence: editRecurrence.value || undefined,
    tags: editTags.value.trim() || undefined,
    due_date: editDueDate.value || null,
  })
  editing.value = false
}

async function handleToggle() {
  const prev = props.todo.completed
  props.todo.completed = !prev
  try {
    const res = await api.toggleTodo(props.todo.id)
    const idx = store.todos.findIndex((t) => t.id === props.todo.id)
    if (idx !== -1) store.todos[idx] = res.data
  } catch {
    props.todo.completed = prev
  }
}
async function handleDelete() { await store.removeTodo(props.todo.id) }

function openDetail() {
  if (isMobile.value) router.push(`/todo/${props.todo.id}`)
}

function formatDate(d: string) { return new Date(d).toLocaleDateString('zh-CN') }
function priorityLabel(p: string) { return { low: '低', medium: '中', high: '高' }[p] || p }
function effortLabel(e: string) { return { easy: '简单', medium: '中等', hard: '困难' }[e] || '' }
function recurrenceLabel(r: string) { return { daily: '每天', weekly: '每周', monthly: '每月' }[r] || '' }
function parseTags(s: string) { return s ? s.split(',').filter(Boolean).map(t => t.trim()) : [] }
</script>

<template>
  <div
    class="bg-white dark:bg-gray-900 rounded-xl border mb-2 shadow-sm transition-colors"
    :class="{
      'border-gray-100 dark:border-gray-800': !isOverdue && !isDueToday,
      'border-red-300 dark:border-red-800 bg-red-50/30 dark:bg-red-950/20': isOverdue,
      'border-amber-300 dark:border-amber-800 bg-amber-50/30 dark:bg-amber-950/20': isDueToday,
      'opacity-60': todo.completed && !editing,
      'border-l-[3px] border-l-gray-300 dark:border-l-gray-600': todo.effort === 'easy',
      'border-l-[3px] border-l-amber-400 dark:border-l-amber-600': todo.effort === 'medium',
      'border-l-[3px] border-l-red-400 dark:border-l-red-600': todo.effort === 'hard',
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
            <span v-if="todo.description" class="text-[12px] text-gray-400 dark:text-gray-500 truncate mt-0.5">{{ todo.description }}</span>
            <div class="flex items-center gap-1.5 mt-0.5 flex-wrap">
              <span v-for="tag in parseTags(todo.tags)" :key="tag"
                class="text-[10px] px-1.5 py-0.5 rounded-md bg-indigo-50 dark:bg-indigo-900/30 text-indigo-600 dark:text-indigo-400">
                {{ tag }}
              </span>
              <span v-if="todo.subtask_count > 0" class="flex items-center gap-1.5">
                <span class="w-12 h-1.5 rounded-full bg-gray-200 dark:bg-gray-700 overflow-hidden">
                  <span class="block h-full rounded-full bg-indigo-400 transition-all" :style="{ width: (todo.subtask_completed / todo.subtask_count * 100) + '%' }" />
                </span>
                <span class="text-[10px] text-gray-400">{{ todo.subtask_completed }}/{{ todo.subtask_count }}</span>
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
          <span v-if="todo.effort" class="text-[11px] font-semibold px-2 py-0.5 rounded-full"
            :class="{
              'bg-gray-100 dark:bg-gray-800 text-gray-600 dark:text-gray-400': todo.effort === 'easy',
              'bg-amber-100 dark:bg-amber-900/40 text-amber-700 dark:text-amber-400': todo.effort === 'medium',
              'bg-red-100 dark:bg-red-900/40 text-red-700 dark:text-red-400': todo.effort === 'hard',
            }">{{ effortLabel(todo.effort) }}</span>
          <span v-if="todo.recurrence" class="text-[11px] text-gray-400" :title="'重复: ' + recurrenceLabel(todo.recurrence)">&#x21bb; {{ recurrenceLabel(todo.recurrence) }}</span>
          <span v-if="todo.due_date" class="text-[11px]"
            :class="{
              'text-red-500 font-semibold': isOverdue,
              'text-amber-500 font-semibold': isDueToday,
              'text-gray-400 dark:text-gray-500': !isOverdue && !isDueToday && !todo.completed,
            }">{{ isOverdue ? '已过期 ' + formatDate(todo.due_date) : isDueToday ? '今天' : formatDate(todo.due_date) }}</span>

          <!-- Subtask toggle -->
          <button @click.stop="toggleSubtasks" class="w-7 h-7 flex items-center justify-center rounded-lg text-gray-300 dark:text-gray-600 hover:text-indigo-500 transition-colors relative"
            :class="{ 'text-indigo-500': showSubtasks }">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="8" y1="6" x2="21" y2="6"/><line x1="8" y1="12" x2="21" y2="12"/><line x1="8" y1="18" x2="21" y2="18"/>
              <line x1="3" y1="6" x2="3.01" y2="6"/><line x1="3" y1="12" x2="3.01" y2="12"/><line x1="3" y1="18" x2="3.01" y2="18"/>
            </svg>
            <span v-if="todo.subtask_count > 0 && !showSubtasks" class="absolute -top-0.5 -right-0.5 min-w-[14px] h-[14px] flex items-center justify-center rounded-full bg-indigo-500 text-white text-[9px] font-bold leading-none px-0.5">{{ todo.subtask_count }}</span>
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
          <Select v-model="editPriority" :options="[
            { label: '低', value: 'low' },
            { label: '中', value: 'medium' },
            { label: '高', value: 'high' },
          ]" />
          <Select v-model="editEffort" :options="[
            { label: '工作量', value: '' },
            { label: '简单', value: 'easy' },
            { label: '中等', value: 'medium' },
            { label: '困难', value: 'hard' },
          ]" />
          <Select v-model="editRecurrence" :options="[
            { label: '不重复', value: '' },
            { label: '每天', value: 'daily' },
            { label: '每周', value: 'weekly' },
            { label: '每月', value: 'monthly' },
          ]" />
          <input v-model="editTags" class="w-[120px] px-2.5 py-2 border border-gray-200 dark:border-gray-700 rounded-lg text-sm bg-white dark:bg-gray-800 dark:text-gray-100 outline-none" placeholder="标签" />
          <button @click="saveEdit" class="px-3 py-2 bg-indigo-500 text-white rounded-lg text-sm font-semibold">保存</button>
          <button @click="cancelEdit" class="px-3 py-2 bg-gray-200 dark:bg-gray-700 text-gray-600 dark:text-gray-300 rounded-lg text-sm font-semibold">取消</button>
        </div>
      </template>
    </div>

    <!-- Subtasks panel -->
    <Transition name="subtask-panel">
      <div v-if="showSubtasks" class="px-4 pb-3 border-t border-gray-50 dark:border-gray-800 pt-2" @click.stop>
        <div v-for="(st, idx) in subtasks" :key="st.id"
          class="flex items-center gap-2 py-1.5 group"
          draggable="true"
          @dragstart="onSubtaskDragstart($event, idx)"
          @dragover="onSubtaskDragover($event, idx)"
          @drop="onSubtaskDrop($event, idx)"
          @dragend="onSubtaskDragend">
          <input type="checkbox" :checked="st.completed" @change="toggleSubtask(st.id)"
            class="w-[15px] h-[15px] accent-indigo-500 cursor-pointer shrink-0" />
        <template v-if="editingSubtask === st.id">
          <input v-model="editSubtitle" @keydown.enter="saveEditSubtask(st.id)" @keydown.escape="cancelEditSubtask()" @blur="saveEditSubtask(st.id)"
            class="flex-1 px-1.5 py-0.5 text-[13px] border border-indigo-400 rounded bg-white dark:bg-gray-800 dark:text-gray-100 outline-none" />
        </template>
        <template v-else>
          <span @dblclick="startEditSubtask(st)" class="text-[13px] flex-1 cursor-default" :class="{ 'line-through text-gray-300 dark:text-gray-600': st.completed, 'text-gray-600 dark:text-gray-400': !st.completed }">{{ st.title }}</span>
          <button @click="startEditSubtask(st)" class="opacity-0 group-hover:opacity-100 text-gray-300 dark:text-gray-600 hover:text-indigo-500 shrink-0 transition-opacity">
            <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/></svg>
          </button>
        </template>
        <button @click="deleteSubtask(st.id)" class="text-gray-300 dark:text-gray-700 hover:text-red-400 shrink-0">
          <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
        </button>
      </div>
      <form @submit.prevent="addSubtask" class="flex gap-1.5 mt-1">
        <input v-model="subInput" class="flex-1 px-2.5 py-1.5 border border-gray-150 dark:border-gray-800 rounded-lg text-xs bg-gray-50 dark:bg-gray-800 dark:text-gray-200 outline-none" placeholder="添加子任务…" />
        <button type="submit" class="px-3 py-1.5 bg-indigo-500 text-white rounded-lg text-xs font-semibold">添加</button>
      </form>
      </div>
    </Transition>
  </div>

  <ConfirmDialog
    :open="showConfirm"
    title="删除待办"
    :message="`确定要删除「${todo.title}」吗？`"
    @confirm="handleDelete(); showConfirm = false"
    @cancel="showConfirm = false"
  />
</template>

<style scoped>
.subtask-panel-enter-active,
.subtask-panel-leave-active {
  transition: all 0.2s ease;
  overflow: hidden;
}
.subtask-panel-enter-from,
.subtask-panel-leave-to {
  opacity: 0;
  max-height: 0;
}
.subtask-panel-enter-to,
.subtask-panel-leave-from {
  opacity: 1;
  max-height: 500px;
}
</style>
