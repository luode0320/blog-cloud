<template>
  <div class="page-picture">
    <el-table class="table-view" ref="tableRef" :data="tableData" height="100%" stripe border v-loading="tableLoading">
      <el-table-column prop="name" label="序号" align="center" width="100">
        <template #default="scope">
          {{ (tableCondition.page.current - 1) * tableCondition.page.size + scope.$index + 1 }}
        </template>
      </el-table-column>
      <el-table-column prop="name" label="图片名称" align="center" />
      <el-table-column prop="path" label="图片地址" align="center">
        <template #default="scope">
          {{ hostUrl + scope.row.picturePrefix + scope.row.path }}
        </template>
      </el-table-column>
      <el-table-column prop="size" label="图片大小" align="center">
        <template #default="scope"> {{ Upload.formatFileSize(scope.row.size) }} </template>
      </el-table-column>
      <el-table-column prop="createTime" label="上传时间" align="center">
        <template #default="scope"> {{ formatTime(scope.row.createTime, "YYYY-MM-DD HH:mm:ss") }} </template>
      </el-table-column>
      <el-table-column prop="path" label="缩略图" align="center" width="150px">
        <template #default="scope">
          <el-image
            class="table-picture"
            :src="hostUrl + scope.row.thumbnailPrefix + scope.row.path"
            :preview-src-list="[hostUrl + scope.row.picturePrefix + scope.row.path]"
            fit="contain"
            preview-teleported
            hide-on-click-modal
          />
        </template>
      </el-table-column>
      <el-table-column prop="path" label="" align="center" width="150px">
        <template #header>
          <el-button type="primary" @click="uploadClick">上传</el-button>
        </template>
        <template #default="scope">
          <el-button type="danger" @click="deleteClick(scope.row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-pagination
      background
      layout="total, sizes, prev, pager, next, jumper"
      :pageSizes="[10, 20, 50, 100]"
      v-model:pageSize="tableCondition.page.size"
      v-model:currentPage="tableCondition.page.current"
      :total="tableTotal"
      @size-change="tablePageSizeChange"
      @current-change="tablePageCurrentChange"
    ></el-pagination>
  </div>
</template>

<script lang="ts" setup>
import { ref, Ref, onMounted, nextTick } from "vue";
import { ElMessage, ElMessageBox, ElTable } from "element-plus";
import PictureApi from "@/api/picture";
import { uploadPicture } from "./util";
import { formatTime, Upload } from "@/utils";
import { host } from "@/config";

const hostUrl = ref("");
const tableData: Ref<PicturePageResult[]> = ref([]);
const tableTotal = ref(0);
const tableLoading = ref(false);
const tableCondition: Ref<PageCondition<null>> = ref({
  page: {
    current: 1,
    size: 20,
  },
  condition: null,
});
const tableRef = ref<InstanceType<typeof ElTable>>();

onMounted(() => {
  hostUrl.value = process.env.NODE_ENV === "production" ? location.origin : host;
  queryTableData();
});

/**
 * 查询表格数据
 */
const queryTableData = () => {
  tableLoading.value = true;
  PictureApi.page(tableCondition.value)
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
 * 删除记录
 */
const deleteClick = (row: PicturePageResult) => {
  ElMessageBox.confirm("是否删除图片：" + row.name + "？", "提示", {
    confirmButtonText: "确定",
    cancelButtonText: "取消",
    type: "warning",
  }).then(() => {
    tableLoading.value = true;
    PictureApi.delete(row.id)
      .then(() => {
        ElMessage.success("删除成功");
        tablePageCurrentChange(1);
      })
      .catch(() => {
        tableLoading.value = false;
      });
  });
};

/**
 * 上传图片
 */
const uploadClick = () => {
  Upload.openFiles(false, Upload.InputAccept.uploadImage).then((fileList) => {
    tableLoading.value = true;
    uploadPicture(fileList[0])
      .then(() => {
        tablePageCurrentChange(1);
      })
      .catch(() => {
        tableLoading.value = false;
      });
  });
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
</script>

<style lang="scss">
.page-picture {
  display: flex;
  flex-direction: column;
  align-items: center;
  .table-view {
    .el-scrollbar__wrap {
      display: flex;
    }
    .table-picture {
      width: 100px;
      height: 100px;
    }
  }
}
</style>
