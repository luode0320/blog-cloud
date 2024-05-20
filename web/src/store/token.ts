/**
 * 定义一个接口来规范Token对象的结构，包含三个字符串类型的属性分别代表：
 * - 用户名
 * - 访问令牌（AccessToken）
 * - 刷新令牌（RefreshToken）
 */
interface TokenResult {
  name: string;
  accessToken: string;
  refreshToken: string;
}

/**
 * Token管理类，负责在浏览器的localStorage中存储和获取用户的认证Token信息。
 */
class Token {
  /**
   * 定义存储Token信息时使用的键名。
   */
  private nameKey: string; // 用户名的存储键
  private accessTokenKey: string; // 访问令牌的存储键
  private refreshTokenKey: string; // 刷新令牌的存储键

  /**
   * 构造函数，初始化存储键名。
   */
  constructor() {
    this.nameKey = "Name"; // 用户名的键
    this.accessTokenKey = "AccessToken"; // 访问令牌的键
    this.refreshTokenKey = "RefreshToken"; // 刷新令牌的键
  }

  /**
   * 将Token信息存储到localStorage中。
   *
   * @param token - 待存储的Token信息对象，需包含name、accessToken、refreshToken属性。
   */
  setToken(token: TokenResult) {
    if (token) {
      // 分别存储用户名、访问令牌、刷新令牌
      localStorage.setItem(this.nameKey, token.name);
      localStorage.setItem(this.accessTokenKey, token.accessToken);
      localStorage.setItem(this.refreshTokenKey, token.refreshToken);
    }
  }

  /**
   * 从localStorage中移除访问令牌和刷新令牌，并刷新当前页面以确保安全地结束用户会话。
   */
  removeToken() {
    localStorage.removeItem(this.accessTokenKey); // 移除访问令牌
    localStorage.removeItem(this.refreshTokenKey); // 移除刷新令牌
    location.reload(); // 刷新页面
  }

  /**
   * 从localStorage中获取用户名。
   *
   * @returns 返回存储的用户名字符串，若未找到则返回null。
   */
  getName() {
    return localStorage.getItem(this.nameKey);
  }

  /**
   * 从localStorage中获取访问令牌（AccessToken）。
   *
   * @returns 返回存储的访问令牌字符串，若未找到则返回null。
   */
  getAccessToken() {
    return localStorage.getItem(this.accessTokenKey);
  }

  /**
   * 从localStorage中获取刷新令牌（RefreshToken）。
   *
   * @returns 返回存储的刷新令牌字符串，若未找到则返回null。
   */
  getRefreshToken() {
    return localStorage.getItem(this.refreshTokenKey);
  }
}

// 导出Token类的实例，以便在整个应用中共享同一个实例进行Token管理。
export default new Token();