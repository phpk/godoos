import { ElMessageBox } from 'element-plus'
import {BrowserWindow} from "./BrowserWindow"
class Dialog {
  constructor() {
    // static class
  }
  public static showProcessDialog(option: {
    message?: string;
    type?: 'info' | 'error' | 'question' | 'warning';
    title?: string;
    buttons?: string[];
  }) {
    console.log(option)
    const opt = Object.assign(
      {
        message: '',
        type: 'info',
        title: '提示',
        buttons: ['OK'],
      },
      option
    );

    const process = ref(0);

    const dialogwin = new BrowserWindow({
      width: 300,
      height: 150,
      content: "DialogProcess",
      title: opt.title,
      resizable: false,
      minimizable: false,
      center: true,
      skipTaskbar: true,
      config: {
        res: process,
        option: opt,
      },
      alwaysOnTop: true,
    });
    dialogwin.show();

    function setProgress(value: number) {
      process.value = value;
      if (value >= 100) {
        //dialogwin.close();
      }
    }

    return {
      setProgress,
    };
  }
  public static showMessageBox(option: {
    message?: string;
    type?: any;
    title?: string;
    buttons?: string[];
  }): Promise<{
    response: number;
  }> {
    const opt = Object.assign(
      {
        message: '',
        type: 'info',
        title: '提示',
        buttons: ['OK'],
      },
      option
    );

    let promres: (value: { response: number }) => void = () => {
      // do nothing
    };

    const porm = new Promise<{
      response: number;
    }>((resolve) => {
      promres = resolve;
    });
    ElMessageBox.confirm(
      opt.message,
      opt.title,
      {
        confirmButtonText: 'OK',
        cancelButtonText: 'Cancel',
        type: opt.type,
      }
    )
      .then(() => {
        promres({
          response: -1,
        });
      })
      .catch(() => {
        promres({
          response: 1,
        });
      })

    return porm;
  }
}

export { Dialog };
