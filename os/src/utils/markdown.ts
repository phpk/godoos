import 'cherry-markdown/dist/cherry-markdown.css';
import CherryEngine from 'cherry-markdown/dist/cherry-markdown.engine.core';

export function renderMarkdown(currentText: string) {
    currentText = currentText.replace("<think>", '> ');
    currentText = currentText.replace("</think>", '\n');
    //console.log('Modified Text:', currentText);
    const cherryEngineInstance:any = new CherryEngine({});
    return cherryEngineInstance.makeHtml(currentText);
}