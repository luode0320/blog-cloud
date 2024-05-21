<template>
  <!-- 页面主体容器，包含文档浏览和编辑的主要界面 -->
  <div class="page-document">
    <!-- 书籍目录组件，接收是否拉伸、加载状态及预览模式等属性 -->
    <book :isStretch="isStretch" :loading="mdLoading" :onlyPreview="onlyPreview" @books="booksFetch"
          @change="bookChange">
    </book>

    <!-- 文档组件，负责展示当前选中的文档内容 -->
    <doc
        ref="docRef"
        :books="books"
        :currentBookId="currentBookId"
        :currentDoc="currentDoc"
        :isStretch="isStretch"
        :onlyPreview="onlyPreview"
        @change="docChange"
        @loading="loadingChange"
    ></doc>

    <!-- 当文档类型为开放API时，展示特定界面 -->
    <div v-if="docType === 'openApi'" class="codemirror-view">
      <!-- 如果只预览，直接显示开放API内容 -->
      <open-api v-if="onlyPreview" :content="currentDoc.content"></open-api>

      <!-- 否则，展示编辑工具栏和CodeMirror编辑器 -->
      <template v-else>
        <!-- 编辑工具栏，包含保存和导出按钮 -->
        <div class="codemirror-toolbar">
          <!-- 自定义SVG图标，用于保存操作 -->
          <div class="icon-outer" title="保存" @click="saveDoc(currentDoc.content)">
            <svg-icon className="icon-save" name="save"></svg-icon>
          </div>
          <!-- 自定义SVG图标，用于导出操作 -->
          <div class="icon-outer" title="导出" @click="exportOpenApi(currentDoc.name, currentDoc.content)">
            <svg-icon className="icon-download" name="download"></svg-icon>
          </div>
        </div>

        <!-- CodeMirror编辑器容器 -->
        <div class="codemirror-inner">
          <!-- CodeMirror编辑器组件，用于编辑文档内容 -->
          <codemirror-editor
              ref="codemirrorRef"
              v-model="currentDoc.content"
              :disabled="onlyPreview || mdLoading"
              :style="{ visibility: codemirrorVisibility }"
              noRadius
              @ready="codemirrorReday"
              @save="saveDoc(currentDoc.content)"
          />
        </div>
      </template>
    </div>

    <!-- 其他文档类型（非开放API），展示Markdown编辑器或预览 -->
    <template v-else>
      <!-- 预览模式下显示Markdown预览 -->
      <md-preview v-if="onlyPreview" :key="'preview' + mdKey" :content="currentDoc.content" class="editor-view"/>

      <!-- 编辑模式下显示Markdown编辑器 -->
      <md-editor
          v-else
          :key="'editor' + mdKey"
          v-model="currentDoc.content"
          v-loading="mdLoading"
          class="editor-view"
          @export="exporMarkdown(currentDoc.name, currentDoc.content)"
          @onSave="saveDoc"
          @onUploadImg="uploadImage"
      />
    </template>
  </div>
</template>

<script lang="ts" setup>
// 导入Vue相关功能以及项目内的组件和服务
import {ref, Ref, onMounted, onBeforeUnmount, nextTick, computed, watch} from "vue";
import MdEditor from "@/components/md-editor";
import MdPreview from "@/components/md-editor/preview";
import CodemirrorEditor from "@/components/codemirror-editor";
import OpenApi from "@/components/open-api/index.vue";
import SvgIcon from "@/components/svg-icon";
import crypto from "crypto-js";

import DocCache from "@/store/doc-cache";
import Token from "@/store/token";
import {host} from "@/config";
import {uploadPicture} from "../picture/util"; // 上传图片
import Book from "./components/book.vue";
import Tree from "./components/tree.vue";
import Doc from "./components/doc.vue";
import {exporMarkdown, exportOpenApi} from "./util";

// 定义组件Props，接收外部传入的配置项
defineProps({
  onlyPreview: {
    type: Boolean,
    default: true,
  },
  isStretch: {
    type: Boolean,
    default: true,
  },
});

// 定义组件内的响应式变量，这些变量用于数据绑定和业务逻辑处理
const docRef = ref<InstanceType<typeof Doc>>(); // 引用Doc组件的实例，用于直接调用其方法
const codemirrorRef = ref(); // 引用Codemirror编辑器组件的实例
const hostUrl = ref(""); // 动态获取或设置基础URL，根据环境不同而变化
const books = ref<Book[]>([]); // 存储书籍列表的响应式数组
const currentBookId = ref(""); // 当前选中书籍的ID
const currentDoc = ref<CurrentDoc>({ // 当前文档的详细信息，包括ID、名称、内容等
  id: "",
  name: "",
  content: "",
  originMD5: "",
  type: "",
  updateTime: "",
});
const mdLoading = ref(false); // 控制Markdown加载状态的标志
const mdKey = ref(0); // 用于强制更新Markdown预览组件的key，以便在文档改变时重新渲染
const codemirrorVisibility = ref("hidden"); // 控制CodeMirror编辑器的初始可见性

// 计算属性，根据当前文档的类型动态返回
const docType = computed(() => {
  return currentDoc.value.type;
});

// 监听器，当文档类型变化且从非"openApi"变为"openApi"时，隐藏CodeMirror编辑器
watch(docType, (newVal, oldVal) => {
  if (oldVal && oldVal !== newVal && newVal === "openApi") {
    codemirrorVisibility.value = "hidden";
  }
});

// 组件挂载后执行的逻辑，设置基础URL并尝试从缓存中恢复当前文档
onMounted(() => {
  hostUrl.value = host;
  DocCache.getDoc().then((res) => {
    if (res) {
      currentDoc.value = res; // 设置从缓存中获取的文档信息
    }
  });
});

// 组件卸载前执行的逻辑，保存当前文档到缓存中（如果用户已登录）
onBeforeUnmount(() => {
  if (Token.getAccessToken()) {
    DocCache.setDoc(currentDoc.value);
  }
});

// 窗口即将刷新时触发，确保文档状态被保存
window.onbeforeunload = () => {
  if (Token.getAccessToken()) {
    DocCache.setDoc(currentDoc.value);
  }
};

// 文档加载状态变化时更新UI显示
const loadingChange = (val: boolean) => {
  mdLoading.value = val;
};

// 处理一级目录选择变化，更新当前选中的书籍ID
const bookChange = (bookId: string) => {
  currentBookId.value = bookId;
};

// 接收一级目录列表更新，并更新组件内的书籍列表
const booksFetch = (bookList: Book[]) => {
  books.value = bookList;
};

// 文档选择或内容变更时的处理逻辑，更新文档信息并考虑是否重新渲染
const docChange = (id: string, name: string, content: string, type: string, updateTime: string, noRender?: boolean) => {
  currentDoc.value.id = id;
  currentDoc.value.name = name;
  currentDoc.value.content = content;
  currentDoc.value.type = type;
  currentDoc.value.originMD5 = crypto.MD5(content).toString();
  currentDoc.value.updateTime = updateTime;
  // 如果不是禁止重新渲染，则强制更新预览
  if (!noRender) {
    mdKey.value++;  // 更新key以强制重绘
    nextTick(() => { // 等待DOM更新后滚动到顶部
      if (codemirrorRef.value) {
        codemirrorRef.value.$el.getElementsByClassName("cm-scroller")[0].scrollTop = 0;
      }
    });
  }
  // 保存文档到缓存
  DocCache.setDoc(currentDoc.value);
};

// Codemirror编辑器加载完毕后的处理，调整滚动位置和显示
const codemirrorReday = () => {
  setTimeout(() => {
    if (codemirrorRef.value) {
      codemirrorRef.value.$el.getElementsByClassName("cm-scroller")[0].scrollTop = 0;
      codemirrorVisibility.value = "unset"; // 显示编辑器
    }
  }, 100);
};

// 图片上传处理函数，接受文件列表并返回上传后的URL列表
const uploadImage = async (files: File[], callback: (urls: string[]) => void) => {
  const pathList: string[] = [];
  for (let file of files) {
    try {
      // // 尝试上传图片并拼接URL
      pathList.push(hostUrl.value + (await uploadPicture(file)));
    } catch (e) {
    }
  }
  // 通知调用者图片上传完成并提供URL列表
  callback(pathList);
};

// 保存文档内容到当前文档实例的方法
const saveDoc = (content: string) => {
  if (mdLoading.value) {
    return;
  }
  docRef.value?.saveDoc(content); // 调用Doc组件的保存方法
};
</script>

<style lang="scss">
.page-document {
  display: flex;
  overflow: auto;

  .editor-view {
    height: 100%;
    flex: 1;
    min-width: 720px;
  }

  .editor-view.md-fullscreen {
    min-width: unset;
  }

  .codemirror-view {
    height: 100%;
    flex: 1;
    min-width: 720px;
    overflow: hidden;
  }

  .codemirror-toolbar {
    height: 34px;
    display: flex;
    align-items: center;
    justify-content: flex-end;
    border: #e6e6e6 1px solid;
    border-bottom: none;
    padding-right: 10px;

    .icon-outer {
      width: 30px;
      height: 24px;
      color: #3f4a54;
      cursor: pointer;

      .icon-save {
        width: 16px;
        height: 16px;
        margin: 4px 7px;
      }

      .icon-download {
        width: 20px;
        height: 20px;
        margin: 2px 5px;
      }
    }

    .icon-outer:hover {
      background: #f6f6f6;
    }
  }

  .codemirror-inner {
    height: calc(100% - 35px);
    overflow: hidden;
  }
}

@media (max-width: 720px) {
  .page-document {
    .editor-view {
      min-width: 100%;

      .catalog-view {
        display: none;
      }
    }

    .codemirror-view {
      min-width: 100%;
    }
  }
}
</style>
