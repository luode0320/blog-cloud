<template>
  <div class="page-encryption">
    <div class="button-group">
      <el-button type="warning" @click="encrypt('count')">字符数</el-button>
      <el-button type="warning" @click="encrypt('countNotBlank')">字符数(无空白)</el-button>
      <el-button type="warning" @click="encrypt('removeBlank')">去除空白</el-button>
      <el-button type="warning" @click="encrypt('removeWrap')">去除换行</el-button>
      <el-button type="success" @click="encrypt('timestampToDate')">时间戳转日期</el-button>
      <el-button type="success" @click="encrypt('dateToTimestamp')">日期转时间戳</el-button>
      <el-button type="success" @click="encrypt('lowerToUpper')">小写转大写</el-button>
      <el-button type="success" @click="encrypt('upperToLower')">大写转小写</el-button>
      <el-button type="danger" @click="encrypt('URLEncode')">URL编码</el-button>
      <el-button type="danger" @click="encrypt('URLDecode')">URL解码</el-button>
      <el-button type="danger" @click="encrypt('Base64Encode')">Base64编码</el-button>
      <el-button type="danger" @click="encrypt('Base64Decode')">Base64解码</el-button>
      <el-button type="danger" @click="encrypt('Base64ToHex')">Base64解码(16进制)</el-button>
      <el-button type="primary" @click="encrypt('MD5')">MD5</el-button>
      <el-button type="primary" @click="encrypt('SHA1')">SHA1</el-button>
      <el-button type="primary" @click="encrypt('SHA256')">SHA256</el-button>
      <el-button type="primary" @click="encrypt('SHA512')">SHA512</el-button>
      <el-button type="primary" @click="encrypt('SHA3-256')">SHA3-256</el-button>
      <el-button type="primary" @click="encrypt('SHA3-512')">SHA3-512</el-button>
    </div>
    <el-input v-model="inputText" type="textarea" :rows="10" resize="none" placeholder="请输入文本"></el-input>
    <el-input v-model="outputText" type="textarea" :rows="10" resize="none" placeholder="结果文本" readonly></el-input>
  </div>
</template>

<script lang="ts" setup>
import { ref } from "vue";
import { ElMessage } from "element-plus";
import crypto from "crypto-js";
import { encode, decode, toUint8Array } from "js-base64";
import { formatTime, parseTime } from "@/utils";

const inputText = ref("");
const outputText = ref("");

const encrypt = (type: string) => {
  if (!inputText.value) {
    outputText.value = "";
    return;
  }
  try {
    switch (type) {
      case "count":
        outputText.value = inputText.value.length.toString();
        break;
      case "countNotBlank":
        outputText.value = inputText.value.replace(/\s+/g, "").length.toString();
        break;
      case "removeBlank":
        outputText.value = inputText.value.replace(/\s+/g, "");
        break;
      case "removeWrap":
        outputText.value = inputText.value.replace(/\n+/g, "");
        break;
      case "timestampToDate":
        outputText.value = formatTime(parseInt(inputText.value), "YYYY/MM/DD HH:mm:ss.SSS");
        break;
      case "dateToTimestamp":
        outputText.value = parseTime(inputText.value).getTime().toString();
        break;
      case "lowerToUpper":
        outputText.value = inputText.value.toUpperCase();
        break;
      case "upperToLower":
        outputText.value = inputText.value.toLowerCase();
        break;
      case "URLEncode":
        outputText.value = encodeURIComponent(inputText.value);
        break;
      case "URLDecode":
        outputText.value = decodeURIComponent(inputText.value);
        break;
      case "Base64Encode":
        outputText.value = encode(inputText.value);
        break;
      case "Base64Decode":
        outputText.value = decode(inputText.value);
        break;
      case "Base64ToHex":
        outputText.value = Array.prototype.map.call(toUint8Array(inputText.value), (t) => ("00" + t.toString(16)).slice(-2)).join("");
        break;
      case "MD5":
        outputText.value = crypto.MD5(inputText.value).toString();
        break;
      case "SHA1":
        outputText.value = crypto.SHA1(inputText.value).toString();
        break;
      case "SHA256":
        outputText.value = crypto.SHA256(inputText.value).toString();
        break;
      case "SHA512":
        outputText.value = crypto.SHA512(inputText.value).toString();
        break;
      case "SHA3-256":
        outputText.value = crypto.SHA3(inputText.value, { outputLength: 256 }).toString();
        break;
      case "SHA3-512":
        outputText.value = crypto.SHA3(inputText.value, { outputLength: 512 }).toString();
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
.page-encryption {
  display: flex;
  flex-direction: column;
  align-items: center;
  .el-textarea + .el-textarea {
    margin-top: 10px;
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
