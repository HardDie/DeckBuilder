import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '@/api'

export const useSystemStore = defineStore('system', () => {
  const version = ref('')

  function fetchCheckStatus() {
    return api.system.checkStatus().then(response => response.data)
  }

  function fetchVersion() {
    return api.system.getVersion().then(response => {
      version.value = response.data
    })
  }

  return { fetchCheckStatus, fetchVersion, version }
})
