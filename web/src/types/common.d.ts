interface Page {
  current: number;
  size: number;
}

interface PageCondition<T> {
  page: Page;
  condition: T;
}

interface PageResult<T> {
  records: Array<T>;
  total: number;
}
