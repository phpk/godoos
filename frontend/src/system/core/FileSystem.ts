import * as fspath from './Path';
import { OsFileInterface } from './FIleInterface';
import { SystemOptions } from '../type/type';
import { InitSystemFile, InitUserFile } from './SystemFileConfig';
import { createInitFile } from './createInitFile';
import { OsFileMode } from './FileMode';
type DateLike = Date | string | number;
// Os文件模式枚举
class OsFileInfo {
  // 是否是文件
  isFile = true;
  // 是否是目录
  isDirectory = false;
  // 是否是符号链接
  isSymlink = false;
  // 文件大小
  size = 0;

  // 最后一次修改此文件的时间戳
  mtime: DateLike = new Date();
  // 最后一次访问此文件的时间戳
  atime: DateLike = new Date();
  // 此文件创建时间的时间戳
  birthtime: DateLike = new Date();
  // 文件权限
  mode = 0o777;
  // 设备编号
  rdev = 0;

  // 构造函数
  constructor(
    // 是否是文件
    isFile?: boolean,
    // 是否是目录
    isDirectory?: boolean,
    // 是否是符号链接
    isSymlink?: boolean,
    // 文件大小
    size?: number,
    // 最后一次修改此文件的时间戳
    mtime?: DateLike,
    // 最后一次访问此文件的时间戳
    atime?: DateLike,
    // 此文件创建时间的时间戳
    birthtime?: DateLike,
    // 文件权限
    mode?: OsFileMode,
    // 设备编号
    rdev?: number
  ) {
    // 设置是否是文件
    if (isFile !== undefined) {
      this.isFile = isFile;
    }
    // 设置是否是目录
    if (isDirectory !== undefined) {
      this.isDirectory = isDirectory;
    }
    // 设置是否是符号链接
    if (isSymlink !== undefined) {
      this.isSymlink = isSymlink;
    }
    // 设置文件大小
    if (size !== undefined) {
      this.size = size;
    }
    // 设置最后一次修改此文件的时间戳
    if (mtime !== undefined) {
      this.mtime = mtime;
    }
    // 设置最后一次访问此文件的时间戳
    if (atime !== undefined) {
      this.atime = atime;
    }
    // 设置此文件创建时间的时间戳
    if (birthtime !== undefined) {
      this.birthtime = birthtime;
    }
    // 设置文件权限
    if (mode !== undefined) {
      this.mode = mode;
    }
    // 设置设备编号
    if (rdev !== undefined) {
      this.rdev = rdev;
    }
  }
}

/**
 * OsFile 类是继承自 OsFileInfo 的一个类。
 */
class OsFile extends OsFileInfo {
  name = ''; // 文件名
  path: string; // 文件路径
  oldPath?: string; // 旧的文件路径
  parentPath: string; // 父目录路径
  content: any; // 文件内容
  ext?: string; // 文件扩展名
  title?: string; // 文件名（不包含扩展名）
  id?: number; // 文件ID（可选）
  isSys?: number; // 文件是否是系统文件（可选）

  /**
   * OsFile 类的构造函数。
   * @param path 文件路径
   * @param content 文件内容
   * @param info 文件信息
   * @param id 文件ID（可选）
   */
  constructor(
    path: string,
    content: any,
    info: Partial<OsFileInfo>,
    id?: number
  ) {
    if (info.isFile) {
      info.isDirectory = false;
      info.isSymlink = false;
    }
    if (info.isDirectory) {
      info.isFile = false;
      info.isSymlink = false;
    }
    if (info.isSymlink) {
      info.isFile = false;
      info.isDirectory = false;
    }
    super(info.isFile, info.isDirectory, info.isSymlink, info.size, info.mtime, info.atime, info.birthtime);
    this.path = path; // 设置文件路径
    this.parentPath = fspath.dirname(path); // 获取文件所在目录路径
    this.name = fspath.basename(path); // 获取文件名
    if (!this.oldPath) {
      this.oldPath = this.path; // 设置旧的文件路径
    }
    if (info.isFile) {
      const titleArr = this.name.split('.');
      this.ext = titleArr.pop(); // 获取文件扩展名
      this.title = titleArr.join("."); // 获取文件名（不包含扩展名）
    }
    this.content = content; // 设置文件内容
    // this.icon = icon; // 文件图标
    // this.type = type; // 文件类型

    this.id = id; // 设置文件ID

    if (id === undefined) {
      delete this.id; // 如果文件ID为空，则删除该属性
    }
    if (this.isSys === undefined) {
      this.isSys = 1; // 如果文件系统属性为空，则设置为1（系统文件）
    }
  }
}
export type OsFileWithoutContent = Omit<OsFile, 'content'>;
/**
 * Os文件系统类
 * 
 * @implements {OsFileInterface}
 */
class OsFileSystem implements OsFileInterface {
  /**
  * 私有的数据库实例
  * @type {IDBDatabase}
  */
  private db!: IDBDatabase;

  /**
   * 当数据库准备就绪时的回调函数
   * @type {((value: OsFileSystem) => void) | null}
   */
  private _ready: ((value: OsFileSystem) => void) | null = null;

  /**
   * 监听路径和内容变化的回调函数映射表
   * @type {Map<RegExp, (path: string, content: string) => void>}
   */
  private _watchMap: Map<RegExp, (path: string, content: string) => void> = new Map();

  /**
   * 文件卷映射表
   * @type {Map<string, OsFileInterface>}
   */
  private volumeMap: Map<string, OsFileInterface> = new Map();
  /**
  * 当打开数据库失败时触发的错误处理函数
  * @param e - 错误对象
  */
  onerror: (e: any) => void = () => {
    console.error('Failed to open database');
  };
  /**
 * 构造函数
 * @param {string} rootPath - 根路径，默认为'/'
 * @param {string} id - id，默认为空字符串
 */
  constructor(rootPath = '/', id = '') {
    // 打开IndexedDB数据库
    const request = window.indexedDB.open('FileSystemDB' + id, 1);
    // 打开失败时的错误处理
    request.onerror = () => {
      console.error('Failed to open database');
    };

    // 打开成功时的处理
    request.onsuccess = () => {
      this.db = request.result;
      this._ready?.(this);
    };

    // 当数据库升级时的处理
    request.onupgradeneeded = () => {
      this.db = request.result;
      // 创建对象存储空间'files'
      const objectStore = this.db.createObjectStore('files', { keyPath: 'id', autoIncrement: true });
      // 创建索引'parentPath'
      objectStore.createIndex('parentPath', 'parentPath');
      // 创建索引'path'，设置唯一性
      objectStore.createIndex('path', 'path', { unique: true });
      // 创建索引'name'
      objectStore.createIndex('name', 'name');

      // 创建根目录
      const rootDir = new OsFile(rootPath, '', {
        mode: 0o111,
        isDirectory: true,
      });
      rootDir.parentPath = rootPath === '/' ? '' : fspath.dirname(rootPath);

      // 将根目录添加到对象存储空间
      objectStore.add(rootDir);
    };
  }

  /**
   * 初始化文件系统
   * @param option 系统选项
   * @returns 返回文件系统实例
   */
  async initFileSystem(option: SystemOptions) {
    await this.whenReady();
    await this.mkdir('/C');
    await this.chmod('/C', OsFileMode.Read);
    await createInitFile(this, option.initFile || InitUserFile, option.userLocation);
    await createInitFile(this, option.initFile || InitSystemFile, option.systemLocation);

    await this.mkdir('/D');
    await this.chmod('/D', OsFileMode.Read);
    await this.mkdir('/E');
    await this.chmod('/E', OsFileMode.Read);

    await this.mkdir('/B');
    await this.chmod('/B', OsFileMode.Read);

    return this;
  }
  on(_: 'error', func: (e: any) => void) {
    this.onerror = func;
  }
  /**
  * 序列化文件系统
  * @returns {Promise<any>} 返回一个Promise对象，包含文件系统的数据
  */
  serializeFileSystem() {
    return new Promise((resolve, reject) => {
      // 创建一个只读事务
      const transaction = this.db.transaction('files', 'readonly');
      // 获取文件对象存储
      const objectStore = transaction.objectStore('files');
      // 获取所有文件数据的请求
      const request = objectStore.getAll();
      // 请求失败时的处理函数
      request.onerror = () => {
        reject('Failed to read file');
      };
      // 请求成功时的处理函数
      request.onsuccess = () => {
        resolve(request.result);
      };
    });
  }
  /**
  * 反序列化文件系统
  * @param files 要添加到文件系统的文件数组
  * @returns 返回一个Promise对象，成功时无返回值，失败时返回错误信息
  */
  deserializeFileSystem(files: OsFile[]) {
    return new Promise((resolve, reject) => {
      const transaction = this.db.transaction('files', 'readwrite');
      const objectStore = transaction.objectStore('files');
      const request = objectStore.clear();
      request.onerror = () => {
        reject('Failed to clear file');
      };
      request.onsuccess = () => {
        files.forEach((file) => {
          objectStore.add(file);
        });
        resolve(void 0);
      };
    });
  }
  /**
  * 当准备就绪时调用的函数，返回一个Promise对象，该对象在OsFileSystem实例准备就绪后会被resolve。
  * 如果OsFileSystem实例已经准备就绪，则直接返回该实例。
  * 
  * @returns {Promise<OsFileSystem>} 返回一个Promise对象，该对象在OsFileSystem实例准备就绪后会被resolve。
  */
  whenReady(): Promise<OsFileSystem> {
    if (this.db) {
      return Promise.resolve(this);
    }
    return new Promise<OsFileSystem>((resolve) => {
      this._ready = resolve;
    });
  }
  /**
 * 注册观察者
 * @param path - 观察的路径，使用正则表达式表示
 * @param callback - 当路径发生变化时调用的回调函数，接收路径和内容作为参数
 */
  registerWatcher(path: RegExp, callback: (path: string, content: string) => void) {
    this._watchMap.set(path, callback);
  }
  /**
 * 提交文件变更观察
 * @param {string} path - 文件路径
 * @param {string} content - 文件内容
 */
  commitWatch(path: string, content: any) {
    // 遍历观察回调函数
    this._watchMap.forEach((callback, reg) => {
      // 如果文件路径匹配正则表达式
      if (reg.test(path)) {
        // 调用回调函数，传入文件路径和内容
        callback(path, content);
      }
    });
  }
  /**
  * 删除文件系统数据库
  * @returns {Promise} 删除完成的Promise对象
  */
  async removeFileSystem() {
    window.indexedDB.deleteDatabase('FileSystemDB');
    return Promise.resolve();
  }

  /**
 * 挂载卷到指定路径
 * @param path - 路径
 * @param volume - 卷对象
 */
  mountVolume(path: string, volume: OsFileInterface) {
    this.volumeMap.set(path, volume);
  }

  /**
  * 检查给定路径下的卷文件
  * @param path 要检查的路径
  * @returns 卷文件或undefined
  */
  checkVolumeChild(path: string): OsFileInterface | undefined {
    let volume: OsFileInterface | undefined;
    this.volumeMap.forEach((volumem, key) => {
      if (fspath.dirname(key) === path) {
        volume = volumem;
      }
    });
    return volume;
  }
  /**
   * 判断指定路径是否为卷的路径
   * @param path
   * @returns
   */
  checkVolumePath(path: string): OsFileInterface | undefined {
    if (this.volumeMap.has(path)) {
      return this.volumeMap.get(path);
    }
    let volume: OsFileInterface | undefined;
    this.volumeMap.forEach((volumem, key) => {
      if (fspath.isChildPath(key, path)) {
        volume = volumem;
      }
    });
    return volume;
  }

  /**
   * 使用对应的卷的文件系统
   */
  beforeGuard<T extends keyof OsFileInterface>(
    volume: OsFileInterface,
    opt: T,
    ...args: Parameters<OsFileInterface[T]>
  ) {
    return (volume[opt] as (...args: Parameters<OsFileInterface[T]>) => ReturnType<OsFileInterface[T]>)(
      ...args
    );
  }

  /**
   * 读取指定路径的文件内容
   * @param path 文件路径
   * @returns 文件内容
   */
  async readFile(path: string): Promise<string | null> {
    try {
      const volume = this.checkVolumePath(path);
      if (volume) {
        return this.beforeGuard(volume, 'readFile', path);
      }

      const transaction = this.db.transaction('files', 'readonly');
      const objectStore = transaction.objectStore('files');

      const index = objectStore.index('path');
      const range = IDBKeyRange.only(path);
      const request = index.get(range);

      return new Promise((resolve, reject) => {
        request.onerror = () => {
          reject('Failed to read file');
        };
        request.onsuccess = () => {
          const file: OsFile = request.result;
          resolve(file ? file.content : null);
        };
      });
    } catch (e: any) {
      this.onerror(e.toString());
      return Promise.reject(e);
    }
  }


  /**
 * 写入文件 不存在则创建，存在则覆盖
 * @param path 文件路径
 * @param data 文件内容
 * @param opt 可选参数
 * @param opt.flag 写入模式，可选值为 'w' (覆盖写入)，'a' (追加写入)，'wx' (排他写入)
 * @returns Promise
 */
  async writeFile(
    path: string,
    data: any,
    opt?: {
      flag?: 'w' | 'a' | 'wx';
    }
  ): Promise<void> {
    const volume = this.checkVolumePath(path);
    if (volume) {
      return this.beforeGuard(volume, 'writeFile', path, data, opt);
    }

    const parentPath = fspath.dirname(path);
    // 判断文件是否存在
    const exists = await this.exists(parentPath);
    if (!exists) {
      this.onerror('无法写入一个不存在的路径上的文件:' + path);
      return Promise.reject('无法写入一个不存在的路径上的文件:' + path);
    }

    const transaction = this.db.transaction('files', 'readwrite');
    const objectStore = transaction.objectStore('files');

    const stats: OsFile | null = await new Promise((resolve, reject) => {
      objectStore.index('path').openCursor(IDBKeyRange.only(path)).onsuccess = (event: any) => {
        const cursor: IDBCursorWithValue = event.target.result;
        if (cursor) {
          const file: OsFile = cursor.value;
          if (file.isDirectory) {
            reject('无法写入一个目录');
          } else {
            resolve(file);
          }
        } else {
          resolve(null);
        }
      };
    });

    if (!stats) {
      const request = objectStore.add(
        new OsFile(path, data, {
          isFile: true,
          size: (typeof data === 'string' ? data.length : data.byteLength),
        })
      );
      return new Promise((resolve, reject) => {
        request.onerror = () => {
          this.onerror('写入文件失败');
          reject('写入文件失败');
        };
        request.onsuccess = () => {
          this.commitWatch(path, data);
          resolve();
        };
      });
    } else {
      if (opt?.flag === 'wx') {
        // 排他模式
        return Promise.resolve();
      }
      if (opt?.flag === 'a') {
        // 追加模式
        data = stats.content + data;
      }

      const request = objectStore.put({
        ...stats,
        content: data,
        size: (typeof data === 'string' ? data.length : data.byteLength),
        mtime: new Date(),
      });
      return new Promise((resolve, reject) => {
        request.onerror = () => {
          this.onerror('写入文件失败');
          reject('写入文件失败');
        };
        request.onsuccess = () => {
          this.commitWatch(path, data);
          resolve();
        };
      });
    }
  }
  /**
 * 在指定路径的文件末尾追加内容。
 * @param path 文件路径
 * @param content 要追加的内容
 * @returns 返回一个Promise，表示操作的异步结果
 */
  async appendFile(path: string, content: string): Promise<void> {
    const volume = this.checkVolumePath(path);
    if (volume) {
      return this.beforeGuard(volume, 'appendFile', path, content);
    }

    const transaction = this.db.transaction('files', 'readwrite');
    const objectStore = transaction.objectStore('files');

    const index = objectStore.index('path');
    const range = IDBKeyRange.only(path);
    const request = index.get(range);

    return new Promise((resolve, reject) => {
      request.onerror = () => {
        this.onerror('Failed to write file');
        reject('Failed to read file');
      };
      request.onsuccess = () => {
        const file: OsFile = request.result;
        if (file) {
          file.content += content;
          file.size = file.content.length;
          file.mtime = new Date();
          const request = objectStore.put(file);
          request.onerror = () => {
            this.onerror('Failed to write file');
            reject('Failed to write file');
          };
          request.onsuccess = () => {
            this.commitWatch(path, file.content);
            resolve();
          };
        } else {
          this.onerror('File not found');
          reject('File not found');
        }
      };
    });
  }
  /**
   * 读取指定路径下的所有文件和文件夹
   * @param fpath 目录路径
   * @returns 文件和文件夹列表
   */
  async readdir(fpath: string): Promise<OsFileWithoutContent[]> {
    const path = fspath.resolve(fpath);
    const volume = this.checkVolumePath(path);
    if (volume) {
      return this.beforeGuard(volume, 'readdir', path);
    }

    const volume2 = this.checkVolumeChild(path);
    let vol: OsFileWithoutContent[] = [];
    if (volume2) {
      try {
        vol = await this.beforeGuard(volume2, 'readdir', path);
      } catch {
        this.onerror('Failed to read volume directory:' + path);
      }
    }

    const transaction = this.db.transaction('files', 'readonly');
    const objectStore = transaction.objectStore('files');

    const index = objectStore.index('parentPath');
    const range = IDBKeyRange.only(path);
    //console.log(range)
    const request = index.getAll(range);
    //console.log(request)
    return new Promise((resolve, reject) => {
      request.onerror = () => {
        this.onerror('Failed to read directory');
        reject('Failed to read directory');
      };
      request.onsuccess = () => {
        const files = request.result;
        //console.log(files)
        resolve([...files, ...vol]);
      };
    });
  }

  async exists(path: string): Promise<boolean> {
    const volume = this.checkVolumePath(path);
    if (volume) {
      try {
        return this.beforeGuard(volume, 'exists', path);
      } catch {
        this.onerror('Failed to read volume directory:' + path);
      }
    }

    try {
      const transaction = this.db.transaction('files', 'readonly');
      const objectStore = transaction.objectStore('files');

      const index = objectStore.index('path');
      const range = IDBKeyRange.only(path);
      const request = index.getAll(range);

      return new Promise((resolve, reject) => {
        request.onerror = () => {
          this.onerror('Failed to read file');
          reject('Failed to read file');
        };
        request.onsuccess = () => {
          const fileArray: OsFile[] = request.result;
          resolve(fileArray.length ? true : false);
        };
      });
    } catch (e) {
      return false;
    }
  }

  /**
 * 获取文件信息
 * @param path 文件路径
 * @returns 返回文件信息的Promise
 */
  async stat(path: string): Promise<OsFileWithoutContent | null> {
    const volume = this.checkVolumePath(path);
    if (volume) {
      try {
        return this.beforeGuard(volume, 'stat', path);
      } catch {
        this.onerror('Failed to read volume directory:' + path);
      }
    }

    const transaction = this.db.transaction('files', 'readonly');
    const objectStore = transaction.objectStore('files');

    const index = objectStore.index('path');
    const range = IDBKeyRange.only(path);
    const request = index.get(range);

    return new Promise((resolve, reject) => {
      request.onerror = () => {
        this.onerror('Failed to read file');
        reject('Failed to read file');
      };
      request.onsuccess = () => {
        const file: OsFile = request.result;
        resolve(file);
      };
    });
  }

  /**
   * 删除指定路径的文件
   * @param path 文件路径
   */
  async unlink(path: string): Promise<void> {
    const volume = this.checkVolumePath(path);
    if (volume) {
      try {
        return this.beforeGuard(volume, 'unlink', path);
      } catch {
        this.onerror('Failed to unlink volume file:' + path);
      }
    }

    const transaction = this.db.transaction('files', 'readwrite');
    const objectStore = transaction.objectStore('files');

    const index = objectStore.index('path');
    const range = IDBKeyRange.only(path);
    const request = index.get(range);

    return new Promise((resolve, reject) => {
      request.onerror = () => {
        this.onerror('Failed to delete file');
        reject('Failed to delete file');
      };
      request.onsuccess = () => {
        const file: OsFile = request.result;
        if (file) {
          if (file.isDirectory) {
            reject('Cannot delete a directory');
          } else {
            objectStore.delete(request.result.id);
            this.commitWatch(path, file.content);
            resolve();
          }
        } else {
          reject('File not found');
        }
      };
    });
  }
  /**
 * 深度遍历文件系统并重命名文件或目录
 * @param vfile - 要重命名的文件或目录对象
 * @param objectStore - IDBObjectStore 对象存储
 * @param newPath - 新路径
 */
  private async dfsRename(vfile: OsFile, objectStore: IDBObjectStore, newPath: string) {
    if (vfile.isDirectory) {
      // 对于目录，遍历其所有父目录的文件和子目录
      objectStore.index('parentPath').openCursor(IDBKeyRange.only(vfile.path)).onsuccess = (event: any) => {
        const cursor: IDBCursorWithValue = event.target.result;
        if (cursor) {
          const tempfile = cursor.value;
          const tempNewPath = fspath.join(newPath, fspath.basename(tempfile.path));
          // 递归调用 dfsRename 函数重命名文件或目录
          this.dfsRename(tempfile, objectStore, tempNewPath);
          cursor.continue();
        }
      };
    }
    const vParentPath = fspath.dirname(newPath);
    vfile.oldPath = vfile.path;
    vfile.path = newPath;
    vfile.parentPath = vParentPath;
    vfile.mtime = new Date();

    // 更新文件或目录信息到 objectStore
    objectStore.put(vfile);
  }
  /**
 * 重命名文件或目录
 * @param path - 要重命名的文件或目录的路径
 * @param newPath - 新的文件或目录路径
 * @returns 返回一个Promise，当重命名成功时，Promise将解析为undefined
 * @throws 当重命名失败时，Promise将被拒绝并抛出错误
 */
  async rename(path: string, newPath: string): Promise<void> {
    const volume = this.checkVolumePath(path);
    const volume2 = this.checkVolumePath(newPath);
    if (!!volume && !!volume2) {
      return this.beforeGuard(volume, 'rename', path, newPath);
    } else if ((!!volume && !volume2) || (!volume && !!volume2)) {
      this.onerror('Cannot rename between volumes');
      return Promise.reject('Cannot rename between volumes');
    }

    // this.beforeGuard('rename', path, newPath);
    // 不能重命名为子路径
    // /C/Users /C/Users/Desktop/Users
    if (path === newPath) {
      return Promise.resolve();
    }
    if (fspath.isChildPath(path, newPath)) {
      this.onerror('Cannot rename to child path');
      return Promise.reject('Cannot rename to child path');
    }
    // if (newPath.startsWith(path)) {//bug
    //   return Promise.reject('Cannot rename to child path');
    // }

    const transaction = this.db.transaction('files', 'readwrite');
    const objectStore = transaction.objectStore('files');

    const index = objectStore.index('path');
    const range = IDBKeyRange.only(path);
    const request = index.get(range);

    return new Promise((resolve, reject) => {
      request.onerror = () => {
        this.onerror('Failed to read file');
        reject('Failed to read file');
      };
      request.onsuccess = () => {
        const file: OsFile = request.result;
        if (file) {
          this.dfsRename(file, objectStore, newPath);
          this.commitWatch(path, file.content);
        }

        resolve();
      };
    });
  }
  /**
 * 递归删除指定目录及其子目录
 * @param vfile 要删除的目录文件对象
 * @param objectStore 对象存储
 * @returns 返回一个Promise对象，表示删除操作是否成功
 */
  private async dfsRmdir(vfile: OsFile, objectStore: IDBObjectStore) {
    if (vfile.mode) {
      if (vfile.mode <= 0o111) {
        this.onerror('Cannot delete a readonly file');
        return Promise.reject('Cannot delete a readonly file');
      }
    }
    if (vfile.isDirectory) {
      // 获取指定目录的所有父目录
      objectStore.index('parentPath').openCursor(IDBKeyRange.only(vfile.path)).onsuccess = (event: any) => {
        const cursor: IDBCursorWithValue = event.target.result;
        if (cursor) {
          const tempfile = cursor.value;
          // 递归删除父目录及其子目录
          this.dfsRmdir(tempfile, objectStore);
          cursor.continue();
        }
      };
    }
    // 删除指定目录的所有文件
    objectStore.index('path').openCursor(IDBKeyRange.only(vfile.path)).onsuccess = (event: any) => {
      const cursor: IDBCursorWithValue = event.target.result;
      if (cursor) {
        objectStore.delete(cursor.value.id);
        cursor.continue();
      }
    };
    // 提交对指定目录的监视
    this.commitWatch(vfile.path, vfile.content);
  }
  /**
   * 删除指定路径的文件夹及其内容
   * @param path 文件夹路径
   */
  async rmdir(path: string): Promise<void> {
    const volume = this.checkVolumePath(path);
    if (volume) {
      return this.beforeGuard(volume, 'rmdir', path);
    }

    const transaction = this.db.transaction('files', 'readwrite');
    const objectStore = transaction.objectStore('files');

    const index = objectStore.index('path');
    const range = IDBKeyRange.only(path);
    const request = index.get(range);

    return new Promise((resolve, reject) => {
      request.onerror = () => {
        this.onerror('Failed to read file');
        reject('Failed to read file');
      };
      request.onsuccess = () => {
        const file: OsFile = request.result;
        if (file) {
          this.dfsRmdir(file, objectStore);
        }
        this.commitWatch(path, file.content);
        resolve();
      };
    });
  }

  /**
   * 创建新的文件夹
   * @param path 文件夹路径
   */
  async mkdir(path: string): Promise<void> {
    const volume = this.checkVolumePath(path);
    if (volume) {
      return this.beforeGuard(volume, 'mkdir', path);
    }

    // 转换路径
    const transedPath = fspath.transformPath(path);
    // 获取父路径
    let parentPath = fspath.dirname(transedPath);
    // 如果父路径为空，设置为根路径
    if (parentPath === '') parentPath = '/';
    // 判断文件是否存在
    const exists = await this.exists(parentPath);
    // 如果文件不存在，则抛出错误
    if (!exists) {
      this.onerror('Cannot create directory to a non-exist path:' + parentPath);
      return Promise.reject('Cannot create directory to a non-exist path:' + parentPath);
    }

    // 检查目录是否已存在
    const res = await this.exists(transedPath);
    if (res) {
      // 目录已存在，无需操作
      // console.error("Directory already exists");
      return Promise.resolve();
    }

    // 开始数据库事务
    const transaction = this.db.transaction('files', 'readwrite');
    // 获取文件对象存储
    const objectStore = transaction.objectStore('files');

    // 添加新目录到文件对象存储
    const request = objectStore.add(
      new OsFile(transedPath, '', {
        isDirectory: true,
      })
    );

    return new Promise((resolve, reject) => {
      request.onerror = () => {
        this.onerror('Failed to create directory');
        reject('Failed to create directory');
      };
      request.onsuccess = () => {
        this.commitWatch(transedPath, '');
        resolve();
      };
    });
  }

  /**
 * 深度复制文件
 * @param vfile - 虚拟文件对象
 * @param objectStore - IDBObjectStore对象
 * @param newPath - 新路径
 */
  private async dfsCopFile(vfile: OsFile, objectStore: IDBObjectStore, newPath: string) {
    if (vfile.isDirectory) {
      // 获取所有父目录的文件
      objectStore.index('parentPath').openCursor(IDBKeyRange.only(vfile.path)).onsuccess = (event: any) => {
        const cursor: IDBCursorWithValue = event.target.result;
        if (cursor) {
          const tempfile = cursor.value;
          const tempNewPath = fspath.join(newPath, fspath.basename(tempfile.path));
          // 递归复制父目录下的所有文件
          this.dfsCopFile(tempfile, objectStore, tempNewPath);
          cursor.continue();
        }
      };
    }
    // 复制文件
    const newFile = {
      ...vfile,
      path: newPath,
      parentPath: fspath.dirname(newPath),
      mtime: new Date(),
    };
    delete newFile.id;
    objectStore.put(newFile);
  }
  /**
  * 复制文件
  * @param src 文件源路径
  * @param dest 文件目标路径
  * @returns 返回一个Promise
  */
  async copyFile(src: string, dest: string): Promise<void> {
    const volume = this.checkVolumePath(src);
    const volume2 = this.checkVolumePath(dest);

    if (!!volume && !!volume2) {
      return this.beforeGuard(volume, 'copyFile', src, dest);
    } else if ((!!volume && !volume2) || (!volume && !!volume2)) {
      /**
       * 无法在卷之间复制文件
       */
      this.onerror('Cannot copyFile between volumes');
      return Promise.reject('Cannot copyFile between volumes');
    }

    const transaction = this.db.transaction('files', 'readwrite');
    const objectStore = transaction.objectStore('files');

    const index = objectStore.index('path');
    const range = IDBKeyRange.only(src);
    const request = index.get(range);

    return new Promise((resolve, reject) => {
      request.onerror = () => {
        /**
         * 读取文件失败
         */
        this.onerror('Failed to read file');
        reject('Failed to read file');
      };
      request.onsuccess = () => {
        const file: OsFile = request.result;
        if (file) {
          this.dfsCopFile(file, objectStore, dest);
          this.commitWatch(src, file.content);
        }

        resolve();
      };
    });
  }

  /**
 * 设置指定文件或目录的权限
 * @param path 文件或目录的路径
 * @param mode 文件或目录的权限模式
 * @returns 返回一个Promise，成功时无返回值，失败时抛出错误
 */
  async chmod(path: string, mode: OsFileMode): Promise<void> {
    const respath = fspath.resolve(path); // 解析路径
    const volume = this.checkVolumePath(respath); // 检查路径是否属于某个卷
    if (volume) {
      return this.beforeGuard(volume, 'chmod', respath, mode); // 如果属于卷路径，则调用beforeGuard方法进行鉴权
    }

    const transaction = this.db.transaction('files', 'readwrite'); // 创建数据库事务
    const objectStore = transaction.objectStore('files'); // 获取文件对象存储

    const index = objectStore.index('path'); // 获取路径索引
    const range = IDBKeyRange.only(respath); // 创建键值范围
    const request = index.get(range); // 获取指定范围内的数据

    return new Promise((resolve, reject) => {
      request.onerror = () => {
        this.onerror('Failed to read file'); // 请求错误时调用onerror方法
        reject('Failed to read file'); // 抛出错误
      };
      request.onsuccess = () => {
        const file: OsFile = request.result; // 获取请求结果
        if (file) {
          file.mode = mode; // 设置文件权限
          file.mtime = new Date(); // 设置文件的修改时间
          objectStore.put(file); // 更新数据库中的文件记录
        } else {
          reject('File not found'); // 文件不存在时抛出错误
        }
        resolve(); // 请求成功
      };
    });
  }
  /**
 * 根据关键字搜索文件
 * @param keyword 搜索关键字
 * @returns 返回一个Promise，Promise的结果是一个OsFileWithoutContent类型的数组，表示搜索到的文件列表
 */
  async search(keyword: string): Promise<OsFileWithoutContent[]> {
    const transaction = this.db.transaction('files', 'readonly');
    const objectStore = transaction.objectStore('files');

    const index = objectStore.index('name');
    const range = IDBKeyRange.bound(keyword, keyword + '\uffff');
    const request = index.getAll(range);

    return new Promise((resolve, reject) => {
      request.onerror = () => {
        this.onerror('搜索文件失败');
        reject('搜索文件失败');
      };
      request.onsuccess = () => {
        const files = request.result;
        resolve(files);
      };
    });
  }
}

export { OsFile, OsFileInfo, OsFileSystem };
