import { Upload } from "@/utils";
import Compressor from "compressorjs";
import PictureApi from "@/api/picture";
import { ElMessage } from "element-plus";

/**
 * 上传图片
 * @param file
 * @returns
 */
export const uploadPicture = (file: File): Promise<string> => {
  return new Promise((resolve, reject) => {
    const fileName = Upload.getFileName(file);
    const extArr = ["apng", "bmp", "gif", "ico", "jfif", "jpeg", "jpg", "png", "webp"];
    if (extArr.indexOf(fileName.ext) < 0) {
      ElMessage.warning("仅支持以下格式的图片：APNG、BMP、GIF、ICO、JPEG、PNG、WebP");
      reject();
      return;
    }
    if (file.size > 1000 * 1000 * 20) {
      ElMessage.warning("图片大小不可超过20MB");
      reject();
      return;
    }
    // 压缩图片
    new Compressor(file, {
      quality: 0.8,
      maxWidth: 100,
      maxHeight: 100,
      success(result) {
        // 上传图片
        const formData = new FormData();
        formData.append("picture", file);
        formData.append("thumbnail", new File([result], file.name));
        PictureApi.upload(formData)
          .then((res) => {
            ElMessage.success(res.message);
            resolve(res.data);
          })
          .catch(() => {
            reject();
          });
      },
      error(err) {
        ElMessage.warning("图片压缩失败：" + err.message);
        reject();
      },
    });
  });
};
