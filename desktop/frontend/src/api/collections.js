export default {
  list(requestData) {
    return fetch(
      `/api/games/${requestData.gameId}/collections?${new URLSearchParams(requestData.config)}`,
    ).then(response => response.json())
  },
  read(requestData) {
    return fetch(`/api/games/${requestData.gameId}/collections/${requestData.collectionId}`).then(
      response => response.json(),
    )
  },
  create(requestData) {
    return fetch(`/api/games/${requestData.gameId}/collections`, {
      method: 'POST',
      body: requestData.body,
    }).then(response => response.json())
  },
  update(requestData) {
    return fetch(`/api/games/${requestData.gameId}/collections/${requestData.collectionId}`, {
      method: 'PATCH',
      body: requestData.body,
    }).then(response => response.json())
  },
  delete(requestData) {
    return fetch(`/api/games/${requestData.gameId}/collections/${requestData.collectionId}`, {
      method: 'DELETE',
    })
  },
}
