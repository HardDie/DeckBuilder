import Vue from 'vue'
import VueRouter from 'vue-router'
import Games from '@/views/Games'
import Game from '@/views/Game'

Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    name: 'Games',
    component: Games
  },
  {
    path: '/games',
    name: 'Games',
    component: Games
  },
  {
    path: '/games/:name',
    name: 'Game',
    component: Game
  }
]

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes
})

export default router
