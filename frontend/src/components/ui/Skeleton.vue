<script setup lang="ts">
defineProps<{ count?: number }>()
</script>

<template>
  <div v-for="i in (count || 4)" :key="i" class="flex items-center gap-3 p-3.5 bg-white dark:bg-gray-900 rounded-xl border border-gray-100 dark:border-gray-800 mb-2 animate-pulse">
    <div class="w-[18px] h-[18px] rounded bg-gray-200 dark:bg-gray-700 shrink-0" />
    <div class="flex flex-col gap-2 flex-1">
      <div class="h-4 bg-gray-200 dark:bg-gray-700 rounded w-3/4" />
      <div class="h-3 bg-gray-100 dark:bg-gray-800 rounded w-1/2" />
    </div>
    <div class="w-10 h-5 bg-gray-100 dark:bg-gray-800 rounded-full" />
  </div>
</template>
