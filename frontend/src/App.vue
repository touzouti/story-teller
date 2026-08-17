<template>
  <div class="relative min-h-screen text-slate-100">
    <div class="absolute inset-0 overflow-hidden pointer-events-none">
      <div class="absolute -left-10 -top-10 h-64 w-64 rounded-full bg-cyan-400/10 blur-3xl"></div>
      <div class="absolute right-10 top-8 h-48 w-48 rounded-full bg-fuchsia-500/10 blur-3xl"></div>
      <div class="absolute right-1/3 bottom-12 h-72 w-72 rounded-full bg-amber-400/10 blur-3xl"></div>
    </div>

    <LandingPage />
  </div>
  
</template>

<script setup lang="ts">
import LandingPage from "./components/LandingPage.vue";
</script>

<style>
/* Add any global styles if needed */
</style>
