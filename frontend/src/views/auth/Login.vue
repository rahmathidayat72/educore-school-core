<template>
  <div class="login-wrapper min-vh-100 d-flex align-items-center justify-content-center py-5 position-relative overflow-hidden">
    <!-- Living Decorative Orbs di Belakang -->
    <div class="bg-orb bg-orb-1"></div>
    <div class="bg-orb bg-orb-2"></div>
    <div class="bg-orb bg-orb-3"></div>

    <div class="container position-relative" style="z-index: 2;">
      <div class="row justify-content-center">
        <div class="col-md-6 col-lg-5">
          <div class="text-center mb-4 login-header">
            <router-link to="/" class="d-inline-flex align-items-center gap-2 text-decoration-none mb-2 brand-link">
              <div class="brand-icon-box bg-white shadow-sm rounded-circle d-flex align-items-center justify-content-center animate-float">
                <i class="fas fa-graduation-cap text-success fs-2"></i>
              </div>
              <span class="fs-4 fw-bold text-dark">SEKOLAH KITA</span>
            </router-link>
            <h5 class="fw-semibold text-secondary">Portal Sistem Informasi Akademik</h5>
          </div>

          <div class="card card-custom card-living-shine p-4 p-md-5 login-card shadow-lg">
            <h4 class="fw-bold text-dark mb-1">Masuk ke Akun Anda</h4>
            <p class="text-muted small mb-4">Masukkan email dan kata sandi untuk mengakses dashboard.</p>

            <div v-if="errorMessage" class="alert alert-danger d-flex align-items-center gap-2 small py-2 error-shake" role="alert">
              <i class="fas fa-exclamation-circle fs-5 flex-shrink-0"></i>
              <div>{{ errorMessage }}</div>
            </div>

            <div v-if="isSessionExpired" class="alert alert-warning d-flex align-items-center gap-2 small py-2 animate-breathe" role="alert">
              <i class="fas fa-clock fs-5 flex-shrink-0"></i>
              <div>Sesi Anda telah berakhir. Silakan login kembali.</div>
            </div>

            <form @submit.prevent="handleLogin">
              <div class="mb-3 input-field-group">
                <label for="email" class="form-label fw-medium small text-secondary">Alamat Email</label>
                <div class="input-group">
                  <span class="input-group-text bg-white text-muted"><i class="fas fa-envelope"></i></span>
                  <input
                    id="email"
                    v-model="email"
                    type="email"
                    class="form-control"
                    placeholder="nama@sekolah.id"
                    autocomplete="email"
                    required
                  />
                </div>
              </div>

              <div class="mb-4 input-field-group">
                <label for="password" class="form-label fw-medium small text-secondary">Kata Sandi</label>
                <div class="input-group">
                  <span class="input-group-text bg-white text-muted"><i class="fas fa-lock"></i></span>
                  <input
                    id="password"
                    v-model="password"
                    type="password"
                    class="form-control"
                    placeholder="••••••••"
                    autocomplete="current-password"
                    required
                  />
                </div>
              </div>

              <button
                type="submit"
                class="btn btn-success btn-living w-100 py-2 fw-semibold d-flex align-items-center justify-content-center gap-2 shadow-sm"
                :disabled="loading"
              >
                <span v-if="loading" class="spinner-border spinner-border-sm" role="status" aria-hidden="true"></span>
                <span v-else><i class="fas fa-sign-in-alt me-1"></i> Masuk Sekarang</span>
              </button>
            </form>

            <div class="mt-4 pt-3 border-top text-center">
              <div class="text-muted small mb-2">Akun Demo Default:</div>
              <div class="bg-light p-2 rounded small text-secondary font-monospace border">
                <div>Email: <strong>admin@sekolah.id</strong></div>
                <div>Pass: <strong>adminpassword</strong></div>
              </div>
              <div class="mt-3">
                <router-link to="/" class="text-decoration-none text-success small fw-medium back-link">
                  <i class="fas fa-arrow-left me-1 back-arrow"></i> Kembali ke Halaman Utama
                </router-link>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const email = ref('')
const password = ref('')
const loading = ref(false)
const errorMessage = ref('')

const isSessionExpired = computed(() => route.query.expired === '1')

const handleLogin = async () => {
  loading.value = true
  errorMessage.value = ''

  const result = await authStore.login(email.value, password.value)
  loading.value = false

  if (result.success) {
    const redirectPath = route.query.redirect || '/admin'
    router.push(redirectPath)
  } else {
    errorMessage.value = result.message || 'Login gagal. Silakan periksa kembali email dan password Anda.'
  }
}
</script>

<style scoped>
.login-wrapper {
  background: radial-gradient(circle at top left, #f1f5f9 0%, #e2e8f0 100%);
}

/* Background Living Orbs */
.bg-orb {
  position: absolute;
  border-radius: 50%;
  filter: blur(65px);
  opacity: 0.6;
  pointer-events: none;
}

.bg-orb-1 {
  width: 380px;
  height: 380px;
  background: rgba(15, 118, 110, 0.22);
  top: -80px;
  left: -100px;
  animation: floatOrb 12s ease-in-out infinite alternate;
}

.bg-orb-2 {
  width: 420px;
  height: 420px;
  background: rgba(2, 132, 199, 0.2);
  bottom: -100px;
  right: -100px;
  animation: floatOrb 14s ease-in-out infinite alternate-reverse;
}

.bg-orb-3 {
  width: 250px;
  height: 250px;
  background: rgba(16, 185, 129, 0.15);
  top: 40%;
  right: 15%;
  animation: floatOrb 9s ease-in-out infinite alternate;
}

@keyframes floatOrb {
  0% { transform: translate(0, 0) scale(1); }
  50% { transform: translate(30px, -40px) scale(1.08); }
  100% { transform: translate(-20px, 30px) scale(0.95); }
}

.brand-icon-box {
  width: 52px;
  height: 52px;
  transition: transform 0.3s cubic-bezier(0.16, 1, 0.3, 1);
}

.brand-link:hover .brand-icon-box {
  transform: scale(1.1) rotate(6deg);
}

.login-card {
  animation: loginCardPop 0.6s cubic-bezier(0.16, 1, 0.3, 1);
  border-radius: 1.15rem;
  border: 1px solid rgba(226, 232, 240, 0.8);
}

@keyframes loginCardPop {
  0% {
    opacity: 0;
    transform: translateY(28px) scale(0.96);
  }
  100% {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

.input-field-group .form-control {
  border-left: none;
  transition: border-color 0.2s ease, box-shadow 0.2s ease;
}

.input-field-group .input-group-text {
  border-right: none;
  transition: border-color 0.2s ease;
}

.input-field-group:focus-within .input-group-text,
.input-field-group:focus-within .form-control {
  border-color: #198754;
}

.input-field-group:focus-within .form-control {
  box-shadow: 0 0 0 0.25rem rgba(25, 135, 84, 0.15);
}

/* Error Shake */
.error-shake {
  animation: errorShake 0.45s cubic-bezier(0.36, 0.07, 0.19, 0.97) both;
}

@keyframes errorShake {
  10%, 90% { transform: translate3d(-1px, 0, 0); }
  20%, 80% { transform: translate3d(2px, 0, 0); }
  30%, 50%, 70% { transform: translate3d(-4px, 0, 0); }
  40%, 60% { transform: translate3d(4px, 0, 0); }
}

.back-arrow {
  transition: transform 0.25s cubic-bezier(0.16, 1, 0.3, 1);
  display: inline-block;
}

.back-link:hover .back-arrow {
  transform: translateX(-4px);
}
</style>
