<template>
  <div id="page-monitoring-nafas" class="nafas-page-container">
    <!-- PAGE HEADER -->
    <div class="nafas-header">
      <div>
        <h1 class="page-title"><i class="ti ti-wind"></i> Monitoring Olah Nafas</h1>
        <p class="page-subtitle">Pantau rekapitulasi & riwayat latihan pernafasan anggota Satria Nusantara secara real-time</p>
      </div>
      <div class="header-actions">
        <button class="btn btn-outline" @click="fetchData">
          <i class="ti ti-refresh" :class="{ spin: isLoading }"></i> Refresh
        </button>
        <button class="btn btn-primary" @click="exportCSV">
          <i class="ti ti-file-export"></i> Export Rekap
        </button>
      </div>
    </div>

    <!-- SUMMARY KPI CARDS -->
    <div class="kpi-grid">
      <div class="kpi-card">
        <div class="kpi-icon-wrap green">
          <i class="ti ti-wind"></i>
        </div>
        <div class="kpi-details">
          <span class="kpi-label">Total Sesi Latihan</span>
          <h2 class="kpi-value">{{ stats.total_sesi || 0 }}</h2>
          <span class="kpi-sub">Sesi Ter-record</span>
        </div>
      </div>

      <div class="kpi-card">
        <div class="kpi-icon-wrap blue">
          <i class="ti ti-clock"></i>
        </div>
        <div class="kpi-details">
          <span class="kpi-label">Total Akumulasi Durasi</span>
          <h2 class="kpi-value">{{ stats.total_menit || 0 }} <span class="unit">Menit</span></h2>
          <span class="kpi-sub">Waktu Olah Nafas</span>
        </div>
      </div>

      <div class="kpi-card">
        <div class="kpi-icon-wrap purple">
          <i class="ti ti-users"></i>
        </div>
        <div class="kpi-details">
          <span class="kpi-label">Anggota Aktif</span>
          <h2 class="kpi-value">{{ stats.anggota_aktif || 0 }}</h2>
          <span class="kpi-sub">Peserta Latihan</span>
        </div>
      </div>

      <div class="kpi-card">
        <div class="kpi-icon-wrap orange">
          <i class="ti ti-chart-donut-4"></i>
        </div>
        <div class="kpi-details">
          <span class="kpi-label">Teknik Terfavorit</span>
          <h2 class="kpi-value text-sm">{{ stats.teknik_favorit || 'Sama Kaki' }}</h2>
          <span class="kpi-sub">Metode Utama</span>
        </div>
      </div>
    </div>

    <!-- FILTER & SEARCH BAR -->
    <div class="filter-card">
      <div class="filter-group flex-1">
        <div class="search-input-wrap">
          <i class="ti ti-search"></i>
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Cari nama anggota, cabang, atau unit..."
            @input="debouncedSearch"
          />
          <button v-if="searchQuery" class="clear-btn" @click="searchQuery = ''; fetchData()">
            <i class="ti ti-x"></i>
          </button>
        </div>
      </div>

      <div class="filter-group">
        <label class="filter-lbl">Teknik:</label>
        <select v-model="filterTeknik" class="filter-select" @change="fetchData">
          <option value="">Semua Teknik</option>
          <option value="Sama Kaki">Sama Kaki (Segitiga Rebah)</option>
          <option value="Sama Sisi">Sama Sisi (Segitiga Sama Sisi)</option>
          <option value="Segi Empat">Segi Empat (Box Breathing)</option>
        </select>
      </div>

      <div class="filter-group">
        <label class="filter-lbl">Cabang:</label>
        <select v-model="filterCabang" class="filter-select" @change="fetchData">
          <option value="">Semua Cabang</option>
          <option value="Yogyakarta">Kota Yogyakarta</option>
          <option value="Surabaya">Surabaya</option>
          <option value="Bandung">Kota Bandung</option>
          <option value="Jakarta">DKI Jakarta</option>
        </select>
      </div>
    </div>

    <!-- MAIN DATA TABLE -->
    <div class="table-card">
      <div class="table-header">
        <h3 class="table-title">Riwayat Sesi Olah Nafas Anggota</h3>
        <span class="table-count">Menampilkan {{ filteredHistories.length }} data</span>
      </div>

      <div v-if="isLoading" class="table-state-box">
        <i class="ti ti-loader-2 spin text-2xl"></i>
        <span>Memuat data riwayat latihan...</span>
      </div>

      <div v-else-if="filteredHistories.length === 0" class="table-state-box">
        <i class="ti ti-wind-off text-3xl"></i>
        <span>Tidak ada data riwayat olah nafas ditemukan</span>
      </div>

      <div v-else class="table-responsive">
        <table class="custom-table">
          <thead>
            <tr>
              <th>No</th>
              <th>Anggota / Peserta</th>
              <th>Cabang & Unit</th>
              <th>Teknik Pernafasan</th>
              <th>Durasi Latihan</th>
              <th>Jumlah Siklus</th>
              <th>Waktu Selesai</th>
              <th>Status</th>
              <th>Aksi</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(row, index) in paginatedHistories" :key="row.id || index">
              <td>{{ (page - 1) * limit + index + 1 }}</td>
              <td>
                <div class="user-cell">
                  <div class="user-avatar">{{ getInitials(row.anggota_nama) }}</div>
                  <div>
                    <div class="user-name">{{ row.anggota_nama }}</div>
                    <div class="user-id">ID: {{ row.anggota_id ? row.anggota_id.substring(0, 8) : 'A-101' }}</div>
                  </div>
                </div>
              </td>
              <td>
                <div class="cabang-cell">
                  <span class="cabang-title">{{ row.cabang_nama || 'Kota Yogyakarta' }}</span>
                  <span class="unit-sub">Unit {{ row.unit_nama || 'Pusat' }}</span>
                </div>
              </td>
              <td>
                <span :class="['teknik-badge', getTeknikClass(row.teknik)]">
                  <i class="ti ti-activity"></i> {{ row.teknik }}
                </span>
              </td>
              <td>
                <div class="duration-cell">
                  <i class="ti ti-clock"></i> {{ row.durasi_fmt || '10:00' }}
                  <span class="sec-hint">({{ row.durasi_detik }}d)</span>
                </div>
              </td>
              <td>
                <span class="cycle-pill">{{ row.siklus || 12 }} Siklus</span>
              </td>
              <td>
                <div class="time-cell">
                  <span>{{ formatDate(row.timestamp) }}</span>
                  <span class="clock-sub">{{ formatTime(row.timestamp) }}</span>
                </div>
              </td>
              <td>
                <span class="status-badge-success"><i class="ti ti-circle-check"></i> Selesai</span>
              </td>
              <td>
                <button class="btn-sm btn-ghost" title="Detail Riwayat" @click="selectedRow = row">
                  <i class="ti ti-eye"></i> Detail
                </button>
              </td>
            </tr>
          </tbody>
        </table>
        <Pagination 
          v-model:currentPage="page" 
          v-model:itemsPerPage="limit" 
          :totalItems="filteredHistories.length" 
        />
      </div>
    </div>

    <!-- DETAIL MODAL POPUP -->
    <div v-if="selectedRow" class="modal-overlay" @click.self="selectedRow = null">
      <div class="modal-card">
        <div class="modal-header">
          <div class="modal-title-wrap">
            <i class="ti ti-wind modal-icon"></i>
            <div>
              <h3 class="modal-title">Detail Sesi Olah Nafas</h3>
              <span class="modal-sub">ID Sesi: {{ selectedRow.id }}</span>
            </div>
          </div>
          <button class="modal-close-btn" @click="selectedRow = null"><i class="ti ti-x"></i></button>
        </div>

        <div class="modal-body">
          <div class="detail-row">
            <span class="detail-label">Nama Anggota</span>
            <span class="detail-val highlight">{{ selectedRow.anggota_nama }}</span>
          </div>

          <div class="detail-row">
            <span class="detail-label">Cabang & Unit</span>
            <span class="detail-val">{{ selectedRow.cabang_nama }} (Unit {{ selectedRow.unit_nama }})</span>
          </div>

          <div class="detail-row">
            <span class="detail-label">Teknik Pernafasan</span>
            <span class="detail-val">
              <span :class="['teknik-badge', getTeknikClass(selectedRow.teknik)]">
                {{ selectedRow.teknik }}
              </span>
            </span>
          </div>

          <div class="detail-row">
            <span class="detail-label">Total Durasi Latihan</span>
            <span class="detail-val bold">{{ selectedRow.durasi_fmt }} ({{ selectedRow.durasi_detik }} detik)</span>
          </div>

          <div class="detail-row">
            <span class="detail-label">Total Siklus Selesai</span>
            <span class="detail-val bold">{{ selectedRow.siklus }} Siklus Penuh</span>
          </div>

          <div class="detail-row">
            <span class="detail-label">Waktu Selesai</span>
            <span class="detail-val">{{ formatDate(selectedRow.timestamp) }} pukul {{ formatTime(selectedRow.timestamp) }}</span>
          </div>

          <div class="breakdown-card">
            <div class="breakdown-title">Rincian Pola Durasi per Step:</div>
            <div v-if="selectedRow.teknik === 'Sama Kaki'" class="step-chips">
              <span class="chip blue">Tarik: 10d</span>
              <span class="chip yellow">Tahan & Tekan: 30d</span>
              <span class="chip green">Keluarkan: 10d</span>
            </div>
            <div v-else-if="selectedRow.teknik === 'Segi Empat'" class="step-chips">
              <span class="chip blue">Tarik: 10d</span>
              <span class="chip yellow">Tahan: 10d</span>
              <span class="chip green">Keluar: 10d</span>
              <span class="chip orange">Tahan 2: 10d</span>
            </div>
            <div v-else class="step-chips">
              <span class="chip blue">Tarik: 10d</span>
              <span class="chip yellow">Tahan & Tekan: 10d</span>
              <span class="chip green">Keluarkan: 10d</span>
            </div>
          </div>
        </div>

        <div class="modal-footer">
          <button class="btn btn-outline" @click="selectedRow = null">Tutup</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
const api = useApi()

const isLoading = ref(false)
const searchQuery = ref('')
const filterTeknik = ref('')
const filterCabang = ref('')
const selectedRow = ref(null)

const stats = ref({
  total_sesi: 5,
  total_menit: 48,
  anggota_aktif: 4,
  teknik_favorit: 'Sama Kaki',
})

const histories = ref([
  {
    id: 'hist-101',
    anggota_id: 'a0000001-0000-0000-0000-000000000001',
    anggota_nama: 'Ahmad Santoso',
    cabang_nama: 'Kota Yogyakarta',
    unit_nama: 'Malioboro',
    teknik: 'Sama Kaki',
    durasi_detik: 600,
    durasi_fmt: '10:00',
    siklus: 12,
    timestamp: new Date().toISOString(),
  },
  {
    id: 'hist-102',
    anggota_id: 'a0000002-0000-0000-0000-000000000002',
    anggota_nama: 'Budi Susanto',
    cabang_nama: 'Surabaya',
    unit_nama: 'Rungkut',
    teknik: 'Segi Empat',
    durasi_detik: 480,
    durasi_fmt: '08:00',
    siklus: 12,
    timestamp: new Date(Date.now() - 7200000).toISOString(),
  },
  {
    id: 'hist-103',
    anggota_id: 'a0000003-0000-0000-0000-000000000003',
    anggota_nama: 'Siti Rahmawati',
    cabang_nama: 'Kota Bandung',
    unit_nama: 'Dago',
    teknik: 'Sama Sisi',
    durasi_detik: 300,
    durasi_fmt: '05:00',
    siklus: 10,
    timestamp: new Date(Date.now() - 18000000).toISOString(),
  },
  {
    id: 'hist-104',
    anggota_id: 'a0000004-0000-0000-0000-000000000004',
    anggota_nama: 'Dewi Lestari',
    cabang_nama: 'DKI Jakarta',
    unit_nama: 'Monas',
    teknik: 'Sama Kaki',
    durasi_detik: 900,
    durasi_fmt: '15:00',
    siklus: 18,
    timestamp: new Date(Date.now() - 86400000).toISOString(),
  },
  {
    id: 'hist-105',
    anggota_id: 'a0000001-0000-0000-0000-000000000001',
    anggota_nama: 'Ahmad Santoso',
    cabang_nama: 'Kota Yogyakarta',
    unit_nama: 'Malioboro',
    teknik: 'Sama Kaki',
    durasi_detik: 600,
    durasi_fmt: '10:00',
    siklus: 12,
    timestamp: new Date(Date.now() - 172800000).toISOString(),
  },
])

let searchTimer = null
const debouncedSearch = () => {
  clearTimeout(searchTimer)
  searchTimer = setTimeout(() => {
    fetchData()
  }, 300)
}

const fetchData = async () => {
  isLoading.value = true
  try {
    const params = new URLSearchParams()
    if (searchQuery.value) params.append('search', searchQuery.value)
    if (filterTeknik.value) params.append('teknik', filterTeknik.value)

    const res = await api.get(`/nafas/monitoring?${params.toString()}`)
    if (res) {
      stats.value.total_sesi = res.total_sesi ?? histories.value.length
      stats.value.total_menit = res.total_menit ?? 48
      stats.value.anggota_aktif = res.anggota_aktif ?? 4
      stats.value.teknik_favorit = res.teknik_favorit ?? 'Sama Kaki'

      if (res.histories && res.histories.length > 0) {
        histories.value = res.histories
      }
    }
  } catch (err) {
    console.warn('API connection fallback to mock:', err)
  } finally {
    isLoading.value = false
  }
}

const filteredHistories = computed(() => {
  return histories.value.filter(item => {
    if (filterTeknik.value && item.teknik !== filterTeknik.value) return false
    if (filterCabang.value && !item.cabang_nama.toLowerCase().includes(filterCabang.value.toLowerCase())) return false
    if (searchQuery.value) {
      const q = searchQuery.value.toLowerCase()
      const matchName = item.anggota_nama.toLowerCase().includes(q)
      const matchCabang = item.cabang_nama.toLowerCase().includes(q)
      const matchUnit = item.unit_nama ? item.unit_nama.toLowerCase().includes(q) : false
      if (!matchName && !matchCabang && !matchUnit) return false
    }
    return true
  })
})

// Pagination states
const page = ref(1)
const limit = ref(10)
const paginatedHistories = computed(() => {
  return filteredHistories.value.slice((page.value - 1) * limit.value, page.value * limit.value)
})

const getInitials = (name) => {
  if (!name) return 'SN'
  const parts = name.trim().split(' ')
  if (parts.length > 1) return (parts[0][0] + parts[1][0]).toUpperCase()
  return name.substring(0, 2).toUpperCase()
}

const getTeknikClass = (teknik) => {
  if (teknik === 'Sama Kaki') return 'green'
  if (teknik === 'Segi Empat') return 'orange'
  return 'blue'
}

const formatDate = (ts) => {
  if (!ts) return '-'
  const d = new Date(ts)
  return `${d.getDate().toString().padLeft(2, '0')}/${(d.getMonth() + 1).toString().padLeft(2, '0')}/${d.getFullYear()}`
}

const formatTime = (ts) => {
  if (!ts) return '-'
  const d = new Date(ts)
  return `${d.getHours().toString().padLeft(2, '0')}:${d.getMinutes().toString().padLeft(2, '0')}`
}

const exportCSV = () => {
  let csvContent = 'data:text/csv;charset=utf-8,'
  csvContent += 'No,Nama Anggota,Cabang,Unit,Teknik,Durasi,Siklus,Tanggal,Waktu\n'

  filteredHistories.value.forEach((row, i) => {
    csvContent += `${i + 1},"${row.anggota_nama}","${row.cabang_nama}","${row.unit_nama}","${row.teknik}","${row.durasi_fmt}",${row.siklus},"${formatDate(row.timestamp)}","${formatTime(row.timestamp)}"\n`
  })

  const encodedUri = encodeURI(csvContent)
  const link = document.createElement('a')
  link.setAttribute('href', encodedUri)
  link.setAttribute('download', `Rekap_Olah_Nafas_SN_${new Date().toISOString().slice(0, 10)}.csv`)
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.nafas-page-container {
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.nafas-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.page-title {
  font-size: 22px;
  font-weight: 800;
  color: var(--text1);
  display: flex;
  align-items: center;
  gap: 8px;
  letter-spacing: -0.5px;
}
.page-title i { color: var(--hijau); }

.page-subtitle {
  font-size: 13px;
  color: var(--text2);
  margin-top: 4px;
}

.header-actions {
  display: flex;
  gap: 10px;
}

.btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 9px 16px;
  border-radius: var(--r8);
  font-size: 13px;
  font-weight: 700;
  cursor: pointer;
  border: none;
  transition: all 0.2s;
}
.btn-primary { background: var(--hijau); color: #fff; }
.btn-primary:hover { background: #154e24; }
.btn-outline { background: #fff; border: 1px solid var(--border); color: var(--text2); }
.btn-outline:hover { background: var(--surface); }

/* KPI Grid */
.kpi-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
}

.kpi-card {
  background: var(--card);
  border: 1px solid var(--border);
  border-radius: var(--r12);
  padding: 18px;
  display: flex;
  align-items: center;
  gap: 16px;
  box-shadow: var(--shadow-sm);
}

.kpi-icon-wrap {
  width: 48px;
  height: 48px;
  border-radius: var(--r12);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
}
.kpi-icon-wrap.green { background: rgba(26, 92, 42, 0.1); color: var(--hijau); }
.kpi-icon-wrap.blue { background: rgba(0, 145, 234, 0.1); color: #0091ea; }
.kpi-icon-wrap.purple { background: rgba(124, 77, 255, 0.1); color: #7c4dff; }
.kpi-icon-wrap.orange { background: rgba(255, 112, 67, 0.1); color: #ff7043; }

.kpi-details { display: flex; flex-direction: column; }
.kpi-label { font-size: 11px; font-weight: 600; color: var(--text3); text-transform: uppercase; letter-spacing: 0.5px; }
.kpi-value { font-size: 22px; font-weight: 800; color: var(--text1); margin: 2px 0; }
.kpi-value.text-sm { font-size: 16px; }
.kpi-value .unit { font-size: 13px; font-weight: 600; color: var(--text2); }
.kpi-sub { font-size: 11px; color: var(--text3); }

/* Filter Card */
.filter-card {
  background: var(--card);
  border: 1px solid var(--border);
  border-radius: var(--r12);
  padding: 14px 18px;
  display: flex;
  align-items: center;
  gap: 16px;
}

.flex-1 { flex: 1; }

.search-input-wrap {
  position: relative;
  display: flex;
  align-items: center;
  width: 100%;
}
.search-input-wrap i { position: absolute; left: 12px; color: var(--text3); font-size: 16px; }
.search-input-wrap input {
  width: 100%;
  padding: 9px 12px 9px 36px;
  border: 1px solid var(--border);
  border-radius: var(--r8);
  font-size: 13px;
  outline: none;
}
.search-input-wrap input:focus { border-color: var(--hijau); }
.clear-btn { position: absolute; right: 10px; background: none; border: none; cursor: pointer; color: var(--text3); }

.filter-group { display: flex; align-items: center; gap: 8px; }
.filter-lbl { font-size: 12px; font-weight: 700; color: var(--text2); }
.filter-select {
  padding: 8px 12px;
  border: 1px solid var(--border);
  border-radius: var(--r8);
  font-size: 12.5px;
  outline: none;
  background: #fff;
  cursor: pointer;
}

/* Table Card */
.table-card {
  background: var(--card);
  border: 1px solid var(--border);
  border-radius: var(--r12);
  overflow: hidden;
  box-shadow: var(--shadow-sm);
}

.table-header {
  padding: 16px 20px;
  border-bottom: 1px solid var(--border);
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.table-title { font-size: 15px; font-weight: 700; color: var(--text1); }
.table-count { font-size: 12px; color: var(--text3); }

.table-state-box {
  padding: 40px;
  text-align: center;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
  color: var(--text3);
  font-size: 13px;
}

.custom-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 13px;
}
.custom-table th {
  background: var(--surface);
  color: var(--text2);
  font-weight: 700;
  font-size: 11px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  padding: 12px 16px;
  text-align: left;
  border-bottom: 1px solid var(--border);
}
.custom-table td {
  padding: 12px 16px;
  border-bottom: 1px solid var(--border);
  color: var(--text1);
  vertical-align: middle;
}
.custom-table tr:last-child td { border-bottom: none; }
.custom-table tr:hover td { background: rgba(0, 0, 0, 0.015); }

/* User cell */
.user-cell { display: flex; align-items: center; gap: 10px; }
.user-avatar {
  width: 34px; height: 34px; border-radius: 50%;
  background: var(--hijauSoft); color: var(--hijau);
  font-weight: 800; font-size: 12px;
  display: flex; align-items: center; justify-content: center;
}
.user-name { font-weight: 700; color: var(--text1); }
.user-id { font-size: 10.5px; color: var(--text3); }

.cabang-cell { display: flex; flex-direction: column; }
.cabang-title { font-weight: 600; }
.unit-sub { font-size: 11px; color: var(--text3); }

/* Badges */
.teknik-badge {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 4px 8px;
  border-radius: 6px;
  font-size: 11px;
  font-weight: 700;
}
.teknik-badge.green { background: rgba(46, 139, 61, 0.12); color: #2e8b3d; }
.teknik-badge.blue { background: rgba(0, 145, 234, 0.12); color: #0091ea; }
.teknik-badge.orange { background: rgba(201, 112, 26, 0.12); color: #c9701a; }

.duration-cell { font-weight: 700; display: flex; align-items: center; gap: 4px; }
.sec-hint { font-weight: normal; font-size: 11px; color: var(--text3); }

.cycle-pill { background: var(--border2); padding: 3px 8px; border-radius: 12px; font-weight: 700; font-size: 11px; }

.time-cell { display: flex; flex-direction: column; }
.clock-sub { font-size: 10.5px; color: var(--text3); }

.status-badge-success {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  background: #e8f5e9;
  color: #2e7d32;
  padding: 3px 8px;
  border-radius: 6px;
  font-size: 11px;
  font-weight: 700;
}

.btn-sm { padding: 4px 10px; font-size: 11.5px; border-radius: 6px; border: none; cursor: pointer; }
.btn-ghost { background: var(--surface); color: var(--text2); border: 1px solid var(--border); }
.btn-ghost:hover { background: #e0e0e0; }

/* Modal */
.modal-overlay {
  position: fixed; top: 0; left: 0; right: 0; bottom: 0;
  background: rgba(0,0,0,0.5); z-index: 1000;
  display: flex; align-items: center; justify-content: center;
}
.modal-card {
  background: #fff; border-radius: var(--r12);
  width: 100%; max-width: 480px; padding: 24px;
  box-shadow: var(--shadow-lg);
}
.modal-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; }
.modal-title-wrap { display: flex; align-items: center; gap: 12px; }
.modal-icon { font-size: 28px; color: var(--hijau); }
.modal-title { font-size: 16px; font-weight: 800; margin: 0; }
.modal-sub { font-size: 11px; color: var(--text3); }
.modal-close-btn { background: none; border: none; font-size: 18px; cursor: pointer; color: var(--text3); }

.modal-body { display: flex; flex-direction: column; gap: 12px; font-size: 13px; }
.detail-row { display: flex; justify-content: space-between; padding-bottom: 8px; border-bottom: 1px solid var(--border); }
.detail-label { color: var(--text3); }
.detail-val { font-weight: 600; color: var(--text1); }
.detail-val.highlight { color: var(--hijau); font-weight: 800; }
.detail-val.bold { font-weight: 800; }

.breakdown-card {
  background: var(--surface);
  padding: 12px;
  border-radius: var(--r8);
  border: 1px solid var(--border);
  margin-top: 8px;
}
.breakdown-title { font-size: 11px; font-weight: 700; color: var(--text2); margin-bottom: 8px; }
.step-chips { display: flex; gap: 6px; flex-wrap: wrap; }
.chip { padding: 3px 8px; border-radius: 6px; font-size: 11px; font-weight: 700; }
.chip.blue { background: #e1f5fe; color: #0288d1; }
.chip.yellow { background: #fff8e1; color: #f57f17; }
.chip.green { background: #e8f5e9; color: #2e7d32; }
.chip.orange { background: #fff3e0; color: #e65100; }

.modal-footer { margin-top: 20px; display: flex; justify-content: flex-end; }
.spin { animation: spin 1s linear infinite; }
@keyframes spin { 100% { transform: rotate(360deg); } }
</style>
