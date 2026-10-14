<template>
  <div class="admin-layout-root d-flex bg-light">
    <!-- Sidebar Desktop & Mobile Wrapper -->
    <div :class="['sidebar-wrapper', { 'd-none d-lg-block': !showMobileSidebar }]">
      <SidebarAdmin />
    </div>

    <!-- Backdrop for mobile sidebar -->
    <div
      v-if="showMobileSidebar"
      class="sidebar-backdrop d-lg-none"
      @click="showMobileSidebar = false"
    ></div>

    <!-- Main Content Area: Navbar stays fixed at top, only main scrolls -->
    <div class="admin-main-wrapper flex-grow-1 d-flex flex-column">
      <NavbarAdmin @toggle-sidebar="showMobileSidebar = !showMobileSidebar" />

      <main ref="mainScrollRef" class="admin-main-scroll flex-grow-1 p-3 p-md-4">
        <div class="container-fluid max-w-1400">
          <router-view v-slot="{ Component, route: childRoute }">
            <transition name="admin-content-fade" mode="out-in">
              <component :is="Component" :key="childRoute.path" />
            </transition>
          </router-view>
        </div>
      </main>
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import SidebarAdmin from './SidebarAdmin.vue'
import NavbarAdmin from './NavbarAdmin.vue'

const showMobileSidebar = ref(false)
const mainScrollRef = ref(null)
const route = useRoute()

// Ketika berpindah fitur/rute admin, tutup drawer mobile jika terbuka dan reset scroll ke atas secara mulus
watch(
  () => route.path,
  () => {
    showMobileSidebar.value = false
    if (mainScrollRef.value) {
      mainScrollRef.value.scrollTo({ top: 0, behavior: 'smooth' })
    }
  }
)
</script>

<style scoped>
.admin-layout-root {
  height: 100vh;
  height: 100dvh;
  width: 100%;
  overflow: hidden;
  position: relative;
}

.sidebar-wrapper {
  z-index: 1040;
  height: 100%;
  flex-shrink: 0;
  overflow-y: auto;
}

@media (max-width: 991.98px) {
  .sidebar-wrapper {
    position: fixed;
    top: 0;
    left: 0;
    bottom: 0;
    height: 100dvh;
  }
}

.sidebar-backdrop {
  position: fixed;
  inset: 0;
  background-color: rgba(15, 23, 42, 0.6);
  backdrop-filter: blur(2px);
  z-index: 1030;
}

.admin-main-wrapper {
  height: 100%;
  min-width: 0;
  overflow: hidden;
  position: relative;
}

.admin-main-scroll {
  flex: 1 1 0%;
  min-height: 0;
  overflow-y: auto;
  overflow-x: hidden;
  scroll-behavior: smooth;
  -webkit-overflow-scrolling: touch;
}

.max-w-1400 {
  max-width: 1400px;
}

/* Transisi Halus Konten Fitur Admin */
.admin-content-fade-enter-active,
.admin-content-fade-leave-active {
  transition: opacity 0.22s cubic-bezier(0.25, 0.8, 0.25, 1), transform 0.22s cubic-bezier(0.25, 0.8, 0.25, 1);
}

.admin-content-fade-enter-from {
  opacity: 0;
  transform: translateY(8px);
}

.admin-content-fade-leave-to {
  opacity: 0;
  transform: translateY(-4px);
}
</style>
