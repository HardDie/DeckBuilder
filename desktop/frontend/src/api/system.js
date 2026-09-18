export default {
  quit() {
    return fetch('/api/system/quit', { method: 'DELETE' })
  },
  checkStatus() {
    return fetch('/api/system/status').then(response => response.json())
  },
  getVersion() {
    return fetch('/api/system/version').then(response => response.json())
  },
}
