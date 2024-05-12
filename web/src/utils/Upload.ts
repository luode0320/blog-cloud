interface HTMLInputEvent extends Event {
  target: HTMLInputElement & EventTarget;
}

/**
 * 文件上传
 */
export class Upload {
  /**
   * 上传文件类型
   */
  static InputAccept = {
    xlsx: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
    image: "image/*",
    uploadImage: ".apng,.bmp,.gif,.ico,.jfif,.jpeg,.jpg,.png,.webp",
  };

  /**
   * 上传文件
   * @param multiple 多选
   * @param accept 支持类型
   * @returns
   */
  static openFiles(multiple?: boolean, accept?: string) {
    return new Promise<FileList>((resolve) => {
      const input = document.createElement("input");
      input.type = "file";
      if (accept) {
        input.accept = accept;
      }
      if (multiple) {
        input.multiple = true;
      }
      input.onchange = (e: Event) => {
        let event = e as HTMLInputEvent;
        const elem = event.target;
        resolve(elem.files ?? new FileList());
      };
      input.click();
    });
  }

  /**
   * 获取文件名和后缀
   * @param file
   * @returns
   */
  static getFileName(file: File) {
    let splits = file.name.split(".");
    let name = splits.slice(0, splits.length - 1).join(".");
    let ext = "";
    if (splits.length > 1) {
      ext = splits[splits.length - 1].toLowerCase();
    }
    return { name, ext };
  }

  /**
   * 将Blob读取为ArrayBuffer
   * @param blob
   * @returns
   */
  static readBlobToArrayBuffer(blob: Blob) {
    return new Promise<ArrayBuffer>((resolve, reject) => {
      const reader = new FileReader();
      reader.onload = (e) => {
        const data = e.target?.result;
        if (data) {
          resolve(data as ArrayBuffer);
        } else {
          reject();
        }
      };
      reader.onerror = () => {
        reject();
      };
      reader.readAsArrayBuffer(blob);
    });
  }

  /**
   * 将Blob读取为String
   * @param blob
   * @returns
   */
  static readBlobToString(blob: Blob) {
    return new Promise<string>((resolve, reject) => {
      const reader = new FileReader();
      reader.onload = (e) => {
        const text = (e.target?.result as string) || "";
        resolve(text);
      };
      reader.onerror = () => {
        reject();
      };
      reader.readAsText(blob);
    });
  }

  /**
   * 将Blob读取为Base64
   * @param blob
   * @returns
   */
  static readBlobToBase64(blob: Blob) {
    return new Promise<string>((resolve, reject) => {
      const reader = new FileReader();
      reader.onload = (e) => {
        const text = (e.target?.result as string) || "";
        resolve(text);
      };
      reader.onerror = () => {
        reject();
      };
      reader.readAsDataURL(blob);
    });
  }

  /**
   * 格式化文件大小
   * @param size
   */
  static formatFileSize(size: number) {
    if (size < 1000) {
      return size + "B";
    }
    if (size < 1000 * 1000) {
      return (size / 1000).toFixed(2) + "KB";
    }
    if (size < 1000 * 1000 * 1000) {
      return (size / 1000 / 1000).toFixed(2) + "MB";
    }
    return (size / 1000 / 1000 / 1000).toFixed(2) + "GB";
  }
}
