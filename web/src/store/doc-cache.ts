import Token from "./token";
import localforage from "localforage";

const store = localforage.createInstance({
  name: "doc",
});

const getKey = () => {
  return "DocCache_" + Token.getName();
};

class DocCache {
  /**
   * 缓存文档内容
   * @param currentDoc
   * @returns
   */
  setDoc(currentDoc: CurrentDoc) {
    return store.setItem<CurrentDoc>(getKey(), JSON.parse(JSON.stringify(currentDoc)));
  }

  /**
   * 获取文档内容
   * @returns
   */
  async getDoc(): Promise<CurrentDoc> {
    return new Promise((resolve, reject) => {
      store
        .getItem<CurrentDoc>(getKey())
        .then((res) => {
          if (res) {
            resolve(res);
          } else {
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
   * 清空文档缓存
   * @returns
   */
  removeDoc() {
    return store.removeItem(getKey());
  }
}

export default new DocCache();
