<template>
  <header class="glass rounded-3xl p-6 md:p-8 flex flex-col gap-6">
    <div class="flex flex-col gap-6 md:flex-row md:items-center md:justify-between">
      <div class="space-y-2">
        <p class="mono text-sm uppercase tracking-[0.3em] text-[var(--accent)]">Obsidian Radar</p>
        <h1 class="title-font text-3xl md:text-4xl font-semibold text-[var(--text)]">
          BountyOS Live Ops Console
        </h1>
        <p class="text-[var(--muted)] max-w-xl">
          Real-time bounty intelligence powered by WebSockets, prioritizing funded work across the ecosystem.
        </p>
      </div>
      <div class="flex flex-wrap items-center gap-4">
        <div
          class="flex items-center gap-2 rounded-full px-4 py-2 text-sm mono border"
          :class="connected ? 'border-[var(--accent)] text-[var(--accent)]' : 'border-[#ff8a3d] text-[#ff8a3d]'"
        >
          <span
            class="h-2 w-2 rounded-full"
            :class="connected ? 'bg-[var(--accent)] animate-[pulseGlow_3s_ease-in-out_infinite]' : 'bg-[#ff8a3d]'"
          ></span>
          {{ connected ? 'WS CONNECTED' : 'WS RETRYING' }}
        </div>
        <div class="text-xs mono text-[var(--muted)]">
          Updated: {{ lastUpdated || 'syncing...' }}
        </div>
      </div>
    </div>
    <nav class="flex flex-wrap gap-3">
      <RouterLink
        to="/"
        class="px-4 py-2 rounded-full mono text-sm border border-transparent"
        :class="$route.name === 'dashboard' ? 'bg-[var(--accent)] text-black' : 'border-[rgba(255,255,255,0.08)] text-[var(--muted)]'"
      >
        Overview
      </RouterLink>
      <RouterLink
        to="/feed"
        class="px-4 py-2 rounded-full mono text-sm border border-transparent"
        :class="$route.name === 'feed' ? 'bg-[var(--accent-2)] text-black' : 'border-[rgba(255,255,255,0.08)] text-[var(--muted)]'"
      >
        Live Feed
      </RouterLink>
    </nav>
  </header>
</template>

<script setup>
import { computed } from 'vue'
import { useBountiesStore } from '../stores/bounties'

const store = useBountiesStore()

const connected = computed(() => store.connected)
const lastUpdated = computed(() =>
  store.lastUpdated ? store.lastUpdated.toLocaleTimeString() : ''
)
</script>
