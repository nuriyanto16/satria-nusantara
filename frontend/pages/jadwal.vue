<template>
  <div id="page-jadwal">
    <!-- Header -->
    <div class="page-header">
      <div>
        <h1 class="page-title">Kelola Jadwal Latihan</h1>
        <div class="page-subtitle">Semua sesi latihan aktif dan pembuatan absensi QR</div>
      </div>
      <div class="header-actions-group">
        <div class="view-toggle">
          <button :class="['view-toggle-btn', { active: currentView === 'list' }]" @click="currentView = 'list'">
            <i class="ti ti-list"></i> Daftar
          </button>
          <button :class="['view-toggle-btn', { active: currentView === 'calendar' }]" @click="currentView = 'calendar'">
            <i class="ti ti-calendar"></i> Kalender
          </button>
        </div>
        <button class="btn btn-primary" @click="openCreateModal"><i class="ti ti-plus"></i> Tambah Sesi</button>
      </div>
    </div>

    <!-- Filters -->
    <div class="filter-bar">
      <select v-model="filterCabang" class="filter-select" @change="onFilterCabangChange">
        <option value="">Semua Cabang</option>
        <option v-for="c in cabangList" :key="c.id" :value="c.id">{{ c.nama }}</option>
      </select>
      <select v-model="filterUnit" class="filter-select" :disabled="!filterCabang" @change="fetchSesi">
        <option value="">Semua Unit</option>
        <option v-for="u in filterUnitList" :key="u.id" :value="u.id">{{ u.nama }}</option>
      </select>
      <input v-if="currentView === 'list'" v-model="filterTanggal" type="date" class="filter-select" @change="fetchSesi" />
    </div>

    <!-- Loading Skeletons -->
    <div v-if="loading" class="table-card">
      <table class="data-table">
        <thead>
          <tr>
            <th>Tanggal</th>
            <th>Unit Latihan</th>
            <th>Pelatih</th>
            <th>Waktu</th>
            <th>Lokasi</th>
            <th>Jenis</th>
            <th>Materi</th>
            <th style="width:140px">Aksi</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="n in 5" :key="n">
            <td><div class="skeleton-box" style="width:80px; height:14px;"></div></td>
            <td><div class="skeleton-box" style="width:100px; height:14px;"></div></td>
            <td><div class="skeleton-box" style="width:80px; height:14px;"></div></td>
            <td><div class="skeleton-box" style="width:90px; height:14px;"></div></td>
            <td><div class="skeleton-box" style="width:120px; height:14px;"></div></td>
            <td><div class="skeleton-box" style="width:80px; height:14px; border-radius:10px;"></div></td>
            <td><div class="skeleton-box" style="width:150px; height:14px;"></div></td>
            <td><div class="skeleton-box" style="width:60px; height:24px; border-radius:4px;"></div></td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Calendar View -->
    <div v-else-if="currentView === 'calendar'" class="calendar-view">
      <div class="calendar-header">
        <button class="btn btn-outline" @click="prevMonth"><i class="ti ti-chevron-left"></i></button>
        <h3 class="calendar-month-title">{{ monthNames[calendarMonth] }} {{ calendarYear }}</h3>
        <button class="btn btn-outline" @click="nextMonth"><i class="ti ti-chevron-right"></i></button>
      </div>

      <div class="calendar-grid">
        <!-- Weekday Headers -->
        <div class="calendar-day-header">Sen</div>
        <div class="calendar-day-header">Sel</div>
        <div class="calendar-day-header">Rab</div>
        <div class="calendar-day-header">Kam</div>
        <div class="calendar-day-header">Jum</div>
        <div class="calendar-day-header">Sab</div>
        <div class="calendar-day-header">Min</div>

        <!-- Cells -->
        <div v-for="(day, index) in daysInMonth" :key="index" :class="['calendar-cell', { empty: !day.day }]">
          <template v-if="day.day">
            <span class="calendar-day-num">{{ day.day }}</span>
            <div class="calendar-sessions">
              <div v-for="s in day.sessions" :key="s.id" class="calendar-session-item" :title="`${s.unit_nama} - ${s.jenis}`" @click="showCalendarDetail(s)">
                <span class="csi-time">{{ s.jam_mulai.substring(0, 5) }}</span>
                <span class="csi-title">{{ s.unit_nama }}</span>
              </div>
            </div>
          </template>
        </div>
      </div>
    </div>

    <!-- Empty State -->
    <div v-else-if="sesiData.length === 0" class="empty-state">
      <i class="ti ti-calendar-off"></i> Belum ada jadwal latihan terdaftar untuk filter ini.
    </div>

    <!-- Table View -->
    <div v-else class="table-card">
      <table class="data-table">
        <thead>
          <tr>
            <th>Tanggal</th>
            <th>Unit Latihan</th>
            <th>Pelatih</th>
            <th>Waktu</th>
            <th>Lokasi</th>
            <th>Jenis</th>
            <th>Materi</th>
            <th style="width:140px">Aksi</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="s in paginatedSesi" :key="s.id">
            <td>
              <div class="sesi-date">{{ formatDate(s.tanggal) }}</div>
            </td>
            <td>{{ s.unit_nama || '-' }}</td>
            <td>{{ s.pelatih_nama || (s.pelatih_id && !s.pelatih_id.includes('-') ? s.pelatih_id : 'Pak Budi Susanto') }}</td>
            <td class="time-cell">{{ s.jam_mulai }} – {{ s.jam_selesai || 'Selesai' }}</td>
            <td class="text-muted">{{ s.lokasi }}</td>
            <td><span :class="['jenis-badge', getJenisClass(s.jenis)]">{{ s.jenis }}</span></td>
            <td class="text-muted">{{ s.materi || '-' }}</td>
            <td>
              <div class="action-btns">
                <button class="action-btn qr" @click="generateQR(s.id)" title="Generate QR Absensi"><i class="ti ti-qrcode"></i></button>
                <button class="action-btn edit" @click="openEditModal(s)" title="Edit Sesi" style="background:#e0f5fb; color:var(--biru);"><i class="ti ti-edit"></i></button>
                <button class="action-btn delete" @click="deleteSesi(s.id)" title="Hapus Sesi"><i class="ti ti-trash"></i></button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
      <Pagination 
        v-model:currentPage="page" 
        v-model:itemsPerPage="limit" 
        :totalItems="sesiData.length" 
      />
    </div>

    <!-- MODAL DETAIL KALENDER -->
    <div v-if="selectedCalendarSession" class="modal-overlay" @click.self="selectedCalendarSession = null">
      <div class="modal-card">
        <div class="modal-header">
          <h2 class="modal-title">Detail Sesi Latihan</h2>
          <button class="modal-close" @click="selectedCalendarSession = null"><i class="ti ti-x"></i></button>
        </div>
        <div class="modal-detail-content">
          <div class="mdc-row"><strong>Tanggal:</strong> {{ formatDate(selectedCalendarSession.tanggal) }}</div>
          <div class="mdc-row"><strong>Unit Latihan:</strong> {{ selectedCalendarSession.unit_nama }}</div>
          <div class="mdc-row"><strong>Pelatih:</strong> {{ selectedCalendarSession.pelatih_nama || (selectedCalendarSession.pelatih_id && !selectedCalendarSession.pelatih_id.includes("-") ? selectedCalendarSession.pelatih_id : "Pak Budi Susanto") }}</div>
          <div class="mdc-row"><strong>Waktu:</strong> {{ selectedCalendarSession.jam_mulai }} – {{ selectedCalendarSession.jam_selesai || 'Selesai' }}</div>
          <div class="mdc-row"><strong>Lokasi:</strong> {{ selectedCalendarSession.lokasi }}</div>
          <div class="mdc-row"><strong>Jenis:</strong> <span :class="['jenis-badge', getJenisClass(selectedCalendarSession.jenis)]">{{ selectedCalendarSession.jenis }}</span></div>
          <div class="mdc-row"><strong>Materi:</strong> {{ selectedCalendarSession.materi || '-' }}</div>
        </div>
        <div class="modal-footer" style="margin-top: 20px;">
          <button class="btn btn-outline" @click="selectedCalendarSession = null">Tutup</button>
          <button class="btn btn-primary" @click="generateQR(selectedCalendarSession.id); selectedCalendarSession = null">
            <i class="ti ti-qrcode"></i> QR Absensi
          </button>
          <button class="btn btn-outline" style="border-color:var(--biru); color:var(--biru);" @click="openEditModal(selectedCalendarSession); selectedCalendarSession = null">
            <i class="ti ti-edit"></i> Edit
          </button>
          <button class="action-btn delete" style="width: auto; padding: 0 12px; border-radius: var(--r8)" @click="deleteSesi(selectedCalendarSession.id); selectedCalendarSession = null">
            <i class="ti ti-trash"></i> Hapus
          </button>
        </div>
      </div>
    </div>

    <!-- MODAL TAMBAH SESI -->
    <div v-if="showModal" class="modal-overlay">
      <div class="modal-card">
        <div class="modal-header">
          <h2 class="modal-title">{{ editingId ? 'Edit Sesi Latihan' : 'Tambah Sesi Latihan Baru' }}</h2>
          <button class="modal-close" @click="showModal = false"><i class="ti ti-x"></i></button>
        </div>
        <form @submit.prevent="saveSesi" class="modal-form">
          <div class="form-group">
            <label class="form-label">Cabang Latihan</label>
            <select v-model="form.cabang_id" class="form-select" @change="onFormCabangChange" required>
              <option value="" disabled>Pilih Cabang</option>
              <option v-for="c in cabangList" :key="c.id" :value="c.id">{{ c.nama }}</option>
            </select>
          </div>
          <div class="form-group">
            <label class="form-label">Unit Latihan</label>
            <select v-model="form.unit_id" class="form-select" :disabled="!form.cabang_id || loadingUnits" required>
              <option value="" disabled>{{ loadingUnits ? 'Memuat unit...' : 'Pilih Unit' }}</option>
              <option v-for="u in formUnitList" :key="u.id" :value="u.id">{{ u.nama }}</option>
            </select>
          </div>
          <div class="form-group">
            <label class="form-label">Pelatih Sesi</label>
            <select v-model="form.pelatih_id" class="form-select" :disabled="!form.cabang_id || loadingPelatih" required>
              <option value="" disabled>{{ loadingPelatih ? 'Memuat pelatih...' : 'Pilih Pelatih' }}</option>
              <option v-for="p in pelatihList" :key="p.id" :value="p.id">{{ p.nama_lengkap }} ({{ p.jenis }})</option>
            </select>
          </div>
          <div class="form-row">
            <div class="form-group half">
              <label class="form-label">Jenis Sesi</label>
              <select v-model="form.jenis" class="form-select" required>
                <option value="Latihan Rutin">Latihan Rutin</option>
                <option value="Latihan Khusus">Latihan Khusus</option>
                <option value="Latihan Pelatih">Latihan Pelatih</option>
              </select>
            </div>
            <div class="form-group half">
              <label class="form-label">Tanggal</label>
              <input v-model="form.tanggal" type="date" class="form-input" required />
            </div>
          </div>
          <div class="form-row">
            <div class="form-group half">
              <label class="form-label">Jam Mulai</label>
              <input v-model="form.jam_mulai" type="time" class="form-input" required />
            </div>
            <div class="form-group half">
              <label class="form-label">Jam Selesai</label>
              <input v-model="form.jam_selesai" type="time" class="form-input" />
            </div>
          </div>
          <div class="form-group">
            <label class="form-label">Lokasi</label>
            <input v-model="form.lokasi" type="text" class="form-input" placeholder="Contoh: Lapangan Malioboro..." required />
          </div>
          <div class="form-group">
            <label class="form-label">Materi Latihan</label>
            <input v-model="form.materi" type="text" class="form-input" placeholder="Contoh: Jurus 1 - 3, Napas Pembinaan..." />
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-outline" @click="showModal = false">Batal</button>
            <button type="submit" class="btn btn-primary" :disabled="submitting">
              <i v-if="submitting" class="ti ti-loader-2 spin"></i> Simpan Sesi
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- MODAL QR CODE ABSENSI -->
    <div v-if="showQRModal" class="modal-overlay">
      <div class="modal-card qr-modal">
        <div class="modal-header">
          <h2 class="modal-title">QR Code Absensi Kehadiran</h2>
          <button class="modal-close" @click="showQRModal = false"><i class="ti ti-x"></i></button>
        </div>
        <div class="qr-container">
          <img :src="`https://api.qrserver.com/v1/create-qr-code/?size=250x250&data=${qrToken}`" alt="QR Code" class="qr-image" />
          <div class="qr-token-label">Token: <code>{{ qrToken }}</code></div>
          <div class="qr-hint">Masa berlaku QR code ini terbatas (4 jam). Tunjukkan ke anggota untuk dipindai via aplikasi mobile.</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({ title: 'Kelola Jadwal Latihan' })

const api = useApi()

const currentView = ref('list') // 'list' or 'calendar'
const filterCabang = ref('')
const filterUnit = ref('')
const filterTanggal = ref('')
const showModal = ref(false)
const showQRModal = ref(false)
const submitting = ref(false)
const editingId = ref<string | null>(null)
const loading = ref(false)
const loadingUnits = ref(false)
const loadingPelatih = ref(false)

const sesiData = ref<any[]>([])
const page = ref(1)
const limit = ref(10)
const paginatedSesi = computed(() => {
  return sesiData.value.slice((page.value - 1) * limit.value, page.value * limit.value)
})
const cabangList = ref<any[]>([])
const filterUnitList = ref<any[]>([])
const formUnitList = ref<any[]>([])
const pelatihList = ref<any[]>([])
const qrToken = ref('')

const selectedCalendarSession = ref<any>(null)

// Calendar States
const calendarDate = ref(new Date())
const calendarYear = computed(() => calendarDate.value.getFullYear())
const calendarMonth = computed(() => calendarDate.value.getMonth())
const monthNames = ["Januari", "Februari", "Maret", "April", "Mei", "Juni", "Juli", "Agustus", "September", "Oktober", "November", "Desember"]

const daysInMonth = computed(() => {
  const year = calendarYear.value
  const month = calendarMonth.value
  
  const firstDay = new Date(year, month, 1)
  let startDayIndex = firstDay.getDay() - 1
  if (startDayIndex < 0) startDayIndex = 6
  
  const totalDays = new Date(year, month + 1, 0).getDate()
  const days: any[] = []
  
  for (let i = 0; i < startDayIndex; i++) {
    days.push({ day: null, dateStr: '', sessions: [] })
  }
  
  for (let d = 1; d <= totalDays; d++) {
    const dStr = `${year}-${String(month + 1).padStart(2, '0')}-${String(d).padStart(2, '0')}`
    const matchedSessions = sesiData.value.filter(s => {
      return s.tanggal && s.tanggal.startsWith(dStr)
    })
    
    days.push({
      day: d,
      dateStr: dStr,
      sessions: matchedSessions
    })
  }
  return days
})

const prevMonth = () => {
  calendarDate.value = new Date(calendarYear.value, calendarMonth.value - 1, 1)
}
const nextMonth = () => {
  calendarDate.value = new Date(calendarYear.value, calendarMonth.value + 1, 1)
}

const showCalendarDetail = (s: any) => {
  selectedCalendarSession.value = s
}

const form = ref({
  cabang_id: '',
  unit_id: '',
  pelatih_id: '',
  jenis: 'Latihan Rutin',
  tanggal: '',
  jam_mulai: '',
  jam_selesai: '',
  lokasi: '',
  materi: ''
})

const fetchCabang = async () => {
  try {
    const res = await api.get('/organization/cabang?limit=100')
    cabangList.value = res.data || []
    
    // Auto-select first branch
    if (cabangList.value.length > 0) {
      filterCabang.value = cabangList.value[0].id
      await onFilterCabangChange()
    }
  } catch (e) {
    console.error(e)
  }
}

const onFilterCabangChange = async () => {
  filterUnit.value = ''
  filterUnitList.value = []
  if (!filterCabang.value) {
    fetchSesi()
    return
  }
  try {
    const res = await api.get(`/organization/cabang/${filterCabang.value}/unit`)
    filterUnitList.value = res || []
    
    // Auto-select first unit
    if (filterUnitList.value.length > 0) {
      filterUnit.value = filterUnitList.value[0].id
    }
    fetchSesi()
  } catch (e) {
    console.error(e)
  }
}

const onFormCabangChange = async () => {
  form.value.unit_id = ''
  form.value.pelatih_id = ''
  formUnitList.value = []
  pelatihList.value = []
  if (!form.value.cabang_id) return

  loadingUnits.value = true
  loadingPelatih.value = true
  try {
    const units = await api.get(`/organization/cabang/${form.value.cabang_id}/unit`)
    formUnitList.value = units || []

    const pelatih = await api.get(`/organization/pelatih?cabang_id=${form.value.cabang_id}`)
    pelatihList.value = pelatih || []
  } catch (e) {
    console.error(e)
  } finally {
    loadingUnits.value = false
    loadingPelatih.value = false
  }
}

const fetchSesi = async () => {
  loading.value = true
  try {
    const dateQuery = currentView.value === 'list' ? filterTanggal.value : ''
    const data = await api.get(`/training/sesi?unit_id=${filterUnit.value}&tanggal=${dateQuery}`)
    sesiData.value = data || []
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

const openCreateModal = () => {
  editingId.value = null
  form.value = {
    cabang_id: '',
    unit_id: '',
    pelatih_id: '',
    jenis: 'Latihan Rutin',
    tanggal: new Date().toISOString().substring(0, 10),
    jam_mulai: '07:00',
    jam_selesai: '09:00',
    lokasi: '',
    materi: ''
  }
  formUnitList.value = []
  pelatihList.value = []
  showModal.value = true
}

const openEditModal = async (s: any) => {
  editingId.value = s.id
  form.value = {
    cabang_id: s.cabang_id || '',
    unit_id: s.unit_id || '',
    pelatih_id: s.pelatih_id || '',
    jenis: s.jenis || 'Latihan Rutin',
    tanggal: s.tanggal ? s.tanggal.substring(0, 10) : '',
    jam_mulai: s.jam_mulai ? s.jam_mulai.substring(0, 5) : '',
    jam_selesai: s.jam_selesai ? s.jam_selesai.substring(0, 5) : '',
    lokasi: s.lokasi || '',
    materi: s.materi || ''
  }
  
  if (form.value.cabang_id) {
    loadingUnits.value = true
    loadingPelatih.value = true
    try {
      const units = await api.get(`/organization/cabang/${form.value.cabang_id}/unit`)
      formUnitList.value = units || []

      const pelatih = await api.get(`/organization/pelatih?cabang_id=${form.value.cabang_id}`)
      pelatihList.value = pelatih || []
    } catch (e) {
      console.error(e)
    } finally {
      loadingUnits.value = false
      loadingPelatih.value = false
    }
  }
  showModal.value = true
}

const saveSesi = async () => {
  submitting.value = true
  try {
    const payload = {
      unit_id: form.value.unit_id,
      pelatih_id: form.value.pelatih_id,
      jenis: form.value.jenis,
      tanggal: form.value.tanggal,
      jam_mulai: form.value.jam_mulai,
      jam_selesai: form.value.jam_selesai,
      lokasi: form.value.lokasi,
      materi: form.value.materi
    }
    
    if (editingId.value) {
      await api.put(`/training/sesi/${editingId.value}`, payload)
    } else {
      await api.post('/training/sesi', payload)
    }
    
    showModal.value = false
    fetchSesi()
  } catch (e: any) {
    alert(e.message || 'Gagal menyimpan sesi latihan')
  } finally {
    submitting.value = false
  }
}

const deleteSesi = async (id: string) => {
  if (!confirm('Apakah Anda yakin ingin menghapus sesi latihan ini?')) return
  try {
    await api.delete(`/training/sesi/${id}`)
    fetchSesi()
  } catch (e: any) {
    alert(e.message || 'Gagal menghapus sesi')
  }
}

const generateQR = async (id: string) => {
  try {
    const res = await api.post(`/training/sesi/${id}/qr`)
    qrToken.value = res.qr_token
    showQRModal.value = true
  } catch (e: any) {
    alert(e.message || 'Gagal generate QR Code')
  }
}

const getJenisClass = (jenis: string) => {
  if (jenis === 'Latihan Rutin') return 'rutin'
  if (jenis === 'Latihan Khusus') return 'khusus'
  return 'pelatih'
}

const formatDate = (dateStr: string) => {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  return d.toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' })
}

onMounted(() => {
  fetchCabang()
  fetchSesi()
})

// Refetch sessions when changing view
watch(currentView, () => {
  fetchSesi()
})
</script>

<style scoped>
.page-header { display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 20px; }
.page-title { font-size: 20px; font-weight: 800; margin: 0; }
.page-subtitle { font-size: 12px; color: var(--text3); margin-top: 2px; }

.header-actions-group { display: flex; gap: 10px; align-items: center; }

.btn { display: inline-flex; align-items: center; gap: 6px; padding: 8px 14px; border-radius: var(--r8); font-size: 12px; font-weight: 600; cursor: pointer; border: none; }
.btn-primary { background: var(--hijau); color: #fff; }
.btn-outline { background: #fff; border: 1px solid var(--border); color: var(--text2); }

.filter-bar { display: flex; gap: 8px; margin-bottom: 14px; flex-wrap: wrap; }
.filter-select { padding: 7px 10px; border: 1px solid var(--border); border-radius: var(--r8); font-size: 12px; background: var(--card); color: var(--text2); }

.table-card { background: var(--card); border: 1px solid var(--border); border-radius: var(--r12); overflow: hidden; }
.data-table { width: 100%; border-collapse: collapse; }
.data-table th { padding: 10px 14px; font-size: 10px; font-weight: 700; color: var(--text3); text-transform: uppercase; text-align: left; background: var(--surface); border-bottom: 1px solid var(--border); }
.data-table td { padding: 10px 14px; font-size: 12px; border-bottom: 1px solid var(--border); vertical-align: middle; }
.data-table tr:last-child td { border-bottom: none; }
.data-table tr:hover td { background: var(--surface); }

.sesi-date { font-weight: 700; font-size: 12px; }
.time-cell { font-weight: 600; }
.text-muted { color: var(--text3); }

.jenis-badge { display: inline-block; padding: 2px 8px; border-radius: 10px; font-size: 10px; font-weight: 700; }
.jenis-badge.rutin { background: var(--hijau3); color: var(--hijau); }
.jenis-badge.khusus { background: #fff8e0; color: #a07000; }
.jenis-badge.pelatih { background: #e0f5fb; color: var(--biru); }

.action-btns { display: flex; gap: 4px; }
.action-btn { width: 28px; height: 28px; border: 1px solid var(--border); border-radius: var(--r6); background: #fff; cursor: pointer; display: flex; align-items: center; justify-content: center; font-size: 14px; color: var(--text2); transition: all .15s; }
.action-btn:hover { border-color: var(--hijau); color: var(--hijau); background: var(--hijau3); }
.action-btn.qr { color: var(--hijau); border-color: var(--hijau4); }
.action-btn.delete { color: var(--merah); border-color: #fca5a5; }
.action-btn.delete:hover { background: #fde8e8; }

/* Modal Styles */
.modal-overlay { position: fixed; top: 0; left: 0; right: 0; bottom: 0; background: rgba(0,0,0,0.5); z-index: 1000; display: flex; align-items: center; justify-content: center; }
.modal-card { background: var(--card); border-radius: var(--r12); width: 100%; max-width: 480px; padding: 20px; box-shadow: var(--shadow-lg); }
.modal-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px; }
.modal-title { font-size: 16px; font-weight: 700; }
.modal-close { background: none; border: none; font-size: 16px; cursor: pointer; color: var(--text3); }
.modal-form { display: flex; flex-direction: column; gap: 12px; }
.form-group { display: flex; flex-direction: column; gap: 4px; }
.form-row { display: flex; gap: 10px; }
.form-group.half { flex: 1; }
.form-label { font-size: 11px; font-weight: 600; color: var(--text2); }
.form-input, .form-select { padding: 8px 12px; border: 1.5px solid var(--border2); border-radius: var(--r8); font-size: 13px; outline: none; background: #fff; width: 100%; }
.form-input:focus, .form-select:focus { border-color: var(--hijau); }
.modal-footer { display: flex; justify-content: flex-end; gap: 8px; margin-top: 12px; }

.modal-detail-content { display: flex; flex-direction: column; gap: 10px; font-size: 13px; }
.mdc-row { border-bottom: 1px solid var(--border); padding-bottom: 8px; }

.qr-modal { max-width: 320px; text-align: center; }
.qr-container { display: flex; flex-direction: column; align-items: center; gap: 12px; padding: 10px 0; }
.qr-image { width: 180px; height: 180px; border: 1px solid var(--border); padding: 6px; border-radius: var(--r8); background: #fff; }
.qr-token-label { font-size: 12px; color: var(--text2); }
.qr-token-label code { background: var(--surface); padding: 2px 6px; border-radius: 4px; font-weight: 700; }
.qr-hint { font-size: 10px; color: var(--text3); line-height: 1.4; }

.loading-state, .empty-state { display: flex; flex-direction: column; align-items: center; justify-content: center; padding: 40px; font-size: 13px; color: var(--text3); gap: 10px; background: var(--card); border-radius: var(--r12); border: 1px dashed var(--border); margin-top: 10px; }
.loading-state i { font-size: 20px; }
.empty-state i { font-size: 32px; color: var(--text3); }

@keyframes spin { from { transform: rotate(0deg); } to { transform: rotate(360deg); } }
.spin { animation: spin .8s linear infinite; }

/* Skeleton Loader */
@keyframes pulse {
  0% { opacity: 0.6; }
  50% { opacity: 1; }
  100% { opacity: 0.6; }
}
.skeleton-box {
  animation: pulse 1.5s infinite ease-in-out;
  background-color: rgba(0, 0, 0, 0.08);
  border-radius: var(--r6);
}

/* View Toggle Button Style */
.view-toggle { display: flex; border: 1px solid var(--border); border-radius: var(--r8); overflow: hidden; background: #fff; }
.view-toggle-btn { background: none; border: none; padding: 6px 12px; font-size: 11px; font-weight: 600; cursor: pointer; color: var(--text2); display: flex; align-items: center; gap: 4px; }
.view-toggle-btn.active { background: var(--hijau); color: #fff; }

/* Calendar CSS styling */
.calendar-view { background: var(--card); border: 1px solid var(--border); border-radius: var(--r12); padding: 20px; box-shadow: var(--shadow-sm); }
.calendar-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; }
.calendar-month-title { font-size: 16px; font-weight: 800; margin: 0; color: var(--text1); }
.calendar-grid { display: grid; grid-template-columns: repeat(7, 1fr); gap: 6px; }
.calendar-day-header { text-align: center; font-weight: 700; font-size: 11px; color: var(--text3); text-transform: uppercase; padding-bottom: 8px; border-bottom: 1px solid var(--border); }
.calendar-cell { min-height: 100px; border: 1px solid var(--border); border-radius: var(--r8); padding: 8px; background: var(--surface); display: flex; flex-direction: column; gap: 6px; transition: border-color .15s; }
.calendar-cell:hover:not(.empty) { border-color: var(--hijau); }
.calendar-cell.empty { background: transparent; border: none; }
.calendar-day-num { font-size: 12px; font-weight: 700; color: var(--text2); align-self: flex-start; }
.calendar-sessions { display: flex; flex-direction: column; gap: 4px; overflow-y: auto; max-height: 70px; }
.calendar-session-item { font-size: 9px; font-weight: 700; background: var(--hijau3); color: var(--hijau); padding: 3px 6px; border-radius: 4px; cursor: pointer; display: flex; flex-direction: column; gap: 1px; line-height: 1.2; text-align: left; transition: all .15s; border-left: 2.5px solid var(--hijau); }
.calendar-session-item:hover { background: var(--hijau4); transform: translateY(-1px); }
.csi-time { font-size: 8px; color: var(--hijau2); }
.csi-title { font-weight: 800; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }

/* Responsive tweaks */
@media(max-width: 768px) {
  .page-header { flex-direction: column; gap: 12px; }
  .calendar-grid { grid-template-columns: repeat(7, minmax(40px, 1fr)); font-size: 10px; }
  .calendar-cell { min-height: 70px; padding: 4px; }
  .calendar-session-item { padding: 2px 4px; font-size: 8px; }
}
</style>
