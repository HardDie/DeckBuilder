import { Collection, Game, Root } from '../../wailsjs/go/search/Search'
import { withBindingError } from '@/api/wails'

export default {
  root(requestData) {
    const config = requestData?.config || {}
    return withBindingError(Root(config.sort || '', config.search || ''))
  },
  game(requestData) {
    const config = requestData?.config || {}
    return withBindingError(Game(requestData?.gameId || '', config.sort || '', config.search || ''))
  },
  collection(requestData) {
    const config = requestData?.config || {}
    return withBindingError(
      Collection(
        requestData?.gameId || '',
        requestData?.collectionId || '',
        config.sort || '',
        config.search || '',
      ),
    )
  },
}
