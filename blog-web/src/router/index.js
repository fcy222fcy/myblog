import { createRouter, createWebHashHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: () => import('../views/Home.vue')
  },
  {
    path: '/archives',
    name: 'Archives',
    component: () => import('../views/Archives.vue')
  },

  {
    path: '/about',
    name: 'About',
    component: () => import('../views/About.vue')
  },
  {
    path: '/search',
    name: 'Search',
    component: () => import('../views/Search.vue')
  },
  {
    path: '/post/:slug',
    name: 'Article',
    component: () => import('../views/Article.vue')
  },
  {
    path: '/daily/:date?',
    name: 'DailyQuestion',
    component: () => import('../views/DailyQuestion.vue')
  },
  {
    path: '/tag/:id',
    name: 'Tag',
    component: () => import('../views/Tag.vue')
  },
  {
    path: '/category/:id',
    name: 'Category',
    component: () => import('../views/Category.vue')
  }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes,
  scrollBehavior(to, from, savedPosition) {
    if (savedPosition) {
      return savedPosition
    } else {
      return { top: 0 }
    }
  }
})

// 记录上一路由信息，供文章详情页「返回列表」按钮精确返回来源页
router.afterEach((to, from) => {
  router.__prevRouteName = from.name
  router.__prevRoutePath = from.fullPath
})

export default router
