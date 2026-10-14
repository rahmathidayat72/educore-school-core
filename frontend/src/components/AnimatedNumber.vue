<template>
  <span class="animated-number font-monospace">{{ displayValue }}</span>
</template>

<script setup>
import { ref, watch, onMounted, onBeforeUnmount } from 'vue'

const props = defineProps({
  value: {
    type: [Number, String],
    default: 0,
  },
  duration: {
    type: Number,
    default: 850,
  },
  format: {
    type: Boolean,
    default: true,
  },
})

const displayValue = ref('0')
let animationFrameId = null

const easeOutQuart = (t) => {
  return 1 - Math.pow(1 - t, 4)
}

const animate = (from, to) => {
  if (animationFrameId) {
    cancelAnimationFrame(animationFrameId)
    animationFrameId = null
  }

  const start = performance.now()
  const diff = to - from

  if (diff === 0) {
    displayValue.value = props.format ? to.toLocaleString('id-ID') : to.toString()
    return
  }

  // Jika user menyalakan preferensi reduce-motion, langsung pasang angka target
  if (typeof window !== 'undefined' && window.matchMedia && window.matchMedia('(prefers-reduced-motion: reduce)').matches) {
    displayValue.value = props.format ? to.toLocaleString('id-ID') : to.toString()
    return
  }

  const step = (currentTime) => {
    const elapsed = currentTime - start
    const progress = Math.min(elapsed / props.duration, 1)
    const easedProgress = easeOutQuart(progress)
    const current = Math.round(from + diff * easedProgress)

    displayValue.value = props.format ? current.toLocaleString('id-ID') : current.toString()

    if (progress < 1) {
      animationFrameId = requestAnimationFrame(step)
    } else {
      displayValue.value = props.format ? to.toLocaleString('id-ID') : to.toString()
      animationFrameId = null
    }
  }

  animationFrameId = requestAnimationFrame(step)
}

watch(
  () => props.value,
  (newVal, oldVal) => {
    const from = Number(oldVal) || 0
    const to = Number(newVal) || 0
    animate(from, to)
  }
)

onMounted(() => {
  const target = Number(props.value) || 0
  animate(0, target)
})

onBeforeUnmount(() => {
  if (animationFrameId) {
    cancelAnimationFrame(animationFrameId)
    animationFrameId = null
  }
})
</script>

<style scoped>
.animated-number {
  display: inline-block;
  transition: transform 0.2s cubic-bezier(0.16, 1, 0.3, 1);
}
</style>
