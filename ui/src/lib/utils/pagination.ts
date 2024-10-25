export function pageToOffset(page: number, pageSize: number): number {
  return (page - 1) * pageSize;
}
