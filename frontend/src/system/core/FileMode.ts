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