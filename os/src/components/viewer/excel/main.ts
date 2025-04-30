import Spreadsheet from 'x-data-spreadsheet';
import {getData, readExcelData, transferExcelToSpreadSheet} from './excel';
import {renderImage, clearCache} from './media';
import {readOnlyInput,download as downloadFile} from './hack';
import {debounce} from 'lodash';
const defaultOptions = {
    minColLength: 20
};
class JsExcelPreview {
    container:any = null;
    wrapper:any = null;
    wrapperMain:any = null;
    options:any = {};
    requestOptions:any = {};
    mediasSource:any = [];
    workbookDataSource:any = {
        _worksheets:[]
    };
    sheetIndex = 1;
    ctx:any = null;
    xs:any = null;
    offset:any = null;
    observer:any = null;
    fileData:any = null;

    constructor(container:any, options={}, requestOptions={}) {
        this.container = container;
        this.options = {...defaultOptions, ...options};
        this.requestOptions = requestOptions;
        this.createWrapper();
        this.initSpreadsheet();
        this.hack();
    }
    createWrapper(){
        this.wrapper = document.createElement('div');
        this.wrapper.className = 'view-excel';
        this.wrapperMain = document.createElement('div');
        this.wrapperMain.className = 'view-excel-main';
        this.wrapper.appendChild(this.wrapperMain);
        this.container.appendChild(this.wrapper);
    }
    initSpreadsheet(){
        this.xs = new Spreadsheet(this.wrapperMain, {
            mode: 'read',
            showToolbar: false,
            showContextmenu: this.options.showContextmenu || false,
            view: {
                height: () => this.wrapper && this.wrapper.clientHeight || 300,
                width: () => this.wrapper && this.wrapper.clientWidth || 1200,
            },
            row: {
                height: 24,
                len: 100
            },
            col: {
                len: 26,
                width: 80,
                indexWidth: 60,
                minWidth: 60,
            },
           // autoFocus: false
        }).loadData({});

        let that = this;

        let swapFunc = this.xs.bottombar.swapFunc;
        this.xs.bottombar.swapFunc = function (index:any) {
            swapFunc.call(that.xs.bottombar, index);
            that.sheetIndex = index + 1;
            setTimeout(()=>{
                that.xs.reRender();
                renderImage(that.ctx, that.mediasSource,that.workbookDataSource._worksheets[that.sheetIndex], that.offset);
            });

        };
        let clear = this.xs.sheet.editor.clear;
        this.xs.sheet.editor.clear = function (...args:any){
            clear.apply(that.xs.sheet.editor, args);
            setTimeout(()=>{
                renderImage(that.ctx, that.mediasSource,that.workbookDataSource._worksheets[that.sheetIndex], that.offset);
            });
        };
        let setOffset = this.xs.sheet.editor.setOffset;
        this.xs.sheet.editor.setOffset = function (...args:any){
            setOffset.apply(that.xs.sheet.editor, args);
            that.offset = args[0];
            renderImage(that.ctx, that.mediasSource,that.workbookDataSource._worksheets[that.sheetIndex], that.offset);
        };
        const canvas = this.wrapperMain.querySelector('canvas');
        this.ctx = canvas.getContext('2d');
    }
    renderExcel(buffer:any){
        this.fileData = buffer;
        return readExcelData(buffer).then((workbook:any) => {
            if (!workbook._worksheets || workbook._worksheets.length === 0) {
                throw new Error('未获取到数据，可能文件格式不正确或文件已损坏');
            }
            const {workbookData, medias, workbookSource} = transferExcelToSpreadSheet(workbook, this.options);
            this.mediasSource = medias;
            this.workbookDataSource = workbookSource;
            this.offset = null;
            this.sheetIndex = 1;
            clearCache();
            this.xs.loadData(workbookData);
            renderImage(this.ctx, this.mediasSource,this.workbookDataSource._worksheets[this.sheetIndex], this.offset);

        }).catch(e => {
            this.mediasSource = [];
            this.workbookDataSource = {
                _worksheets:[]
            };
            clearCache();
            this.xs.loadData({});
            return Promise.reject(e);
        });
    }
    hack(){
        const observerCallback:any = debounce(readOnlyInput, 200).bind(this, this.wrapperMain);
        this.observer = new MutationObserver(observerCallback);
        const observerConfig = { attributes: true, childList: true, subtree: true };
        this.observer.observe(this.wrapperMain, observerConfig);
        observerCallback(this.wrapperMain);
    }

    setOptions(options:any) {
        this.options = options;
    }
    setRequestOptions(requestOptions:any) {
        this.requestOptions = requestOptions;
    }
    preview(src:any){
        return new Promise(((resolve, reject) => {
            getData(src, this.requestOptions).then((res)=>{
                this.renderExcel(res).then(resolve);
            }).catch(e => {
                this.xs.loadData({});
                reject(e);
            });
        }));
    }
    save(fileName:any){
        downloadFile(fileName || `js-preview-excel-${new Date().getTime()}.xlsx`,this.fileData);
    }
    destroy(){
        this.observer.disconnect();
        this.container.removeChild(this.wrapper);
        this.container = null;
        this.wrapper = null;
        this.wrapperMain = null;
        this.ctx = null;
        this.xs = null;
        this.observer = null;
        this.options = null;
        this.requestOptions = null;
        this.mediasSource = null;
        this.workbookDataSource = null;
    }
}
export function init(container:any, options:any, requestOptions:any){
    return new JsExcelPreview(container, options, requestOptions);
}