import { ElMessage } from "element-plus"; // 导入Element Plus的消息提示组件
import Compressor from "compressorjs"; // 导入Compressor.js库，用于图片压缩

import PictureApi from "@/api/picture"; // 导入图片上传的API服务
import { Upload } from "@/utils"; // 导入自定义的上传工具模块

/**
 * 上传图片的异步函数，支持文件大小检查、格式验证及图片压缩后上传至服务器
 * @param {File} file - 用户选择的文件对象，预期为图片格式
 * @returns {Promise<string>} - 成功上传后的图片URL或在失败情况下抛出异常
 */
export const uploadPicture = (file: File): Promise<string> => {
  return new Promise((resolve, reject) => {
    // 从文件名中提取扩展名
    const fileName = Upload.getFileName(file);
    // 支持的图片扩展名数组
    const extArr = ["apng", "bmp", "gif", "ico", "jfif", "jpeg", "jpg", "png", "webp"];

    // 检查文件扩展名是否在允许的范围内
    if (extArr.indexOf(fileName.ext) < 0) {
      // 如果不是支持的图片格式，显示警告消息并拒绝承诺
      ElMessage.warning("仅支持以下格式的图片：APNG、BMP、GIF、ICO、JPEG、PNG、WebP");
      reject();
      return;
    }

    // 检查文件大小是否超过20MB
    if (file.size > 1000 * 1000 * 20) {
      ElMessage.warning("图片大小不可超过20MB");
      reject();
      return;
    }

    // 使用Compressor进行图片压缩
    new Compressor(file, {
      quality: 0.8, // 压缩质量，范围0-1
      maxWidth: 100, // 最大宽度，单位px
      maxHeight: 100, // 最大高度，单位px
      success(result) {
        // 压缩成功后，构造表单数据准备上传
        const formData = new FormData();
        formData.append("picture", file); // 原始图片
        formData.append("thumbnail", new File([result], file.name)); // 压缩后的图片作为缩略图

        // 调用图片上传API
        PictureApi.upload(formData)
          .then((res) => {
            // // 上传成功，显示成功消息并返回图片URL
            ElMessage.success(res.message);
            resolve(res.data);
          })
          .catch(() => {
            // 上传失败，拒绝
            reject();
          });
      },
      error(err) {
        // 压缩失败，显示错误消息并拒绝承诺
        ElMessage.warning("图片压缩失败：" + err.message);
        reject();
      },
    });
  });
};
