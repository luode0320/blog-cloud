import axios, { AxiosRequestConfig } from "axios"; // 引入axios库，用于发送HTTP请求
import { ElMessage } from "element-plus"; // 引入Element Plus的弹框消息组件
import { host, context, basicTokenKey } from "@/config"; // 导入配置文件中的主机地址、上下文路径和基础令牌密钥
import Token from "@/store/token"; // 导入token存储模块，用于管理用户认证信息
import sha256 from "crypto-js/sha256"; // 引入sha256加密库，用于生成请求头中的Basic认证信息

// 定义一个泛型接口，用于规范后端返回的数据结构
interface ResponseResult<T> {
    code: number; // HTTP状态码，200表示成功，非200表示有错误
    message: string; // 错误信息或成功提示信息
    data: T; // 返回的数据实体，具体类型由泛型T决定
}

// 声明一个标志位，用于记录当前是否正处于刷新token的过程中，防止重复刷新
let refreshing = false;

// 创建一个用于数据请求的axios实例，配置基础URL和超时时间
const dataInstance = axios.create({
    baseURL: host + context + "/data",
    timeout: 20000,
});

// 创建一个用于授权请求的axios实例，配置基础URL和超时时间
const authInstance = axios.create({
    baseURL: host + context + "/token",
    timeout: 10000,
});


// 创建一个用于开放接口请求的axios实例，配置基础URL和超时时间
const openInstance = axios.create({
    baseURL: host + context + "/open",
    timeout: 20000,
});

// 数据接口的请求拦截器
dataInstance.interceptors.request.use(
    (config) => {
        // 如果请求头中没有Authorization字段，则添加Bearer Token认证信息
        if (!config.headers.Authorization) {
            config.headers.Authorization = "Bearer " + Token.getAccessToken();
        }
        return config; // 返回修改后的配置
    },
    (error) => {
        // 处理请求错误，直接返回reject
        return Promise.reject(error);
    }
);

// 授权接口的请求拦截器
authInstance.interceptors.request.use(
    (config) => {
        // 如果请求头中没有Authorization字段，生成基于时间戳和基本令牌密钥的哈希值作为Basic认证信息
        if (!config.headers.Authorization) {
            let token = sha256(basicTokenKey + Math.floor(new Date().getTime() / 600000)).toString();
            config.headers.Authorization = "Basic " + token;
        }
        return config;
    },
    (error) => {
        return Promise.reject(error);
    }
);

// 数据接口的响应拦截器
dataInstance.interceptors.response.use(
    (response) => {
        // 成功状态码为200且业务代码也为200时直接返回响应
        if (response.status === 200 && response.data.code === 200) {
            return response;
        }
        // 当接收到401响应码时，尝试刷新token
        else if (response.data.code === 401) {
            return new Promise((resolve, reject) => {
                // 避免并发刷新token
                if (!refreshing) {
                    refreshing = true;
                    // 发起刷新token的请求
                    authInstance.post("/refresh", { refreshToken: Token.getRefreshToken() })
                        .then((tokenResult) => {
                            refreshing = false;
                            // 刷新成功，更新token，并重新发起之前被拒绝的请求
                            Token.setToken(tokenResult.data.data);
                            response.config.headers.Authorization = "Bearer " + Token.getAccessToken();
                            dataInstance(response.config).then(resolve).catch(reject);
                        })
                        .catch(() => {
                            refreshing = false;
                            // 刷新失败，清除token信息
                            Token.removeToken();
                            reject();
                        });
                } else {
                    // 若已在刷新，等待刷新完成再重试
                    const interval = setInterval(() => {
                        if (!refreshing) {
                            clearInterval(interval);
                            response.config.headers.Authorization = "Bearer " + Token.getAccessToken();
                            dataInstance(response.config).then(resolve).catch(reject);
                        }
                    }, 500);
                }
            });
        }
        // 其他错误情况，显示错误信息
        else {
            ElMessage.error(response.data.message || "服务器错误");
            return Promise.reject(response);
        }
    },
    (error) => {
        // 网络错误或超时，显示错误信息
        ElMessage.error(error.statusText || "连接超时");
        return Promise.reject(error);
    }
);

/**
 * 授权接口respone拦截器
 */
authInstance.interceptors.response.use(
  (response) => {
    // 请求结果码
    if (response.status === 200) {
      // 业务结果码
      if (response.data.code === 200) {
        return response;
      }
      ElMessage.error(response.data.message ? response.data.message : "服务器错误");
    } else {
      ElMessage.error(response.statusText ? response.statusText : "连接超时");
    }
    return Promise.reject(response);
  },
  (error) => {
    ElMessage.error("连接超时");
    return Promise.reject(error);
  }
);

/**
 * 开放接口respone拦截器
 */
openInstance.interceptors.response.use(
  (response) => {
    // 请求结果码
    if (response.status === 200) {
      // 业务结果码
      if (response.data.code === 200) {
        return response;
      }
      ElMessage.error(response.data.message ? response.data.message : "服务器错误");
    } else {
      ElMessage.error(response.statusText ? response.statusText : "连接超时");
    }
    return Promise.reject(response);
  },
  (error) => {
    ElMessage.error("连接超时");
    return Promise.reject(error);
  }
);

// 封装对外提供的数据接口请求函数
export default function request<T>(config: AxiosRequestConfig): Promise<ResponseResult<T>> {
    return new Promise((resolve, reject) => {
        dataInstance.request<ResponseResult<T>>(config)
            .then((res) => resolve(res.data))
            .catch((err) => reject(err.data || err));
    });
}

// 封装对外提供的授权接口请求函数
export function authRequest<T>(config: AxiosRequestConfig): Promise<ResponseResult<T>> {
    return new Promise((resolve, reject) => {
        authInstance.request<ResponseResult<T>>(config)
            .then((res) => resolve(res.data))
            .catch((err) => reject(err.data || err));
    });
}


// 封装对外提供的开放接口请求函数
export function openRequest<T>(config: AxiosRequestConfig): Promise<ResponseResult<T>> {
    return new Promise((resolve, reject) => {
        openInstance.request<ResponseResult<T>>(config)
            .then((res) => resolve(res.data))
            .catch((err) => reject(err.data || err));
    });
}