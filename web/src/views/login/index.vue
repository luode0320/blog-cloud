<template>
  <!-- 页面主体容器，包含动态文字雨背景、标题、登录/注册表单 -->
  <div class="page-login">
    <!-- 引入自定义的动态文字雨组件作为背景装饰 -->
    <text-rain class="text-rain-background"></text-rain>

    <!-- 标题区域，可点击跳转至公开博客页面 -->
    <div class="title-view" title="查看公开博客" @click="publishClick">博客</div>

    <!-- 利用Vue的transition组件实现登录/注册表单的淡入淡出效果 -->
    <transition name="fade">
      <!-- 登录状态下的表单 -->
      <div v-if="isLogin" class="content-view">
        <form>
          <!-- 用户名输入框，使用Element Plus的输入组件 -->
          <el-input
              v-model.trim="inputData.name"
              class="input-view"
              clearable
              placeholder="请输入用户名"
              size="large"
              @keyup.enter.native="loginClick"
          ></el-input>
          <!-- 密码输入框 -->
          <el-input
              v-model.trim="inputData.password"
              class="input-view"
              clearable
              placeholder="请输入密码"
              size="large"
              type="password"
              @keyup.enter.native="loginClick"
          ></el-input>
        </form>

        <!-- 注册按钮，点击后切换到注册界面 -->
        <el-button :disabled="loading" class="register-button" link size="small" text type="primary"
                   @click="registerClick">注册
        </el-button>

        <!-- 登录按钮，提交登录请求 -->
        <el-button :disabled="loading" class="login-button" size="large" type="primary" @click="loginClick">登录
        </el-button>
      </div>

      <!-- 注册状态下的表单 -->
      <div v-else class="content-view">
        <form>
          <!-- 用户名输入框 -->
          <el-input v-model.trim="inputData.name" class="input-view" clearable placeholder="请输入用户名"
                    size="large"></el-input>
          <!-- 密码输入框 -->
          <el-input v-model.trim="inputData.password" class="input-view" clearable placeholder="请输入密码" size="large"
                    type="password"></el-input>
          <!-- 确认密码输入框 -->
          <el-input
              v-model.trim="inputData.confirmPassword"
              class="input-view"
              clearable
              placeholder="请再次输入密码"
              size="large"
              type="password"
          ></el-input>
        </form>

        <!-- 返回登录按钮 -->
        <el-button :disabled="loading" class="register-button" link size="small" text type="primary"
                   @click="registerClick">返回登录
        </el-button>

        <!-- 提交注册请求的按钮 -->
        <el-button :disabled="loading" class="login-button" size="large" type="primary" @click="loginClick">注册
        </el-button>
      </div>
    </transition>
  </div>
</template>

<script lang="ts" setup>
import {ref, onMounted} from 'vue'; // Vue的响应式引用和生命周期钩子
import {ElMessage} from 'element-plus'; // Element Plus的消息提示组件
import {useRouter} from 'vue-router'; // 路由导航工具

import Token from '@/store/token'; // 用户认证信息存储模块
import TokenApi from '@/api/token'; // 用户认证相关的API请求模块
import TextRain from '@/components/text-rain/index.vue'; // 文字雨组件导入

// 路由实例，用于页面跳转
const router = useRouter();
// 控制按钮加载状态的标志
const loading = ref(false);
// 控制界面展示登录还是注册表单的标志
const isLogin = ref(true);
// 存储用户输入的表单数据（用户名、密码、确认密码）
const inputData = ref({
  name: "",
  password: "",
  confirmPassword: "",
});

// 页面挂载后执行，尝试从缓存中恢复用户名
onMounted(() => {
  let nameCache = Token.getName();
  if (nameCache) {
    inputData.value.name = nameCache;
  }
});

// 处理注册/登录切换逻辑
const registerClick = () => {
  // 清空输入框并切换界面状态
  inputData.value.name = "";
  inputData.value.password = "";
  inputData.value.confirmPassword = "";
  isLogin.value = !isLogin.value;
};

// 处理登录/注册请求
const loginClick = () => {
  // 验证必填项
  if (!inputData.value.name) {
    ElMessage.warning("请输入用户名");
    return;
  }
  if (!inputData.value.password) {
    ElMessage.warning("请输入密码");
    return;
  }

  // 根据当前界面状态执行登录或注册操作
  if (isLogin.value) {
    // 登录
    loading.value = true;
    TokenApi.signIn(inputData.value.name, inputData.value.password)
        .then((res) => {
          // 登录成功，保存令牌并跳转至主页
          Token.setToken(res.data);
          router.push({name: "layout"});
        })
        .finally(() => {
          loading.value = false;
        });
  } else {
    // 注册处理，验证密码一致性
    if (inputData.value.password !== inputData.value.confirmPassword) {
      ElMessage.warning("两次密码不一致");
      return;
    }
    loading.value = true;
    TokenApi.signUp(inputData.value.name, inputData.value.password)
        .then(() => {
          // 注册成功提示
          ElMessage.success("注册成功");
          // 清空重置密码, 保留密码，准备返回登录界面
          inputData.value.confirmPassword = "";
          isLogin.value = true;
        })
        .finally(() => {
          loading.value = false;
        });
  }
};

// 点击标题打开公开博客页面
const publishClick = () => {
  let url = "https://blog.luode.vip";
  window.open(url, "_blank");
};
</script>

<style lang="scss" scoped>
/* 页面登录样式 */
.page-login {
  /* 页面布局定位及尺寸 */
  position: fixed;
  width: 100%;
  height: 100%;
  left: 0;
  right: 0;
  top: 0;
  bottom: 0;
  /* 背景渐变色 */
  background: linear-gradient(0deg, rgba(255, 238, 213, 1) 0%, rgba(148, 210, 233, 1) 70%);
  display: flex;
  flex-direction: column;
  align-items: center;
  user-select: none;
  overflow: auto;

  /* 标题样式 */

  .title-view {
    /* 位置调整 */
    margin-top: 10vh;
    /* 字体样式 */
    font-size: 24px;
    font-weight: bold;
    color: #3f3f3f;
    /* 可点击指针 */
    cursor: pointer;
  }

  /* 内容视图容器 */

  .content-view {
    /* 位置及尺寸调整 */
    margin-top: 20px;
    width: 300px;
    /* 表单布局 */
    display: flex;
    flex-direction: column;
    align-items: flex-end;
    /* 背景色 */
    background: rgba(255, 255, 255, 0.2);
    padding: 30px 20px;

    /* 输入框样式 */

    .input-view {
      width: 100%;
      margin-top: 10px;
    }

    /* 注册按钮样式调整 */

    .register-button {
      margin: 5px 0;
    }

    .login-button {
      width: 100%;
    }

    /* 注册按钮间间距调整 */

    .register-button + .login-button {
      margin-left: 0;
    }
  }

  /* 动态文字雨组件的绝对定位 */

  .text-rain-background {
    position: absolute;
    z-index: -1;
  }
}
</style>
