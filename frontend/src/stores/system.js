import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '@/api'

export const useSystemStore = defineStore('system', () => {
  const version = ref('')

  function fetchCheckStatus() {
    return api.system.checkStatus().then(response => response.data)
  }

  function fetchDownloadStatus() {
    return api.system.downloadStatus().then(response => response.data)
  }

  function fetchVersion() {
    return api.system.getVersion().then(response => {
      version.value = response.data
    })
  }

  function fetchSettings() {
    return api.system.getSettings().then(response => response.data)
  }

  function saveSettings(settings) {
    return api.system.updateSettings(settings).then(response => response.data)
  }

  return {
    fetchCheckStatus,
    fetchDownloadStatus,
    fetchVersion,
    fetchSettings,
    saveSettings,
    version,
  }
})
