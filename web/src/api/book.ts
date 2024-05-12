import request from "@/utils/request";

class BookApi {
  /**
   * 查询文集列表
   * @returns
   */
  list() {
    return request<Book[]>({
      method: "post",
      url: "/book/list",
    });
  }

  /**
   * 添加文集
   * @param book
   * @returns
   */
  add(book: Book) {
    return request({
      method: "post",
      url: "/book/add",
      data: book,
    });
  }

  /**
   * 修改文集
   * @param book
   * @returns
   */
  update(book: Book) {
    return request({
      method: "post",
      url: "/book/update",
      data: book,
    });
  }

  /**
   * 根据id删除文集
   * @param id
   * @returns
   */
  delete(id: string) {
    return request({
      method: "post",
      url: "/book/delete",
      data: { id: id },
    });
  }
}

export default new BookApi();
