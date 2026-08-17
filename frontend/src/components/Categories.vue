<template>
  <div class="container mx-auto max-w-5xl px-5 py-4 pb-16 lg:py-8">
    <transition name="fade-story" mode="out-in">
      <div
        v-if="!activeStory"
        key="hub"
        class="relative overflow-hidden rounded-3xl glass-surface border border-white/10 p-5 lg:p-8"
      >
        <div
          class="absolute inset-0 bg-gradient-to-br from-white/2 via-white/0 to-white/5"
        ></div>
        <div class="relative z-10 flex flex-col gap-6">
          <div
            class="flex flex-col lg:flex-row items-start gap-6 lg:items-center justify-between"
          >
            <div class="flex items-center gap-4">
              <button
                @click="goBack"
                class="btn btn-circle btn-ghost border border-white/10 text-slate-200"
                aria-label="Go back"
              >
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  class="h-5 w-5"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M10 19l-7-7m0 0l7-7m-7 7h18"
                  />
                </svg>
              </button>
              <div>
                <p class="text-sm uppercase tracking-[0.2em] text-slate-400">
                  Salon des histoires
                </p>
                <h2 class="text-3xl md:text-4xl font-bold text-white">
                  {{
                    showSavedStories
                      ? "Mes aventures sauvegardées"
                      : "Choisis l'aventure de ce soir"
                  }}
                </h2>
              </div>
            </div>

            <button
              @click="toggleSavedStories"
              class="btn btn-ghost gap-2 text-slate-200 border border-white/10 hover:bg-white/10"
              :class="{ 'bg-white/10 text-white': showSavedStories }"
            >
              <BookOpen class="w-4 h-4" />
              <span class="hidden sm:inline">{{
                showSavedStories ? "Explorer" : "Histoires sauvegardées"
              }}</span>
            </button>
          </div>

          <div
            v-if="showSavedStories"
            class="grid gap-4 max-h-[65vh] overflow-y-auto pr-1"
          >
            <div
              v-if="savedStories.length === 0"
              class="text-center py-12 text-slate-400"
            >
              <p>Aucune histoire gardée pour le moment.</p>
            </div>
            <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <article
                v-for="story in savedStories"
                :key="story.id"
                class="relative rounded-2xl border border-white/10 p-5 shadow-lg cursor-pointer hover:bg-white/5 transition bg-black/20"
                @click="openSavedStory(story)"
              >
                <h3 class="text-lg font-semibold text-white mb-1">
                  {{ story.title }}
                </h3>
                <div class="flex gap-2 text-xs text-slate-400 mb-3">
                  <span>{{ formatDate(story.createdAt) }}</span>
                  <span>•</span>
                  <span>{{ story.length }}</span>
                </div>
                <p class="text-sm text-slate-300 line-clamp-3">
                  {{ story.story }}
                </p>
              </article>
            </div>
          </div>

          <div
            v-else
            class="grid gap-4 max-h-[65vh] overflow-y-auto pr-1"
            ref="listContainer"
          >
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4 lg:gap-5">
              <article
                v-for="category in curatedCategories"
                :key="category.name"
                class="relative rounded-2xl border border-white/10 p-4 shadow-xl overflow-hidden cursor-pointer transition hover:-translate-y-0.5 hover:border-white/20 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-400/60"
                :class="category.bg"
                role="button"
                tabindex="0"
                @click="toggleExpand(category)"
                @keyup.enter="toggleExpand(category)"
                :ref="(el) => setCardRef(category.name, el)"
              >
                <div
                  class="absolute inset-0 bg-gradient-to-br from-white/5 via-transparent to-white/0 pointer-events-none"
                ></div>

                <transition name="fade-card" mode="out-in">
                  <template v-if="!isExpanded(category)">
                    <div class="relative z-10 flex flex-col gap-3">
                      <div class="flex items-center justify-between">
                        <div>
                          <p
                            class="text-xs uppercase tracking-[0.2em] text-slate-300"
                          >
                            {{ category.pill }}
                          </p>
                          <h3 class="text-xl font-semibold text-white">
                            {{ category.name }}
                          </h3>
                        </div>
                        <div
                          class="px-3 py-1 rounded-full bg-black/30 border border-white/10 text-xs text-slate-200"
                        >
                          {{ category.length }}
                        </div>
                      </div>

                      <p class="text-sm text-slate-200 leading-relaxed">
                        {{ category.description }}
                      </p>

                      <div
                        class="flex items-center gap-3 pt-3 text-sm border-t border-white/10"
                      >
                        <div
                          class="h-10 w-10 rounded-xl bg-white/10 border border-white/10 flex items-center justify-center"
                        >
                          <span class="text-lg">{{ category.icon }}</span>
                        </div>
                        <div
                          class="flex-1 text-xs text-slate-200 font-semibold text-right"
                        >
                          Appuie pour ouvrir
                        </div>
                      </div>
                    </div>
                  </template>
                  <template v-else>
                    <div class="relative z-10 flex flex-col gap-3">
                      <div class="flex items-start justify-between gap-3">
                        <div>
                          <p
                            class="text-xs uppercase tracking-[0.2em] text-slate-300"
                          >
                            {{ category.pill }}
                          </p>
                          <h3 class="text-xl font-semibold text-white">
                            {{ category.name }}
                          </h3>
                          <p
                            class="text-sm text-slate-200 leading-relaxed mt-2"
                          >
                            {{ category.description }}
                          </p>
                        </div>
                        <button
                          @click.stop="clearSelection"
                          class="btn btn-circle btn-ghost btn-sm text-slate-400 hover:text-white -mt-1 -mr-1 flex-shrink-0"
                        >
                          <X class="w-5 h-5" />
                        </button>
                      </div>

                      <div
                        class="grid grid-cols-2 gap-2 text-xs text-slate-100/90"
                      >
                        <span
                          class="px-3 py-2 rounded-xl bg-white/10 border border-white/10"
                        >
                          Ton : {{ category.tone }}
                        </span>
                        <span
                          class="px-3 py-2 rounded-xl bg-white/10 border border-white/10"
                        >
                          Décor : {{ category.setting }}
                        </span>
                        <span
                          class="px-3 py-2 rounded-xl bg-white/10 border border-white/10"
                        >
                          Compagnon : {{ category.companion }}
                        </span>

                        <!-- Duration Display / Editor -->
                        <div
                          v-if="editingDuration"
                          class="col-span-2 flex flex-col gap-3 px-3 py-3 rounded-xl bg-white/10 border border-white/10 animate-fade-in sm:flex-row sm:items-center"
                          @click.stop
                        >
                          <div class="flex-1 flex flex-col gap-2">
                            <div
                              class="flex items-center justify-between text-xs text-slate-200"
                            >
                              <span>Durée</span>
                              <span class="font-mono"
                                >{{ durationValue }} min</span
                              >
                            </div>
                            <input
                              type="range"
                              min="5"
                              max="30"
                              v-model.number="durationValue"
                              @input="updateDurationFromSlider"
                              class="w-full text-primary h-8 touch-none"
                            />
                          </div>
                          <div class="flex gap-2 sm:flex-col sm:gap-1">
                            <button
                              @click.stop="confirmDuration"
                              class="btn btn-primary btn-sm w-full sm:w-auto sm:px-4"
                            >
                              Valider
                            </button>
                            <button
                              @click.stop="clearSelection"
                              class="btn btn-ghost btn-sm w-full sm:w-auto"
                            >
                              Annuler
                            </button>
                          </div>
                        </div>

                        <button
                          v-else
                          class="group col-span-2 px-3 py-2 rounded-xl bg-white/10 border border-white/10 text-left hover:bg-white/20 transition active:scale-95 select-none flex items-center justify-between"
                          @click.stop="enableDurationEdit"
                          title="Changer la durée"
                        >
                          <span class="flex items-center gap-2">
                            <Clock class="w-3 h-3 inline-block -mt-0.5" />
                            {{ selectedDuration }}
                          </span>
                          <Pencil class="w-3 h-3 text-slate-300" />
                        </button>
                      </div>

                      <div class="flex gap-3 pt-2">
                        <button
                          class="btn btn-primary btn-sm shadow-md shadow-cyan-500/20 w-full"
                          @click.stop="startStory(category)"
                        >
                          Lancer l'histoire
                        </button>
                      </div>
                    </div>
                  </template>
                </transition>
              </article>
            </div>
          </div>
        </div>
      </div>
      <div v-else key="story">
        <div class="story-stage story-overlay-active">
          <StoryPage
            :category="activeStory"
            :story-content="storyText"
            :story-title="storyTitle"
            :loading="storyLoading"
            :error="storyError"
            :estimated-minutes="storyEstimate"
            class="story-fullscreen"
            @close="closeStory"
            @retry="activeStory && startStory(activeStory)"
          />
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, type ComponentPublicInstance } from "vue";
import { BookOpen, Clock, Pencil, X } from "lucide-vue-next";
import StoryPage from "./StoryPage.vue";

const emit = defineEmits(["go-back", "story-overlay"]);

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

type SavedStory = {
  id: string;
  title: string;
  story: string;
  tone: string;
  length: string;
  createdAt: string;
};

const expandedCategory = ref<Category | null>(null);
const activeStory = ref<Category | null>(null);
const storyLoading = ref(false);
const storyError = ref<string | null>(null);
const storyText = ref<string | null>(null);
const storyTitle = ref<string | null>(null);
const cardRefs = ref<Record<string, HTMLElement | null>>({});
const listContainer = ref<HTMLElement | null>(null);
const showSavedStories = ref(false);
const savedStories = ref<SavedStory[]>([]);
const durationOptions: string[] = [
  "Courte (3-5 min)",
  "Moyenne (5-8 min)",
  "Longue (8-12 min)",
];
const selectedDuration = ref(durationOptions[1]!); // Default to Moyenne
const editingDuration = ref(false);
const durationValue = ref(8);

const enableDurationEdit = () => {
  editingDuration.value = true;
  
  // Try to parse range first (e.g., "5-8 min")
  const rangeMatch = selectedDuration.value.match(/(\d+)\s*-\s*(\d+)/);
  if (rangeMatch) {
    const min = parseInt(rangeMatch[1]!);
    const max = parseInt(rangeMatch[2]!);
    const avg = Math.round((min + max) / 2);
    durationValue.value = Math.max(5, Math.min(30, avg));
    return;
  }

  // Fallback to single number
  const match = selectedDuration.value.match(/(\d+)/);
  if (match) {
    const val = parseInt(match[0]);
    durationValue.value = Math.max(5, Math.min(30, val));
  } else {
    durationValue.value = 8;
  }
};

const updateDurationFromSlider = () => {
  selectedDuration.value = `${durationValue.value} min`;
};

const confirmDuration = () => {
  editingDuration.value = false;
};

const serverAddr = import.meta.env.VITE_SERVER_ADDR || "";
const defaultBase = serverAddr
  ? `http://backend${serverAddr}`
  : `${window.location.origin}/api`;
const apiBase =
  import.meta.env.VITE_API_BASE_URL ||
  import.meta.env.VITE_API_BASE ||
  defaultBase;
const storyEstimate = computed(() => {
  if (!activeStory.value) return null;
  const val = durationValue.value;
  if (val <= 10) return 1;
  if (val <= 20) return 2;
  return 3;
});

const curatedCategories: Category[] = [
  {
    name: "Lanternes au port des bassins",
    description:
      "Douce brise océane, quais scintillants et chants de baleine apaisants qui accompagnent vers le sommeil.",
    tone: "Calme",
    length: "5-7 min",
    prompt:
      "Une lanterne garde une étoile allumée jusqu'à ce que ton rêve arrive",
    companion: "Crabe lumineux",
    setting: "Quais côtiers",
    pill: "Océan douillet",
    bg: "bg-gradient-to-br from-cyan-500/20 via-blue-900/40 to-indigo-900/50",
    icon: "🌊",
  },
  {
    name: "Le tramway de la bibliothèque de minuit",
    description:
      "Un tram lent glisse entre des bibliothèques flottantes où chaque livre murmure une berceuse en tournant les pages.",
    tone: "Merveille",
    length: "6-8 min",
    prompt: "Une chouette bibliothécaire tamponne un billet en clair de lune",
    companion: "Chouette bibliothécaire",
    setting: "Voies dans le ciel",
    pill: "Aventure douce",
    bg: "bg-gradient-to-br from-violet-500/20 via-purple-900/40 to-slate-900/60",
    icon: "🚋",
  },
  {
    name: "Camp du verger de nuages",
    description:
      "Des fruits qui brillent comme de petits soleils, des hamacs entre des branches de nuages et des lucioles berceuses qui respirent avec toi.",
    tone: "Léger",
    length: "4-6 min",
    prompt: "Des lucioles apprennent à la tente à respirer avec toi",
    companion: "Trio de lucioles",
    setting: "Verger flottant",
    pill: "Calme respiré",
    bg: "bg-gradient-to-br from-amber-400/25 via-orange-900/40 to-rose-900/50",
    icon: "⛺",
  },
  {
    name: "Signaux de la prairie argentée",
    description:
      "Une prairie remplie de satellites mélodiques et de petits robots qui collectent les vœux tandis que le ciel s'assombrit.",
    tone: "Curieux",
    length: "7-9 min",
    prompt: "De petits robots tapotent leurs antennes au rythme de ton cœur",
    companion: "Robot bricoleur",
    setting: "Champ étoilé",
    pill: "Science douce",
    bg: "bg-gradient-to-br from-cyan-300/20 via-emerald-900/40 to-slate-950/70",
    icon: "🛰️",
  },
  {
    name: "Échos de la caverne lumineuse",
    description:
      "Des cristaux bioluminescents fredonnent des berceuses tandis que des gouttes régulières marquent le temps sous la colline.",
    tone: "Apaisant",
    length: "6-8 min",
    prompt: "Un cristal guide tes pas en calant sa lueur sur ta respiration",
    companion: "Renard lumineux",
    setting: "Grotte cachée",
    pill: "Découverte calme",
    bg: "bg-gradient-to-br from-emerald-400/20 via-teal-900/40 to-slate-950/70",
    icon: "🦊",
  },
  {
    name: "Cerfs-volants d'aurore au champ givré",
    description:
      "Des cerfs-volants de papier attrapent des rubans d'aurore, planant au-dessus de collines de neige douce et de feux de cacao.",
    tone: "Rêveur",
    length: "5-7 min",
    prompt: "Un cerf-volant trace des couleurs qui écrivent ton vœu",
    companion: "Lièvre des neiges",
    setting: "Prairie givrée",
    pill: "Hiver doux",
    bg: "bg-gradient-to-br from-sky-400/20 via-blue-900/40 to-indigo-950/70",
    icon: "🪁",
  },
];

const goBack = () => {
  if (showSavedStories.value) {
    showSavedStories.value = false;
  } else {
    emit("go-back");
  }
};

const toggleSavedStories = async () => {
  showSavedStories.value = !showSavedStories.value;
  if (showSavedStories.value) {
    await fetchSavedStories();
  }
};

const fetchSavedStories = async () => {
  try {
    const res = await fetch(`${apiBase}/stories`);
    if (res.ok) {
      savedStories.value = await res.json();
    }
  } catch (e) {
    console.error("Failed to fetch stories", e);
  }
};

const formatDate = (dateStr: string) => {
  return new Date(dateStr).toLocaleDateString("fr-FR", {
    day: "numeric",
    month: "short",
  });
};

const openSavedStory = (story: SavedStory) => {
  activeStory.value = {
    name: story.title,
    description: "",
    tone: story.tone,
    length: story.length,
    prompt: "",
    companion: "",
    setting: "",
    pill: "Histoire gardée",
    bg: "bg-gradient-to-br from-slate-700 via-slate-800 to-slate-900",
    icon: "📖",
  };
  storyTitle.value = story.title;
  storyText.value = story.story;
  storyError.value = null;
  storyLoading.value = false;
  emit("story-overlay", true);
};

const toggleExpand = (category: Category) => {
  if (activeStory.value) return;
  const next = expandedCategory.value?.name === category.name ? null : category;
  expandedCategory.value = next;

  if (next) {
    // Reset duration to default (Moyenne)
    selectedDuration.value = durationOptions[1]!;
    editingDuration.value = false;

    // Wait for the "out-in" transition to swap content (leave 200ms + render)
    // The card height updates only after the old content leaves.
    setTimeout(() => {
      const el = cardRefs.value[next.name];
      scrollCardIntoView(el);
    }, 300);
  }
};

const clearSelection = () => {
  expandedCategory.value = null;
};

const isExpanded = (category: Category) =>
  expandedCategory.value?.name === category.name;

const startStory = async (category: Category) => {
  activeStory.value = { ...category }; // Clone to avoid mutating the source list
  // Use the selected duration instead of category.length
  activeStory.value.length = selectedDuration.value;

  expandedCategory.value = null;
  storyLoading.value = true;
  storyError.value = null;
  storyText.value = null;
  storyTitle.value = category.name;
  emit("story-overlay", true);
  window.scrollTo({ top: 0, behavior: "smooth" });

  const payload = {
    title: category.name,
    tone: category.tone,
    setting: category.setting,
    companion: category.companion,
    // Make length more explicit for the AI
    length: `${selectedDuration.value} reading time`,
  };

  try {
    const res = await fetch(`${apiBase}/story`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(payload),
    });

    if (!res.ok) {
      throw new Error(`Request failed with status ${res.status}`);
    }

    const data = await res.json();
    storyTitle.value = data.title || category.name;
    storyText.value = data.story || data.text || data.content || "";

    // Refresh saved stories in background so list is up to date when they go back
    fetchSavedStories();
  } catch (err) {
    console.error("Story request failed", err);
    storyError.value =
      "Impossible de récupérer ton histoire pour le moment. Réessaie dans un instant.";
  } finally {
    storyLoading.value = false;
  }
};

const closeStory = () => {
  activeStory.value = null;
  expandedCategory.value = null;
  storyLoading.value = false;
  storyError.value = null;
  storyText.value = null;
  storyTitle.value = null;
  emit("story-overlay", false);
};

const setCardRef = (
  name: string,
  el: Element | null | ComponentPublicInstance | undefined
) => {
  const dom = (el as ComponentPublicInstance | null)?.$el ?? el;
  cardRefs.value[name] = (dom as HTMLElement | null) ?? null;
};

const scrollCardIntoView = (el?: HTMLElement | null) => {
  const container = listContainer.value;
  if (!el || !container) return;

  // Calculate positions relative to the scrolling container
  const targetTop = el.offsetTop - container.offsetTop;
  const targetBottom = targetTop + el.clientHeight;

  const containerHeight = container.clientHeight;
  const currentScroll = container.scrollTop;

  // Center alignment target
  let desired = targetTop - containerHeight / 2 + el.clientHeight / 2;

  // Adjustment: Ensure the bottom of the card is visible (with padding)
  // especially important for the last items or tall expanded cards
  const padding = 20;
  if (targetBottom > currentScroll + containerHeight - padding) {
    // If bottom is cut off, try to scroll down to show it
    // But don't scroll so far that the top is hidden, unless the card is taller than the viewport
    const alignBottom = targetBottom - containerHeight + padding;
    desired = Math.max(desired, alignBottom);
  }

  container.scrollTo({
    top: Math.max(0, desired),
    behavior: "smooth",
  });
};
</script>

<style scoped>
:deep(.fade-story-enter-active),
:deep(.fade-story-leave-active) {
  transition: opacity 0.35s ease, transform 0.35s ease;
}
:deep(.fade-story-enter-from),
:deep(.fade-story-leave-to) {
  opacity: 0;
  transform: translateY(12px);
}

:deep(.fade-card-enter-active),
:deep(.fade-card-leave-active) {
  transition: opacity 0.2s ease, transform 0.2s ease;
}
:deep(.fade-card-enter-from),
:deep(.fade-card-leave-to) {
  opacity: 0;
  transform: translateY(6px);
}

:deep(.story-fullscreen) {
  max-height: 80vh;
  overflow-y: auto;
  width: min(90vw, 820px);
}

.story-stage {
  position: fixed;
  inset: 0;
  display: grid;
  place-items: center;
  z-index: 120;
  background: transparent;
}

/* Custom range slider styling for better touch targets */
input[type="range"] {
  -webkit-appearance: none;
  appearance: none;
  background: transparent;
  cursor: pointer;
}

/* Webkit (Chrome, Safari, iOS) */
input[type="range"]::-webkit-slider-thumb {
  -webkit-appearance: none;
  height: 20px;
  width: 20px;
  border-radius: 50%;
  background: currentColor;
  margin-top: -8px; /* Center thumb on track */
  box-shadow: 0 0 0 4px rgba(255, 255, 255, 0.1);
}

input[type="range"]::-webkit-slider-runnable-track {
  width: 100%;
  height: 4px;
  cursor: pointer;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 2px;
}

/* Firefox */
input[type="range"]::-moz-range-thumb {
  height: 20px;
  width: 20px;
  border: none;
  border-radius: 50%;
  background: currentColor;
  box-shadow: 0 0 0 4px rgba(255, 255, 255, 0.1);
}

input[type="range"]::-moz-range-track {
  width: 100%;
  height: 4px;
  cursor: pointer;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 2px;
}
</style>
