import { openRequest } from "@/utils/request";

class OpenApi {
  /**
   * 查询公开发布文档
   * @param id
   * @returns
   */
  getDoc(id: string) {
    return openRequest<Doc>({
      method: "get",
      url: "/doc/get/" + id,
    });
  }

  /**
   * 分页查询公开发布文档列表
   * @param pageCondition
   * @returns
   */
  pageDoc(pageCondition: PageCondition<DocPageCondition>) {
    return openRequest<PageResult<DocPageResult>>({
      method: "post",
      url: "/doc/page",
      data: pageCondition,
    });
  }
}

export default new OpenApi();
