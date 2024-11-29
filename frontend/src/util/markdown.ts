import 'cherry-markdown/dist/cherry-markdown.css';
import CherryEngine from 'cherry-markdown/dist/cherry-markdown.engine.core';

export function renderMarkdown(currentText: string) {
    //console.log(currentText)
    const cherryEngineInstance:any = new CherryEngine({});
    return cherryEngineInstance.makeHtml(currentText);
    // const currentCherry = new Cherry(cherryConfig);
    // currentCherry.setMarkdown(currentText);
    // console.log(currentCherry.getHtml())
    // return currentCherry.getHtml();
}