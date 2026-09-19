import { List, Create, Read, Update, Delete, Duplicate } from '../../wailsjs/go/game/Game'
import { withBindingError, writeRequestFromBody } from '@/api/wails'

export default {
  list(requestData) {
    const config = requestData.config || {}
    return withBindingError(List(config.sort || '', config.search || ''))
  },
  read(requestData) {
    return withBindingError(Read(requestData.gameId || ''))
  },
  create(requestData) {
    return withBindingError(writeRequestFromBody(requestData.body).then(req => Create(req)))
  },
  update(requestData) {
    return withBindingError(
      writeRequestFromBody(requestData.body).then(req =>
        Update(requestData.gameId || requestData.id || '', req),
      ),
    )
  },
  delete(requestData) {
    return withBindingError(Delete(requestData.gameId || requestData.id || ''))
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
      Duplicate(requestData.gameId || requestData.id || '', requestData.body?.name || ''),
    )
  },
  generate(requestData) {
    return fetch(`/api/games/${requestData.gameId}/generate`, {
      method: 'POST',
      body: JSON.stringify(requestData.body),
    })
  },
}
