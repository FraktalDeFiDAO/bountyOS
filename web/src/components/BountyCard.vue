<template>
  <article class="glass rounded-2xl p-5 flex flex-col gap-4 hover:translate-y-[-2px] transition">
    <div class="flex items-start justify-between gap-4">
      <div>
        <p class="mono text-xs uppercase tracking-[0.3em] text-[var(--muted)]">{{ bounty.platform }}</p>
        <h3 class="title-font text-xl font-semibold mt-2">{{ bounty.title }}</h3>
      </div>
      <div class="text-right">
        <span class="mono text-xs uppercase tracking-[0.3em] text-[var(--muted)]">Score</span>
        <div class="text-2xl font-semibold" :class="scoreColor">{{ bounty.score }}</div>
      </div>
    </div>
    <div class="flex flex-wrap items-center gap-3">
      <span class="mono text-sm px-3 py-1 rounded-full bg-[rgba(110,231,216,0.18)] text-[var(--accent)]">
        {{ bounty.reward }} {{ bounty.currency }}
      </span>
      <span v-if="bounty.payment_type" class="mono text-xs uppercase tracking-[0.2em] text-[var(--muted)]">
        {{ bounty.payment_type }}
      </span>
      <span v-for="tag in bounty.tags || []" :key="tag" class="text-xs px-2 py-1 rounded-full bg-[rgba(255,255,255,0.06)]">
        {{ tag }}
      </span>
    </div>
    <a
      class="mono text-xs uppercase tracking-[0.3em] text-[var(--accent-3)] hover:text-[var(--accent)]"
      :href="bounty.url"
      target="_blank"
      rel="noreferrer"
    >
      Open bounty →
    </a>
  </article>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  bounty: {
    type: Object,
    required: true
  }
})

const scoreColor = computed(() => {
  const score = props.bounty.score || 0
  if (score >= 80) return 'text-[#ff8a3d]'
  if (score >= 50) return 'text-[var(--accent)]'
  return 'text-[var(--accent-3)]'
})
</script>
