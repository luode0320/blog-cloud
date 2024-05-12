<template>
  <el-tabs v-model="activeName" class="page-tool" stretch>
    <el-tab-pane label="JSON格式化" name="JsonFormat">
      <json-format class="tool-card" />
    </el-tab-pane>
    <el-tab-pane label="文本转换" name="Encryption">
      <encryption class="tool-card" />
    </el-tab-pane>
    <el-tab-pane label="RSA加密" name="Rsa">
      <rsa class="tool-card" />
    </el-tab-pane>
    <el-tab-pane label="随机数" name="Generate">
      <generate class="tool-card" />
    </el-tab-pane>
    <el-tab-pane label="正则表达式" name="Regular">
      <regular class="tool-card" />
    </el-tab-pane>
    <el-tab-pane label="Cron表达式" name="Cron">
      <cron class="tool-card" />
    </el-tab-pane>
    <el-tab-pane label="docker-compose转换" name="DockerConvert">
      <docker-convert class="tool-card" />
    </el-tab-pane>
    <el-tab-pane label="二维码" name="QrCode">
      <qr-code class="tool-card" />
    </el-tab-pane>
    <el-tab-pane label="图片编辑" name="ImageEditor">
      <image-editor />
    </el-tab-pane>
    <el-tab-pane label="图片压缩" name="ImageCompression">
      <image-compression class="tool-card" />
    </el-tab-pane>
  </el-tabs>
</template>

<script lang="ts" setup>
import { ref } from "vue";
import JsonFormat from "./components/json-format.vue";
import Encryption from "./components/encryption.vue";
import Rsa from "./components/rsa.vue";
import Generate from "./components/generate.vue";
import Regular from "./components/regular.vue";
import ImageCompression from "./components/image-compression.vue";
import QrCode from "./components/qr-code.vue";
import DockerConvert from "./components/docker-convert.vue";
import Cron from "./components/cron.vue";
import ImageEditor from "./components/image-editor.vue";

const activeName = ref("JsonFormat");
</script>

<style lang="scss">
.page-tool {
  overflow: hidden !important;
  display: flex;
  flex-direction: column;
  .el-tabs__item {
    user-select: none;
  }
  .el-tabs__content {
    overflow: auto;
    flex: 1;
    .el-tab-pane {
      height: 100%;
    }
  }
  .tool-card {
    height: 100%;
    width: 96%;
    margin: 0 2%;
  }
}
</style>
