<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useLists } from '../../stores/lists'
import { useTodos } from '../../stores/todo'

const router = useRouter()
const route = useRoute()
const listStore = useLists()
const todoStore = useTodos()

const showNewList = ref(false)
const newListName = ref('')

const colors = ['#6366f1', '#ec4899', '#f59e0b', '#10b981', '#3b82f6', '#8b5cf6', '#ef4444', '#06b6d4']

function selectList(id: number) {
  todoStore.setList(id)
  todoStore.fetchTodos()
  if (route.path !== '/') router.push('/')
}

function isActive(id: number) {
  return todoStore.currentListId === id
}

async function createList() {
  if (!newListName.value.trim()) return
  const color = colors[listStore.lists.length % colors.length]
  await listStore.addList(newListName.value.trim(), color)
  newListName.value = ''
  showNewList.value = false
}

async function deleteList(id: number) {
  await listStore.removeList(id)
  if (todoStore.currentListId === id) {
    selectList(0)
  }
}

onMounted(() => {
  listStore.fetchLists()
})
</script>

<template>
  <aside class="hidden md:flex md:flex-col md:w-56 md:shrink-0 md:border-r md:border-gray-200 dark:md:border-gray-800 md:h-[calc(100vh-57px)] md:sticky md:top-[57px] md:bg-white dark:md:bg-gray-900 md:px-3 md:py-4 md:overflow-y-auto transition-colors">
    <div class="flex items-center justify-between mb-3 px-2">
      <span class="text-xs font-semibold text-gray-400 dark:text-gray-500 uppercase tracking-wider">列表</span>
      <button @click="showNewList = !showNewList" class="text-gray-400 hover:text-indigo-500 transition-colors">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
      </button>
    </div>

    <!-- New list form -->
    <form v-if="showNewList" @submit.prevent="createList" class="flex gap-1 mb-2 px-2">
      <input v-model="newListName" placeholder="列表名称" class="flex-1 px-2 py-1.5 text-xs border border-gray-200 dark:border-gray-700 rounded-lg bg-white dark:bg-gray-800 dark:text-gray-100 outline-none" />
      <button type="submit" class="px-2 py-1 text-xs bg-indigo-500 text-white rounded-lg">创建</button>
    </form>

    <!-- All todos -->
    <button
      @click="selectList(0)"
      class="w-full text-left px-3 py-2 rounded-xl text-sm font-medium mb-0.5 transition-colors"
      :class="isActive(0) ? 'bg-indigo-50 dark:bg-indigo-900/30 text-indigo-700 dark:text-indigo-300' : 'text-gray-600 dark:text-gray-400 hover:bg-gray-50 dark:hover:bg-gray-800'"
    >
      全部待办
    </button>

    <!-- Custom lists -->
    <button
      v-for="list in listStore.lists"
      :key="list.id"
      @click="selectList(list.id)"
      class="w-full text-left px-3 py-2 rounded-xl text-sm font-medium mb-0.5 transition-colors group flex items-center justify-between"
      :class="isActive(list.id) ? 'bg-indigo-50 dark:bg-indigo-900/30 text-indigo-700 dark:text-indigo-300' : 'text-gray-600 dark:text-gray-400 hover:bg-gray-50 dark:hover:bg-gray-800'"
    >
      <span class="flex items-center gap-2">
        <span class="w-2.5 h-2.5 rounded-full shrink-0" :style="{ background: list.color }" />
        {{ list.name }}
      </span>
      <button
        @click.stop="deleteList(list.id)"
        class="opacity-0 group-hover:opacity-100 text-gray-400 hover:text-red-500 transition-all"
      >
        <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
      </button>
    </button>

    <!-- Stats link -->
    <div class="mt-auto pt-4 border-t border-gray-100 dark:border-gray-800 mt-4">
      <button
        @click="router.push('/stats')"
        class="w-full text-left px-3 py-2 rounded-xl text-sm font-medium transition-colors flex items-center gap-2"
        :class="route.path === '/stats' ? 'bg-indigo-50 dark:bg-indigo-900/30 text-indigo-700 dark:text-indigo-300' : 'text-gray-600 dark:text-gray-400 hover:bg-gray-50 dark:hover:bg-gray-800'"
      >
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="20" x2="18" y2="10"/><line x1="12" y1="20" x2="12" y2="4"/><line x1="6" y1="20" x2="6" y2="14"/></svg>
        统计
      </button>
    </div>
  </aside>
</template>
