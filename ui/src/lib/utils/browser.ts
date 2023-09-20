export function getBrowserTheme(): 'light' | 'dark' {
  const isDark = window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches;
  if (isDark) {
    return 'dark';
  }
  return 'light';
}
