import { read, utils } from "xlsx";

export class Excel {
  /**
   * 获取excel的sheet
   * @param data
   * @param sheetNum sheet序号，从1开始
   * @returns
   */
  static getExcelSheet(data: string | ArrayBuffer, sheetNum: number = 1) {
    const workbook = read(data, { type: "array" });
    const firstSheetName = workbook.SheetNames[sheetNum - 1];
    return workbook.Sheets[firstSheetName];
  }

  /**
   * 获取表格数据
   * @param data
   * @param headerRow 表头行号，从1开始
   * @param sheetNum sheet序号，从1开始
   * @returns
   */
  static getExcelJson<T>(data: string | ArrayBuffer, headerRow: number = 1, sheetNum: number = 1) {
    const sheet = Excel.getExcelSheet(data, sheetNum);
    return utils.sheet_to_json<T>(sheet, { range: headerRow - 1 });
  }

  /**
   * 获取表头行
   * @param data
   * @param headerRow 表头行号，从1开始
   * @param sheetNum sheet序号，从1开始
   * @returns
   */
  static getExcelHeaders(data: string | ArrayBuffer, headerRow: number = 1, sheetNum: number = 1) {
    const sheet = Excel.getExcelSheet(data, sheetNum);
    const headers = [];
    if (sheet["!ref"]) {
      const range = utils.decode_range(sheet["!ref"]);
      const row = range.s.r + headerRow - 1;
      for (let i = range.s.c; i <= range.e.c; ++i) {
        let cell = sheet[utils.encode_cell({ c: i, r: row })];
        let hdr = "UNKNOWN " + i;
        if (cell && cell.t) {
          hdr = utils.format_cell(cell);
        }
        headers.push(hdr);
      }
    }
    return headers;
  }

  /**
   * 格式化excel中的时间
   * @param date
   * @returns
   */
  static formatExcelDate(date: number) {
    const utcDays = Math.floor(date - 25569);
    const utcValue = utcDays * 86400;
    const dateInfo = new Date(utcValue * 1000);
    const fractionalDay = date - Math.floor(date) + 0.0000001;
    let totalSeconds = Math.floor(86400 * fractionalDay);
    const seconds = totalSeconds % 60;
    totalSeconds -= seconds;
    const hours = Math.floor(totalSeconds / (60 * 60));
    const minutes = Math.floor(totalSeconds / 60) % 60;
    const result = new Date(dateInfo.getFullYear(), dateInfo.getMonth(), dateInfo.getDate(), hours, minutes, seconds);
    if (result instanceof Date && !isNaN(result.getTime())) {
      return result;
    }
    return null;
  }
}
