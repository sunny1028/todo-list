<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'
import AppHeader from './components/layout/AppHeader.vue'
import AppNav from './components/layout/AppNav.vue'
import Sidebar from './components/layout/Sidebar.vue'
import ToastContainer from './components/ui/ToastContainer.vue'
import ShortcutPanel from './components/ui/ShortcutPanel.vue'
import { useNotifications } from './composables/useNotifications'
import { useAuth } from './stores/auth'
import { useRouter } from 'vue-router'

useNotifications()

const auth = useAuth()
const router = useRouter()
const showBanner = computed(() => auth.initialized && !auth.hasPassword)

const showShortcuts = ref(false)

function onKeydown(e: KeyboardEvent) {
  if (e.key === '?' && !isInputFocused()) {
    e.preventDefault()
    showShortcuts.value = !showShortcuts.value
  }
}

function isInputFocused() {
  const el = document.activeElement
  return el instanceof HTMLInputElement || el instanceof HTMLTextAreaElement || el instanceof HTMLSelectElement
}

onMounted(() => window.addEventListener('keydown', onKeydown))
onUnmounted(() => window.removeEventListener('keydown', onKeydown))
</script>

<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-950 transition-colors">
    <AppHeader />
    <div v-if="showBanner" class="bg-amber-50 dark:bg-amber-950 border-b border-amber-200 dark:border-amber-800 px-4 py-2 text-sm text-amber-800 dark:text-amber-200 flex items-center justify-between">
      <span>匿名使用中，换浏览器或清缓存将丢失数据。</span>
      <button @click="router.push('/login')" class="underline font-medium hover:text-amber-600 dark:hover:text-amber-300 ml-2 whitespace-nowrap">设置账号</button>
    </div>
    <div class="flex">
      <Sidebar />
      <main class="flex-1 mx-auto max-w-2xl px-4 py-6 pb-20 md:pb-6">
        <RouterView />
      </main>
    </div>
    <AppNav />
    <ToastContainer />
    <ShortcutPanel v-if="showShortcuts" @close="showShortcuts = false" />
  </div>
</template>
