import { version } from '../../../package.json';
import { OsFileMode } from './FileMode';
//console.log(OsFileMode)
const InitSystemFile = {
  type: 'dir',
  name: '',
  mode: OsFileMode.ReadWriteExecute,
  children: [
    {
      name: 'plugin',
      type: 'dir',
      mode: OsFileMode.ReadWriteExecute,
      children: [],
    },
    {
      name: 'os',
      type: 'dir',
      mode: OsFileMode.ReadWriteExecute,
      children: [
        {
          name: 'version.txt',
          type: 'file',
          content: version,
          mode: OsFileMode.ReadWriteExecute,
        },
      ],
    },
  ],
};
const InitUserFile = {
  type: 'dir',
  name: '',
  mode: OsFileMode.ReadWriteExecute,
  children: [
    {
      name: 'Desktop',
      type: 'dir',
      mode: OsFileMode.ReadWriteExecute,
    },
    {
      name: 'Magnet',
      type: 'dir',
      mode: OsFileMode.ReadWriteExecute,
    },
    {
      name: 'Menulist',
      type: 'dir',
      mode: OsFileMode.ReadWriteExecute,
    },
    {
      name: 'Doc',
      type: 'dir',
      mode: OsFileMode.ReadWriteExecute,
    },
    {
      name: 'Note',
      type: 'dir',
      mode: OsFileMode.ReadWriteExecute,
    },
    {
      name: 'Markdown',
      type: 'dir',
      mode: OsFileMode.ReadWriteExecute,
    },
    {
      name: 'PPT',
      type: 'dir',
      mode: OsFileMode.ReadWriteExecute,
    },
    {
      name: 'Baiban',
      type: 'dir',
      mode: OsFileMode.ReadWriteExecute,
    },
    {
      name: 'Kanban',
      type: 'dir',
      mode: OsFileMode.ReadWriteExecute,
    },
    {
      name: 'Execl',
      type: 'dir',
      mode: OsFileMode.ReadWriteExecute,
    },
    {
      name: 'Photo',
      type: 'dir',
      mode: OsFileMode.ReadWriteExecute,
    },
    {
      name: 'Music',
      type: 'dir',
      mode: OsFileMode.ReadWriteExecute,
    },
    {
      name: 'Mind',
      type: 'dir',
      mode: OsFileMode.ReadWriteExecute,
    },
    {
      name: 'Schedule',
      type: 'dir',
      mode: OsFileMode.ReadWriteExecute,
    },
    {
      name: 'Reciv',
      type: 'dir',
      mode: OsFileMode.ReadWriteExecute,
    },
  ],
};

export { InitSystemFile, InitUserFile };
