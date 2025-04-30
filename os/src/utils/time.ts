export function formatTimetoMD(timestamp: string | number) {
	const date = new Date(Number(timestamp)); // 确保时间戳为数字
	const month = date.getMonth() + 1; // 月份是从0开始的
	const day = date.getDate();
	const hours = date.getHours();
	const minutes = date.getMinutes();

	// 格式化小时和分钟为两位数
	const formattedHours = hours.toString().padStart(2, "0");
	const formattedMinutes = minutes.toString().padStart(2, "0");

	return `${month}-${day} ${formattedHours}:${formattedMinutes}`;
}
export function formatTimetoYMD(timestamp: number) {
    const date = new Date(timestamp)
    const year = date.getFullYear();
    const month = (date.getMonth() + 1).toString().padStart(2, '0');
    const day = date.getDate().toString().padStart(2, '0');
    const hours = date.getHours().toString().padStart(2, '0');
    const minutes = date.getMinutes().toString().padStart(2, '0');
    const seconds = date.getSeconds().toString().padStart(2, '0');

    return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`
}