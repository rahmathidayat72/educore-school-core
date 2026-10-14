<template>
  <header class="navbar navbar-expand bg-white border-bottom px-4 py-2 admin-navbar shadow-sm">
    <div class="d-flex align-items-center gap-3">
      <button class="btn btn-outline-secondary btn-sm d-lg-none" @click="$emit('toggle-sidebar')">
        <i class="fas fa-bars"></i>
      </button>
      <div class="fw-bold text-dark fs-6 d-none d-sm-block">
        Sistem Informasi Sekolah
      </div>
    </div>

    <div class="ms-auto d-flex align-items-center gap-3">
      <router-link to="/" class="btn btn-sm btn-outline-success btn-living d-flex align-items-center gap-2">
        <i class="fas fa-globe"></i>
        <span class="d-none d-md-inline">Lihat Web Utama</span>
      </router-link>

      <div class="dropdown">
        <button
          class="btn btn-light btn-sm dropdown-toggle d-flex align-items-center gap-2 border user-dropdown-btn shadow-sm"
          type="button"
          data-bs-toggle="dropdown"
          aria-expanded="false"
        >
          <div class="bg-success text-white rounded-circle d-flex align-items-center justify-content-center user-avatar" style="width: 24px; height: 24px; font-size: 0.75rem;">
            <i class="fas fa-user"></i>
          </div>
          <span class="fw-medium small">{{ authStore.currentUser?.username || 'User' }}</span>
        </button>
        <ul class="dropdown-menu dropdown-menu-end shadow-sm admin-dropdown-menu">
          <li class="px-3 py-2 border-bottom">
            <div class="small fw-bold text-dark">{{ authStore.currentUser?.username }}</div>
            <div class="text-secondary small">{{ authStore.currentUser?.email }}</div>
          </li>
          <li>
            <router-link to="/admin/users" class="dropdown-item small py-2">
              <i class="fas fa-user-cog me-2 text-secondary"></i>Manajemen User
            </router-link>
          </li>
          <li><hr class="dropdown-divider my-1"></li>
          <li>
            <button @click="logout" class="dropdown-item small py-2 text-danger">
              <i class="fas fa-sign-out-alt me-2"></i>Keluar
            </button>
          </li>
        </ul>
      </div>
    </div>
  </header>
</template>

<script setup>
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'

defineEmits(['toggle-sidebar'])

const authStore = useAuthStore()
const router = useRouter()

const logout = () => {
  authStore.logout()
  router.push('/login')
}
</script>

<style scoped>
.admin-navbar {
  flex-shrink: 0;
  z-index: 1020;
  width: 100%;
  min-height: 58px;
  background-color: #ffffff !important;
  border-bottom: 1px solid #e2e8f0;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
}

.admin-dropdown-menu {
  z-index: 1060 !important;
}

.user-dropdown-btn {
  transition: transform 0.2s cubic-bezier(0.16, 1, 0.3, 1), box-shadow 0.2s ease;
}

.user-dropdown-btn:hover {
  transform: translateY(-1px);
}

.user-avatar {
  transition: transform 0.25s cubic-bezier(0.16, 1, 0.3, 1);
}

.user-dropdown-btn:hover .user-avatar {
  transform: scale(1.15);
}
</style>
