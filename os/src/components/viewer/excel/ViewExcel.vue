<!-- <template>
    <div class="view-excel" ref="wrapperRef">
        <div class="view-excel-main" ref="rootRef"></div>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount, watch, nextTick } from 'vue';
import Spreadsheet from 'x-data-spreadsheet';
import 'x-data-spreadsheet/dist/xspreadsheet.css'
import { getData, readExcelData, transferExcelToSpreadSheet } from './excel';
import { renderImage, clearCache } from './media';
import { readOnlyInput, download as downloadFile } from './hack';
import { debounce } from 'lodash';
import { isBase64, base64ToBuffer } from '@/utils/file';

const props = defineProps<{
    src: string | ArrayBuffer | Blob | null;
    requestOptions: any;
    options: any;
}>();

const emit = defineEmits<{
    (e: 'rendered'): void;
    (e: 'error', error: Error): void;
}>();

const defaultOptions: any = {
    minColLength: 25
};

const wrapperRef = ref<HTMLElement | null>(null);
const rootRef: any = ref<HTMLElement | null>(null);
let workbookDataSource: any = {
    _worksheets: []
};
let mediasSource: MediaSource[] = [];
let sheetIndex = 0;
let ctx: CanvasRenderingContext2D | null = null;
let xs: any = null;
let offset: any = null;
let fileData: ArrayBuffer | null = null;

const renderExcel: any = (buffer: ArrayBuffer) => {
    fileData = buffer;
    readExcelData(buffer).then((workbook: any) => {
        console.log(workbook);

        if (!workbook._worksheets || workbook._worksheets.length === 0) {
            throw new Error('未获取到数据，可能文件格式不正确或文件已损坏');
        }
        let { workbookData, medias, workbookSource } = transferExcelToSpreadSheet(workbook, { ...defaultOptions, ...props.options });
        if (props.options?.transformData && typeof props.options?.transformData === 'function') {
            workbookData = props.options.transformData(workbookData);
        }
        mediasSource = medias;
        workbookDataSource = workbookSource;
        offset = null;
        sheetIndex = 0;
        clearCache();
        xs?.loadData(workbookData);
        renderImage(ctx, mediasSource, workbookDataSource._worksheets[sheetIndex], offset);
        emit('rendered');
    }).catch(e => {
        console.warn(e);
        mediasSource = [];
        workbookDataSource = {
            _worksheets: []
        };
        clearCache();
        xs?.loadData({});
        emit('error', e);
    });
};

const observerCallback: any = debounce(readOnlyInput, 200).bind(this, rootRef.value);
const observer = new MutationObserver(observerCallback);
const observerConfig = { attributes: true, childList: true, subtree: true };

onMounted(() => {
    nextTick(() => {
        observer.observe(rootRef.value!, observerConfig);
        observerCallback(rootRef.value!);

        xs = new Spreadsheet(rootRef.value!, {
            mode: 'read',
            showToolbar: false,
            showContextmenu: props.options?.showContextmenu || false,
            view: {
                height: () => wrapperRef.value?.clientHeight || 500,
                width: () => wrapperRef.value?.clientWidth || 1200,
            },
            row: {
                height: 25,
                len: 100
            },
            col: {
                len: 26,
                width: 100,
                indexWidth: 60,
                minWidth: 60,
            },
            style: {
                // 背景颜色
                bgcolor: '#ffffff',
                // 水平对齐方式
                align: 'left',
                // 垂直对齐方式
                valign: 'middle',
                // 是否需要换行
                textwrap: false,
                // 虚线边框
                strike: false,
                // 下画线
                underline: false,
                // 文字颜色
                color: '#0a0a0a',
                // 字体设置
                font: {
                    // 字体
                    name: 'Helvetica',
                    // 字号大小
                    size: 10,
                    // 是否加粗
                    bold: false,
                    // 斜体
                    italic: false,
                },
            }
            //autoFocus: false
        }).loadData({});

        let swapFunc = xs.bottombar.swapFunc;
        xs.bottombar.swapFunc = function (index: number) {
            swapFunc.call(xs.bottombar, index);
            sheetIndex = index;
            setTimeout(() => {
                xs.reRender();
                renderImage(ctx, mediasSource, workbookDataSource._worksheets[sheetIndex], offset);
            });
        };

        let clear = xs.sheet.editor.clear;
        xs.sheet.editor.clear = function (...args: any[]) {
            clear.apply(xs.sheet.editor, args);
            setTimeout(() => {
                renderImage(ctx, mediasSource, workbookDataSource._worksheets[sheetIndex], offset);
            });
        };

        let setOffset = xs.sheet.editor.setOffset;
        xs.sheet.editor.setOffset = function (...args: any[]) {
            setOffset.apply(xs.sheet.editor, args);
            offset = args[0];
            renderImage(ctx, mediasSource, workbookDataSource._worksheets[sheetIndex], offset);
        };

        const canvas = rootRef.value?.querySelector('canvas');
        ctx = canvas?.getContext('2d') || null;
        if (props.src) {
            if (typeof props.src === 'string' && isBase64(props.src)) {
                const buffer = base64ToBuffer(props.src);
                renderExcel(buffer);
            } else {
                getData(props.src, props.requestOptions).then(renderExcel).catch(e => {
                    xs?.loadData({});
                    emit('error', e);
                });
            }
        }
    });
});

onBeforeUnmount(() => {
    observer.disconnect();
    xs = null;
});

watch(() => props.src, (newSrc) => {
    if (newSrc) {
        getData(newSrc, props.requestOptions).then(renderExcel).catch(e => {
            xs?.loadData({});
            emit('error', e);
        });
    } else {
        xs?.loadData({});
    }
});

const save = (fileName?: string) => {
    downloadFile(fileName || `view-excel-${new Date().getTime()}.xlsx`, fileData);
};

defineExpose({
    save
});
</script>

<style lang="scss">
.view-excel {
    height: 550px;
}
</style> -->

<template></template>
