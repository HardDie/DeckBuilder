import { Prepare, Replace } from '../../wailsjs/go/replace/Replace'
import { withBindingError } from '@/api/wails'

async function blobBytes(value, emptyMessage) {
  if (!(value instanceof Blob) || value.size === 0) {
    throw new Error(emptyMessage)
  }
  return Array.from(new Uint8Array(await value.arrayBuffer()))
}

function mappingBytes(mapping) {
  return Array.from(new TextEncoder().encode(JSON.stringify({ data: mapping })))
}

export default {
  prepare(requestData) {
    return withBindingError(
      blobBytes(requestData?.file, 'The file must be passed as an argument').then(bytes =>
        Prepare(bytes),
      ),
    )
  },
  replace(requestData) {
    return withBindingError(
      blobBytes(requestData?.file, 'The file must be passed as an argument').then(fileBytes =>
        Replace(fileBytes, mappingBytes(requestData?.mapping)),
      ),
    )
  },
}
