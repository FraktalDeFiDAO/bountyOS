<template>
  <Transition name="slide-up">
    <div v-if="visible" class="fixed bottom-6 left-1/2 -translate-x-1/2 z-50 glass rounded-xl px-6 py-4 flex items-center gap-4 shadow-lg">
      <span class="text-lg">{{ icon }}</span>
      <p class="text-sm">{{ message }}</p>
      <button
        @click="$emit('retry')"
        class="px-4 py-1 rounded-full bg-[var(--accent)] text-black text-sm font-medium hover:opacity-90 transition"
      >
        Retry
      </button>
      <button
        @click="$emit('dismiss')"
        class="text-[var(--muted)] hover:text-white transition"
      >
        ✕
      </button>
    </div>
  </Transition>
</template>

<script setup>
import { computed } from 'vue'
import { ErrorTypes } from '../utils/errors'

const props = defineProps({
  error: { type: Error, default: null },
  errorType: { type: String, default: null }
})

defineEmits(['retry', 'dismiss'])

const visible = computed(() => !!props.error)

const icon = computed(() => {
  switch (props.errorType) {
    case ErrorTypes.NETWORK:
      return '📡'
    case ErrorTypes.TIMEOUT:
      return '⏱️'
    case ErrorTypes.SERVER:
      return '🖥️'
    default:
      return '⚠️'
  }
})

const message = computed(() => {
  switch (props.errorType) {
    case ErrorTypes.NETWORK:
      return 'Connection issue detected'
    case ErrorTypes.TIMEOUT:
      return 'Request timed out'
    case ErrorTypes.SERVER:
      return 'Server error'
    default:
      return 'Something went wrong'
  }
})
</script>

<style scoped>
.slide-up-enter-active,
.slide-up-leave-active {
  transition: all 0.3s ease;
}

.slide-up-enter-from,
.slide-up-leave-to {
  opacity: 0;
  transform: translate(-50%, 20px);
}
</style>
