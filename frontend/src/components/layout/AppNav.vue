<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()

const tabs = [
  { path: '/', label: '全部' },
  { path: '/?filter=active', label: '未完成' },
  { path: '/?filter=completed', label: '已完成' },
  { path: '/stats', label: '统计' },
]

function isActive(path: string) {
  if (path === '/stats') return route.path === '/stats'
  if (path === '/') return !route.query.filter && route.path === '/'
  return path.includes(`filter=${route.query.filter}`)
}
</script>

<template>
  <nav class="fixed bottom-0 left-0 right-0 bg-white dark:bg-gray-900 border-t border-gray-200 dark:border-gray-800 flex z-10 md:hidden transition-colors">
    <button
      v-for="tab in tabs"
      :key="tab.label"
      @click="router.push(tab.path)"
      class="flex-1 py-3 text-center text-sm font-medium transition-colors"
      :class="isActive(tab.path) ? 'text-indigo-600 dark:text-indigo-400' : 'text-gray-400 dark:text-gray-600'"
    >
      {{ tab.label }}
    </button>
  </nav>
</template>
