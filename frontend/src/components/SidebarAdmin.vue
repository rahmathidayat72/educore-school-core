<template>
  <aside class="admin-sidebar shadow">
    <div class="p-3 border-bottom border-secondary d-flex align-items-center gap-2 sidebar-brand">
      <div class="brand-icon-box bg-success bg-opacity-25 p-2 rounded-circle d-flex align-items-center justify-content-center">
        <i class="fas fa-graduation-cap text-success fs-4 animate-icon-float"></i>
      </div>
      <div>
        <div class="fw-bold text-white fs-6">SEKOLAH KITA</div>
        <div class="text-secondary" style="font-size: 0.75rem;">{{ authStore.portalTitle }}</div>
      </div>
    </div>

    <div class="py-3 flex-grow-1 overflow-auto">
      <div class="px-3 mb-2 text-uppercase text-secondary fw-semibold" style="font-size: 0.7rem; letter-spacing: 0.05em;">
        Menu Utama
      </div>

      <nav class="nav flex-column">
        <router-link to="/admin" exact class="nav-link">
          <i class="fas fa-tachometer-alt fa-fw"></i>
          <span>Dashboard</span>
        </router-link>

        <router-link v-if="authStore.canManageSiswa" to="/admin/siswa" class="nav-link">
          <i class="fas fa-user-graduate fa-fw"></i>
          <span>Data Siswa</span>
        </router-link>

        <router-link v-if="authStore.canManageGuru" to="/admin/guru" class="nav-link">
          <i class="fas fa-chalkboard-teacher fa-fw"></i>
          <span>Data Guru</span>
        </router-link>

        <router-link to="/admin/kelas" class="nav-link">
          <i class="fas fa-school fa-fw"></i>
          <span>Data Kelas</span>
        </router-link>

        <router-link to="/admin/mapel" class="nav-link">
          <i class="fas fa-book-open fa-fw"></i>
          <span>Mata Pelajaran</span>
        </router-link>

        <router-link to="/admin/jadwal" class="nav-link">
          <i class="fas fa-calendar-alt fa-fw"></i>
          <span>Jadwal Pelajaran</span>
        </router-link>

        <div v-if="authStore.isAdmin" class="px-3 my-2 text-uppercase text-secondary fw-semibold" style="font-size: 0.7rem; letter-spacing: 0.05em;">
          Sistem
        </div>

        <router-link v-if="authStore.canManageUsers" to="/admin/users" class="nav-link">
          <i class="fas fa-users-cog fa-fw"></i>
          <span>Kelola Pengguna</span>
        </router-link>

        <router-link v-if="authStore.isAdmin" to="/admin/logs" class="nav-link">
          <i class="fas fa-clipboard-list fa-fw"></i>
          <span>Log Aktivitas</span>
        </router-link>

        <router-link to="/" class="nav-link">
          <i class="fas fa-external-link-alt fa-fw"></i>
          <span>Halaman Publik</span>
        </router-link>
      </nav>
    </div>

    <div class="p-3 border-top border-secondary">
      <div class="d-flex align-items-center justify-content-between">
        <div class="d-flex align-items-center gap-2 overflow-hidden me-2">
          <div class="bg-success text-white rounded-circle p-2 d-flex align-items-center justify-content-center" style="width: 34px; height: 34px;">
            <i class="fas fa-user small"></i>
          </div>
          <div class="overflow-hidden">
            <div class="text-white text-truncate small fw-bold">{{ authStore.currentUser?.username || 'Pengguna' }}</div>
            <div class="text-secondary text-truncate" style="font-size: 0.7rem;">
              {{ authStore.isAdmin ? 'Admin' : (authStore.isGuru ? 'Guru' : 'Pengguna') }} &bull; {{ authStore.currentUser?.email || '' }}
            </div>
          </div>
        </div>
        <button @click="handleLogout" class="btn btn-sm btn-outline-danger p-1 px-2" title="Keluar">
          <i class="fas fa-sign-out-alt"></i>
        </button>
      </div>
    </div>
  </aside>
</template>

<script setup>
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'

const authStore = useAuthStore()
const router = useRouter()

const handleLogout = () => {
  if (confirm('Apakah Anda yakin ingin keluar dari sistem?')) {
    authStore.logout()
    router.push('/login')
  }
}
</script>
