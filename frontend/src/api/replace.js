import { Prepare, Replace } from '../../wailsjs/go/replace/Replace'
import { fileToBase64, withBindingError } from '@/api/wails'

async function blobBytes(value, emptyMessage) {
  if (!(value instanceof Blob) || value.size === 0) {
    throw new Error(emptyMessage)
  }
  return fileToBase64(value)
}

function mappingBytes(mapping) {
  return fileToBase64(new Blob([JSON.stringify({ data: mapping })], { type: 'application/json' }))
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
      Promise.all([
        blobBytes(requestData?.file, 'The file must be passed as an argument'),
        mappingBytes(requestData?.mapping),
      ]).then(([fileBytes, mapping]) => Replace(fileBytes, mapping)),
    )
  },
}
