// 声明.vue文件模块的类型
declare module "*.vue" {
  import { defineComponent } from "vue"; // 导入Vue的defineComponent函数
  const Component: ReturnType<typeof defineComponent>; // 声明一个Component常量，其类型为defineComponent函数返回值的类型
  export default Component; // 导出Component常量
}
