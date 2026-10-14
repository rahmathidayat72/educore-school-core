<template>
  <Teleport to="body">
    <Transition name="modal-fade">
      <div v-if="isOpen" class="app-modal-overlay calendar-modal-backdrop" tabindex="-1" @click.self="close" @keydown.esc="close">
        <div class="modal-dialog modal-dialog-centered modal-lg calendar-dialog">
          <div class="modal-content border-0 shadow-lg rounded-4 overflow-hidden">
            <!-- Modal Header -->
            <div class="modal-header bg-primary text-white py-3 px-4 d-flex justify-content-between align-items-center">
              <div class="d-flex align-items-center gap-2">
                <div class="bg-white text-primary rounded-circle p-2 d-flex align-items-center justify-content-center" style="width: 38px; height: 38px;">
                  <i class="fas fa-calendar-alt fs-5"></i>
                </div>
                <div>
                  <h5 class="modal-title fw-bold mb-0">Kalender Indonesia</h5>
                  <div class="small text-white-50">Kalender Masehi & Hari Libur Nasional</div>
                </div>
              </div>
              <button type="button" class="btn-close btn-close-white" @click="close" aria-label="Tutup"></button>
            </div>

            <!-- Modal Body -->
            <div class="modal-body p-3 p-md-4">
          <!-- Navigation Bar -->
          <div class="d-flex flex-wrap align-items-center justify-content-between gap-2 mb-3 pb-3 border-bottom">
            <div class="d-flex align-items-center gap-2">
              <button class="btn btn-outline-secondary btn-sm" @click="prevMonth" title="Bulan Sebelumnya">
                <i class="fas fa-chevron-left"></i>
              </button>
              <h4 class="fw-bold text-dark mb-0 text-capitalize font-monospace" style="min-width: 220px;">
                {{ monthNames[viewMonth] }} {{ viewYear }}
              </h4>
              <button class="btn btn-outline-secondary btn-sm" @click="nextMonth" title="Bulan Berikutnya">
                <i class="fas fa-chevron-right"></i>
              </button>
            </div>

            <div class="d-flex align-items-center gap-2">
              <button class="btn btn-primary btn-sm d-flex align-items-center gap-1 px-3" @click="goToToday" title="Kembali ke Hari Ini">
                <i class="fas fa-calendar-check"></i>
                <span>Hari Ini</span>
              </button>
              <span class="badge bg-danger text-white py-2 px-3 small">
                Waktu Indonesia Barat (WIB)
              </span>
            </div>
          </div>

          <!-- Calendar Grid -->
          <div class="calendar-wrapper">
            <!-- Day of Week Headers (Senin - Minggu) -->
            <div class="calendar-grid text-center fw-bold small text-secondary mb-2">
              <div class="p-2 text-dark">Sen</div>
              <div class="p-2 text-dark">Sel</div>
              <div class="p-2 text-dark">Rab</div>
              <div class="p-2 text-dark">Kam</div>
              <div class="p-2 text-dark">Jum</div>
              <div class="p-2 text-dark">Sab</div>
              <div class="p-2 text-danger">Min</div>
            </div>

            <!-- Date Cells -->
            <div class="calendar-grid">
              <div
                v-for="(cell, idx) in calendarCells"
                :key="idx"
                class="calendar-cell p-2 border rounded-3 position-relative d-flex flex-column align-items-center justify-content-between"
                :class="{
                  'other-month bg-light text-muted opacity-50': !cell.isCurrentMonth,
                  'is-sunday text-danger bg-danger-subtle bg-opacity-25': cell.isCurrentMonth && cell.isSunday,
                  'is-holiday text-danger bg-danger-subtle': cell.isCurrentMonth && cell.holiday,
                  'is-today shadow-sm': cell.isToday,
                  'selected-cell': isSelected(cell),
                }"
                @click="selectCell(cell)"
              >
                <!-- Nomor Tanggal -->
                <div class="d-flex align-items-center justify-content-between w-100 mb-1">
                  <span
                    class="date-number fw-bold"
                    :class="{
                      'today-badge': cell.isToday,
                      'text-danger': cell.isCurrentMonth && (cell.isSunday || cell.holiday),
                    }"
                  >
                    {{ cell.day }}
                  </span>
                  <span v-if="cell.holiday && cell.isCurrentMonth" class="badge-dot" :title="cell.holiday"></span>
                </div>

                <!-- Label Libur Miniatur -->
                <div v-if="cell.holiday && cell.isCurrentMonth" class="holiday-mini text-truncate small" :title="cell.holiday">
                  {{ cell.holiday }}
                </div>
                <div v-else class="holiday-empty"></div>
              </div>
            </div>
          </div>

          <!-- Detail Tanggal yang Dipilih atau Hari Ini -->
          <div v-if="selectedCell" class="mt-3 p-3 rounded-3 border bg-light d-flex align-items-center justify-content-between flex-wrap gap-2">
            <div class="d-flex align-items-center gap-2">
              <div class="bg-primary text-white rounded-circle p-2 d-flex align-items-center justify-content-center" style="width: 36px; height: 36px;">
                <i class="far fa-clock"></i>
              </div>
              <div>
                <div class="fw-bold text-dark">
                  {{ formatSelectedDate(selectedCell) }}
                </div>
                <div v-if="selectedCell.holiday" class="text-danger small fw-semibold">
                  <i class="fas fa-flag me-1"></i>{{ selectedCell.holiday }}
                </div>
                <div v-else-if="selectedCell.isSunday" class="text-danger small">
                  Hari Minggu (Libur Mingguan)
                </div>
                <div v-else class="text-secondary small">
                  Hari Kerja / Jam Belajar Mengajar Aktif
                </div>
              </div>
            </div>
            <span v-if="selectedCell.isToday" class="badge bg-success-subtle text-success px-3 py-2">
              <i class="fas fa-check-circle me-1"></i>Tanggal Hari Ini
            </span>
          </div>

          <!-- Daftar Hari Libur Nasional di Bulan yang Sedang Dilihat -->
          <div class="mt-3 pt-3 border-top">
            <div class="d-flex align-items-center justify-content-between mb-2">
              <div class="fw-bold text-dark small text-uppercase" style="letter-spacing: 0.05em;">
                <i class="fas fa-umbrella-beach text-danger me-1"></i>
                Hari Libur & Peringatan Nasional ({{ monthNames[viewMonth] }} {{ viewYear }})
              </div>
              <span class="badge bg-secondary-subtle text-dark small">
                {{ monthHolidays.length }} Hari Libur
              </span>
            </div>

            <div v-if="monthHolidays.length > 0" class="row g-2">
              <div v-for="h in monthHolidays" :key="h.date" class="col-md-6">
                <div class="p-2 border border-danger-subtle rounded-3 bg-danger-subtle bg-opacity-25 d-flex align-items-center gap-2">
                  <span class="badge bg-danger text-white font-monospace">{{ h.day }}</span>
                  <div class="small fw-semibold text-danger text-truncate" :title="h.description">
                    {{ h.description }}
                  </div>
                </div>
              </div>
            </div>
            <div v-else class="text-secondary small fst-italic py-1">
              Tidak ada hari libur nasional resmi pada bulan {{ monthNames[viewMonth] }} {{ viewYear }}.
            </div>
          </div>
        </div>

        <!-- Modal Footer -->
        <div class="modal-footer bg-light py-2 px-4 d-flex justify-content-between align-items-center">
          <div class="small text-secondary">
            <span class="badge bg-primary me-1">&nbsp;</span> Tanggal Hari Ini &nbsp;|&nbsp;
            <span class="badge bg-danger me-1">&nbsp;</span> Hari Minggu / Libur Nasional
          </div>
          <button type="button" class="btn btn-secondary btn-sm px-4" @click="close">
            Tutup
          </button>
        </div>
      </div>
    </div>
  </div>
  </Transition>
  </Teleport>
</template>

<script setup>
import { ref, computed, watch, onMounted, onBeforeUnmount } from 'vue'
import { getHoliday, getHolidaysInMonth } from '@/data/holidaysIndonesia'

const props = defineProps({
  isOpen: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['close'])

const monthNames = [
  'Januari', 'Februari', 'Maret', 'April', 'Mei', 'Juni',
  'Juli', 'Agustus', 'September', 'Oktober', 'November', 'Desember',
]

const dayNames = [
  'Minggu', 'Senin', 'Selasa', 'Rabu', 'Kamis', 'Jumat', 'Sabtu',
]

// Waktu Sekarang berbasis WIB (Asia/Jakarta)
const getTodayWIB = () => {
  const now = new Date()
  const formatter = new Intl.DateTimeFormat('en-CA', {
    timeZone: 'Asia/Jakarta',
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
  })
  const parts = formatter.format(now).split('-')
  return {
    year: parseInt(parts[0], 10),
    month: parseInt(parts[1], 10) - 1,
    day: parseInt(parts[2], 10),
  }
}

const today = getTodayWIB()
const viewYear = ref(today.year)
const viewMonth = ref(today.month)
const selectedCell = ref(null)

// Inisialisasi sel hari ini saat dibuka
watch(
  () => props.isOpen,
  (val) => {
    if (val) {
      const current = getTodayWIB()
      viewYear.value = current.year
      viewMonth.value = current.month
      selectedCell.value = {
        year: current.year,
        month: current.month,
        day: current.day,
        isCurrentMonth: true,
        isToday: true,
        isSunday: new Date(current.year, current.month, current.day).getDay() === 0,
        holiday: getHoliday(current.year, current.month, current.day),
      }
    }
  }
)

const prevMonth = () => {
  if (viewMonth.value === 0) {
    viewMonth.value = 11
    viewYear.value -= 1
  } else {
    viewMonth.value -= 1
  }
}

const nextMonth = () => {
  if (viewMonth.value === 11) {
    viewMonth.value = 0
    viewYear.value += 1
  } else {
    viewMonth.value += 1
  }
}

const goToToday = () => {
  const current = getTodayWIB()
  viewYear.value = current.year
  viewMonth.value = current.month
  selectedCell.value = {
    year: current.year,
    month: current.month,
    day: current.day,
    isCurrentMonth: true,
    isToday: true,
    isSunday: new Date(current.year, current.month, current.day).getDay() === 0,
    holiday: getHoliday(current.year, current.month, current.day),
  }
}

const close = () => {
  emit('close')
}

// Daftar hari libur pada bulan yang sedang dilihat
const monthHolidays = computed(() => {
  return getHolidaysInMonth(viewYear.value, viewMonth.value)
})

// Susunan sel kalender (Senin - Minggu)
const calendarCells = computed(() => {
  const year = viewYear.value
  const month = viewMonth.value
  const currentToday = getTodayWIB()

  // Hari pertama di bulan ini (0 = Minggu, 1 = Senin, ...)
  const firstDayOfWeek = new Date(year, month, 1).getDay()
  // Normalisasi Senin sebagai kolom pertama (0 = Sen, 6 = Min)
  const startOffset = (firstDayOfWeek + 6) % 7

  const daysInMonth = new Date(year, month + 1, 0).getDate()
  const daysInPrevMonth = new Date(year, month, 0).getDate()

  const cells = []

  // Hari dari bulan sebelumnya
  for (let i = startOffset - 1; i >= 0; i--) {
    const day = daysInPrevMonth - i
    const prevMonthIdx = month === 0 ? 11 : month - 1
    const prevYear = month === 0 ? year - 1 : year
    const dayOfWeek = new Date(prevYear, prevMonthIdx, day).getDay()
    cells.push({
      year: prevYear,
      month: prevMonthIdx,
      day,
      isCurrentMonth: false,
      isToday: false,
      isSunday: dayOfWeek === 0,
      holiday: null,
    })
  }

  // Hari di bulan berjalan
  for (let d = 1; d <= daysInMonth; d++) {
    const isToday =
      year === currentToday.year &&
      month === currentToday.month &&
      d === currentToday.day
    const dayOfWeek = new Date(year, month, d).getDay()
    const holidayDesc = getHoliday(year, month, d)

    cells.push({
      year,
      month,
      day: d,
      isCurrentMonth: true,
      isToday,
      isSunday: dayOfWeek === 0,
      holiday: holidayDesc,
    })
  }

  // Lengkapi sisa grid hingga kelipatan 7 (maks 42 sel)
  const remaining = (7 - (cells.length % 7)) % 7
  for (let i = 1; i <= remaining; i++) {
    const nextMonthIdx = month === 11 ? 0 : month + 1
    const nextYear = month === 11 ? year + 1 : year
    const dayOfWeek = new Date(nextYear, nextMonthIdx, i).getDay()
    cells.push({
      year: nextYear,
      month: nextMonthIdx,
      day: i,
      isCurrentMonth: false,
      isToday: false,
      isSunday: dayOfWeek === 0,
      holiday: null,
    })
  }

  return cells
})

const selectCell = (cell) => {
  if (!cell.isCurrentMonth) {
    viewYear.value = cell.year
    viewMonth.value = cell.month
  }
  selectedCell.value = cell
}

const isSelected = (cell) => {
  if (!selectedCell.value) return false
  return (
    selectedCell.value.year === cell.year &&
    selectedCell.value.month === cell.month &&
    selectedCell.value.day === cell.day
  )
}

const formatSelectedDate = (cell) => {
  if (!cell) return ''
  const dt = new Date(cell.year, cell.month, cell.day)
  const dayName = dayNames[dt.getDay()]
  return `${dayName}, ${cell.day} ${monthNames[cell.month]} ${cell.year}`
}

const handleKeydown = (e) => {
  if (e.key === 'Escape' && props.isOpen) {
    close()
  }
}

onMounted(() => {
  window.addEventListener('keydown', handleKeydown)
})

onBeforeUnmount(() => {
  window.removeEventListener('keydown', handleKeydown)
})
</script>

<style scoped>
.calendar-modal-backdrop {
  background-color: rgba(15, 23, 42, 0.65);
  backdrop-filter: blur(6px);
  z-index: 1060;
  transition: backdrop-filter 0.3s ease;
}

.calendar-dialog {
  animation: calendarPop 0.35s cubic-bezier(0.16, 1, 0.3, 1);
}

@keyframes calendarPop {
  0% {
    opacity: 0;
    transform: scale(0.92) translateY(-20px);
  }
  100% {
    opacity: 1;
    transform: scale(1) translateY(0);
  }
}

.calendar-grid {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: 6px;
}

.calendar-cell {
  min-height: 64px;
  cursor: pointer;
  background-color: #ffffff;
  border-color: #e2e8f0 !important;
  transition: transform 0.22s cubic-bezier(0.16, 1, 0.3, 1), box-shadow 0.22s ease, border-color 0.2s ease, background-color 0.2s ease;
}

.calendar-cell:hover {
  border-color: #0d6efd !important;
  transform: translateY(-3px) scale(1.02);
  box-shadow: 0 6px 14px rgba(13, 110, 253, 0.15) !important;
  z-index: 1;
}

.calendar-cell.selected-cell {
  border-color: #0d6efd !important;
  background-color: #eff6ff !important;
  box-shadow: 0 0 0 2.5px rgba(13, 110, 253, 0.3);
}

.date-number {
  font-size: 0.95rem;
}

.today-badge {
  background-color: #0d6efd;
  color: #ffffff !important;
  width: 26px;
  height: 26px;
  border-radius: 50%;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 2px 6px rgba(13, 110, 253, 0.4);
  animation: breathSubtle 3s ease-in-out infinite;
}

.badge-dot {
  width: 7px;
  height: 7px;
  background-color: #dc3545;
  border-radius: 50%;
}

.holiday-mini {
  font-size: 0.65rem;
  color: #dc3545;
  font-weight: 600;
  width: 100%;
  text-align: center;
  line-height: 1.1;
}

.holiday-empty {
  height: 10px;
}
</style>
