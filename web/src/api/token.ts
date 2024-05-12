import { authRequest } from "@/utils/request";
import sha256 from "crypto-js/sha256";
import Token from "@/store/token";

class TokenApi {
  /**
   * 注册
   * @param name
   * @param password
   * @returns
   */
  signUp(name: string, password: string) {
    return authRequest({
      method: "post",
      url: "/sign-up",
      data: { name: name, password: sha256(password).toString() },
    });
  }

  /**
   * 登录
   * @param name
   * @param password
   * @returns
   */
  signIn(name: string, password: string) {
    return authRequest<TokenResult>({
      method: "post",
      url: "/sign-in",
      data: { name: name, password: sha256(password).toString() },
    });
  }

  /**
   * 退出登录
   * @returns
   */
  signOut() {
    return authRequest({
      method: "post",
      url: "/sign-out",
      data: { refreshToken: Token.getRefreshToken() },
    });
  }
}

export default new TokenApi();
