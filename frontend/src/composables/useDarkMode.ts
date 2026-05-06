import { ref } from 'vue'

const isDark = ref(false)
const mode = ref<'light' | 'dark' | 'auto'>('auto')

function apply() {
  document.documentElement.classList.toggle('dark', isDark.value)
}

export function useDarkMode() {
  const saved = localStorage.getItem('theme') as 'light' | 'dark' | null

  if (saved === 'dark') {
    mode.value = 'dark'
    isDark.value = true
  } else if (saved === 'light') {
    mode.value = 'light'
    isDark.value = false
  } else {
    mode.value = 'auto'
    isDark.value = window.matchMedia('(prefers-color-scheme: dark)').matches
  }

  apply()

  window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', (e) => {
    if (mode.value === 'auto') {
      isDark.value = e.matches
      apply()
    }
  })

  function toggle() {
    if (mode.value === 'light') {
      mode.value = 'dark'
      isDark.value = true
      localStorage.setItem('theme', 'dark')
    } else if (mode.value === 'dark') {
      mode.value = 'auto'
      localStorage.removeItem('theme')
      isDark.value = window.matchMedia('(prefers-color-scheme: dark)').matches
    } else {
      mode.value = 'light'
      isDark.value = false
      localStorage.setItem('theme', 'light')
    }
    apply()
  }

  return { isDark, mode, toggle }
}
