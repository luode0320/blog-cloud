import request from "@/utils/request";
import sha256 from "crypto-js/sha256";

class UserApi {
  /**
   * 修改密码
   * @param password
   * @param newPassword
   * @returns
   */
  updatePassword(password: string, newPassword: string) {
    return request({
      method: "post",
      url: "/user/update-password",
      data: { password: sha256(password).toString(), newPassword: sha256(newPassword).toString() },
    });
  }
}

export default new UserApi();
