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

export async function writeRequestFromBody(body) {
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
