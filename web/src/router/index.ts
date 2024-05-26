import { createRouter, createWebHashHistory } from "vue-router"; // 导入 Vue Router 相关模块
import Token from "@/store/token"; // 导入 Token 模块，用于验证用户身份

// 定义路由配置
const routes = [
  { path: "/login", name: "login", component: () => import("@/views/login/index.vue") }, // 登录页面路由
  {
    path: "/", // 根路由
    name: "layout",
    component: () => import("@/views/layout/index.vue"), // 布局组件路由
    redirect: "/document", // 重定向到文档页面
    children: [
      { path: "/document", name: "document", component: () => import("@/views/document/index.vue") }, // 文档页面路由
      // { path: "/picture", name: "picture", component: () => import("@/views/picture/index.vue") }, // 图片页面路由
      { path: "/tool", name: "tool", component: () => import("@/views/tool/index.vue") }, // 工具页面路由
    ],
  },
  // { path: "/open/document", name: "openDocument", component: () => import("@/views/open/doc.vue") }, // 公开文档页面路由
  // { path: "/open/publish", name: "openPublish", component: () => import("@/views/open/publish.vue") }, // 公开发布页面路由
];

// 创建路由实例
const router = createRouter({
  history: createWebHashHistory(),
  //@ts-ignore
  routes,
});

// 路由导航守卫，用于身份验证和路由拦截
router.beforeEach((to, from) => {
  // 如果是访问公开文档或公开发布页面，则放行
  if (to.name === "openDocument" || to.name === "openPublish") {
    return;
  }
  // 如果用户未登录且不是访问登录页面，则跳转至登录页面
  if (!Token.getAccessToken() && to.name !== "login") {
    router.push({ name: "login" });
  }
});

export default router; // 导出路由实例
