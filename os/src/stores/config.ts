
const CONFIG_KEY = 'Godoos_systemConfig';
export const API_URL = 'http://localhost:8188';
export function getSystemConfig() {
    const systemConfig = localStorage.getItem(CONFIG_KEY) || '{}';
    const config = JSON.parse(systemConfig);
    if (!config.serviceType) config.serviceType = 'local';
    if (!config.serviceUrl || config.serviceType === 'local') config.serviceUrl = API_URL;
    return config;
}
export function setSystemKey(key: string, value: any) {
    const config = getSystemConfig();
    config[key] = value;
    if(key === 'serviceType' && value === 'local') config.serviceUrl = API_URL;
    localStorage.setItem(CONFIG_KEY, JSON.stringify(config));
}
export function setSystemConfig(config: any) {
    if(config.serviceType === 'local') config.serviceUrl = API_URL;
    localStorage.setItem(CONFIG_KEY, JSON.stringify(config));
}
export function getSystemKey(key: string) {
    const config = getSystemConfig();
    if(config.serviceType === 'local') config.serviceUrl = API_URL;
    return config[key];
}