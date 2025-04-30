type DateLike = Date | string | number;
export interface File {
  id?: number; // 文件ID（可选）
  size?: number;
  created: Date;
  modified: Date;
  name : string; // 文件名
  path: string; // 文件路径
  knowledgeId?:number;//知识库id
  oldPath?: string; // 旧的文件路径
  parentPath: string; // 父目录路径
  content: any; // 文件内容
  ext?: string; // 文件扩展名
  title: string; // 文件名（不包含扩展名）
  isSys?: boolean; // 文件是否是系统文件（可选）
  isShare?: boolean; // 文件是否为共享文件
  isPwd?: boolean; //文件是否上锁
  isFile?: boolean; // 文件是否为文件
  // 是否是目录
  isDirectory?: boolean;
  // 是否是符号链接
  isSymlink?: boolean;
   // 最后一次修改此文件的时间戳
   mtime?: DateLike;
   // 最后一次访问此文件的时间戳
   atime?: DateLike;
   // 此文件创建时间的时间戳
   birthtime?: DateLike;
   // 文件权限
   mode?: OsFileMode;
}

export enum OsFileMode {
  // 读取权限
  Read = 0b001, //1
  // 写入权限
  Write = 0b010, //2
  // 执行权限
  Execute = 0b100, //4
  // 读取和写入权限
  ReadWrite = Read | Write, //3
  // 读取和执行权限
  ReadExecute = Read | Execute, //5
  // 写入和执行权限
  WriteExecute = Write | Execute, // 6
  // 读取、写入和执行权限
  ReadWriteExecute = Read | Write | Execute, //7
}