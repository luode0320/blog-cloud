<template>
  <div class="page-qrcode">
    <el-form label-width="100px">
      <el-form-item label="文字内容">
        <el-input
          v-model="content"
          placeholder="请输入二维码文字内容"
          type="textarea"
          :rows="3"
          resize="none"
          :maxlength="400"
          show-word-limit
        ></el-input>
      </el-form-item>
      <el-form-item label="大小">
        <el-slider v-model="size" :min="10" :max="2000" show-input></el-slider>
      </el-form-item>
      <el-form-item label="边距">
        <el-slider v-model="margin" :min="0" :max="100" show-input></el-slider>
      </el-form-item>
      <el-form-item label="背景颜色">
        <el-color-picker v-model="background"></el-color-picker>
      </el-form-item>
      <el-form-item label="二维码颜色">
        <el-color-picker v-model="foreground"></el-color-picker>
      </el-form-item>
      <el-form-item label="容错等级">
        <el-select v-model="level">
          <el-option label="低" value="L"></el-option>
          <el-option label="中" value="M"></el-option>
          <el-option label="高" value="Q"></el-option>
          <el-option label="极高" value="H"></el-option>
        </el-select>
      </el-form-item>
    </el-form>
    <qrcode-vue :value="content" :size="size" :level="level" :background="background" :foreground="foreground" :margin="margin"></qrcode-vue>
  </div>
</template>

<script lang="ts" setup>
import { ref, Ref } from "vue";
import QrcodeVue, { Level } from "qrcode.vue";

const content = ref("");
const size = ref(360);
const level: Ref<Level> = ref("H");
const background = ref("#ffffff");
const foreground = ref("#000000");
const margin = ref(2);
</script>

<style lang="scss" scoped>
.page-qrcode {
  display: flex;
  flex-direction: column;
  align-items: center;
  .el-form {
    width: 400px;
    max-width: 100%;
  }
}
</style>
