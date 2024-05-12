<template>
  <div class="page-json-format">
    <el-form label-width="60px">
      <el-form-item label="空格数">
        <el-select v-model="spaceCount" style="width: 60px; margin-right: 10px">
          <el-option label="0" :value="0"></el-option>
          <el-option label="2" :value="2"></el-option>
          <el-option label="4" :value="4"></el-option>
        </el-select>
        <el-button type="primary" @click="formatClick">格式化</el-button>
      </el-form-item>
    </el-form>
    <div class="json-view">
      <codemirror-editor v-model="content" placeholder="请输入json" />
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref } from "vue";
import { jsonrepair } from "jsonrepair";
import CodemirrorEditor from "@/components/codemirror-editor";
import { ElMessage } from "element-plus";

const content = ref("");
const spaceCount = ref(4);

const formatClick = () => {
  try {
    content.value = JSON.stringify(JSON.parse(jsonrepair(content.value)), null, spaceCount.value);
  } catch (e) {
    ElMessage.warning("json格式有误");
    console.error(e);
  }
};
</script>

<style lang="scss">
.page-json-format {
  display: flex;
  flex-direction: column;
  align-items: center;
  .json-view {
    margin-top: 5px;
    height: calc(100% - 60px);
    width: calc(100% - 2px);
    overflow: hidden;
  }
}
</style>
