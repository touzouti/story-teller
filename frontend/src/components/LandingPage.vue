<template>
  <div class="min-h-screen flex flex-col">
    <!-- SCENE: always present (background, trees, moon, etc.) -->
    <div
      class="hero min-h-screen bg-base-200 relative overflow-hidden"
      :class="{ animating, arrived, landing, 'story-mode': showStoryOverlay }"
      :style="{
        '--finalStart': finalStart + 's',
        '--finalDuration': finalDuration + 's',
      }"
    >
      <!-- Gravity stars background -->
      <GravityStarsBackground
        class="absolute inset-0"
        :starsCount="140"
        :mouseInfluence="80"
        :gravityStrength="1"
        glowAnimation="spring"
        :starsInteraction="true"
        starsInteractionType="bounce"
      />

      <div class="mist-layer" aria-hidden="true"></div>

      <!-- Landing hero text -->
      <transition name="fade-hero" appear>
        <div
          v-if="showHeroContent"
          class="hero-content text-center z-50 pointer-events-auto"
        >
          <div class="max-w-2xl text-center">
            <img
              src="/logo.png"
              alt="Logo Story Teller"
              class="w-24 mx-auto h-auto mb-4 drop-shadow"
            />
            <h1 class="text-5xl font-bold text-primary-focus">
              Bienvenue, petite exploratrice !
            </h1>
            <p
              class="py-6 text-lg text-slate-200 max-w-2xl mx-auto leading-relaxed"
            >
              Embarque pour un voyage magique : lance l'animation pour entrer
              dans le hub des chapitres et choisir l'histoire de ce soir.
            </p>
            <div class="flex flex-wrap justify-center gap-3">
              <button
                @click="playAnimation"
                class="btn btn-primary btn-lg shadow-lg shadow-cyan-500/20"
              >
                Commencer l'aventure
              </button>
            </div>
          </div>
        </div>
      </transition>

      <!-- Scene props -->
      <img
        src="/moon.png"
        alt="Lune"
        class="absolute top-10 left-10 w-48 h-auto z-20 moon"
      />

      <!-- Tent: hidden by default, appears with finals or when 'arrived' -->
      <img
        src="/tent.png"
        alt="Maison"
        class="absolute bottom-0 w-32 h-auto z-40 tent"
        style="right: 20%"
      />

      <!-- Static trees (fade in on load & when coming back) -->
      <transition name="fade-static" appear>
        <div v-if="showStaticTrees && !arrived" class="static-trees">
          <img
            src="/tree-spark-1.png"
            alt="Arbre"
            class="absolute bottom-0 w-92 h-auto z-20 tree-static"
            id="tree1"
            style="left: -10rem"
          />
          <img
            src="/tree-spark-2.png"
            alt="Arbre"
            class="absolute bottom-0 w-80 h-auto z-20 tree-static"
            id="tree2"
            style="right: -5rem"
          />
          <img
            src="/tree-spark-3.png"
            alt="Arbre"
            class="absolute bottom-0 w-92 h-auto z-10 tree-static"
            id="tree3"
            style="right: -10rem"
          />
          <img
            src="/tree-1.png"
            alt="Arbre"
            class="absolute bottom-0 w-92 h-auto z-20 tree-static"
            id="tree4"
          />
          <img
            src="/tree-2.png"
            alt="Arbre"
            class="absolute bottom-0 w-80 h-auto z-20 tree-static"
            id="tree5"
          />
          <img
            src="/tree-3.png"
            alt="Arbre"
            class="absolute bottom-0 w-92 h-auto z-10 tree-static"
            id="tree5"
          />
        </div>
      </transition>

      <!-- Final trees: animate from sides and land at final spots -->
      <img
        src="/tree-1.png"
        alt="Arbre"
        class="absolute bottom-0 w-92 h-auto final-tree"
        id="final-tree1"
        style="--to-left: 14%; --from-left: -32%"
      />
      <img
        src="/tree-2.png"
        alt="Arbre"
        class="absolute bottom-0 w-80 h-auto final-tree"
        id="final-tree2"
        style="--to-left: 84%; --from-left: 132%"
      />
      <img
        src="/tree-3.png"
        alt="Arbre"
        class="absolute bottom-0 w-92 h-auto final-tree"
        id="final-tree3"
        style="--to-left: 92%; --from-left: 142%"
      />
    </div>

    <!-- Categories + Lottie overlay (fades over the scene) -->
    <transition name="fade-categories" @after-leave="onCategoriesLeave">
      <div
        v-if="categoriesVisible"
        class="categories-container"
        :class="{ 'story-overlay-active': showStoryOverlay }"
      >
        <div v-if="showLottie && !showStoryOverlay" class="lottie-overlay">
          <Vue3Lottie
            :animationLink="lottieLink"
            :loop="true"
            :autoplay="true"
            :speed="1"
            :width="320"
            :height="320"
          />
        </div>
        <Categories
          @go-back="reverseAnimation"
          @story-overlay="toggleStoryOverlay"
        />
      </div>
    </transition>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, computed } from "vue";
import Categories from "./Categories.vue";
import { Vue3Lottie } from "vue3-lottie";
import GravityStarsBackground from "./GravityStarsBackground.vue";

const STAGE_KEY = "storyteller:stage";

// Initialize state synchronously from localStorage to prevent FOUC/flashing
const savedStage = localStorage.getItem(STAGE_KEY);
const isCategoriesStage = savedStage === "categories";

const animating = ref(false);
const arrived = ref(isCategoriesStage);
const landing = computed(() => !animating.value && !arrived.value);
const showStoryOverlay = ref(false);

const isReversing = ref(false);
const categoriesVisible = ref(isCategoriesStage);
const showHeroContent = ref(!isCategoriesStage);
const showStaticTrees = ref(!isCategoriesStage); // Hide static trees if already arrived

const finalStart = ref(0);
const finalDuration = ref(2.6);

const showLottie = ref(isCategoriesStage);
const lottieLink = "/lottie/owls.json";

const playAnimation = () => {
  showHeroContent.value = false;
  animating.value = true;
  showStaticTrees.value = false;

  finalStart.value = 0.15;
  finalDuration.value = 2.6;

  const showCategoriesAt = finalStart.value + finalDuration.value + 0.3;
  window.setTimeout(() => {
    arrived.value = true;
    categoriesVisible.value = true;
  }, showCategoriesAt * 1000);
};

const reverseAnimation = () => {
  isReversing.value = true;
  categoriesVisible.value = false;
};

const toggleStoryOverlay = (on: boolean) => {
  showStoryOverlay.value = on;
};

const onCategoriesLeave = () => {
  if (!isReversing.value) return;

  isReversing.value = false;
  localStorage.setItem(STAGE_KEY, "landing");

  animating.value = false;
  arrived.value = false;
  showHeroContent.value = true;
  showStaticTrees.value = true; // ensure trees are shown again
};

// Save stage whenever categories visibility changes, and trigger Lottie
watch(categoriesVisible, (vis) => {
  if (vis) {
    localStorage.setItem(STAGE_KEY, "categories");
    showLottie.value = true;
  } else {
    localStorage.setItem(STAGE_KEY, "landing");
    showLottie.value = false;
  }
});
</script>

<style scoped>
/* Scene background */
.hero {
  background-image: url("https://daisyui.com/images/stock/photo-1507851035232-5677ad759701.jpg");
  background-size: cover;
  background-position: center;
  position: fixed;
  inset: 0;
  overflow: hidden;
}

.hero.animating {
  background-position: center 12%;
  filter: brightness(1.05) saturate(1.04);
}

/* Click-through safety */
.tree-static,
.final-tree {
  pointer-events: none;
}

.mist-layer {
  position: absolute;
  inset: 0;
  background: radial-gradient(
      circle at 50% 60%,
      rgba(255, 255, 255, 0.08),
      transparent 45%
    ),
    linear-gradient(180deg, rgba(10, 17, 31, 0.05), rgba(10, 17, 31, 0.45));
  mix-blend-mode: screen;
  opacity: 0.55;
  animation: fog-drift 14s ease-in-out infinite alternate;
  z-index: 5;
  pointer-events: none;
}

/* Fade for hero content (landing text) */
.fade-hero-enter-active,
.fade-hero-leave-active {
  transition: opacity 0.5s;
}
.fade-hero-enter-from,
.fade-hero-leave-to {
  opacity: 0;
}

/* Fade for static trees group */
.fade-static-enter-active,
.fade-static-leave-active {
  transition: opacity 0.7s;
}
.fade-static-enter-from,
.fade-static-leave-to {
  opacity: 0;
}

/* Fade for categories overlay ONLY (scene stays visible) */
.fade-categories-enter-active,
.fade-categories-leave-active {
  transition: opacity 0.5s, transform 0.6s ease;
}
.fade-categories-enter-from,
.fade-categories-leave-to {
  opacity: 0;
  transform: translateY(24px) scale(0.92);
}

/* Base transitions */
.moon,
.tree-static,
.final-tree,
.tent {
  transition: transform 4.8s ease-in-out, opacity 4.8s ease-in-out;
}

/* Moon subtle motion */
.animating .moon {
  transform: rotate(12deg) scale(1.05);
}

/* ===== Static trees positioning ===== */
.static-trees #tree4 {
  left: -30rem;
  bottom: -5rem;
  transform: scale(0.8);
}
.static-trees #tree5 {
  right: -30rem;
  bottom: -5rem;
  transform: scale(0.9);
}
.static-trees #tree6 {
  left: 50%;
  transform: translateX(-50%) scale(0.7);
  bottom: -10rem;
}

/* Drift to sides SLOWLY and fade out when animating */
.animating #tree1 {
  transform: translateX(-120vw) scale(1.6);
  opacity: 0;
  transition-duration: 5s;
}
.animating #tree2 {
  transform: translateX(120vw) scale(1.8);
  opacity: 0;
  transition-duration: 5.3s;
}
.animating #tree3 {
  transform: translateX(120vw) scale(2);
  opacity: 0;
  transition-duration: 5.6s;
}
.animating #tree4 {
  transform: translateX(95vw) scale(1.9);
  opacity: 0;
  transition-delay: 1s;
  transition-duration: 5s;
}
.animating #tree5 {
  transform: translateX(-95vw) scale(2.1);
  opacity: 0;
  transition-delay: 1.3s;
  transition-duration: 5.3s;
}
.animating #tree6 {
  transform: translateX(-50%) translateY(0) scale(1.9);
  opacity: 0;
  transition-delay: 0.8s;
  transition-duration: 4.8s;
}

@keyframes fog-drift {
  0% {
    transform: translate3d(-12px, 6px, 0) scale(1);
    opacity: 0.7;
  }
  50% {
    transform: translate3d(12px, -4px, 0) scale(1.02);
    opacity: 0.9;
  }
  100% {
    transform: translate3d(-8px, 10px, 0) scale(1.03);
    opacity: 0.65;
  }
}

/* ===== Final trees ===== */
.final-tree {
  opacity: 0;
  visibility: hidden; /* Prevent FOUC */
  position: absolute;
  bottom: 0;
  z-index: 40;
  filter: drop-shadow(0 8px 12px rgba(0, 0, 0, 0.25));
  left: var(--from-left, 50%);
  transform: translateX(-50%) scale(0.9);
}

.animating #final-tree1 {
  visibility: visible;
  animation: final-move var(--finalDuration) ease-out var(--finalStart) forwards;
}
.animating #final-tree2 {
  visibility: visible;
  animation: final-move var(--finalDuration) ease-out
    calc(var(--finalStart) + 0.18s) forwards;
}
.animating #final-tree3 {
  visibility: visible;
  animation: final-move var(--finalDuration) ease-out
    calc(var(--finalStart) + 0.36s) forwards;
}

@keyframes final-move {
  0% {
    opacity: 0;
    left: var(--from-left, 50%);
    transform: translateX(-50%) scale(0.9);
  }
  25% {
    opacity: 0.9;
  }
  100% {
    opacity: 1;
    left: var(--to-left);
    transform: translateX(-50%) scale(1);
  }
}

/* ===== Tent ===== */
.tent {
  opacity: 0;
  visibility: hidden; /* Prevent FOUC */
  transform: translate(30%, 40%) scale(0.85);
}
.story-mode .moon,
.story-mode .tree-static,
.story-mode .final-tree,
.story-mode .tent,
.story-mode .lottie-overlay {
  opacity: 0 !important;
  visibility: hidden;
  pointer-events: none;
  transition: opacity 0.4s ease;
}
.animating .tent {
  visibility: visible;
  animation: tent-appear var(--finalDuration) ease-out
    calc(var(--finalStart) + 0.25s) forwards;
}
@keyframes tent-appear {
  0% {
    opacity: 0;
    transform: translate(30%, 40%) scale(0.85);
  }
  20% {
    opacity: 1;
  }
  100% {
    opacity: 1;
    transform: translate(0, 0) scale(1);
  }
}

/* ===== Arrived state ===== */
.arrived .final-tree {
  opacity: 1;
  visibility: visible;
  left: var(--to-left);
  transform: translateX(-50%) scale(1);
}
.arrived .tent {
  opacity: 1;
  visibility: visible;
  transform: translate(0, 0) scale(1);
}

.hero.story-mode .moon,
.hero.story-mode .tree-static,
.hero.story-mode .final-tree,
.hero.story-mode .tent,
.hero.story-mode .lottie-overlay {
  opacity: 0;
  pointer-events: none;
  transition: opacity 0.6s ease;
}

/* Categories overlay page */
.categories-container {
  position: fixed;
  inset: 0;
  display: grid;
  place-items: center;
  z-index: 100;
}

/* Lottie overlay */
.lottie-overlay {
  position: absolute;
  left: 12%;
  bottom: 22%;
  transform: translateX(-50%);
  z-index: 120;
  pointer-events: none;
  transition: opacity 0.6s ease;
}

@media (max-width: 640px) {
  .lottie-overlay {
    display: none;
  }
}

@media (min-width: 1024px) and (max-width: 1439px) {
  .lottie-overlay {
    left: 12%;
    bottom: 21%;
  }
}

@media (min-width: 1440px) {
  .lottie-overlay {
    left: 10%;
    bottom: 22%;
  }
}

/* Small screens (phones) */
@media (max-width: 640px) {
  .moon {
    width: 5rem !important;
    top: 1rem !important;
    left: 1rem !important;
  }

  .static-trees {
    display: none !important;
  }

  /* Force bottom anchor and origin for proper scaling */
  .tree-static, .final-tree {
    transform-origin: bottom center !important;
    bottom: 0 !important;
  }

  #final-tree1 {
    --to-left: 8% !important;
    transform: translateX(-50%) scale(0.45) !important;
  }
  #final-tree2 {
    --to-left: 92% !important;
    transform: translateX(-50%) scale(0.45) !important;
  }
  #final-tree3 {
    --to-left: 96% !important;
    transform: translateX(-50%) scale(0.45) !important;
  }

  /* Static trees smaller */
  #tree1 {
    left: -2rem !important;
    transform: scale(0.45) !important;
  }
  #tree2 {
    right: -1rem !important;
    transform: scale(0.45) !important;
  }
  #tree3 {
    right: -2rem !important;
    transform: scale(0.45) !important;
  }
  #tree4, #tree5, #tree6 {
    transform: scale(0.4) !important;
  }

  .tent {
    display: none !important;
  }
  
  .arrived .tent {
    display: none !important;
  }
  
  @keyframes final-move {
    0% {
      opacity: 0;
      left: var(--from-left, 50%);
      transform: translateX(-50%) scale(0.4);
    }
    25% {
      opacity: 0.9;
    }
    100% {
      opacity: 1;
      left: var(--to-left);
      transform: translateX(-50%) scale(0.45);
    }
  }
}

/* Medium screens (tablets) */
@media (min-width: 641px) and (max-width: 1024px) {
  #final-tree1 {
    --to-left: 14% !important;
  }
  #final-tree2 {
    --to-left: 86% !important;
  }
  #final-tree3 {
    --to-left: 93% !important;
  }

  .moon {
    width: 36vw !important;
  }

  #tree1 {
    left: -8rem !important;
  }
  #tree2 {
    right: -4rem !important;
  }
  #tree3 {
    right: -8rem !important;
  }

  .tent {
    right: 16% !important;
  }
}

/* Very large screens (desktops/TV) */
@media (min-width: 1440px) {
  #final-tree1 {
    --to-left: 10% !important;
    transform: translateX(-50%) scale(1.05);
  }
  #final-tree2 {
    --to-left: 90% !important;
    transform: translateX(-50%) scale(1.05);
  }
  #final-tree3 {
    --to-left: 96% !important;
    transform: translateX(-50%) scale(1.05);
  }
}
</style>
