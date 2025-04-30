    function HexToRgb(str: string): number[] {
        str = str.replace('#', '');
        const hxs = str.match(/../g);
        if (!hxs || hxs.length !== 3) {
            throw new Error('Invalid hex color format');
        }
        const rgb: number[] = hxs.map(hex => parseInt(hex, 16));
        return rgb;
    }
    
    function RgbToHex(a: number, b: number, c: number): string {
        const hexs = [a.toString(16), b.toString(16), c.toString(16)];
        for (let i = 0; i < 3; i++) if (hexs[i].length === 1) hexs[i] = '0' + hexs[i];
        return '#' + hexs.join('');
    }
    
    export function getDarkColor(color: string, level: number): string {
        const rgbc = HexToRgb(color);
        for (let i = 0; i < 3; i++) rgbc[i] = Math.floor(rgbc[i] * (1 - level));
        return RgbToHex(rgbc[0], rgbc[1], rgbc[2]);
    }
    
    export function getLightColor(color: string, level: number): string {
        const rgbc = HexToRgb(color);
        for (let i = 0; i < 3; i++) rgbc[i] = Math.floor((255 - rgbc[i]) * level + rgbc[i]);
        return RgbToHex(rgbc[0], rgbc[1], rgbc[2]);
    }
    export const themeColor = [
        '#FFFFFF',
        '#000000',
        '#BFBFBF',
        '#323232',
        '#4472C4',
        '#ED7D31',
        '#A5A5A5',
        '#FFC000',
        '#5B9BD5',
        '#71AD47'
    ];
    
    export const indexedColor = [
        '#000000',
        '#FFFFFF',
        '#FF0000',
        '#00FF00',
        '#0000FF',
        '#FFFF00',
        '#FF00FF',
        '#00FFFF',
        '#000000',
        '#FFFFFF',
        '#FF0000',
        '#00FF00',
        '#0000FF',
        '#FFFF00',
        '#FF00FF',
        '#00FFFF',
        '#800000',
        '#008000',
        '#000080',
        '#808000',
        '#800080',
        '#008080',
        '#C0C0C0',
        '#808080',
        '#9999FF',
        '#993366',
        '#FFFFCC',
        '#CCFFFF',
        '#660066',
        '#FF8080',
        '#0066CC',
        '#CCCCFF',
        '#000080',
        '#FF00FF',
        '#FFFF00',
        '#00FFFF',
        '#800080',
        '#800000',
        '#008080',
        '#0000FF',
        '#00CCFF',
        '#CCFFFF',
        '#CCFFCC',
        '#FFFF99',
        '#99CCFF',
        '#FF99CC',
        '#CC99FF',
        '#FFCC99',
        '#3366FF',
        '#33CCCC',
        '#99CC00',
        '#FFCC00',
        '#FF9900',
        '#FF6600',
        '#666699',
        '#969696',
        '#003366',
        '#339966',
        '#003300',
        '#333300',
        '#993300',
        '#993366',
        '#333399',
        '#333333',
        '#FFFFFF'
    ];