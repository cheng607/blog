export const formatTime = (time: string): string => {
    if (!time) return '';
    const date = new Date(time);
    if (isNaN(date.getTime())) return '无效时间';
    const padZero = (num: number) => num.toString().padStart(2, '0');

    const year = date.getFullYear();
    const month = padZero(date.getMonth() + 1);
    const day = padZero(date.getDate());
    const hours = padZero(date.getHours());
    const minutes = padZero(date.getMinutes());

    return `${year}年${month}月${day}日 ${hours}:${minutes}`;
};