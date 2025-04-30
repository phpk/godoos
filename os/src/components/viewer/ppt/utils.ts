
import { SVGPathData } from 'svg-pathdata'
export const calculateRotatedPosition = (
    x: number,
    y: number,
    w: number,
    h: number,
    ox: number,
    oy: number,
    k: number,
  ) => {
    const radians = k * (Math.PI / 180)
  
    const containerCenterX = x + w / 2
    const containerCenterY = y + h / 2
  
    const relativeX = ox - w / 2
    const relativeY = oy - h / 2
  
    const rotatedX = relativeX * Math.cos(radians) + relativeY * Math.sin(radians)
    const rotatedY = -relativeX * Math.sin(radians) + relativeY * Math.cos(radians)
  
    const graphicX = containerCenterX + rotatedX
    const graphicY = containerCenterY + rotatedY
  
    return { x: graphicX, y: graphicY }
  }
  export const parseLineElement = (el: any) => {
    let start: [number, number] = [0, 0]
    let end: [number, number] = [0, 0]
  
    if (!el.isFlipV && !el.isFlipH) { // 右下
      start = [0, 0]
      end = [el.width, el.height]
    }
    else if (el.isFlipV && el.isFlipH) { // 左上
      start = [el.width, el.height]
      end = [0, 0]
    }
    else if (el.isFlipV && !el.isFlipH) { // 右上
      start = [0, el.height]
      end = [el.width, 0]
    }
    else { // 左下
      start = [el.width, 0]
      end = [0, el.height]
    }
  
    const data: any = {
      type: 'line',
      width: el.borderWidth || 1,
      left: el.left,
      top: el.top,
      start,
      end,
      style: el.borderType,
      color: el.borderColor,
      points: ['', /straightConnector/.test(el.shapType) ? 'arrow' : '']
    }
    if (/bentConnector/.test(el.shapType)) {
      data.broken2 = [
        Math.abs(start[0] - end[0]) / 2,
        Math.abs(start[1] - end[1]) / 2,
      ]
    }
  
    return data
  }
  export const getSvgPathRange = (path: string) => {
    try {
      const pathData = new SVGPathData(path)
      const xList = []
      const yList = []
      for (const item of pathData.commands) {
        const x = ('x' in item) ? item.x : 0
        const y = ('y' in item) ? item.y : 0
        xList.push(x)
        yList.push(y)
      }
      return {
        minX: Math.min(...xList),
        minY: Math.min(...yList),
        maxX: Math.max(...xList),
        maxY: Math.max(...yList),
      }
    }
    catch {
      return {
        minX: 0,
        minY: 0,
        maxX: 0,
        maxY: 0,
      }
    }
  }
  export const convertFontSizePtToPx = (html: string, ratio: number) => {
    return html.replace(/font-size:\s*([\d.]+)pt/g, (match, p1) => {
      return `font-size: ${(parseFloat(p1) * ratio).toFixed(1)}px`;
    });
  };
  