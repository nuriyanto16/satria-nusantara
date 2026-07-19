<template>
  <div id="page-kebugaran" class="kebugaran-split-layout">
    <!-- ===== PANEL KIRI: FORM BUAT EVENT TES ===== -->
    <div class="panel-form">
      <div class="pf-head">
        <div class="pf-title">Buat Sesi Tes Baru</div>
        <div class="pf-sub">Jadwalkan sesi tes kebugaran per unit</div>
      </div>
      <div class="pf-scroll">
        <!-- LOKASI -->
        <div class="form-section">
          <div class="form-section-title"><i class="ti ti-map-pin"></i> Lokasi Tes</div>
          <div class="form-group">
            <label class="form-label">Cabang</label>
            <select v-model="formCabang" class="form-input" @change="onCabangChange">
              <option v-for="c in cabangList" :key="c.id" :value="c.nama">{{ c.nama }}</option>
            </select>
          </div>
          <div class="form-group">
            <label class="form-label">Unit Latihan</label>
            <select v-model="formUnit" class="form-input">
              <option v-for="u in unitList" :key="u.id" :value="u.nama">{{ u.nama }}</option>
            </select>
          </div>
        </div>

        <!-- JADWAL & PERIODE -->
        <div class="form-section">
          <div class="form-section-title"><i class="ti ti-calendar"></i> Jadwal & Periode</div>
          <div class="form-group">
            <label class="form-label">Tanggal Tes</label>
            <input v-model="formTanggal" type="date" class="form-input" />
          </div>
          <div class="form-group">
            <label class="form-label">Periode</label>
            <select v-model="formPeriode" class="form-input">
              <option value="Periode 1 — 2026">Periode 1 — 2026</option>
              <option value="Periode 2 — 2026">Periode 2 — 2026</option>
              <option value="Periode 3 — 2026">Periode 3 — 2026</option>
            </select>
          </div>
          <div class="form-group">
            <label class="form-label">Catatan (opsional)</label>
            <textarea v-model="formNotes" class="form-input" rows="2" style="height:auto;resize:none;" placeholder="Misal: Tes rutin Q2 2026"></textarea>
          </div>
        </div>

        <!-- STANDAR TES INFO -->
        <div class="form-section">
          <div class="form-section-title"><i class="ti ti-clipboard-list"></i> Parameter Tes Kebugaran (5 Item)</div>
          <div class="standar-tes-box">
            <div class="test-item-row"><span>💨</span> Nafas dalam Air (detik)</div>
            <div class="test-item-row"><span>💪</span> Push Up (repetisi)</div>
            <div class="test-item-row"><span>🏃</span> Sit Up (repetisi)</div>
            <div class="test-item-row"><span>🤸</span> Sit & Reach (cm)</div>
            <div class="test-item-row"><span>⚡</span> Shuttle Run (detik)</div>
            <div style="font-size: 10px; color: var(--text3); border-top: 1px solid var(--border); padding-top: 6px; margin-top: 6px;">
              Kategori & standar nilai disesuaikan berdasarkan usia & gender anggota.
            </div>
          </div>
        </div>

        <button class="btn btn-primary" style="width: 100%; justify-content: center;" @click="createSesiTes">
          <i class="ti ti-plus"></i> Buat Event Tes
        </button>
      </div>
    </div>

    <!-- ===== PANEL KANAN: TABS & DATA ===== -->
    <div class="panel-right">
      <!-- Tabs header -->
      <div class="tab-bar">
        <button :class="['tab-btn', { active: activeTab === 'list' }]" @click="activeTab = 'list'">
          <i class="ti ti-list"></i> List Event Tes
        </button>
        <button :class="['tab-btn', { active: activeTab === 'rekap' }]" @click="activeTab = 'rekap'">
          <i class="ti ti-chart-bar"></i> Rekap & Ranking
        </button>
      </div>

      <!-- TAB 1: LIST EVENT TES -->
      <div v-if="activeTab === 'list'" class="tab-content scrollable-content">
        <!-- Summary Cards -->
        <div class="kb-summary-cards">
          <div class="kb-summary-card">
            <div class="kb-sum-lbl">Total Event</div>
            <div class="kb-sum-val" style="color: var(--hijau);">{{ sesiTesList.length }}</div>
            <div class="kb-sum-sub">Periode Terpilih</div>
          </div>
          <div class="kb-summary-card">
            <div class="kb-sum-lbl">Peserta Dites</div>
            <div class="kb-sum-val" style="color: var(--biru);">247</div>
            <div class="kb-sum-sub">Anggota Aktif</div>
          </div>
          <div class="kb-summary-card">
            <div class="kb-sum-lbl">Rata-rata Skor</div>
            <div class="kb-sum-val" style="color: var(--kuning);">72%</div>
            <div class="kb-sum-sub">Kategori Cukup</div>
          </div>
          <div class="kb-summary-card">
            <div class="kb-sum-lbl">Input Pending</div>
            <div class="kb-sum-val" style="color: var(--merah);">1</div>
            <div class="kb-sum-sub">Belum Selesai</div>
          </div>
        </div>

        <!-- Table Grid -->
        <div class="table-card">
          <table class="event-table">
            <thead>
              <tr>
                <th>Tanggal</th>
                <th>Cabang / Unit</th>
                <th>Periode</th>
                <th>Peserta</th>
                <th>Rata-rata</th>
                <th>Status</th>
                <th>Aksi</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="s in sesiTesList" :key="s.id" class="event-row" @click="viewDetail(s)">
                <td>
                  <div style="font-weight: 700;">{{ formatDate(s.tanggal) }}</div>
                  <div style="font-size: 10px; color: var(--text3)">Sabtu</div>
                </td>
                <td>
                  <div style="font-weight: 600;">{{ s.cabang }}</div>
                  <div style="font-size: 11px; color: var(--text3)">{{ s.unit }}</div>
                </td>
                <td style="font-size: 12px;">{{ s.periode }}</td>
                <td style="font-weight: 700;">
                  {{ s.selesai }} <span style="font-weight: normal; color: var(--text3); font-size: 10px;">/ {{ s.total }}</span>
                </td>
                <td>
                  <div class="pct-cell">
                    <div class="pct-bar"><div class="pct-fill" :style="{ width: s.avgScore + '%', background: getScoreColor(s.avgScore) }"></div></div>
                    <span style="font-size: 11px; font-weight: 700;" :style="{ color: getScoreColor(s.avgScore) }">{{ s.avgScore }}%</span>
                  </div>
                </td>
                <td>
                  <span :class="['status-badge', getStatusClass(s)]">{{ getStatusLabel(s) }}</span>
                </td>
                <td @click.stopPropagation>
                  <div style="display: flex; gap: 4px;">
                    <button v-if="s.selesai < s.total" class="icon-btn" title="Input Data" @click="openInput(s)">
                      <i class="ti ti-edit"></i>
                    </button>
                    <button class="icon-btn" title="Lihat Detail" @click="viewDetail(s)">
                      <i class="ti ti-eye"></i>
                    </button>
                    <button class="icon-btn danger" title="Hapus" @click="deleteSesi(s.id)">
                      <i class="ti ti-trash"></i>
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- TAB 2: REKAP & RANKING -->
      <div v-else-if="activeTab === 'rekap'" class="tab-content scrollable-content" style="gap: 16px;">
        <div class="rekap-grid">
          <!-- Distribusi Kategori -->
          <div class="rekap-card">
            <div class="card-title">Distribusi Kategori Kebugaran</div>
            <div class="distribusi-list">
              <div class="dist-row">
                <div class="dist-header"><span>🟢 Baik</span> <strong>68 (28%)</strong></div>
                <div class="prog-track"><div class="prog-fill" style="width: 28%; background: var(--hijau)"></div></div>
              </div>
              <div class="dist-row">
                <div class="dist-header"><span>🟡 Cukup</span> <strong>134 (54%)</strong></div>
                <div class="prog-track"><div class="prog-fill" style="width: 54%; background: var(--kuning)"></div></div>
              </div>
              <div class="dist-row">
                <div class="dist-header"><span>🔴 Kurang</span> <strong>45 (18%)</strong></div>
                <div class="prog-track"><div class="prog-fill" style="width: 18%; background: var(--merah)"></div></div>
              </div>
            </div>
          </div>

          <!-- Rata-rata Tes -->
          <div class="rekap-card">
            <div class="card-title">Rata-rata per Jenis Tes</div>
            <div class="averages-list">
              <div class="avg-row"><span>💨 Nafas dalam Air</span> <span class="badge badge-green">Baik · 87 dtk</span></div>
              <div class="avg-row"><span>💪 Push Up</span> <span class="badge badge-yellow">Cukup · 18 rep</span></div>
              <div class="avg-row"><span>🏃 Sit Up</span> <span class="badge badge-yellow">Cukup · 21 rep</span></div>
              <div class="avg-row"><span>🤸 Sit & Reach</span> <span class="badge badge-green">Baik · 24 cm</span></div>
              <div class="avg-row"><span>⚡ Shuttle Run</span> <span class="badge badge-red">Kurang · 14.2 dtk</span></div>
            </div>
          </div>
        </div>

        <!-- Ranking Unit -->
        <div class="table-card">
          <div class="table-title-bar">🏆 Ranking Unit Paling Bugar</div>
          <table class="event-table">
            <thead>
              <tr>
                <th>Rank</th>
                <th>Unit</th>
                <th>Cabang</th>
                <th>Peserta</th>
                <th>Skor Rata-rata</th>
                <th>Kategori</th>
              </tr>
            </thead>
            <tbody>
              <tr>
                <td><div class="rank-badge gold">1</div></td>
                <td><strong>Unit Malioboro</strong></td>
                <td>Kota Yogyakarta</td>
                <td>32</td>
                <td>
                  <div class="pct-cell">
                    <div class="pct-bar"><div class="pct-fill" style="width: 78%; background: var(--hijau)"></div></div>
                    <strong>78%</strong>
                  </div>
                </td>
                <td><span class="status-badge selesai">Baik</span></td>
              </tr>
              <tr>
                <td><div class="rank-badge silver">2</div></td>
                <td><strong>Unit Kotagede</strong></td>
                <td>Kota Yogyakarta</td>
                <td>28</td>
                <td>
                  <div class="pct-cell">
                    <div class="pct-bar"><div class="pct-fill" style="width: 65%; background: var(--kuning)"></div></div>
                    <strong>65%</strong>
                  </div>
                </td>
                <td><span class="status-badge warn">Cukup</span></td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- MODAL INPUT HASIL TES -->
    <div v-if="showInputModal && selectedSesi" class="modal-overlay" @click.self="showInputModal = false">
      <div class="modal-box large">
        <div class="modal-head">
          <div class="modal-title">Input Hasil Tes Kebugaran</div>
          <button class="modal-close" @click="showInputModal = false"><i class="ti ti-x"></i></button>
        </div>
        <div class="modal-body" style="padding: 20px;">
          <div class="event-detail-header-card kebugaran">
            <div class="edh-type">TES KEBUGARAN · INPUT DATA</div>
            <div class="edh-name">{{ selectedSesi.unit }} · {{ selectedSesi.cabang }}</div>
            <div class="edh-meta">
              <span class="edh-meta-item"><i class="ti ti-calendar"></i> {{ formatDate(selectedSesi.tanggal) }}</span>
              <span class="edh-meta-item"><i class="ti ti-clock"></i> {{ selectedSesi.periode }}</span>
              <span class="edh-meta-item"><i class="ti ti-users"></i> {{ selectedSesi.selesai }} / {{ selectedSesi.total }} Diinput</span>
            </div>
          </div>

          <div class="table-card" style="margin-top: 16px;">
            <table class="part-table">
              <thead>
                <tr>
                  <th>Anggota</th>
                  <th style="width: 90px;">Nafas (dtk)</th>
                  <th style="width: 80px;">Push Up</th>
                  <th style="width: 80px;">Sit Up</th>
                  <th style="width: 80px;">Sit & Reach</th>
                  <th style="width: 90px;">Shuttle Run</th>
                  <th style="width: 90px;">Kategori</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="p in inputPesertaList" :key="p.id">
                  <td>
                    <div style="font-weight: 700;">{{ p.nama }}</div>
                    <div style="font-size: 10px; color: var(--text3)">{{ p.noAnggota }}</div>
                  </td>
                  <td><input v-model="p.nafas" type="number" class="form-input-table" @change="recalculateCategory(p)" /></td>
                  <td><input v-model="p.pushUp" type="number" class="form-input-table" @change="recalculateCategory(p)" /></td>
                  <td><input v-model="p.sitUp" type="number" class="form-input-table" @change="recalculateCategory(p)" /></td>
                  <td><input v-model="p.sitReach" type="number" class="form-input-table" @change="recalculateCategory(p)" /></td>
                  <td><input v-model="p.shuttleRun" type="number" step="0.1" class="form-input-table" @change="recalculateCategory(p)" /></td>
                  <td><span :class="['kategori-badge', p.kategori.toLowerCase()]">{{ p.kategori }}</span></td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
        <div class="modal-footer" style="padding: 12px 20px;">
          <button class="btn btn-outline" @click="showInputModal = false">Batal</button>
          <button class="btn btn-primary" @click="saveInputResults">Simpan Hasil Tes</button>
        </div>
      </div>
    </div>

    <!-- MODAL DETAIL EVENT TES -->
    <div v-if="showDetailModal && selectedSesi" class="modal-overlay" @click.self="showDetailModal = false">
      <div class="modal-box large">
        <div class="modal-head">
          <div class="modal-title">Hasil Kebugaran Anggota</div>
          <button class="modal-close" @click="showDetailModal = false"><i class="ti ti-x"></i></button>
        </div>
        <div class="modal-body" style="padding: 20px;">
          <div class="event-detail-header-card kebugaran">
            <div class="edh-type">TES KEBUGARAN · LAPORAN LENGKAP</div>
            <div class="edh-name">{{ selectedSesi.unit }} · {{ selectedSesi.cabang }}</div>
            <div class="edh-meta">
              <span class="edh-meta-item"><i class="ti ti-calendar"></i> {{ formatDate(selectedSesi.tanggal) }}</span>
              <span class="edh-meta-item"><i class="ti ti-clock"></i> {{ selectedSesi.periode }}</span>
            </div>
          </div>

          <div style="display: flex; justify-content: space-between; align-items: center; margin-top: 16px; margin-bottom: 8px;">
            <div style="font-size: 13px; font-weight: 700;">Skor Pengujian Kebugaran</div>
            <button class="export-btn" @click="exportResults"><i class="ti ti-file-spreadsheet"></i> Export Excel</button>
          </div>

          <div class="table-card">
            <table class="part-table">
              <thead>
                <tr>
                  <th>No</th>
                  <th>Nama Anggota</th>
                  <th>Nafas</th>
                  <th>Push Up</th>
                  <th>Sit Up</th>
                  <th>Sit & Reach</th>
                  <th>Shuttle Run</th>
                  <th>Kategori</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="(p, idx) in inputPesertaList" :key="p.id">
                  <td>{{ idx + 1 }}</td>
                  <td>
                    <div style="font-weight: 700;">{{ p.nama }}</div>
                    <div style="font-size: 10px; color: var(--text3)">{{ p.noAnggota }}</div>
                  </td>
                  <td>{{ p.nafas || '-' }} dtk</td>
                  <td>{{ p.pushUp || '-' }} rep</td>
                  <td>{{ p.sitUp || '-' }} rep</td>
                  <td>{{ p.sitReach || '-' }} cm</td>
                  <td>{{ p.shuttleRun || '-' }} dtk</td>
                  <td><span :class="['kategori-badge', p.kategori.toLowerCase()]">{{ p.kategori }}</span></td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
        <div class="modal-footer" style="padding: 12px 20px;">
          <button class="btn btn-outline" @click="showDetailModal = false">Tutup</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({ title: 'Tes Kebugaran' })

const api = useApi()
const activeTab = ref('list')

// Modals
const showInputModal = ref(false)
const showDetailModal = ref(false)
const selectedSesi = ref<any>(null)

// Form
const formCabang = ref('Kota Yogyakarta')
const formUnit = ref('Unit Malioboro')
const formTanggal = ref('2026-07-25')
const formPeriode = ref('Periode 2 — 2026')
const formNotes = ref('')

const cabangList = ref<any[]>([])
const unitList = ref<any[]>([])

const fetchCabang = async () => {
  try {
    const res = await api.get('/organization/cabang?limit=100')
    cabangList.value = res.data || []
    if (cabangList.value.length > 0) {
      formCabang.value = cabangList.value[0].nama
      onCabangChange()
    }
  } catch (e) {
    console.error(e)
  }
}

const onCabangChange = async () => {
  const matching = cabangList.value.find(c => c.nama === formCabang.value)
  if (!matching) return
  try {
    const res = await api.get(`/organization/cabang/${matching.id}/unit`)
    unitList.value = res || []
    if (unitList.value.length > 0) {
      formUnit.value = unitList.value[0].nama
    }
  } catch (e) {
    console.error(e)
  }
}

// Dummy Test Session Database
const sesiTesList = ref<any[]>([
  { id: 1, tanggal: '2026-06-14', cabang: 'Kota Yogyakarta', unit: 'Unit Malioboro', periode: 'Periode 2 — 2026', selesai: 32, total: 32, avgScore: 78 },
  { id: 2, tanggal: '2026-06-21', cabang: 'Kota Yogyakarta', unit: 'Unit Kotagede', periode: 'Periode 2 — 2026', selesai: 28, total: 31, avgScore: 65 },
  { id: 3, tanggal: '2026-06-22', cabang: 'Kota Bandung', unit: 'Unit Dago', periode: 'Periode 2 — 2026', selesai: 0, total: 42, avgScore: 0 }
])

const createSesiTes = () => {
  sesiTesList.value.unshift({
    id: Date.now(),
    tanggal: formTanggal.value,
    cabang: formCabang.value,
    unit: formUnit.value,
    periode: formPeriode.value,
    selesai: 0,
    total: 30,
    avgScore: 0
  })
  alert('Event tes kebugaran baru berhasil dijadwalkan!')
}

const deleteSesi = (id: number) => {
  if (confirm('Apakah Anda yakin ingin menghapus sesi tes ini?')) {
    sesiTesList.value = sesiTesList.value.filter(s => s.id !== id)
  }
}

const getScoreColor = (score: number) => {
  if (score === 0) return 'var(--text3)'
  if (score >= 75) return 'var(--hijau)'
  if (score >= 60) return 'var(--kuning)'
  return 'var(--merah)'
}

const getStatusClass = (s: any) => {
  if (s.selesai === 0) return 'red'
  if (s.selesai < s.total) return 'warn'
  return 'aktif'
}

const getStatusLabel = (s: any) => {
  if (s.selesai === 0) return 'Belum Mulai'
  if (s.selesai < s.total) return 'Input Terbuka'
  return 'Selesai'
}

const formatDate = (dateStr: string) => {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  return d.toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' })
}

// Participants for Test Results Input
const inputPesertaList = ref<any[]>([
  { id: 1, nama: 'Budi Santoso', noAnggota: 'YO-YGY-00034', nafas: 120, pushUp: 35, sitUp: 40, sitReach: 22, shuttleRun: 13.2, kategori: 'Baik' },
  { id: 2, nama: 'Rina Wulandari', noAnggota: 'JB-BDG-00142', nafas: 95, pushUp: 28, sitUp: 32, sitReach: 18, shuttleRun: 14.8, kategori: 'Cukup' },
  { id: 3, nama: 'Ahmad Santoso', noAnggota: 'YO-YGY-00111', nafas: 140, pushUp: 42, sitUp: 50, sitReach: 26, shuttleRun: 12.5, kategori: 'Baik' },
  { id: 4, nama: 'Dewi Wardani', noAnggota: 'JT-SMG-00089', nafas: 80, pushUp: 20, sitUp: 25, sitReach: 12, shuttleRun: 16.1, kategori: 'Kurang' }
])

const openInput = (s: any) => {
  selectedSesi.value = s
  showInputModal.value = true
}

const viewDetail = (s: any) => {
  selectedSesi.value = s
  showDetailModal.value = true
}

const recalculateCategory = (p: any) => {
  if (!p.nafas || !p.pushUp || !p.sitUp) {
    p.kategori = 'Kurang'
    return
  }
  const score = (p.nafas / 120 * 30) + (p.pushUp / 35 * 25) + (p.sitUp / 40 * 25) + (p.sitReach / 20 * 20)
  if (score >= 85) p.kategori = 'Baik'
  else if (score >= 60) p.kategori = 'Cukup'
  else p.kategori = 'Kurang'
}

const saveInputResults = () => {
  if (selectedSesi.value) {
    const idx = sesiTesList.value.findIndex(s => s.id === selectedSesi.value.id)
    if (idx !== -1) {
      sesiTesList.value[idx].selesai = sesiTesList.value[idx].total // Complete it
      sesiTesList.value[idx].avgScore = 74 // Set average score
    }
  }
  showInputModal.value = false
  alert('Hasil tes kebugaran berhasil disimpan!')
}

const exportResults = () => {
  alert('Data hasil tes kebugaran berhasil diekspor ke Excel!')
}

onMounted(() => {
  fetchCabang()
})
</script>

<style scoped>
.kebugaran-split-layout { display: flex; height: calc(100vh - 65px); margin: -16px; background: var(--surface); }

/* Left panel form */
.panel-form { width: 350px; background: var(--card); border-right: 1px solid var(--border); display: flex; flex-direction: column; height: 100%; flex-shrink: 0; }
.pf-head { padding: 16px 20px; border-bottom: 1px solid var(--border); }
.pf-title { font-size: 15px; font-weight: 800; color: var(--text1); }
.pf-sub { font-size: 11px; color: var(--text3); margin-top: 2px; }
.pf-scroll { flex: 1; overflow-y: auto; padding: 20px; }

.form-section { margin-bottom: 18px; }
.form-section-title { font-size: 11px; font-weight: 700; color: var(--text2); display: flex; align-items: center; gap: 6px; padding-bottom: 6px; border-bottom: 1px solid var(--border); margin-bottom: 12px; text-transform: uppercase; letter-spacing: 0.5px; }
.form-section-title i { color: var(--hijau); }
.form-group { margin-bottom: 12px; }
.form-label { font-size: 11px; font-weight: 600; color: var(--text2); margin-bottom: 5px; display: block; }
.form-input { width: 100%; padding: 7px 10px; border: 1px solid var(--border2); border-radius: var(--r8); font-size: 12px; background: #fff; color: var(--text1); outline: none; }
.form-input:focus { border-color: var(--hijau); }

.standar-tes-box { background: var(--surface); border-radius: var(--r8); padding: 10px 12px; border: 1px solid var(--border); }
.test-item-row { display: flex; align-items: center; gap: 8px; font-size: 12px; margin-bottom: 6px; }

/* Right panel tabs & content */
.panel-right { flex: 1; display: flex; flex-direction: column; height: 100%; min-width: 0; background: var(--surface); }
.tab-bar { display: flex; border-bottom: 1px solid var(--border); background: var(--card); padding: 0 16px; }
.tab-btn { display: flex; align-items: center; gap: 6px; padding: 12px 16px; font-size: 12px; font-weight: 600; background: none; border: none; cursor: pointer; color: var(--text2); border-bottom: 2px solid transparent; }
.tab-btn.active { color: var(--hijau); border-bottom-color: var(--hijau); }

.tab-content { flex: 1; display: flex; flex-direction: column; padding: 20px; }
.scrollable-content { overflow-y: auto; }

/* Summary Cards */
.kb-summary-cards { display: grid; grid-template-columns: repeat(4, 1fr); gap: 10px; margin-bottom: 16px; flex-shrink: 0; }
.kb-summary-card { background: var(--card); border-radius: var(--r8); padding: 12px; border-left: 3px solid var(--border); box-shadow: var(--shadow-sm); border: 1px solid var(--border); border-left-width: 3px; }
.kb-sum-lbl { font-size: 10px; color: var(--text3); text-transform: uppercase; }
.kb-sum-val { font-size: 20px; font-weight: 800; margin-top: 4px; }
.kb-sum-sub { font-size: 10px; color: var(--text3); margin-top: 2px; }

/* Tables & Lists */
.table-card { background: var(--card); border: 1px solid var(--border); border-radius: var(--r12); overflow-x: auto; flex-shrink: 0; }
.table-title-bar { padding: 12px 16px; border-bottom: 1px solid var(--border); font-size: 13px; font-weight: 700; }
.event-table { width: 100%; border-collapse: collapse; }
.event-table th { padding: 10px 14px; font-size: 10px; font-weight: 700; color: var(--text3); text-transform: uppercase; text-align: left; background: var(--surface); border-bottom: 1px solid var(--border); }
.event-table td { padding: 12px 14px; font-size: 12px; border-bottom: 1px solid var(--border); vertical-align: middle; }
.event-row { cursor: pointer; transition: background .15s; }
.event-row:hover { background: var(--surface); }

.pct-cell { display: flex; align-items: center; gap: 8px; font-size: 11px; }
.pct-bar { width: 60px; height: 5px; background: var(--border); border-radius: 3px; overflow: hidden; }
.pct-fill { height: 100%; border-radius: 3px; }

.status-badge { display: inline-block; padding: 2px 8px; border-radius: 10px; font-size: 10px; font-weight: 700; }
.status-badge.aktif { background: #fff8e0; color: #a07000; }
.status-badge.warn { background: #fff8e0; color: #a07000; }
.status-badge.selesai { background: var(--hijau3); color: var(--hijau); }
.status-badge.red { background: #fde8e8; color: var(--merah); }

.icon-btn { width: 26px; height: 26px; border: 1px solid var(--border); border-radius: var(--r6); background: #fff; cursor: pointer; display: inline-flex; align-items: center; justify-content: center; font-size: 12px; transition: all .15s; }
.icon-btn:hover { border-color: var(--hijau); color: var(--hijau); background: var(--hijau3); }
.icon-btn.danger:hover { border-color: var(--merah); color: var(--merah); background: #fdecea; }

/* Rekap & Averages Tab style */
.rekap-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 16px; margin-bottom: 16px; }
.rekap-card { background: var(--card); border: 1px solid var(--border); border-radius: var(--r8); padding: 16px; }
.card-title { font-size: 13px; font-weight: 700; margin-bottom: 14px; }
.distribusi-list { display: flex; flex-direction: column; gap: 10px; }
.dist-row { display: flex; flex-direction: column; gap: 4px; }
.dist-header { display: flex; justify-content: space-between; font-size: 11px; }
.prog-track { height: 7px; background: var(--border); border-radius: 4px; }
.prog-fill { height: 100%; border-radius: 4px; }

.averages-list { display: flex; flex-direction: column; gap: 8px; }
.avg-row { display: flex; justify-content: space-between; align-items: center; font-size: 12px; }
.badge { padding: 2px 8px; border-radius: 10px; font-size: 10px; font-weight: 700; border: 1px solid transparent; }
.badge-green { background: var(--hijau3); color: var(--hijau); }
.badge-yellow { background: #fff8e0; color: #a07000; border-color: #e8c42a; }
.badge-red { background: #fde8e8; color: var(--merah); }

.rank-badge { width: 24px; height: 24px; border-radius: 50%; display: flex; align-items: center; justify-content: center; font-size: 11px; font-weight: 800; color: #fff; }
.rank-badge.gold { background: #f59e0b; }
.rank-badge.silver { background: #cbd5e1; color: var(--text2); }

/* Modals */
.modal-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.5); z-index: 1000; display: flex; align-items: center; justify-content: center; padding: 20px; }
.modal-box { background: var(--card); border-radius: var(--r12); width: 100%; max-width: 500px; max-height: 90vh; overflow-y: auto; box-shadow: var(--shadow-lg); }
.modal-box.large { max-width: 760px; }
.modal-head { padding: 16px 20px; border-bottom: 1px solid var(--border); display: flex; justify-content: space-between; align-items: center; }
.modal-title { font-size: 15px; font-weight: 800; }
.modal-close { background: none; border: none; font-size: 16px; cursor: pointer; color: var(--text3); }
.modal-footer { padding: 12px 20px; border-top: 1px solid var(--border); display: flex; justify-content: flex-end; gap: 8px; }

/* Detail Event Card inside Modal */
.event-detail-header-card { padding: 16px; border-radius: var(--r8); color: #fff; }
.event-detail-header-card.kebugaran { background: linear-gradient(135deg, var(--hijau), #51af61); }
.edh-type { font-size: 9px; font-weight: 800; letter-spacing: 0.5px; opacity: 0.8; }
.edh-name { font-size: 16px; font-weight: 800; margin-top: 4px; }
.edh-meta { display: flex; gap: 12px; margin-top: 10px; font-size: 11px; opacity: 0.9; }
.edh-meta-item { display: flex; align-items: center; gap: 4px; }

.part-table { width: 100%; border-collapse: collapse; font-size: 11px; }
.part-table th { padding: 6px 10px; background: var(--surface); color: var(--text3); font-weight: 700; text-align: left; }
.part-table td { padding: 8px 10px; border-bottom: 1px solid var(--border); vertical-align: middle; }
.part-table tr:last-child td { border-bottom: none; }

.form-input-table { width: 100%; padding: 4px 6px; border: 1px solid var(--border2); border-radius: var(--r6); font-size: 11px; text-align: center; }
.kategori-badge { display: inline-block; padding: 2px 6px; border-radius: 8px; font-size: 9px; font-weight: 700; }
.kategori-badge.baik { background: var(--hijau3); color: var(--hijau); }
.kategori-badge.cukup { background: #fff8e0; color: #a07000; }
.kategori-badge.kurang { background: #fde8e8; color: var(--merah); }

.export-btn { border: 1px solid var(--border); padding: 4px 8px; border-radius: 6px; font-size: 10px; font-weight: 700; cursor: pointer; background: #fff; display: inline-flex; align-items: center; gap: 4px; }
.export-btn:hover { color: var(--hijau); border-color: var(--hijau); }
</style>
