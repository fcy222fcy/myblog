import test from 'node:test'
import assert from 'node:assert/strict'
import { normalizeSiteProfile } from './siteProfile.js'

test('uses the site owner profile when the public profile is unavailable', () => {
  const profile = normalizeSiteProfile(null)

  assert.deepEqual(profile, {
    nickname: 'Fu ChengYan',
    bio: '日常落灰的个人博客，擅长面向搜索引擎编程。分享 Golang 开发、AI 和 NAS 折腾经验',
    avatar: '/avatar.jpg'
  })
})

test('keeps public profile values while filling missing fields from the site owner profile', () => {
  const profile = normalizeSiteProfile({
    nickname: 'FCY',
    bio: '',
    avatar: '/uploads/avatar.png'
  })

  assert.deepEqual(profile, {
    nickname: 'FCY',
    bio: '日常落灰的个人博客，擅长面向搜索引擎编程。分享 Golang 开发、AI 和 NAS 折腾经验',
    avatar: '/uploads/avatar.png'
  })
})
