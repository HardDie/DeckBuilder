import { List, Create, Read, Update, Delete } from '../../wailsjs/go/card/Card'
import { withBindingError, writeRequestFromBody } from '@/api/wails'

function cardID(requestData) {
  return Number(requestData.cardId || requestData.id || 0)
}

export default {
  list(requestData) {
    const config = requestData.config || {}
    return withBindingError(
      List(
        requestData.gameId || '',
        requestData.collectionId || '',
        requestData.deckId || '',
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
        requestData.deckId || '',
        cardID(requestData),
      ),
    )
  },
  create(requestData) {
    return withBindingError(
      writeRequestFromBody(requestData.body).then(req =>
        Create(
          requestData.gameId || '',
          requestData.collectionId || '',
          requestData.deckId || '',
          req,
        ),
      ),
    )
  },
  update(requestData) {
    return withBindingError(
      writeRequestFromBody(requestData.body).then(req =>
        Update(
          requestData.gameId || '',
          requestData.collectionId || '',
          requestData.deckId || '',
          cardID(requestData),
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
        requestData.deckId || '',
        cardID(requestData),
      ),
    )
  },
}
