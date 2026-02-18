<template>
  <div class="glass rounded-2xl p-8 text-center">
    <div class="text-5xl mb-4">{{ icon }}</div>
    <h3 class="title-font text-xl font-semibold mb-2">{{ title }}</h3>
    <p class="text-[var(--muted)] mb-6">{{ message }}</p>
    <div v-if="retryable" class="flex flex-col items-center gap-3">
      <button
        @click="$emit('retry')"
        class="px-6 py-2 rounded-full bg-[var(--accent)] text-black font-medium hover:opacity-90 transition"
      >
        Retry
      </button>
      <p v-if="retryCount > 0" class="mono text-xs text-[var(--muted)]">
        Attempt {{ retryCount }} of {{ maxRetries }}
      </p>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { ErrorTypes } from '../utils/errors'

const props = defineProps({
  error: { type: Error, default: null },
  errorType: { type: String, default: null },
  retryable: { type: Boolean, default: false },
  retryCount: { type: Number, default: 0 },
  maxRetries: { type: Number, default: 5 }
})

defineEmits(['retry'])

const icon = computed(() => {
  switch (props.errorType) {
    case ErrorTypes.NETWORK:
      return '📡'
    case ErrorTypes.TIMEOUT:
      return '⏱️'
    case ErrorTypes.SERVER:
      return '🖥️'
    case ErrorTypes.VALIDATION:
      return '⚠️'
    default:
      return '❌'
  }
})

const title = computed(() => {
  switch (props.errorType) {
    case ErrorTypes.NETWORK:
      return 'Connection Error'
    case ErrorTypes.TIMEOUT:
      return 'Request Timed Out'
    case ErrorTypes.SERVER:
      return 'Server Error'
    case ErrorTypes.VALIDATION:
      return 'Invalid Request'
    default:
      return 'Something Went Wrong'
  }
})

const message = computed(() => {
  switch (props.errorType) {
    case ErrorTypes.NETWORK:
      return 'Unable to connect to the server. Please check your internet connection.'
    case ErrorTypes.TIMEOUT:
      return 'The request took too long to complete. Please try again.'
    case ErrorTypes.SERVER:
      return 'The server is experiencing issues. Please try again later.'
    case ErrorTypes.VALIDATION:
      return 'There was a problem with your request.'
    default:
      return props.error?.message || 'An unexpected error occurred.'
  }
})
</script>
