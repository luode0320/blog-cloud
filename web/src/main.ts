import { createApp } from "vue";
import "./styles/index.scss"; // 导入应用程序的全局样式文件
import ElementPlus from "element-plus"; // 导入 Element Plus UI 库
import locale from "element-plus/es/locale/lang/zh-cn"; // 导入 Element Plus 的中文语言包
import App from "./App.vue"; // 导入根组件 App.vue
import router from "./router"; // 导入路由配置
import "virtual:svg-icons-register"; // 导入 SVG 图标

const app = createApp(App); // 创建 Vue 应用实例
app.use(router); // 注册路由
app.use(ElementPlus, { locale }); // 使用 Element Plus UI 库，并设置为中文语言
app.mount("#app"); // 将根组件挂载到 id 为 app 的 HTML 元素上
