<template>
  <div v-loading="bookLoading" :class="{ 'page-book-shrink': !isStretch }" class="page-book">
    <div v-if="loading" class="mask-view"></div>

    <!--
    这段 Vue 代码是用于实现一个弹出框，用于创建新的一级目录。
    v-if 指令用于根据条件判断是否渲染元素，如果 !onlyPreview 为真，则渲染下面的内容。
    :visible 属性用于控制弹出框的显示与隐藏，这里使用 addBookVisible 变量来控制。
    placement 属性用于设置弹出框的位置，这里设置为底部。
    trigger 属性用于设置触发弹出框的方式，这里设置为点击触发。
    width 属性用于设置弹出框的宽度。

    el-input 组件用于渲染一个输入框。
    v-model 指令用于实现双向绑定，将 newBookName 与输入框的值进行绑定。
    el-button 组件用于渲染一个按钮。
    @click 事件监听器用于在点击取消按钮时触发 addBookCancel 方法。
    @click 事件监听器用于在点击确定按钮时触发 addBookSave 方法。

    template 标签用于定义一个模板，这里用于定义弹出框的触发按钮。
    el-button 组件用于渲染一个按钮。
    :icon 属性用于设置按钮的图标，这里使用 Plus 图标。
    class 属性用于设置按钮的类名。
    link 属性用于设置按钮的样式为链接样式。
    size 属性用于设置按钮的大小。
    @click 事件监听器用于在点击按钮时触发 addBookVisible = true，显示弹出框。
    -->
    <el-popover v-if="!onlyPreview" :visible="addBookVisible" placement="bottom" trigger="click" width="200px">

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

    <el-scrollbar class="scroll-view">
      <!--
      这段 Vue 代码是用于渲染一个列表，根据 books 数组的内容生成对应的项。
      v-for 指令用于遍历 books 数组，将数组中的每个元素赋值给 item。
      :key 指令用于为每个项设置唯一的标识，这里使用 item.id 作为标识。
      :class 指令用于根据条件动态设置元素的类名，如果 currentBookId 等于 item.id，则设置类名为 'selected'，否则为空字符串。
      @click 事件监听器用于在点击项时触发 bookClick 方法，并将 item 作为参数传递。

      v-if 指令用于根据条件判断是否渲染元素，如果 updateBookId 存在且等于 item.id，则渲染下面的内容，否则渲染下面的 text-tip 组件。
      v-model 指令用于实现双向绑定，将 updateBookName 与输入框的值进行绑定。
      @click 事件监听器用于在点击保存按钮时触发 updateBookSave 方法。
      @click 事件监听器用于在点击取消按钮时触发 updateBookCancel 方法。

      v-else 指令用于在上面的条件不满足时渲染下面的内容。

      v-if 指令用于根据条件判断是否渲染元素，如果 !onlyPreview 为真且 item.id 存在，则渲染下面的内容。
      @click.stop 修饰符用于阻止事件冒泡。
      @click 事件监听器用于在点击设置按钮时触发空函数。

      el-dropdown 组件用于实现下拉菜单的功能。
      el-icon 组件用于渲染一个图标。
      :title 属性用于设置图标的标题。
      @click.stop 修饰符用于阻止事件冒泡。
      el-dropdown-menu 组件用于渲染下拉菜单的内容。
      el-dropdown-item 组件用于渲染下拉菜单的项。
      -->
      <div v-for="(item, index) in books" :key="item.id"
           :class="[(currentBookId === '' && index === 0) || currentBookId === item.id ? 'selected' : '', 'item-view']"
           @click="bookClick(item,index)">

        <div v-if="updateBookId && updateBookId === item.id" class="update-view">
          <el-input v-model="updateBookName" placeholder="请输入一级目录"></el-input>
          <el-button :icon="CircleCheckFilled" link style="margin-left: 12px" type="success"
                     @click="updateBookSave"></el-button>
          <el-button :icon="CircleCloseFilled" link type="danger" @click="updateBookCancel"></el-button>
        </div>
        <text-tip v-else :content="item.name"></text-tip>

        <el-dropdown v-if="!onlyPreview && item.id" trigger="click">
          <el-icon class="setting-button" title="操作" @click.stop="() => {}">
            <Tools/>
          </el-icon>
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
import {ref, Ref, onMounted, watch} from "vue";
import {ElMessage, ElMessageBox} from "element-plus";
import {Plus, Tools, CircleCheckFilled, CircleCloseFilled} from "@element-plus/icons-vue";
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
  // 点击添加一级目录取消，执行 addBookCancel 方法
  addBookCancel();
  // 点击修改一级目录取消，执行 updateBookCancel 方法
  updateBookCancel();
  // 设置书籍加载状态为 true
  bookLoading.value = true;
  // 调用 BookApi.list() 方法获取书籍列表
  BookApi.list()
      .then((res) => {
        // 将获取到的书籍列表数据赋值给 books
        books.value = res.data;
        // 触发 books 事件，传递书籍列表数据
        emit("books", res.data);

        if (currentBookId.value === '' && books.value.length > 0) {
          currentBookId.value = books.value[0].id;
        }
      })
      .finally(() => {
        // 设置书籍加载状态为 false
        bookLoading.value = false;
      });
};


/**
 * 点击一级目录
 */
const bookClick = (book: Book, index: number) => {
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
  BookApi.add({id: "", name: name})
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
 * 点击添加二级目录
 */
const addTwoBookClick = (book: Book) => {
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
  BookApi.update({id: updateBookId.value, name: name})
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
