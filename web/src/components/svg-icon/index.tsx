import { defineComponent, computed } from "vue";  // 导入Vue的相关功能以定义组件
import "./style.scss";  // 导入样式文件，假设是用于自定义SVG图标的样式

// 定义一个名为"SvgIcon"的Vue组件
export default defineComponent({
  // 组件名称，便于调试和识别
  name: "SvgIcon",

  // 接收的属性(props)，定义了图标前缀、图标名称、自定义类名和样式
  props: {
    // 图标类别的前缀，默认为'icon'
    prefix: {
      type: String,
      default: "icon",
    },
    // 必填的图标名称，用于定位具体的SVG图标
    name: {
      type: String,
      required: true,
    },
    // 可选的自定义类名，用于扩展样式
    className: {
      type: String,
    },
    // 自定义内联样式
    customStyle: {
      type: String,
    },
  },

  // 组件可以向外触发的事件，这里定义了一个点击事件
  emits: ["click"],

  // setup函数是Vue组件的核心选项，用于设置组件的响应式数据、计算属性、方法等
  setup(props, { emit }) {
    // 计算属性，组合图标ID，格式为'#prefix-name'
    const svgName = computed(() => {
      return `#${props.prefix}-${props.name}`;
    });

    // 计算属性，组合SVG元素的类名，包括默认的'svg-icon'和传入的自定义类名
    const svgClass = computed(() => {
      return "svg-icon " + props.className;
    });

    // 返回渲染函数，定义组件的模板结构
    return () => {
      return (
        <>
          {/* SVG元素结构 */}
          <svg
            class={svgClass.value}  // 应用计算出的类名
            aria-hidden={true}  // 隐藏于辅助技术（如屏幕阅读器）
            style={props.customStyle} // 应用自定义样式
            onClick={(ev: MouseEvent) => {
              emit("click", ev);  // 触发click事件，并传递点击事件对象
            }}
          >
            {/* 使用<use>元素引用外部SVG符号，href指向计算出的图标ID */}
            <use xlinkHref={svgName.value}></use>
          </svg>
        </>
      );
    };
  },
});
