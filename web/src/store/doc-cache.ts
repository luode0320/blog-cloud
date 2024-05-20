// 导入Token模块，该模块通常用于管理用户的认证信息，如获取用户名等
import Token from "./token";
// 导入localforage库，这是一个客户端存储库，支持IndexedDB, WebSQL, or localStorage存储，提供异步API和更友好的Promise接口
import localforage from "localforage";

// 使用localforage创建一个名为"doc"的独立存储实例，便于管理文档相关的缓存数据
const store = localforage.createInstance({
  name: "doc",
});

// 定义一个函数用于生成缓存的键名，这里结合用户名来确保每个用户的缓存数据相互独立
const getKey = () => {
  return "DocCache_" + Token.getName();
};

// 定义一个类DocCache，用于处理文档的缓存逻辑
class DocCache {
  /**
   * 缓存文档内容的方法。
   * 接受一个CurrentDoc类型的对象作为参数，将其存储到localforage中。
   * @param currentDoc 当前需要缓存的文档对象
   * @returns 存储操作的Promise对象
   */
  setDoc(currentDoc: CurrentDoc) {
    return store.setItem<CurrentDoc>(getKey(), currentDoc);
  }

  /**
   * 获取文档内容的方法。
   * 异步返回Promise，解析后得到缓存中的文档对象，如果没有缓存则返回一个默认的文档对象结构。
   * @returns 包含文档信息的Promise对象
   */
  async getDoc(): Promise<CurrentDoc> {
    return new Promise((resolve, reject) => {
      // 尝试从缓存中获取文档信息
      store.getItem<CurrentDoc>(getKey())
        .then((res) => {
          // 如果存在缓存，则解析并返回
          if (res) {
            resolve(res);
          } else {
            // 如果没有缓存内容，则返回一个默认构造的文档对象
            resolve({
              id: "",
              name: "",
              content: "",
              originMD5: "",
              type: "",
              updateTime: "",
            });
          }
        })
        .catch((err) => {
          // 如果在获取过程中出现错误，则捕获异常并返回默认文档对象
          console.error(err);
          resolve({
            id: "",
            name: "",
            content: "",
            originMD5: "",
            type: "",
            updateTime: "",
          });
        });
    });
  }

  /**
   * 清空文档缓存的方法。
   * 删除与当前用户关联的文档缓存项。
   * @returns 删除操作的Promise对象
   */
  removeDoc() {
    return store.removeItem(getKey());
  }
}

export default new DocCache();
