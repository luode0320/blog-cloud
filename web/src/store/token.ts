interface Token {
  nameKey: string;
  accessTokenKey: string;
  refreshTokenKey: string;
}

class Token {
  constructor() {
    this.nameKey = "Name";
    this.accessTokenKey = "AccessToken";
    this.refreshTokenKey = "RefreshToken";
  }

  /**
   * 缓存token
   * @param token
   */
  setToken(token: TokenResult) {
    if (token) {
      localStorage.setItem(this.nameKey, token.name);
      localStorage.setItem(this.accessTokenKey, token.accessToken);
      localStorage.setItem(this.refreshTokenKey, token.refreshToken);
    }
  }

  /**
   * 清空token
   */
  removeToken() {
    localStorage.removeItem(this.accessTokenKey);
    localStorage.removeItem(this.refreshTokenKey);
    location.reload();
  }

  /**
   * 获取用户名
   */
  getName() {
    return localStorage.getItem(this.nameKey);
  }

  /**
   * 获取AccessToken
   */
  getAccessToken() {
    return localStorage.getItem(this.accessTokenKey);
  }

  /**
   * 获取RefreshToken
   */
  getRefreshToken() {
    return localStorage.getItem(this.refreshTokenKey);
  }
}

export default new Token();
