<script setup lang="ts">
import { useToast } from '../../stores/toast'

const toast = useToast()
</script>

<template>
  <TransitionGroup
    name="toast"
    tag="div"
    class="fixed top-4 right-4 z-50 flex flex-col gap-2 pointer-events-none"
  >
    <div
      v-for="t in toast.toasts"
      :key="t.id"
      class="pointer-events-auto px-4 py-2.5 rounded-xl text-sm font-medium shadow-lg backdrop-blur-sm transition-all flex items-center gap-2"
      :class="{
        'bg-emerald-500/90 text-white': t.type === 'success',
        'bg-red-500/90 text-white': t.type === 'error',
        'bg-gray-800/90 text-white dark:bg-white/90 dark:text-gray-800': t.type === 'info',
      }"
    >
      <span>{{ t.message }}</span>
      <button v-if="t.action" @click="t.action.onClick" class="underline font-bold hover:opacity-80 shrink-0">{{ t.action.label }}</button>
    </div>
  </TransitionGroup>
</template>

<style scoped>
.toast-enter-active { transition: all 0.3s ease; }
.toast-leave-active { transition: all 0.2s ease; }
.toast-enter-from { opacity: 0; transform: translateX(40px); }
.toast-leave-to { opacity: 0; transform: translateX(40px); }
</style>
