import { createApp } from "vue";
import "./styles/index.scss";
import ElementPlus from "element-plus";
import locale from "element-plus/es/locale/lang/zh-cn";
import App from "./App.vue";
import router from "./router";
import "virtual:svg-icons-register";

const app = createApp(App);
app.use(router);
app.use(ElementPlus, { locale });
app.mount("#app");
