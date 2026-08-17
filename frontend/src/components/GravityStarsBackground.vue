<template>
  <div
    ref="containerRef"
    data-slot="gravity-stars-background"
    class="relative w-full h-full overflow-hidden"
    @mousemove="handleMouseMove"
    @touchmove.passive="handleTouchMove"
  >
    <canvas ref="canvasRef" class="block w-full h-full" />
  </div>
</template>

<script setup lang="ts">
import {
  ref,
  onMounted,
  onBeforeUnmount,
  withDefaults,
  defineProps,
  watch,
} from "vue";

type MouseGravity = "attract" | "repel";
type GlowAnimation = "instant" | "ease" | "spring";
type StarsInteractionType = "bounce" | "merge";

interface GravityStarsProps {
  starsCount?: number;
  starsSize?: number;
  starsOpacity?: number;
  glowIntensity?: number;
  glowAnimation?: GlowAnimation;
  movementSpeed?: number;
  mouseInfluence?: number;
  mouseGravity?: MouseGravity;
  gravityStrength?: number;
  starsInteraction?: boolean;
  starsInteractionType?: StarsInteractionType;
}

interface Particle {
  x: number;
  y: number;
  vx: number;
  vy: number;
  size: number;
  opacity: number;
  baseOpacity: number;
  mass: number;
  glowMultiplier: number;
  glowVelocity: number;
}

const props = withDefaults(defineProps<GravityStarsProps>(), {
  starsCount: 75,
  starsSize: 2,
  starsOpacity: 0.75,
  glowIntensity: 15,
  glowAnimation: "ease",
  movementSpeed: 0.3,
  mouseInfluence: 100,
  mouseGravity: "attract",
  gravityStrength: 75,
  starsInteraction: false,
  starsInteractionType: "bounce",
});

const containerRef = ref<HTMLDivElement | null>(null);
const canvasRef = ref<HTMLCanvasElement | null>(null);
const animId = ref<number | null>(null);
const stars = ref<Particle[]>([]);
const mouse = ref<{ x: number; y: number }>({ x: -9999, y: -9999 }); // far away initially

const dpr = ref(1);
const canvasSize = ref({ width: 800, height: 600 });

let resizeObserver: ResizeObserver | null = null;

const resetMousePosition = () => {
  mouse.value = { x: -9999, y: -9999 };
};

const updateMouseFromClient = (clientX: number, clientY: number) => {
  const canvas = canvasRef.value;
  if (!canvas) return;
  const rect = canvas.getBoundingClientRect();
  mouse.value = {
    x: clientX - rect.left,
    y: clientY - rect.top,
  };
};

const readColor = (): string => {
  const el = containerRef.value;
  if (!el) return "#ffffff";
  const cs = getComputedStyle(el);
  // We use the text color of the container as star color
  return cs.color || "#ffffff";
};

const initStars = (w: number, h: number) => {
  const count = props.starsCount;
  const newStars: Particle[] = [];

  for (let i = 0; i < count; i++) {
    const angle = Math.random() * Math.PI * 2;
    const speed = props.movementSpeed * (0.5 + Math.random() * 0.5);
    newStars.push({
      x: Math.random() * w,
      y: Math.random() * h,
      vx: Math.cos(angle) * speed,
      vy: Math.sin(angle) * speed,
      size: Math.random() * props.starsSize + 1,
      opacity: props.starsOpacity,
      baseOpacity: props.starsOpacity,
      mass: Math.random() * 0.5 + 0.5,
      glowMultiplier: 1,
      glowVelocity: 0,
    });
  }

  stars.value = newStars;
};

const redistributeStars = (w: number, h: number) => {
  for (const p of stars.value) {
    p.x = Math.random() * w;
    p.y = Math.random() * h;
  }
};

const resizeCanvas = () => {
  const canvas = canvasRef.value;
  const container = containerRef.value;
  if (!canvas || !container) return;

  const rect = container.getBoundingClientRect();
  const nextDpr = Math.max(1, Math.min(window.devicePixelRatio || 1, 2));
  dpr.value = nextDpr;

  canvas.width = Math.max(1, Math.floor(rect.width * nextDpr));
  canvas.height = Math.max(1, Math.floor(rect.height * nextDpr));
  canvas.style.width = `${rect.width}px`;
  canvas.style.height = `${rect.height}px`;

  canvasSize.value = { width: rect.width, height: rect.height };

  if (stars.value.length === 0) {
    initStars(rect.width, rect.height);
  } else {
    redistributeStars(rect.width, rect.height);
  }
};

const handleMouseMove = (event: MouseEvent) =>
  updateMouseFromClient(event.clientX, event.clientY);

const handleTouchMove = (event: TouchEvent) => {
  const t = event.touches[0];
  if (!t) return;

  updateMouseFromClient(t.clientX, t.clientY);
};

const updateStars = () => {
  const w = canvasSize.value.width;
  const h = canvasSize.value.height;
  const m = mouse.value;
  const list = stars.value;

  for (let i = 0; i < list.length; i++) {
    const p = list[i];
    if (!p) continue; // TS safety

    const dx = m.x - p.x;
    const dy = m.y - p.y;
    const dist = Math.hypot(dx, dy);

    // Mouse influence
    if (dist < props.mouseInfluence && dist > 0) {
      const force = (props.mouseInfluence - dist) / props.mouseInfluence;
      const nx = dx / dist;
      const ny = dy / dist;
      // slightly boosted strength for a more visible effect
      const g = force * (props.gravityStrength * 0.0015);

      if (props.mouseGravity === "attract") {
        p.vx += nx * g;
        p.vy += ny * g;
      } else if (props.mouseGravity === "repel") {
        p.vx -= nx * g;
        p.vy -= ny * g;
      }

      // Highlight / glow when close to mouse
      p.opacity = Math.min(1, p.baseOpacity + force * 0.5);

      const targetGlow = 1 + force * 2.5;
      const currentGlow = p.glowMultiplier || 1;

      if (props.glowAnimation === "instant") {
        p.glowMultiplier = targetGlow;
      } else if (props.glowAnimation === "ease") {
        const ease = 0.15;
        p.glowMultiplier = currentGlow + (targetGlow - currentGlow) * ease;
      } else {
        const spring = (targetGlow - currentGlow) * 0.2;
        const damping = 0.85;
        p.glowVelocity = (p.glowVelocity || 0) * damping + spring;
        p.glowMultiplier = currentGlow + (p.glowVelocity || 0);
      }
    } else {
      // Fade back to base opacity & glow
      p.opacity = Math.max(p.baseOpacity * 0.3, p.opacity - 0.02);
      const targetGlow = 1;
      const currentGlow = p.glowMultiplier || 1;

      if (props.glowAnimation === "instant") {
        p.glowMultiplier = targetGlow;
      } else if (props.glowAnimation === "ease") {
        const ease = 0.08;
        p.glowMultiplier = Math.max(
          1,
          currentGlow + (targetGlow - currentGlow) * ease
        );
      } else {
        const spring = (targetGlow - currentGlow) * 0.15;
        const damping = 0.9;
        p.glowVelocity = (p.glowVelocity || 0) * damping + spring;
        p.glowMultiplier = Math.max(1, currentGlow + (p.glowVelocity || 0));
      }
    }

    // Star-star interaction (optional)
    if (props.starsInteraction) {
      for (let j = i + 1; j < list.length; j++) {
        const o = list[j];
        if (!o) continue;

        const dx2 = o.x - p.x;
        const dy2 = o.y - p.y;
        const d = Math.hypot(dx2, dy2);
        const minD = p.size + o.size + 5;

        if (d < minD && d > 0) {
          if (props.starsInteractionType === "bounce") {
            const nx = dx2 / d;
            const ny = dy2 / d;
            const rvx = p.vx - o.vx;
            const rvy = p.vy - o.vy;
            const speed = rvx * nx + rvy * ny;
            if (speed < 0) continue;
            const impulse = (2 * speed) / (p.mass + o.mass);
            p.vx -= impulse * o.mass * nx;
            p.vy -= impulse * o.mass * ny;
            o.vx += impulse * p.mass * nx;
            o.vy += impulse * p.mass * ny;
            const overlap = minD - d;
            const sx = nx * overlap * 0.5;
            const sy = ny * overlap * 0.5;
            p.x -= sx;
            p.y -= sy;
            o.x += sx;
            o.y += sy;
          } else {
            const mergeForce = (minD - d) / minD;
            p.glowMultiplier = (p.glowMultiplier || 1) + mergeForce * 0.5;
            o.glowMultiplier = (o.glowMultiplier || 1) + mergeForce * 0.5;
            const af = mergeForce * 0.01;
            p.vx += dx2 * af;
            p.vy += dy2 * af;
            o.vx -= dx2 * af;
            o.vy -= dy2 * af;
          }
        }
      }
    }

    // Motion integration
    p.x += p.vx;
    p.y += p.vy;

    // Tiny jitter & damping
    p.vx += (Math.random() - 0.5) * 0.001;
    p.vy += (Math.random() - 0.5) * 0.001;

    p.vx *= 0.999;
    p.vy *= 0.999;

    // Wrap around edges
    if (p.x < 0) p.x = w;
    if (p.x > w) p.x = 0;
    if (p.y < 0) p.y = h;
    if (p.y > h) p.y = 0;
  }
};

const drawStars = (ctx: CanvasRenderingContext2D) => {
  ctx.clearRect(0, 0, ctx.canvas.width, ctx.canvas.height);
  const color = readColor();

  for (const p of stars.value) {
    ctx.save();
    ctx.shadowColor = color;
    ctx.shadowBlur = props.glowIntensity * (p.glowMultiplier || 1) * 2;
    ctx.globalAlpha = p.opacity;
    ctx.fillStyle = color;
    ctx.beginPath();
    ctx.arc(
      p.x * dpr.value,
      p.y * dpr.value,
      p.size * dpr.value,
      0,
      Math.PI * 2
    );
    ctx.fill();
    ctx.restore();
  }
};

const animate = () => {
  const canvas = canvasRef.value;
  if (!canvas) return;
  const ctx = canvas.getContext("2d");
  if (!ctx) return;

  updateStars();
  drawStars(ctx);

  animId.value = requestAnimationFrame(animate);
};

const handleWindowResize = () => resizeCanvas();
const handlePointerMoveGlobal = (event: PointerEvent) =>
  updateMouseFromClient(event.clientX, event.clientY);
const handleTouchMoveGlobal = (event: TouchEvent) => {
  const t = event.touches[0];
  if (!t) return;
  updateMouseFromClient(t.clientX, t.clientY);
};
const handlePointerLeave = () => resetMousePosition();

onMounted(() => {
  resizeCanvas();

  const container = containerRef.value;
  if (typeof ResizeObserver !== "undefined" && container) {
    resizeObserver = new ResizeObserver(() => resizeCanvas());
    resizeObserver.observe(container);
  }

  window.addEventListener("resize", handleWindowResize);
  window.addEventListener("pointermove", handlePointerMoveGlobal);
  window.addEventListener("pointerleave", handlePointerLeave);
  window.addEventListener("touchmove", handleTouchMoveGlobal, { passive: true });

  if (stars.value.length === 0) {
    initStars(canvasSize.value.width, canvasSize.value.height);
  }

  animId.value = requestAnimationFrame(animate);
});

onBeforeUnmount(() => {
  window.removeEventListener("resize", handleWindowResize);
  window.removeEventListener("pointermove", handlePointerMoveGlobal);
  window.removeEventListener("pointerleave", handlePointerLeave);
  window.removeEventListener("touchmove", handleTouchMoveGlobal);

  if (resizeObserver) {
    resizeObserver.disconnect();
    resizeObserver = null;
  }

  if (animId.value) {
    cancelAnimationFrame(animId.value);
    animId.value = null;
  }
});

// React to opacity/speed changes at runtime
watch(
  () => [props.starsOpacity, props.movementSpeed] as const,
  () => {
    for (const p of stars.value) {
      p.baseOpacity = props.starsOpacity;
      p.opacity = props.starsOpacity;
      const spd = Math.hypot(p.vx, p.vy);
      if (spd > 0) {
        const ratio = props.movementSpeed / spd;
        p.vx *= ratio;
        p.vy *= ratio;
      }
    }
  }
);

// ⭐ React to starsCount changes (allow parent to control star count)
watch(
  () => props.starsCount,
  () => {
    initStars(canvasSize.value.width, canvasSize.value.height);
  }
);
</script>

<style scoped>
/* All visuals are in canvas, nothing special needed here */
</style>
