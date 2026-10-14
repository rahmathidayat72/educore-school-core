<template>
  <nav class="navbar navbar-expand-lg navbar-light bg-white border-bottom shadow-sm sticky-top public-navbar">
    <div class="container">
      <router-link to="/" class="navbar-brand d-flex align-items-center gap-2 brand-logo">
        <div class="brand-icon-wrap bg-success text-white rounded-circle d-flex align-items-center justify-content-center shadow-sm">
          <i class="fas fa-graduation-cap fs-5"></i>
        </div>
        <span class="fw-bold text-dark fs-5 brand-title">SEKOLAH KITA</span>
      </router-link>
      <button
        class="navbar-toggler border-0"
        type="button"
        data-bs-toggle="collapse"
        data-bs-target="#publicNav"
        aria-controls="publicNav"
        aria-expanded="false"
        aria-label="Toggle navigation"
      >
        <span class="navbar-toggler-icon"></span>
      </button>
      <div class="collapse navbar-collapse" id="publicNav">
        <ul class="navbar-nav me-auto mb-2 mb-lg-0 ms-lg-3 nav-menu">
          <li class="nav-item">
            <router-link to="/" class="nav-link text-dark nav-animated-link active">Beranda</router-link>
          </li>
          <li class="nav-item">
            <a class="nav-link text-dark nav-animated-link" href="#tentang">Tentang Kami</a>
          </li>
          <li class="nav-item">
            <a class="nav-link text-dark nav-animated-link" href="#keunggulan">Keunggulan</a>
          </li>
          <li class="nav-item">
            <a class="nav-link text-dark nav-animated-link" href="#kontak">Kontak</a>
          </li>
        </ul>
        <div class="d-flex align-items-center gap-2">
          <template v-if="authStore.isAuthenticated">
            <router-link to="/admin" class="btn btn-success btn-living d-flex align-items-center gap-2 shadow-sm">
              <i class="fas fa-columns"></i>
              <span>Dashboard Portal</span>
            </router-link>
          </template>
          <template v-else>
            <router-link to="/login" class="btn btn-outline-success btn-living d-flex align-items-center gap-2">
              <i class="fas fa-sign-in-alt"></i>
              <span>Login Portal</span>
            </router-link>
          </template>
        </div>
      </div>
    </div>
  </nav>
</template>

<script setup>
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
</script>

<style scoped>
.public-navbar {
  transition: all 0.3s ease;
  backdrop-filter: blur(8px);
  background-color: rgba(255, 255, 255, 0.96) !important;
  top: 0;
  z-index: 1030;
}

.brand-icon-wrap {
  width: 38px;
  height: 38px;
  transition: transform 0.35s cubic-bezier(0.16, 1, 0.3, 1);
}

.brand-logo:hover .brand-icon-wrap {
  transform: scale(1.12) rotate(-8deg);
}

.nav-animated-link {
  position: relative;
  font-weight: 500;
  padding-left: 0.75rem !important;
  padding-right: 0.75rem !important;
  transition: color 0.2s ease;
}

.nav-animated-link::after {
  content: '';
  position: absolute;
  bottom: 4px;
  left: 0.75rem;
  right: 0.75rem;
  height: 2px;
  background-color: #198754;
  transform: scaleX(0);
  transition: transform 0.25s cubic-bezier(0.16, 1, 0.3, 1);
  transform-origin: center;
}

.nav-animated-link:hover::after,
.nav-animated-link.active::after {
  transform: scaleX(1);
}

.nav-animated-link:hover {
  color: #198754 !important;
}
</style>
