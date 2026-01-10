import { defineStore } from 'pinia'

const okTypes = new Set(['bounty'])

export const useBountiesStore = defineStore('bounties', {
  state: () => ({
    bounties: [],
    connected: false,
    lastUpdated: null,
    error: null,
    wsBackoff: 1500,
    ws: null
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
    }
  },
  actions: {
    async fetchInitial() {
      try {
        const res = await fetch('/api/bounties')
        if (!res.ok) throw new Error(`Failed to fetch bounties: ${res.status}`)
        const data = await res.json()
        this.bounties = data
        this.lastUpdated = new Date()
      } catch (err) {
        this.error = err.message
      }
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
      const protocol = window.location.protocol === 'https:' ? 'wss' : 'ws'
      const wsUrl = `${protocol}://${window.location.host}/ws`

      if (this.ws) {
        this.ws.close()
      }

      const ws = new WebSocket(wsUrl)
      this.ws = ws

      ws.onopen = () => {
        this.connected = true
        this.error = null
        this.wsBackoff = 1500
      }

      ws.onmessage = (event) => {
        try {
          const payload = JSON.parse(event.data)
          if (!payload || !okTypes.has(payload.type)) return
          this.upsertBounty(payload.data)
        } catch (err) {
          this.error = 'Stream parse error'
        }
      }

      ws.onclose = () => {
        this.connected = false
        this.scheduleReconnect()
      }

      ws.onerror = () => {
        this.connected = false
        this.scheduleReconnect()
      }
    },
    scheduleReconnect() {
      const delay = Math.min(this.wsBackoff, 10000)
      setTimeout(() => this.connectWS(), delay)
      this.wsBackoff = Math.min(this.wsBackoff * 1.5, 12000)
    }
  }
})
