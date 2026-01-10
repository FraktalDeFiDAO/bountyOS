<template>
  <div class="min-h-screen px-4 pb-16 pt-10 md:px-10">
    <TopNav />
    <main class="mt-10 space-y-10 animate-[riseIn_0.8s_ease-out]">
      <RouterView />
    </main>
  </div>
</template>

<script setup>
import TopNav from './components/TopNav.vue'
</script>
