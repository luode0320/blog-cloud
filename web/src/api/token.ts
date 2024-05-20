// 导入自定义封装的带有鉴权处理的请求函数
import { authRequest } from "@/utils/request";
// 导入SHA256加密库用于密码加密
import sha256 from "crypto-js/sha256";
// 导入Token管理类，用于获取刷新令牌
import Token from "@/store/token";

/**
 * TokenApi类，封装了与用户认证相关的API请求方法。
 */
class TokenApi {
  /**
   * 用户注册方法，向服务器发送注册请求。
   *
   * @param name - 用户名
   * @param password - 明文密码，将被加密后发送
   * @returns 注册操作的Promise，包含服务器响应数据
   */
  signUp(name: string, password: string): Promise<any> {
    // 使用SHA256加密用户密码
    const encryptedPassword = sha256(password).toString();
    // 发送POST请求到服务器的注册接口，包含加密后的密码
    return authRequest({
      method: "post",
      url: "/sign-up",
      data: { name, password: encryptedPassword },
    });
  }

  /**
   * 用户登录方法，向服务器发送登录请求。
   *
   * @param name - 用户名
   * @param password - 明文密码，将被加密后发送
   * @returns 登录操作的Promise，解析后为TokenResult类型的数据
   */
  signIn(name: string, password: string): Promise<any> {
    // 加密用户密码
    const encryptedPassword = sha256(password).toString();
    // 发送登录请求，包含加密后的密码
    return authRequest<TokenResult>({
      method: "post",
      url: "/sign-in",
      data: { name, password: encryptedPassword },
    });
  }

  /**
   * 用户登出方法，向服务器发送登出请求以注销当前会话。
   *
   * @returns 登出操作的Promise
   */
  signOut(): Promise<any> {
    // 获取本地存储的刷新令牌
    const refreshToken = Token.getRefreshToken();
    // 如果有刷新令牌，则发送登出请求
    if (refreshToken) {
      return authRequest({
        method: "post",
        url: "/sign-out",
        data: { refreshToken },
      });
    } else {
      // 若无刷新令牌，理论上不应到达此分支，但可适当处理或抛出错误
      console.error("找不到用于注销的刷新令牌.");
      return Promise.reject("没有可用的刷新令牌.");
    }
  }
}

// 导出TokenApi类的实例，使得外部可以使用此类提供的方法而无需创建新实例
export default new TokenApi();