<template>
  <div
    class="rounded-2xl border border-white/10 p-4 sm:p-6 bg-slate-950/80 backdrop-blur-xl shadow-2xl space-y-4 max-h-[82vh] overflow-y-auto w-full max-w-3xl relative z-10"
  >
    <button 
      class="btn btn-circle btn-ghost btn-sm absolute top-2 right-2 text-slate-400 hover:text-white z-10"
      @click.stop="$emit('close')"
    >
      <X class="w-5 h-5" />
    </button>

    <div class="flex flex-col sm:flex-row items-start sm:items-center gap-3 sm:gap-4">
      <div class="flex-1 pr-8">
        <p class="text-xs uppercase tracking-[0.2em] text-slate-400">
          {{ category.pill }}
        </p>
        <h3 class="text-2xl font-semibold text-white leading-tight">
          {{ resolvedTitle }}
        </h3>
        <p class="text-sm sm:text-base text-slate-200 leading-relaxed mt-1">
          {{ category.description }}
        </p>
        <div class="flex flex-wrap gap-2 mt-3 text-xs sm:text-sm text-slate-100/90">
          <span class="px-3 py-1 rounded-full bg-white/10 border border-white/10">
            Ton : {{ category.tone }}
          </span>
          <span class="px-3 py-1 rounded-full bg-white/10 border border-white/10">
            Décor : {{ category.setting }}
          </span>
          <span class="px-3 py-1 rounded-full bg-white/10 border border-white/10">
            Compagnon : {{ category.companion }}
          </span>
          <span class="px-3 py-1 rounded-full bg-white/10 border border-white/10">
            Durée : {{ category.length }}
          </span>
        </div>
      </div>
      <div class="flex flex-col items-end gap-2 self-start sm:self-center mt-2 sm:mt-0 w-full sm:w-auto">
        <button
          v-if="error"
          class="btn btn-outline btn-xs text-slate-200 w-full sm:w-auto"
          @click.stop="$emit('retry')"
        >
          Réessayer
        </button>
        <button
          v-else-if="isPlaceholder"
          class="btn btn-ghost btn-xs text-slate-200 flex items-center gap-1 w-full sm:w-auto justify-center sm:justify-end"
          @click.stop="$emit('retry')"
          title="Relancer la génération"
        >
          <RefreshCw class="w-4 h-4" />
          <span class="inline">Relancer</span>
        </button>
      </div>
    </div>

    <div
      class="rounded-xl border border-white/10 bg-black/20 p-4 text-slate-100 leading-relaxed space-y-3 transition-all duration-200"
      :class="currentFontSize"
    >
      <div class="flex items-center justify-between mb-1 gap-2">
        <p class="font-semibold text-white text-base sm:text-lg">L'histoire de ce soir</p>

        <div class="flex gap-1 bg-white/5 rounded-lg p-0.5 border border-white/10" v-if="!loading && !error">
          <button 
            class="btn btn-ghost btn-xs btn-square h-7 w-7 min-h-0 text-slate-300 hover:bg-white/10 hover:text-white" 
            @click.stop.prevent="decreaseFontSize"
            :disabled="fontSizeIndex === 0"
            title="Diminuer le texte"
          >
            <Minus class="w-3.5 h-3.5" />
          </button>
          <button 
            class="btn btn-ghost btn-xs btn-square h-7 w-7 min-h-0 text-slate-300 hover:bg-white/10 hover:text-white" 
            @click.stop.prevent="increaseFontSize"
            :disabled="fontSizeIndex === fontSizes.length - 1"
            title="Agrandir le texte"
          >
            <Plus class="w-3.5 h-3.5" />
          </button>
        </div>
      </div>

      <template v-if="loading">
        <div class="flex items-center gap-3">
          <span class="loading loading-ring loading-md text-primary"></span>
          <div>
            <p class="text-white font-semibold">On invoque ton histoire...</p>
            <p class="text-xs text-slate-300" v-if="estimatedMinutes">
              Estimation ~{{ estimatedMinutes }} min
            </p>
          </div>
        </div>
        <p class="text-xs text-slate-400">
          Tes choix partent vers les tisserands d'histoires...
        </p>
      </template>

      <template v-else-if="error">
        <p class="text-rose-200">{{ error }}</p>
        <p class="text-xs text-slate-400">Tu peux réessayer dans un instant.</p>
      </template>

      <template v-else>
        <p v-if="heading" class="font-semibold text-white leading-snug mb-4 text-base sm:text-lg">
          {{ heading }}
        </p>
        <p v-for="(para, idx) in bodyParagraphs" :key="idx" class="text-sm sm:text-base">
          {{ para }}
        </p>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from "vue";
import { RefreshCw, Minus, Plus, X } from "lucide-vue-next";

type Category = {
  name: string;
  description: string;
  tone: string;
  length: string;
  prompt: string;
  companion: string;
  setting: string;
  pill: string;
  bg: string;
  icon: string;
};

const props = defineProps<{
  category: Category;
  storyContent?: string | null;
  storyTitle?: string | null;
  loading?: boolean;
  error?: string | null;
  estimatedMinutes?: number | null;
}>();

defineEmits<{
  (e: "close"): void;
  (e: "retry"): void;
}>();

const fontSizes = ["text-sm", "text-base", "text-lg", "text-xl", "text-2xl"];
const fontSizeIndex = ref(0); // Default to text-sm

const currentFontSize = computed(() => fontSizes[fontSizeIndex.value]);

const increaseFontSize = () => {
  if (fontSizeIndex.value < fontSizes.length - 1) {
    fontSizeIndex.value++;
  }
};

const decreaseFontSize = () => {
  if (fontSizeIndex.value > 0) {
    fontSizeIndex.value--;
  }
};

const normalizedText = computed(() => {
  if (props.loading || props.error) return "";
  return (props.storyContent || "").trim();
});

const heading = computed(() => {
  const text = normalizedText.value;
  if (!text.startsWith("###")) return "";
  const firstLine = text.split("\n")[0] || "";
  return firstLine.replace(/^#+\s*/, "").trim();
});

const resolvedTitle = computed(() => {
  return props.storyTitle?.trim() || heading.value || props.category.name;
});

const defaultText =
  "Les conteurs ne sont pas encore prêts, mais une aventure arrive bientôt. Respire doucement, ferme les yeux et imagine le monde qui se dessine pour toi.";

const bodyParagraphs = computed(() => {
  const text = normalizedText.value;
  if (!text) return defaultText.split(/\n/).filter(Boolean);
  const lines = text.split(/\n/);
  const rest =
    text.startsWith("###") && lines.length > 1 ? lines.slice(1).join("\n") : text;
  return rest
    .split(/\n\s*\n/)
    .map((p) => p.trim())
    .filter(Boolean);
});

const isPlaceholder = computed(() => {
  if (props.loading || props.error) return false;
  const content = (props.storyContent || "").trim();
  if (!content) return true;
  return content.includes("Les tisserands d'histoires");
});
</script>

<style scoped>
/* Component-specific styling lives in utility classes above */
</style>
