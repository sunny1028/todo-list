<script setup lang="ts">
defineProps<{
  open: boolean
  title: string
  message: string
}>()

const emit = defineEmits<{
  confirm: []
  cancel: []
}>()
</script>

<template>
  <Teleport to="body">
    <div
      v-if="open"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/30 backdrop-blur-sm"
      @click.self="emit('cancel')"
    >
      <div class="bg-white dark:bg-gray-800 rounded-2xl shadow-2xl p-6 w-[90%] max-w-sm mx-auto">
        <h3 class="text-lg font-bold text-gray-800 dark:text-gray-100 mb-2">{{ title }}</h3>
        <p class="text-sm text-gray-500 dark:text-gray-400 mb-6">{{ message }}</p>
        <div class="flex gap-3 justify-end">
          <button
            @click="emit('cancel')"
            class="px-4 py-2 text-sm font-semibold rounded-xl bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-300 hover:bg-gray-200 dark:hover:bg-gray-600 transition-colors"
          >
            取消
          </button>
          <button
            @click="emit('confirm')"
            class="px-4 py-2 text-sm font-semibold rounded-xl bg-red-500 text-white hover:bg-red-600 transition-colors"
          >
            删除
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>
