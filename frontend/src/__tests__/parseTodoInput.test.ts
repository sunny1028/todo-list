import { describe, it, expect, beforeEach, afterEach } from 'vitest'
import { parseTodoInput } from '../utils/parseTodoInput'

function todayStr(): string {
  const d = new Date()
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
}
function tomorrowStr(): string {
  const d = new Date(Date.now() + 86400000)
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
}
function dayAfterTomorrowStr(): string {
  const d = new Date(Date.now() + 172800000)
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
}

describe('parseTodoInput', () => {
  it('empty string returns empty title with null fields', () => {
    const r = parseTodoInput('')
    expect(r.title).toBe('')
    expect(r.dueDate).toBeNull()
    expect(r.priority).toBeNull()
    expect(r.tags).toBeNull()
    expect(r.effort).toBeNull()
    expect(r.recurrence).toBeNull()
  })

  it('plain text returns as title', () => {
    const r = parseTodoInput('买牛奶')
    expect(r.title).toBe('买牛奶')
    expect(r.priority).toBeNull()
    expect(r.tags).toBeNull()
  })

  // Tags
  it('extracts single tag', () => {
    const r = parseTodoInput('#购物')
    expect(r.tags).toBe('购物')
    expect(r.title).toBe('')
  })

  it('extracts multiple tags', () => {
    const r = parseTodoInput('#购物 #工作')
    expect(r.tags).toBe('购物,工作')
  })

  it('tag in context', () => {
    const r = parseTodoInput('买牛奶 #购物')
    expect(r.title).toBe('买牛奶')
    expect(r.tags).toBe('购物')
  })

  it('mid-word # is not parsed as tag', () => {
    const r = parseTodoInput('买牛奶#购物')
    expect(r.title).toBe('买牛奶#购物')
    expect(r.tags).toBeNull()
  })

  // Priority
  it('extracts !高 as high', () => {
    expect(parseTodoInput('!高').priority).toBe('high')
  })
  it('extracts !中 as medium', () => {
    expect(parseTodoInput('!中').priority).toBe('medium')
  })
  it('extracts !低 as low', () => {
    expect(parseTodoInput('!低').priority).toBe('low')
  })
  it('extracts !h as high', () => {
    expect(parseTodoInput('!h').priority).toBe('high')
  })
  it('extracts !1 as high', () => {
    expect(parseTodoInput('!1').priority).toBe('high')
  })
  it('last priority wins', () => {
    expect(parseTodoInput('!高 !低').priority).toBe('low')
  })
  it('mid-word ! is not parsed', () => {
    const r = parseTodoInput('买牛奶!高')
    expect(r.priority).toBeNull()
    expect(r.title).toBe('买牛奶!高')
  })

  // Effort
  it('extracts @简单 as easy', () => {
    expect(parseTodoInput('@简单').effort).toBe('easy')
  })
  it('extracts @easy as easy', () => {
    expect(parseTodoInput('@easy').effort).toBe('easy')
  })
  it('extracts @困难 as hard', () => {
    expect(parseTodoInput('@困难').effort).toBe('hard')
  })

  // Recurrence
  it('extracts *每天 as daily', () => {
    expect(parseTodoInput('*每天').recurrence).toBe('daily')
  })
  it('extracts *weekly as weekly', () => {
    expect(parseTodoInput('*weekly').recurrence).toBe('weekly')
  })
  it('extracts *每月 as monthly', () => {
    expect(parseTodoInput('*每月').recurrence).toBe('monthly')
  })

  // Dates
  it('resolves 今天', () => {
    expect(parseTodoInput('今天').dueDate).toBe(todayStr())
  })
  it('resolves 明天', () => {
    expect(parseTodoInput('明天').dueDate).toBe(tomorrowStr())
  })
  it('resolves 后天', () => {
    expect(parseTodoInput('后天').dueDate).toBe(dayAfterTomorrowStr())
  })
  it('resolves YYYY-MM-DD', () => {
    expect(parseTodoInput('2026-05-08').dueDate).toBe('2026-05-08')
  })
  it('resolves M/D', () => {
    const r = parseTodoInput('5/8')
    const y = new Date().getFullYear()
    expect(r.dueDate).toBe(`${y}-05-08`)
  })
  it('resolves M月D日', () => {
    const r = parseTodoInput('5月8日')
    const y = new Date().getFullYear()
    expect(r.dueDate).toBe(`${y}-05-08`)
  })

  // Combined
  it('parses full expression', () => {
    const r = parseTodoInput('明天 买牛奶 #购物 !高')
    expect(r.title).toBe('买牛奶')
    expect(r.dueDate).toBe(tomorrowStr())
    expect(r.priority).toBe('high')
    expect(r.tags).toBe('购物')
  })

  it('parses with all markers', () => {
    const r = parseTodoInput('明天 写报告 #工作 !高 @困难 *每周')
    expect(r.title).toBe('写报告')
    expect(r.dueDate).toBe(tomorrowStr())
    expect(r.priority).toBe('high')
    expect(r.tags).toBe('工作')
    expect(r.effort).toBe('hard')
    expect(r.recurrence).toBe('weekly')
  })

  it('only markers no title yields empty title', () => {
    const r = parseTodoInput('!高 #购物')
    expect(r.title).toBe('')
  })

  it('normalizes whitespace', () => {
    const r = parseTodoInput('  买牛奶    !高  #购物  ')
    expect(r.title).toBe('买牛奶')
    expect(r.priority).toBe('high')
    expect(r.tags).toBe('购物')
  })
})
