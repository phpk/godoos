export async function OpenDirDialog(){
    if((window as any).go) {
        //(window as any).go.OpenDirDialog();
        return (window as any)['go']['main']['App']['OpenDirDialog']();
    }else {
        return ""
    }
}