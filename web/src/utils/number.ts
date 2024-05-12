/**
 * 判断非数字
 * @param value
 * @returns
 */
export const checkNaN = (value: any): boolean => {
  return value === undefined || value === null || isNaN(value) || value === Infinity || value === -Infinity;
};

/**
 * 添加千分位
 * @param value
 */
export const comdify = (value: any) => {
  if (checkNaN(value)) {
    return "";
  }
  let n = value + "";
  let re = /\d{1,3}(?=(\d{3})+$)/g;
  let n1 = n.replace(/^(\d+)((\.\d+)?)$/, function (s, s1, s2) {
    return s1.replace(re, "$&,") + s2;
  });
  return n1;
};

/**
 * 保留n位小数
 * @param value 待格式化数字
 * @param digit 小数位数
 * @param noComdify 不加千分位
 */
export const formatFloat = (value: any, digit: number = 2, noComdify?: boolean) => {
  if (checkNaN(value)) {
    return "";
  }
  const multiple = Math.pow(10, digit);
  let result = Math.round(parseFloat(value) * multiple) / multiple + "";
  let xsd = result.split(".");
  if (xsd.length === 1) {
    // 没有小数的情况
    let completion = "";
    if (digit > 0) {
      completion += ".";
      for (let i = 0; i < digit; i++) {
        completion += "0";
      }
    }
    result = result + completion;
  } else if (xsd.length > 1) {
    // 有小数的情况
    for (let i = 0; i < digit - xsd[1].length; i++) {
      result = result + "0";
    }
  }
  return noComdify ? result : comdify(result);
};
