<template>
  <!-- 这个模板包含应用程序的根组件 -->
  <router-view></router-view> <!-- 路由视图，根据当前路由渲染不同的页面/组件 -->
</template>

<script lang="ts">
import { defineComponent } from "vue";
export default defineComponent({
  name: "App", // 定义组件的名称
});
</script>

<style>
/* 应用程序根元素的样式 */
#app {
  position: fixed; /* 相对于视口定位元素 */
  left: 0; /* 将元素与视口左边对齐 */
  right: 0; /* 将元素与视口右边对齐 */
  top: 0; /* 将元素与视口顶部对齐 */
  bottom: 0; /* 将元素与视口底部对齐 */
}
</style>
