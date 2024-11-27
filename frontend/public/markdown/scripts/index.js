var markdownTitle = ""
var CustomHookA = Cherry.createSyntaxHook('codeBlock', Cherry.constants.HOOKS_TYPE_LIST.PAR, {
  makeHtml(str) {
    console.warn('custom hook', 'hello');
    return str;
  },
  rule(str) {
    const regex = {
      begin: '',
      content: '',
      end: '',
    };
    regex.reg = new RegExp(regex.begin + regex.content + regex.end, 'g');
    return regex;
  },
});

/**
 * 自定义一个自定义菜单
 * 点第一次时，把选中的文字变成同时加粗和斜体
 * 保持光标选区不变，点第二次时，把加粗斜体的文字变成普通文本
 */
var customMenuA = Cherry.createMenuHook('加粗斜体', {
  iconName: 'font',
  onClick: function (selection) {
    // 获取用户选中的文字，调用getSelection方法后，如果用户没有选中任何文字，会尝试获取光标所在位置的单词或句子
    let $selection = this.getSelection(selection) || '同时加粗斜体';
    // 如果是单选，并且选中内容的开始结束内没有加粗语法，则扩大选中范围
    if (!this.isSelections && !/^\s*(\*\*\*)[\s\S]+(\1)/.test($selection)) {
      this.getMoreSelection('***', '***', () => {
        const newSelection = this.editor.editor.getSelection();
        const isBoldItalic = /^\s*(\*\*\*)[\s\S]+(\1)/.test(newSelection);
        if (isBoldItalic) {
          $selection = newSelection;
        }
        return isBoldItalic;
      });
    }
    // 如果选中的文本中已经有加粗语法了，则去掉加粗语法
    if (/^\s*(\*\*\*)[\s\S]+(\1)/.test($selection)) {
      return $selection.replace(/(^)(\s*)(\*\*\*)([^\n]+)(\3)(\s*)($)/gm, '$1$4$7');
    }
    /**
     * 注册缩小选区的规则
     *    注册后，插入“***TEXT***”，选中状态会变成“***【TEXT】***”
     *    如果不注册，插入后效果为：“【***TEXT***】”
     */
    this.registerAfterClickCb(() => {
      this.setLessSelection('***', '***');
    });
    return $selection.replace(/(^)([^\n]+)($)/gm, '$1***$2***$3');
  }
});
/**
 * 定义一个空壳，用于自行规划cherry已有工具栏的层级结构
 */
var customMenuB = Cherry.createMenuHook('实验室', {
  icon: {
    type: 'svg',
    content: '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10" /><path d="M8 14s1.5 2 4 2 4-2 4-2" /><line x1="9" y1="9" x2="9.01" y2="9" /><line x1="15" y1="9" x2="15.01" y2="9" /></svg>',
    iconStyle: 'width: 15px; height: 15px; vertical-align: middle;',
  },
});
var exportDataHook = Cherry.createMenuHook('导出', {
  subMenuConfig: [
    {
      noIcon: true,
      name: '思维导图',
      onclick: () => {
        if (!markdownTitle || markdownTitle == "") {
          markdownTitle = window.prompt("请输入文稿标题");
        }
        const postData = { title: markdownTitle, content: cherry.getMarkdown() }
        window.parent.postMessage({ type: 'saveMind', data: postData }, '*')
      }
    },
    // {
    //   noIcon: true,
    //   name: 'PPTX',
    //   onclick: () => {
        
    //   }
    // },
    {
      noIcon: true,
      name: 'PDF',
      onclick: () => {
        cherry.export();
      }
    },
    {
      noIcon: true,
      name: '长图',
      onclick: () => {
        cherry.export('img');
      }
    },
    {
      noIcon: true,
      name: 'markdown',
      onclick: () => {
        cherry.export('markdown');
      }
    },
    {
      noIcon: true,
      name: 'html',
      onclick: () => {
        cherry.export('html');
      }
    },
  ]
});

class AiDialogClass {
  actionArr = {
    creation_leader: "生成大纲",
    creation_builder: "根据主题和提纲进行撰写",
    creation_continuation: "续写",
    creation_optimization: "优化",
    creation_proofreading: "纠错",
    creation_summarize: "总结",
    creation_translation: "翻译"
  }
  buttons = [
    {
      title: '追加',
      action: 'aiAdd',
      method: this.addAiContent
    },
    {
      title: '替换',
      action: 'aiReplace',
      method: this.replaceAiContent
    },
    {
      title: '舍弃',
      action: 'aiCancle',
      method: this.closeDialog
    },
    {
      title: '生成大纲',
      action: 'aiOutline',
      method: this.createAiOutline
    },
    {
      title: '生成文章',
      action: 'aiArticle',
      method: this.createAiArticle
    },
  ]
  constructor() {
    this.container = document.querySelector('#aiDialog')
    this.createDialog()
  }
  // 生成弹窗的基本框架
  createDialog() {
    const dialog = document.createElement('div')
    dialog.innerHTML = `
      <div class="ai-markdown-dialog" style="padding: 0px 15px;">
        <div class="ai-dialog-title">${this.action}</div>
        <textarea class="ai-dialog-content" readonly></textarea>
        <div class="button-box">
          <button class="ai-dialog-button" data-type="aiAdd">追加</button>
          <button class="ai-dialog-button" data-type="aiReplace">替换</button>
          <button class="ai-dialog-button"  data-type="aiCancle">取消</button>
        </div>
        <div class="button-box hide">
          <button class="ai-dialog-button" data-type="aiReplace">生成大纲</button>
          <button class="ai-dialog-button"  data-type="aiCancle">取消</button>
        </div>
        <div class="button-box hide">
          <button class="ai-dialog-button" data-type="aiArticle">生成文章</button>
          <button class="ai-dialog-button"  data-type="aiCancle">取消</button>
        </div>
      </div>`
    this.container.appendChild(dialog)
     // 绑定事件
    const aiOptions = Array.from(this.container.querySelectorAll('.ai-dialog-button'))
    aiOptions.forEach(item => {
      item.addEventListener('click', () => {
        const actionType = item.getAttribute('data-type')
        const btn = this.buttons.find(item => item.action == actionType)
        if (btn && btn.method) {
          btn.method(this.container)
        }
      })
    })
  }
  // 添加不同按钮
  addButton(action) {
    let pos = 0
    action == '大纲' ? pos = 1 : pos = 0
    const btnArr = Array.from(this.container.querySelectorAll('.button-box'))
    for (let i = 0; i < btnArr.length; i++) {
      if (i == pos) {
        btnArr[i].classList.remove('hide')
      }else {
        btnArr[i].classList.add('hide')
      }
    }
  }
  // 点击不同按钮生成不同弹窗
  showDialog(action, content) {
    // this.sendRequest(action,content)
    this.container.querySelector('textarea').value = content
    this.addButton(action)
  }
  // 发送请求
  sendRequest(action, data) {
    let title = 'ai助手'
    Object.keys(this.actionArr).forEach(item => {
      if (item == action) {
        title = 'AI' + this.actionArr[item]
      }
    })
    this.container.querySelector('.ai-dialog-title').innerText = title
    window.parent.postMessage({
      type: 'aiCreater',
      data,
      action
    }, '*')
    this.container.classList.remove('hide')
  }
  // 关闭弹窗
  closeDialog(dialog) {
    dialog?.classList.add('hide')
  }
  // 追加
  addAiContent(dialog) {
    const content = dialog.querySelector('textarea').value
    // cherry.insert(content,false)
    cherry.setMarkdown(content)
    // this.closeDialog(dialog)
  }
  // 替换
  replaceAiContent() {
    console.log('替换');
    
  }
  // 生成大纲
  createAiOutline() {
    console.log('生成大纲');
    
  }
  // 生成文章
  createAiArticle() {
    console.log('生成文章');
  }
}
const aiDialog = new AiDialogClass()
var aiEditMenu = Cherry.createMenuHook('AI', {
  subMenuConfig: [
    {
      noIcon: true,
      name: '优化',
      onclick: () => {
        console.log(' AI优化', cherry.getMarkdown());
        aiDialog.sendRequest('creation_optimization', cherry.getMarkdown())
      }
    },
    {
      noIcon: true,
      name: '纠错',
      onclick: () => {
        aiDialog.sendRequest('creation_proofreading', cherry.getMarkdown())
      }
    },
    {
      noIcon: true,
      name: '续写',
      onclick: () => {
        aiDialog.sendRequest('creation_continuation', cherry.getMarkdown())
      }
    },
    {
      noIcon: true,
      name: '翻译',
      onclick: () => {
        aiDialog.sendRequest('creation_translation', cherry.getMarkdown())
      }
    },
    {
      noIcon: true,
      name: '总结',
      onclick: () => {
        aiDialog.sendRequest('creation_summarize', cherry.getMarkdown())
      }
    },
    {
      noIcon: true,
      name: '大纲',
      onclick: () => {
        aiDialog.sendRequest('creation_leader', cherry.getMarkdown())
      }
    }
  ]
})
const saveData = () => {
  //console.log(markdownTitle)
  if (!markdownTitle || markdownTitle == "") {
    markdownTitle = window.prompt("请输入文稿标题");
  }
  const postData = { title: markdownTitle, content: cherry.getMarkdown() }
  window.parent.postMessage({ type: 'exportMd', data: JSON.stringify(postData) }, '*')
}
var saveMenu = Cherry.createMenuHook('保存', {
  onClick: function () {
    //console.log(selection)
    saveData()
  }
});
/**
 * 定义带图表表格的按钮
 */
var customMenuTable = Cherry.createMenuHook('图表', {
  iconName: 'trendingUp',
  subMenuConfig: [
    { noIcon: true, name: '折线图', onclick: (event) => { cherry.insert('\n| :line:{x,y} | Header1 | Header2 | Header3 | Header4 |\n| ------ | ------ | ------ | ------ | ------ |\n| Sample1 | 11 | 11 | 4 | 33 |\n| Sample2 | 112 | 111 | 22 | 222 |\n| Sample3 | 333 | 142 | 311 | 11 |\n'); } },
    { noIcon: true, name: '柱状图', onclick: (event) => { cherry.insert('\n| :bar:{x,y} | Header1 | Header2 | Header3 | Header4 |\n| ------ | ------ | ------ | ------ | ------ |\n| Sample1 | 11 | 11 | 4 | 33 |\n| Sample2 | 112 | 111 | 22 | 222 |\n| Sample3 | 333 | 142 | 311 | 11 |\n'); } },
  ]
});

var basicConfig = {
  id: 'markdown',
  externals: {
    echarts: window.echarts,
    katex: window.katex,
    MathJax: window.MathJax,
  },
  isPreviewOnly: false,
  engine: {
    global: {
      urlProcessor(url, srcType) {
        console.log(`url-processor`, url, srcType);
        return url;
      },
    },
    syntax: {
      image: {
        videoWrapper: (link, type, defaultWrapper) => {
          console.log(type);
          return defaultWrapper;
        },
      },
      autoLink: {
        /** 生成的<a>标签追加target属性的默认值 空：在<a>标签里不会追加target属性， _blank：在<a>标签里追加target="_blank"属性 */
        target: '',
        /** 生成的<a>标签追加rel属性的默认值 空：在<a>标签里不会追加rel属性， nofollow：在<a>标签里追加rel="nofollow：在"属性*/
        rel: '',
        /** 是否开启短链接 */
        enableShortLink: true,
        /** 短链接长度 */
        shortLinkLength: 20,
      },
      codeBlock: {
        theme: 'twilight',
        wrap: true, // 超出长度是否换行，false则显示滚动条
        lineNumber: true, // 默认显示行号
        copyCode: true, // 是否显示“复制”按钮
        editCode: true, // 是否显示“编辑”按钮
        changeLang: true, // 是否显示“切换语言”按钮
        expandCode: false, // 是否展开/收起代码块，当代码块行数大于10行时，会自动收起代码块
        selfClosing: true, // 自动闭合，为true时，当md中有奇数个```时，会自动在md末尾追加一个```
      },
      table: {
        enableChart: true,
      },
      fontEmphasis: {
        allowWhitespace: false, // 是否允许首尾空格
      },
      strikethrough: {
        needWhitespace: false, // 是否必须有前后空格
      },
      mathBlock: {
        engine: 'katex', // katex或MathJax
        src: '/markdown/katex/katex.min.js'
        //src: 'https://cdn.jsdelivr.net/npm/mathjax@3/es5/tex-svg.js', // 如果使用MathJax plugins，则需要使用该url通过script标签引入
      },
      inlineMath: {
        engine: 'katex', // katex或MathJax
      },
      // emoji: {
      //   useUnicode: true,
      //   customResourceURL: 'https://github.githubassets.com/images/icons/emoji/unicode/${code}.png?v8',
      //   upperCase: false,
      // },
      // htmlBlock: {
      //   filterStyle: true,
      // }
      toc: {
        tocStyle: 'nested'
      }
      // 'header': {
      //   strict: false
      // }
    },
    customSyntax: {
      // SyntaxHookClass
      CustomHook: {
        syntaxClass: CustomHookA,
        force: false,
        after: 'br',
      },
    },
  },
  multipleFileSelection: {
    video: true,
    audio: false,
    image: true,
    word: false,
    pdf: true,
    file: true,
  },
  toolbars: {
    toolbar: [
      'undo', 'redo', '|',
      // 把字体样式类按钮都放在加粗按钮下面
      { bold: ['bold', 'italic', 'underline', 'strikethrough', 'sub', 'sup', 'ruby', 'customMenuAName'] },
      // 'bold',
      // 'italic',
      // {
      //   strikethrough: ['strikethrough', 'underline', 'sub', 'sup', 'ruby'],
      // },
      'size',
      '|',
      'color',
      'header',
      {
        ol: ['ol',
          'ul',
          'checklist',
          'panel']
      },

      'justify',
      //'formula',
      {
        insert: ['drawIo', 'image', 'audio', 'video', 'link', 'ruby', 'detail', 'hr', 'br', 'code', 'inlineCode', 'formula', 'toc', 'table', 'pdf', 'word', 'file'],
      },
      'graph',
      'customMenuTable',
      //'togglePreview',
      //'codeTheme',
      'fullScreen',
      'search',
      'togglePreview'
    ],
    toolbarRight: ['saveMenu', '|', 'exportDataHook', 'changeLocale', 'wordCount', 'aiEditMenu'],
    bubble: ['bold', 'italic', 'underline', 'strikethrough', 'sub', 'sup', 'quote', 'ruby', '|', 'size', 'color'], // array or false
    //sidebar: ['mobilePreview', 'copy', 'theme', 'publish'],
    sidebar: ['mobilePreview', 'copy', 'shortcutKey', 'theme'],
    toc: {
      updateLocationHash: true, // 要不要更新URL的hash
      defaultModel: 'pure', // pure: 精简模式/缩略模式，只有一排小点； full: 完整模式，会展示所有标题
    },
    //toc: false, 
    customMenu: {
      customMenuAName: customMenuA,
      customMenuBName: customMenuB,
      saveMenu,
      customMenuTable,
      exportDataHook,
      aiEditMenu
    },
    shortcutKeySettings: {
      /** 是否替换已有的快捷键, true: 替换默认快捷键； false： 会追加到默认快捷键里，相同的shortcutKey会覆盖默认的 */
      isReplace: false,
      shortcutKeyMap: {
        'Control-S': {
          hookName: 'saveMenu',
          aliasName: '保存',
        },
        'Alt-Digit1': {
          hookName: 'header',
          aliasName: '标题',
        },
        'Control-Shift-KeyX': {
          hookName: 'bold',
          aliasName: '加粗',
        },
      },
    },
  },
  drawioIframeUrl: './drawio.html',
  previewer: {
    // 自定义markdown预览区域class
    // className: 'markdown'
    floatWhenClosePreviewer: true,
  },
  keydown: [],
  //extensions: [],
  callback: {
    changeString2Pinyin: pinyin,
    onClickPreview: (event) => {
      console.log("onClickPreview", event);
    },
  },
  editor: {
    id: 'cherry-text',
    name: 'cherry-text',
    autoSave2Textarea: true,
    //defaultModel: 'editOnly',
    defaultModel: 'edit&preview',
    showFullWidthMark: true, // 是否高亮全角符号 ·|￥|、|：|“|”|【|】|（|）|《|》
    showSuggestList: true, // 是否显示联想框
  },
  // cherry初始化后是否检查 location.hash 尝试滚动到对应位置
  autoScrollByHashAfterInit: true,
  // locale: 'en_US',
  themeSettings: {
    mainTheme: 'light',
  },
};
var config = Object.assign({}, basicConfig, { value: "" });
window.cherry = new Cherry(config);
// fetch('./markdown/basic.md').then((response) => response.text()).then((value) => {
//   var config = Object.assign({}, basicConfig, { value: value });
//   window.cherry = new Cherry(config);
// });

const debouncedHandleKeyDown = (event) => {
  // 确保仅在我们的按钮获得焦点时处理快捷键
  if (
    (event.metaKey || event.ctrlKey) &&
    event.key.toLowerCase() === "s"
  ) {
    event.stopPropagation(); // 先阻止事件冒泡
    event.preventDefault(); // 再阻止默认行为
    saveData();
  }
};
function isBase64(str) {
  if (str === '' || str.trim() === '') {
    return false;
  }
  try {
    return btoa(atob(str)) == str;
  } catch (err) {
    return false;
  }
}
function decodeBase64(base64String) {
  // 将Base64字符串分成每64个字符一组
  const padding = (base64String.length % 4) === 0 ? 0 : 4 - (base64String.length % 4);
  base64String += '='.repeat(padding);

  // 使用atob()函数解码Base64字符串
  const binaryString = atob(base64String);

  // 将二进制字符串转换为TypedArray
  const bytes = new Uint8Array(binaryString.length);
  for (let i = 0; i < binaryString.length; i++) {
      bytes[i] = binaryString.charCodeAt(i);
  }

  // 将TypedArray转换为字符串
  return new TextDecoder('utf-8').decode(bytes);
}
const eventHandler = (e) => {
  const eventData = e.data
  if (eventData.type === 'start') {
    markdownTitle = eventData.title || '未命名文稿'
    return
  }
  // console.log(markdownTitle)
  if (eventData.type === 'init') {
    const data = eventData.data
    markdownTitle = data.title || '未命名文稿'
    if (!data) {
      return;
    }
    let content = data.content;
    
    if (isBase64(content)) {
      content = decodeBase64(content);
    } else if (content instanceof ArrayBuffer) {
      content = new TextDecoder('utf-8').decode(content);
    }
    
    cherry.setMarkdown(content);
  }
  if (eventData.type == 'aiReciver') {
    aiDialog.showDialog(eventData.action, eventData.data)
  }
}
window.addEventListener('load', () => {
  window.parent.postMessage({ type: 'initSuccess' }, '*')
  window.addEventListener('message', eventHandler)
  document.addEventListener("keydown", debouncedHandleKeyDown);
})
window.addEventListener('unload', () => {
  window.removeEventListener('message', eventHandler)
  document.removeEventListener("keydown", debouncedHandleKeyDown);
})
