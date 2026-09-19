import {
  ListGames,
  CreateGame,
  ReadGame,
  UpdateGame,
  DeleteGame,
  DuplicateGame,
} from '../../wailsjs/go/main/App'
import { withBindingError, writeRequestFromBody } from '@/api/wails'

export default {
  list(requestData) {
    const config = requestData.config || {}
    return withBindingError(ListGames(config.sort || '', config.search || ''))
  },
  read(requestData) {
    return withBindingError(ReadGame(requestData.gameId || ''))
  },
  create(requestData) {
    return withBindingError(writeRequestFromBody(requestData.body).then(req => CreateGame(req)))
  },
  update(requestData) {
    return withBindingError(
      writeRequestFromBody(requestData.body).then(req =>
        UpdateGame(requestData.gameId || requestData.id || '', req),
      ),
    )
  },
  delete(requestData) {
    return withBindingError(DeleteGame(requestData.gameId || requestData.id || ''))
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
    return withBindingError(
      DuplicateGame(requestData.gameId || requestData.id || '', requestData.body?.name || ''),
    )
  },
  generate(requestData) {
    return fetch(`/api/games/${requestData.gameId}/generate`, {
      method: 'POST',
      body: JSON.stringify(requestData.body),
    })
  },
}
