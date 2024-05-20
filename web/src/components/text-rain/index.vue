<!-- Vue组件模板部分 -->
<template>
  <!-- 定义一个div容器用于展示文本雨效果，并通过ref绑定到JavaScript中以进行DOM操作 -->
  <div ref="textRainRef" class="text-rain"></div>
</template>

<script lang="ts" setup>
// 引入Vue的Composition API所需功能
import {ref, onMounted, onBeforeUnmount} from "vue";

// 定义组件接收的属性(props)
const props = defineProps({
  // 字符集合，用于随机选取构成文本雨的字符，默认包含数字和大写字母
  letters: {
    type: String,
    default: "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ",
  },
  // 文本出现的间隔时间，单位为毫秒，默认每100毫秒出现一个新字符
  speed: {
    type: Number,
    default: 100,
  },
});

// 使用ref创建对.text-rain元素的引用
const textRainRef = ref();
// 创建一个ref用来存储定时器的引用，便于清理
const interval = ref();


// onMounted生命周期钩子，在组件被挂载到DOM后执行
onMounted(() => {
  // 启动一个定时器，周期性地调用rain函数来生成文本雨效果
  interval.value = setInterval(() => rain(), props.speed);
});

// onBeforeUnmount生命周期钩子，在组件卸载之前执行
onBeforeUnmount(() => {
  // 清理定时器，避免内存泄漏
  if (interval.value) clearInterval(interval.value);
});

// 生成随机字符的函数
const randomText = () => {
  // 从props.letters中随机选取一个字符并返回
  return props.letters[Math.round(Math.random() * (props.letters.length - 1))];
};

// 实现文本雨逻辑的核心函数
const rain = () => {
  // 获取textRain元素
  const textRain = textRainRef.value;
  // 创建一个新的div元素表示单个“雨滴”
  const text = document.createElement("div");
  // 设置雨滴的内容为随机字符
  text.innerText = randomText();
  // 添加类名以应用样式, 雨滴
  text.className = "raindrop";
  // 随机设置雨滴的起始横向位置
  text.style.left = textRain.clientWidth * Math.random() + "px";
  // 随机设置雨滴的字体大小
  text.style.fontSize = 12 + 10 * Math.random() + "px";
  // 将雨滴添加到容器中
  textRain.appendChild(text);
  // 2秒后移除雨滴，模拟其消失
  setTimeout(() => textRain.removeChild(text), 2000);
};
</script>

<style lang="scss">
/* 定义文本雨容器的样式 */
.text-rain {
  width: 100%; /* 宽度100%，占据父元素全部宽度 */
  height: 100%; /* 高度100%，占据父元素全部高度 */
  position: relative; /* 相对定位，作为雨滴的定位上下文 */
  overflow: hidden; /* 隐藏超出容器的内容，使雨滴在底部消失 */

  .raindrop { /* 雨滴样式 */
    color: #ffffff; /* 文字颜色为白色 */
    height: 100%; /* 雨滴高度充满容器 */
    position: absolute; /* 绝对定位，可自由移动 */
    right: 0; /* 初始位置设在容器最右侧 */
    text-shadow: 0 0 5px #ffffff, 0 0 15px #ffffff, 0 0 30px #ffffff; /* 设置文字阴影增强视觉效果 */
    transform-origin: bottom; /* 变换的原点设在底部，用于动画效果 */
    animation: animate 2s linear forwards; /* 定义雨滴下落的动画 */
  }

  /* 动画关键帧定义 */
  @keyframes animate {
    /* 开始时水平位置不变 */
    0% {
      transform: translateX(0);
    }
    /* 下落至接近底部 */
    70% {
      transform: translateY(calc(100% - 22px));
    }
    /* 停留在底部 */
    100% {
      transform: translateY(calc(100% - 22px));
    }
  }
}
</style>
