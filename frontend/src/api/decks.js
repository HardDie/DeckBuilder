import { List, ListAllUnique, Create, Read, Update, Delete } from '../../wailsjs/go/deck/Deck'
import { withBindingError, writeRequestFromBody } from '@/api/wails'

export default {
  list(requestData) {
    const config = requestData.config || {}
    return withBindingError(
      List(
        requestData.gameId || '',
        requestData.collectionId || '',
        config.sort || '',
        config.search || '',
      ),
    )
  },
  read(requestData) {
    return withBindingError(
      Read(
        requestData.gameId || '',
        requestData.collectionId || '',
        requestData.deckId || requestData.id || '',
      ),
    )
  },
  create(requestData) {
    return withBindingError(
      writeRequestFromBody(requestData.body).then(req =>
        Create(requestData.gameId || '', requestData.collectionId || '', req),
      ),
    )
  },
  update(requestData) {
    return withBindingError(
      writeRequestFromBody(requestData.body).then(req =>
        Update(
          requestData.gameId || '',
          requestData.collectionId || '',
          requestData.deckId || requestData.id || '',
          req,
        ),
      ),
    )
  },
  delete(requestData) {
    return withBindingError(
      Delete(
        requestData.gameId || '',
        requestData.collectionId || '',
        requestData.deckId || requestData.id || '',
      ),
    )
  },
  deckSuggestions(requestData) {
    return withBindingError(ListAllUnique(requestData.gameId || ''))
  },
}
