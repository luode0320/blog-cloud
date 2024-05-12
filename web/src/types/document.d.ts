interface Doc {
  id: string;
  name: string;
  content: string;
  type?: string;
  published?: boolean;
  bookId: string;
  createTime?: number;
  updateTime?: number;
}

interface CurrentDoc {
  id: string;
  name: string;
  content: string;
  originMD5: string;
  type: string;
  updateTime: string;
}

interface DocPageResult {
  id: string;
  name: string;
  type: string;
  username: string;
  bookName: string;
  createTime: number;
  updateTime: number;
}

interface DocPageCondition {
  name?: string;
  type?: string;
  bookName?: string;
  username?: string;
}
