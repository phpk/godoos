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
function transformElement(element:any) {
  if (element.children) {
      var children = _.map(element.children, transformElement);
      element = {...element, children: children};
  }

  if (element.type === "paragraph") {
      element = transformParagraph(element);
  }

  return element;
}

function transformParagraph(element:any) {
  if (element.alignment === "center" && !element.styleId) {
      return {...element, styleId: "Heading2"};
  } else {
      return element;
  }
}

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
        "comment-reference => sup"
      ],
      convertImage: mammoth.images.imgElement(function (image) {
        return image.read("base64").then(function (imageBuffer) {
          return {
            src: "data:" + image.contentType + ";base64," + imageBuffer
          };
        });
      })
    })
    command.executeSetHTML({
      main: result.value
    })
  }
}
