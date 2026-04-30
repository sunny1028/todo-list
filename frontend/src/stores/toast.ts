import { defineStore } from 'pinia'
import { ref } from 'vue'

export interface Toast {
  id: number
  message: string
  type: 'success' | 'error' | 'info'
  action?: { label: string; onClick: () => void }
  duration?: number
}

let nextId = 0

export const useToast = defineStore('toast', () => {
  const toasts = ref<Toast[]>([])

  function show(message: string, type: Toast['type'] = 'info', action?: Toast['action'], duration?: number) {
    const id = ++nextId
    const d = duration ?? (action ? 10000 : 3000)
    toasts.value.push({ id, message, type, action, duration })
    setTimeout(() => remove(id), d)
  }

  function remove(id: number) {
    toasts.value = toasts.value.filter((t) => t.id !== id)
  }

  return { toasts, show, remove }
})
