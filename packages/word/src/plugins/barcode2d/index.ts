import { Editor, ElementType } from '../../editor'
import { BrowserQRCodeReader, BrowserQRCodeSvgWriter } from '@zxing/browser'
import { EncodeHintType } from '@zxing/library'
import { Dialog } from '../../components/dialog/Dialog'

function convertSvgElementToBase64(svgElement: HTMLElement | SVGSVGElement) {
  return `data:image/svg+xml;base64,${btoa(
    decodeURIComponent(svgElement.outerHTML)
  )}`
}

declare module '../../editor' {
  interface Command {
    executeInsertBarcode2D(
      content: string,
      width: number,
      height: number,
      hints?: Map<EncodeHintType, any>
    ): void
  }
}

export interface IBarcode2DOption {
  isRegisterDetectorContextMenu?: boolean
}

export default function barcode2DPlugin(
  editor: Editor
) {
  const command = editor.command

  // 插入二维码
  command.executeInsertBarcode2D = (
    content: string,
    width: number,
    height: number,
    hints?: Map<EncodeHintType, any>
  ) => {
    const codeWriter = new BrowserQRCodeSvgWriter()
    // 配置默认边距
    if (!hints) {
      hints = new Map<EncodeHintType, any>()
    }
    if (!hints.has(EncodeHintType.MARGIN)) {
      hints.set(EncodeHintType.MARGIN, 0)
    }
    // 生成svg元素并增加命名空间
    const svgElement = codeWriter.write(content, width, height, hints)
    svgElement.setAttribute('xmlns', 'http://www.w3.org/2000/svg')
    // 插入
    command.executeInsertElementList([
      {
        type: ElementType.IMAGE,
        value: convertSvgElementToBase64(svgElement),
        width,
        height
      }
    ])
  }

  // 识别二维码
  //if (options.isRegisterDetectorContextMenu) {
    editor.register.contextMenuList([
      {
        name: '识别二维码链接',
        icon: 'qrcode',
        when: payload => {
          return (
            payload.startElement === payload.endElement &&
            payload.startElement?.type === ElementType.IMAGE
          )
        },
        callback: async (_, context) => {
          const value = context.startElement?.value
          if (!value) return
          const codeReader = new BrowserQRCodeReader()
          const result = await codeReader.decodeFromImageUrl(value)
          const text = result.getText()
          // 如果是超链接则跳转
          // if (text && /^(http|https):.*/.test(text)) {
          //   window.open(text, '_blank')
          // }
          new Dialog({
            title: '识别内容',
            data: [
              {
                type: 'text',
                label: '识别内容',
                name: 'name',
                placeholder: '请输入文本',
                value: text
              },
            ]
          })
        }
      }
    ])
  //}
}
