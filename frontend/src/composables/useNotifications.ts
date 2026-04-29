import { onMounted, onUnmounted } from 'vue'

export function useNotifications() {
  let interval: ReturnType<typeof setInterval> | null = null

  function requestPermission() {
    if ('Notification' in window && Notification.permission === 'default') {
      Notification.requestPermission()
    }
  }

  async function checkDueDates() {
    if (!('Notification' in window) || Notification.permission !== 'granted') return

    try {
      // Fetch fresh data
      const res = await import('../api/todos').then(m => m.listTodos())
      const todos = res.data.filter(t => !t.completed && !t.archived && t.due_date)
      const now = new Date()
      const today = now.toISOString().slice(0, 10)

      for (const todo of todos) {
        if (todo.due_date === today && !notificationSent(todo.id, 'today')) {
          new Notification('今天截止', { body: todo.title, icon: '/favicon.svg', tag: `today-${todo.id}` })
          markSent(todo.id, 'today')
        }
        if (todo.due_date && todo.due_date < today && !notificationSent(todo.id, 'overdue')) {
          new Notification('已过期', { body: todo.title, icon: '/favicon.svg', tag: `overdue-${todo.id}` })
          markSent(todo.id, 'overdue')
        }
      }
    } catch { /* silent */ }
  }

  function notificationSent(id: number, type: string): boolean {
    const key = `notif-${type}-${id}`
    const sent = localStorage.getItem(key)
    if (!sent) return false
    const sentDate = new Date(sent).toISOString().slice(0, 10)
    return sentDate === new Date().toISOString().slice(0, 10)
  }

  function markSent(id: number, type: string) {
    localStorage.setItem(`notif-${type}-${id}`, new Date().toISOString())
  }

  onMounted(() => {
    requestPermission()
    checkDueDates()
    interval = setInterval(checkDueDates, 300000) // every 5 min
    document.addEventListener('visibilitychange', () => {
      if (document.visibilityState === 'visible') checkDueDates()
    })
  })

  onUnmounted(() => {
    if (interval) clearInterval(interval)
  })
}
