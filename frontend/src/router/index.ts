import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/login',
      name: 'Login',
      component: () => import('../views/login/index.vue'),
      meta: { requiresAuth: false }
    },
    {
      path: '/',
      component: () => import('../layouts/MainLayout.vue'),
      redirect: '/dashboard',
      children: [
        {
          path: 'dashboard',
          name: 'Dashboard',
          component: () => import('../views/dashboard/index.vue'),
          meta: { title: '仪表盘', icon: 'Odometer' }
        },
        {
          path: 'customers',
          name: 'Customers',
          component: () => import('../views/customer/index.vue'),
          meta: { title: '客户管理', icon: 'User' }
        },
        {
          path: 'customers/:id',
          name: 'CustomerDetail',
          component: () => import('../views/customer/detail.vue'),
          meta: { title: '客户详情', hidden: true }
        },
        {
          path: 'contracts',
          name: 'Contracts',
          component: () => import('../views/contract/index.vue'),
          meta: { title: '合同管理', icon: 'Document' }
        },
        {
          path: 'contracts/:id',
          name: 'ContractDetail',
          component: () => import('../views/contract/detail.vue'),
          meta: { title: '合同详情', hidden: true }
        },
        {
          path: 'interactions',
          name: 'Interactions',
          component: () => import('../views/interaction/index.vue'),
          meta: { title: '互动日志', icon: 'ChatDotRound' }
        },
        {
          path: 'users',
          name: 'Users',
          component: () => import('../views/system/users.vue'),
          meta: { title: '用户管理', icon: 'UserFilled' }
        }
      ]
    }
  ]
})

router.beforeEach((to, _from, next) => {
  const token = localStorage.getItem('token')
  if (to.meta.requiresAuth !== false && !token) {
    next('/login')
  } else {
    next()
  }
})

export default router
