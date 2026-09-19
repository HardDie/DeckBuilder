import { useToast } from 'vue-toastification'

export function withBindingError(promise) {
  return promise.catch(err => {
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
  if (body instanceof FormData) {
    const file = body.get('imageFile')
    const hasFile = file instanceof Blob && file.size > 0
    return {
      name: body.get('name') || '',
      description: body.get('description') || '',
      image: hasFile ? '' : body.get('image') || '',
      imageFile: hasFile ? await imageFileBytes(file) : null,
      count: formCount(body.get('count')),
      variables: formVariables(body.get('variables')),
    }
  }
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
