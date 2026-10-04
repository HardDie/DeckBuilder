// Fixture catalog for `make screenshots`. Loaded only when the page opens with ?screenshot=1.
// Pictures are the Four Souls files under docs/wiki/images/foursouls, served by Vite.

const asset = name => `/screenshot-assets/${name}`

const games = [
  {
    id: 'four_souls',
    name: 'Four Souls',
    description: '',
    cachedImage: asset('game-cover.png'),
  },
]

const collections = [
  {
    id: 'base_game',
    name: 'Base Game V2',
    description: '',
    cachedImage: asset('collection-base.png'),
  },
  {
    id: 'four_souls_plus',
    name: 'Four Souls+ V2',
    description: '',
    cachedImage: asset('collection-plus.png'),
  },
  {
    id: 'requiem',
    name: 'Requiem',
    description: '',
    cachedImage: asset('collection-requiem.png'),
  },
]

const decks = [
  {
    id: 'character',
    name: 'Character',
    description: '',
    cachedImage: asset('deck-character.png'),
  },
  {
    id: 'monster',
    name: 'Monster',
    description: '',
    cachedImage: asset('deck-monster.png'),
  },
  {
    id: 'treasure',
    name: 'Treasure',
    description: '',
    cachedImage: asset('deck-treasure.png'),
  },
  {
    id: 'loot',
    name: 'Loot',
    description: '',
    cachedImage: asset('deck-loot.png'),
  },
]

const cards = [
  {
    id: 1,
    name: 'Isaac',
    description: '',
    count: 1,
    cachedImage: asset('card-isaac.png'),
  },
  {
    id: 2,
    name: 'Maggy',
    description: '',
    count: 2,
    cachedImage: asset('card-maggy.png'),
  },
  {
    id: 3,
    name: 'Cain',
    description: '',
    count: 1,
    cachedImage: asset('card-cain.png'),
  },
]

function listOf(items, cardsTotal) {
  const meta = { total: items.length }
  if (cardsTotal != null) {
    meta.cardsTotal = cardsTotal
  }
  return Promise.resolve({ data: items, meta })
}

function one(items, id) {
  const found = items.find(item => String(item.id) === String(id))
  if (!found) {
    return Promise.reject(new Error(`screenshot fixture has no ${id}`))
  }
  return Promise.resolve({ data: found })
}

const gamesApi = {
  list() {
    return listOf(games)
  },
  read(requestData) {
    return one(games, requestData.gameId || requestData.id)
  },
}

const collectionsApi = {
  list() {
    return listOf(collections)
  },
  read(requestData) {
    return one(collections, requestData.collectionId || requestData.id)
  },
}

const decksApi = {
  list() {
    return listOf(decks)
  },
  read(requestData) {
    return one(decks, requestData.deckId || requestData.id)
  },
}

const cardsApi = {
  list() {
    const cardsTotal = cards.reduce((sum, card) => sum + (card.count || 0), 0)
    return listOf(cards, cardsTotal)
  },
  read(requestData) {
    return one(cards, requestData.cardId || requestData.id)
  },
}

const idle = () => Promise.resolve({ data: null })

export default {
  games: {
    ...gamesApi,
    create: idle,
    update: idle,
    delete: idle,
    export: idle,
    import: idle,
    duplicate: idle,
    generate: idle,
  },
  collections: { ...collectionsApi, create: idle, update: idle, delete: idle },
  decks: { ...decksApi, create: idle, update: idle, delete: idle },
  cards: { ...cardsApi, create: idle, update: idle, delete: idle },
  system: {
    getSettings: () =>
      Promise.resolve({ data: { lang: 'en', enable_back_shadow: false, card_scale: 1 } }),
    updateSettings: idle,
    checkStatus: () => Promise.resolve({ data: { status: 'empty', progress: 0 } }),
    getVersion: () => Promise.resolve({ data: 'dev' }),
  },
  search: { list: idle },
  replace: { prepare: idle, replace: idle },
}
