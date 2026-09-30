import games from '@/api/games'
import collections from '@/api/collections'
import decks from '@/api/decks'
import cards from '@/api/cards'
import system from '@/api/system'
import search from '@/api/search'
import replace from '@/api/replace'
import screenshotApi from '@/api/screenshot'
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

const live = {
  games,
  collections,
  decks,
  cards,
  system,
  search,
  replace,
}

// Chosen once at load. Client-side routing drops the query; the module stays stubbed.
const screenshot = new URLSearchParams(window.location.search).get('screenshot') === '1'

export default screenshot ? screenshotApi : live
