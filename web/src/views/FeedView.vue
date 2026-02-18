<template>
  <section class="space-y-6">
    <OfflineBanner v-if="store.offline" />

    <div class="glass rounded-3xl p-6 md:p-8 flex flex-col gap-6">
      <div class="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
        <div>
          <p class="mono text-xs uppercase tracking-[0.35em] text-[var(--muted)]">Live Feed</p>
          <h2 class="title-font text-2xl md:text-3xl font-semibold">All Active Bounties</h2>
        </div>
        <div class="flex items-center gap-4">
          <ConnectionQuality :quality="store.connectionQuality" :connected="store.connected" />
          <div class="mono text-xs text-[var(--muted)]">Total: {{ store.stats.total }}</div>
        </div>
      </div>
      <div class="relative">
        <input
          v-model="query"
          class="w-full rounded-full bg-[rgba(255,255,255,0.06)] border border-transparent focus:border-[var(--accent)] px-5 py-3 text-sm outline-none"
          placeholder="Search by title, platform, or tag..."
          :disabled="store.isLoading"
        />
      </div>
    </div>

    <ErrorState
      v-if="store.error && !store.isLoading && store.bounties.length === 0"
      :error="store.error"
      :error-type="store.errorType"
      :retryable="store.isRetryable"
      :retry-count="store.retryCount"
      :max-retries="store.maxRetries"
      @retry="store.retry"
    />

    <div v-else-if="store.isLoading" class="grid gap-4 md:grid-cols-2">
      <BountyCardSkeleton v-for="i in 6" :key="i" />
    </div>

    <template v-else>
      <div class="grid gap-4 md:grid-cols-2">
        <BountyCard v-for="bounty in filteredBounties" :key="bounty.url" :bounty="bounty" />
      </div>

      <ErrorToast
        v-if="store.error && store.bounties.length > 0"
        :error="store.error"
        :error-type="store.errorType"
        @retry="store.refresh"
        @dismiss="store.error = null"
      />
    </template>
  </section>
</template>

<script setup>
import { computed, onMounted, onUnmounted } from 'vue'
import { useBountiesStore } from '../stores/bounties'
import BountyCard from '../components/BountyCard.vue'
import BountyCardSkeleton from '../components/BountyCardSkeleton.vue'
import ErrorState from '../components/ErrorState.vue'
import ErrorToast from '../components/ErrorToast.vue'
import OfflineBanner from '../components/OfflineBanner.vue'
import ConnectionQuality from '../components/ConnectionQuality.vue'

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
  store.setupOfflineDetection()
  if (!store.bounties.length) {
    store.fetchInitial()
  }
  if (!store.connected) {
    store.connectWS()
  }
})

onUnmounted(() => {
  store.cleanup()
})
</script>

<script>
import { ref } from 'vue'
export default { name: 'FeedView' }
</script>
