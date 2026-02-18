<template>
  <div class="flex items-center gap-2 mono text-xs">
    <span
      class="w-2 h-2 rounded-full"
      :class="indicatorClass"
      :title="tooltip"
    ></span>
    <span class="text-[var(--muted)]">{{ label }}</span>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  quality: { type: String, default: 'unknown' },
  connected: { type: Boolean, default: false }
})

const indicatorClass = computed(() => {
  if (!props.connected) return 'bg-red-500'
  switch (props.quality) {
    case 'excellent':
      return 'bg-green-500'
    case 'good':
      return 'bg-yellow-500'
    case 'poor':
      return 'bg-orange-500'
    default:
      return 'bg-gray-500'
  }
})

const label = computed(() => {
  if (!props.connected) return 'Disconnected'
  switch (props.quality) {
    case 'excellent':
      return 'Excellent'
    case 'good':
      return 'Good'
    case 'poor':
      return 'Poor'
    default:
      return 'Connecting...'
  }
})

const tooltip = computed(() => {
  if (!props.connected) return 'Not connected to server'
  switch (props.quality) {
    case 'excellent':
      return 'Connection is fast and stable'
    case 'good':
      return 'Connection is stable'
    case 'poor':
      return 'Connection is slow'
    default:
      return 'Establishing connection...'
  }
})
</script>
