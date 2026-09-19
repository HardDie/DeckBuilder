import {
  ListGames,
  CreateGame,
  ReadGame,
  UpdateGame,
  DeleteGame,
  DuplicateGame,
} from '../../wailsjs/go/main/App'
import { useToast } from 'vue-toastification'

function toastBindingError(err) {
  const toast = useToast()
  const message = typeof err === 'string' ? err : err?.message || 'Unknown error'
  toast.error(message)
}

function withBindingError(promise) {
  return promise.catch(err => {
    toastBindingError(err)
    throw err
  })
}

async function imageFileBytes(file) {
  if (!(file instanceof Blob)) {
    return []
  }
  return Array.from(new Uint8Array(await file.arrayBuffer()))
}

async function gameWriteRequestFromBody(body) {
  if (body instanceof FormData) {
    const file = body.get('imageFile')
    const hasFile = file instanceof Blob
    return {
      name: body.get('name') || '',
      description: body.get('description') || '',
      image: hasFile ? '' : body.get('image') || '',
      imageFile: hasFile ? await imageFileBytes(file) : [],
    }
  }
  return {
    name: body?.name || '',
    description: body?.description || '',
    image: body?.image || '',
    imageFile: body?.imageFile || [],
  }
}

export default {
  list(requestData) {
    const config = requestData.config || {}
    return withBindingError(ListGames(config.sort || '', config.search || ''))
  },
  read(requestData) {
    return withBindingError(ReadGame(requestData.gameId || ''))
  },
  create(requestData) {
    return withBindingError(gameWriteRequestFromBody(requestData.body).then(req => CreateGame(req)))
  },
  update(requestData) {
    return withBindingError(
      gameWriteRequestFromBody(requestData.body).then(req =>
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
