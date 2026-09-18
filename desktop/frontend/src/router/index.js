import { createRouter, createWebHistory } from 'vue-router'
import MainView from '@/views/MainView.vue'

import { setBar } from './hooks'

const routes = [
  {
    path: '/',
    name: 'games',
    beforeEnter: setBar,
    component: MainView,
  },
  {
    path: '/games/:gameId',
    name: 'collections',
    beforeEnter: setBar,
    component: MainView,
  },
  {
    path: '/games/:gameId/collections/:collectionId',
    name: 'decks',
    beforeEnter: setBar,
    component: MainView,
  },
  {
    path: '/games/:gameId/collections/:collectionId/decks/:deckId',
    name: 'cards',
    beforeEnter: setBar,
    component: MainView,
  },
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
})

export default router
