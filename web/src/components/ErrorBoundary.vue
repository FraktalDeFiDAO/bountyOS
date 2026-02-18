<template>
  <div v-if="error" class="error-boundary">
    <slot name="error" :error="error" :retry="retry">
      <div class="error-content glass rounded-2xl p-6 text-center">
        <div class="text-4xl mb-4">⚠️</div>
        <p class="text-lg mb-4">{{ errorMessage }}</p>
        <button
          @click="retry"
          class="px-6 py-2 rounded-full bg-[var(--accent)] text-black font-medium hover:opacity-90 transition"
        >
          Try Again
        </button>
      </div>
    </slot>
  </div>
  <slot v-else />
</template>

<script setup>
import { ref, computed, onErrorCaptured } from 'vue'
import { formatError } from '../utils/errors'

const error = ref(null)

const errorMessage = computed(() => {
  if (!error.value) return ''
  return formatError(error.value)
})

onErrorCaptured((err) => {
  error.value = err
  console.error('Error captured by boundary:', err)
  return false
})

const retry = () => {
  error.value = null
}

defineExpose({ retry })
</script>

<style scoped>
.error-boundary {
  display: contents;
}
</style>
