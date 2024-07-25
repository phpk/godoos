export async function OpenDirDialog(){
    if((window as any).go) {
        return (window as any)['go']['app']['App']['OpenDirDialog']();
    }else {
        return ""
    }
}

export function RestartApp(){
    if(!(window as any).go){
        window.location.reload();
    }else{
        return (window as any)['go']['app']['App']['RestartApp']();
    }
   
}
