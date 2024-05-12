<template>
  <el-form class="page-regular" label-width="100px">
    <el-form-item label="常用正则">
      <el-button style="margin-bottom: 3px" v-for="(item, index) in regList" :key="index" @click="regClick(item.value)">{{ item.label }}</el-button>
    </el-form-item>
    <el-form-item label="正则表达式">
      <el-input v-model="regText" type="textarea" :rows="3" resize="none" placeholder="请输入正则表达式"></el-input>
      <el-checkbox-group v-model="modifierList">
        <el-checkbox label="g">全局匹配(g)</el-checkbox>
        <el-checkbox label="m">多行匹配(m)</el-checkbox>
        <el-checkbox label="i">不区分大小写(i)</el-checkbox>
        <el-checkbox label="group">分组</el-checkbox>
      </el-checkbox-group>
    </el-form-item>
    <el-form-item label="匹配内容">
      <div class="codemirror-view">
        <codemirror-editor v-model="inText" placeholder="请输入待匹配内容" />
      </div>
    </el-form-item>
    <el-form-item label="替换内容">
      <el-input v-model="replaceText" type="textarea" :rows="3" resize="none" placeholder="请输入替换内容"></el-input>
      <div class="button-group">
        <el-button class="my-button" @click="testClick" type="primary">匹配</el-button>
        <el-button class="my-button" @click="execClick" type="success">检索</el-button>
        <el-button class="my-button" @click="replaceClick" type="warning">替换</el-button>
      </div>
    </el-form-item>
    <el-form-item label="结果">
      <div class="codemirror-view">
        <codemirror-editor v-model="outText" />
      </div>
    </el-form-item>
  </el-form>
</template>

<script lang="ts" setup>
import { ref } from "vue";
import CodemirrorEditor from "@/components/codemirror-editor";

const regList = ref([
  { label: "空白", value: "\\s+" },
  { label: "数字", value: "[0-9]+" },
  { label: "英文", value: "[A-Za-z]+" },
  { label: "汉字", value: "[\\u4e00-\\u9fa5]+" },
  { label: "手机号码", value: "1\\d{10}" },
  { label: "邮箱", value: "\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*" },
  { label: "日期时间", value: "[1-9]\\d{3}-(0[1-9]|1[0-2])-(0[1-9]|[1-2][0-9]|3[0-1])\\s+(20|21|22|23|[0-1]\\d):[0-5]\\d:[0-5]\\d" },
  { label: "18位身份证", value: "[1-9]\\d{5}(18|19|([23]\\d))\\d{2}((0[1-9])|(10|11|12))(([0-2][1-9])|10|20|30|31)\\d{3}[0-9Xx]" },
]);
const modifierList = ref(["g"]);
const regText = ref("");
const inText = ref("");
const outText = ref("");
const replaceText = ref("");

/**
 * 初始化正则表达式
 */
const initPatten = () => {
  outText.value = "";
  let modifiers = "";
  for (let i in modifierList.value) {
    // 排除group
    if (modifierList.value[i] !== "group") {
      modifiers += modifierList.value[i];
    }
  }
  return new RegExp(regText.value, modifiers);
};

/**
 * 匹配是否存在
 */
const testClick = () => {
  const pattern = initPatten();
  outText.value = String(pattern.test(inText.value));
};

/**
 * 检索
 */
const execClick = () => {
  const pattern = initPatten();
  // 判断是否分组
  let groupStatus = false;
  for (let i in modifierList.value) {
    if (modifierList.value[i] === "group") {
      groupStatus = true;
      break;
    }
  }
  let result = pattern.exec(inText.value);
  if (groupStatus) {
    // 分组
    // 全局匹配
    if (pattern.global) {
      let first = true;
      while (result !== null && result[0] !== "") {
        if (first) {
          first = false;
        } else {
          outText.value += "\r\n";
        }
        for (let i in result) {
          if (parseInt(i) < result.length) {
            outText.value += "[" + i + "]" + result[i];
          }
        }
        result = pattern.exec(inText.value);
      }
    } else {
      // 非全局匹配
      if (result && result[0] !== "") {
        for (let i in result) {
          if (parseInt(i) <= result.length) {
            outText.value += "[" + i + "]" + result[i];
          }
        }
      }
    }
  } else {
    // 不分组
    // 全局匹配
    if (pattern.global) {
      let first = true;
      while (result !== null && result[0] !== "") {
        if (first) {
          outText.value += result[0];
          first = false;
        } else {
          outText.value += "\r\n" + result[0];
        }
        result = pattern.exec(inText.value);
      }
    } else {
      // 非全局匹配
      if (result && result[0] !== "") {
        outText.value = result[0];
      }
    }
  }
};

/**
 * 替换匹配值
 */
const replaceClick = () => {
  const pattern = initPatten();
  outText.value = inText.value.replace(pattern, replaceText.value);
};

/**
 * 常用正则按钮点击
 * @param val
 */
const regClick = (val: string) => {
  regText.value = val;
};
</script>

<style lang="scss">
.page-regular {
  width: 100%;
  .button-group {
    margin-top: 5px;
  }
  .codemirror-view {
    width: 100%;
    height: 200px;
  }
}
</style>
