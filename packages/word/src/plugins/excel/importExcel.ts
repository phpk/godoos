import {
  Command,
  ElementType,
  IElement,
  VerticalAlign
} from '../../editor'
import { IColgroup } from '../../editor/interface/table/Colgroup'
import { ITd } from '../../editor/interface/table/Td'
import ExcelJS, { CellRichTextValue } from 'exceljs'

declare module '../../editor' {
  interface Command {
    executeImportExcel(options: IImportExcelOption): void
  }
}

declare module 'exceljs' {
  interface WorksheetModel {
    rows: Row[]
  }
  interface Row {
    cells: Cell[]
  }
}

export interface IImportExcelOption {
  arrayBuffer: ArrayBuffer
}

const ANCHOR_COL_WIDTH = 6 // excel是以字符串”0“，字体大小11作为宽度基础设置
const DEFAULT_COL_WIDTH_COUNT = 8.6 // 默认8.6个字符宽度
const DEFAULT_COL_HEIGHT = 40 // 默认40px

// 垂直布局映射
const EXCEL_EDITOR_VERTICAL_MAPPING = {
  top: VerticalAlign.TOP,
  middle: VerticalAlign.MIDDLE,
  bottom: VerticalAlign.BOTTOM,
  distributed: VerticalAlign.TOP,
  justify: VerticalAlign.TOP
}

export default function (command: Command) {
  return async function (options: IImportExcelOption) {
    const workbook = new ExcelJS.Workbook()
    await workbook.xlsx.load(options.arrayBuffer)
    const elementList: IElement[] = []
    // 循环每个sheet
    workbook.eachSheet(function (worksheet) {
      // 列宽
      const colgroup: IColgroup[] = worksheet.columns.map(col => ({
        width: (col.width || DEFAULT_COL_WIDTH_COUNT) * ANCHOR_COL_WIDTH
      }))
      const tableElement: IElement = {
        type: ElementType.TABLE,
        value: '',
        colgroup,
        trList: []
      }
      // 按列处理
      worksheet.eachRow((row, rowNumber) => {
        const rowIndex = rowNumber - 1
        const model = row.model
        const tdList: ITd[] = []
        if (!Array.isArray(row.values)) return
        // 行处理
        for (let v = 1; v < row.values.length; v++) {
          const cell = model?.cells?.[v - 1]
          if (cell?.master) continue
          const cellStyle = cell?.style
          const cellValue = cell?.value
          // 垂直布局
          const verticalAlign =
            EXCEL_EDITOR_VERTICAL_MAPPING[
              cellStyle?.alignment?.vertical || 'top'
            ]
          // 数据样式
          const value: IElement[] = []
          const richtext = (<CellRichTextValue>cellValue)?.richText
          if (richtext) {
            richtext.forEach(item => {
              value.push({
                value: item.text,
                bold: item.font?.bold,
                italic: item.font?.italic,
                size: item.font?.size,
                strikeout: item.font?.strike,
                underline: !!item.font?.underline
              })
            })
          } else {
            value.push({
              value: cellValue?.toString() || '',
              bold: cellStyle?.font?.bold,
              italic: cellStyle?.font?.italic,
              size: cellStyle?.font?.size,
              strikeout: cellStyle?.font?.strike
            })
          }
          // 合并单元格信息
          let colspan = 1
          let rowspan = 1
          const rowList = worksheet.model.rows
          for (let r = rowIndex; r < rowList.length; r++) {
            const nextCells = rowList[r].cells
            for (let c = 0; c < nextCells.length; c++) {
              // 忽略当前元素
              if (r === rowIndex && c === v - 1) continue
              const nextCell = nextCells[c]
              // 附属单元格
              if (nextCell.master === cell?.address) {
                if (r === rowIndex) {
                  // 相同行则为跨列
                  colspan += 1
                } else if (r > rowIndex && c === v - 1) {
                  // 下一行则为跨行
                  rowspan += 1
                }
              }
            }
          }
          tdList.push({
            colspan,
            rowspan,
            verticalAlign,
            value
          })
        }
        if (tdList.length) {
          // 行高
          const height = model?.height || DEFAULT_COL_HEIGHT
          tableElement.trList!.push({
            height,
            minHeight: height,
            tdList
          })
        }
      })
      if (tableElement.trList?.length) {
        elementList.push(tableElement)
      }
    })
    // 设置值
    if (elementList.length) {
      command.executeSetValue({
        main: elementList
      })
    }
  }
}
