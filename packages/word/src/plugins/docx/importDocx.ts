import { Command } from '../../editor'
import mammoth from 'mammoth'
import _ from 'lodash'

declare module '../../editor' {
  interface Command {
    executeImportDocx(options: IImportDocxOption): void
  }
}

export interface IImportDocxOption {
  arrayBuffer: ArrayBuffer
}

function transformElement(element: any) {
  if (element.children) {
    const children = _.map(element.children, transformElement);
    element = { ...element, children: children };
  }

  if (element.type === "paragraph") {
    element = transformParagraph(element);
  }

  if (element.type === "table") {
    element = transformTable(element);
  }

  if (element.type === "hyperlink") {
    element = transformHyperlink(element);
  }
  // if (element.type === "listItem") {
  //   element = transformListItem(element);
  // }


  return element;
}

function transformParagraph(element: any) {
  if (element.alignment === "center" && !element.styleId) {
    return { ...element, styleId: "Heading2" };
  }

  // 处理换行符
  if (element.text) {
    element.text = element.text.replace(/\n/g, '<br>');
  }

  // 保留居中 h1 的样式
  if (element.styleId === "Heading1" && element.alignment === "center") {
    return { ...element, alignment: "center" };
  }
  if (element.styleId === "Heading2" && element.alignment === "center") {
    return { ...element, alignment: "center" };
  }
  if (element.styleId === "Heading3" && element.alignment === "center") {
    return { ...element, alignment: "center" };
  }

  return element;
}

function transformTable(element: any) {
  const rows = element.children.map((row: any) => {
    const cells = row.children.map((cell: any) => {
      let cellStyle = '';

      // 处理单元格的背景色
      if (cell.backgroundColor) {
        cellStyle += `background-color: ${cell.backgroundColor}; `;
      }

      // 处理单元格的文字对齐方式
      if (cell.alignment) {
        cellStyle += `text-align: ${cell.alignment}; `;
      }

      // 处理单元格的边框
      if (cell.border) {
        cellStyle += `border: ${cell.border}; `;
      }

      return `<td style="${cellStyle}">${cell.text}</td>`;
    });
    return `<tr>${cells.join('')}</tr>`;
  });

  // 处理表格的整体样式
  let tableStyle = '';
  if (element.border) {
    tableStyle += `border-collapse: collapse; border: ${element.border}; `;
  }
  if (element.backgroundColor) {
    tableStyle += `background-color: ${element.backgroundColor}; `;
  }

  return { ...element, text: `<table style="${tableStyle}">${rows.join('')}</table>` };
}
function transformHyperlink(element: any) {
  return { ...element, text: `<a href="${element.href}">${element.text}</a>` };
}

const waitImgWH = (value: string) => new Promise((resolve) => {
  const dom = document.createElement('div')
  dom.innerHTML = value
  setTimeout(() => {
    resolve(dom.innerHTML)
  }, 0)
})

export default function (command: Command) {
  return async function (options: IImportDocxOption) {
    const { arrayBuffer } = options
    const result = await mammoth.convertToHtml({
      arrayBuffer,
    }, {
      transformDocument: transformElement,
      styleMap: [
        "p[style-name='Normal'] => p",
        "p[style-name='Heading 1'] => h1:fresh",
        "p[style-name='Heading 2'] => h2:fresh",
        "p[style-name='Heading 3'] => h3:fresh",
        "p[style-name='Heading 4'] => h4:fresh",
        "p[style-name='Heading 5'] => h5:fresh",
        "p[style-name='Heading 6'] => h6:fresh",
        "p[style-name='Section Title'] => h1:fresh\n",
        "p[style-name='Subsection Title'] => h2:fresh",
        "p[style-name='Equation'] => span.math-display:fresh",
        "p[style-name='Title'] => h1:fresh",
        "p[style-name='Subtitle'] => h2:fresh",
        "p[style-name='Quote'] => blockquote",
        "p[style-name='List Paragraph'] => li",
        "r[style-name='Emphasis'] => em",
        "r[style-name='Strong'] => strong",
        "r[style-name='Underline'] => u",
        "r[style-name='Strikethrough'] => del",
        "r[style-name='Comment Reference'] => sup",
        "b => strong",
        "i => em",
        "u => u",
        "strike => del",
        "comment-reference => sup",
        "a => a"
      ],
      convertImage: mammoth.images.imgElement(function (image) {
        return image.read("base64").then(function (imageBuffer) {
          return {
            src: "data:" + image.contentType + ";base64," + imageBuffer
          };
        });
      })
    });
    const value: any = result.value.includes('<img') ? await waitImgWH(result.value) : result.value
    command.executeSetHTML({
      main: value
    })
  }
}