import '@simonwep/pickr/dist/themes/nano.min.css'
import Pickr from '@simonwep/pickr'
import Editor from '../../editor'
import './style/index.scss'
import { ToolbarType } from './enum'
import { IToolbarRegister } from './interface'
import { PLUGIN_PREFIX } from './constant'
import { Svgs } from './icons/Svgs'

function createPickerToolbar(
  container: HTMLDivElement,
  toolbarType: ToolbarType,
  changed: (color: string) => void
) {
  const toolbarItem = document.createElement('div')
  toolbarItem.classList.add(`${PLUGIN_PREFIX}-picker`)
  toolbarItem.classList.add(`${PLUGIN_PREFIX}-${toolbarType}`)
  // 颜色选择容器
  const pickerContainer = document.createElement('div')
  pickerContainer.classList.add(`${PLUGIN_PREFIX}-picker-container`)
  const pickerDom = document.createElement('div')
  pickerContainer.append(pickerDom)
  toolbarItem.append(pickerContainer)
  container.append(toolbarItem)
  // 实例化颜色选择器
  const currentColor = '#000000'
  const pickr = new Pickr({
    el: pickerDom,
    theme: 'nano',
    useAsButton: true,
    inline: true,
    default: currentColor,
    i18n: {
      'btn:save': '✓'
    },
    components: {
      preview: true,
      opacity: true,
      hue: true,
      interaction: {
        input: true,
        save: true
      }
    }
  })
  const icon = document.createElement('i')
  toolbarItem.append(icon)
  const colorBar = document.createElement('span')
  colorBar.style.backgroundColor = currentColor
  toolbarItem.append(colorBar)
  toolbarItem.onclick = evt => {
    const target = evt.target as HTMLElement
    if (pickerContainer !== target && !pickerContainer.contains(target)) {
      pickr.show()
    }
  }
  pickr.on('save', (cb: any) => {
    pickr.hide()
    const color = cb.toHEXA().toString()
    colorBar.style.backgroundColor = color
    changed(color)
  })
}
export const defaultAiPanelMenus = [
  {
    prompt: `<content>{content}</content>\n请帮我优化一下这段内容，并直接返回优化后的结果。\n注意：你应该先判断一下这句话是中文还是英文，如果是中文，请给我返回中文的内容，如果是英文，请给我返回英文内容，只需要返回内容即可，不需要告知我是中文还是英文。`,
    icon: Svgs.optimize,
    title: '改进写作'
  },
  {
    prompt: `<content>{content}</content>\n请帮我检查一下这段内容，是否有拼写错误或者语法上的错误。\n注意：你应该先判断一下这句话是中文还是英文，如果是中文，请给我返回中文的内容，如果是英文，请给我返回英文内容，只需要返回内容即可，不需要告知我是中文还是英文。`,
    icon: Svgs.checkGrammar,
    title: '检查拼写和语法'
  },
  {
    prompt: `<content>{content}</content>\n这句话的内容较长，帮我简化一下这个内容，并直接返回简化后的内容结果。\n注意：你应该先判断一下这句话是中文还是英文，如果是中文，请给我返回中文的内容，如果是英文，请给我返回英文内容，只需要返回内容即可，不需要告知我是中文还是英文。`,
    icon: Svgs.simplification,
    title: '简化内容'
  },
  {
    prompt: `<content>{content}</content>\n这句话的内容较简短，帮我简单的优化和丰富一下内容，并直接返回优化后的结果。注意：优化的内容不能超过原来内容的 2 倍。\n注意：你应该先判断一下这句话是中文还是英文，如果是中文，请给我返回中文的内容，如果是英文，请给我返回英文内容，只需要返回内容即可，不需要告知我是中文还是英文。`,
    icon: Svgs.richContent,
    title: '丰富内容'
  },
  '<hr/>',
  {
    prompt: `<content>{content}</content>\n请帮我翻译以上内容，在翻译之前，想先判断一下这个内容是不是中文，如果是中文，则翻译问英文，如果是其他语言，则需要翻译为中文，注意，你只需要返回翻译的结果，不需要对此进行任何解释，不需要除了翻译结果以外的其他任何内容。`,
    icon: Svgs.translation,
    title: '翻译'
  },
  {
    prompt: `<content>{content}</content>\n请帮我总结以上内容，并直接返回总结的结果\n注意：你应该先判断一下这句话是中文还是英文，如果是中文，请给我返回中文的内容，如果是英文，请给我返回英文内容，只需要返回内容即可，不需要告知我是中文还是英文。`,
    icon: Svgs.summary,
    title: '总结'
  }
]
// 创建ai选择栏
function createAIToolbar(
  container: HTMLDivElement,
  toolbarType: ToolbarType,
  changed: (color: string) => void
) {
  const toolbarItem = document.createElement('div')
  toolbarItem.classList.add(`${PLUGIN_PREFIX}-picker-ai`)
  toolbarItem.classList.add(`${PLUGIN_PREFIX}-${toolbarType}`)
  const chooseBox = document.createElement('div')
  chooseBox.classList.add('ce-picker-container')
  toolbarItem.append(chooseBox)
  const aiBubbleMenuItems = defaultAiPanelMenus
  toolbarItem.innerHTML = `
      <i></i>
      <div class="aie-container">
        <div class="aie-ai-panel-body">
            <div class="aie-ai-panel-body-content" style="display: none"><div class="loader">
              ${Svgs.refresh}
            </div><textarea readonly></textarea></div>
            <div class="aie-ai-panel-body-input"><input id="prompt" placeholder="告诉 AI 下一步应该如何？比如：帮我翻译成英语" type="text" />
            <button type="button" id="go" style="width: 30px;height: 30px">
              ${Svgs.aiPanelStart}
            </button></div>
            <div class="aie-ai-panel-body-tips">
              ${Svgs.tips}
              提示：您可以在上面输入文字或者选择下方的操作
            </div>
        </div>
        <div class="aie-ai-panel-footer" style="display: none">
          <div class="aie-ai-panel-footer-tips">
            您可以进行以下操作:
          </div>
          <p id="insert">${Svgs.addContent}追加</p>
          <p id="replace">${Svgs.replace}替换</p>
          <hr/>
          <p id="hide">${Svgs.cancle}舍弃</p>
        </div>
        
        <!--aie-ai-panel-actions-->
        <div class="aie-ai-panel-footer aie-ai-panel-actions" >
        <div class="aie-ai-panel-footer-tips">您可以进行以下操作:</div>
        ${aiBubbleMenuItems
      .map(menuItem => {
        return typeof menuItem === 'string'
          ? menuItem
          : `<p data-prompt="${menuItem.prompt}">${menuItem.icon} ${menuItem.title} </p>`
      })
      .join('')}
        </div>
      </div>
      `
  container.append(toolbarItem)
}

// 工具栏列表
const toolbarRegisterList: IToolbarRegister[] = [
  {
    // key: ToolbarType.AI_EDIT,
    // callback(editor) {
    //   console.log('editor:', editor)
    // }
    render(container, editor) {
      createAIToolbar(container, ToolbarType.AI_EDIT, color => {
        editor.command.executeColor(color)
      })
    }
  },
  {
    isDivider: true
  },
  {
    key: ToolbarType.SIZE_ADD,
    callback(editor) {
      editor.command.executeSizeAdd()
    }
  },
  {
    key: ToolbarType.SIZE_MINUS,
    callback(editor) {
      editor.command.executeSizeMinus()
    }
  },
  {
    isDivider: true
  },
  {
    key: ToolbarType.BOLD,
    callback(editor) {
      editor.command.executeBold()
    }
  },
  {
    key: ToolbarType.ITALIC,
    callback(editor) {
      editor.command.executeItalic()
    }
  },
  {
    key: ToolbarType.UNDERLINE,
    callback(editor) {
      editor.command.executeUnderline()
    }
  },
  {
    key: ToolbarType.STRIKEOUT,
    callback(editor) {
      editor.command.executeStrikeout()
    }
  },
  {
    isDivider: true
  },
  {
    render(container, editor) {
      createPickerToolbar(container, ToolbarType.COLOR, color => {
        editor.command.executeColor(color)
      })
    }
  },
  {
    render(container, editor) {
      createPickerToolbar(container, ToolbarType.HIGHLIGHT, color => {
        editor.command.executeHighlight(color)
      })
    }
  }
]

function createToolbar(editor: Editor): HTMLDivElement {
  const toolbarContainer = document.createElement('div')
  toolbarContainer.classList.add(`${PLUGIN_PREFIX}-floating-toolbar`)
  for (const toolbar of toolbarRegisterList) {
    if (toolbar.render) {
      toolbar.render(toolbarContainer, editor)
    } else if (toolbar.isDivider) {
      const divider = document.createElement('div')
      divider.classList.add(`${PLUGIN_PREFIX}-divider`)
      toolbarContainer.append(divider)
    } else {
      const { key, callback } = toolbar
      const toolbarItem = document.createElement('div')
      toolbarItem.classList.add(`${PLUGIN_PREFIX}-${key}`)
      const icon = document.createElement('i')
      toolbarItem.append(icon)
      toolbarItem.onclick = () => {
        callback?.(editor)
      }
      toolbarContainer.append(toolbarItem)
    }
  }
  return toolbarContainer
}

function toggleToolbarVisible(toolbar: HTMLDivElement, visible: boolean) {
  visible ? toolbar.classList.remove('hide') : toolbar.classList.add('hide')
}

function toggleToolbarItemActive(toolbarItem: HTMLDivElement, active: boolean) {
  active
    ? toolbarItem.classList.add('active')
    : toolbarItem.classList.remove('active')
}

export default function floatingToolbarPlugin(editor: Editor) {
  // 创建工具栏
  const toolbarContainer = createToolbar(editor)
  const editorContainer = editor.command.getContainer()
  editorContainer.append(toolbarContainer)

  // 监听选区样式变化
  editor.eventBus.on('rangeStyleChange', rangeStyle => {
    if (rangeStyle.type === null) {
      toggleToolbarVisible(toolbarContainer, false)
      return
    }
    const context = editor.command.getRangeContext()
    if (!context || context.isCollapsed || !context.rangeRects[0]) {
      toggleToolbarVisible(toolbarContainer, false)
      return
    }
    // 定位
    const position = context.rangeRects[0]
    toolbarContainer.style.left = `${position.x}px`
    toolbarContainer.style.top = `${position.y + position.height}px`
    // 样式回显
    const aiDom = toolbarContainer.querySelector<HTMLDivElement>(
      `.${PLUGIN_PREFIX}-ai-edit`
    )
    if (aiDom) {
      toggleToolbarItemActive(aiDom, rangeStyle.aiEdit)
    }
    const boldDom = toolbarContainer.querySelector<HTMLDivElement>(
      `.${PLUGIN_PREFIX}-bold`
    )
    if (boldDom) {
      toggleToolbarItemActive(boldDom, rangeStyle.bold)
    }
    const italicDom = toolbarContainer.querySelector<HTMLDivElement>(
      `.${PLUGIN_PREFIX}-italic`
    )
    if (italicDom) {
      toggleToolbarItemActive(italicDom, rangeStyle.italic)
    }
    const underlineDom = toolbarContainer.querySelector<HTMLDivElement>(
      `.${PLUGIN_PREFIX}-underline`
    )
    if (underlineDom) {
      toggleToolbarItemActive(underlineDom, rangeStyle.underline)
    }
    const strikeoutDom = toolbarContainer.querySelector<HTMLDivElement>(
      `.${PLUGIN_PREFIX}-strikeout`
    )
    if (strikeoutDom) {
      toggleToolbarItemActive(strikeoutDom, rangeStyle.strikeout)
    }
    toggleToolbarVisible(toolbarContainer, true)
  })
}
