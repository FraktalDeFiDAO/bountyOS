<template>
  <section class="space-y-6">
    <div class="glass rounded-3xl p-6 md:p-8 flex flex-col gap-6">
      <div class="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
        <div>
          <p class="mono text-xs uppercase tracking-[0.35em] text-[var(--muted)]">Live Feed</p>
          <h2 class="title-font text-2xl md:text-3xl font-semibold">All Active Bounties</h2>
        </div>
        <div class="mono text-xs text-[var(--muted)]">Total: {{ store.stats.total }}</div>
      </div>
      <div class="relative">
        <input
          v-model="query"
          class="w-full rounded-full bg-[rgba(255,255,255,0.06)] border border-transparent focus:border-[var(--accent)] px-5 py-3 text-sm outline-none"
          placeholder="Search by title, platform, or tag..."
        />
      </div>
    </div>

    <div class="grid gap-4 md:grid-cols-2">
      <BountyCard v-for="bounty in filteredBounties" :key="bounty.url" :bounty="bounty" />
    </div>
  </section>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import { useBountiesStore } from '../stores/bounties'
import BountyCard from '../components/BountyCard.vue'

const store = useBountiesStore()
const query = ref('')

const filteredBounties = computed(() => {
  const term = query.value.trim().toLowerCase()
  const list = store.sortedBounties
  if (!term) return list

  return list.filter((b) => {
    const haystack = `${b.title} ${b.platform} ${(b.tags || []).join(' ')}`.toLowerCase()
    return haystack.includes(term)
  })
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
