import test from 'node:test'
import assert from 'node:assert/strict'
import { getOrCreateVisitorID } from './visitorIdentity.js'

const createStorage = () => {
  const values = new Map()
  return {
    getItem: key => values.get(key) || null,
    setItem: (key, value) => values.set(key, value)
  }
}

test('creates and reuses one anonymous visitor id', () => {
  const storage = createStorage()
  let calls = 0
  const cryptoProvider = {
    randomUUID: () => {
      calls += 1
      return '550e8400-e29b-41d4-a716-446655440000'
    }
  }

  const first = getOrCreateVisitorID(storage, cryptoProvider)
  const second = getOrCreateVisitorID(storage, cryptoProvider)

  assert.equal(first, '550e8400-e29b-41d4-a716-446655440000')
  assert.equal(second, first)
  assert.equal(calls, 1)
})

test('returns an empty id when browser storage is unavailable', () => {
  const brokenStorage = {
    getItem: () => { throw new Error('blocked') },
    setItem: () => { throw new Error('blocked') }
  }

  assert.equal(getOrCreateVisitorID(brokenStorage, {}), '')
})
