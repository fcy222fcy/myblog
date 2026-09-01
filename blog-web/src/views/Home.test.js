import test from 'node:test'
import assert from 'node:assert/strict'
import { readFileSync } from 'node:fs'

const source = readFileSync(new URL('./Home.vue', import.meta.url), 'utf8')

test('keeps the home article summary at three lines with overflow ellipsis', () => {
  assert.match(source, /<p class="article-summary">\{\{ displaySummary\(article\) \}\}<\/p>/)
  assert.match(source, /\.article-summary\s*\{[\s\S]*height:\s*4\.95em;/)
  assert.match(source, /\.article-summary\s*\{[\s\S]*-webkit-line-clamp:\s*3;/)
  assert.match(source, /\.article-summary\s*\{[\s\S]*overflow:\s*hidden;/)
})
