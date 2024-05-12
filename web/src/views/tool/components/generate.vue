<template>
  <div class="page-generate">
    <el-form label-width="100px">
      <el-form-item label="随机位数">
        <el-input-number v-model="inputText" placeholder="请输入" :min="1" :max="9999"></el-input-number>
      </el-form-item>
    </el-form>
    <div class="button-group">
      <el-button type="success" @click="generate('randomNum')">数字</el-button>
      <el-button type="success" @click="generate('randomStr')">字符串</el-button>
      <el-button type="primary" @click="generate('UUID')">UUID</el-button>
      <el-button type="primary" @click="generate('date')">日期</el-button>
      <el-button type="primary" @click="generate('timestamp')">时间戳</el-button>
      <el-button type="primary" @click="generate('color')">颜色</el-button>
    </div>
    <el-input v-model="outputText" type="textarea" :rows="10" resize="none" placeholder="结果文本" readonly></el-input>
  </div>
</template>

<script lang="ts" setup>
import { ref } from "vue";
import { ElMessage } from "element-plus";
import { formatTime, checkNaN, randomStr, random } from "@/utils";
import { uuid } from "@/utils";

const inputText = ref(6);
const outputText = ref("");
const colorText = ref("0123456789abcdef");

const generate = (type: string) => {
  try {
    let date = new Date();
    switch (type) {
      case "UUID":
        let uuidArr = [];
        for (let i = 0; i < 10; i++) {
          uuidArr.push(uuid());
        }
        outputText.value = uuidArr.join("\n");
        break;
      case "date":
        outputText.value = formatTime(date, "YYYY/MM/DD HH:mm:ss.SSS");
        break;
      case "timestamp":
        outputText.value = date.getTime().toString();
        break;
      case "randomNum":
        if (!inputText.value || checkNaN(inputText.value)) {
          outputText.value = "";
        } else {
          outputText.value = randomStr(inputText.value, true);
        }
        break;
      case "randomStr":
        if (!inputText.value || checkNaN(inputText.value)) {
          outputText.value = "";
        } else {
          outputText.value = randomStr(inputText.value);
        }
        break;
      case "color":
        let colorArr = [];
        for (let i = 0; i < 10; i++) {
          colorArr.push("#" + random(colorText.value, 6));
        }
        outputText.value = colorArr.join("\n");
        break;
    }
  } catch (e) {
    console.error(e);
    ElMessage.error("文本处理失败");
    outputText.value = "";
  }
};
</script>

<style lang="scss">
.page-generate {
  display: flex;
  flex-direction: column;
  align-items: center;
  .el-form {
    .el-form-item {
      margin-bottom: 0;
    }
  }
  .button-group {
    width: 100%;
    margin-bottom: 10px;
    display: flex;
    flex-wrap: wrap;
    justify-content: center;
    .el-button {
      margin: 10px 10px 0 10px;
    }
    .el-button + .el-button {
      margin-left: 0;
    }
  }
}
</style>
