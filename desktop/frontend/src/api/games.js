export default {
  list(requestData) {
    return fetch(`/api/games?${new URLSearchParams(requestData.config)}`).then(response =>
      response.json(),
    )
  },
  read(requestData) {
    return fetch(`/api/games/${requestData.gameId}`).then(response => response.json())
  },
  create(requestData) {
    return fetch('/api/games', {
      method: 'POST',
      body: requestData.body,
    }).then(response => response.json())
  },
  update(requestData) {
    return fetch(`/api/games/${requestData.gameId}`, {
      method: 'PATCH',
      body: requestData.body,
    }).then(response => response.json())
  },
  delete(requestData) {
    return fetch(`/api/games/${requestData.gameId}`, { method: 'DELETE' })
  },
  export(requestData) {
    return fetch(`/api/games/${requestData.gameId}/export`)
  },
  import(requestData) {
    return fetch('/api/games/import', {
      method: 'POST',
      body: requestData,
    }).then(response => response.json())
  },
  duplicate(requestData) {
    return fetch(`/api/games/${requestData.gameId}/duplicate`, {
      method: 'POST',
      body: JSON.stringify(requestData.body),
    }).then(response => response.json())
  },
  generate(requestData) {
    return fetch(`/api/games/${requestData.gameId}/generate`, {
      method: 'POST',
      body: JSON.stringify(requestData.body),
    })
  },
}
