import test from 'node:test'
import assert from 'node:assert/strict'
import { createPageViewRecorder } from './viewRecorder.js'

test('records each content item only once per page lifecycle', async () => {
  const calls = []
  const record = createPageViewRecorder(async (contentType, contentID) => {
    calls.push([contentType, contentID])
    return { data: { counted: true, view_count: 8 } }
  })

  const first = record('daily_question', 3)
  const duplicate = record('daily_question', 3)
  const other = record('daily_question', 4)

  assert.equal((await first).data.view_count, 8)
  assert.equal(await duplicate, null)
  await other
  assert.deepEqual(calls, [['daily_question', 3], ['daily_question', 4]])
})

test('allows retry after a failed recording request', async () => {
  let calls = 0
  const record = createPageViewRecorder(async () => {
    calls += 1
    if (calls === 1) throw new Error('network')
    return { data: { counted: true, view_count: 2 } }
  })

  await assert.rejects(record('article', 1), /network/)
  const retried = await record('article', 1)

  assert.equal(retried.data.view_count, 2)
  assert.equal(calls, 2)
})
