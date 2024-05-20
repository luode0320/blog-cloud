<template>
  <!-- 页面布局容器 -->
  <div class="page-layout">
    <!-- 顶部视图，包含左侧操作栏和右侧用户信息下拉菜单 -->
    <div class="top-view">
      <!-- 左侧视图，含有标题和菜单 -->
      <div class="left-view">
        <!-- 标题区域，含有一个可点击的图标和博客标题 -->
        <div :style="isDocument ? 'cursor: pointer' : ''" :title="isDocument ? (isStretch ? '收起侧栏' : '弹出侧栏') : ''"
             class="title-view" @click="iconClick">
          <span>
            <!-- 图标点击事件会根据当前页面状态切换侧边栏显示或隐藏 -->
            <svg-icon customStyle="width: 20px; height: 20px; margin: 5px 5px 0 0" name="md"></svg-icon>
          </span>
          <!-- 标题 -->
          <span>
            博客
          </span>
        </div>

        <!-- 菜单项，使用路由链接进行导航 -->
        <div class="menu-view">
          <router-link to="/document">文档</router-link>
          <router-link to="/picture">图片</router-link>
          <router-link to="/tool">工具</router-link>
        </div>
      </div>

      <!-- 右侧视图，包含用户姓名和下拉菜单 -->
      <el-dropdown class="right-view">
        <!-- 显示用户姓名，点击后展示下拉菜单 -->
        <div class="text-view">{{ name }}</div>
        <!-- 下拉菜单内容定义 -->
        <template #dropdown>
          <el-dropdown-menu>
            <!-- 菜单项：博客主页 -->
            <el-dropdown-item style="user-select: none" @click="publishClick">博客主页</el-dropdown-item>
            <!-- 修改密码对话框触发 -->
            <el-dropdown-item style="user-select: none" @click="dialogVisible = true">修改密码</el-dropdown-item>
            <!-- 退出登录功能 -->
            <el-dropdown-item style="user-select: none" @click="logout">退出登录</el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>

    <!-- 动态路由视图，传递是否拉伸和预览模式的状态 -->
    <router-view :isStretch="isStretch" :onlyPreview="onlyPreview" class="content-view"></router-view>
    <!-- 修改密码对话框 -->
    <el-dialog v-model="dialogVisible" :before-close="dialogClose" :show-close="false" title="修改密码" width="400px">
      <!-- 表单用于输入原密码和新密码 -->
      <form>
        <el-input v-model.trim="form.password" clearable placeholder="请输入原密码" size="large"
                  type="password"></el-input>
        <el-input v-model.trim="form.newPassword" clearable placeholder="请输入新密码" size="large"
                  style="margin: 10px 0" type="password"></el-input>
        <el-input v-model.trim="form.confirmPassword" clearable placeholder="请再次输入密码" size="large"
                  type="password"></el-input>
      </form>

      <!-- 对话框底部按钮 -->
      <template #footer>
        <span class="dialog-footer">
        <!-- 取消按钮 -->
          <el-button :loading="dialogLoading" @click="dialogClose">取消</el-button>
          <!-- 保存按钮，执行密码更新操作 -->
          <el-button :loading="dialogLoading" type="primary" @click="updatePassword">保存</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script lang="ts" setup>
import {ref, watch} from "vue";
import {ElMessage, ElMessageBox} from "element-plus";
import SvgIcon from "@/components/svg-icon";  // SVG图标组件
import Token from "@/store/token";  // 用户Token存储
import TokenApi from "@/api/token"; // Token API
import UserApi from "@/api/user"; // 用户API
import router from "@/router";  // 路由
import DocCache from "@/store/doc-cache"; // 文档缓存

// 初始化引用
const name = ref(Token.getName());  // 获取用户姓名
const dialogVisible = ref(false); // 控制修改密码对话框显示, 默认不显示
const dialogLoading = ref(false); // 加载Loading指示器，用于表单提交时
const form = ref({password: "", newPassword: "", confirmPassword: ""}); // 一个默认的表单对象
const onlyPreview = ref(false); // 是否只读预览模式
const isStretch = ref(true);  // 侧边栏是否展开
const isDocument = ref(router.currentRoute.value.name === "document");  // 判断是否在文档页面

// 监听路由变化，更新文档页面标识
watch(
    // 当前路由的页面, 是否是文档页面
    () => router.currentRoute.value.name,
    (val) => {
      if (val === "document") {
        isDocument.value = true;
      } else {
        isDocument.value = false;
      }
    }
);

// 退出登录功能
const logout = () => {
  ElMessageBox.confirm("是否退出登录？", "提示", {
    confirmButtonText: "退出登录",
    cancelButtonText: "取消",
    type: "info",
  }).then(() => {
    DocCache.removeDoc(); // 清除文档缓存
    TokenApi.signOut(); // 执行登出操作
    Token.removeToken();  // 移除Token
  }).catch(() => {
  });
};

// 更新密码功能
const updatePassword = () => {
  if (form.value.password === "" || form.value.newPassword === "" || form.value.confirmPassword === "") {
    ElMessage.warning("请填写密码");
    return;
  }
  if (form.value.newPassword !== form.value.confirmPassword) {
    ElMessage.warning("两次密码不一致");
    return;
  }

  dialogLoading.value = true;
  UserApi.updatePassword(form.value.password, form.value.newPassword)
      .then((res) => {
        ElMessage.success("修改成功");
        dialogLoading.value = false;
        dialogClose();
      })
      .catch(() => {
        dialogLoading.value = false;
      });
};

// 关闭对话框
const dialogClose = () => {
  // 如果还在转圈圈加载, 则不动
  if (dialogLoading.value) {
    return;
  }
  dialogVisible.value = false; // 关闭修改密码对话框
  form.value.password = ""; // 清除密码
  form.value.newPassword = "";  // 清楚新密码
  form.value.confirmPassword = "";  // 清楚确认密码
};

// 点击标题图标切换侧边栏显示
const iconClick = () => {
  if (isDocument.value) {
    isStretch.value = !isStretch.value;
  }
};

// 点击公开文档, 在新窗口打开公开文档页面
const publishClick = () => {
  let url = "https://blog.luode.vip";
  window.open(url, "_blank");
};
</script>

<style lang="scss" scoped>
// 页面布局样式
.page-layout {
  position: fixed;
  width: 100%;
  height: 100%;
  left: 0;
  right: 0;
  top: 0;
  bottom: 0;
}

// 顶部视图样式
.top-view {
  height: 49px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  user-select: none;
  border-bottom: 1px #eee solid;
  padding: 0 2%;
  white-space: nowrap;

// 左右侧视图样式设置

  .left-view {
    display: flex;
    align-items: center;
    height: 100%;
    flex: 1;


    .title-view {
      font-size: 18px;
      font-weight: bold;
      display: flex;
      align-items: center;
    }

  // 标题和菜单样式

    .menu-view {
      margin-left: 2%;
      display: flex;
      align-items: center;
      height: 100%;
      font-size: 14px;

      a {
        text-decoration: none;
        color: #303133;
        height: 100%;
        display: flex;
        align-items: center;
        justify-content: center;
        border-bottom: 2px solid transparent;
        box-sizing: border-box;
        padding: 4px 20px 0 20px;
        transition: all 0.3s;
      }

      a:hover {
        background: rgba(0, 148, 193, 0.1);
      }

      .router-link-active {
        color: #0094c1;
        border-color: #0094c1;
      }
    }
  }

  .right-view {
    height: 100%;

    .text-view {
      height: 100%;
      display: flex;
      align-items: center;
      justify-content: center;
      cursor: pointer;
      padding: 0 20px;
      transition: all 0.3s;
      color: #303133;
      outline: none;
    }

    .text-view:hover {
      color: #0094c1;
    }
  }
}

// 内容区域样式
.content-view {
  width: 100%;
  height: calc(100% - 50px);
  overflow: auto;
  background-color: #fcfcfc;
}
</style>
