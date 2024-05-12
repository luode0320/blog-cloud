<template>
  <div class="page-docker-convert">
    <el-button type="primary" @click="convert">转换</el-button>
    <div class="codemirror-view" style="height: 15%">
      <codemirror-editor v-model="dockerRun" placeholder="docker run 命令" />
    </div>
    <div class="codemirror-view" style="height: 15%">
      <codemirror-editor v-model="existCompose" :extensions="[javascript()]" placeholder="待合并的docker-compose" />
    </div>
    <div class="codemirror-view" style="height: calc(70% - 57px)">
      <codemirror-editor v-model="composeResult" :extensions="[javascript()]" placeholder="转换后的docker-compose" />
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref } from "vue";
import composerize from "composerize";
import CodemirrorEditor from "@/components/codemirror-editor";
import { javascript } from "@codemirror/lang-javascript";
import { ElMessage } from "element-plus";

const dockerRun = ref("");
const existCompose = ref("");
const composeResult = ref("");

const convert = () => {
  try {
    composeResult.value = composerize(dockerRun.value, existCompose.value);
  } catch (e) {
    console.error(e);
    ElMessage.warning("docker run 命令有误");
  }
};
</script>

<style lang="scss" scoped>
.page-docker-convert {
  display: flex;
  flex-direction: column;
  align-items: center;

  .codemirror-view {
    width: 100%;
    margin-top: 5px;
  }
}
</style>
