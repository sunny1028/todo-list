import { ref } from 'vue'

const isDark = ref(false)

function apply() {
  document.documentElement.classList.toggle('dark', isDark.value)
}

export function useDarkMode() {
  // Init from localStorage or system preference
  const saved = localStorage.getItem('theme')
  if (saved === 'dark') isDark.value = true
  else if (saved === 'light') isDark.value = false
  else isDark.value = window.matchMedia('(prefers-color-scheme: dark)').matches

  apply()

  // Listen for system changes (only when no manual override)
  window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', (e) => {
    if (!localStorage.getItem('theme')) {
      isDark.value = e.matches
      apply()
    }
  })

  function toggle() {
    isDark.value = !isDark.value
    localStorage.setItem('theme', isDark.value ? 'dark' : 'light')
    apply()
  }

  return { isDark, toggle }
}
