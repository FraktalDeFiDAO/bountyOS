import { ref, onMounted, onUnmounted } from 'vue'

const isOffline = ref(!navigator.onLine)
const listeners = new Set()

function handleOnline() {
  isOffline.value = false
  listeners.forEach((fn) => fn(false))
}

function handleOffline() {
  isOffline.value = true
  listeners.forEach((fn) => fn(true))
}

export function useOfflineDetection() {
  const offline = ref(isOffline.value)

  function updateOffline(value) {
    offline.value = value
  }

  onMounted(() => {
    listeners.add(updateOffline)
    offline.value = isOffline.value
  })

  onUnmounted(() => {
    listeners.delete(updateOffline)
  })

  return { offline }
}

export function getOfflineState() {
  return isOffline
}

if (typeof window !== 'undefined') {
  window.addEventListener('online', handleOnline)
  window.addEventListener('offline', handleOffline)
}

export function setupOfflineListeners(store) {
  const handleOnlineEvent = () => {
    store.offline = false
    if (store.bounties.length === 0) {
      store.fetchInitial()
    }
    if (!store.connected) {
      store.connectWS()
    }
  }

  const handleOfflineEvent = () => {
    store.offline = true
    if (store.ws) {
      store.ws.close()
    }
  }

  window.addEventListener('online', handleOnlineEvent)
  window.addEventListener('offline', handleOfflineEvent)

  return () => {
    window.removeEventListener('online', handleOnlineEvent)
    window.removeEventListener('offline', handleOfflineEvent)
  }
}
