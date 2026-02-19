import { describe, it, expect } from 'vitest'
import { formatPrice, cn } from '@/lib/utils'

describe('formatPrice', () => {
  it('formats a number as USD currency', () => {
    expect(formatPrice(9.99)).toBe('$9.99')
  })

  it('formats zero correctly', () => {
    expect(formatPrice(0)).toBe('$0.00')
  })

  it('formats large numbers with commas', () => {
    expect(formatPrice(1234.56)).toBe('$1,234.56')
  })
})

describe('cn', () => {
  it('merges class names', () => {
    expect(cn('foo', 'bar')).toBe('foo bar')
  })

  it('handles conditional classes', () => {
    expect(cn('foo', false && 'bar', 'baz')).toBe('foo baz')
  })
})
