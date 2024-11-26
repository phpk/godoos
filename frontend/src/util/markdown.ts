import 'cherry-markdown/dist/cherry-markdown.css';
import Cherry from 'cherry-markdown';
const cherryConfig:any = {
    editor: {
        height: 'auto',
        defaultModel: 'previewOnly',
    },
    engine: {
        global: {
            // 开启流式模式 （默认 true）
            flowSessionContext: true,
        },
        syntax: {
            codeBlock: {
                selfClosing: false,
            },
            header: {
                anchorStyle: 'none',
            },
            table: {
                selfClosing: false,
            },
            fontEmphasis: {
                selfClosing: false,
            }
        }
    },
    previewer: {
        enablePreviewerBubble: false,
    },
    isPreviewOnly: true,
};

export function renderMarkdown(currentText: string) {
    const currentCherry = new Cherry(cherryConfig);
    currentCherry.setMarkdown(currentText);
    return currentCherry.getHtml();
}