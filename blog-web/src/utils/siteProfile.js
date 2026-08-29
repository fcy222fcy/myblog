export const DEFAULT_SITE_PROFILE = Object.freeze({
  nickname: 'Fu ChengYan',
  bio: '日常落灰的个人博客，擅长面向搜索引擎编程。分享 Golang 开发、AI 和 NAS 折腾经验',
  avatar: '/avatar.jpg'
})

const cleanText = (value) => typeof value === 'string' ? value.trim() : ''

export const normalizeSiteProfile = (profile) => ({
  nickname: cleanText(profile?.nickname) || DEFAULT_SITE_PROFILE.nickname,
  bio: cleanText(profile?.bio) || DEFAULT_SITE_PROFILE.bio,
  avatar: cleanText(profile?.avatar) || DEFAULT_SITE_PROFILE.avatar
})
