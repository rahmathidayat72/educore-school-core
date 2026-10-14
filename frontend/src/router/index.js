import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

import LandingPage from '@/views/public/LandingPage.vue'
import Login from '@/views/auth/Login.vue'
import AdminLayout from '@/components/AdminLayout.vue'
import Dashboard from '@/views/admin/Dashboard.vue'
import Users from '@/views/admin/Users.vue'
import Guru from '@/views/admin/Guru.vue'
import Kelas from '@/views/admin/Kelas.vue'
import Siswa from '@/views/admin/Siswa.vue'
import MataPelajaran from '@/views/admin/MataPelajaran.vue'
import JadwalPelajaran from '@/views/admin/JadwalPelajaran.vue'
import Logs from '@/views/admin/Logs.vue'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: LandingPage,
  },
  {
    path: '/login',
    name: 'Login',
    component: Login,
    meta: { guestOnly: true },
  },
  {
    path: '/admin',
    component: AdminLayout,
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        name: 'AdminDashboard',
        component: Dashboard,
      },
      {
        path: 'users',
        name: 'AdminUsers',
        component: Users,
        meta: { roles: ['admin'] },
      },
      {
        path: 'guru',
        name: 'AdminGuru',
        component: Guru,
        meta: { roles: ['admin'] },
      },
      {
        path: 'kelas',
        name: 'AdminKelas',
        component: Kelas,
        meta: { roles: ['admin', 'guru'] },
      },
      {
        path: 'siswa',
        name: 'AdminSiswa',
        component: Siswa,
        meta: { roles: ['admin'] },
      },
      {
        path: 'mapel',
        name: 'AdminMapel',
        component: MataPelajaran,
        meta: { roles: ['admin', 'guru'] },
      },
      {
        path: 'jadwal',
        name: 'AdminJadwal',
        component: JadwalPelajaran,
        meta: { roles: ['admin', 'guru'] },
      },
      {
        path: 'logs',
        name: 'AdminLogs',
        component: Logs,
        meta: { roles: ['admin'] },
      },
    ],
  },
  {
    path: '/:pathMatch(.*)*',
    redirect: '/',
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior() {
    return { top: 0 }
  },
})

router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()

  if (to.matched.some((record) => record.meta.requiresAuth)) {
    if (!authStore.isAuthenticated) {
      return next({ name: 'Login', query: { redirect: to.fullPath } })
    }

    // Role-based protection
    const roleRecord = to.matched.find((record) => record.meta && (record.meta.roles || record.meta.requiresRole))
    if (roleRecord) {
      const allowedRoles = roleRecord.meta.roles || (roleRecord.meta.requiresRole ? [roleRecord.meta.requiresRole] : [])
      const currentRole = authStore.currentUser?.role ? String(authStore.currentUser.role).trim().toLowerCase() : (authStore.isAdmin ? 'admin' : 'user')
      if (allowedRoles.length > 0 && !allowedRoles.includes(currentRole)) {
        return next({ name: 'AdminDashboard' })
      }
    }
  }

  if (to.matched.some((record) => record.meta.guestOnly)) {
    if (authStore.isAuthenticated) {
      return next({ name: 'AdminDashboard' })
    }
  }

  next()
})

export default router
