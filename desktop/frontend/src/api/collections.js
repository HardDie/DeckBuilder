import { List, Create, Read, Update, Delete } from '../../wailsjs/go/collection/Collection'
import { withBindingError, writeRequestFromBody } from '@/api/wails'

export default {
  list(requestData) {
    const config = requestData.config || {}
    return withBindingError(List(requestData.gameId || '', config.sort || '', config.search || ''))
  },
  read(requestData) {
    return withBindingError(
      Read(requestData.gameId || '', requestData.collectionId || requestData.id || ''),
    )
  },
  create(requestData) {
    return withBindingError(
      writeRequestFromBody(requestData.body).then(req => Create(requestData.gameId || '', req)),
    )
  },
  update(requestData) {
    return withBindingError(
      writeRequestFromBody(requestData.body).then(req =>
        Update(requestData.gameId || '', requestData.collectionId || requestData.id || '', req),
      ),
    )
  },
  delete(requestData) {
    return withBindingError(
      Delete(requestData.gameId || '', requestData.collectionId || requestData.id || ''),
    )
  },
}
