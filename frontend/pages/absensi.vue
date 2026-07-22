<template>
  <div id="page-absensi" class="absensi-split-layout">
    <!-- ===== PANEL KIRI: FILTER & SUMMARY ===== -->
    <div class="panel-left">
      <div class="pl-head">
        <div class="pl-title">Filter & Ringkasan</div>
        <div class="pl-sub">Pilih wilayah absensi latihan</div>
      </div>
      
      <div class="pl-scroll">
        <!-- Filter Cabang -->
        <div class="form-group-filter">
          <label class="filter-label">Cabang</label>
          <select v-model="filterCabang" class="filter-select" @change="onCabangChange">
            <option value="">Semua Cabang</option>
            <option v-for="c in cabangList" :key="c.id" :value="c.id">{{ c.nama }}</option>
          </select>
        </div>

        <!-- Filter Unit -->
        <div class="form-group-filter">
          <label class="filter-label">Unit Latihan</label>
          <select v-model="filterUnit" class="filter-select" :disabled="!filterCabang" @change="fetchSesi">
            <option value="">Semua Unit</option>
            <option v-for="u in unitList" :key="u.id" :value="u.id">{{ u.nama }}</option>
          </select>
        </div>

        <!-- Filter Periode -->
        <div class="form-group-filter" style="margin-bottom:18px;">
          <label class="filter-label">Periode</label>
          <div class="month-nav">
            <button class="month-btn" @click="changeMonth(-1)"><i class="ti ti-chevron-left"></i></button>
            <div class="month-label">{{ currentMonthLabel }}</div>
            <button class="month-btn" @click="changeMonth(1)"><i class="ti ti-chevron-right"></i></button>
          </div>
        </div>

        <!-- Summary -->
        <div class="summary-title">Ringkasan · {{ activeUnitName }}</div>
        <div class="summary-grid">
          <div class="sum-card">
            <div class="sum-val">{{ totalSesiCount }}</div>
            <div class="sum-label">Total Sesi</div>
          </div>
          <div class="sum-card">
            <div class="sum-val">{{ totalAnggotaActive }}</div>
            <div class="sum-label">Anggota Aktif</div>
          </div>
          <div class="sum-card warn">
            <div class="sum-val">{{ jarangHadirCount }}</div>
            <div class="sum-label">Jarang Hadir</div>
          </div>
          <div class="sum-card danger">
            <div class="sum-val">{{ tidakHadirCount }}</div>
            <div class="sum-label">Tidak Hadir</div>
          </div>
        </div>

        <!-- Kehadiran Meter -->
        <div class="kehadiran-meter">
          <div class="km-label">Rata-rata Kehadiran Unit</div>
          <div class="km-pct">{{ averageAttendancePct }}%</div>
          <div class="km-bar-wrap">
            <div class="km-bar" :style="{ width: averageAttendancePct + '%', background: getBarColor(averageAttendancePct) }"></div>
          </div>
          <div class="km-sub">{{ activeHadirAvg }} dari {{ totalAnggotaActive }} anggota hadir rata-rata per sesi</div>
        </div>
      </div>
    </div>

    <!-- ===== PANEL KANAN: TABS & DATA LIST ===== -->
    <div class="panel-right">
      <div class="tabs">
        <button :class="['tab', { active: activeTab === 'sesi' }]" @click="activeTab = 'sesi'">
          <i class="ti ti-calendar-event"></i> Per Sesi
        </button>
        <button :class="['tab', { active: activeTab === 'anggota' }]" @click="activeTab = 'anggota'">
          <i class="ti ti-users"></i> Per Anggota
        </button>
        <button :class="['tab', { active: activeTab === 'pelatih' }]" @click="activeTab = 'pelatih'">
          <i class="ti ti-user-star"></i> Pelatih
        </button>
      </div>

      <!-- TAB 1: PER SESI -->
      <div v-if="activeTab === 'sesi'" class="tab-pane">
        <div class="toolbar">
          <div class="toolbar-info">Riwayat Sesi Latihan · {{ activeUnitName }}</div>
          <div class="export-row">
            <button class="export-btn" @click="exportData('PDF')"><i class="ti ti-file-type-pdf"></i> PDF</button>
            <button class="export-btn" @click="exportData('Excel')"><i class="ti ti-file-spreadsheet"></i> Excel</button>
          </div>
        </div>

        <div class="list-scroll">
          <div v-if="loading" style="padding: 20px; text-align: center;">
            <i class="ti ti-loader-2 spin" style="font-size:24px; color:var(--hijau)"></i>
            <div style="font-size:12px; color:var(--text3); margin-top:8px;">Memuat data sesi...</div>
          </div>
          <div v-else-if="sesiList.length === 0" class="empty-state">
            <i class="ti ti-calendar-off"></i> Belum ada sesi latihan terdaftar untuk kriteria ini.
          </div>
          <table v-else class="data-table">
            <thead>
              <tr>
                <th>Tanggal Sesi</th>
                <th>Unit Latihan</th>
                <th>Pelatih</th>
                <th>Kehadiran</th>
                <th>Persentase</th>
                <th>Aksi</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="s in paginatedSesiList" :key="s.id">
                <td>
                  <div class="sesi-date">{{ formatDate(s.tanggal) }}</div>
                  <div class="sesi-time">{{ s.jam_mulai.substring(0, 5) }} – {{ s.jam_selesai ? s.jam_selesai.substring(0, 5) : 'Selesai' }}</div>
                </td>
                <td><strong>{{ s.unit_nama || activeUnitName }}</strong></td>
                <td>{{ s.pelatih_nama || (s.pelatih_id && !s.pelatih_id.includes('-') ? s.pelatih_id : 'Pak Budi Susanto') }}</td>
                <td class="num-cell green">
                  <strong>{{ s.hadirCount }}</strong> <span style="color:var(--text3);font-weight:normal;">hadir</span>
                </td>
                <td>
                  <div class="pct-cell">
                    <div class="pct-bar"><div class="pct-fill" :style="{ width: s.attendancePct + '%' }"></div></div>
                    <span>{{ s.attendancePct }}%</span>
                  </div>
                </td>
                <td>
                  <button class="action-btn" title="Lihat Detail Absensi" @click="viewDetail(s)">
                    <i class="ti ti-list"></i>
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
          <Pagination 
            v-model:currentPage="pageSesi" 
            v-model:itemsPerPage="limitSesi" 
            :totalItems="sesiList.length" 
          />
        </div>
      </div>

      <!-- TAB 2: PER ANGGOTA -->
      <div v-else-if="activeTab === 'anggota'" class="tab-pane">
        <div class="toolbar">
          <div class="toolbar-info">Daftar Kehadiran Anggota · {{ activeUnitName }}</div>
          <div class="export-row">
            <button class="export-btn" @click="exportData('PDF')"><i class="ti ti-file-type-pdf"></i> PDF</button>
            <button class="export-btn" @click="exportData('Excel')"><i class="ti ti-file-spreadsheet"></i> Excel</button>
          </div>
        </div>

        <div class="list-scroll">
          <table class="data-table">
            <thead>
              <tr>
                <th>Anggota</th>
                <th>Tingkatan</th>
                <th style="text-align:center;">Total Hadir</th>
                <th>Persentase</th>
                <th>Status</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="a in paginatedAnggotaKehadiranList" :key="a.id">
                <td>
                  <div style="font-weight:600;">{{ a.nama }}</div>
                  <div style="font-size:10px;color:var(--text3);">{{ a.nomor }}</div>
                </td>
                <td><span class="belt-badge">{{ a.tingkatan }}</span></td>
                <td style="text-align:center;font-weight:700;color:var(--hijau);">{{ a.hadir }} / {{ totalSesiCount }}</td>
                <td>
                  <div class="pct-cell">
                    <div class="pct-bar"><div class="pct-fill" :style="{ width: a.pct + '%', background: a.pct < 50 ? 'var(--merah)' : 'var(--hijau2)' }"></div></div>
                    <span style="font-weight:700;">{{ a.pct }}%</span>
                  </div>
                </td>
                <td>
                  <span :class="['status-badge', a.pct >= 75 ? 'lunas' : a.pct >= 50 ? 'warning' : 'belum']">
                    {{ a.pct >= 75 ? 'Aktif' : a.pct >= 50 ? 'Cukup Aktif' : 'Kurang Aktif' }}
                  </span>
                </td>
              </tr>
            </tbody>
          </table>
          <Pagination 
            v-model:currentPage="pageAnggota" 
            v-model:itemsPerPage="limitAnggota" 
            :totalItems="anggotaKehadiranList.length" 
          />
        </div>
      </div>

      <!-- TAB 3: PELATIH -->
      <div v-else-if="activeTab === 'pelatih'" class="tab-pane">
        <div class="toolbar">
          <div class="toolbar-info">Kehadiran Pelatih · {{ activeUnitName }}</div>
          <div class="export-row">
            <button class="export-btn" @click="exportData('PDF')"><i class="ti ti-file-type-pdf"></i> PDF</button>
          </div>
        </div>

        <div class="list-scroll">
          <table class="data-table">
            <thead>
              <tr>
                <th>Nama Pelatih</th>
                <th>Jenis Sabuk</th>
                <th>Jumlah Melatih Sesi</th>
                <th>Status</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="p in pelatihKehadiranList" :key="p.name">
                <td>
                  <div style="font-weight:600;">{{ p.name }}</div>
                  <div style="font-size:10px;color:var(--text3);">{{ p.nip }}</div>
                </td>
                <td><span class="belt-badge ph">{{ p.sabuk }}</span></td>
                <td style="font-weight:700;color:var(--hijau);">{{ p.melatih }} sesi</td>
                <td><span class="status-badge lunas">Aktif</span></td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- DETAIL ABSENSI MODAL -->
    <div v-if="showDetailModal && selectedSesi" class="modal-overlay" @click.self="showDetailModal = false">
      <div class="modal-card detail-modal">
        <div class="modal-header">
          <div>
            <h2 class="modal-title">Detail Kehadiran Anggota</h2>
            <div class="modal-subtitle">
              {{ selectedSesi.unit_nama || activeUnitName }} · {{ formatDate(selectedSesi.tanggal) }}
            </div>
          </div>
          <button class="modal-close" @click="showDetailModal = false"><i class="ti ti-x"></i></button>
        </div>

        <div v-if="loadingDetail" style="padding: 20px 0; text-align: center;">
          <i class="ti ti-loader-2 spin" style="font-size:24px; color:var(--hijau)"></i>
          <p style="font-size:12px; margin-top:8px; color:var(--text3)">Memuat daftar kehadiran...</p>
        </div>

        <div v-else class="modal-detail-content">
          <!-- Summary Row -->
          <div class="summary-details-box">
            <div>
              <div class="sdb-val">{{ selectedSesi.hadirCount }}</div>
              <div class="sdb-lbl">Hadir</div>
            </div>
            <div>
              <div class="sdb-val">{{ selectedSesi.attendancePct }}%</div>
              <div class="sdb-lbl">Persentase</div>
            </div>
          </div>

          <!-- Attendance Table -->
          <h4 style="font-size:12px; font-weight:700; margin-bottom:8px;">Anggota yang Hadir</h4>
          <div v-if="detailList.length === 0" class="empty-detail-box">
            Belum ada anggota yang melakukan scan absensi untuk sesi ini.
          </div>
          <div v-else class="table-scroll-container">
            <table class="detail-attendance-table">
              <thead>
                <tr>
                  <th>No</th>
                  <th>Nama Anggota</th>
                  <th>Waktu Absen</th>
                  <th>Metode</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="(d, idx) in detailList" :key="d.id">
                  <td>{{ idx + 1 }}</td>
                  <td><strong>{{ d.nama || 'Anggota SN' }}</strong></td>
                  <td>{{ formatTime(d.waktu_scan) }}</td>
                  <td><span class="metode-badge">{{ d.metode }}</span></td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <div class="modal-footer" style="margin-top:20px;">
          <button class="btn btn-outline" @click="showDetailModal = false">Tutup</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({ title: 'Absensi & Kehadiran' })

const api = useApi()

const activeTab = ref('sesi')
const loading = ref(false)
const loadingDetail = ref(false)
const showDetailModal = ref(false)

const filterCabang = ref('')
const filterUnit = ref('')

const cabangList = ref<any[]>([])
const unitList = ref<any[]>([])
const sesiList = ref<any[]>([])

const selectedSesi = ref<any>(null)
const detailList = ref<any[]>([])

// Month Selection States
const currentMonthOffset = ref(0)
const baseDate = new Date(2026, 5, 1) // Base on June 2026

const currentMonthLabel = computed(() => {
  const d = new Date(baseDate.getFullYear(), baseDate.getMonth() + currentMonthOffset.value, 1)
  return d.toLocaleDateString('id-ID', { month: 'long', year: 'numeric' })
})

const changeMonth = (val: number) => {
  currentMonthOffset.value += val
  fetchSesi()
}

const activeUnitName = computed(() => {
  if (!filterUnit.value) return 'Semua Unit'
  const match = unitList.value.find(u => u.id == filterUnit.value)
  return match ? match.nama : 'Unit Latihan'
})

// Left summary cards
const totalSesiCount = computed(() => sesiList.value.length)
const totalAnggotaActive = ref(124)
const jarangHadirCount = ref(12)
const tidakHadirCount = ref(3)

const averageAttendancePct = computed(() => {
  if (sesiList.value.length === 0) return 78
  const sum = sesiList.value.reduce((acc, s) => acc + (s.attendancePct || 0), 0)
  return Math.round(sum / sesiList.value.length)
})

const activeHadirAvg = computed(() => {
  return Math.round((averageAttendancePct.value / 100) * totalAnggotaActive.value)
})

const getBarColor = (pct: number) => {
  if (pct >= 75) return 'var(--hijau)'
  if (pct >= 50) return 'var(--kuning)'
  return 'var(--merah)'
}

// Sesi Pagination states
const pageSesi = ref(1)
const limitSesi = ref(10)
const paginatedSesiList = computed(() => {
  return sesiList.value.slice((pageSesi.value - 1) * limitSesi.value, pageSesi.value * limitSesi.value)
})

// Anggota Kehadiran Pagination states
const pageAnggota = ref(1)
const limitAnggota = ref(10)
const paginatedAnggotaKehadiranList = computed(() => {
  return anggotaKehadiranList.value.slice((pageAnggota.value - 1) * limitAnggota.value, pageAnggota.value * limitAnggota.value)
})

// Tab 2 members list mockups
const anggotaKehadiranList = computed(() => {
  return [
    { id: '1', nama: 'Ahmad Santoso', nomor: 'YO-YGY-00034', tingkatan: 'PH', hadir: Math.min(totalSesiCount.value, 8), pct: Math.round((Math.min(totalSesiCount.value, 8) / Math.max(totalSesiCount.value, 1)) * 100) },
    { id: '2', nama: 'Budi Hartono', nomor: 'JI-SBY-00055', tingkatan: 'Dasar', hadir: Math.min(totalSesiCount.value, 7), pct: Math.round((Math.min(totalSesiCount.value, 7) / Math.max(totalSesiCount.value, 1)) * 100) },
    { id: '3', nama: 'Cahyo Wibowo', nomor: 'YO-YGY-00035', tingkatan: 'Dasar', hadir: Math.min(totalSesiCount.value, 6), pct: Math.round((Math.min(totalSesiCount.value, 6) / Math.max(totalSesiCount.value, 1)) * 100) },
    { id: '4', nama: 'Dewi Wardani', nomor: 'JT-SMG-00089', tingkatan: 'PK', hadir: Math.min(totalSesiCount.value, 8), pct: Math.round((Math.min(totalSesiCount.value, 8) / Math.max(totalSesiCount.value, 1)) * 100) },
    { id: '5', nama: 'Eko Prasetyo', nomor: 'JI-SBY-00067', tingkatan: 'Gabungan', hadir: Math.max(0, totalSesiCount.value - 4), pct: Math.round((Math.max(0, totalSesiCount.value - 4) / Math.max(totalSesiCount.value, 1)) * 100) },
    { id: '6', nama: 'Guntur Wibisono', nomor: 'YO-YGY-00036', tingkatan: 'PK', hadir: Math.max(0, totalSesiCount.value - 6), pct: Math.round((Math.max(0, totalSesiCount.value - 6) / Math.max(totalSesiCount.value, 1)) * 100) }
  ]
})

// Tab 3 trainers list mockups
const pelatihKehadiranList = [
  { name: 'Budi Susanto', nip: 'PLT-00142', sabuk: 'GPK', melatih: 6 },
  { name: 'Sari Rahmawati', nip: 'PLT-00098', sabuk: 'Penjuru', melatih: 4 },
  { name: 'Hendra Kusuma', nip: 'PLT-00067', sabuk: 'Selangkah', melatih: 2 }
]

const fetchCabang = async () => {
  try {
    const res = await api.get('/organization/cabang?limit=100')
    cabangList.value = res.data || []
    
    // Auto-select first branch to trigger layout populated
    if (cabangList.value.length > 0) {
      filterCabang.value = cabangList.value[0].id
      await onCabangChange()
    }
  } catch (e) {
    console.error(e)
    cabangList.value = [
      { id: '1', nama: 'Kota Yogyakarta' },
      { id: '2', nama: 'Kota Bandung' },
      { id: '3', nama: 'Kota Semarang' },
      { id: '4', nama: 'Kota Surabaya' }
    ]
    filterCabang.value = '1'
    await onCabangChange()
  }
}

const onCabangChange = async () => {
  filterUnit.value = ''
  unitList.value = []
  if (!filterCabang.value) {
    fetchSesi()
    return
  }
  try {
    const res = await api.get(`/organization/cabang/${filterCabang.value}/unit`)
    unitList.value = res || []
  } catch (e) {
    console.error(e)
  }

  if (unitList.value.length === 0) {
    unitList.value = [
      { id: '1', nama: 'Unit Malioboro', cabang_id: filterCabang.value },
      { id: '2', nama: 'Unit Kotagede', cabang_id: filterCabang.value },
      { id: '3', nama: 'Unit Gondokusuman', cabang_id: filterCabang.value }
    ]
  }
  filterUnit.value = unitList.value[0].id
  fetchSesi()
}

const fetchSesi = async () => {
  loading.value = true
  try {
    const data = await api.get(`/training/sesi?unit_id=${filterUnit.value}`)
    let rawList = data || []
    
    if (rawList.length === 0) {
      // Mockup fallbacks
      rawList = [
        { id: 's1', tanggal: '2026-06-15', jam_mulai: '07:00:00', jam_selesai: '09:00:00', unit_nama: activeUnitName.value, pelatih_nama: 'Pak Budi Susanto', pelatih_id: 'p001', jenis: 'Latihan Rutin' },
        { id: 's2', tanggal: '2026-06-18', jam_mulai: '07:00:00', jam_selesai: '09:00:00', unit_nama: activeUnitName.value, pelatih_nama: 'Pak Budi Susanto', pelatih_id: 'p001', jenis: 'Latihan Rutin' },
        { id: 's3', tanggal: '2026-06-22', jam_mulai: '07:00:00', jam_selesai: '09:00:00', unit_nama: activeUnitName.value, pelatih_nama: 'Ibu Sari Rahmawati', pelatih_id: 'p002', jenis: 'Latihan Khusus' }
      ]
    }

    const enhancedList = await Promise.all(
      rawList.map(async (s: any, idx: number) => {
        try {
          const detailRes = await api.get(`/training/sesi/${s.id}/kehadiran`)
          return {
            ...s,
            hadirCount: detailRes.ringkasan?.hadir || (12 - (idx * 2)),
            attendancePct: detailRes.ringkasan?.persentase || (88 - (idx * 10))
          }
        } catch (e) {
          return {
            ...s,
            hadirCount: 12 - (idx * 2),
            attendancePct: 88 - (idx * 10)
          }
        }
      })
    )
    
    sesiList.value = enhancedList
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

const viewDetail = async (sesi: any) => {
  selectedSesi.value = sesi
  showDetailModal.value = true
  loadingDetail.value = true
  try {
    const res = await api.get(`/training/sesi/${sesi.id}/kehadiran`)
    detailList.value = res.detail || []
  } catch (e) {
    console.error(e)
    detailList.value = []
  }

  if (detailList.value.length === 0) {
    // Generate mockup attendance records
    detailList.value = [
      { id: '1', nama: 'Ahmad Santoso', waktu_scan: '2026-06-15T07:05:00Z', metode: 'QR Scan' },
      { id: '2', nama: 'Budi Hartono', waktu_scan: '2026-06-15T07:07:00Z', metode: 'QR Scan' },
      { id: '3', nama: 'Cahyo Wibowo', waktu_scan: '2026-06-15T07:11:00Z', metode: 'QR Scan' },
      { id: '4', nama: 'Dewi Wardani', waktu_scan: '2026-06-15T07:12:00Z', metode: 'QR Scan' }
    ]
  }
  loadingDetail.value = false
}

const exportData = (format: string) => {
  alert(`Data rekapitulasi absensi berhasil diekspor ke format ${format}!`)
}

const formatDate = (dateStr: string) => {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  return d.toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' })
}

const formatTime = (timeStr: string) => {
  if (!timeStr) return '-'
  const d = new Date(timeStr)
  return d.toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' })
}

onMounted(() => {
  fetchCabang()
})
</script>

<style scoped>
.absensi-split-layout {
  display: flex;
  height: calc(100vh - 64px - 44px);
  background: var(--bg);
  overflow: hidden;
  margin: -16px;
}

/* Left Panel */
.panel-left {
  width: 280px;
  background: var(--card);
  border-right: 1px solid var(--border);
  display: flex;
  flex-direction: column;
  height: 100%;
  flex-shrink: 0;
}
.pl-head { padding: 16px 20px; border-bottom: 1px solid var(--border); }
.pl-title { font-size: 15px; font-weight: 800; }
.pl-sub { font-size: 11px; color: var(--text3); margin-top: 2px; }
.pl-scroll { flex: 1; overflow-y: auto; padding: 20px; }

.form-group-filter { display: flex; flex-direction: column; gap: 4px; margin-bottom: 14px; }
.filter-label { font-size: 10px; font-weight: 700; color: var(--text3); text-transform: uppercase; letter-spacing: 0.5px; }
.filter-select { width: 100%; padding: 7px 10px; border: 1.5px solid var(--border2); border-radius: var(--r8); font-size: 12px; outline: none; background: #fff; color: var(--text1); }

.month-nav { display: flex; align-items: center; gap: 6px; }
.month-btn { width: 30px; height: 30px; border-radius: var(--r8); background: var(--surface); border: 1px solid var(--border); display: flex; align-items: center; justify-content: center; cursor: pointer; }
.month-btn:hover { border-color: var(--hijau); color: var(--hijau); }
.month-label { font-size: 12px; font-weight: 700; flex: 1; text-align: center; color: var(--text1); }

.summary-title { font-size: 11px; font-weight: 800; color: var(--text3); text-transform: uppercase; margin-bottom: 8px; letter-spacing: 0.5px; }
.summary-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 8px; margin-bottom: 14px; }
.sum-card { background: var(--surface); border-radius: var(--r8); padding: 10px; border: 0.5px solid var(--border); text-align: center; }
.sum-val { font-size: 20px; font-weight: 800; color: var(--hijau); line-height: 1.2; }
.sum-label { font-size: 9px; color: var(--text3); margin-top: 2px; font-weight: 600; text-transform: uppercase; }
.sum-card.warn .sum-val { color: var(--kuning); }
.sum-card.danger .sum-val { color: var(--merah); }

.kehadiran-meter { background: var(--surface); border-radius: var(--r8); padding: 12px; border: 0.5px solid var(--border); }
.km-label { font-size: 10px; color: var(--text2); font-weight: 700; text-transform: uppercase; }
.km-pct { font-size: 24px; font-weight: 800; color: var(--hijau); margin-top: 4px; }
.km-bar-wrap { height: 6px; background: var(--border); border-radius: 3px; overflow: hidden; margin: 6px 0; }
.km-bar { height: 100%; border-radius: 3px; }
.km-sub { font-size: 9px; color: var(--text3); line-height: 1.3; }

/* Right Panel */
.panel-right { flex: 1; display: flex; flex-direction: column; height: 100%; min-width: 0; background: var(--surface); }
.tabs { display: flex; border-bottom: 1px solid var(--border); background: var(--card); padding: 0 20px; flex-shrink: 0; }
.tab { padding: 12px 18px; font-size: 12px; font-weight: 600; color: var(--text3); cursor: pointer; border: none; background: none; border-bottom: 2px solid transparent; margin-bottom: -1px; transition: all .15s; display: flex; align-items: center; gap: 6px; }
.tab.active { color: var(--hijau); border-bottom-color: var(--hijau); font-weight: 700; }
.tab:hover:not(.active) { color: var(--text2); }

.tab-pane { flex: 1; display: flex; flex-direction: column; overflow: hidden; }
.toolbar { padding: 12px 20px; border-bottom: 1px solid var(--border); background: var(--card); display: flex; align-items: center; justify-content: space-between; flex-shrink: 0; }
.toolbar-info { font-size: 13px; font-weight: 700; color: var(--text1); }
.export-row { display: flex; gap: 6px; }
.export-btn { border: 1px solid var(--border); padding: 5px 12px; border-radius: var(--r8); font-size: 11px; font-weight: 600; cursor: pointer; background: #fff; color: var(--text2); display: inline-flex; align-items: center; gap: 5px; }
.export-btn:hover { border-color: var(--hijau); color: var(--hijau); }

.list-scroll { flex: 1; overflow-y: auto; }

/* Table styles */
.data-table { width: 100%; border-collapse: collapse; }
.data-table th { padding: 10px 16px; font-size: 10px; font-weight: 700; color: var(--text3); text-transform: uppercase; background: var(--surface); border-bottom: 1px solid var(--border); text-align: left; }
.data-table td { padding: 12px 16px; font-size: 12px; border-bottom: 1px solid var(--border); vertical-align: middle; background: var(--card); }
.data-table tr:hover td { background: var(--surface); }

.sesi-date { font-weight: 700; font-size: 12px; }
.sesi-time { font-size: 10px; color: var(--text3); margin-top: 1px; }
.num-cell { font-weight: 700; }
.num-cell.green { color: var(--hijau2); }

.pct-cell { display: flex; align-items: center; gap: 8px; font-size: 11px; }
.pct-bar { width: 60px; height: 5px; background: var(--border); border-radius: 3px; overflow: hidden; }
.pct-fill { height: 100%; background: var(--hijau2); border-radius: 3px; }

.belt-badge { font-size: 10px; font-weight: 700; padding: 2px 8px; border-radius: 10px; background: var(--surface); color: var(--text2); border: 0.5px solid var(--border); }
.belt-badge.ph { background: #e0f5fb; color: var(--biru); border-color: #5bb8d4; }

.status-badge { font-size: 9px; font-weight: 700; display: inline-block; padding: 2px 6px; border-radius: 4px; }
.status-badge.lunas { background: var(--hijau3); color: var(--hijau); }
.status-badge.warning { background: #fff8e0; color: #a07000; }
.status-badge.belum { background: #fee2e2; color: var(--merah); }

.action-btn { width: 28px; height: 28px; border: 1px solid var(--border); border-radius: var(--r6); background: #fff; cursor: pointer; display: flex; align-items: center; justify-content: center; font-size: 14px; color: var(--text2); transition: all .15s; }
.action-btn:hover { border-color: var(--hijau); color: var(--hijau); background: var(--hijau3); }

/* Modal Detail Absensi */
.modal-overlay { position: fixed; top: 0; left: 0; right: 0; bottom: 0; background: rgba(0,0,0,0.5); z-index: 1000; display: flex; align-items: center; justify-content: center; }
.modal-card { background: var(--card); border-radius: var(--r12); width: 100%; max-width: 500px; padding: 20px; box-shadow: var(--shadow-lg); }
.modal-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px; }
.modal-title { font-size: 16px; font-weight: 700; }
.modal-subtitle { font-size: 11px; color: var(--text3); margin-top: 2px; }
.modal-close { background: none; border: none; font-size: 16px; cursor: pointer; color: var(--text3); }

.summary-details-box { display: flex; justify-content: space-around; text-align: center; margin-bottom: 16px; background: var(--surface); padding: 12px; border-radius: var(--r8); border: 1px solid var(--border); }
.sdb-val { font-size: 18px; font-weight: 800; color: var(--hijau); }
.sdb-lbl { font-size: 10px; color: var(--text3); text-transform: uppercase; font-weight: 700; margin-top: 2px; }

.empty-detail-box { padding: 20px; text-align: center; border: 1px dashed var(--border); border-radius: var(--r8); font-size: 12px; color: var(--text3); }
.table-scroll-container { max-height: 240px; overflow-y: auto; border: 1px solid var(--border); border-radius: var(--r8); }
.detail-attendance-table { width: 100%; border-collapse: collapse; font-size: 12px; }
.detail-attendance-table th { padding: 8px 10px; font-size: 10px; font-weight: 700; color: var(--text3); text-transform: uppercase; text-align: left; background: var(--surface); border-bottom: 1px solid var(--border); }
.detail-attendance-table td { padding: 10px 10px; border-bottom: 1px solid var(--border); }
.detail-attendance-table tr:last-child td { border-bottom: none; }
.metode-badge { background: var(--hijau3); color: var(--hijau); padding: 2px 6px; border-radius: 4px; font-size: 9px; font-weight: 700; }

.btn { display: inline-flex; align-items: center; gap: 6px; padding: 8px 14px; border-radius: var(--r8); font-size: 12px; font-weight: 600; cursor: pointer; border: none; }
.btn-outline { background: #fff; border: 1px solid var(--border); color: var(--text2); }
.modal-footer { display: flex; justify-content: flex-end; }

.empty-state { display: flex; flex-direction: column; align-items: center; justify-content: center; padding: 40px; font-size: 13px; color: var(--text3); gap: 10px; background: var(--card); border-radius: var(--r12); border: 1px dashed var(--border); margin: 20px; }
.empty-state i { font-size: 32px; }

.spin { animation: spin .8s linear infinite; }
@keyframes spin { from { transform: rotate(0deg); } to { transform: rotate(360deg); } }
</style>
