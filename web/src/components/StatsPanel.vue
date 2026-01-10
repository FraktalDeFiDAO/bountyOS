<template>
  <section class="grid gap-4 md:grid-cols-4">
    <div v-for="stat in stats" :key="stat.label" class="glass rounded-2xl p-5">
      <p class="mono text-xs uppercase tracking-[0.35em] text-[var(--muted)]">{{ stat.label }}</p>
      <p class="title-font text-3xl font-semibold mt-3 text-[var(--text)]">{{ stat.value }}</p>
      <p class="text-xs text-[var(--muted)] mt-2">{{ stat.caption }}</p>
    </div>
  </section>
</template>

<script setup>
import { computed } from 'vue'
import { useBountiesStore } from '../stores/bounties'

const store = useBountiesStore()

const stats = computed(() => {
  const summary = store.stats
  return [
    {
      label: 'Total Bounties',
      value: summary.total,
      caption: 'Validated, reachable links only'
    },
    {
      label: 'Crypto-Settled',
      value: summary.crypto,
      caption: 'Fast settlement, highest priority'
    },
    {
      label: 'Avg. Urgency',
      value: summary.avgScore.toFixed(1),
      caption: 'Score derived from Obsidian rules'
    },
    {
      label: 'Sources',
      value: summary.platforms,
      caption: 'Active platforms in the stream'
    }
  ]
})
</script>
