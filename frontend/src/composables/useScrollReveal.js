// Composable untuk mengaktifkan animasi scroll reveal menggunakan IntersectionObserver
import { onMounted, onBeforeUnmount } from 'vue'

export function useScrollReveal(selector = '.reveal-on-scroll', options = {}) {
  let observer = null

  const defaultOptions = {
    threshold: 0.12,
    rootMargin: '0px 0px -40px 0px',
    ...options,
  }

  onMounted(() => {
    if (typeof window === 'undefined' || !('IntersectionObserver' in window)) {
      // Fallback untuk browser lawas: langsung tampilkan
      document.querySelectorAll(selector).forEach((el) => {
        el.classList.add('is-revealed')
      })
      return
    }

    observer = new IntersectionObserver((entries, obs) => {
      entries.forEach((entry) => {
        if (entry.isIntersecting) {
          entry.target.classList.add('is-revealed')
          obs.unobserve(entry.target)
        }
      })
    }, defaultOptions)

    const elements = document.querySelectorAll(selector)
    elements.forEach((el) => observer.observe(el))
  })

  onBeforeUnmount(() => {
    if (observer) {
      observer.disconnect()
      observer = null
    }
  })
}
