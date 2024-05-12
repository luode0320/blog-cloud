export * from "./number";
export * from "./string";
export * from "./time";
export * from "./Upload";

/**
 * 复制对象
 * @param obj
 * @returns
 */
export const clone = <T>(obj: T): T => {
  return JSON.parse(JSON.stringify(obj));
};

/**
 * 判空
 * @param obj
 * @returns
 */
export const isEmpty = (obj: any) => {
  return obj === undefined || obj === null || obj === "" || obj.length === 0;
};

/**
 * 递归树
 * @param array
 * @param callback
 * @param childName
 * @param stop
 */
export const recursive = (array: any[], callback: Function, childName: string = "children", stop?: boolean) => {
  if (array && array.length > 0 && !stop) {
    array.forEach((item) => {
      stop = callback(item);
      recursive(item[childName], callback, childName, stop);
    });
  }
};

/**
 * 延时
 * @param ms 延时毫秒数
 * @returns
 */
export const timeout = (ms: number) => {
  return new Promise((resolve) => setTimeout(resolve, ms));
};

/**
 * 防抖函数
 * @param func 执行函数
 * @param ms 延迟毫秒数
 * @param immediate 是否立即执行一次
 * @returns
 */
export const debounce = (func: Function, ms: number, immediate?: boolean) => {
  let timer: NodeJS.Timeout | null = null;
  let debounced = function (context: any) {
    if (timer) {
      clearTimeout(timer);
    }
    if (immediate) {
      if (!timer) {
        func.apply(context, arguments);
      }
      timer = setTimeout(() => {
        func.apply(context, arguments);
        if (timer) {
          clearTimeout(timer);
          timer = null;
        }
      }, ms);
    } else {
      timer = setTimeout(() => {
        func.apply(context, arguments);
      }, ms);
    }
  };
  return debounced;
};
