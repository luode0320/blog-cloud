import dayjs, { ManipulateType } from "dayjs";

/**
 * 格式化时间
 * @param date
 * @param pattern
 * @returns
 */
export const formatTime = (date: dayjs.ConfigType, pattern: string = "YYYY-MM-DD HH:mm:ss"): string => {
  return dayjs(date).format(pattern);
};

/**
 * 转为时间格式
 * @param date
 * @returns
 */
export const parseTime = (date: string | number) => {
  return dayjs(date).toDate();
};

/**
 * 增加指定时间
 * @param date
 * @param value
 * @param unit
 * @returns
 */
export const addTime = (date: Date, value: number, unit: ManipulateType) => {
  return dayjs(date).add(value, unit).toDate();
};

/**
 * 获取时间是周几
 * @param date
 * @returns
 */
export const getWeek = (date: Date) => {
  const weekNum = dayjs(date).day();
  var week = ["日", "一", "二", "三", "四", "五", "六"];
  return "周" + week[weekNum];
};
