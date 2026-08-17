<template>
  <div
    class="gravity-stars-bg"
    ref="container"
    @mousemove="onMove"
    @mouseleave="onLeave"
  >
    <div
      v-for="star in stars"
      :key="star.id"
      class="gravity-star"
      :style="{
        left: star.x + '%',
        top: star.y + '%',
        width: star.size + 'px',
        height: star.size + 'px',
        opacity: star.opacity,
        transform: `translate3d(${star.parallaxX}px, ${star.parallaxY}px, 0)`,
      }"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";

interface Star {
  id: number;
  x: number;
  y: number;
  size: number;
  opacity: number;
  depth: number;
  parallaxX: number;
  parallaxY: number;
}

const props = withDefaults(
  defineProps<{
    starsCount?: number;
  }>(),
  {
    starsCount: 75,
  }
);

const container = ref<HTMLElement | null>(null);
const stars = ref<Star[]>([]);

function rand(min: number, max: number) {
  return Math.random() * (max - min) + min;
}

onMounted(() => {
  const arr: Star[] = [];
  for (let i = 0; i < props.starsCount; i++) {
    const depth = rand(0.3, 1); // closer vs further stars
    arr.push({
      id: i,
      x: rand(0, 100),
      y: rand(0, 100),
      size: rand(1, 3.2),
      opacity: rand(0.4, 1),
      depth,
      parallaxX: 0,
      parallaxY: 0,
    });
  }
  stars.value = arr;
});

function onMove(event: MouseEvent) {
  const el = container.value;
  if (!el) return;

  const rect = el.getBoundingClientRect();
  const cx = rect.left + rect.width / 2;
  const cy = rect.top + rect.height / 2;
  const dx = (event.clientX - cx) / rect.width; // -0.5..0.5
  const dy = (event.clientY - cy) / rect.height;

  const strength = 40; // px max offset

  for (const star of stars.value) {
    const factor = star.depth; // closer stars move more
    star.parallaxX = dx * strength * factor;
    star.parallaxY = dy * strength * factor;
  }
}

function onLeave() {
  for (const star of stars.value) {
    star.parallaxX = 0;
    star.parallaxY = 0;
  }
}
</script>

<style scoped>
.gravity-stars-bg {
  position: absolute;
  inset: 0;
  overflow: hidden;
  pointer-events: none;
  background: radial-gradient(
      circle at 20% 20%,
      rgba(255, 255, 255, 0.08),
      transparent 55%
    ),
    radial-gradient(
      circle at 80% 80%,
      rgba(255, 255, 255, 0.04),
      transparent 60%
    );
}

.gravity-star {
  position: absolute;
  border-radius: 9999px;
  background: rgba(255, 255, 255, 0.9);
  box-shadow: 0 0 6px rgba(255, 255, 255, 0.6),
    0 0 16px rgba(140, 180, 255, 0.4);
  animation: twinkle 2.8s ease-in-out infinite;
}

.gravity-star:nth-child(3n) {
  animation-duration: 3.3s;
}
.gravity-star:nth-child(5n) {
  animation-duration: 4.1s;
  animation-delay: -1s;
}

@keyframes twinkle {
  0%,
  100% {
    opacity: 0.45;
  }
  50% {
    opacity: 1;
  }
}
</style>
