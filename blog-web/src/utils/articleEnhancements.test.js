import test from 'node:test'
import assert from 'node:assert/strict'
import { renderArticleMarkdown } from './articleEnhancements.js'

test('builds a nested table of contents with stable unique heading ids', () => {
  const result = renderArticleMarkdown(`## 安装环境

### Windows

## 安装环境`)

  assert.deepEqual(result.toc, [
    { id: '安装环境', text: '安装环境', level: 2, number: '1.' },
    { id: 'windows', text: 'Windows', level: 3, number: '1.1.' },
    { id: '安装环境-2', text: '安装环境', level: 2, number: '2.' }
  ])
  assert.match(result.html, /<h2 id="安装环境">安装环境<\/h2>/)
  assert.match(result.html, /<h3 id="windows">Windows<\/h3>/)
  assert.match(result.html, /<h2 id="安装环境-2">安装环境<\/h2>/)
})

test('renders fenced code with language label copy action and line numbers', () => {
  const result = renderArticleMarkdown('```javascript\nconst answer = 42\nconsole.log(answer)\n```')

  assert.match(result.html, /class="code-block"/)
  assert.match(result.html, /class="code-language">JavaScript</)
  assert.match(result.html, /data-copy-code/)
  assert.match(result.html, /data-line="1"/)
  assert.match(result.html, /data-line="2"/)
  assert.match(result.html, /class="hljs language-javascript"/)
})

test('renders article images lazily with an accessible failure placeholder', () => {
  const result = renderArticleMarkdown('![架构图](images/architecture.png "系统架构")')

  assert.match(result.html, /class="article-image-frame"/)
  assert.match(result.html, /loading="lazy"/)
  assert.match(result.html, /decoding="async"/)
  assert.match(result.html, /alt="架构图"/)
  assert.match(result.html, /class="article-image-fallback"/)
  assert.match(result.html, />图片暂时无法显示</)
  assert.doesNotMatch(result.html, /图片加载失败：架构图/)
})
