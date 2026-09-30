import { Prepare, Replace } from '../../wailsjs/go/replace/Replace'
import { withBindingError } from '@/api/wails'

async function blobBytes(value, emptyMessage) {
  if (!(value instanceof Blob) || value.size === 0) {
    throw new Error(emptyMessage)
  }
  return Array.from(new Uint8Array(await value.arrayBuffer()))
}

export default {
  prepare(requestData) {
    const file = requestData instanceof FormData ? requestData.get('file') : requestData?.file
    return withBindingError(
      blobBytes(file, 'The file must be passed as an argument').then(bytes => Prepare(bytes)),
    )
  },
  replace(requestData) {
    const file = requestData instanceof FormData ? requestData.get('file') : requestData?.file
    const mapping =
      requestData instanceof FormData ? requestData.get('mapping') : requestData?.mapping
    return withBindingError(
      Promise.all([
        blobBytes(file, 'The file must be passed as an argument'),
        blobBytes(mapping, 'The mapping must be passed as an argument'),
      ]).then(([fileBytes, mappingBytes]) => Replace(fileBytes, mappingBytes)),
    )
  },
}
