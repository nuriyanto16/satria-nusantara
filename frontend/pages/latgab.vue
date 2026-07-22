<template>
  <div id="page-latgab" class="latgab-split-layout">
    <!-- ===== PANEL KIRI: FORM ===== -->
    <div class="panel-form">
      <div class="pf-head">
        <div class="pf-title">{{ editingEventId ? 'Edit Event' : 'Buat Event Baru' }}</div>
        <div class="pf-sub">Isi detail event Latgab, EKT, atau Pelatnas</div>
      </div>
      <div class="pf-scroll">
        <!-- JENIS EVENT -->
        <div class="form-group">
          <label class="form-label">Jenis Event</label>
          <div class="jenis-toggle">
            <div :class="['jenis-btn', { active: formJenis === 'Latgab' }]" @click="setJenis('Latgab')">
              <i class="ti ti-users-group"></i>
              <span class="jenis-label">Latgab</span>
            </div>
            <div :class="['jenis-btn', { active: formJenis === 'EKT Jurus' }]" @click="setJenis('EKT Jurus')">
              <i class="ti ti-award"></i>
              <span class="jenis-label">EKT Jurus</span>
            </div>
            <div :class="['jenis-btn', { active: formJenis === 'EKT Non Jurus' }]" @click="setJenis('EKT Non Jurus')">
              <i class="ti ti-certificate"></i>
              <span class="jenis-label">EKT NJ</span>
            </div>
            <div :class="['jenis-btn', { active: formJenis === 'Pelatnas' }]" @click="setJenis('Pelatnas')">
              <i class="ti ti-star"></i>
              <span class="jenis-label">Pelatnas</span>
            </div>
          </div>
        </div>

        <!-- PENYELENGGARA & PESERTA -->
        <div class="form-section">
          <div class="form-section-title"><i class="ti ti-map-pin"></i> Penyelenggara & Peserta</div>

          <!-- Cabang Penyelenggara (All except Pelatnas) -->
          <div v-if="formJenis !== 'Pelatnas'" class="form-group">
            <label class="form-label">Cabang Penyelenggara</label>
            <select v-model="formCabang" class="form-input">
              <option value="Kota Yogyakarta">Kota Yogyakarta</option>
              <option value="Kota Bandung">Kota Bandung</option>
              <option value="Jakarta Pusat">Jakarta Pusat</option>
              <option value="Kota Surabaya">Kota Surabaya</option>
              <option value="Kota Semarang">Kota Semarang</option>
            </select>
          </div>

          <!-- Peserta (Latgab & EKT) -->
          <div v-if="formJenis !== 'Pelatnas'" class="form-group">
            <label class="form-label">Peserta</label>
            <div class="peserta-tabs">
              <div :class="['peserta-tab', { active: formPesertaMode === 'provinsi' }]" @click="formPesertaMode = 'provinsi'">
                <i class="ti ti-map"></i> Per Provinsi
              </div>
              <div :class="['peserta-tab', { active: formPesertaMode === 'multi' }]" @click="formPesertaMode = 'multi'">
                <i class="ti ti-building"></i> Pilih Cabang
              </div>
            </div>

            <!-- Mode Provinsi -->
            <div v-if="formPesertaMode === 'provinsi'" style="margin-top: 8px;">
              <select v-model="formProvinsi" class="form-input">
                <option value="DI Yogyakarta">DI Yogyakarta</option>
                <option value="Jawa Barat">Jawa Barat</option>
                <option value="DKI Jakarta">DKI Jakarta</option>
                <option value="Jawa Timur">Jawa Timur</option>
                <option value="Jawa Tengah">Jawa Tengah</option>
              </select>
              <div class="form-hint">Semua cabang di provinsi terpilih otomatis terdaftar.</div>
            </div>

            <!-- Mode Multi-Cabang -->
            <div v-else style="margin-top: 8px;">
              <div class="cabang-checkboxes">
                <label v-for="c in ['Kota Yogyakarta', 'Kota Bandung', 'Jakarta Pusat', 'Kota Surabaya', 'Kota Semarang']" :key="c" class="cb-label">
                  <input type="checkbox" :value="c" v-model="formSelectedCabangs" /> {{ c }}
                </label>
              </div>
            </div>
          </div>

          <!-- Pelatnas Scope -->
          <div v-else class="peny-nasional">
            <i class="ti ti-flag"></i>
            <div>
              <div class="pn-title">Berlaku Nasional · SN Pusat</div>
              <div class="pn-sub">Pelatnas diselenggarakan langsung oleh Pusat untuk seluruh Indonesia.</div>
            </div>
          </div>

          <!-- Penerima Notifikasi Otomatis -->
          <div class="form-group" style="margin-top: 12px;">
            <label class="form-label">Penerima Notifikasi Otomatis</label>
            <div class="notif-scope-card">
              <i class="ti ti-bell"></i>
              <div>
                <div class="ns-title">{{ notifScopeTitle }}</div>
                <div class="ns-sub">{{ notifScopeSub }}</div>
              </div>
            </div>
          </div>
        </div>

        <!-- INFO EVENT -->
        <div class="form-section">
          <div class="form-section-title"><i class="ti ti-info-circle"></i> Info Event</div>
          <div class="form-group">
            <label class="form-label">Nama Event</label>
            <input v-model="formNama" class="form-input" placeholder="Cth: Latgab DIY Jul 2026" />
          </div>
          <div class="form-group">
            <label class="form-label">Tanggal</label>
            <input v-model="formTanggal" class="form-input" type="date" />
          </div>
          <div class="form-row-2">
            <div class="form-group">
              <label class="form-label">Jam Mulai</label>
              <input v-model="formJamMulai" class="form-input" type="time" />
            </div>
            <div class="form-group">
              <label class="form-label">Jam Selesai</label>
              <input v-model="formJamSelesai" class="form-input" type="time" />
            </div>
          </div>
          <div class="form-group">
            <label class="form-label">Lokasi</label>
            <input v-model="formLokasi" class="form-input" placeholder="Cth: GOR Universitas Negeri Yogyakarta" />
          </div>
          <div class="form-group">
            <label class="form-label">Batas Pendaftaran</label>
            <input v-model="formBatas" class="form-input" type="date" />
          </div>
        </div>

        <!-- BIAYA PESERTA -->
        <div class="form-section">
          <div class="form-section-title"><i class="ti ti-wallet"></i> Biaya Peserta</div>

          <!-- Latgab: Free / Paid Toggle -->
          <div v-if="formJenis === 'Latgab'">
            <div class="biaya-toggle-row">
              <div :class="['toggle-btn', { active: formBiayaFree }]" @click="formBiayaFree = true">Gratis</div>
              <div :class="['toggle-btn', { active: !formBiayaFree }]" @click="formBiayaFree = false">Berbayar</div>
            </div>
            <div v-if="!formBiayaFree" class="form-group" style="margin-top: 8px;">
              <label class="form-label">Nominal Biaya (Rp)</label>
              <input v-model="formBiayaNominal" type="number" class="form-input" placeholder="27500" />
            </div>
          </div>

          <!-- EKT Jurus: Multi levels breakdown -->
          <div v-else-if="formJenis === 'EKT Jurus'" class="ekt-levels-card">
            <div class="level-row">
              <span>Tingkat Dasar – Gabungan</span>
              <input v-model="formEktFees.dasar" class="form-input-sm" />
            </div>
            <div class="level-row">
              <span>Tingkat PK</span>
              <input v-model="formEktFees.pk" class="form-input-sm" />
            </div>
            <div class="level-row">
              <span>Tingkat GPK – Penjuru</span>
              <input v-model="formEktFees.gpk" class="form-input-sm" />
            </div>
          </div>

          <!-- EKT Non Jurus checkboxes -->
          <div v-else-if="formJenis === 'EKT Non Jurus'" class="ekt-nj-checkboxes">
            <label class="nj-cb-card">
              <input type="checkbox" v-model="formEktNjOpts" value="Pengayaan I" />
              <div>
                <strong>Pengayaan I</strong>
                <span class="subtext">Syarat Gabungan · PH Jurus 6</span>
              </div>
              <span class="fee">Rp 30.000</span>
            </label>
            <label class="nj-cb-card">
              <input type="checkbox" v-model="formEktNjOpts" value="Pengayaan II" />
              <div>
                <strong>Pengayaan II</strong>
                <span class="subtext">Syarat PK · Gabungan Jurus 4</span>
              </div>
              <span class="fee">Rp 45.000</span>
            </label>
          </div>

          <!-- Pelatnas breakdown -->
          <div v-else-if="formJenis === 'Pelatnas'">
            <div class="form-group">
              <label class="form-label">Biaya Pokok Pelatnas (Rp)</label>
              <input v-model="formBiayaNominal" class="form-input" type="number" placeholder="350000" />
            </div>
            <div class="ekt-levels-card">
              <div class="level-row">
                <span>Zona Jawa</span>
                <input v-model="formPelatnasZones.jawa" class="form-input-sm" />
              </div>
              <div class="level-row">
                <span>Luar Jawa</span>
                <input v-model="formPelatnasZones.luar" class="form-input-sm" />
              </div>
            </div>
          </div>
        </div>

        <!-- REKENING & DESKRIPSI -->
        <div v-if="formJenis !== 'Latgab' || !formBiayaFree" class="form-section">
          <div class="form-section-title"><i class="ti ti-building-bank"></i> Rekening Pembayaran</div>
          <div class="form-row-2">
            <div class="form-group">
              <label class="form-label">Bank</label>
              <select v-model="formBank" class="form-input">
                <option value="BCA">BCA</option>
                <option value="BRI">BRI</option>
                <option value="Mandiri">Mandiri</option>
                <option value="BSI">BSI</option>
              </select>
            </div>
            <div class="form-group">
              <label class="form-label">No. Rekening</label>
              <input v-model="formRekening" class="form-input" placeholder="123456789" />
            </div>
          </div>
        </div>
      </div>

      <div class="pf-footer">
        <button class="btn btn-primary" style="width: 100%; justify-content: center;" @click="saveEvent">
          <i class="ti ti-calendar-plus"></i> {{ editingEventId ? 'Simpan Perubahan' : 'Simpan & Publikasikan Event' }}
        </button>
      </div>
    </div>

    <!-- ===== PANEL KANAN: LIST EVENT ===== -->
    <div class="panel-list">
      <div class="pl-toolbar">
        <h2 style="font-size: 16px; font-weight: 800; color: var(--text1)">Daftar Event Terdaftar</h2>
        <div style="display: flex; gap: 8px;">
          <select v-model="filterType" class="filter-select">
            <option value="">Semua Jenis</option>
            <option value="Latgab">Latgab</option>
            <option value="EKT Jurus">EKT Jurus</option>
            <option value="EKT Non Jurus">EKT NJ</option>
            <option value="Pelatnas">Pelatnas</option>
          </select>
          <select v-model="filterStatus" class="filter-select">
            <option value="">Semua Status</option>
            <option value="aktif">Aktif</option>
            <option value="selesai">Selesai</option>
          </select>
        </div>
      </div>

      <div class="list-scroll">
        <table class="event-table">
          <thead>
            <tr>
              <th>Event</th>
              <th>Tanggal</th>
              <th>Cabang</th>
              <th>Peserta</th>
              <th>Biaya</th>
              <th>Status</th>
              <th>Aksi</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="ev in paginatedEvents" :key="ev.id" class="event-row" @click="openDetail(ev)">
              <td>
                <div class="event-name">{{ ev.nama }}</div>
                <div class="event-meta">
                  <span :class="['evt-badge', ev.jenis.toLowerCase().replace(' ', '-')]">{{ ev.jenis }}</span>
                </div>
              </td>
              <td>
                <div style="font-size:12px; font-weight:600;">{{ formatDate(ev.tanggal) }}</div>
                <div style="font-size:10px; color:var(--text3);">{{ ev.waktu }}</div>
              </td>
              <td>
                <div style="font-size:12px;">{{ ev.penyelenggara }}</div>
              </td>
              <td>
                <div class="peserta-bar">
                  <div class="pb-wrap">
                    <div class="pb-fill" :style="{ width: (ev.peserta / ev.targetPeserta * 100) + '%', background: getProgressBarColor(ev.jenis) }"></div>
                  </div>
                  <span style="font-size:11px; font-weight:700;">{{ ev.peserta }} / {{ ev.targetPeserta }}</span>
                </div>
              </td>
              <td>
                <div style="font-size:12px; font-weight:700;">{{ formatBiayaLabel(ev) }}</div>
              </td>
              <td>
                <span :class="['status-badge', ev.status]">{{ ev.status === 'aktif' ? 'Buka' : 'Selesai' }}</span>
              </td>
              <td @click.stopPropagation>
                <div style="display: flex; gap: 4px;">
                  <button class="icon-btn" title="Edit" @click="editEvent(ev)"><i class="ti ti-edit"></i></button>
                  <button class="icon-btn danger" title="Hapus" @click="deleteEvent(ev.id)"><i class="ti ti-trash"></i></button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
        <Pagination 
          v-model:currentPage="pageEvents" 
          v-model:itemsPerPage="limitEvents" 
          :totalItems="filteredEvents.length" 
        />
      </div>
    </div>

    <!-- DETAIL MODAL -->
    <div v-if="showDetailModal && selectedEvent" class="modal-overlay" @click.self="showDetailModal = false">
      <div class="modal-box-detail">
        <div class="modal-head">
          <div class="modal-title">Detail Event</div>
          <button class="modal-close-btn" @click="showDetailModal = false"><i class="ti ti-x"></i></button>
        </div>
        <div class="modal-body" style="padding: 20px;">
          <!-- Detail Event Header -->
          <div :class="['event-detail-header-card', selectedEvent.jenis.toLowerCase().replace(' ', '-')]">
            <div class="edh-type">{{ selectedEvent.jenis.toUpperCase() }}</div>
            <div class="edh-name">{{ selectedEvent.nama }}</div>
            <div class="edh-meta">
              <span class="edh-meta-item"><i class="ti ti-calendar"></i> {{ formatDate(selectedEvent.tanggal) }}</span>
              <span class="edh-meta-item"><i class="ti ti-clock"></i> {{ selectedEvent.waktu }}</span>
              <span class="edh-meta-item"><i class="ti ti-map-pin"></i> {{ selectedEvent.lokasi }}</span>
            </div>
          </div>

          <!-- Info Grid -->
          <div class="info-grid">
            <div class="info-box"><div class="info-box-label">Penyelenggara</div><div class="info-box-val">{{ selectedEvent.penyelenggara }}</div></div>
            <div class="info-box"><div class="info-box-label">Batas Daftar</div><div class="info-box-val">{{ formatDate(selectedEvent.batasDaftar) }}</div></div>
            <div class="info-box"><div class="info-box-label">Biaya Pokok</div><div class="info-box-val">{{ formatBiayaLabel(selectedEvent) }}</div></div>
            <div class="info-box"><div class="info-box-label">Status</div><div class="info-box-val"><span :class="['status-badge', selectedEvent.status]">{{ selectedEvent.status === 'aktif' ? 'Pendaftaran Buka' : 'Selesai' }}</span></div></div>
          </div>

          <!-- Keikutsertaan / Keuangan Summary -->
          <div class="keikutsertaan-box">
            <div class="keikutsertaan-title"><i class="ti ti-users"></i> Keikutsertaan & Keuangan Event</div>
            <div class="keikutsertaan-grid">
              <div class="kb-summary-card">
                <div class="kb-lbl">Peserta Terdaftar</div>
                <div class="kb-val" style="color: var(--hijau);">{{ selectedEvent.peserta }} <span style="font-size: 11px; color: var(--text3); font-weight: normal;">dari {{ selectedEvent.targetPeserta }}</span></div>
              </div>
              <div class="kb-summary-card">
                <div class="kb-lbl">Estimasi Keuangan</div>
                <div class="kb-val">{{ calculateKeuangan(selectedEvent) }}</div>
              </div>
            </div>
          </div>

          <!-- Participant table listing -->
          <div class="participant-section">
            <div class="part-header">
              <div style="font-size: 12px; font-weight: 700;">Daftar Peserta Terdaftar ({{ selectedEvent.peserta }} orang)</div>
              <div style="display: flex; gap: 6px;">
                <button class="export-btn" @click="exportData('PDF')"><i class="ti ti-file-type-pdf"></i> PDF</button>
                <button class="export-btn" @click="exportData('Excel')"><i class="ti ti-file-spreadsheet"></i> Excel</button>
              </div>
            </div>
            <div class="part-search-bar">
              <i class="ti ti-search"></i>
              <input v-model="pesertaSearch" placeholder="Cari nama peserta..." class="part-search-input" />
            </div>

            <div class="part-table-wrapper">
              <table class="part-table">
                <thead>
                  <tr>
                    <th>Nama</th>
                    <th>Cabang / Unit</th>
                    <th>Tingkatan</th>
                    <th style="text-align: right;">Status</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="p in paginatedParticipants" :key="p.id">
                    <td>
                      <div style="font-weight: 600;">{{ p.nama }}</div>
                      <div style="font-size: 10px; color: var(--text3)">{{ p.noAnggota }}</div>
                    </td>
                    <td>{{ p.cabang }} · {{ p.unit }}</td>
                    <td><span class="bdg-level">{{ p.tingkatan }}</span></td>
                    <td style="text-align: right;"><span class="status-badge aktif">Terdaftar</span></td>
                  </tr>
                </tbody>
              </table>
              <Pagination 
                v-model:currentPage="pageParticipants" 
                v-model:itemsPerPage="limitParticipants" 
                :totalItems="filteredParticipants.length" 
              />
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({ title: 'Latgab / EKT / Pelatnas' })

const filterType = ref('')
const filterStatus = ref('')

// Form State
const editingEventId = ref<number | null>(null)
const formJenis = ref('Latgab')
const formCabang = ref('Kota Yogyakarta')
const formPesertaMode = ref('provinsi')
const formProvinsi = ref('DI Yogyakarta')
const formSelectedCabangs = ref<string[]>([])
const formNama = ref('')
const formTanggal = ref('2026-07-23')
const formJamMulai = ref('07:00')
const formJamSelesai = ref('17:00')
const formLokasi = ref('')
const formBatas = ref('2026-07-18')
const formBiayaFree = ref(true)
const formBiayaNominal = ref<number | null>(null)
const formEktFees = ref({ dasar: '57.500', pk: '77.500', gpk: '85.000' })
const formEktNjOpts = ref<string[]>(['Pengayaan I'])
const formPelatnasZones = ref({ jawa: '150.000', luar: '400.000' })
const formBank = ref('BCA')
const formRekening = ref('')

const setJenis = (jenis: string) => {
  formJenis.value = jenis
  if (jenis === 'Latgab') {
    formNama.value = 'Latgab Bersama Cabang Yogyakarta'
    formBiayaFree.value = true
  } else if (jenis === 'EKT Jurus') {
    formNama.value = 'Latgab + EKT Jurus Wilayah DIY'
    formBiayaFree.value = false
    formBiayaNominal.value = 57500
  } else if (jenis === 'EKT Non Jurus') {
    formNama.value = 'EKT Non Jurus (Pengayaan I) DIY'
    formBiayaFree.value = false
    formBiayaNominal.value = 30000
  } else {
    formNama.value = 'Pelatihan Nasional SN 2026'
    formBiayaFree.value = false
    formBiayaNominal.value = 350000
  }
}

// Notification Scopes descriptions (matches mockup updateNotifScope)
const notifScopeTitle = computed(() => {
  if (formJenis.value === 'Pelatnas') return 'Nasional (Seluruh Indonesia)'
  if (formPesertaMode.value === 'provinsi') return `Seluruh Cabang di ${formProvinsi.value}`
  return `${formSelectedCabangs.value.length} Cabang Terpilih`
})

const notifScopeSub = computed(() => {
  if (formJenis.value === 'Pelatnas') return 'Anggota aktif di seluruh Indonesia akan menerima pengumuman di aplikasi'
  return 'Hanya anggota terdaftar di lingkup di atas yang menerima notifikasi'
})

// Hardcoded initial event data to match the mockup
const events = ref<any[]>([
  {
    id: 1,
    jenis: 'Latgab',
    nama: 'Latgab Se-Jawa Barat 2026',
    tanggal: '2026-06-22',
    waktu: '07.00–12.00 WIB',
    lokasi: 'Lapangan Gasibu, Bandung',
    penyelenggara: 'Kota Bandung',
    peserta: 187,
    targetPeserta: 300,
    batasDaftar: '2026-06-18',
    status: 'aktif',
    biayaFree: true,
    biayaNominal: 0
  },
  {
    id: 2,
    jenis: 'EKT Jurus',
    nama: 'EKT Cabang Bandung 2026',
    tanggal: '2026-06-15',
    waktu: '08.00–12.00 WIB',
    lokasi: 'Gedung Silat Bandung',
    penyelenggara: 'Kota Bandung',
    peserta: 42,
    targetPeserta: 50,
    batasDaftar: '2026-06-12',
    status: 'aktif',
    biayaFree: false,
    biayaNominal: 75000
  },
  {
    id: 3,
    jenis: 'Pelatnas',
    nama: 'Pelatnas Tahunan 2026 — Yogyakarta',
    tanggal: '2026-08-10',
    waktu: '08.00–17.00 WIB',
    lokasi: 'Pendopo Pusat, Yogyakarta',
    penyelenggara: 'SN Pusat',
    peserta: 312,
    targetPeserta: 500,
    batasDaftar: '2026-07-31',
    status: 'aktif',
    biayaFree: false,
    biayaNominal: 350000
  },
  {
    id: 4,
    jenis: 'Latgab',
    nama: 'Latgab Se-Jawa Tengah 2026',
    tanggal: '2026-06-08',
    waktu: '07.00–11.00 WIB',
    lokasi: 'Lapangan Simpang Lima, Semarang',
    penyelenggara: 'Kota Semarang',
    peserta: 210,
    targetPeserta: 200,
    batasDaftar: '2026-06-05',
    status: 'selesai',
    biayaFree: false,
    biayaNominal: 30000
  }
])

const filteredEvents = computed(() => {
  return events.value.filter(e => {
    const matchType = !filterType.value || e.jenis === filterType.value
    const matchStatus = !filterStatus.value || e.status === filterStatus.value
    return matchType && matchStatus
  })
})

// Event List Pagination states
const pageEvents = ref(1)
const limitEvents = ref(10)
const paginatedEvents = computed(() => {
  return filteredEvents.value.slice((pageEvents.value - 1) * limitEvents.value, pageEvents.value * limitEvents.value)
})

const getProgressBarColor = (jenis: string) => {
  if (jenis === 'Latgab') return 'var(--biru)'
  if (jenis === 'EKT Jurus') return 'var(--merah)'
  if (jenis === 'EKT Non Jurus') return '#8e24aa'
  return '#7c3aed'
}

const formatBiayaLabel = (ev: any) => {
  if (ev.biayaFree) return 'Gratis'
  return `Rp ${ev.biayaNominal.toLocaleString('id-ID')}`
}

const formatDate = (dateStr: string) => {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  return d.toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' })
}

const saveEvent = () => {
  if (!formNama.value) {
    alert('Nama event harus diisi!')
    return
  }

  if (editingEventId.value) {
    const idx = events.value.findIndex(e => e.id === editingEventId.value)
    if (idx !== -1) {
      events.value[idx] = {
        ...events.value[idx],
        nama: formNama.value,
        jenis: formJenis.value,
        tanggal: formTanggal.value,
        waktu: `${formJamMulai.value}–${formJamSelesai.value} WIB`,
        lokasi: formLokasi.value || 'Gedung Serbaguna',
        penyelenggara: formJenis.value === 'Pelatnas' ? 'SN Pusat' : formCabang.value,
        batasDaftar: formBatas.value,
        biayaFree: formJenis.value === 'Latgab' ? formBiayaFree.value : false,
        biayaNominal: formBiayaFree.value && formJenis.value === 'Latgab' ? 0 : (formBiayaNominal.value || 0)
      }
    }
    editingEventId.value = null
    alert('Event berhasil diperbarui!')
  } else {
    events.value.unshift({
      id: Date.now(),
      jenis: formJenis.value,
      nama: formNama.value,
      tanggal: formTanggal.value,
      waktu: `${formJamMulai.value}–${formJamSelesai.value} WIB`,
      lokasi: formLokasi.value || 'Lapangan',
      penyelenggara: formJenis.value === 'Pelatnas' ? 'SN Pusat' : formCabang.value,
      peserta: 0,
      targetPeserta: formJenis.value === 'Pelatnas' ? 500 : 100,
      batasDaftar: formBatas.value,
      status: 'aktif',
      biayaFree: formJenis.value === 'Latgab' ? formBiayaFree.value : false,
      biayaNominal: formBiayaFree.value && formJenis.value === 'Latgab' ? 0 : (formBiayaNominal.value || 0)
    })
    alert('Event baru berhasil diterbitkan!')
  }

  // Reset form
  formNama.value = ''
  formLokasi.value = ''
  formRekening.value = ''
}

const editEvent = (ev: any) => {
  editingEventId.value = ev.id
  formNama.value = ev.nama
  formJenis.value = ev.jenis
  formTanggal.value = ev.tanggal
  formLokasi.value = ev.lokasi
  formBatas.value = ev.batasDaftar
  formBiayaFree.value = ev.biayaFree
  formBiayaNominal.value = ev.biayaNominal
}

const deleteEvent = (id: number) => {
  if (confirm('Apakah Anda yakin ingin menghapus event ini?')) {
    events.value = events.value.filter(e => e.id !== id)
  }
}

// Modal Detail Details
const showDetailModal = ref(false)
const selectedEvent = ref<any>(null)
const pesertaSearch = ref('')

const openDetail = (ev: any) => {
  selectedEvent.value = ev
  showDetailModal.value = true
}

const calculateKeuangan = (ev: any) => {
  if (ev.biayaFree) return 'Gratis · Rp 0'
  const total = ev.peserta * ev.biayaNominal
  return `Rp ${total.toLocaleString('id-ID')}`
}

const participants = ref<any[]>([
  { id: 1, nama: 'Budi Santoso', noAnggota: 'YO-YGY-00034', cabang: 'Kota Yogyakarta', unit: 'Unit Kotagede', tingkatan: 'PH Jurus 6' },
  { id: 2, nama: 'Rina Wulandari', noAnggota: 'JB-BDG-00142', cabang: 'Kota Bandung', unit: 'Unit Dago', tingkatan: 'Gabungan Jurus 3' },
  { id: 3, nama: 'Ahmad Santoso', noAnggota: 'YO-YGY-00111', cabang: 'Kota Yogyakarta', unit: 'Unit Malioboro', tingkatan: 'Dasar Jurus 5' },
  { id: 4, nama: 'Dewi Wardani', noAnggota: 'JT-SMG-00089', cabang: 'Kota Semarang', unit: 'Unit Simpang Lima', tingkatan: 'PK Jurus 2' },
  { id: 5, nama: 'Hendra Nugraha', noAnggota: 'JK-JKT-00201', cabang: 'Jakarta Pusat', unit: 'Unit Menteng', tingkatan: 'Pra Dasar Jurus 5' }
])

const filteredParticipants = computed(() => {
  if (!pesertaSearch.value) return participants.value
  return participants.value.filter(p => p.nama.toLowerCase().includes(pesertaSearch.value.toLowerCase()))
})

// Participant List Pagination states
const pageParticipants = ref(1)
const limitParticipants = ref(10)
const paginatedParticipants = computed(() => {
  return filteredParticipants.value.slice((pageParticipants.value - 1) * limitParticipants.value, pageParticipants.value * limitParticipants.value)
})

const exportData = (format: string) => {
  alert(`Data peserta berhasil diekspor ke format ${format}!`)
}
</script>

<style scoped>
.latgab-split-layout { display: flex; height: calc(100vh - 65px); margin: -16px; background: var(--surface); }

/* Left panel form */
.panel-form { width: 360px; background: var(--card); border-right: 1px solid var(--border); display: flex; flex-direction: column; height: 100%; flex-shrink: 0; }
.pf-head { padding: 16px 20px; border-bottom: 1px solid var(--border); }
.pf-title { font-size: 15px; font-weight: 800; color: var(--text1); }
.pf-sub { font-size: 11px; color: var(--text3); margin-top: 2px; }
.pf-scroll { flex: 1; overflow-y: auto; padding: 20px; }
.pf-footer { padding: 12px 20px; border-top: 1px solid var(--border); background: var(--card); }

/* Form inputs & styles */
.form-section { margin-bottom: 18px; }
.form-section-title { font-size: 11px; font-weight: 700; color: var(--text2); display: flex; align-items: center; gap: 6px; padding-bottom: 8px; border-bottom: 1px solid var(--border); margin-bottom: 12px; text-transform: uppercase; letter-spacing: 0.5px; }
.form-section-title i { color: var(--hijau); }
.form-group { margin-bottom: 12px; }
.form-label { font-size: 11px; font-weight: 600; color: var(--text2); margin-bottom: 5px; display: block; }
.form-input { width: 100%; padding: 7px 10px; border: 1px solid var(--border2); border-radius: var(--r8); font-size: 12px; background: #fff; color: var(--text1); outline: none; }
.form-input:focus { border-color: var(--hijau); }
.form-row-2 { display: grid; grid-template-columns: 1fr 1fr; gap: 10px; }

/* Toggles & Cards */
.jenis-toggle { display: grid; grid-template-columns: repeat(4, 1fr); gap: 4px; border: 1px solid var(--border); border-radius: var(--r8); padding: 2px; background: var(--surface); }
.jenis-btn { display: flex; flex-direction: column; align-items: center; justify-content: center; padding: 6px 0; border-radius: 6px; cursor: pointer; text-align: center; }
.jenis-btn i { font-size: 16px; color: var(--text3); }
.jenis-btn .jenis-label { font-size: 9px; font-weight: 600; margin-top: 3px; color: var(--text2); }
.jenis-btn.active { background: var(--hijau); }
.jenis-btn.active i, .jenis-btn.active .jenis-label { color: #fff; }

.peserta-tabs { display: flex; border: 1px solid var(--border); border-radius: var(--r8); overflow: hidden; }
.peserta-tab { flex: 1; padding: 6px; text-align: center; font-size: 11px; font-weight: 600; cursor: pointer; background: var(--surface); color: var(--text2); }
.peserta-tab.active { background: var(--hijau3); color: var(--hijau); }

.cabang-checkboxes { display: grid; grid-template-columns: 1fr; gap: 6px; padding: 8px; border: 1px solid var(--border); border-radius: var(--r8); background: var(--surface); max-height: 100px; overflow-y: auto; }
.cb-label { display: flex; align-items: center; gap: 6px; font-size: 11px; cursor: pointer; }

.peny-nasional { display: flex; align-items: center; gap: 10px; padding: 10px; border-radius: var(--r8); border: 1px solid var(--border); background: var(--surface); }
.peny-nasional i { font-size: 20px; color: var(--hijau); }
.pn-title { font-size: 11px; font-weight: 700; }
.pn-sub { font-size: 10px; color: var(--text3); }

.notif-scope-card { display: flex; align-items: center; gap: 8px; padding: 10px; border-radius: var(--r8); border: 1px solid var(--border); background: var(--surface); }
.notif-scope-card i { font-size: 16px; color: var(--hijau); }
.ns-title { font-size: 11px; font-weight: 700; }
.ns-sub { font-size: 10px; color: var(--text3); line-height: 1.2; }

.biaya-toggle-row { display: flex; border: 1px solid var(--border); border-radius: var(--r8); overflow: hidden; }
.toggle-btn { flex: 1; padding: 6px; text-align: center; font-size: 11px; font-weight: 600; cursor: pointer; background: var(--surface); }
.toggle-btn.active { background: var(--hijau3); color: var(--hijau); }

.ekt-levels-card { display: flex; flex-direction: column; gap: 6px; padding: 8px 12px; border: 1px solid var(--border); border-radius: var(--r8); background: var(--surface); }
.level-row { display: flex; justify-content: space-between; align-items: center; font-size: 11px; }
.form-input-sm { width: 70px; padding: 4px 6px; border: 1px solid var(--border2); border-radius: var(--r6); font-size: 11px; text-align: right; }

.ekt-nj-checkboxes { display: flex; flex-direction: column; gap: 6px; }
.nj-cb-card { display: flex; align-items: center; gap: 8px; padding: 8px 12px; border: 1px solid var(--border); border-radius: var(--r8); cursor: pointer; }
.nj-cb-card strong { font-size: 11px; }
.nj-cb-card .subtext { font-size: 9px; color: var(--text3); display: block; }
.nj-cb-card .fee { font-size: 11px; font-weight: 700; margin-left: auto; color: var(--hijau2); }

/* Right panel list */
.panel-list { flex: 1; display: flex; flex-direction: column; height: 100%; min-width: 0; background: var(--surface); }
.pl-toolbar { padding: 16px 20px; border-bottom: 1px solid var(--border); display: flex; justify-content: space-between; align-items: center; background: var(--card); }
.filter-select { padding: 5px 8px; border: 1px solid var(--border); border-radius: var(--r6); font-size: 11px; background: var(--card); }
.list-scroll { flex: 1; overflow-y: auto; padding: 20px; }

/* Event Table */
.event-table { width: 100%; border-collapse: collapse; background: var(--card); border: 1px solid var(--border); border-radius: var(--r12); overflow: hidden; }
.event-table th { padding: 10px 14px; font-size: 10px; font-weight: 700; color: var(--text3); text-transform: uppercase; text-align: left; background: var(--surface); border-bottom: 1px solid var(--border); }
.event-table td { padding: 12px 14px; font-size: 12px; border-bottom: 1px solid var(--border); vertical-align: middle; }
.event-row { cursor: pointer; transition: background .15s; }
.event-row:hover { background: var(--surface); }
.event-name { font-weight: 700; font-size: 13px; color: var(--text1); }
.event-meta { margin-top: 4px; display: flex; gap: 4px; }

/* Badges */
.evt-badge { display: inline-block; padding: 2px 6px; border-radius: 4px; font-size: 9px; font-weight: 700; color: #fff; }
.evt-badge.latgab { background: var(--biru); }
.evt-badge.ekt-jurus { background: var(--merah); }
.evt-badge.ekt-non-jurus { background: #8e24aa; }
.evt-badge.pelatnas { background: #7c3aed; }

.status-badge { display: inline-block; padding: 2px 8px; border-radius: 10px; font-size: 10px; font-weight: 700; }
.status-badge.aktif { background: var(--hijau3); color: var(--hijau); }
.status-badge.selesai { background: var(--surface); color: var(--text3); }

.peserta-bar { display: flex; align-items: center; gap: 6px; font-size: 11px; }
.pb-wrap { width: 60px; height: 5px; background: var(--border); border-radius: 3px; overflow: hidden; }
.pb-fill { height: 100%; border-radius: 3px; }

.icon-btn { width: 26px; height: 26px; border: 1px solid var(--border); border-radius: var(--r6); background: #fff; cursor: pointer; display: inline-flex; align-items: center; justify-content: center; font-size: 12px; transition: all .15s; }
.icon-btn:hover { border-color: var(--hijau); color: var(--hijau); background: var(--hijau3); }
.icon-btn.danger:hover { border-color: var(--merah); color: var(--merah); background: #fdecea; }

.form-hint { font-size: 9px; color: var(--text3); margin-top: 3px; }

/* Modal overlay & popup */
.modal-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.5); z-index: 1000; display: flex; align-items: center; justify-content: center; padding: 20px; }
.modal-box-detail { background: var(--card); border-radius: var(--r12); width: 100%; max-width: 600px; max-height: 90vh; overflow-y: auto; box-shadow: var(--shadow-lg); }
.modal-head { padding: 16px 20px; border-bottom: 1px solid var(--border); display: flex; justify-content: space-between; align-items: center; }
.modal-title { font-size: 15px; font-weight: 800; }
.modal-close-btn { background: none; border: none; font-size: 16px; cursor: pointer; color: var(--text3); }

/* Detail Event Card inside Modal */
.event-detail-header-card { padding: 16px; border-radius: var(--r8); color: #fff; margin-bottom: 16px; }
.event-detail-header-card.latgab { background: linear-gradient(135deg, var(--biru), #4dabf7); }
.event-detail-header-card.ekt-jurus { background: linear-gradient(135deg, var(--merah), #ff8787); }
.event-detail-header-card.ekt-non-jurus { background: linear-gradient(135deg, #8e24aa, #ba68c8); }
.event-detail-header-card.pelatnas { background: linear-gradient(135deg, #7c3aed, #9f7aea); }

.edh-type { font-size: 9px; font-weight: 800; letter-spacing: 0.5px; opacity: 0.8; }
.edh-name { font-size: 16px; font-weight: 800; margin-top: 4px; }
.edh-meta { display: flex; gap: 12px; margin-top: 10px; font-size: 11px; opacity: 0.9; }
.edh-meta-item { display: flex; align-items: center; gap: 4px; }

.info-grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 8px; margin-bottom: 16px; }
.info-box { padding: 8px 10px; border: 1px solid var(--border); border-radius: var(--r8); background: var(--surface); text-align: center; }
.info-box-label { font-size: 9px; color: var(--text3); text-transform: uppercase; }
.info-box-val { font-size: 12px; font-weight: 700; margin-top: 4px; }

.keikutsertaan-box { background: var(--surface); border: 1px solid var(--border); border-radius: var(--r8); padding: 12px; margin-bottom: 16px; }
.keikutsertaan-title { font-size: 11px; font-weight: 700; display: flex; align-items: center; gap: 6px; margin-bottom: 8px; }
.keikutsertaan-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 8px; }
.kb-summary-card { background: #fff; border: 1px solid var(--border); border-radius: var(--r8); padding: 8px 12px; }
.kb-lbl { font-size: 9px; color: var(--text3); }
.kb-val { font-size: 16px; font-weight: 800; margin-top: 3px; }

.part-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 10px; }
.export-btn { border: 1px solid var(--border); padding: 4px 8px; border-radius: 6px; font-size: 10px; font-weight: 700; cursor: pointer; background: #fff; display: inline-flex; align-items: center; gap: 4px; }
.export-btn:hover { color: var(--hijau); border-color: var(--hijau); }

.part-search-bar { position: relative; margin-bottom: 10px; }
.part-search-bar i { position: absolute; left: 10px; top: 50%; transform: translateY(-50%); font-size: 12px; color: var(--text3); }
.part-search-input { width: 100%; padding: 6px 10px 6px 28px; border: 1px solid var(--border); border-radius: var(--r8); font-size: 12px; outline: none; }
.part-search-input:focus { border-color: var(--hijau); }

.part-table-wrapper { border: 1px solid var(--border); border-radius: var(--r8); overflow: hidden; max-height: 150px; overflow-y: auto; }
.part-table { width: 100%; border-collapse: collapse; font-size: 11px; }
.part-table th { padding: 6px 10px; background: var(--surface); color: var(--text3); font-weight: 700; text-align: left; }
.part-table td { padding: 8px 10px; border-bottom: 1px solid var(--border); }
.part-table tr:last-child td { border-bottom: none; }

.bdg-level { background: var(--hijau3); color: var(--hijau); padding: 1px 4px; border-radius: 4px; font-size: 8px; font-weight: 700; }
</style>
