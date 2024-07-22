import { Ref, defineComponent, markRaw, ref } from 'vue';
import { Menu } from './Menu';

export interface TrayOptions {
  image?: string | ReturnType<typeof defineComponent>;
  component: string;
}
export class Tray {
  image?: string | ReturnType<typeof defineComponent>;
  component?:string;
  _id?: string;
  _contextMenu: ReturnType<typeof defineComponent> | Menu;
  _contextMenuShow = false;
  _contextMenuWidth = 100;
  _contextMenuHeight = 100;
  public static trayList: Ref<Tray[]> = ref<Tray[]>([]);
  constructor(options: TrayOptions) {
    if (typeof options.image === 'string') {
      this.image = options.image;
    } else {
      this.image = markRaw(options.image);
    }
    if (typeof options.component === 'string') {
      this.component = options.component;
    }
    if(!Tray.trayList.value.includes(this)){
      Tray.trayList.value.push(this);
      this._id = Tray.trayList.value.length.toString();
    }
    
  }
  setContextMenu(content: ReturnType<typeof defineComponent> | Menu, width = 100, height = 100) {
    this._contextMenu = markRaw(content);
    this._contextMenuWidth = width;
    this._contextMenuHeight = height;
  }
  setImage(image: string | ReturnType<typeof defineComponent>) {
    if (typeof image === 'string') {
      this.image = image;
    } else {
      this.image = markRaw(image);
    }
    Tray.trayList.value = Tray.trayList.value.slice();
  }
  destroy() {
    const index = Tray.trayList.value.findIndex((item) => item._id === this._id);
    if (index !== -1) {
      Tray.trayList.value.splice(index, 1);
    }
  }
}
