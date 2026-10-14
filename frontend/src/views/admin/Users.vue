<template>
  <div class="animate-fade-up">
    <!-- Header -->
    <div class="d-flex flex-column flex-md-row justify-content-between align-items-md-center gap-3 mb-4">
      <div>
        <div class="d-flex align-items-center gap-2 mb-1">
          <h3 class="fw-bold text-dark mb-0">Manajemen Pengguna</h3>
          <span class="badge bg-success-subtle text-success px-2 py-1 small">
            {{ filteredUsers.length }} Akun
          </span>
        </div>
        <p class="text-secondary small mb-0">Kelola akun, kata sandi, dan tingkat hak akses pengguna sistem informasi sekolah.</p>
      </div>
      <button @click="openCreateModal" class="btn btn-success btn-shine d-flex align-items-center gap-2 shadow-sm">
        <i class="fas fa-user-plus"></i>
        <span>Tambah Pengguna</span>
      </button>
    </div>

    <!-- Alert Notifikasi -->
    <Transition name="alert-pop">
      <div v-if="alertMessage" :class="['alert alert-dismissible fade show small py-2 d-flex align-items-center gap-2 mb-4 shadow-sm', alertClass]">
        <i :class="alertIcon"></i>
        <div>{{ alertMessage }}</div>
        <button type="button" class="btn-close" @click="alertMessage = ''"></button>
      </div>
    </Transition>

    <!-- Search & Filter Card -->
    <div class="card card-custom p-3 mb-4 border-0 shadow-sm">
      <div class="row g-2 align-items-center">
        <div class="col-md-6">
          <div class="input-group">
            <span class="input-group-text bg-white text-muted border-end-0"><i class="fas fa-search"></i></span>
            <input
              v-model="searchQuery"
              type="text"
              class="form-control border-start-0 ps-0"
              placeholder="Cari berdasarkan nama atau email..."
            />
            <button v-if="searchQuery" class="btn btn-link text-secondary p-0 px-2" @click="searchQuery = ''" title="Bersihkan">
              <i class="fas fa-times-circle"></i>
            </button>
          </div>
        </div>
        <div class="col-md-4">
          <select v-model="roleFilter" class="form-select">
            <option value="">Semua Peran (Role)</option>
            <option value="admin">Administrator (Admin)</option>
            <option value="guru">Tenaga Pengajar (Guru)</option>
            <option value="user">Pengguna Standar (User)</option>
          </select>
        </div>
        <div class="col-md-2 text-md-end">
          <button @click="fetchUsers" class="btn btn-outline-secondary w-100 d-flex align-items-center justify-content-center gap-2" :disabled="loading">
            <i class="fas fa-sync-alt" :class="{ 'fa-spin': loading }"></i>
            <span>Segarkan</span>
          </button>
        </div>
      </div>
    </div>

    <!-- Tabel Data Users -->
    <div class="card card-custom overflow-hidden border-0 shadow-sm">
      <div class="table-responsive">
        <table class="table table-custom table-hover mb-0">
          <thead>
            <tr>
              <th style="width: 60px;">No</th>
              <th>Username</th>
              <th>Email</th>
              <th>Peran / Hak Akses</th>
              <th class="text-end" style="width: 140px;">Aksi</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="loading">
              <td colspan="5" class="text-center py-5 text-secondary">
                <div class="spinner-border spinner-border-sm text-success me-2" role="status"></div>
                Memuat data pengguna...
              </td>
            </tr>
            <tr v-else-if="filteredUsers.length === 0">
              <td colspan="5" class="text-center py-5 text-muted">
                <i class="fas fa-users-slash fs-2 mb-2 d-block text-secondary"></i>
                Tidak ada pengguna yang sesuai.
              </td>
            </tr>
            <tr v-for="(user, index) in filteredUsers" :key="user.id">
              <td class="text-secondary fw-semibold">{{ index + 1 }}</td>
              <td>
                <div class="d-flex align-items-center gap-2">
                  <div
                    class="avatar-user rounded-circle d-flex align-items-center justify-content-center flex-shrink-0"
                    :class="{
                      'role-admin': user.role === 'admin',
                      'role-guru': user.role === 'guru',
                      'role-user': user.role === 'user',
                    }"
                  >
                    <i :class="user.role === 'admin' ? 'fas fa-user-shield' : (user.role === 'guru' ? 'fas fa-chalkboard-teacher' : 'fas fa-user')" class="small"></i>
                  </div>
                  <span class="fw-bold text-dark">{{ user.username }}</span>
                </div>
              </td>
              <td class="text-secondary small">{{ user.email }}</td>
              <td>
                <span
                  class="badge px-2 py-1"
                  :class="{
                    'bg-danger-subtle text-danger': user.role === 'admin',
                    'bg-primary-subtle text-primary': user.role === 'guru',
                    'bg-secondary-subtle text-secondary': user.role === 'user',
                  }"
                >
                  <i :class="user.role === 'admin' ? 'fas fa-shield-alt me-1' : (user.role === 'guru' ? 'fas fa-graduation-cap me-1' : 'fas fa-user me-1')"></i>
                  {{ user.role === 'admin' ? 'Admin' : (user.role === 'guru' ? 'Guru' : 'User') }}
                </span>
              </td>
              <td class="text-end">
                <div class="btn-group btn-group-sm">
                  <button @click="openEditModal(user)" class="btn btn-outline-primary" title="Ubah Akun Pengguna">
                    <i class="fas fa-edit"></i>
                  </button>
                  <button @click="openDeleteModal(user)" class="btn btn-outline-danger" title="Hapus Akun Pengguna">
                    <i class="fas fa-trash-alt"></i>
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Modal Form (Create / Edit) -->
    <Teleport to="body">
      <Transition name="modal-fade">
        <div v-if="showModal" class="app-modal-overlay" tabindex="-1" @click.self="closeModal">
          <div class="modal-dialog modal-dialog-centered">
            <div class="modal-content border-0 shadow-lg">
              <div class="modal-header bg-light">
                <h5 class="modal-title fw-bold text-dark d-flex align-items-center gap-2">
                  <i :class="isEditing ? 'fas fa-user-edit text-success' : 'fas fa-user-plus text-success'"></i>
                  <span>{{ isEditing ? 'Ubah Data Pengguna' : 'Tambah Pengguna Baru' }}</span>
                </h5>
                <button type="button" class="btn-close" @click="closeModal"></button>
              </div>
              <form @submit.prevent="handleSubmit">
                <div class="modal-body p-4">
                  <div class="mb-3">
                    <label class="form-label small fw-semibold text-secondary">Username</label>
                    <div class="input-group">
                      <span class="input-group-text bg-light text-muted"><i class="fas fa-user"></i></span>
                      <input v-model="form.username" type="text" class="form-control" required placeholder="Contoh: budi_santoso" />
                    </div>
                  </div>

                  <div class="mb-3">
                    <label class="form-label small fw-semibold text-secondary">Email</label>
                    <div class="input-group">
                      <span class="input-group-text bg-light text-muted"><i class="fas fa-envelope"></i></span>
                      <input v-model="form.email" type="email" class="form-control" required placeholder="budi@sekolah.id" />
                    </div>
                  </div>

                  <div class="mb-3">
                    <label class="form-label small fw-semibold text-secondary">
                      Password {{ isEditing ? '(Kosongkan jika tidak ingin diubah)' : '' }}
                    </label>
                    <div class="input-group">
                      <span class="input-group-text bg-light text-muted"><i class="fas fa-lock"></i></span>
                      <input
                        v-model="form.password"
                        type="password"
                        class="form-control"
                        :required="!isEditing"
                        placeholder="••••••••"
                      />
                    </div>
                  </div>

                  <div class="mb-3">
                    <label class="form-label small fw-semibold text-secondary">Tingkat Hak Akses (Role)</label>
                    <div class="input-group">
                      <span class="input-group-text bg-light text-muted"><i class="fas fa-user-tag"></i></span>
                      <select v-model="form.role" class="form-select" required>
                        <option value="admin">Administrator (Admin)</option>
                        <option value="guru">Tenaga Pengajar (Guru)</option>
                        <option value="user">Pengguna Standar (User)</option>
                      </select>
                    </div>
                  </div>
                </div>
                <div class="modal-footer bg-light">
                  <button type="button" class="btn btn-light" @click="closeModal">Batal</button>
                  <button type="submit" class="btn btn-success btn-shine d-flex align-items-center gap-2" :disabled="submitting">
                    <span v-if="submitting" class="spinner-border spinner-border-sm"></span>
                    <span>{{ isEditing ? 'Simpan Perubahan' : 'Buat Pengguna' }}</span>
                  </button>
                </div>
              </form>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>

    <!-- Modal Konfirmasi Hapus Pengguna -->
    <Teleport to="body">
      <Transition name="modal-fade">
        <div v-if="showDeleteModal" class="app-modal-overlay" tabindex="-1" @click.self="closeDeleteModal">
          <div class="modal-dialog modal-dialog-centered modal-sm">
            <div class="modal-content border-0 shadow-lg text-center p-4">
              <div class="bg-danger-subtle text-danger rounded-circle p-3 d-inline-flex align-items-center justify-content-center mx-auto mb-3" style="width: 58px; height: 58px;">
                <i class="fas fa-exclamation-triangle fs-4"></i>
              </div>
              <h5 class="fw-bold text-dark mb-1">Hapus Akun Pengguna?</h5>
              <p class="text-secondary small mb-3">
                Apakah Anda yakin ingin menghapus akun <strong class="text-dark">{{ userToDelete?.username }}</strong>?
              </p>
              <div class="d-flex gap-2 justify-content-center">
                <button type="button" class="btn btn-light btn-sm px-3" @click="closeDeleteModal" :disabled="deleting">
                  Batal
                </button>
                <button type="button" class="btn btn-danger btn-sm px-3 d-flex align-items-center gap-1 shadow-sm" @click="confirmDelete" :disabled="deleting">
                  <span v-if="deleting" class="spinner-border spinner-border-sm"></span>
                  <i v-else class="fas fa-trash-alt"></i>
                  <span>{{ deleting ? 'Menghapus...' : 'Ya, Hapus' }}</span>
                </button>
              </div>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { userApi } from '@/api/userApi'

const router = useRouter()
const authStore = useAuthStore()

const users = ref([])
const loading = ref(false)
const submitting = ref(false)
const deleting = ref(false)
const searchQuery = ref('')
const roleFilter = ref('')

const alertMessage = ref('')
const alertClass = ref('alert-success')
const alertIcon = ref('fas fa-check-circle')

const showModal = ref(false)
const isEditing = ref(false)
const editingId = ref(null)

const showDeleteModal = ref(false)
const userToDelete = ref(null)

const form = ref({
  username: '',
  email: '',
  password: '',
  role: 'user',
})

const showAlert = (message, type = 'success') => {
  alertMessage.value = message
  if (type === 'success') {
    alertClass.value = 'alert-success'
    alertIcon.value = 'fas fa-check-circle'
  } else {
    alertClass.value = 'alert-danger'
    alertIcon.value = 'fas fa-exclamation-triangle'
  }
}

const fetchUsers = async () => {
  loading.value = true
  try {
    const res = await userApi.getAll()
    if (res.data && res.data.data) {
      users.value = res.data.data
    } else {
      users.value = []
    }
  } catch (err) {
    showAlert('Gagal mengambil daftar pengguna.', 'danger')
  } finally {
    loading.value = false
  }
}

const filteredUsers = computed(() => {
  return users.value.filter((u) => {
    const matchesSearch =
      u.username?.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      u.email?.toLowerCase().includes(searchQuery.value.toLowerCase())
    const matchesRole = roleFilter.value ? u.role === roleFilter.value : true
    return matchesSearch && matchesRole
  })
})

const openCreateModal = () => {
  isEditing.value = false
  editingId.value = null
  form.value = {
    username: '',
    email: '',
    password: '',
    role: 'user',
  }
  showModal.value = true
}

const openEditModal = (user) => {
  isEditing.value = true
  editingId.value = user.id
  form.value = {
    username: user.username,
    email: user.email,
    password: '',
    role: user.role,
  }
  showModal.value = true
}

const closeModal = () => {
  showModal.value = false
}

const openDeleteModal = (user) => {
  userToDelete.value = user
  showDeleteModal.value = true
}

const closeDeleteModal = () => {
  showDeleteModal.value = false
  userToDelete.value = null
}

const handleSubmit = async () => {
  submitting.value = true
  try {
    if (isEditing.value) {
      await userApi.update(editingId.value, form.value)
      showAlert('Pengguna berhasil diperbarui!')
    } else {
      await userApi.create(form.value)
      showAlert('Pengguna baru berhasil ditambahkan!')
    }
    closeModal()
    fetchUsers()
  } catch (err) {
    const errMsg = err.response?.data?.message || err.message || 'Gagal memproses data pengguna.'
    showAlert(errMsg, 'danger')
  } finally {
    submitting.value = false
  }
}

const confirmDelete = async () => {
  if (!userToDelete.value) return
  deleting.value = true
  try {
    await userApi.delete(userToDelete.value.id)
    showAlert(`Pengguna "${userToDelete.value.username}" berhasil dihapus!`)
    closeDeleteModal()
    fetchUsers()
  } catch (err) {
    showAlert('Gagal menghapus pengguna.', 'danger')
  } finally {
    deleting.value = false
  }
}

onMounted(() => {
  if (!authStore.isAdmin) {
    router.replace('/admin')
    return
  }
  fetchUsers()
})
</script>

<style scoped>
.avatar-user {
  width: 34px;
  height: 34px;
  transition: transform 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
}
.role-admin {
  background: linear-gradient(135deg, rgba(220, 53, 69, 0.15), rgba(239, 68, 68, 0.25));
  color: #dc3545;
}
.role-guru {
  background: linear-gradient(135deg, rgba(13, 110, 253, 0.15), rgba(2, 132, 199, 0.25));
  color: #0d6efd;
}
.role-user {
  background: linear-gradient(135deg, rgba(108, 117, 125, 0.15), rgba(100, 116, 139, 0.25));
  color: #64748b;
}
tr:hover .avatar-user {
  transform: scale(1.1) rotate(-6deg);
}
</style>
