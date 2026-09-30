import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '@/api'

export const useStore = defineStore('main', () => {
  const itemType = ref('games')
  const sort = ref('name')
  const search = ref('')

  const breadcrumbPathItems = ref([])
  const pathItems = ref([])

  function setItemType(type) {
    itemType.value = type
  }

  function setSort(val) {
    sort.value = val
  }

  function setSearch(val) {
    search.value = val
  }

  function setPathItems(val) {
    pathItems.value = val
  }

  async function fetchPathItems(path) {
    const prepPath = path.split('/').slice(1)
    pathItems.value = prepPath
    const onlyIds = prepPath.filter((el, idx) => idx % 2 !== 0)
    breadcrumbPathItems.value = await Promise.all(
      onlyIds.map((pathEl, idx, arr) => {
        switch (idx) {
          case 0:
            return api.games.read({ gameId: arr[0] }).then(response => response.data.name)
          case 1:
            return api.collections
              .read({
                gameId: arr[0],
                collectionId: arr[1],
              })
              .then(response => response.data.name)
          case 2:
            return api.decks
              .read({
                gameId: arr[0],
                collectionId: arr[1],
                deckId: arr[2],
              })
              .then(response => response.data.name)
          case 3:
            return api.cards
              .read({
                gameId: arr[0],
                collectionId: arr[1],
                deckId: arr[2],
                cardId: arr[3],
              })
              .then(response => response.data.name)
          default:
            return 'unknown'
        }
      }),
    )
  }

  return {
    itemType,
    sort,
    setItemType,
    breadcrumbPathItems,
    pathItems,
    fetchPathItems,
    setSort,
    setPathItems,
    search,
    setSearch,
  }
})
