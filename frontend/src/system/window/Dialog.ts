import { ElMessageBox } from 'element-plus'
import {BrowserWindow} from "./BrowserWindow"
import { md5 } from "js-md5"
import { notifyError } from "@/util/msg"
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
    console.log('弹窗',option)
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
        dialogwin.close();
      }
    }

    return {
      setProgress,
      dialogwin
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
  public static showInputBox(correctPwd: string): Promise<{response: number}>{
    let promres: (value: { response: number }) => void = () => {
      // do nothing
    };

    const porm = new Promise<{
      response: number;
    }>((resolve) => {
      promres = resolve;
    });
    ElMessageBox.prompt('请输入文件加密密码'   , '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
    }).then(({value}) => {
      // console.log('正确否：', correctPwd === md5(value));
      // console.log('正确密码：', correctPwd);
      const isCorrect = md5(value) === correctPwd ? 1 : 0
      if(!isCorrect) {
        notifyError('密码错误')
      }
      promres({
        response: isCorrect,
      });
    }).catch(()=>{
      console.log('取消');
      promres({
        response: -1,
      });
    })
    return porm;
  }
}

export { Dialog };
