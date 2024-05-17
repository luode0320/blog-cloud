<template>
  <div v-loading="bookLoading" :class="{ 'page-book-shrink': !isStretch }" class="page-book">
    <div v-if="loading" class="mask-view"></div>

    <el-popover v-if="!onlyPreview" :visible="addBookVisible" placement="bottom" trigger="click" width="200px">
      <el-input v-model="newBookName" placeholder="请输入一级目录" style="margin-right: 10px"></el-input>
      <div style="display: flex; margin-top: 8px; justify-content: flex-end">
        <el-button size="small" @click="addBookCancel">取消</el-button>
        <el-button size="small" type="primary" @click="addBookSave">确定</el-button>
      </div>
      <template #reference>
        <el-button :icon="Plus" class="create-button" link size="large" type="warning" @click="addBookVisible = true">创建一级目录</el-button>
      </template>
    </el-popover>

    <el-button v-else class="create-button" link size="large" type="warning">一级目录选择</el-button>

    <el-scrollbar class="scroll-view">
      <div v-for="item in books" :key="item.id" :class="currentBookId === item.id ? 'selected' : ''" class="item-view" @click="bookClick(item)">
        <div v-if="updateBookId && updateBookId === item.id" class="update-view">
          <el-input v-model="updateBookName" placeholder="请输入一级目录"></el-input>
          <el-button :icon="CircleCheckFilled" link style="margin-left: 12px" type="success" @click="updateBookSave"></el-button>
          <el-button :icon="CircleCloseFilled" link type="danger" @click="updateBookCancel"></el-button>
        </div>
        <text-tip v-else :content="item.name"></text-tip>

        <el-dropdown v-if="!onlyPreview && item.id" trigger="click">
          <el-icon class="setting-button" title="操作" @click.stop="() => {}"><Tools /></el-icon>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item style="user-select: none" @click="updateBookClick(item)">修改一级目录</el-dropdown-item>
              <el-dropdown-item style="user-select: none" @click="deleteBookClick(item)">删除一级目录</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>

      </div>
    </el-scrollbar>
  </div>
</template>

<script lang="ts" setup>
import { ref, Ref, onMounted, watch } from "vue";
import { ElMessage, ElMessageBox } from "element-plus";
import { Plus, Tools, CircleCheckFilled, CircleCloseFilled } from "@element-plus/icons-vue";
import TextTip from "@/components/text-tip";
import BookApi from "@/api/book";

const books: Ref<Book[]> = ref([]);
const bookLoading = ref(false);
const currentBookId = ref("");
const addBookVisible = ref(false);
const newBookName = ref("");
const updateBookId = ref("");
const updateBookName = ref("");

const emit = defineEmits<{ change: [bookId: string]; books: [bookList: Book[]] }>();

defineProps({
  onlyPreview: {
    type: Boolean,
    default: true,
  },
  isStretch: {
    type: Boolean,
    default: true,
  },
  loading: {
    type: Boolean,
    default: false,
  },
});

watch(currentBookId, (val) => {
  emit("change", val);
});

onMounted(() => {
  queryBooks();
});

/**
 * 查询一级目录列表
 */
const queryBooks = () => {
  // 点击添加一级目录取消
  addBookCancel();
  // 点击修改一级目录取消
  updateBookCancel();
  bookLoading.value = true;
  BookApi.list()
    .then((res) => {
      books.value = res.data;
      emit("books", res.data);
    })
    .finally(() => {
      bookLoading.value = false;
    });
};

/**
 * 点击一级目录
 */
const bookClick = (book: Book) => {
  if (updateBookId.value) {
    return;
  }
  currentBookId.value = book.id;
};

/**
 * 点击添加一级目录保存
 */
const addBookSave = () => {
  let name = String(newBookName.value).trim();
  if (!name) {
    ElMessage.warning("请填写一级目录");
    return;
  }
  bookLoading.value = true;
  BookApi.add({ id: "", name: name })
    .then(() => {
      ElMessage.success("创建成功");
      queryBooks();
    })
    .catch(() => {
      bookLoading.value = false;
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
  BookApi.update({ id: updateBookId.value, name: name })
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
 * 点击删除一级目录
 */
const deleteBookClick = (book: Book) => {
  ElMessageBox.confirm("是否删除一级目录：" + book.name + "？", "提示", {
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
