<template>
  <template v-if="mode === 'detail'">
    <div class="file-item file-bar mode-detail">
      <div class="file-item_img"></div>
      <div class="file-item_title"></div>
      <div class="file-item_type">
        <span>{{ t("size") }}</span>
      </div>
      <div class="file-item_type">
        <span>{{ t("creation.time") }}</span>
      </div>
      <div class="file-item_type">
        <span>{{ t("modification.time") }}</span>
      </div>
      <div class="file-item_type">
        <span>{{ t("permission") }}</span>
      </div>
    </div>
  </template>
  <div draggable="true" class="file-item" :class="{
    chosen: chosenIndexs.includes(index),
    'no-chosen': !chosenIndexs.includes(index),
    'mode-icon': mode === 'icon',
    'mode-list': mode === 'list',
    'mode-big': mode === 'big',
    'mode-middle': mode === 'middle',
    'mode-detail': mode === 'detail',
    'drag-over': hoverIndex === index,
  }" :style="{
    '--theme-color': theme === 'light' ? '#ffffff6b' : '#3bdbff3d',
  }" v-for="(item, index) in fileList" :key="item.path" @dblclick="handleOnOpen(item)"
    @touchstart.passive="doubleTouch($event, item)" @contextmenu.stop.prevent="handleRightClick($event, item, index)"
    @drop="hadnleDrop($event, item.path)" @dragenter.prevent="handleDragEnter(index)" @dragover.prevent
    @dragleave="handleDragLeave()" @dragstart.stop="startDragApp($event, item)" @click="handleClick(index)"
    @mousedown.stop :ref="(ref: any) => {
      if (ref) {
        appPositions[index] = markRaw(ref as Element);
      }
    }
      ">
    <div class="file-item_img">
      <FileIcon :file="item" />
    </div>
    <span v-if="editIndex !== index" class="file-item_title">
      {{ getName(item) }}
    </span>
    <textarea autofocus draggable="false" @dragover.stop @dragstart.stop @dragenter.stop @mousedown.stop @dblclick.stop
      @click.stop @blur="onEditNameEnd" v-if="editIndex === index" class="file-item_title file-item_editing"
      v-model="editName"></textarea>
    <template v-if="mode === 'detail'">
      <div class="file-item_type">
        <span>{{ item.isDirectory ? "-" : dealSize(item.size) }}</span>
      </div>
      <div class="file-item_type">
        <span>{{ item.birthtime.toLocaleString() }}</span>
      </div>
      <div class="file-item_type">
        <span>{{ item.mtime.toLocaleString() }}</span>
      </div>
      <div class="file-item_type">
        <span>{{ item.mode?.toString?.(8) || "unknow" }}</span>
      </div>
    </template>
  </div>
</template>
<script lang="ts" setup>
import { useSystem, basename, OsFileWithoutContent, Notify, BrowserWindow } from '@/system/index.ts';
import { getSystemKey } from '@/system/config'
import { emitEvent, mountEvent } from '@/system/event';
import { useContextMenu } from '@/hook/useContextMenu.ts';
import { t, dealSystemName } from '@/i18n';
import { useFileDrag } from '@/hook/useFileDrag';
import { useAppMenu } from '@/hook/useAppMenu';
import { onMounted, ref, markRaw } from 'vue';
import { Rect } from '@/hook/useRectChosen';
import { throttle } from '@/util/debounce';
import { dealSize } from '@/util/file';
import { Menu } from '@/system/menu/Menu';
import { useChooseStore } from "@/stores/choose";
const { openPropsWindow, copyFile, createLink, deleteFile } = useContextMenu();
const sys = useSystem();
const { startDrag, folderDrop } = useFileDrag(sys);
const choose = useChooseStore()
const props = defineProps({
  onChosen: {
    type: Function,
    required: true,
  },
  onOpen: {
    type: Function,
    default: () => {
      //
    },
  },
  onRefresh: {
    type: Function,
    default: () => {
      //
    },
  },
  fileList: {
    type: Array<OsFileWithoutContent>,
    default: () => [],
  },
  theme: {
    type: String || Object,
    default: 'light',
  },
  mode: {
    type: String,
    default: 'icon',
  },
});

function getName(item: any) {
  const name = dealSystemName(basename(item.path))
  // console.log(name)
  // console.log(item.path)
  if (name.endsWith('.exe')) {
    return t(name.replace('.exe', ""))
  } else {
    return name
  }
}
function handleOnOpen(item: OsFileWithoutContent) {
  // props.onOpen(item);
  // emitEvent('desktop.app.open');
  chosenIndexs.value = [];
  if (choose.ifShow && !item.isDirectory) {
    choose.path.push(item.path)
    choose.close()
  } else {
    // console.log(' file list:',props.fileList);
    
    props.onOpen(item);
    emitEvent('desktop.app.open');
  }
}
function hadnleDrop(mouse: DragEvent, path: string) {
  hoverIndex.value = -1;
  folderDrop(mouse, path);
  chosenIndexs.value = [];
}
let expired: number | null = null;
function doubleTouch(e: TouchEvent, item: OsFileWithoutContent) {
  if (e.touches.length === 1) {
    if (!expired) {
      expired = e.timeStamp + 400;
    } else if (e.timeStamp <= expired) {
      // remove the default of this event ( Zoom )
      handleOnOpen(item);
      e.preventDefault();
      // then reset the variable for other "double Touches" event
      expired = null;
    } else {
      // if the second touch was expired, make it as it's the first
      expired = e.timeStamp + 400;
    }
  }
}

const editIndex = ref<number>(-1);
const editName = ref<string>('');
function onEditNameEnd() {
  const editEndName = editName.value.trim();
  if (editEndName && editIndex.value >= 0) {
    const editpath: any = props.fileList[editIndex.value].path.toString()
    let newPath: any;
    let sp = "/"
    if (editpath.indexOf("/") === -1) {
      sp = "\\"
    }
    newPath = editpath?.split(sp);
    newPath.pop()
    newPath.push(editEndName)
    newPath = newPath.join(sp)
    sys?.fs.rename(
      editpath,
      newPath
    );
    props.onRefresh();
    if (newPath.indexOf("Desktop") !== -1) {
      sys.refershAppList()
    }

  }
  editIndex.value = -1;
}
mountEvent('edit.end', () => {
  onEditNameEnd();
});

const hoverIndex = ref<number>(-1);
const appPositions = ref<Array<Element>>([]);

const chosenIndexs = ref<Array<number>>([]);
function handleClick(index: number) {
  chosenIndexs.value = [index];
}
onMounted(() => {
  chosenIndexs.value = [];
  props.onChosen(
    throttle((rect: Rect) => {
      const tempChosen: number[] = [];
      appPositions.value.forEach((el, index) => {
        const rect2 = el.getBoundingClientRect();
        const rect2Center = {
          x: rect2.left + rect2.width / 2,
          y: rect2.top + rect2.height / 2,
        };
        if (
          rect2Center.x > rect.left &&
          rect2Center.x < rect.left + rect.width &&
          rect2Center.y > rect.top &&
          rect2Center.y < rect.top + rect.height
        ) {
          tempChosen.push(index);
        }
      });
      chosenIndexs.value = tempChosen;
    }, 100)
  );
});

function startDragApp(mouse: DragEvent, item: OsFileWithoutContent) {
  if (chosenIndexs.value.length) {
    startDrag(
      mouse,
      chosenIndexs.value.map((index) => {
        return props.fileList[index];
      }),
      () => {
        chosenIndexs.value = [];
      }
    );
  } else {
    startDrag(mouse, [item], () => {
      chosenIndexs.value = [];
    });
  }
}

function handleRightClick(mouse: MouseEvent, item: OsFileWithoutContent, index: number) {
  
  if (chosenIndexs.value.length <= 1) {
    chosenIndexs.value = [props.fileList.findIndex((app) => app.path === item.path)];
  }
  const ext = item.name.split(".").pop();

  const zipSucess = (res: any) => {
    if (!res || res.code < 0) {
      new Notify({
        title: t('tips'),
        content: t('error'),
      });
    } else {
      props.onRefresh();
      new Notify({
        title: t('tips'),
        content: t('file.zip.success'),
      });
      if (item.parentPath == '/C/Users/Desktop') {
        sys.refershAppList()
      }
    }

  };
  // eslint-disable-next-line prefer-const
  let menuArr: any = [
    {
      label: t('open'),
      click: () => {
        chosenIndexs.value = [];
        props.onOpen(item);
      },
    },
    // {
    //   label: t('open.with'),
    //   click: () => {
    //     chosenIndexs.value = [];
    //     openWith(item);
    //   },
    // },

  ];
  if (item.isDirectory && !item.isShare) {
    if (getSystemKey('storeType') == 'local') {
      menuArr.push({
        label: t('zip'),
        submenu: [
          {
            label: 'zip',
            click: () => {
              sys.fs.zip(item.path, 'zip').then((res: any) => {
                zipSucess(res)
              });
            },
          },
          {
            label: 'tar',
            click: () => {
              sys.fs.zip(item.path, 'tar').then((res: any) => {
                zipSucess(res)
              });
            },
          },
          {
            label: 'gz',
            click: () => {
              sys.fs.zip(item.path, 'gz').then((res: any) => {
                zipSucess(res)
              });
            },
          },
        ],
      })
    } else {
      menuArr.push({
        label: t('zip'),
        click: () => {
          sys.fs.zip(item.path, 'zip').then((res: any) => {
            zipSucess(res)
          });
        },
      })
    }

  }
  if (choose.ifShow) {
    menuArr.push({
      label: "选中发送",
      click: () => {
        const paths: any = []
        chosenIndexs.value.forEach((index) => {
          const item = props.fileList[index];
          paths.push(item.path)
        })
        if (paths.length > 0) {
          choose.path = paths
          choose.close()
        }
        chosenIndexs.value = [];
      },
    })
  }
  // eslint-disable-next-line prefer-const
  let extMenus = useAppMenu(item, sys, props);
  if (extMenus && extMenus.length > 0) {
    // eslint-disable-next-line prefer-spread
    menuArr.push.apply(menuArr, extMenus)
  }
  if (ext != 'exe') {
    const fileMenus = [
      {
        label: t('rename'),
        click: () => {
          editIndex.value = index;
          editName.value = basename(item.path);
          chosenIndexs.value = [];
        },
      },
      {
        label: t('copy'),
        click: () => {
          //if(["/","/B"].includes(item.path)) return;
          copyFile(chosenIndexs.value.map((index) => props.fileList[index]));
          chosenIndexs.value = [];
        },
      },

    ];
    if (!item.isShare || item.path.indexOf('/F/othershare') !== 0) {
      fileMenus.push({
        label: t('delete'),
        click: async () => {
          for (let i = 0; i < chosenIndexs.value.length; i++) {
            await deleteFile(props.fileList[chosenIndexs.value[i]]);
          }
          chosenIndexs.value = [];
          props.onRefresh();
          if (item.path.indexOf('Desktop') > -1) {
            sys.refershAppList()
          }
        },
      })
    }
    const userType = sys.getConfig('userType');
    if (userType == 'member' && !item.isShare && !item.isDirectory) {

      menuArr.push(
        {
          label: '分享给...',
          click: () => {
            const win = new BrowserWindow({
              title: '分享',
              content: "ShareFiles",
              config: {
                path: item.path,
              },
              width: 500,
              height: 500,
              center: true,
            });
            win.show();
          },
        }
      )
      menuArr.push(
        {
          label: '评论',
          click: () => {
            const win = new BrowserWindow({
              title: '评论',
              content: "CommentsFiles",
              config: {
                path: item.path,
              },
              width: 350,
              height: 400,
              center: true,
            });
            win.show();
          },
        }
      )
    }
    // eslint-disable-next-line prefer-spread
    menuArr.push.apply(menuArr, fileMenus)
  }
  const sysEndMenu = [

    {
      label: t('create.shortcut'),
      click: () => {
        createLink(item.path)?.then(() => {
          chosenIndexs.value = [];
          props.onRefresh();
        });
      },
    },
    {
      label: t('props'),
      click: () => {
        chosenIndexs.value.forEach((index) => {
          openPropsWindow(props.fileList[index].path);
          chosenIndexs.value = [];
        });
      },
    }
  ];
  // eslint-disable-next-line prefer-spread
  menuArr.push.apply(menuArr, sysEndMenu)
  //console.log(item)

  //console.log(ext)

  Menu.buildFromTemplate(menuArr).popup(mouse);
}

function handleDragEnter(index: number) {
  hoverIndex.value = index;
}

function handleDragLeave() {
  hoverIndex.value = -1;
}

// function dealtName(name: string) {
//   return name;
// }
</script>
<style lang="scss" scoped>
.file-item {
  position: relative;
  display: flex;
  flex-direction: column;
  justify-content: flex-start;
  align-items: center;
  width: var(--desk-item-size);
  height: var(--desk-item-size);
  font-size: var(--ui-font-size);
  color: var(--icon-title-color);
  padding-top: 4px;
  border: 1px solid transparent;
  margin: 6px;

  .file-item_img {
    width: 60%;
    height: 60%;
    pointer-events: none;
  }

  .file-item_type {
    display: none;
  }

  .file-item_title {
    pointer-events: none;
  }

  .file-item_editing {
    display: inline-block !important;
    outline: none;
    pointer-events: all;
    padding: 0;
    margin: 0;
    min-width: 0;
    height: min-content !important;
    width: min-content !important;
    resize: none;
    border-radius: 0;
  }
}

.file-item:hover {
  background-color: #b1f1ff4c;
}

.chosen {
  border: 1px dashed #3bdbff3d;
  // background-color: #ffffff6b;
  background-color: var(--theme-color);

  .file-item_title {
    overflow: hidden;
    text-overflow: ellipsis;
    word-break: break-all;
    display: -webkit-box;
    -webkit-box-orient: vertical;
    -webkit-line-clamp: 2;
  }
}

.no-chosen {
  .file-item_title {
    overflow: hidden;
    text-overflow: ellipsis;
    word-break: break-all;
    display: -webkit-box;
    -webkit-box-orient: vertical;
    -webkit-line-clamp: 2;
  }
}

.drag-over {
  border: 1px dashed #3bdbff3d;
  // background-color: #ffffff6b;
  background-color: var(--theme-color);
}

.mode-icon {
  .file-item_img {
    width: 60%;
    height: calc(0.6 * var(--desk-item-size));
    margin: 0px auto;
    user-select: none;
    flex-shrink: 0;
  }

  .file-item_title {
    // color: var(--color-ui-desk-item-title);
    // height: calc(0.4 * var(--desk-item-size));
    // display: block;
    text-align: center;
    word-break: break-all;
    flex-grow: 0;
  }
}

.mode-list {
  display: flex;
  flex-direction: row;
  justify-content: flex-start;
  height: var(--menulist-item-height);
  width: var(--menulist-width);

  .file-item_img {
    width: var(--menulist-item-height);
    height: calc(0.6 * var(--menulist-item-height));

    flex-shrink: 0;
    user-select: none;
  }

  .file-item_title {
    height: min-content;
    word-break: break-all;
  }
}

.mode-icon {
  width: var(--desk-item-size);
  height: var(--desk-item-size);
}

.mode-big {
  width: calc(var(--desk-item-size) * 2.5);
  height: calc(var(--desk-item-size) * 2.5);
}

.mode-middle {
  width: calc(var(--desk-item-size) * 1.5);
  height: calc(var(--desk-item-size) * 1.5);
}

.mode-detail {
  display: flex;
  flex-direction: row;
  justify-content: flex-start;
  height: var(--menulist-item-height);
  width: 100%;
  margin: 2px;

  .file-item_img {
    width: 30px;
  }

  .file-item_title {
    width: 40%;
    display: flex;
    align-items: center;
    word-break: break-all;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .file-item_type {
    display: block;
    color: var(--color-dark-hover);
    width: 20%;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
}

.file-bar:hover {
  background-color: unset;
  user-select: none;
}
</style>
