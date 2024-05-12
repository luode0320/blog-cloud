import { createRouter, createWebHashHistory } from "vue-router";
import Token from "@/store/token";

const routes = [
  { path: "/login", name: "login", component: () => import("@/views/login/index.vue") },
  {
    path: "/",
    name: "layout",
    component: () => import("@/views/layout/index.vue"),
    redirect: "/document",
    children: [
      { path: "/document", name: "document", component: () => import("@/views/document/index.vue") },
      { path: "/picture", name: "picture", component: () => import("@/views/picture/index.vue") },
      { path: "/tool", name: "tool", component: () => import("@/views/tool/index.vue") },
    ],
  },
  { path: "/open/document", name: "openDocument", component: () => import("@/views/open/doc.vue") },
  { path: "/open/publish", name: "openPublish", component: () => import("@/views/open/publish.vue") },
];

const router = createRouter({
  history: createWebHashHistory(),
  //@ts-ignore
  routes,
});

router.beforeEach((to, from) => {
  if (to.name === "openDocument" || to.name === "openPublish") {
    return;
  }
  if (!Token.getAccessToken() && to.name !== "login") {
    router.push({ name: "login" });
  }
});

export default router;
