import { useStore } from '@/stores/main'

export const setBar = (to, from, next) => {
  const mainStore = useStore()
  mainStore.fetchPathItems(to.path)
  next()
}
