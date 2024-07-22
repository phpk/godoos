import { createWebHistory, createRouter } from "vue-router";

const routes = [
  // {
  //   path: '/',
  //   component: () => import('../pages/test.vue'),  // 首页组件
  // }, 
  // {
  //   path: '/install',
  //   name: "install",
  //   component: () => import('../views/install/install.vue')  // 关于我们组件
  // },
  {
    path: '/:pathMatch(.*)*',
    name: "not-found",
    component: () => import('../components/window/NotFound.vue')
  }
]

const router = createRouter({
  history: createWebHistory(), // 路由类型
  routes // short for `routes: routes`
})


export default router