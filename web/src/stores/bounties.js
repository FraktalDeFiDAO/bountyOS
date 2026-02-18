import { defineStore } from 'pinia'
import { api } from '../utils/api'
import { classifyError, ErrorTypes } from '../utils/errors'
import { setupOfflineListeners } from '../utils/offline'

const okTypes = new Set(['bounty'])
const MAX_RETRIES = 5
const INITIAL_BACKOFF = 1000
const MAX_BACKOFF = 30000

export const useBountiesStore = defineStore('bounties', {
  state: () => ({
    bounties: [],
    connected: false,
    connectionQuality: 'unknown',
    loading: 'idle',
    error: null,
    errorType: null,
    lastUpdated: null,
    retryCount: 0,
    maxRetries: MAX_RETRIES,
    wsBackoff: INITIAL_BACKOFF,
    ws: null,
    offline: !navigator.onLine,
    wsLatency: null,
    _cleanup: null
  }),
  getters: {
    sortedBounties: (state) =>
      [...state.bounties].sort((a, b) => (b.score || 0) - (a.score || 0)),
    topBounties: (state) =>
      [...state.bounties].sort((a, b) => (b.score || 0) - (a.score || 0)).slice(0, 8),
    stats: (state) => {
      const total = state.bounties.length
      const byPlatform = {}
      let totalScore = 0
      let crypto = 0

      state.bounties.forEach((b) => {
        const platform = b.platform || 'UNKNOWN'
        byPlatform[platform] = (byPlatform[platform] || 0) + 1
        totalScore += b.score || 0
        if (b.payment_type === 'crypto') {
          crypto += 1
        }
      })

      return {
        total,
        crypto,
        avgScore: total ? totalScore / total : 0,
        platforms: Object.keys(byPlatform).length,
        byPlatform
      }
    },
    isRetryable: (state) => {
      if (!state.error) return false
      const classified = classifyError(state.error)
      return classified.retryable
    },
    isLoading: (state) => state.loading === 'initial',
    isRefreshing: (state) => state.loading === 'refreshing',
    isReady: (state) => state.loading === 'idle' && state.bounties.length > 0
  },
  actions: {
    setupOfflineDetection() {
      if (this._cleanup) {
        this._cleanup()
      }
      this._cleanup = setupOfflineListeners(this)
    },

    updateConnectionQuality(latency) {
      this.wsLatency = latency
      if (latency === null) {
        this.connectionQuality = 'unknown'
      } else if (latency < 100) {
        this.connectionQuality = 'excellent'
      } else if (latency < 300) {
        this.connectionQuality = 'good'
      } else {
        this.connectionQuality = 'poor'
      }
    },

    async fetchInitial() {
      if (this.loading === 'initial') return

      this.loading = 'initial'
      this.error = null
      this.errorType = null

      try {
        const startTime = performance.now()
        const data = await api.get('/api/bounties')
        const latency = performance.now() - startTime

        this.bounties = data
        this.lastUpdated = new Date()
        this.loading = 'idle'
        this.retryCount = 0
        this.updateConnectionQuality(latency)
      } catch (err) {
        const classified = classifyError(err)
        this.error = err
        this.errorType = classified.type
        this.loading = 'idle'

        if (this.retryCount < this.maxRetries) {
          this.retryCount++
          const delay = this.getRetryDelay(this.retryCount)
          setTimeout(() => this.fetchInitial(), delay)
        }
      }
    },

    async refresh() {
      if (this.loading !== 'idle') return

      this.loading = 'refreshing'
      this.error = null
      this.errorType = null

      try {
        const startTime = performance.now()
        const data = await api.get('/api/bounties')
        const latency = performance.now() - startTime

        this.bounties = data
        this.lastUpdated = new Date()
        this.loading = 'idle'
        this.retryCount = 0
        this.updateConnectionQuality(latency)
      } catch (err) {
        const classified = classifyError(err)
        this.error = err
        this.errorType = classified.type
        this.loading = 'idle'
      }
    },

    retry() {
      this.retryCount = 0
      this.error = null
      this.errorType = null
      this.fetchInitial()
    },

    getRetryDelay(attempt) {
      const baseDelay = INITIAL_BACKOFF
      const exponentialDelay = baseDelay * Math.pow(2, attempt - 1)
      const jitter = Math.random() * 500
      return Math.min(exponentialDelay + jitter, MAX_BACKOFF)
    },

    upsertBounty(bounty) {
      if (!bounty || !bounty.url) return
      const idx = this.bounties.findIndex((b) => b.url === bounty.url)
      if (idx === -1) {
        this.bounties.unshift(bounty)
      } else {
        this.bounties[idx] = bounty
      }
      this.lastUpdated = new Date()
    },

    connectWS() {
      if (this.offline) return

      const protocol = window.location.protocol === 'https:' ? 'wss' : 'ws'
      const wsUrl = `${protocol}://${window.location.host}/ws`

      if (this.ws) {
        this.ws.close()
      }

      const ws = new WebSocket(wsUrl)
      this.ws = ws
      const connectTime = Date.now()

      ws.onopen = () => {
        this.connected = true
        this.error = null
        this.errorType = null
        this.wsBackoff = INITIAL_BACKOFF
        const latency = Date.now() - connectTime
        this.updateConnectionQuality(latency)
      }

      ws.onmessage = (event) => {
        try {
          const payload = JSON.parse(event.data)
          if (!payload || !okTypes.has(payload.type)) return
          this.upsertBounty(payload.data)
        } catch (err) {
          this.error = 'Stream parse error'
          this.errorType = ErrorTypes.SERVER
        }
      }

      ws.onclose = (event) => {
        this.connected = false
        this.updateConnectionQuality(null)
        if (!this.offline && event.code !== 1000) {
          this.scheduleReconnect()
        }
      }

      ws.onerror = () => {
        this.connected = false
        this.updateConnectionQuality(null)
        if (!this.offline) {
          this.scheduleReconnect()
        }
      }
    },

    scheduleReconnect() {
      const delay = Math.min(this.wsBackoff, MAX_BACKOFF)
      setTimeout(() => {
        if (!this.offline) {
          this.connectWS()
        }
      }, delay)
      this.wsBackoff = Math.min(this.wsBackoff * 2, MAX_BACKOFF)
    },

    disconnect() {
      if (this.ws) {
        this.ws.close(1000)
        this.ws = null
      }
      this.connected = false
      this.updateConnectionQuality(null)
    },

    cleanup() {
      this.disconnect()
      if (this._cleanup) {
        this._cleanup()
        this._cleanup = null
      }
    }
  }
})
