<template>
  <Teleport to="body">
    <Transition name="modal-fade">
      <div v-if="isOpen" class="app-modal-overlay" tabindex="-1" @click.self="closeModal">
        <div class="modal-dialog modal-dialog-centered modal-lg modal-dialog-scrollable">
          <div class="modal-content border-0 shadow-lg bg-white d-flex flex-column" style="max-height: 90vh;">
            
            <!-- Modal Header -->
            <div class="modal-header bg-dark text-white flex-shrink-0">
              <h5 class="modal-title fw-bold d-flex align-items-center gap-2">
                <i class="fas fa-file-import text-warning"></i>
                <span>Import Data Masal (Excel / CSV) — {{ title }}</span>
              </h5>
              <button type="button" class="btn-close btn-close-white" @click="closeModal" :disabled="importing"></button>
            </div>

            <!-- Modal Body -->
            <div class="modal-body p-4 overflow-y-auto">
              <!-- Langkah 1: Download Template Resmi -->
              <div class="alert alert-info d-flex flex-column flex-sm-row align-items-sm-center justify-content-between gap-3 shadow-sm mb-4 border-0">
                <div class="d-flex align-items-center gap-3">
                  <div class="fs-2 text-info"><i class="fas fa-info-circle"></i></div>
                  <div>
                    <div class="fw-bold text-dark">Gunakan Format Template Baku</div>
                    <div class="small text-secondary">
                      Unduh template resmi agar struktur kolom dan format data (seperti tanggal YYYY-MM-DD) seragam serta valid.
                    </div>
                  </div>
                </div>
                <button @click="onDownloadTemplate" type="button" class="btn btn-primary btn-sm flex-shrink-0 shadow-sm d-flex align-items-center gap-2">
                  <i class="fas fa-download"></i>
                  <span>Unduh Template CSV</span>
                </button>
              </div>

              <!-- Langkah 2: Area Upload File -->
              <div class="mb-4">
                <label class="form-label fw-bold small text-dark mb-2">Unggah Berkas CSV yang Telah Diisi</label>
                <div 
                  class="upload-dropzone p-4 text-center rounded-3 border-2 border-dashed bg-light"
                  :class="{ 'border-primary bg-primary-subtle': isDragging }"
                  @dragover.prevent="isDragging = true"
                  @dragleave.prevent="isDragging = false"
                  @drop.prevent="handleDrop"
                >
                  <i class="fas fa-file-csv fs-1 text-success mb-2"></i>
                  <div class="fw-semibold text-dark mb-1">Pilih atau seret berkas CSV ke sini</div>
                  <div class="small text-muted mb-3">Mendukung berkas format .csv (kompatibel Excel & Spreadsheet)</div>
                  
                  <input 
                    type="file" 
                    ref="fileInput" 
                    accept=".csv,text/csv" 
                    class="d-none" 
                    @change="handleFileSelect"
                  />
                  <button type="button" class="btn btn-outline-success btn-sm px-3 shadow-sm" @click="$refs.fileInput.click()">
                    <i class="fas fa-folder-open me-1"></i> Telusuri Berkas
                  </button>
                </div>
              </div>

              <!-- Info Berkas Terpilih -->
              <div v-if="selectedFileName" class="d-flex align-items-center justify-content-between p-2 px-3 bg-light rounded-3 border mb-3 small">
                <div class="d-flex align-items-center gap-2 text-dark">
                  <i class="fas fa-check-circle text-success"></i>
                  <span class="fw-bold">{{ selectedFileName }}</span>
                  <span class="badge bg-secondary-subtle text-secondary">{{ parsedRows.length }} Baris Data Ditemukan</span>
                </div>
                <button type="button" class="btn btn-link btn-sm text-danger text-decoration-none p-0" @click="resetFile">
                  <i class="fas fa-trash-alt me-1"></i>Ganti
                </button>
              </div>

              <!-- Langkah 3: Pratinjau Data (Preview 5 Baris Pertama) -->
              <div v-if="parsedRows.length > 0" class="mb-3">
                <h6 class="fw-bold text-dark small text-uppercase mb-2 d-flex align-items-center justify-content-between">
                  <span>Pratinjau Data (Maks. 5 Baris Pertama)</span>
                  <span class="badge bg-success-subtle text-success">Siap Diimpor</span>
                </h6>
                <div class="table-responsive border rounded-3 bg-white">
                  <table class="table table-sm table-striped table-hover mb-0" style="font-size: 0.75rem;">
                    <thead class="table-light">
                      <tr>
                        <th style="width: 35px;">No</th>
                        <th v-for="col in previewColumns" :key="col">{{ col }}</th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr v-for="(row, idx) in parsedRows.slice(0, 5)" :key="idx">
                        <td class="text-secondary fw-semibold">{{ idx + 1 }}</td>
                        <td v-for="col in previewColumns" :key="col" class="text-truncate" style="max-width: 150px;">
                          {{ row[col] || '-' }}
                        </td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>

              <!-- Progress Bar Saat Proses Import Berjalan -->
              <div v-if="importing" class="mt-3">
                <div class="d-flex justify-content-between small fw-semibold text-secondary mb-1">
                  <span>Sedang memproses penyimpanan data...</span>
                  <span>{{ progressPercent }}% ({{ processedCount }}/{{ parsedRows.length }})</span>
                </div>
                <div class="progress" style="height: 8px;">
                  <div class="progress-bar progress-bar-striped progress-bar-animated bg-success" :style="{ width: progressPercent + '%' }"></div>
                </div>
              </div>

              <!-- Log Error / Peringatan Baris Gagal -->
              <div v-if="importErrors.length > 0" class="alert alert-warning mt-3 small p-3">
                <div class="fw-bold text-dark mb-1"><i class="fas fa-exclamation-triangle me-1"></i>Peringatan / Data Tidak Lengkap ({{ importErrors.length }} baris):</div>
                <ul class="mb-0 ps-3">
                  <li v-for="(err, i) in importErrors.slice(0, 5)" :key="i">{{ err }}</li>
                  <li v-if="importErrors.length > 5">...dan {{ importErrors.length - 5 }} baris lainnya.</li>
                </ul>
              </div>

            </div>

            <!-- Modal Footer -->
            <div class="modal-footer bg-light flex-shrink-0">
              <button type="button" class="btn btn-secondary btn-sm" @click="closeModal" :disabled="importing">
                Tutup
              </button>
              <button 
                type="button" 
                class="btn btn-success btn-sm px-3 shadow-sm d-flex align-items-center gap-2" 
                @click="startImport" 
                :disabled="parsedRows.length === 0 || importing"
              >
                <span v-if="importing" class="spinner-border spinner-border-sm"></span>
                <i v-else class="fas fa-check"></i>
                <span>{{ importing ? 'Mengimpor Data...' : `Mulai Import (${parsedRows.length} Data)` }}</span>
              </button>
            </div>

          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup>
import { ref, computed } from 'vue'
import { downloadCSVTemplate, parseCSV } from '@/utils/exportHelper'

const props = defineProps({
  isOpen: {
    type: Boolean,
    default: false
  },
  title: {
    type: String,
    default: 'Data'
  },
  templateFilename: {
    type: String,
    default: 'Template_Import'
  },
  schema: {
    type: Array,
    required: true
    // Array of { key: 'nama', label: 'Nama Lengkap', example: 'Budi Santoso' }
  },
  sampleRows: {
    type: Array,
    default: () => []
  },
  // Handler async callback untuk menyimpan data per baris atau batch: (row, index) => Promise<any>
  onImportRow: {
    type: Function,
    required: true
  }
})

const emit = defineEmits(['close', 'success'])

const isDragging = ref(false)
const selectedFileName = ref('')
const parsedRows = ref([])
const importing = ref(false)
const processedCount = ref(0)
const importErrors = ref([])

const previewColumns = computed(() => {
  return props.schema.map(s => s.label || s.key)
})

const progressPercent = computed(() => {
  if (parsedRows.value.length === 0) return 0
  return Math.round((processedCount.value / parsedRows.value.length) * 100)
})

const onDownloadTemplate = () => {
  downloadCSVTemplate(props.templateFilename, props.schema, props.sampleRows)
}

const handleFileSelect = (e) => {
  const file = e.target.files[0]
  if (file) processFile(file)
}

const handleDrop = (e) => {
  isDragging.value = false
  const file = e.dataTransfer.files[0]
  if (file) processFile(file)
}

const processFile = (file) => {
  if (!file.name.endsWith('.csv')) {
    alert('Harap unggah berkas dengan format .csv')
    return
  }

  selectedFileName.value = file.name
  importErrors.value = []
  
  const reader = new FileReader()
  reader.onload = (event) => {
    try {
      const text = event.target.result
      const rawRows = parseCSV(text)
      if (rawRows.length <= 1) {
        alert('Berkas CSV kosong atau hanya berisi baris header.')
        return
      }

      // Baris pertama adalah header
      const headers = rawRows[0]
      // Petakan header CSV ke schema keys
      const dataRows = rawRows.slice(1)
      const mapped = []

      dataRows.forEach((rowValues) => {
        const rowObj = {}
        props.schema.forEach((colSchema) => {
          // Cari indeks kolom yang cocok dengan label atau key
          const targetIndex = headers.findIndex(h => 
            h.toLowerCase() === (colSchema.label || '').toLowerCase() ||
            h.toLowerCase() === colSchema.key.toLowerCase()
          )
          if (targetIndex !== -1 && rowValues[targetIndex] !== undefined) {
            rowObj[colSchema.key] = rowValues[targetIndex]
            rowObj[colSchema.label || colSchema.key] = rowValues[targetIndex]
          } else {
            rowObj[colSchema.key] = ''
          }
        })
        mapped.push(rowObj)
      })

      parsedRows.value = mapped
    } catch (err) {
      alert('Gagal membaca berkas CSV: ' + err.message)
    }
  }
  reader.readAsText(file)
}

const resetFile = () => {
  selectedFileName.value = ''
  parsedRows.value = []
  importErrors.value = []
  processedCount.value = 0
}

const closeModal = () => {
  if (importing.value) return
  resetFile()
  emit('close')
}

const startImport = async () => {
  if (parsedRows.value.length === 0) return
  importing.value = true
  processedCount.value = 0
  importErrors.value = []

  let successCount = 0

  for (let i = 0; i < parsedRows.value.length; i++) {
    const row = parsedRows.value[i]
    try {
      await props.onImportRow(row, i)
      successCount++
    } catch (err) {
      const errMsg = err.response?.data?.message || err.message || 'Gagal menyimpan'
      const identifier = row.nama || row.nama_lengkap || row.nis || row.nip || `Baris ${i + 2}`
      importErrors.value.push(`${identifier}: ${errMsg}`)
    } finally {
      processedCount.value++
    }
  }

  importing.value = false

  if (successCount > 0) {
    emit('success', { total: parsedRows.value.length, success: successCount, failed: importErrors.value.length })
    if (importErrors.value.length === 0) {
      closeModal()
    }
  }
}
</script>

<style scoped>
.app-modal-overlay {
  position: fixed;
  inset: 0;
  background-color: rgba(15, 23, 42, 0.65);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1055;
  padding: 1.5rem 1rem;
}

.upload-dropzone {
  cursor: pointer;
  transition: all 0.2s ease-in-out;
}
.upload-dropzone:hover {
  background-color: #f1f5f9;
  border-color: #0d6efd;
}
</style>
