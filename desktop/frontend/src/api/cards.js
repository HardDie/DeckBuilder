export default {
  list(requestData) {
    return fetch(
      `/api/games/${requestData.gameId}/collections/${requestData.collectionId}/decks/${
        requestData.deckId
      }/cards?${new URLSearchParams(requestData.config)}`,
    ).then(response => response.json())
  },
  read(requestData) {
    return fetch(
      `/api/games/${requestData.gameId}/collections/${requestData.collectionId}/decks/${requestData.deckId}/cards/${requestData.cardId}`,
    ).then(response => response.json())
  },
  create(requestData) {
    return fetch(
      `/api/games/${requestData.gameId}/collections/${requestData.collectionId}/decks/${requestData.deckId}/cards`,
      {
        method: 'POST',
        body: requestData.body,
      },
    ).then(response => response.json())
  },
  update(requestData) {
    return fetch(
      `/api/games/${requestData.gameId}/collections/${requestData.collectionId}/decks/${requestData.deckId}/cards/${requestData.cardId}`,
      {
        method: 'PATCH',
        body: requestData.body,
      },
    ).then(response => response.json())
  },
  delete(requestData) {
    return fetch(
      `/api/games/${requestData.gameId}/collections/${requestData.collectionId}/decks/${requestData.deckId}/cards/${requestData.cardId}`,
      { method: 'DELETE' },
    )
  },
}
