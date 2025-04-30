export function readOnlyInput(root: HTMLElement | null): void {
    if (root) {
        const nodes = root.querySelectorAll<HTMLInputElement>('input');
        for (const node of nodes) {
            if (!node.readOnly) {
                node.readOnly = true;
            }
        }
        if (document.activeElement) {
            (document.activeElement as HTMLElement).blur();
        }
    }
}

export async function getUrl(src: string | Blob | ArrayBuffer | Response): Promise<string> {
    if (typeof src === 'string') {
        return src;
    } else if (src instanceof Blob) {
        return URL.createObjectURL(src);
    } else if (src instanceof ArrayBuffer) {
        return URL.createObjectURL(new Blob([src]));
    } else if (src instanceof Response) {
        const blob = await src.blob();
        return URL.createObjectURL(blob);
    } else {
        throw new Error('Unsupported source type');
    }
}
export async function download(filename:string, data:any){
    if(!data){
        return; 
    }
   if (data instanceof ArrayBuffer) {
       data = new Blob([data]);
    }
    downloadFile(filename, URL.createObjectURL(data));
}

export function downloadFile(filename:string, href:string){
    let eleLink = document.createElement('a');
    eleLink.download = filename;
    eleLink.style.display = 'none';
    eleLink.href = href;
    document.body.appendChild(eleLink);
    eleLink.click();
    document.body.removeChild(eleLink);
}