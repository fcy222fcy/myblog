import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/login/Index.vue')
  },
  {
    path: '/',
    component: () => import('../components/layout/AdminLayout.vue'),
    redirect: '/dashboard',
    children: [
      { path: 'dashboard', name: 'Dashboard', component: () => import('../views/dashboard/Index.vue'), meta: { title: '仪表盘' } },
      { path: 'articles', name: 'Articles', component: () => import('../views/article/List.vue'), meta: { title: '文章管理' } },
      { path: 'articles/edit/:id?', name: 'ArticleEdit', component: () => import('../views/article/Edit.vue'), meta: { title: '文章编辑' } },
      { path: 'categories', name: 'Categories', component: () => import('../views/category/List.vue'), meta: { title: '分类管理' } },
      { path: 'tags', name: 'Tags', component: () => import('../views/tag/List.vue'), meta: { title: '标签管理' } },
      { path: 'comments', name: 'Comments', component: () => import('../views/comment/List.vue'), meta: { title: '评论管理' } },
      { path: 'daily-question', name: 'DailyQuestion', component: () => import('../views/daily/List.vue'), meta: { title: '每日一问' } },
      { path: 'about', name: 'About', component: () => import('../views/about/Index.vue'), meta: { title: '关于我' } },
      { path: 'audit', name: 'Audit', component: () => import('../views/audit/Index.vue'), meta: { title: '操作日志' } },
    ]
  }
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes
})

router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  if (to.path !== '/login' && !token) {
    next('/login')
  } else {
    next()
  }
})

export default router
