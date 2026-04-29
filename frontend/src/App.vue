<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'
import AppHeader from './components/layout/AppHeader.vue'
import AppNav from './components/layout/AppNav.vue'
import Sidebar from './components/layout/Sidebar.vue'
import ToastContainer from './components/ui/ToastContainer.vue'
import ShortcutPanel from './components/ui/ShortcutPanel.vue'
import { useNotifications } from './composables/useNotifications'

useNotifications()

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
