<template>
  <!-- 页面加载指示器 -->
  <div v-loading="bookLoading" :class="{ 'page-book-shrink': !isStretch }" class="page-book">

    <!-- 弹出层，用于创建新一级目录 -->
    <div v-if="loading" class="mask-view"></div>
    <el-popover v-if="!onlyPreview" :visible="addBookVisible" placement="bottom" trigger="click" width="200px">
      <!-- 输入框和确认/取消按钮 -->
      <el-input v-model="newBookName" placeholder="请输入一级目录" style="margin-right: 10px"></el-input>
      <div style="display: flex; margin-top: 8px; justify-content: flex-end">
        <el-button size="small" @click="addBookCancel">取消</el-button>
        <el-button size="small" type="primary" @click="addBookSave">确定</el-button>
      </div>

      <template #reference>
        <el-button :icon="Plus" class="create-button" link size="large" type="warning" @click="addBookVisible = true">
          创建一级目录
        </el-button>
      </template>

    </el-popover>

    <!-- 滚动视图容器 -->
    <el-scrollbar class="scroll-view">
      <div v-for="(item, index) in books" :key="item.id">
        <!-- 循环渲染每一本书籍项 -->
        <div v-if="!item.parentId || currentBookId == item.parentId || currentBookId == item.id"
             :class="[(currentBookId === '' && index === 0) || currentBookId === item.id ? 'selected' : '', 'item-view']"
             @click="bookClick(item,index)">

          <!-- 修改一级目录的输入框和按钮 -->
          <div v-if="updateBookId && updateBookId === item.id" class="update-view">
            <el-input v-model="updateBookName" placeholder="请输入一级目录"></el-input>
            <el-button :icon="CircleCheckFilled" link style="margin-left: 12px" type="success"
                       @click="updateBookSave"></el-button>
            <el-button :icon="CircleCloseFilled" link type="danger" @click="updateBookCancel"></el-button>
          </div>
          <!-- 若非修改状态则显示书名 -->
          <text-tip v-else :content="item.name"></text-tip>

          <!-- 这里可以设计二级目录的简单展示或操作 -->
          <el-popover v-if="addSecondLevelVisible && currentBookId === item.id" :visible="addSecondLevelVisible"
                      placement="bottom" trigger="click"
                      width="200px">
            <!-- 输入框和确认/取消按钮 -->
            <el-input v-model="newSecondLevelName" placeholder="请输入二级目录" style="margin-right: 10px"></el-input>
            <div style="display: flex; margin-top: 8px; justify-content: flex-end">
              <el-button size="small" @click="addTwoBookCancel">取消</el-button>
              <el-button size="small" type="primary" @click="addSecondLevelSave">确定</el-button>
            </div>

            <template #reference>
              <el-button link size="large">
              </el-button>
            </template>
          </el-popover>

          <!-- 非预览模式下的操作下拉菜单 -->
          <el-dropdown v-if="!onlyPreview && item.id" trigger="click">
            <el-icon class="setting-button" title="操作" @click.stop="() => {}">
              <Tools/>
            </el-icon>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item v-if="!item.parentId" style="user-select: none"
                                  @click="addSecondLevelBookClick(item)">添加二级目录
                </el-dropdown-item>
                <el-dropdown-item style="user-select: none" @click="updateBookClick(item)">修改一级目录
                </el-dropdown-item>
                <el-dropdown-item style="user-select: none" @click="deleteBookClick(item)">删除一级目录
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </div>
    </el-scrollbar>
  </div>
</template>

<script lang="ts" setup>
import {ref, Ref, onMounted, watch} from "vue";
import {ElMessage, ElMessageBox} from "element-plus";
import {Plus, Tools, CircleCheckFilled, CircleCloseFilled} from "@element-plus/icons-vue";

import TextTip from "@/components/text-tip";
import BookApi from "@/api/book";

// 定义响应式数据
const books: Ref<Book[]> = ref([]);
const twoBooks: Ref<Book[]> = ref([]);
const bookLoading = ref(false); // 书籍列表加载状态
const currentBook: Ref<Book> = ref(null); // 当前选中的一级目录
const currentBookId = ref(""); // 当前选中的一级目录ID
const currentTwoBookId = ref(""); // 当前选中的一级目录ID
const addBookVisible = ref(false); // 添加一级目录弹窗是否显示
const newBookName = ref(""); // 新增一级目录的名称
const updateBookId = ref(""); // 正在修改的一级目录ID
const updateBookName = ref(""); // 修改后的一级目录名称

// 定义响应式数据结构调整，为每个一级目录添加二级目录数组

// 定义二级目录添加相关数据
const addSecondLevelVisible = ref(false);
const selectedBookForSecondLevel = ref(null as Book | null);
const newSecondLevelName = ref('');

// 定义事件发射器
const emit = defineEmits<{ change: [bookId: string]; books: [bookList: Book[]] }>();

defineProps({
  onlyPreview: { // 是否只预览模式
    type: Boolean,
    default: true,
  },
  isStretch: { // 页面是否伸展
    type: Boolean,
    default: true,
  },
  loading: { // 页面整体加载状态
    type: Boolean,
    default: false,
  },
});

// 监听currentBookId变化，触发外部change事件
watch(currentBookId, (val) => {
  emit("change", val);
});

// 组件挂载后查询所有一级目录
onMounted(() => {
  queryBooks();
});

// 查询一级目录列表
const queryBooks = async () => {
  // 点击添加一级目录取消
  addBookCancel();
  // 点击添加二级目录取消
  addTwoBookCancel()
  // 点击修改一级目录取消，执行 updateBookCancel 方法
  updateBookCancel();
  // 开始加载动画
  bookLoading.value = true;
  // 调用API获取书籍列表
  try {
    // 调用API获取书籍列表
    const res = await BookApi.list();
    // 更新书籍列表
    // todo 这里既有一个也有二级
    console.log("res.data:", res.data)
    books.value = res.data;
    // 触发books事件传递数据给父组件
    emit('books', res.data);

    // 如果当前没有选中目录且存在目录，则默认选中第一个
    if (!currentBookId.value && books.value.length) {
      currentBookId.value = books.value[0].id;
      currentBook.value = books.value[0];
    }
  } finally {
    // 结束加载动画
    bookLoading.value = false;
  }
};

// 点击一级目录的处理
const bookClick = (book: Book, index: number) => {
  addTwoBookCancel()
  // 如果正在修改目录则不响应点击
  if (updateBookId.value) {
    return;
  }
  currentBookId.value = book.id;
  currentBook.value = book;
};

// 点击二级目录的处理
const bookTwoClick = (book: Book, index: number) => {
  addTwoBookCancel()
  // 如果正在修改目录则不响应点击
  if (updateBookId.value) {
    return;
  }
  currentBookId.value = book.id;
  currentBook.value = book;
};

// 添加一级目录保存逻辑
const addBookSave = () => {
  // 检查输入是否为空
  let name = String(newBookName.value).trim();
  if (!name) {
    ElMessage.warning("请填写一级目录");
    return;
  }
  // 开始保存动画
  bookLoading.value = true;
  BookApi.add({id: "", parentId: "", name: name})
      .then(() => {
        ElMessage.success("创建成功");
        queryBooks();
      })
      .catch(() => {
        bookLoading.value = false;
      });
};

// 定义添加二级目录的方法
const addSecondLevelBookClick = (book: Book) => {
  selectedBookForSecondLevel.value = book;
  addSecondLevelVisible.value = true;
};

const addSecondLevelSave = () => {
  let name = String(newSecondLevelName.value).trim();
  if (!name) {
    ElMessage.warning("请填写二级目录名称");
    return;
  }

  bookLoading.value = true;
  BookApi.add({id: "", parentId: selectedBookForSecondLevel.value.id, name})
      .then(() => {
        ElMessage.success("二级目录创建成功");
        queryBooks(); // 刷新一级和二级目录列表
      })
      .catch(() => {
        bookLoading.value = false;
      })
      .finally(() => {
        addSecondLevelVisible.value = false;
        newSecondLevelName.value = "";
      });
};

/**
 * 点击添加一级目录取消
 */
const addBookCancel = () => {
  addBookVisible.value = false;
  newBookName.value = "";
};

/**
 * 点击添加二级目录取消
 */
const addTwoBookCancel = () => {
  addSecondLevelVisible.value = false;
  newSecondLevelName.value = "";
};

/**
 * 点击修改一级目录
 */
const updateBookClick = (book: Book) => {
  updateBookId.value = book.id;
  updateBookName.value = book.name;
};

/**
 * 点击修改一级目录保存
 */
const updateBookSave = () => {
  let name = String(updateBookName.value).trim();
  if (!name) {
    ElMessage.warning("请填写一级目录");
    return false;
  }
  bookLoading.value = true;
  BookApi.update({id: updateBookId.value, parentId: "", name: name})
      .then(() => {
        ElMessage.success("修改成功");
        queryBooks();
      })
      .catch(() => {
        bookLoading.value = false;
      });
};

/**
 * 点击修改一级目录取消
 */
const updateBookCancel = () => {
  updateBookId.value = "";
  updateBookName.value = "";
};

/**
 * 点击删除目录
 */
const deleteBookClick = (book: Book) => {
  ElMessageBox.confirm("是否删除目录：" + book.name + "？", "提示", {
    confirmButtonText: "确定",
    cancelButtonText: "取消",
    type: "warning",
  }).then(() => {
    BookApi.delete(book.id).then(() => {
      ElMessage.success("删除成功");
      if (currentBookId.value === book.id) {
        currentBookId.value = "";
      }
      queryBooks();
    });
  });
};
</script>

<style lang="scss">
.page-book {
  height: 100%;
  min-width: 220px;
  width: 220px;
  background: #404040;
  display: flex;
  flex-direction: column;
  overflow-x: hidden;
  transition: margin-left 0.3s;

  .mask-view {
    position: absolute;
    width: 100%;
    height: 100%;
    z-index: 1000;
  }

  .create-button {
    height: 60px;
    border-bottom: 1px solid #555 !important;
  }

  .el-button--large [class*="el-icon"] + span {
    margin-left: 3px;
  }

  .scroll-view {
    color: #f2f2f2;
    font-size: 13px;

    .item-view {
      display: flex;
      align-items: center;
      justify-content: space-between;
      padding: 18px 15px;
      cursor: pointer;
      border-left: 3px #404040 solid;
      transition: 0.05s;
      border-bottom: 1px solid #555;

      .update-view {
        display: flex;
        align-items: center;
      }
    }

    .item-view:hover {
      background: #666;
      border-left-color: #666;
    }

    .item-view.selected {
      background: #666;
      border-left-color: #e6a23c;
    }

    .setting-button {
      margin-left: 10px;
      color: #f2f2f2;
    }

    .setting-button:hover {
      color: #ccc;
    }
  }

  .el-loading-mask {
    background: #404040;
  }
}

.page-book-shrink {
  margin-left: -220px;
}

@media (max-width: 480px) {
  .page-book {
    min-width: 45%;
    width: 45%;
  }

  .page-book-shrink {
    margin-left: -45%;
  }
}
</style>
