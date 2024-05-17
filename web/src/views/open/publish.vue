<template>
  <div class="page-open-publish">
    <div class="title-view">公开文档</div>
    <el-table ref="tableRef" v-loading="tableLoading" :data="tableData" border class="table-view" height="100%" stripe>
      <el-table-column align="center" label="序号" prop="name" width="100">
        <template #default="scope">
          {{ (tableCondition.page.current - 1) * tableCondition.page.size + scope.$index + 1 }}
        </template>
      </el-table-column>
      <el-table-column :label-class-name="columnClass.name" align="center" label="文档名称" prop="name">
        <template #header="scope">
          <el-popover v-model:visible="namePopover" :hide-after="0" trigger="click" width="170" @hide="tablePopoverHide">
            <el-input v-model="tableCondition.condition.name" clearable placeholder="文档名称筛选" @clear="namePopover = false"></el-input>
            <template #reference>
              <div style="cursor: pointer">文档名称</div>
            </template>
          </el-popover>
        </template>
      </el-table-column>
      <el-table-column :label-class-name="columnClass.type" align="center" label="文档类型" prop="type">
        <template #header="scope">
          <el-popover v-model:visible="typePopover" :hide-after="0" trigger="click" width="170" @hide="tablePopoverHide">
            <el-select v-model="tableCondition.condition.type" clearable placeholder="文档类型筛选" @clear="typePopover = false">
              <el-option label="Markdown" value="md" />
              <el-option label="OpenAPI" value="openApi" />
            </el-select>
            <template #reference>
              <div style="cursor: pointer">文档类型</div>
            </template>
          </el-popover>
        </template>
        <template #default="scope"> {{ scope.row.type === "md" ? "Markdown" : "OpenAPI" }} </template>
      </el-table-column>
      <el-table-column :label-class-name="columnClass.bookName" align="center" label="一级目录" prop="bookName">
        <template #header="scope">
          <el-popover v-model:visible="bookNamePopover" :hide-after="0" trigger="click" width="170" @hide="tablePopoverHide">
            <el-input v-model="tableCondition.condition.bookName" clearable placeholder="一级目录筛选" @clear="bookNamePopover = false"></el-input>
            <template #reference>
              <div style="cursor: pointer">一级目录</div>
            </template>
          </el-popover>
        </template>
      </el-table-column>
      <el-table-column :label-class-name="columnClass.username" align="center" label="作者" prop="username">
        <template #header="scope">
          <el-popover v-model:visible="usernamePopover" :hide-after="0" trigger="click" width="170" @hide="tablePopoverHide">
            <el-input v-model="tableCondition.condition.username" clearable placeholder="作者筛选" @clear="usernamePopover = false"></el-input>
            <template #reference>
              <div style="cursor: pointer">作者</div>
            </template>
          </el-popover>
        </template>
      </el-table-column>
      <el-table-column align="center" label="创建时间" prop="createTime">
        <template #default="scope"> {{ formatTime(scope.row.createTime, "YYYY-MM-DD HH:mm:ss") }} </template>
      </el-table-column>
      <el-table-column align="center" label="修改时间" prop="createTime">
        <template #default="scope"> {{ formatTime(scope.row.updateTime, "YYYY-MM-DD HH:mm:ss") }} </template>
      </el-table-column>
      <el-table-column align="center" label="文档地址" width="160">
        <template #default="scope">
          <el-button text type="primary" @click="copyClick(scope.row.id)">复制</el-button>
          <el-button text type="primary" @click="hrefClick(scope.row.id)">跳转</el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-pagination
      v-model:currentPage="tableCondition.page.current"
      v-model:pageSize="tableCondition.page.size"
      :pageSizes="[10, 20, 50, 100]"
      :total="tableTotal"
      background
      layout="total, sizes, prev, pager, next, jumper"
      @size-change="tablePageSizeChange"
      @current-change="tablePageCurrentChange"
    ></el-pagination>
  </div>
</template>

<script lang="ts" setup>
import { ref, Ref, nextTick, onMounted } from "vue";
import { ElTable, ElMessage } from "element-plus";
import OpenApi from "@/api/open";
import { formatTime } from "@/utils";
import copy from "copy-to-clipboard";

const hostUrl = ref(location.origin);
const tableCondition = ref({
  page: { current: 1, size: 100 },
  condition: { name: "", type: "", bookName: "", username: "" },
});
const columnClass = ref({
  name: "",
  type: "",
  bookName: "",
  username: "",
});
const lastCondition = ref("");
const tableData: Ref<DocPageResult[]> = ref([]);
const tableTotal = ref(0);
const tableLoading = ref(false);
const tableRef = ref<InstanceType<typeof ElTable>>();
const namePopover = ref(false);
const typePopover = ref(false);
const bookNamePopover = ref(false);
const usernamePopover = ref(false);

onMounted(() => {
  queryTableData();
});

/**
 * 查询表格数据
 */
const queryTableData = () => {
  lastCondition.value = JSON.stringify(tableCondition.value.condition);
  tableLoading.value = true;
  OpenApi.pageDoc(tableCondition.value)
    .then((res) => {
      tableData.value = res.data.records;
      tableTotal.value = res.data.total;
      nextTick(() => {
        tableRef.value?.setScrollTop(0);
      });
    })
    .finally(() => {
      tableLoading.value = false;
    });
};

/**
 * 表格自定义筛选
 */
const tablePopoverHide = () => {
  if (tableCondition.value.condition.name) {
    columnClass.value.name = "column-active";
  } else {
    columnClass.value.name = "";
  }
  if (tableCondition.value.condition.type) {
    columnClass.value.type = "column-active";
  } else {
    columnClass.value.type = "";
  }
  if (tableCondition.value.condition.bookName) {
    columnClass.value.bookName = "column-active";
  } else {
    columnClass.value.bookName = "";
  }
  if (tableCondition.value.condition.username) {
    columnClass.value.username = "column-active";
  } else {
    columnClass.value.username = "";
  }
  if (JSON.stringify(tableCondition.value.condition) === lastCondition.value) {
    return;
  }
  tableCondition.value.page.current = 1;
  queryTableData();
};

/**
 * 每页显示条数变化
 * @param size
 */
const tablePageSizeChange = (size: number) => {
  tableCondition.value.page.current = 1;
  tableCondition.value.page.size = size;
  queryTableData();
};

/**
 * 当前页码变化
 * @param current
 */
const tablePageCurrentChange = (current: number) => {
  tableCondition.value.page.current = current;
  queryTableData();
};

/**
 * 点击复制
 * @param id
 */
const copyClick = (id: string) => {
  let url = hostUrl.value + "/#/open/document?id=" + id;
  const result = copy(url);
  if (result) {
    ElMessage.success("发布地址已复制到剪切板");
  } else {
    ElMessage.error("复制到剪切板失败，地址：" + url);
  }
};

/**
 * 点击跳转
 * @param id
 */
const hrefClick = (id: string) => {
  let url = hostUrl.value + "/#/open/document?id=" + id;
  window.open(url, "_blank");
};
</script>

<style lang="scss">
.page-open-publish {
  height: 100%;
  width: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  .title-view {
    margin: 10px 0;
    font-size: 18px;
    color: #5f5f5f;
    font-weight: bold;
    user-select: none;
  }
  .el-button.is-text + .el-button.is-text {
    margin-left: 0;
  }
  .column-active {
    color: #0094c1;
  }
}
</style>
