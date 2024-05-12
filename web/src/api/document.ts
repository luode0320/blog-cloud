import request from "@/utils/request";

class DocumentApi {
  /**
   * 查询文档列表
   * @param bookId
   * @returns
   */
  list(bookId: string) {
    return request<Doc[]>({
      method: "post",
      url: "/doc/list",
      data: { bookId: bookId },
    });
  }

  /**
   * 查询文档
   * @param id
   * @returns
   */
  get(id: string) {
    return request<Doc>({
      method: "post",
      url: "/doc/get",
      data: { id: id },
    });
  }

  /**
   * 添加文档
   * @param doc
   * @returns
   */
  add(doc: Doc) {
    return request<Doc>({
      method: "post",
      url: "/doc/add",
      data: doc,
    });
  }

  /**
   * 修改文档基础信息
   * @param doc
   * @returns
   */
  update(doc: Doc) {
    return request({
      method: "post",
      url: "/doc/update",
      data: doc,
    });
  }

  /**
   * 修改文档内容
   * @param doc
   * @returns
   */
  updateContent(doc: Doc) {
    return request<Doc>({
      method: "post",
      url: "/doc/update-content",
      data: doc,
    });
  }

  /**
   * 根据id删除文档
   * @param id
   * @returns
   */
  delete(id: string) {
    return request({
      method: "post",
      url: "/doc/delete",
      data: { id: id },
    });
  }
}

export default new DocumentApi();
