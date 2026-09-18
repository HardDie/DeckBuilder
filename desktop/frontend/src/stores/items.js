import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '@/api'
import { useStore } from '@/stores/main'

export const useItemsStore = defineStore('items', () => {
  const mainStore = useStore()

  const items = ref([])
  const itemsMeta = ref({})
  const isItemsLoading = ref(false)
  const isApiPending = ref(false)

  function fetchItems(requestData) {
    isItemsLoading.value = true
    return api[mainStore.itemType]
      .list({
        ...requestData,
        config: { sort: mainStore.sort, search: mainStore.search },
      })
      .then(response => {
        items.value = response.data
        itemsMeta.value = response.meta
      })
      .finally(() => {
        isItemsLoading.value = false
      })
  }

  function fetchItem(requestData) {
    return api[mainStore.itemType].read(requestData).then(response => response.data)
  }

  const addItem = payload => {
    items.value = items.value.concat([payload])
    if (!itemsMeta.value.total) {
      itemsMeta.value.total = 0
    }
    itemsMeta.value.total += 1

    if (payload.count) {
      if (!itemsMeta.value.cardsTotal) {
        itemsMeta.value.cardsTotal = 0
      }
      itemsMeta.value.cardsTotal += payload.count
    }
  }

  const replaceItem = payload => {
    items.value = items.value.map(item => {
      if (item.id === payload.oldId) {
        return payload
      }
      return item
    })
  }

  function clearItems() {
    items.value = []
    itemsMeta.value = {}
  }

  const deleteItem = payload => {
    const deleted = items.value.find(el => String(el.id) === payload.id)

    if (!deleted) {
      return
    }

    items.value = items.value.filter(el => el.id !== deleted.id)

    if (itemsMeta.value?.total) {
      itemsMeta.value.total -= 1
    }

    if (itemsMeta.value?.cardsTotal && deleted.count) {
      itemsMeta.value.cardsTotal -= deleted.count
    }
  }

  function fetchCreateItem(requestData) {
    isApiPending.value = true
    return api[mainStore.itemType]
      .create(requestData)
      .then(response => addItem(response.data))
      .finally(() => {
        isApiPending.value = false
      })
  }

  function fetchEditItem(requestData) {
    isApiPending.value = true
    return api[mainStore.itemType]
      .update(requestData)
      .then(response => replaceItem({ ...response.data, oldId: requestData.id }))
      .finally(() => {
        isApiPending.value = false
      })
  }

  function fetchDeleteItem(requestData) {
    return api[mainStore.itemType].delete(requestData).then(() => deleteItem(requestData))
  }

  function fetchImportGame(requestData) {
    isApiPending.value = true
    return api.games
      .import(requestData)
      .then(response => addItem(response.data))
      .finally(() => {
        isApiPending.value = false
      })
  }

  function fetchExportGame(requestData) {
    isApiPending.value = true
    return api.games
      .export(requestData)
      .then(response => {
        response.blob().then(data => {
          const a = document.createElement('a')
          a.href = URL.createObjectURL(data)
          a.download = `${requestData.gameId}.zip`
          a.click()
        })
      })
      .finally(() => (isApiPending.value = false))
  }

  function fetchDuplicateGame(requestData) {
    isApiPending.value = true
    return api.games
      .duplicate(requestData)
      .then(response => addItem(response.data))
      .finally(() => (isApiPending.value = false))
  }

  function fetchGenerateGame(requestData) {
    return api.games.generate({ ...requestData, body: { sortOrder: mainStore.sort } })
  }

  return {
    items,
    itemsMeta,
    isItemsLoading,
    clearItems,
    fetchItems,
    fetchItem,
    fetchCreateItem,
    fetchEditItem,
    fetchDeleteItem,

    fetchImportGame,
    fetchExportGame,
    fetchDuplicateGame,
    fetchGenerateGame,
    isApiPending,
  }
})
