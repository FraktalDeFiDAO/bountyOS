<template>
  <section class="space-y-8">
    <StatsPanel />

    <div class="grid gap-6 lg:grid-cols-[2fr,1fr]">
      <div class="glass rounded-3xl p-6 md:p-8 space-y-4">
        <p class="mono text-xs uppercase tracking-[0.35em] text-[var(--muted)]">Priority Picks</p>
        <div class="space-y-4">
          <BountyCard v-for="bounty in store.topBounties" :key="bounty.url" :bounty="bounty" />
        </div>
      </div>
      <div class="glass rounded-3xl p-6 md:p-8 space-y-4">
        <p class="mono text-xs uppercase tracking-[0.35em] text-[var(--muted)]">Platform Mix</p>
        <div class="space-y-3">
          <div v-for="platform in platformStats" :key="platform.name" class="flex items-center justify-between">
            <span class="mono text-sm text-[var(--text)]">{{ platform.name }}</span>
            <span class="mono text-xs text-[var(--muted)]">{{ platform.count }}</span>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup>
import { computed, onMounted } from 'vue'
import { useBountiesStore } from '../stores/bounties'
import StatsPanel from '../components/StatsPanel.vue'
import BountyCard from '../components/BountyCard.vue'

const store = useBountiesStore()

const platformStats = computed(() => {
  const entries = Object.entries(store.stats.byPlatform || {})
  return entries
    .map(([name, count]) => ({ name, count }))
    .sort((a, b) => b.count - a.count)
    .slice(0, 8)
})

onMounted(() => {
  if (!store.bounties.length) {
    store.fetchInitial()
  }
  if (!store.connected) {
    store.connectWS()
  }
})
</script>
