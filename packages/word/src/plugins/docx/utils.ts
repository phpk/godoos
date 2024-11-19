export function saveAs(blob: Blob, name: string) {
  const a = document.createElement('a')
  a.href = window.URL.createObjectURL(blob)
  a.download = name
  a.click()
  window.URL.revokeObjectURL(a.href)
}
