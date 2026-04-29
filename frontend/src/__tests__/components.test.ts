import { describe, it, expect, beforeEach, afterEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import ConfirmDialog from '../components/ui/ConfirmDialog.vue'
import Skeleton from '../components/ui/Skeleton.vue'

beforeEach(() => {
  setActivePinia(createPinia())
})

afterEach(() => {
  document.body.innerHTML = ''
})

describe('ConfirmDialog', () => {
  it('renders when open', () => {
    mount(ConfirmDialog, {
      props: { open: true, title: '确认', message: '确定要删除吗？' },
    })
    expect(document.body.textContent).toContain('确认')
  })

  it('does not render when closed', () => {
    mount(ConfirmDialog, {
      props: { open: false, title: '确认', message: '确定？' },
    })
    expect(document.body.textContent).not.toContain('确认')
  })

  it('emits confirm on delete button click', async () => {
    const wrapper = mount(ConfirmDialog, {
      props: { open: true, title: '确认', message: '确定？' },
    })
    const buttons = document.body.querySelectorAll('button')
    ;(buttons[buttons.length - 1] as HTMLButtonElement).click()
    expect(wrapper.emitted('confirm')).toBeTruthy()
  })

  it('emits cancel on cancel button click', async () => {
    const wrapper = mount(ConfirmDialog, {
      props: { open: true, title: '确认', message: '确定？' },
    })
    const buttons = document.body.querySelectorAll('button')
    ;(buttons[buttons.length - 2] as HTMLButtonElement).click()
    expect(wrapper.emitted('cancel')).toBeTruthy()
  })
})

describe('Skeleton', () => {
  it('renders correct count', () => {
    const wrapper = mount(Skeleton, { props: { count: 3 } })
    expect(wrapper.findAll('.animate-pulse').length).toBe(3)
  })

  it('defaults to 4 items', () => {
    const wrapper = mount(Skeleton)
    expect(wrapper.findAll('.animate-pulse').length).toBe(4)
  })
})
