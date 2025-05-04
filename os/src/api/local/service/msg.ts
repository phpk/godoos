interface APIResponse {
    message?: string;
    code: number;
    success: boolean;
    data?: any;
    error?: string;
    time: number;
}

// 请求成功的时候使用该方法返回信息
export function success(msg: string, data: any): APIResponse {
    return {
        code: 0,
        success: true,
        message: msg,
        data: data,
        time: Math.floor(Date.now() / 1000),
    };
}


// 请求失败的时候，使用该方法返回信息
export function error(msg: string): APIResponse {
    return {
        code: -1,
        success: false,
        message: msg,
        time: Math.floor(Date.now() / 1000),
    };
}