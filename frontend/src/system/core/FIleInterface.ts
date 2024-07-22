import type { OsFile, OsFileWithoutContent } from './FileSystem';

export interface OsFileInterface {
  readFile: (path: string) => Promise<string | null>;
  writeFile: (
    path: string,
    data: string | ArrayBuffer | Blob,
    opt?: {
      flag?: 'w' | 'a' | 'wx';
    }
  ) => Promise<void>;
  appendFile: (path: string, content: string) => Promise<void>;
  readdir: (path: string) => Promise<OsFileWithoutContent[]>;
  exists: (path: string) => Promise<boolean>;
  stat: (path: string) => Promise<OsFileWithoutContent | null>;
  unlink: (path: string) => Promise<void>;
  rename: (oldPath: string, newPath: string) => Promise<void>;
  rmdir: (path: string) => Promise<void>;
  mkdir: (path: string) => Promise<void>;
  copyFile: (src: string, dest: string) => Promise<void>;
  chmod: (path: string, mode: number) => Promise<void>;
  search: (keyword: string) => Promise<OsFileWithoutContent[]>;
  serializeFileSystem: () => Promise<unknown>;
  deserializeFileSystem: (files: OsFile[]) => Promise<unknown>;
  removeFileSystem: () => Promise<void>;
  registerWatcher: (path: RegExp, callback: (path: string, content: string) => void) => void;
}
