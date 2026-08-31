import { marked } from 'marked'
import hljs from 'highlight.js/lib/core'
import bash from 'highlight.js/lib/languages/bash'
import cpp from 'highlight.js/lib/languages/cpp'
import css from 'highlight.js/lib/languages/css'
import dockerfile from 'highlight.js/lib/languages/dockerfile'
import go from 'highlight.js/lib/languages/go'
import java from 'highlight.js/lib/languages/java'
import javascript from 'highlight.js/lib/languages/javascript'
import json from 'highlight.js/lib/languages/json'
import markdown from 'highlight.js/lib/languages/markdown'
import nginx from 'highlight.js/lib/languages/nginx'
import powershell from 'highlight.js/lib/languages/powershell'
import python from 'highlight.js/lib/languages/python'
import sql from 'highlight.js/lib/languages/sql'
import typescript from 'highlight.js/lib/languages/typescript'
import xml from 'highlight.js/lib/languages/xml'
import yaml from 'highlight.js/lib/languages/yaml'

hljs.registerLanguage('bash', bash)
hljs.registerLanguage('cpp', cpp)
hljs.registerLanguage('css', css)
hljs.registerLanguage('dockerfile', dockerfile)
hljs.registerLanguage('go', go)
hljs.registerLanguage('java', java)
hljs.registerLanguage('javascript', javascript)
hljs.registerLanguage('json', json)
hljs.registerLanguage('markdown', markdown)
hljs.registerLanguage('nginx', nginx)
hljs.registerLanguage('powershell', powershell)
hljs.registerLanguage('python', python)
hljs.registerLanguage('sql', sql)
hljs.registerLanguage('typescript', typescript)
hljs.registerLanguage('xml', xml)
hljs.registerLanguage('yaml', yaml)

hljs.registerAliases(['sh', 'shell'], { languageName: 'bash' })
hljs.registerAliases(['c'], { languageName: 'cpp' })
hljs.registerAliases(['html', 'vue'], { languageName: 'xml' })
hljs.registerAliases(['js'], { languageName: 'javascript' })
hljs.registerAliases(['md'], { languageName: 'markdown' })
hljs.registerAliases(['ps1'], { languageName: 'powershell' })
hljs.registerAliases(['py'], { languageName: 'python' })
hljs.registerAliases(['ts'], { languageName: 'typescript' })
hljs.registerAliases(['yml'], { languageName: 'yaml' })

const LANGUAGE_NAMES = {
  bash: 'Bash',
  c: 'C',
  cpp: 'C++',
  css: 'CSS',
  dockerfile: 'Dockerfile',
  go: 'Go',
  html: 'HTML',
  java: 'Java',
  javascript: 'JavaScript',
  js: 'JavaScript',
  json: 'JSON',
  markdown: 'Markdown',
  md: 'Markdown',
  nginx: 'Nginx',
  plaintext: '纯文本',
  powershell: 'PowerShell',
  ps1: 'PowerShell',
  python: 'Python',
  py: 'Python',
  shell: 'Shell',
  sql: 'SQL',
  ts: 'TypeScript',
  typescript: 'TypeScript',
  vue: 'Vue',
  xml: 'XML',
  yaml: 'YAML',
  yml: 'YAML'
}

const escapeHtml = (value = '') => String(value)
  .replace(/&/g, '&amp;')
  .replace(/</g, '&lt;')
  .replace(/>/g, '&gt;')
  .replace(/"/g, '&quot;')
  .replace(/'/g, '&#39;')

const stripHtml = (value = '') => String(value).replace(/<[^>]*>/g, '').trim()

const slugify = (value = '') => {
  const slug = stripHtml(value)
    .toLowerCase()
    .replace(/&[a-z0-9#]+;/gi, '')
    .replace(/[^\p{L}\p{N}\s-]/gu, '')
    .trim()
    .replace(/[\s_-]+/g, '-')
    .replace(/^-+|-+$/g, '')

  return slug || 'section'
}

const normalizeLanguage = (info = '') => String(info).trim().split(/\s+/)[0].toLowerCase()

const highlightLine = (line, language) => {
  if (!line) return '&#8203;'
  if (!language || !hljs.getLanguage(language)) return escapeHtml(line)

  return hljs.highlight(line, { language, ignoreIllegals: true }).value
}

export const renderArticleMarkdown = (markdown = '') => {
  const toc = []
  const headingCounts = new Map()
  let sectionNumber = 0
  let subsectionNumber = 0
  const renderer = new marked.Renderer()

  renderer.heading = (text, level, raw) => {
    const baseID = slugify(raw || text)
    const occurrence = (headingCounts.get(baseID) || 0) + 1
    headingCounts.set(baseID, occurrence)
    const id = occurrence === 1 ? baseID : `${baseID}-${occurrence}`

    if (level === 2 || level === 3) {
      if (level === 2) {
        sectionNumber += 1
        subsectionNumber = 0
      } else {
        subsectionNumber += 1
      }

      const number = level === 2
        ? `${sectionNumber}.`
        : `${sectionNumber || 1}.${subsectionNumber}.`
      toc.push({ id, text: stripHtml(text), level, number })
    }

    return `<h${level} id="${escapeHtml(id)}">${text}</h${level}>\n`
  }

  renderer.code = (code, infostring = '') => {
    const language = normalizeLanguage(infostring)
    const languageName = LANGUAGE_NAMES[language] || (language ? language.toUpperCase() : '纯文本')
    const languageClass = language ? ` language-${escapeHtml(language)}` : ''
    const lines = String(code).replace(/\n$/, '').split('\n')
    const renderedLines = lines.map((line, index) => (
      `<span class="code-line" data-line="${index + 1}">${highlightLine(line, language)}</span>`
    )).join('')

    return `<div class="code-block" data-language="${escapeHtml(language || 'plaintext')}">`
      + '<div class="code-toolbar">'
      + `<span class="code-language">${escapeHtml(languageName)}</span>`
      + '<button class="code-copy-button" type="button" data-copy-code aria-label="复制代码">复制</button>'
      + '</div>'
      + `<pre><code class="hljs${languageClass}">${renderedLines}</code></pre>`
      + '</div>\n'
  }

  renderer.image = (href, title, text) => {
    const alt = text || '文章图片'
    const titleAttribute = title ? ` title="${escapeHtml(title)}"` : ''

    return '<figure class="article-image-frame">'
      + `<img src="${escapeHtml(href || '')}" alt="${escapeHtml(alt)}" loading="lazy" decoding="async"${titleAttribute}>`
      + '<figcaption class="article-image-fallback" role="status">图片暂时无法显示</figcaption>'
      + '</figure>'
  }

  return {
    html: marked.parse(String(markdown), { renderer }),
    toc
  }
}
