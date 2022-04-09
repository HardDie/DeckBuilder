import api from '@/api'

export default {
  fetchGames ({ commit }, requestData) {
    return api.games.gameslist().then(response => commit('setGames', response.games))
  },
  fetchCreateGame ({ commit }, requestData) {
    return api.games.createGame(requestData)
  }
}
