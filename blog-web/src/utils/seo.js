/**
 * SEO 工具函数 - 动态 Meta 标签和结构化数据
 */

/**
 * 更新页面 Meta 标签
 * @param {Object} article - 文章数据
 * @param {Object} siteConfig - 站点配置
 */
export const updateMetaTags = (article, siteConfig = {}) => {
  const {
    siteName = '我的博客',
    siteUrl = window.location.origin
  } = siteConfig

  // 更新标题（优先使用后台配置的 SEO 标题）
  document.title = article?.title
    ? `${article?.seo_title || article.title} - ${siteName}`
    : `${siteName} - 记录生活与技术`

  const metaTags = [
    {
      name: 'description',
      content: article?.seo_description || article?.summary || `${siteName} - 分享技术、生活与思考`
    },
    {
      name: 'keywords',
      content: article?.seo_keywords || article?.tags?.map(t => t.name).join(',') || siteName
    },
    // Open Graph
    { property: 'og:title', content: article?.seo_title || article?.title || siteName },
    { property: 'og:description', content: article?.seo_description || article?.summary || siteName },
    { property: 'og:type', content: article ? 'article' : 'website' },
    { property: 'og:site_name', content: siteName },
    { property: 'og:url', content: window.location.href }
  ]

  if (article?.cover) {
    metaTags.push({ property: 'og:image', content: article.cover })
  }

  metaTags.forEach(({ name, property, content }) => {
    if (!content) return
    const selector = name
      ? `meta[name="${name}"]`
      : `meta[property="${property}"]`
    let meta = document.querySelector(selector)
    if (!meta) {
      meta = document.createElement('meta')
      if (name) meta.name = name
      if (property) meta.setAttribute('property', property)
      document.head.appendChild(meta)
    }
    meta.content = content
  })
}

/**
 * 添加 JSON-LD 结构化数据
 * @param {Object} article - 文章数据
 * @param {Object} siteConfig - 站点配置
 */
export const addStructuredData = (article, siteConfig = {}) => {
  const {
    siteName = '我的博客',
    siteUrl = window.location.origin
  } = siteConfig

  let schema

  if (article) {
    schema = {
      '@context': 'https://schema.org',
      '@type': 'BlogPosting',
      'headline': article.title,
      'description': article.summary,
      'image': article.cover || `${siteUrl}/default-cover.jpg`,
      'datePublished': article.created_at,
      'dateModified': article.updated_at || article.created_at,
      'author': {
        '@type': 'Person',
        'name': siteName
      },
      'publisher': {
        '@type': 'Organization',
        'name': siteName,
        'logo': {
          '@type': 'ImageObject',
          'url': `${siteUrl}/logo.png`
        }
      }
    }
  } else {
    schema = {
      '@context': 'https://schema.org',
      '@type': 'WebSite',
      'name': siteName,
      'url': siteUrl,
      'description': `${siteName} - 分享技术、生活与思考`
    }
  }

  let script = document.querySelector('script[type="application/ld+json"]')
  if (!script) {
    script = document.createElement('script')
    script.type = 'application/ld+json'
    document.head.appendChild(script)
  }
  script.textContent = JSON.stringify(schema)
}

/**
 * 重置 SEO 标签（离开页面时）
 */
export const resetMetaTags = (siteConfig = {}) => {
  const { siteName = '我的博客' } = siteConfig
  updateMetaTags(null, siteConfig)
  const script = document.querySelector('script[type="application/ld+json"]')
  if (script) script.remove()
}
