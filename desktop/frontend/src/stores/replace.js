import { defineStore } from 'pinia'
import api from '@/api'

export const useReplaceStore = defineStore('replace', () => {
  function prepare(data) {
    return api.replace.prepare(data)
  }

  function replace(data) {
    return api.replace.replace(data).then(response => {
      const jsonBlob = new Blob([JSON.stringify(response.data)], { type: 'application/json' })
      const a = document.createElement('a')
      a.href = URL.createObjectURL(jsonBlob)
      a.download = `result-${new Date().toLocaleDateString('en-US')}.json`
      a.click()
    })
  }

  return { prepare, replace }
})
