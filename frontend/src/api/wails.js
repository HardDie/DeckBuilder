import { useToast } from 'vue-toastification'

// withBindingError shows a binding error as a toast and rethrows it.
// A successful result with `warning` (e.g. an image that was not saved) shows a warning toast.
export function withBindingError(promise) {
  return promise
    .then(result => {
      if (result?.warning) {
        useToast().warning(result.warning)
      }
      return result
    })
    .catch(err => {
      const toast = useToast()
      const message = typeof err === 'string' ? err : err?.message || 'Unknown error'
      toast.error(message)
      throw err
    })
}

async function imageFileBytes(file) {
  if (!(file instanceof Blob)) {
    return []
  }
  return Array.from(new Uint8Array(await file.arrayBuffer()))
}

function formCount(value) {
  if (value === null || value === undefined || value === '') {
    return 1
  }
  const n = parseInt(value, 10)
  return Number.isFinite(n) ? n : 1
}

function formVariables(value) {
  if (value == null || value === '') {
    return null
  }
  if (typeof value === 'object' && !(value instanceof Blob)) {
    return value
  }
  try {
    const parsed = JSON.parse(value)
    if (parsed && typeof parsed === 'object') {
      return parsed
    }
  } catch {
    // match HTTP CreateHandler
  }
  throw new Error('Bad variables json')
}

export async function writeRequestFromBody(body) {
  const file = body?.imageFile
  const hasFile = file instanceof Blob && file.size > 0
  return {
    name: body?.name || '',
    description: body?.description || '',
    image: hasFile ? '' : body?.image || '',
    imageFile: hasFile ? await imageFileBytes(file) : null,
    count: formCount(body?.count),
    variables: formVariables(body?.variables),
  }
}
