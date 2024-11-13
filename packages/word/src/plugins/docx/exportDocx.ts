import Color from 'color'
import {
  IElement,
  ElementType,
  TitleLevel,
  ListStyle,
  Command
} from '../../editor'
import {
  Document,
  Packer,
  Paragraph,
  Header,
  Footer,
  Table,
  HeadingLevel,
  ParagraphChild,
  TextRun,
  Tab,
  ExternalHyperlink,
  ImageRun,
  WidthType,
  TableRow,
  TableCell,
  MathRun
} from 'docx'
import { saveAs } from './utils'

// 标题映射
const titleLevelToHeadingLevel = {
  [TitleLevel.FIRST]: HeadingLevel.HEADING_1,
  [TitleLevel.SECOND]: HeadingLevel.HEADING_2,
  [TitleLevel.THIRD]: HeadingLevel.HEADING_3,
  [TitleLevel.FOURTH]: HeadingLevel.HEADING_4,
  [TitleLevel.FIFTH]: HeadingLevel.HEADING_5,
  [TitleLevel.SIXTH]: HeadingLevel.HEADING_6
}

function convertElementToParagraphChild(element: IElement): ParagraphChild {
  if (element.type === ElementType.IMAGE) {
    return new ImageRun({
      type: 'png',
      data: element.value,
      transformation: {
        width: element.width!,
        height: element.height!
      }
    })
  }
  if (element.type === ElementType.HYPERLINK) {
    return new ExternalHyperlink({
      children: [
        new TextRun({
          text: element.valueList?.map(child => child.value).join(''),
          style: 'Hyperlink'
        })
      ],
      link: element.url!
    })
  }
  if (element.type === ElementType.TAB) {
    return new TextRun({
      children: [new Tab()]
    })
  }
  if (element.type === ElementType.LATEX) {
    return new MathRun(element.value)
  }
  return new TextRun({
    font: element.font,
    text: element.value,
    bold: element.bold,
    size: `${(element.size || 16) / 0.75}pt`,
    color: Color(element.color).hex() || '#000000',
    italics: element.italic,
    strike: element.strikeout,
    // highlight: element.highlight ? Color(element.highlight).hex() : undefined,
    superScript: element.type === ElementType.SUPERSCRIPT,
    subScript: element.type === ElementType.SUBSCRIPT,
    underline: element.underline ? {} : undefined
  })
}

type DocxChildren = (Paragraph | Table)[]
function convertElementListToDocxChildren(
  elementList: IElement[]
): DocxChildren {
  const children: DocxChildren = []

  let paragraphChild: ParagraphChild[] = []

  function appendParagraph() {
    if (paragraphChild.length) {
      children.push(
        new Paragraph({
          children: paragraphChild
        })
      )
      paragraphChild = []
    }
  }

  for (let e = 0; e < elementList.length; e++) {
    const element = elementList[e]
    if (element.type === ElementType.TITLE) {
      appendParagraph()
      children.push(
        new Paragraph({
          heading: titleLevelToHeadingLevel[element.level!],
          children:
            element.valueList?.map(child =>
              convertElementToParagraphChild(child)
            ) || []
        })
      )
    } else if (element.type === ElementType.LIST) {
      appendParagraph()
      // 拆分列表
      const listChildren =
        element.valueList
          ?.map(item => item.value)
          .join('')
          .split('\n')
          .map(
            (text, index) =>
              new Paragraph({
                children: [
                  new TextRun({
                    text: `${
                      !element.listStyle ||
                      element.listStyle === ListStyle.DECIMAL
                        ? `${index + 1}. `
                        : `• `
                    }${text}`
                  })
                ]
              })
          ) || []
      children.push(...listChildren)
    } else if (element.type === ElementType.TABLE) {
      appendParagraph()
      const { trList } = element
      const tableRowList: TableRow[] = []
      for (let r = 0; r < trList!.length; r++) {
        const tdList = trList![r].tdList
        const tableCellList: TableCell[] = []
        for (let c = 0; c < tdList.length; c++) {
          const td = tdList[c]
          tableCellList.push(
            new TableCell({
              columnSpan: td.colspan,
              rowSpan: td.rowspan,
              children: convertElementListToDocxChildren(td.value) || []
            })
          )
        }
        tableRowList.push(
          new TableRow({
            children: tableCellList
          })
        )
      }
      children.push(
        new Table({
          rows: tableRowList,
          width: {
            size: '100%',
            type: WidthType.PERCENTAGE
          }
        })
      )
    } else if (element.type === ElementType.DATE) {
      paragraphChild.push(
        ...(element.valueList?.map(child =>
          convertElementToParagraphChild(child)
        ) || [])
      )
    } else {
      if (/^\n/.test(element.value)) {
        appendParagraph()
        element.value = element.value.replace(/^\n/, '')
      }
      paragraphChild.push(convertElementToParagraphChild(element))
    }
  }
  appendParagraph()
  return children
}

export interface IExportDocxOption {
  fileName: string,
  isFile:boolean
}

declare module '../../editor' {
  interface Command {
    executeExportDocx(options: IExportDocxOption): void
  }
}

export default function (command: Command) {
  return function (options: IExportDocxOption) {
    const { fileName,isFile } = options
    const {
      data: { header, main, footer }
    } = command.getValue()

    const doc = new Document({
      sections: [
        {
          headers: {
            default: new Header({
              children: convertElementListToDocxChildren(header || [])
            })
          },
          footers: {
            default: new Footer({
              children: convertElementListToDocxChildren(footer || [])
            })
          },
          children: convertElementListToDocxChildren(main || [])
        }
      ]
    })
    if(isFile){
      Packer.toBlob(doc).then(blob => {
        saveAs(blob, `${fileName}.docx`)
      })
    }else{
      Packer.toBase64String(doc).then((base64:any) => {
        const save = {
          data: JSON.stringify({ content: base64, title: fileName }),
          type: 'exportDocx'
        }
        //console.log(save)
        window.parent.postMessage(save, '*')
      })
    }
    
  }
}
