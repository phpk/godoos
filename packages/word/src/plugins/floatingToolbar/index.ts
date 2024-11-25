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
    icon: Svgs.optimize,
    title: '优化',
    key: 'creation_optimization'
  },
  {
    icon: Svgs.checkGrammar,
    title: '纠错',
    key: 'creation_proofreading'
  },
  {
    icon: Svgs.richContent,
    title: '续写',
    key: 'creation_continuation'
  },
  '<hr/>',
  {
    icon: Svgs.translation,
    title: '翻译',
    key: 'creation_translation'
  },
  {
    icon: Svgs.summary,
    title: '总结',
    key: 'creation_summarize'
  }
]
// 创建ai选择栏
function createAIToolbar(
  container: HTMLDivElement,
  toolbarType: ToolbarType,
  editor: Editor
) {
  const toolbarItem = document.createElement('div')
  toolbarItem.classList.add(`${PLUGIN_PREFIX}-picker-ai`)
  toolbarItem.classList.add(`${PLUGIN_PREFIX}-${toolbarType}`)
  const aiBubbleMenuItems = defaultAiPanelMenus
  toolbarItem.innerHTML = `
      <i id="${toolbarType}-btn"></i>
      <div class="aie-container ce-picker-container ai-hide">
        <div class="aie-ai-panel-body">
            <div class="aie-ai-panel-body-content ai-hide">
              <div class="loader" id="aiLoader">${Svgs.refresh}</div>
              <textarea readonly id="aiTextarea"></textarea>
            </div>
            <div class="aie-ai-panel-body-input"><input id="inputOption" placeholder="告诉 AI 下一步应该如何？比如：帮我翻译成英语" type="text" />
            <button id="goAskAi" >
              <p id="aiStart">${Svgs.aiPanelStart}</p>
              <p id="aiStop" class="ai-hide">${Svgs.aiPanelStop}</p>
            </button></div>
            <div class="aie-ai-panel-body-tips">
              ${Svgs.tips}
              提示：您可以在上面输入文字或者选择下方的操作
            </div>
        </div>
        <div class="aie-ai-panel-footer ai-hide" id="footer-one">
          <div class="aie-ai-panel-footer-tips">
            您可以进行以下操作:
          </div>
          <p id="insert">${Svgs.addContent}追加</p>
          <p id="replace">${Svgs.replace}替换</p>
          <hr/>
          <p id="hide">${Svgs.cancle}舍弃</p>
        </div>
        
        <!--aie-ai-panel-actions-->
        <div class="aie-ai-panel-footer aie-ai-panel-actions" id="footer-two">
          <div class="aie-ai-panel-footer-tips">您可以进行以下操作:</div>
          ${aiBubbleMenuItems
            .map(menuItem => {
              return typeof menuItem === 'string'
                ? menuItem
                : `<p id="ai-operate" data-type="${menuItem.key}">${menuItem.icon} ${menuItem.title} </p>`
            })
            .join('')}
        </div>
      </div>
      `
  container.append(toolbarItem)
  bindAiPanelEvent(container, editor)
}
//ai回答问题，切换视图
function viewChange(container: HTMLDivElement) {
  const aiLoader = container.querySelector('#aiLoader')
  aiLoader?.classList.remove('ai-hide')
  container.querySelector('#footer-one')?.classList.remove('ai-hide')
  container.querySelector('#footer-two')?.classList.add('ai-hide')
  container.querySelector('.aie-ai-panel-body-content')?.classList.remove('ai-hide')
  container.querySelector('#aiStart')!.classList.add('ai-hide')
  container.querySelector('#aiStop')!.classList.remove('ai-hide')
  container.querySelector<HTMLButtonElement>('#goAskAi')!.disabled = true
}
//AI弹窗绑定点击事件
function bindAiPanelEvent(container: HTMLDivElement, editor: Editor) {
  const textarea = container.querySelector<HTMLTextAreaElement>('#aiTextarea')!
  // 菜单栏AI选项
  container.querySelector<HTMLDivElement>(`#ai-edit-btn`)?.addEventListener('click', () => {
    const target = container.querySelector<HTMLDivElement>(`.${PLUGIN_PREFIX}-ai-edit`)!
    const isActive = target.classList.contains('ai-active')
    isActive ? target.classList.remove('ai-active') : target.classList.add('ai-active')
    const aiDialog = container.querySelector<HTMLDivElement>('.aie-container')
    !isActive ? aiDialog?.classList.remove('ai-hide') : aiDialog?.classList.add('ai-hide')
  })
  
  //选中操作，进行处理
  const aiOptions = Array.from(container.querySelectorAll('#ai-operate'))
  aiOptions.forEach((item) => {
    item.addEventListener('click', () => {
      viewChange(container)
      const chooseType = item.getAttribute('data-type')
      editor.command.executeAiEdit(chooseType)
    })
  })

  // 搜索
  container.querySelector('#goAskAi')!.addEventListener('click', () => {
    // console.log('搜索');
    const inputOption = container.querySelector<HTMLInputElement>('#inputOption')!
    if (inputOption.value) {
      viewChange(container)
      editor.command.executeAiEdit('creation_ask', inputOption.value)
      // container.querySelector<HTMLButtonElement>('#goAskAi')!.disabled = true
    }
  })
  // 替换
  container.querySelector('#replace')!.addEventListener('click', () => {
    // console.log('替换', textarea)
    editor.command.executeReplace(textarea.value)
  })
  //  插入
  container.querySelector('#insert')!.addEventListener('click', () => {
    const aiContent = editor.command.executeAiEdit('')
    editor.command.executeReplace(aiContent + textarea.value)
    editor.command.executeSearch('')
    // console.log('插入:')
  })
  // 舍弃
  container.querySelector('#hide')!.addEventListener('click', () => {
    editor.command.executeSearch('')
    initAiDialog(container, editor)
    // console.log('舍弃')
  })
}
export function changeAiTextarea(eventData: {[key: string]: any}) {
  const aiLoader = document.querySelector('#aiLoader')
  const textarea = document.querySelector<HTMLTextAreaElement>('#aiTextarea')!
  if (eventData.data) {
    textarea.value = eventData.data
  }
  aiLoader?.classList.add('ai-hide')
  document.querySelector('#aiStart')!.classList.remove('ai-hide')
  document.querySelector('#aiStop')!.classList.add('ai-hide')
  document.querySelector<HTMLButtonElement>('#goAskAi')!.disabled = false
}
//ai弹窗初始化
function initAiDialog(container: HTMLDivElement, editor: Editor) {
  editor.command.executeSearch('')
  container.querySelector('#aiStart')!.classList.remove('ai-hide')
  container.querySelector('#aiStop')!.classList.add('ai-hide')
  container.querySelector<HTMLButtonElement>('#goAskAi')!.disabled = false
  container.querySelector<HTMLDivElement>('.aie-container')?.classList.add('ai-hide')
  container.querySelector<HTMLDivElement>(`.${PLUGIN_PREFIX}-ai-edit`)?.classList.remove('ai-active')
  container.querySelector<HTMLDivElement>('#footer-one')?.classList.add('ai-hide')
  container.querySelector<HTMLDivElement>('#footer-two')?.classList.remove('ai-hide')
  container.querySelector('.aie-ai-panel-body-content')?.classList.add('ai-hide')
  const inputOption = container.querySelector<HTMLInputElement>('#inputOption')!
  inputOption.value = ''
  const textarea = container.querySelector<HTMLTextAreaElement>('textarea')!
  textarea.value = ""
}
// 工具栏列表
const toolbarRegisterList: IToolbarRegister[] = [
  {
    render(container, editor) {
      createAIToolbar(container, ToolbarType.AI_EDIT, editor)
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
  // let aiShow = false
  // toolbarContainer.querySelector(`.${PLUGIN_PREFIX}-ai-edit`)?.addEventListener('click', () => {
  //   const aiDialog = toolbarContainer.querySelector<HTMLDivElement>('.aie-container')
  //   // aiDialog ? toggleToolbarItemActive(aiDialog, false) : ''
  //   // console.log('点击事件：', aiDialog);
  //   aiShow ? aiDialog?.classList.add('hide') : aiDialog?.classList.remove('hide')
  //   aiShow = !aiShow
  // })
  return toolbarContainer
}

function toggleToolbarVisible(toolbar: HTMLDivElement, visible: boolean) {
  visible ? toolbar.classList.remove('hide') : toolbar.classList.add('hide')
  // console.log('元素：', visible, toolbar);
  
}

function toggleToolbarItemActive(toolbarItem: HTMLDivElement, active: boolean) {
  active
    ? toolbarItem.classList.add('active')
    : toolbarItem.classList.remove('active')
}

export function floatingToolbarPlugin(editor: Editor) {
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
      initAiDialog(toolbarContainer, editor)
      return
    }
    // 定位
    const position = context.rangeRects[0]
    toolbarContainer.style.left = `${position.x}px`
    toolbarContainer.style.top = `${position.y + position.height}px`
    // 样式回显
    // const aiDom = toolbarContainer.querySelector<HTMLDivElement>(
    //   `.${PLUGIN_PREFIX}-ai-edit`
    // )
    // // console.log('样式变化：', rangeStyle);
    
    // if (aiDom) {
    //   toggleToolbarItemActive(aiDom, rangeStyle.aiEdit)
    //   // console.log('AI变化：', rangeStyle.aiEdit);
      
    //   // if (rangeStyle.aiEdit) {
    //   //   const aiDialog = toolbarContainer.querySelector<HTMLDivElement>('.aie-container')
    //   //   aiDialog ? toggleToolbarItemActive(aiDialog, rangeStyle.aiEdit) : ''
    //   // }
    // }
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
