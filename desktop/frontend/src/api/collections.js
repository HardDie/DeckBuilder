import {
  ListCollections,
  CreateCollection,
  ReadCollection,
  UpdateCollection,
  DeleteCollection,
} from '../../wailsjs/go/main/App'
import { withBindingError, writeRequestFromBody } from '@/api/wails'

export default {
  list(requestData) {
    const config = requestData.config || {}
    return withBindingError(
      ListCollections(requestData.gameId || '', config.sort || '', config.search || ''),
    )
  },
  read(requestData) {
    return withBindingError(
      ReadCollection(requestData.gameId || '', requestData.collectionId || requestData.id || ''),
    )
  },
  create(requestData) {
    return withBindingError(
      writeRequestFromBody(requestData.body).then(req =>
        CreateCollection(requestData.gameId || '', req),
      ),
    )
  },
  update(requestData) {
    return withBindingError(
      writeRequestFromBody(requestData.body).then(req =>
        UpdateCollection(
          requestData.gameId || '',
          requestData.collectionId || requestData.id || '',
          req,
        ),
      ),
    )
  },
  delete(requestData) {
    return withBindingError(
      DeleteCollection(requestData.gameId || '', requestData.collectionId || requestData.id || ''),
    )
  },
}
