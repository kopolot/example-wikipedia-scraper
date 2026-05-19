export function unixTimestampToDate(timestamp: number): Date {
    return new Date(timestamp * 1000);
}