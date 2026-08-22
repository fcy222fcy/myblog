export const createPageViewRecorder = (recordFn) => {
  const recorded = new Set()

  return async (contentType, contentID) => {
    const key = `${contentType}:${contentID}`
    if (recorded.has(key)) return null
    recorded.add(key)
    try {
      return await recordFn(contentType, contentID)
    } catch (error) {
      recorded.delete(key)
      throw error
    }
  }
}
