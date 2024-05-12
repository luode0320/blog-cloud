<template>
  <div class="text-rain" ref="textRainRef"></div>
</template>

<script lang="ts" setup>
import { ref, onMounted, onBeforeUnmount } from "vue";

const props = defineProps({
  letters: {
    type: String,
    default: "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ",
  },
  // 多少毫秒出现一个字
  speed: {
    type: Number,
    default: 100,
  },
});

const textRainRef = ref();
const interval = ref();

onMounted(() => {
  interval.value = setInterval(() => {
    rain();
  }, props.speed);
});

onBeforeUnmount(() => {
  if (interval.value) {
    clearInterval(interval.value);
  }
});

const randomText = () => {
  return props.letters[Math.round(Math.random() * (props.letters.length - 1))];
};

const rain = () => {
  const textRain = textRainRef.value;
  const text = document.createElement("div");
  text.innerText = randomText();
  text.className = "raindrop";
  text.style.left = textRain.clientWidth * Math.random() + "px";
  text.style.fontSize = 12 + 10 * Math.random() + "px";
  textRain.appendChild(text);
  setTimeout(() => {
    textRain.removeChild(text);
  }, 2000);
};
</script>

<style lang="scss">
.text-rain {
  width: 100%;
  height: 100%;
  position: relative;
  overflow: hidden;

  .raindrop {
    color: #ffffff;
    height: 100%;
    position: absolute;
    right: 0;
    text-shadow: 0 0 5px #ffffff, 0 0 15px #ffffff, 0 0 30px #ffffff;
    transform-origin: bottom;
    animation: animate 2s linear forwards;
  }

  @keyframes animate {
    0% {
      transform: translateX(0);
    }
    70% {
      transform: translateY(calc(100% - 22px));
    }
    100% {
      transform: translateY(calc(100% - 22px));
    }
  }
}
</style>
