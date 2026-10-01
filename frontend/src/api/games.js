import {
  List,
  Create,
  Read,
  Update,
  Delete,
  Duplicate,
  Export,
  Import,
} from '../../wailsjs/go/game/Game'
import { Game as GenerateGame } from '../../wailsjs/go/generator/Generator'
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
    return withBindingError(Export(requestData.gameId || requestData.id || ''))
  },
  import(requestData) {
    return withBindingError(importZip(requestData).then(({ name, file }) => Import(name, file)))
  },
  duplicate(requestData) {
    return withBindingError(
      Duplicate(requestData.gameId || requestData.id || '', requestData.body?.name || ''),
    )
  },
  generate(requestData) {
    const body = requestData.body || {}
    return withBindingError(
      GenerateGame(requestData.gameId || '', body.sortOrder || '', Number(body.scale) || 0),
    )
  },
}

async function importZip(body) {
  const file = body?.file
  const name = body?.name || ''
  if (!(file instanceof Blob) || file.size === 0) {
    throw new Error('The file must be passed as an argument')
  }
  return {
    name,
    file: Array.from(new Uint8Array(await file.arrayBuffer())),
  }
}
