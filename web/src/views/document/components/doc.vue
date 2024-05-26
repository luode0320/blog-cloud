<template>
  <!-- 页面加载指示器，当docLoading为true时显示加载动画 -->
  <div v-loading="docLoading" :class="{ 'page-doc-shrink': !isStretch }" class="page-doc">

    <!-- 如果文档被禁用，显示遮罩层 -->
    <div v-if="docDisabled" class="mask-view"></div>

    <!-- 创建新文档的弹出层 -->
    <el-popover v-if="!onlyPreview" :visible="addDocVisible" placement="bottom" trigger="click" width="240px">
      <!-- 文档名称输入框 -->
      <el-input v-model="newDocName" placeholder="请输入文档名称" style="margin-right: 10px"></el-input>
      <!-- 确认与取消按钮 -->
      <div style="display: flex; margin-top: 8px; justify-content: flex-end">
        <el-button size="small" @click="addDocCancel">取消</el-button>
        <el-button size="small" type="primary" @click="addDocSave">确定</el-button>
      </div>
      <!-- 弹出层的触发元素，即创建文档按钮 -->
      <template #reference>
        <el-button :icon="Plus" class="create-button" link size="large" type="primary" @click="addDocVisible = true">
          创建文档
        </el-button>
      </template>
    </el-popover>

    <!-- 当文档功能被禁用时，显示“文档选择”替代按钮 -->
    <el-button v-else class="create-button" link size="large" type="primary">文档选择</el-button>

    <!-- 滚动视图容器，用于容纳文档列表 -->
    <el-scrollbar ref="scrollRef" class="scroll-view">
      <!-- 循环遍历文档列表，渲染每一个文档项 -->
      <div
          v-for="item in docs"
          :key="item.id"
          :class="docIdTemp === item.id || (!docIdTemp && currentDoc.id === item.id) ? 'selected' : ''"
          class="item-view"
          @click="docClick(item)"
      >
        <!-- 文档名称 -->
        <text-tip :content="item.name"></text-tip>
        <!-- 更新时间 -->
        <div class="sub-text">{{ formatTime(item.updateTime, "YYYY-MM-DD HH:mm:ss") }}</div>
        <!-- 发布状态图标，点击可复制 -->
        <!--        <div v-if="item.published" class="published-view" title="已发布" @click.stop="copyPublishedClick(item)"></div>-->

        <!-- 操作下拉菜单，包含修改和删除文档选项 -->
        <el-dropdown v-if="!onlyPreview && item.id" trigger="click">
          <el-icon class="setting-button" title="操作" @click.stop="() => {}">
            <Tools/>
          </el-icon>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item style="user-select: none" @click="updateDocClick(item)">修改文档</el-dropdown-item>
              <el-dropdown-item style="user-select: none" @click="deleteDocClick(item)">删除文档</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </el-scrollbar>

    <!-- 对话框，用于创建或更新文档信息 -->
    <el-dialog
        v-model="dialog.visible"
        :before-close="dialogClose"
        :show-close="false"
        :title="dialog.isAdd ? '创建文档' : '更新文档'"
        width="400px"
    >
      <!-- 表单内容，包括文档名称、类型、是否公开发布和所属一级目录的选择 -->
      <el-form label-width="70px" size="large">
        <el-form-item label="文档名称">
          <el-input v-model="dialog.condition.name" placeholder="请输入文档名称" style="width: 100%"></el-input>
        </el-form-item>
        <!--        <el-form-item label="所属目录">-->
        <!--          <el-select v-model="dialog.condition.bookId" style="width: 100%">-->
        <!--            <el-option v-for="item in books" :key="item.id" :label="item.name" :value="item.id"></el-option>-->
        <!--          </el-select>-->
        <!--        </el-form-item>-->
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button :loading="dialog.loading" @click="dialogClose">取消</el-button>
          <el-button :loading="dialog.loading" type="primary" @click="dialogSave">保存</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script lang="ts" setup>
import {ref, Ref, onMounted, watch, PropType, nextTick} from "vue";
import {ElMessage, ElMessageBox} from "element-plus"; // Element Plus的消息提示和对话框组件
import {Plus, Tools} from "@element-plus/icons-vue"; // Element Plus的图标
import copy from "copy-to-clipboard"; // 复制到剪贴板的库
import crypto from "crypto-js"; // 加密库，用于MD5计算
import NProgress from "nprogress"; // 进度条插件
import "nprogress/nprogress.css"; // 进度条样式

import TextTip from "@/components/text-tip"; // 自定义文本提示组件
import DocumentApi from "@/api/document"; // 文档相关的API请求模块
import {formatTime} from "@/utils"; // 时间格式化工具函数

// 定义一些响应式变量
const hostUrl = ref(location.origin); // 当前页面的主机地址
const docs: Ref<Doc[]> = ref([]); // 文档列表数据
const docLoading = ref(false); // 是否正在加载文档列表
const docDisabled = ref(false); // 是否禁用文档操作
const addDocVisible = ref(false); // 添加文档弹窗是否可见
const newDocName = ref(""); // 新建文档的名称
const newDocType = ref("md"); // 新建文档的类型，默认markdown
const docIdTemp = ref(""); // 临时存储点击的文档ID，用于高亮处理
const dialog = ref({ // 弹窗的状态与数据
  isAdd: true, // 是否为新增模式
  loading: false, // 弹窗内的加载状态
  visible: false, // 弹窗是否显示
  condition: { // 弹窗表单中的数据条件
    id: "", // 文档ID
    name: "", // 文档名称
    content: "", // 文档内容
    bookId: "", // 所属书籍ID
    type: "md", // 文档类型
    published: false, // 是否已发布
  },
});
const scrollRef = ref(); // 滚动视图的引用

// 定义组件对外的事件发射器
const emit = defineEmits<{
  change: [id: string, name: string, content: string, type: string, updateTime: string, noRender?: boolean]; // 文档信息变更事件
  loading: [val: boolean]; // 加载状态变更事件
}>();

// 定义组件接收的属性
const props = defineProps({
  onlyPreview: { // 是否仅预览模式
    type: Boolean,
    default: true,
  },
  isStretch: { // 是否拉伸模式
    type: Boolean,
    default: true,
  },
  currentBookId: { // 当前选中书籍的ID
    type: String,
    default: "",
  },
  currentDoc: { // 当前文档的信息
    type: Object as PropType<CurrentDoc>,
    default: {
      id: "",
      content: "",
      originMD5: "",
      updateTime: "",
    },
  },
  books: { // 图书列表
    type: Array as PropType<Book[]>,
    default: [],
  },
});

// 监听currentBookId的变化，当它变化时重新查询文档列表
watch(
    () => props.currentBookId,
    (val) => {
      queryDocs(val);
    }
);

// 监听docLoading的变化，同步更新加载状态
watch(docLoading, (val) => {
  emit("loading", val);
});

// 组件挂载后立即执行查询文档列表
onMounted(() => {
  queryDocs(props.currentBookId);
});

// 查询文档列表的方法
const queryDocs = (bookId: string) => {
  // 重置添加文档弹窗，清空加载状态，发起API请求获取文档列表
  addDocCancel();
  docLoading.value = true;
  DocumentApi.list(bookId)
      .then((res) => {
        // 检查当前文档是否有更新，如有则通知父组件更新
        for (let item of res.data) {
          if (item.id === props.currentDoc.id) {
            if (String(item.updateTime) !== props.currentDoc.updateTime) {
              emitDoc("", "", "", "", "");
            }
            break;
          }
        }
        docs.value = res.data; // 更新文档列表
        // 滚动到当前选中文档的位置
        nextTick(() => {
          scrollRef.value.$el.getElementsByClassName("item-view selected")[0]?.scrollIntoView();
        });
      })
      .finally(() => {
        docLoading.value = false; // 加载完成，恢复加载状态
      });
};

// 校验文档是否已被修改，如果修改则提示用户确认
const checkDocChange = () => {
  return new Promise((resolve, reject) => {
    // 使用MD5校验内容是否有更改
    if (props.currentDoc.originMD5 && crypto.MD5(props.currentDoc.content).toString() !== props.currentDoc.originMD5) {
      // 弹出确认框询问用户是否放弃更改
      ElMessageBox.confirm("文档未保存，是否继续？", "提示", {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning",
      })
          .then(() => {
            resolve(null);
          })
          .catch(() => {
            reject();
          });
    } else {
      resolve(null); // 未修改直接通过
    }
  });
};

// 发送文档信息变更的事件到父组件
const emitDoc = (id: string, name: string, content: string, type: string, updateTime: string, noRender?: boolean) => {
  emit("change", id, name, content, type, updateTime, noRender);
};

// 点击文档项的处理逻辑
const docClick = async (doc: Doc) => {
  // 检查文档是否被修改
  await checkDocChange();
  docIdTemp.value = doc.id; // 设置临时ID，准备高亮
  docDisabled.value = true; // 禁用操作，防止重复点击
  emit("loading", true); // 开始加载
  NProgress.start(); // 显示进度条
  try {
    const res = await DocumentApi.get(doc.id); // 获取文档详情
    // 通知父组件更新文档信息
    emitDoc(res.data.id, res.data.name, res.data.content, res.data.type!, String(res.data.updateTime));
  } finally {
    docIdTemp.value = ""; // 清除临时ID
    docDisabled.value = false; // 启用操作
    emit("loading", false); // 结束加载
    NProgress.done(); // 隐藏进度条
  }
};

// 添加文档保存逻辑
const addDocSave = () => {
  const name = newDocName.value.trim();
  if (!name) {
    ElMessage.warning("请填写文档名称");
    return;
  }
  checkDocChange().then(() => {
    docLoading.value = true; // 开始加载
    // 发起添加文档的API请求
    DocumentApi.add({id: "", name, content: "", type: newDocType.value, bookId: props.currentBookId})
        .then((res) => {
          ElMessage.success("创建成功");
          // 通知父组件文档已创建
          emitDoc(res.data.id, res.data.name, res.data.content, res.data.type!, String(res.data.updateTime));
          queryDocs(props.currentBookId); // 重新查询文档列表
        })
        .finally(() => {
          docLoading.value = false; // 结束加载
        });
  });
};

/**
 * 点击添加文档取消
 */
const addDocCancel = () => {
  addDocVisible.value = false;
  newDocName.value = "";
  newDocType.value = "md";
};

/**
 * 点击修改文档
 */
const updateDocClick = (doc: Doc) => {
  dialog.value.condition.id = doc.id;
  dialog.value.condition.name = doc.name;
  dialog.value.condition.content = "";
  dialog.value.condition.bookId = doc.bookId;
  dialog.value.condition.type = doc.type!;
  dialog.value.condition.published = doc.published!;
  dialog.value.isAdd = false;
  dialog.value.visible = true;
};

/**
 * 点击删除文档
 */
const deleteDocClick = (doc: Doc) => {
  ElMessageBox.confirm("是否删除文档：" + doc.name + "？", "提示", {
    confirmButtonText: "确定",
    cancelButtonText: "取消",
    type: "warning",
  }).then(() => {
    DocumentApi.delete(doc.id).then(() => {
      ElMessage.success("删除成功");
      if (props.currentDoc.id === doc.id) {
        emitDoc("", "", "", "", "");
      }
      queryDocs(props.currentBookId);
    });
  });
};

/**
 * 点击发布地址
 * @param doc
 */
const copyPublishedClick = (doc: Doc) => {
  let url = hostUrl.value + "/#/open/document?id=" + doc.id;
  const result = copy(url);
  if (result) {
    ElMessage.success("发布地址已复制到剪切板");
  } else {
    ElMessage.error("复制到剪切板失败，地址：" + url);
  }
};

/**
 * 弹窗关闭
 */
const dialogClose = () => {
  if (dialog.value.loading) {
    return;
  }
  dialog.value.condition.id = "";
  dialog.value.condition.name = "";
  dialog.value.condition.content = "";
  dialog.value.condition.bookId = "";
  dialog.value.condition.type = "md";
  dialog.value.condition.published = false;
  dialog.value.isAdd = true;
  dialog.value.visible = false;
};

/**
 * 弹窗保存
 */
const dialogSave = () => {
  if (dialog.value.isAdd) {
    // 新增文档
    docLoading.value = true;
    DocumentApi.add(dialog.value.condition)
        .then((res) => {
          ElMessage.success("创建成功");
          emitDoc(res.data.id, res.data.name, res.data.content, res.data.type!, String(res.data.updateTime));
          docLoading.value = false;
          dialogClose();
          queryDocs(props.currentBookId);
        })
        .catch(() => {
          docLoading.value = false;
        });
  } else {
    // 更新基本信息
    let name = String(dialog.value.condition.name).trim();
    if (!name) {
      ElMessage.warning("请填写文档名称");
      return;
    }
    dialog.value.condition.name = name;
    dialog.value.loading = true;
    DocumentApi.update(dialog.value.condition)
        .then(() => {
          ElMessage.success("更新成功");
          dialog.value.loading = false;
          dialogClose();
          queryDocs(props.currentBookId);
        })
        .catch(() => {
          dialog.value.loading = false;
        });
  }
};

/**
 * 保存文档
 */
const saveDoc = (content: string) => {
  if (props.currentDoc.id !== "") {
    // 更新文档内容
    docLoading.value = true;
    DocumentApi.updateContent({id: props.currentDoc.id, name: "", content: content, bookId: ""})
        .then((res) => {
          ElMessage.success("保存成功");
          emitDoc(res.data.id, res.data.name, res.data.content, res.data.type!, String(res.data.updateTime), true);
          // 更新当前文档的更新时间
          for (let item of docs.value) {
            if (item.id === res.data.id) {
              item.updateTime = res.data.updateTime;
              break;
            }
          }
        })
        .finally(() => {
          docLoading.value = false;
        });
  } else {
    // 新增
    dialog.value.condition.id = "";
    dialog.value.condition.name = "";
    dialog.value.condition.content = content;
    dialog.value.condition.bookId = "";
    dialog.value.condition.type = "md";
    dialog.value.condition.published = false;
    dialog.value.isAdd = true;
    dialog.value.visible = true;
  }
};

defineExpose({saveDoc});
</script>

<style lang="scss">
.page-doc {
  height: 100%;
  min-width: 260px;
  width: 260px;
  background: #fafafa;
  display: flex;
  flex-direction: column;
  overflow-x: hidden;
  position: relative;
  transition: margin-left 0.3s;

  .mask-view {
    position: absolute;
    width: 100%;
    height: 100%;
    z-index: 1000;
  }

  .create-button {
    height: 60px;
    border-bottom: 1px solid #e6e6e6 !important;
  }

  .el-button--large [class*="el-icon"] + span {
    margin-left: 3px;
  }

  .scroll-view {
    color: #595959;
    font-size: 13px;

    .item-view {
      display: flex;
      align-items: center;
      justify-content: space-between;
      padding: 18px 15px;
      cursor: pointer;
      border-left: 3px #fafafa solid;
      transition: 0.05s;
      border-bottom: 1px solid #eaeaea;
      position: relative;

      .update-view {
        display: flex;
        align-items: center;
      }

      .sub-text {
        position: absolute;
        font-size: 12px;
        bottom: 3px;
        right: 20px;
        color: #ccc;
      }

      .published-view {
        position: absolute;
        top: 0;
        right: 0;
        width: 0;
        height: 0;
        border-top: 20px solid skyblue;
        border-left: 20px solid transparent;
      }
    }

    .item-view:hover {
      background: #e6e6e6;
      border-left-color: #e6e6e6;
    }

    .item-view.selected {
      background: #e6e6e6;
      border-left-color: #0094c1;
    }

    .setting-button {
      margin-left: 10px;
      color: #595959;
    }

    .setting-button:hover {
      color: #777;
    }
  }

  .el-loading-mask {
    background: #fafafa;
  }
}

.page-doc-shrink {
  margin-left: -260px;
}

@media (max-width: 480px) {
  .page-doc {
    min-width: 55%;
    width: 55%;
  }

  .page-doc-shrink {
    margin-left: -55%;
  }
}
</style>
