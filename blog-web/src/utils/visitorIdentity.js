const STORAGE_KEY = 'blog_anonymous_visitor_id'
const UUID_PATTERN = /^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/i

const createUUID = (cryptoProvider) => {
  if (typeof cryptoProvider?.randomUUID === 'function') {
    return cryptoProvider.randomUUID()
  }
  if (typeof cryptoProvider?.getRandomValues !== 'function') return ''
  const bytes = new Uint8Array(16)
  cryptoProvider.getRandomValues(bytes)
  bytes[6] = (bytes[6] & 0x0f) | 0x40
  bytes[8] = (bytes[8] & 0x3f) | 0x80
  const hex = [...bytes].map(value => value.toString(16).padStart(2, '0')).join('')
  return `${hex.slice(0, 8)}-${hex.slice(8, 12)}-${hex.slice(12, 16)}-${hex.slice(16, 20)}-${hex.slice(20)}`
}

export const getOrCreateVisitorID = (storage, cryptoProvider) => {
  try {
    const existing = storage?.getItem(STORAGE_KEY) || ''
    if (UUID_PATTERN.test(existing)) return existing.toLowerCase()

    const created = createUUID(cryptoProvider)
    if (!UUID_PATTERN.test(created)) return ''
    storage?.setItem(STORAGE_KEY, created)
    return created.toLowerCase()
  } catch {
    return ''
  }
}
