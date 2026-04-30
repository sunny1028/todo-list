<script setup lang="ts" generic="T extends string">
import { computed, ref } from 'vue'

const props = withDefaults(defineProps<{
  modelValue: string
  options: { label: string; value: T }[]
  placeholder?: string
}>(), {
  placeholder: '请选择',
})

const emit = defineEmits<{ 'update:modelValue': [value: string] }>()

const open = ref(false)

const selectedLabel = computed(() => {
  const opt = props.options.find(o => o.value === props.modelValue)
  return opt?.label ?? props.placeholder
})

function select(value: string) {
  emit('update:modelValue', value)
  open.value = false
}

function toggle() {
  open.value = !open.value
}
</script>

<template>
  <div class="relative">
    <button type="button" @click="toggle"
      class="flex items-center gap-1.5 px-3 py-2 border border-gray-200 dark:border-gray-800 rounded-lg text-sm bg-white dark:bg-gray-900 dark:text-gray-100 outline-none hover:border-indigo-400 transition-colors text-left whitespace-nowrap">
      <span :class="{ 'text-gray-400 dark:text-gray-500': !modelValue }">{{ selectedLabel }}</span>
      <svg class="shrink-0 text-gray-400 dark:text-gray-500 transition-transform" :class="{ 'rotate-180': open }" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="6 9 12 15 18 9"/></svg>
    </button>
    <div v-if="open" class="absolute left-0 top-full mt-1 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl shadow-lg z-20 py-1 min-w-full">
      <button v-for="opt in options" :key="opt.value" type="button" @click="select(opt.value)"
        class="w-full text-left px-3 py-2 text-sm transition-colors whitespace-nowrap"
        :class="modelValue === opt.value ? 'text-indigo-600 dark:text-indigo-400 bg-indigo-50 dark:bg-indigo-900/30' : 'text-gray-600 dark:text-gray-400 hover:bg-gray-50 dark:hover:bg-gray-700/50'">
        {{ opt.label }}
      </button>
    </div>
    <div v-if="open" class="fixed inset-0 z-10" @click="open = false" />
  </div>
</template>
