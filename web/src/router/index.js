import { createRouter, createWebHashHistory } from 'vue-router'

/* Layout */
import Layout from '@/layout/index.vue'

/* Router Modules */

/**
 * Note: sub-menu only appear when route children.length >= 1
 * Detail see: https://panjiachen.github.io/vue-element-admin-site/guide/essentials/router-and-nav.html
 *
 * hidden: true                   if set true, item will not show in the sidebar(default is false)
                                  * 如果设置为true, item将不会显示在侧边栏中(默认为false)
 * alwaysShow: true               if set true, will always show the root menu
 *                                if not set alwaysShow, when item has more than one children route,
 *                                it will becomes nested mode, otherwise not show the root menu
                                  * 如果设置为true，将始终显示根菜单

                                  *如果没有设置alwaysShow，当item有多个子路由时，

                                  *它将成为嵌套模式，否则不显示根菜单
 * redirect: noRedirect           if set noRedirect will no redirect in the breadcrumb 如果设置noRedirect将不会在面包屑中重定向
 * name:'router-name'             the name is used by <keep-alive> (must set!!!)
 * meta : {
    roles: ['admin','editor']    control the page roles (you can set multiple roles)控制页面角色(您可以设置多个角色)
    title: 'title'               the name show in sidebar and breadcrumb (recommend set)名称显示在侧边栏和面包屑(建议设置)
    icon: 'svg-name'/'el-icon-x' the icon show in the sidebar图标显示在侧边栏中
    noCache: true                if set true, the page will no be cached(default is false)如果设置为true，页面将不会被缓存(默认为false)
    affix: true                  if set true, the tag will affix in the tags-view如果设置为true，标签将会被粘贴到标签视图中
    breadcrumb: false            if set false, the item will hidden in breadcrumb(default is true)如果设置为false，该项目将隐藏在breadcrumb中(默认为true)
    activeMenu: '/example/list'  if set path, the sidebar will highlight the path you set如果设置了path，侧边栏将突出显示您设置的路径
  }
 */

/**
 * constantRoutes
 * a base page that does not have permission requirements
 * all roles can be accessed
 * 没有权限要求的基本页
  *所有角色都可以访问
 */
export const constantRoutes = [
  {
    path: '/redirect',
    component: Layout,
    hidden: true,
    children: [
      {
        path: '/redirect/:path(.*)',
        component: () => import('@/views/redirect/index.vue')
      },
      
    ]
  },
  {
    path: '/login',
    component: () => import('@/views/login/index.vue'),
    hidden: true
  },
  {
      path: '/changePass',
      name:'changePass',
      hidden: true,
      meta: { title: '忘记密码', icon: 'user', noCache: true },
      component: () => import('@/views/changePassword/index.vue'),
  },
  {
    path: '/404',
    component: () => import('@/views/error-page/404'),
    hidden: true
  },
  {
    path: '/401',
    component: () => import('@/views/error-page/401'),
    hidden: true
  },
  {

    path: '/',
    component: Layout,
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        component: () => import('@/views/dashboard/index.vue'),
        name: 'Dashboard',
        meta: { title: '首页', icon: 'dashboard', affix: true }
      }
      
    ]
  },
  {
    path: '/profile',
    component: Layout,
    redirect: '/profile/index',
    hidden: true,
    children: [
      {
        path: 'index',
        component: () => import('@/views/profile/index.vue'),
        name: 'Profile',
        meta: { title: '个人中心', icon: 'user', noCache: true }
      },
    
    ]
  },
  {
    path: '/groupUser',
    component: Layout,
    redirect: '/groupUser/userList/index',
    hidden: true,
    children: [
      {
        path: '/userList',
        component: () => import('@/views/groupUser/userList/index.vue'),
        name: 'userList',
        meta: { title: '分组成员', icon: 'user', noCache: true }
      }
    ]

  },
  {
    // 设置分区：侧边栏「设置」常驻入口指向的独立配置页面。
    // hidden=true，避免与侧边栏中硬编码的「设置」分组重复出现。
    path: '/settings',
    component: Layout,
    redirect: '/settings/directory',
    hidden: true,
    meta: { title: '设置', icon: 'setting', requireAdmin: true },
    children: [
      {
        path: 'directory',
        component: () => import('@/views/settings/directory/index.vue'),
        name: 'SettingsDirectory',
        meta: { title: '目录配置', icon: 'setting', noCache: true }
      },
      {
        path: 'thirdparty',
        component: () => import('@/views/settings/thirdParty/index.vue'),
        name: 'SettingsThirdParty',
        meta: { title: '平台对接', icon: 'connection', noCache: true }
      },
      {
        path: 'notification',
        component: () => import('@/views/settings/notification/index.vue'),
        name: 'SettingsNotification',
        meta: { title: '通知设置', icon: 'bell', noCache: true }
      },
      {
        path: 'sync',
        component: () => import('@/views/settings/sync/index.vue'),
        name: 'SettingsSync',
        meta: { title: '定时同步', icon: 'refresh', noCache: true }
      },
      {
        path: 'security',
        component: () => import('@/views/settings/security/index.vue'),
        name: 'SettingsSecurity',
        meta: { title: '登录安全', icon: 'lock', noCache: true }
      }
    ]
  },

]

/**
 * asyncRoutes
 * the routes that need to be dynamically loaded based on user roles
 */

const router = createRouter({
  history: createWebHashHistory(),
  scrollBehavior: () => ({ top: 0 }),
  routes: constantRoutes
})

// 收集 constantRoutes 的名字，resetRouter 时只移除动态新增的路由
const constantRouteNames = new Set()
constantRoutes.forEach(r => {
  if (r.name) constantRouteNames.add(r.name)
  ;(r.children || []).forEach(c => c.name && constantRouteNames.add(c.name))
})

export function resetRouter() {
  router.getRoutes().forEach(route => {
    const { name } = route
    if (name && !constantRouteNames.has(name)) {
      router.removeRoute(name)
    }
  })
}

export default router
