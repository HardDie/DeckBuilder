const gamesEndpoint = 'http://localhost:5000/games'

export default {
  gameslist () {
    return fetch(gamesEndpoint).then(res => res.json())
  },
  createGame (requestData) {
    return fetch(gamesEndpoint, { method: 'POST', body: JSON.stringify(requestData) }).then(res => res.json())
  }
}
