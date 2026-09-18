import games from '@/api/games'
import collections from '@/api/collections'
import decks from '@/api/decks'
import cards from '@/api/cards'
import system from '@/api/system'
import replace from '@/api/replace'
import { useToast } from 'vue-toastification'

const { fetch: originalFetch } = window

window.fetch = (...args) =>
  originalFetch(...args).then(response => {
    const toast = useToast()
    if (!response.ok) {
      response.json().then(res => toast.error(res?.error?.message || 'Unknown error'))
    }
    return response
  })

export default {
  games,
  collections,
  decks,
  cards,
  system,
  replace,
}
