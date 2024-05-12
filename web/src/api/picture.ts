import request from "@/utils/request";

class PictureApi {
  /**
   * 分页查询
   * @param tableCondition
   * @returns
   */
  page(tableCondition: PageCondition<null>) {
    return request<PageResult<PicturePageResult>>({
      method: "post",
      url: "/pic/page",
      data: tableCondition,
    });
  }

  /**
   * 根据id删除图片
   * @param id
   * @returns
   */
  delete(id: string) {
    return request({
      method: "post",
      url: "/pic/delete",
      data: { id: id },
    });
  }

  /**
   * 上传图片
   * @param formData
   * @returns
   */
  upload(formData: FormData) {
    return request<string>({
      method: "post",
      url: "/pic/upload",
      data: formData,
      timeout: 60000,
    });
  }
}

export default new PictureApi();
