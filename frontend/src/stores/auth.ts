import { defineStore } from 'pinia'
import { ref } from 'vue'
import axios from 'axios'

export const useAuth = defineStore('auth', () => {
  const uuid = ref('')
  const token = ref('')
  const hasPassword = ref(false)
  const initialized = ref(false)
  const userId = ref(0)

  function generateUUID(): string {
    return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, (c) => {
      const r = (Math.random() * 16) | 0
      return (c === 'x' ? r : (r & 0x3) | 0x8).toString(16)
    })
  }

  async function init() {
    let storedUuid = localStorage.getItem('uuid')
    if (!storedUuid) {
      storedUuid = generateUUID()
      localStorage.setItem('uuid', storedUuid)
    }
    uuid.value = storedUuid

    try {
      const res = await axios.post('/api/auth/uuid', { uuid: uuid.value })
      token.value = res.data.token
      hasPassword.value = res.data.has_password
      localStorage.setItem('auth_token', token.value)
      try {
        const payload = JSON.parse(atob(token.value.split('.')[1]))
        userId.value = payload.user_id || 0
      } catch { /* ignore */ }
    } catch {
      // Degraded mode — API calls will return 401
    } finally {
      initialized.value = true
    }
  }

  async function bind(username: string, password: string) {
    const res = await axios.post('/api/auth/bind', { username, password }, {
      headers: { Authorization: `Bearer ${token.value}` },
    })
    hasPassword.value = true
    return res.data
  }

  async function login(username: string, password: string) {
    const res = await axios.post('/api/auth/login', { username, password })
    token.value = res.data.token
    hasPassword.value = res.data.has_password
    localStorage.setItem('auth_token', token.value)
    uuid.value = '' // will be set from JWT on next init, but we need it now
    // Parse uuid from token payload
    try {
      const payload = JSON.parse(atob(res.data.token.split('.')[1]))
      uuid.value = payload.uuid
      userId.value = payload.user_id || 0
      localStorage.setItem('uuid', payload.uuid)
    } catch { /* ignore */ }
    return res.data
  }

  async function mergeLogin(username: string, password: string) {
    const res = await axios.post('/api/auth/merge', { username, password }, {
      headers: { Authorization: `Bearer ${token.value}` },
    })
    token.value = res.data.token
    hasPassword.value = res.data.has_password
    localStorage.setItem('auth_token', token.value)
    try {
      const payload = JSON.parse(atob(res.data.token.split('.')[1]))
      uuid.value = payload.uuid
      userId.value = payload.user_id || 0
      localStorage.setItem('uuid', payload.uuid)
    } catch { /* ignore */ }
    return res.data
  }

  async function logout() {
    localStorage.removeItem('auth_token')
    localStorage.removeItem('uuid')
    token.value = ''
    hasPassword.value = false
    uuid.value = ''
    await init()
  }

  return { uuid, token, hasPassword, initialized, userId, init, bind, login, mergeLogin, logout }
})
