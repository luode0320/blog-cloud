import request from "@/utils/request";

class RsaApi {
  /**
   * 生成RSA密钥对
   * @param condition
   * @returns
   */
  generate(condition: RsaCondition) {
    return request<RsaResult>({
      method: "post",
      url: "/rsa/generate",
      data: condition,
    });
  }

  /**
   * RSA公钥加密
   * @param condition
   * @returns
   */
  encrypt(condition: RsaCondition) {
    return request<string>({
      method: "post",
      url: "/rsa/encrypt",
      data: condition,
    });
  }

  /**
   * RSA私钥解密
   * @param condition
   * @returns
   */
  decrypt(condition: RsaCondition) {
    return request<string>({
      method: "post",
      url: "/rsa/decrypt",
      data: condition,
    });
  }

  /**
   * RSA私钥签名
   * @param condition
   * @returns
   */
  sign(condition: RsaCondition) {
    return request<string>({
      method: "post",
      url: "/rsa/sign",
      data: condition,
    });
  }

  /**
   * RSA公钥验签
   * @param condition
   * @returns
   */
  verify(condition: RsaCondition) {
    return request<boolean>({
      method: "post",
      url: "/rsa/verify",
      data: condition,
    });
  }
}

export default new RsaApi();
