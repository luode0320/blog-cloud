<template>
  <div class="page-cron">
    <el-button type="primary" @click="cronConvert">解析</el-button>
    <el-input class="input-view" v-model="cronText" placeholder="请输入Cron表达式" clearable @clear="clear"></el-input>
    <el-input class="input-view" v-model="translateText" type="textarea" :rows="3" resize="none" placeholder="Cron解析结果" readonly></el-input>
    <el-input class="input-view" v-model="nextText" type="textarea" :rows="5" resize="none" placeholder="近5次执行时间" readonly></el-input>
  </div>
</template>

<script lang="ts" setup>
import { ref } from "vue";
import cronstrue from "cronstrue";
import "cronstrue/locales/zh_CN";
import cronParser from "cron-parser";
import { formatTime, getWeek } from "@/utils";

const cronText = ref("");
const translateText = ref("");
const nextText = ref("");

const cronConvert = () => {
  try {
    translateText.value = cronstrue.toString(cronText.value, { use24HourTimeFormat: true, locale: "zh_CN" });
  } catch (e) {
    console.error(e);
    translateText.value = "解析失败";
    nextText.value = "解析失败";
    return;
  }
  try {
    const interval = cronParser.parseExpression(cronText.value);
    const nextArr = [];
    for (let i = 0; i < 5; i++) {
      let nextDate = interval.next().toDate();
      nextArr.push(formatTime(nextDate) + " " + getWeek(nextDate));
    }
    nextText.value = nextArr.join("\n");
  } catch (e) {
    console.error(e);
    nextText.value = "解析失败";
  }
};

const clear = () => {
  translateText.value = "";
  nextText.value = "";
};
</script>

<style lang="scss" scoped>
.page-cron {
  display: flex;
  flex-direction: column;
  align-items: center;
  .input-view {
    margin-top: 10px;
    width: 400px;
    max-width: 100%;
  }
}
</style>
