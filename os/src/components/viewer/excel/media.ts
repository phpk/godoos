

let cache: any = [];

export function renderImage(ctx: any, medias: any, sheet: any, offset: any): void {
    if (sheet && sheet._media.length) {
        sheet._media.forEach((media:any) => {
            const { imageId, range, type } = media;
            if (type === 'image') {
                const position = calcPosition(sheet, range, offset);
                drawImage(ctx, imageId, medias[imageId], position);
            }
        });
    }
}

const clipWidth = 60; // 左侧序号列宽
const clipHeight = 25; // 顶部序号行高
const defaultColWidth = 80;
const defaultRowHeight = 24;
const devicePixelRatio = window.devicePixelRatio;

function calcPosition(sheet: any, range: any, offset: any){
    const { tl = {}, br = {} } = range;
    const { nativeCol, nativeColOff, nativeRow, nativeRowOff } = tl;

    let basicX = clipWidth;
    let basicY = clipHeight;
    for (let i = 0; i < nativeCol; i++) {
        basicX += (sheet?._columns?.[i]?.width * 6 || defaultColWidth);
    }
    for (let i = 0; i < nativeRow; i++) {
        basicY += (sheet?._rows?.[i]?.height || defaultRowHeight);
    }
    const x = basicX + nativeColOff / 12700;
    const y = basicY + nativeRowOff / 12700;

    const {
        nativeCol: nativeColEnd,
        nativeColOff: nativeColOffEnd,
        nativeRow: nativeRowEnd,
        nativeRowOff: nativeRowOffEnd
    } = br;
    let width;
    if (nativeCol === nativeColEnd) {
        width = (nativeColOffEnd - nativeColOff) / 12700;
    } else {
        width = (sheet?._columns?.[nativeCol]?.width * 6 || defaultColWidth) - nativeColOff / 12700;
        for (let i = nativeCol + 1; i < nativeColEnd; i++) {
            width += (sheet?._columns?.[i]?.width * 6 || defaultColWidth);
        }
        width += nativeColOffEnd / 12700;
    }
    let height;
    if (nativeRow === nativeRowEnd) {
        height = (nativeRowOffEnd - nativeRowOff) / 12700;
    } else {
        height = (sheet?._rows?.[nativeRow]?.height || defaultRowHeight) - nativeRowOff / 12700;
        for (let i = nativeRow + 1; i < nativeRowEnd; i++) {
            height += (sheet?._rows?.[i]?.height || defaultRowHeight);
        }
        height += nativeRowOffEnd / 12700;
    }

    return {
        x: (x - (offset?.scroll?.x || 0)) * devicePixelRatio,
        y: (y - (offset?.scroll?.y || 0)) * devicePixelRatio,
        width: width * devicePixelRatio,
        height: height * devicePixelRatio
    };
}

export function clearCache(): void {
    cache = [];
}

function drawImage(ctx: CanvasRenderingContext2D, index: string, data: any, position: any): void {
    getImage(index, data).then(image => {
        let sx = 0;
        let sy = 0;
        let sWidth = image.width;
        let sHeight = image.height;
        let dx = position.x;
        let dy = position.y;
        let dWidth = position.width;
        let dHeight = position.height;
        let scaleX = dWidth / sWidth;
        let scaleY = dHeight / sHeight;

        if (dx < clipWidth * devicePixelRatio) {
            let diff = clipWidth * devicePixelRatio - dx;
            dx = clipWidth * devicePixelRatio;
            dWidth -= diff;
            sWidth -= diff / scaleX;
            sx += diff / scaleX;
        }
        if (dy < clipHeight * devicePixelRatio) {
            let diff = clipHeight * devicePixelRatio - dy;
            dy = clipHeight * devicePixelRatio;
            dHeight -= diff;
            sHeight -= diff / scaleY;
            sy += diff / scaleY;
        }
        ctx.drawImage(image, sx, sy, sWidth, sHeight, dx, dy, dWidth, dHeight);
    }).catch(e => {
        console.error(e);
    });
}

function getImage(index: string, data: any): Promise<HTMLImageElement> {
    return new Promise((resolve, reject) => {
        if (cache[index]) {
            return resolve(cache[index]);
        }
        const { buffer } = data.buffer;
        let blob = new Blob([buffer], { type: 'image/' + data.extension });
        let url = URL.createObjectURL(blob);
        let image = new Image();
        image.src = url;
        image.onload = function () {
            resolve(image);
            cache[index] = image;
        };
        image.onerror = function (e) {
            reject(e);
        };
    });
}