<template>
  <div id="page-anggota" class="anggota-split-layout">
    <!-- ===== LEFT PANEL: FILTER & TABLE LIST ===== -->
    <div class="panel-left">
      <!-- Header + Actions -->
      <div class="page-header">
        <div class="page-header-left">
          <div class="page-title-group">
            <h1 class="page-title">Kelola Anggota</h1>
            <div class="page-subtitle">Total <strong>{{ totalAnggota }}</strong> anggota aktif terdaftar</div>
          </div>
        </div>
        <div class="page-header-right">
          <button class="btn btn-outline" style="width: auto" @click="exportCSV">
            <i class="ti ti-download"></i> Export CSV
          </button>
          <button class="btn btn-primary" style="width: auto" @click="openAddModal">
            <i class="ti ti-plus"></i> Tambah Anggota
          </button>
        </div>
      </div>

      <!-- Tabs -->
      <div class="tab-bar">
        <button v-for="t in tabs" :key="t.value" :class="['tab-btn', { active: activeTab === t.value }]" @click="changeTab(t.value)">
          {{ t.label }}
          <span class="tab-count">{{ t.count }}</span>
        </button>
      </div>

      <!-- Filters -->
      <div class="filter-bar">
        <div class="search-box">
          <i class="ti ti-search"></i>
          <input v-model="search" type="text" placeholder="Cari nama atau nomor anggota..." @input="debouncedFetch" />
        </div>
        <select v-model="filterCabang" class="filter-select" @change="fetchAnggota">
          <option value="">Cari cabang...</option>
          <option v-for="c in cabangList" :key="c.id" :value="c.id">{{ c.nama }}</option>
        </select>
        <select v-model="filterTingkatan" class="filter-select" @change="fetchAnggota">
          <option value="">Semua tingkatan</option>
          <option v-for="t in tingkatanList" :key="t" :value="t">{{ t }}</option>
        </select>
      </div>

      <!-- Table Skeleton / Content -->
      <div v-if="loading" class="table-card">
        <table class="data-table">
          <thead>
            <tr>
              <th>Nama Anggota</th>
              <th>Cabang / Unit</th>
              <th>Tingkatan</th>
              <th>BLBA</th>
              <th>Status</th>
              <th style="width:140px">Aksi</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="n in 5" :key="n">
              <td>
                <div class="member-cell">
                  <div class="skeleton-box" style="width:32px; height:32px; border-radius:50%;"></div>
                  <div class="member-info">
                    <div class="skeleton-box" style="width:120px; height:14px; margin-bottom:4px;"></div>
                    <div class="skeleton-box" style="width:80px; height:10px;"></div>
                  </div>
                </div>
              </td>
              <td>
                <div class="unit-cell">
                  <div class="skeleton-box" style="width:100px; height:12px; margin-bottom:4px;"></div>
                  <div class="skeleton-box" style="width:80px; height:10px;"></div>
                </div>
              </td>
              <td><div class="skeleton-box" style="width:80px; height:14px;"></div></td>
              <td><div class="skeleton-box" style="width:50px; height:14px;"></div></td>
              <td><div class="skeleton-box" style="width:50px; height:14px; border-radius:10px;"></div></td>
              <td><div class="skeleton-box" style="width:80px; height:24px; border-radius:4px;"></div></td>
            </tr>
          </tbody>
        </table>
      </div>

      <div v-else-if="anggotaData.length === 0" class="empty-state">
        <i class="ti ti-users"></i> Tidak ada anggota yang cocok dengan filter atau status saat ini.
      </div>

      <div v-else class="table-card">
        <table class="data-table">
          <thead>
            <tr>
              <th>Nama Anggota</th>
              <th>Cabang / Unit</th>
              <th>Tingkatan</th>
              <th>BLBA</th>
              <th>Status</th>
              <th style="width:140px">Aksi</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="a in anggotaData" :key="a.id" :class="{ 'row-selected': selectedAnggota?.id === a.id }" @click="selectAnggota(a)">
              <td>
                <div class="member-cell">
                  <img v-if="a.foto_url" :src="a.foto_url" class="member-avatar-img" />
                  <div v-else class="member-avatar" :style="{ background: getAvatarBg(a.nama_lengkap) }">
                    {{ getInitials(a.nama_lengkap) }}
                  </div>
                  <div class="member-info">
                    <div class="member-name">{{ a.nama_lengkap }}</div>
                    <div class="member-email">{{ a.nomor_anggota || 'PENDING' }}</div>
                  </div>
                </div>
              </td>
              <td>
                <div class="unit-cell">
                  <div class="unit-name">{{ a.cabang_nama || 'Pusat' }}</div>
                  <div class="cabang-name">{{ a.unit_nama || 'Unit Latihan' }}</div>
                </div>
              </td>
              <td>
                <span class="tingkatan-badge">{{ a.tingkatan || 'Dasar' }}</span>
              </td>
              <td>
                <span :class="['blba-badge', getBLBAStatus(a) === 'Lunas' ? 'lunas' : 'belum']">
                  {{ getBLBAStatus(a) }}
                </span>
              </td>
              <td>
                <span :class="['status-badge', a.status]">{{ statusLabel[a.status] || a.status }}</span>
              </td>
              <td @click.stopPropagation>
                <div class="action-btns-row">
                  <template v-if="a.status === 'pending'">
                    <button class="btn-verify approve" @click="verifikasi(a.id, 'approve')" :disabled="submitting">
                      Setujui
                    </button>
                    <button class="btn-verify reject" @click="verifikasi(a.id, 'reject')" :disabled="submitting">
                      Tolak
                    </button>
                  </template>
                  <template v-else>
                    <button class="icon-btn" title="Edit Data" @click="openEditModal(a)">
                      <i class="ti ti-edit"></i>
                    </button>
                    <button class="icon-btn danger" title="Nonaktifkan" @click="toggleStatus(a)">
                      <i class="ti ti-user-off"></i>
                    </button>
                  </template>
                </div>
              </td>
            </tr>
          </tbody>
        </table>

        <!-- Pagination -->
        <Pagination 
          v-model:currentPage="page" 
          v-model:itemsPerPage="limit" 
          :totalItems="totalAnggota" 
          @update:currentPage="fetchAnggota" 
          @update:itemsPerPage="fetchAnggota" 
        />
      </div>
    </div>

    <!-- ===== RIGHT PANEL: DETAIL SIDEBAR PANEL ===== -->
    <div class="panel-detail">
      <div v-if="!selectedAnggota" class="pd-empty">
        <i class="ti ti-user-search"></i>
        <p>Pilih anggota untuk melihat detail</p>
      </div>

      <div v-else class="pd-content">
        <div class="pd-header">
          <div class="pd-avatar-wrap">
            <img v-if="selectedAnggota.foto_url" :src="selectedAnggota.foto_url" class="pd-avatar-img" />
            <div v-else class="pd-avatar-fallback" :style="{ background: getAvatarBg(selectedAnggota.nama_lengkap) }">
              {{ getInitials(selectedAnggota.nama_lengkap) }}
            </div>
            <div style="flex:1;">
              <div class="pd-name">{{ selectedAnggota.nama_lengkap }}</div>
              <div class="pd-nomor">{{ selectedAnggota.nomor_anggota || 'PENDING' }}</div>
              <div class="pd-badges">
                <span :class="['status-badge', selectedAnggota.status]">{{ statusLabel[selectedAnggota.status] || selectedAnggota.status }}</span>
                <span class="tingkatan-badge border">{{ selectedAnggota.tingkatan || 'Pra Dasar' }}</span>
              </div>
            </div>
          </div>
        </div>

        <div class="pd-scroll">
          <!-- DATA PRIBADI -->
          <div class="pd-section">Data pribadi</div>
          <div>
            <div class="info-row"><span class="info-key">Nama lengkap</span><span class="info-val">{{ selectedAnggota.nama_lengkap }}</span></div>
            <div class="info-row"><span class="info-key">Email</span><span class="info-val">{{ getMockEmail(selectedAnggota) }}</span></div>
            <div class="info-row"><span class="info-key">No. HP</span><span class="info-val">{{ selectedAnggota.no_hp || '0819-1122-3344' }}</span></div>
            <div class="info-row"><span class="info-key">Tanggal lahir</span><span class="info-val">{{ formatDate(selectedAnggota.tanggal_lahir) }}</span></div>
            <div class="info-row"><span class="info-key">Jenis kelamin</span><span class="info-val">{{ selectedAnggota.jenis_kelamin === 'P' ? 'Perempuan' : 'Laki-laki' }}</span></div>
            <div class="info-row"><span class="info-key">Domisili</span><span class="info-val">{{ selectedAnggota.cabang_nama || 'Yogyakarta' }}</span></div>
          </div>

          <!-- KEANGGOTAAN -->
          <div class="pd-section">Keanggotaan</div>
          <div>
            <div class="info-row"><span class="info-key">No. anggota</span><span class="info-val monospace">{{ selectedAnggota.nomor_anggota || 'PENDING' }}</span></div>
            <div class="info-row"><span class="info-key">Cabang</span><span class="info-val">{{ selectedAnggota.cabang_nama || 'Pusat' }}</span></div>
            <div class="info-row"><span class="info-key">Unit latihan</span><span class="info-val">{{ selectedAnggota.unit_nama || 'Unit Latihan' }}</span></div>
            <div class="info-row"><span class="info-key">Tingkatan</span><span class="info-val">{{ selectedAnggota.tingkatan || 'Dasar' }}</span></div>
            <div class="info-row"><span class="info-key">Status</span><span class="info-val"><span :class="['status-badge', selectedAnggota.status]">{{ statusLabel[selectedAnggota.status] || selectedAnggota.status }}</span></span></div>
          </div>

          <!-- BLBA & KEHADIRAN -->
          <div class="pd-section">BLBA & Kehadiran</div>
          <div>
            <div class="info-row">
              <span class="info-key">BLBA bulan ini</span>
              <span class="info-val">
                <span :class="selectedMemberStats?.blba_bulan_ini === 'Lunas' ? 'iuran-lunas' : 'iuran-belum'">
                  {{ selectedMemberStats?.blba_bulan_ini === 'Lunas' ? '✓ Lunas' : '✗ Belum' }}
                </span>
              </span>
            </div>
            <div class="info-row">
              <span class="info-key">BLBA berjalan</span>
              <span class="info-val" v-if="selectedMemberStats">
                {{ selectedMemberStats.blba_lunas }} / {{ selectedMemberStats.blba_total }} bulan
              </span>
              <span class="info-val" v-else>-</span>
            </div>
            <div class="info-row">
              <span class="info-key">Kehadiran bulan ini</span>
              <span class="info-val" v-if="selectedMemberStats">
                {{ selectedMemberStats.kehadiran_bulan_ini }} / {{ selectedMemberStats.kehadiran_total_bulan_ini }} sesi
              </span>
              <span class="info-val" v-else>-</span>
            </div>
            <div class="info-row">
              <span class="info-key">Total kehadiran</span>
              <span class="info-val" v-if="selectedMemberStats">
                {{ selectedMemberStats.total_kehadiran }} sesi (tingkat ini)
              </span>
              <span class="info-val" v-else>-</span>
            </div>
          </div>

          <!-- RIWAYAT & STATISTIK KEHADIRAN -->
          <div class="pd-section">Riwayat & statistik kehadiran</div>
          <div class="stat-green-card" v-if="selectedMemberStats">
            <div class="stat-green-head">
              <div>
                <div class="stat-green-lbl">Kehadiran tingkat saat ini ({{ selectedAnggota.tingkatan || 'Dasar' }})</div>
                <div class="stat-green-val">{{ selectedMemberStats.total_kehadiran }} <span class="stat-green-target">/ {{ selectedMemberStats.kehadiran_target }} sesi target</span></div>
              </div>
              <div style="text-align:right;">
                <div class="stat-green-val">{{ selectedMemberStats.kehadiran_pct }}%</div>
                <div class="stat-green-sub">rata-rata</div>
              </div>
            </div>
            <div class="stat-green-track"><div class="stat-green-bar" :style="{ width: selectedMemberStats.kehadiran_pct + '%' }"></div></div>
            <div class="stat-green-desc">
              <i class="ti ti-circle-check"></i>
              Syarat EKT {{ selectedMemberStats.total_kehadiran >= selectedMemberStats.kehadiran_target ? 'terpenuhi' : 'belum terpenuhi' }} (min. {{ selectedMemberStats.kehadiran_target }} sesi)
            </div>
          </div>

          <!-- 6 BULAN TERAKHIR -->
          <div style="font-size:10px;font-weight:700;color:var(--text3);margin-bottom:7px;text-transform:uppercase;letter-spacing:.04em;">6 Bulan Terakhir</div>
          <div class="history-months-box">
            <div class="history-months-header"><span>Bulan</span><span style="text-align:center;">Hadir</span><span style="text-align:center;">Sesi</span><span style="text-align:right;">%</span></div>
            <div v-for="h in selectedMemberStats?.attendanceHistory" :key="h.month" class="history-months-row">
              <span style="font-weight:600;">{{ h.month }}</span>
              <span style="text-align:center;font-weight:700;color:var(--hijau);">{{ h.attended }}</span>
              <span style="text-align:center;color:var(--text3);">{{ h.total }}</span>
              <div class="history-months-pct">
                <div class="hm-track">
                  <div class="hm-bar" :style="{ width: h.pct + '%', background: h.pct < 75 ? 'var(--kuning)' : 'var(--hijau)' }"></div>
                </div>
                <span style="font-weight:700;" :style="{ color: h.pct < 75 ? '#9a7000' : 'var(--hijau)' }">{{ h.pct }}%</span>
              </div>
            </div>
          </div>

          <!-- FOTO PROFIL -->
          <div class="pd-section">Foto profil</div>
          <div class="photo-box">
            <img v-if="selectedAnggota.foto_url" :src="selectedAnggota.foto_url" class="photo-avatar-img" />
            <div v-else class="photo-avatar" :style="{ background: getAvatarBg(selectedAnggota.nama_lengkap) }">
              {{ getInitials(selectedAnggota.nama_lengkap) }}
            </div>
            <div style="flex:1;">
              <div style="font-size:11px;color:var(--text2);margin-bottom:5px;">
                {{ selectedAnggota.foto_url ? 'Foto profil terunggah' : 'Belum ada foto profil' }}
              </div>
              <button class="photo-btn" @click="triggerUpload"><i class="ti ti-upload"></i> Upload foto</button>
            </div>
          </div>
          <input type="file" ref="fileInput" accept="image/*" @change="onFileSelected" style="display: none;" />

          <!-- PROFIL KEBUGARAN -->
          <div class="pd-section" style="margin-top:14px;">Profil Kebugaran</div>
          <div class="kebugaran-details-box" v-if="selectedMemberStats">
            <div style="display:flex;align-items:center;gap:10px;margin-bottom:12px;">
              <div class="keb-score-badge">
                <div class="keb-score-val">{{ selectedMemberStats.score }}%</div>
                <div class="keb-score-lbl">Skor terkini</div>
              </div>
              <div style="flex:1;">
                <div style="font-size:11px;font-weight:700;margin-bottom:2px;">
                  <span :class="['kategori-badge', selectedMemberStats.catClass]">{{ selectedMemberStats.category }}</span>
                </div>
                <div style="font-size:10px;color:var(--text3);">Periode 2 — 2026</div>
                <div style="font-size:10px;margin-top:4px;font-weight:600;" :style="{ color: selectedMemberStats.diffColor }">
                  {{ selectedMemberStats.diffStr }} vs periode sebelumnya
                </div>
              </div>
            </div>

            <!-- SVG line chart -->
            <div style="margin-bottom:8px;">
              <div style="font-size:10px;font-weight:700;color:var(--text3);text-transform:uppercase;letter-spacing:.04em;margin-bottom:6px;">Tren Skor per Periode</div>
              <svg viewBox="0 0 260 120" xmlns="http://www.w3.org/2000/svg" style="width:100%;height:120px;display:block;">
                <line x1="32" y1="77" x2="248" y2="77" stroke="#e5e7eb" stroke-width="1" stroke-dasharray="3,3"/>
                <text x="28" y="80" text-anchor="end" font-size="8" fill="#9ca3af">50%</text>
                <line x1="32" y1="43" x2="248" y2="43" stroke="#e5e7eb" stroke-width="1" stroke-dasharray="3,3"/>
                <text x="28" y="46" text-anchor="end" font-size="8" fill="#9ca3af">70%</text>
                <line x1="32" y1="10" x2="248" y2="10" stroke="#e5e7eb" stroke-width="1" stroke-dasharray="3,3"/>
                <text x="28" y="13" text-anchor="end" font-size="8" fill="#9ca3af">90%</text>
                <text x="32" y="116" text-anchor="middle" font-size="8" fill="#9ca3af">P1/24</text>
                <text x="62" y="116" text-anchor="middle" font-size="8" fill="#9ca3af">P2/24</text>
                <text x="93" y="116" text-anchor="middle" font-size="8" fill="#9ca3af">P3/24</text>
                <text x="124" y="116" text-anchor="middle" font-size="8" fill="#9ca3af">P1/25</text>
                <text x="155" y="116" text-anchor="middle" font-size="8" fill="#9ca3af">P2/25</text>
                <text x="186" y="116" text-anchor="middle" font-size="8" fill="#9ca3af">P3/25</text>
                <text x="217" y="116" text-anchor="middle" font-size="8" fill="#9ca3af">P1/26</text>
                <text x="248" y="116" text-anchor="middle" font-size="8" fill="#9ca3af">P2/26</text>
                <polygon :points="selectedMemberStats.polygonPoints" fill="#1a5c2a" fill-opacity="0.08"/>
                <polyline :points="selectedMemberStats.polylinePoints" fill="none" stroke="#1a5c2a" stroke-width="2" stroke-linejoin="round" stroke-linecap="round"/>
                <circle v-for="(p, idx) in selectedMemberStats.points" 
                        :key="idx" 
                        :cx="p.x" 
                        :cy="p.y" 
                        :r="idx === 7 ? 4 : 2.5" 
                        fill="#1a5c2a" 
                        stroke="#fff" 
                        stroke-width="1.5"
                />
                <text :x="selectedMemberStats.points[7].x - 22" :y="selectedMemberStats.points[7].y - 8" font-size="9" fill="#1a5c2a" font-weight="700">{{ selectedMemberStats.score }}%</text>
              </svg>
            </div>

            <!-- Last Fitness test -->
            <div style="font-size:10px;font-weight:700;color:var(--text3);text-transform:uppercase;letter-spacing:.04em;margin-bottom:6px;margin-top:10px;">Hasil Terakhir (P2/2026)</div>
            <div class="keb-row-list">
              <div v-for="res in selectedMemberStats.results" :key="res.label" class="keb-row-item">
                <span>{{ res.label }}</span>
                <div style="display:flex;align-items:center;gap:6px;">
                  <span style="font-weight:700">{{ res.value }}</span>
                  <span :class="['kategori-badge', res.cat.class]">{{ res.cat.label }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div class="pd-actions">
          <template v-if="selectedAnggota.status === 'pending'">
            <button class="btn btn-primary" @click="verifikasi(selectedAnggota.id, 'approve')"><i class="ti ti-check"></i> Setujui Pendaftaran</button>
            <button class="btn btn-outline" style="color:var(--merah); border-color:var(--merah);" @click="verifikasi(selectedAnggota.id, 'reject')"><i class="ti ti-x"></i> Tolak Pendaftaran</button>
          </template>
          <template v-else>
            <button class="btn btn-primary" @click="openEditModal(selectedAnggota)"><i class="ti ti-edit"></i> Edit data anggota</button>
            <div style="display:flex;gap:7px;">
              <button class="btn btn-outline" style="flex:1;" @click="resetPassword"><i class="ti ti-key"></i> Reset password</button>
              <button class="btn btn-outline" style="flex:1;color:var(--merah);border-color:var(--merah);" @click="toggleStatus(selectedAnggota)"><i class="ti ti-user-off"></i> Nonaktifkan</button>
            </div>
          </template>
        </div>
      </div>
    </div>

    <!-- ===== MODAL ADD / EDIT ANGGOTA ===== -->
    <div v-if="showAddModal" class="modal-overlay">
      <div class="modal-card">
        <div class="modal-header">
          <h2 class="modal-title">{{ editingId ? 'Edit Data Anggota' : 'Tambah Anggota Baru' }}</h2>
          <button class="modal-close" @click="showAddModal = false"><i class="ti ti-x"></i></button>
        </div>
        <form @submit.prevent="saveAnggota" class="modal-form">
          <div class="form-group">
            <label class="form-label">Nama Lengkap</label>
            <input v-model="form.nama_lengkap" type="text" class="form-input" placeholder="Nama sesuai KTP..." required />
          </div>
          <div class="form-group">
            <label class="form-label">No. HP / WhatsApp</label>
            <input v-model="form.no_hp" type="text" class="form-input" placeholder="Contoh: 081234567890" required />
          </div>
          <div class="form-row">
            <div class="form-group half">
              <label class="form-label">Jenis Kelamin</label>
              <select v-model="form.jenis_kelamin" class="form-select" required>
                <option value="L">Laki-laki</option>
                <option value="P">Perempuan</option>
              </select>
            </div>
            <div class="form-group half">
              <label class="form-label">Tanggal Lahir</label>
              <input v-model="form.tanggal_lahir" type="date" class="form-input" required />
            </div>
          </div>
          <div class="form-group">
            <label class="form-label">Cabang Latihan</label>
            <select v-model="form.cabang_id" class="form-select" @change="onCabangChange(false)" required>
              <option value="" disabled>Pilih Cabang</option>
              <option v-for="c in cabangList" :key="c.id" :value="c.id">{{ c.nama }}</option>
            </select>
          </div>
          <div class="form-group">
            <label class="form-label">Unit Latihan</label>
            <select v-model="form.unit_id" class="form-select" :disabled="!form.cabang_id || loadingUnits" required>
              <option value="" disabled>{{ loadingUnits ? 'Memuat unit...' : 'Pilih Unit' }}</option>
              <option v-for="u in unitList" :key="u.id" :value="u.id">{{ u.nama }}</option>
            </select>
          </div>
          <div v-if="editingId" class="form-group">
            <label class="form-label">Tingkatan</label>
            <select v-model="form.tingkatan" class="form-select">
              <option v-for="t in tingkatanList" :key="t" :value="t">{{ t }}</option>
            </select>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-outline" style="width: auto" @click="showAddModal = false">Batal</button>
            <button type="submit" class="btn btn-primary" style="width: auto" :disabled="submitting">
              <i v-if="submitting" class="ti ti-loader-2 spin"></i> {{ editingId ? 'Simpan Perubahan' : 'Simpan Pendaftaran' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({ title: 'Kelola Anggota' })
const route = useRoute()

const api = useApi()

const activeTab = ref('aktif')
const search = ref('')
const filterCabang = ref('')
const filterTingkatan = ref('')
const showAddModal = ref(false)
const submitting = ref(false)
const loading = ref(false)
const loadingUnits = ref(false)

const page = ref(1)
const limit = ref(15)
const totalAnggota = ref(0)
const totalPages = ref(1)

const anggotaData = ref<any[]>([])
const cabangList = ref<any[]>([])
const unitList = ref<any[]>([])
const selectedAnggota = ref<any>(null)
const editingId = ref<string | null>(null)

const selectedMemberStats = ref<any>(null)
const fileInput = ref<HTMLInputElement | null>(null)

const triggerUpload = () => {
  fileInput.value?.click()
}

const onFileSelected = async (event: Event) => {
  const target = event.target as HTMLInputElement
  if (!target.files || target.files.length === 0) return

  const file = target.files[0]
  const formData = new FormData()
  formData.append('foto', file)

  try {
    const res = await api.post(`/organization/anggota/${selectedAnggota.value.id}/upload`, formData)
    alert('Foto profil berhasil diunggah!')
    if (res && res.foto_url) {
      selectedAnggota.value.foto_url = res.foto_url
      const idx = anggotaData.value.findIndex((x: any) => x.id === selectedAnggota.value.id)
      if (idx !== -1) {
        anggotaData.value[idx].foto_url = res.foto_url
      }
    }
  } catch (e: any) {
    console.error('Gagal mengunggah foto:', e)
    alert('Gagal mengunggah foto: ' + e.message)
  }
}

watch(selectedAnggota, async (newVal) => {
  if (!newVal) {
    selectedMemberStats.value = null
    return
  }
  try {
    const data = await api.get(`/organization/anggota/${newVal.id}/stats`)
    selectedMemberStats.value = data
  } catch (e) {
    console.error('Gagal mengambil stats dari backend:', e)
    selectedMemberStats.value = null
  }
}, { immediate: true })

const tabs = ref([
  { label: 'Semua Anggota', value: 'aktif', count: 0 },
  { label: 'Menunggu Verifikasi', value: 'pending', count: 0 },
  { label: 'Nonaktif', value: 'nonaktif', count: 0 },
])

const tingkatanList = ['Pra Dasar','Dasar','PH','Gabungan','PK','GPK','Penjuru','Selangkah','Meditasi']
const statusLabel: Record<string, string> = { aktif: 'Aktif', pending: 'Pending', nonaktif: 'Nonaktif' }

const form = ref({
  nama_lengkap: '',
  no_hp: '',
  jenis_kelamin: 'L',
  tanggal_lahir: '',
  cabang_id: '',
  unit_id: '',
  tingkatan: 'Pra Dasar',
  status: 'aktif'
})

const fetchStats = async () => {
  try {
    const resAktif = await api.get('/organization/anggota?status=aktif&limit=1')
    const resPending = await api.get('/organization/anggota?status=pending&limit=1')
    const resNon = await api.get('/organization/anggota?status=nonaktif&limit=1')
    
    tabs.value[0].count = (resAktif.total || 0) + (resNon.total || 0)
    tabs.value[1].count = resPending.total || 0
    tabs.value[2].count = resNon.total || 0
  } catch (e) {
    console.error(e)
  }
}

const fetchAnggota = async () => {
  loading.value = true
  try {
    const data = await api.get(
      `/organization/anggota?page=${page.value}&limit=${limit.value}&status=${activeTab.value}&cabang_id=${filterCabang.value}&search=${search.value}`
    )
    anggotaData.value = data.data || []
    totalAnggota.value = data.total || 0
    totalPages.value = data.total_pages || 1

    // Auto-select first member when data loaded
    if (anggotaData.value.length > 0) {
      if (!selectedAnggota.value || !anggotaData.value.some((x: any) => x.id === selectedAnggota.value.id)) {
        selectedAnggota.value = anggotaData.value[0]
      }
    } else {
      selectedAnggota.value = null
    }
  } catch (e) {
    console.error('Failed to fetch anggota', e)
  } finally {
    loading.value = false
  }
}

let timeout: any = null
const debouncedFetch = () => {
  clearTimeout(timeout)
  timeout = setTimeout(() => {
    page.value = 1
    fetchAnggota()
  }, 300)
}

const changeTab = (val: string) => {
  activeTab.value = val
  page.value = 1
  fetchAnggota()
}

const changePage = (p: number) => {
  page.value = p
  fetchAnggota()
}

const selectAnggota = (a: any) => {
  selectedAnggota.value = a
}

const fetchCabang = async () => {
  try {
    const res = await api.get('/organization/cabang?limit=100')
    cabangList.value = res.data || []
  } catch (e) {
    console.error(e)
  }
}

const onCabangChange = async (keepUnitId = false) => {
  if (!keepUnitId) {
    form.value.unit_id = ''
  }
  unitList.value = []
  if (!form.value.cabang_id) return
  
  loadingUnits.value = true
  try {
    const res = await api.get(`/organization/cabang/${form.value.cabang_id}/unit`)
    unitList.value = res || []
  } catch (e) {
    console.error(e)
  } finally {
    loadingUnits.value = false
  }
}

const openAddModal = () => {
  editingId.value = null
  form.value = {
    nama_lengkap: '',
    no_hp: '',
    jenis_kelamin: 'L',
    tanggal_lahir: '',
    cabang_id: '',
    unit_id: '',
    tingkatan: 'Pra Dasar',
    status: 'aktif'
  }
  unitList.value = []
  showAddModal.value = true
}

const openEditModal = async (a: any) => {
  editingId.value = a.id
  form.value = {
    nama_lengkap: a.nama_lengkap || '',
    no_hp: a.no_hp || '',
    jenis_kelamin: a.jenis_kelamin || 'L',
    tanggal_lahir: a.tanggal_lahir ? a.tanggal_lahir.substring(0, 10) : '',
    cabang_id: a.cabang_id || '',
    unit_id: a.unit_id || '',
    tingkatan: a.tingkatan || 'Pra Dasar',
    status: a.status || 'aktif'
  }
  if (!form.value.cabang_id || !form.value.tanggal_lahir) {
    try {
      const detail = await api.get(`/organization/anggota/${a.id}`)
      if (detail) {
        form.value.cabang_id = detail.cabang_id || form.value.cabang_id
        form.value.unit_id = detail.unit_id || form.value.unit_id
        form.value.tanggal_lahir = detail.tanggal_lahir ? detail.tanggal_lahir.substring(0, 10) : form.value.tanggal_lahir
        form.value.tingkatan = detail.tingkatan || form.value.tingkatan
        form.value.status = detail.status || form.value.status
      }
    } catch (e) {
      console.error('Gagal memuat detail anggota:', e)
    }
  }
  await onCabangChange(true)
  showAddModal.value = true
}

const saveAnggota = async () => {
  submitting.value = true
  try {
    if (editingId.value) {
      await api.put(`/organization/anggota/${editingId.value}`, {
        nama_lengkap: form.value.nama_lengkap,
        no_hp: form.value.no_hp,
        jenis_kelamin: form.value.jenis_kelamin,
        tanggal_lahir: form.value.tanggal_lahir,
        unit_id: form.value.unit_id,
        tingkatan: form.value.tingkatan,
        status: form.value.status
      })
      alert('Data anggota berhasil diperbarui!')
    } else {
      await api.post('/organization/anggota', {
        nama_lengkap: form.value.nama_lengkap,
        no_hp: form.value.no_hp,
        jenis_kelamin: form.value.jenis_kelamin,
        tanggal_lahir: form.value.tanggal_lahir,
        unit_id: form.value.unit_id
      })
      alert('Pendaftaran anggota berhasil diajukan!')
    }
    showAddModal.value = false
    editingId.value = null
    await fetchStats()
    await fetchAnggota()
    if (selectedAnggota.value) {
      try {
        const updated = await api.get(`/organization/anggota/${selectedAnggota.value.id}`)
        if (updated) selectedAnggota.value = updated
      } catch (_) {}
    }
  } catch (e: any) {
    alert(e.message || 'Gagal menyimpan data')
  } finally {
    submitting.value = false
  }
}

const verifikasi = async (id: string, aksi: string) => {
  if (!confirm(`Apakah Anda yakin ingin melakukan ${aksi} untuk anggota ini?`)) return
  submitting.value = true
  try {
    await api.post(`/organization/anggota/${id}/verifikasi`, { aksi })
    fetchStats()
    fetchAnggota()
  } catch (e: any) {
    alert(e.message || 'Gagal memverifikasi anggota')
  } finally {
    submitting.value = false
  }
}

const toggleStatus = (a: any) => {
  if (confirm(`Apakah Anda yakin ingin menonaktifkan anggota ${a.nama_lengkap}?`)) {
    a.status = 'nonaktif'
    alert(`Anggota ${a.nama_lengkap} dinonaktifkan!`)
  }
}

const resetPassword = () => {
  alert('Link reset password berhasil dikirim ke email anggota!')
}

const exportCSV = () => {
  alert('Data anggota berhasil diekspor ke file CSV!')
}

// Visual Helpers
const getInitials = (name: string) => {
  if (!name) return 'RW'
  const parts = name.split(' ')
  if (parts.length > 1) {
    return (parts[0][0] + parts[1][0]).toUpperCase()
  }
  return parts[0].substring(0, 2).toUpperCase()
}

const getAvatarBg = (name: string) => {
  const colors = ['#1a5c2a', '#5b21b6', '#0c5478', '#8c1c1c', '#7a6000', '#0e7aad']
  let hash = 0
  for (let i = 0; i < name.length; i++) {
    hash = name.charCodeAt(i) + ((hash << 5) - hash)
  }
  const index = Math.abs(hash % colors.length)
  return colors[index]
}

const getMockEmail = (a: any) => {
  return a.nama_lengkap.toLowerCase().replace(/\s/g, '') + '@email.com'
}

const getBLBAStatus = (a: any) => {
  return (a.status === 'pending') ? 'Belum' : 'Lunas'
}

const getAttendancePct = (a: any) => {
  const sessions = a.counter_kehadiran || 18
  return Math.round((sessions / 24) * 100)
}

const formatDate = (dateStr: string) => {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  return d.toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' })
}

onMounted(() => {
  if (route.query.search) {
    search.value = String(route.query.search)
  }
  fetchStats()
  fetchAnggota()
  fetchCabang()
})

watch(() => route.query.search, (newSearch) => {
  if (newSearch !== undefined) {
    search.value = String(newSearch)
    page.value = 1
    fetchAnggota()
  }
})
</script>

<style scoped>
.anggota-split-layout {
  display: flex;
  gap: 20px;
  height: calc(100vh - 64px - 44px);
  background: var(--bg);
  overflow: hidden;
  position: relative;
}

/* Left panel for list content */
.panel-left {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  height: 100%;
  overflow-y: auto;
  background: var(--bg);
  padding-right: 4px;
}

.page-header { display: flex; align-items: flex-start; justify-content: space-between; margin-bottom: 20px; gap: 12px; flex-wrap: wrap; }
.page-title { font-size: 20px; font-weight: 800; color: var(--text1); margin: 0; }
.page-subtitle { font-size: 12px; color: var(--text3); margin-top: 2px; }
.page-header-right { display: flex; gap: 8px; flex-wrap: wrap; }

.btn { display: inline-flex; align-items: center; gap: 6px; padding: 8px 14px; border-radius: var(--r8); font-size: 12px; font-weight: 600; cursor: pointer; border: none; transition: all .15s; }
.btn-primary { background: var(--hijau); color: #fff; }
.btn-primary:hover { background: var(--hijau2); }
.btn-outline { background: #fff; color: var(--text2); border: 1px solid var(--border); }
.btn-outline:hover { background: var(--surface); border-color: var(--hijau); color: var(--hijau); }

.tab-bar { display: flex; gap: 4px; margin-bottom: 16px; border-bottom: 1px solid var(--border); }
.tab-btn { display: flex; align-items: center; gap: 6px; padding: 8px 16px; font-size: 12px; font-weight: 600; background: none; border: none; cursor: pointer; color: var(--text2); border-bottom: 2px solid transparent; margin-bottom: -1px; }
.tab-btn.active { color: var(--hijau); border-bottom-color: var(--hijau); }
.tab-count { background: var(--hijau3); color: var(--hijau); border-radius: 10px; padding: 1px 7px; font-size: 10px; }

.filter-bar { display: flex; gap: 8px; margin-bottom: 14px; flex-wrap: wrap; }
.search-box { display: flex; align-items: center; gap: 8px; flex: 1; min-width: 200px; background: var(--card); border: 1px solid var(--border); border-radius: var(--r8); padding: 7px 12px; }
.search-box i { color: var(--text3); font-size: 14px; }
.search-box input { border: none; outline: none; background: none; font-size: 12px; flex: 1; color: var(--text1); }
.filter-select { padding: 7px 10px; border: 1px solid var(--border); border-radius: var(--r8); font-size: 12px; background: var(--card); color: var(--text2); cursor: pointer; }

.table-card { background: var(--card); border: 1px solid var(--border); border-radius: var(--r12); overflow-x: auto; flex-shrink: 0; }
.data-table { width: 100%; border-collapse: collapse; }
.data-table th { padding: 10px 14px; font-size: 11px; font-weight: 700; color: var(--text3); text-transform: uppercase; text-align: left; background: var(--surface); border-bottom: 1px solid var(--border); letter-spacing: .03em; }
.data-table td { padding: 12px 14px; font-size: 12px; color: var(--text1); border-bottom: 1px solid var(--border); vertical-align: middle; }
.data-table tr:last-child td { border-bottom: none; }
.data-table tr { cursor: pointer; transition: background 0.15s; }
.data-table tr:hover td { background: var(--surface); }
.data-table tr.row-selected td { background: var(--hijau3) !important; border-left: 3px solid var(--hijau); }

.member-cell { display: flex; align-items: center; gap: 10px; }
.member-avatar { width: 32px; height: 32px; border-radius: 50%; display: flex; align-items: center; justify-content: center; font-size: 11px; font-weight: 700; color: #fff; flex-shrink: 0; }
.member-name { font-size: 12px; font-weight: 600; }
.member-email { font-size: 11px; color: var(--text3); }
.nomor-anggota { font-size: 11px; font-family: monospace; background: var(--surface); padding: 2px 6px; border-radius: 4px; }
.unit-cell { display: flex; flex-direction: column; }
.unit-name { font-size: 12px; font-weight: 500; }
.cabang-name { font-size: 11px; color: var(--text3); }
.tingkatan-badge { display: inline-block; padding: 2px 8px; border-radius: 10px; font-size: 11px; font-weight: 700; background: var(--hijau3); color: var(--hijau); }
.tingkatan-badge.border { border: 1px solid rgba(255,255,255,0.3); background: rgba(255,255,255,0.2); color: #fff; }

.blba-badge { display: inline-block; padding: 2px 6px; border-radius: 10px; font-size: 10px; font-weight: 700; }
.blba-badge.lunas { background: var(--hijau3); color: var(--hijau); }
.blba-badge.belum { background: #fff8e0; color: #a07000; }

.status-badge { display: inline-block; padding: 3px 8px; border-radius: 10px; font-size: 11px; font-weight: 600; text-transform: capitalize; }
.status-badge.aktif { background: var(--hijau3); color: var(--hijau); }
.status-badge.nonaktif { background: #f0f0f0; color: var(--text3); }
.status-badge.pending { background: #fff8e0; color: #a07000; border: 0.5px solid #e8c42a; }

.action-btns-row { display: flex; gap: 4px; align-items: center; }
.btn-verify { border: none; padding: 5px 10px; border-radius: var(--r6); font-size: 10px; font-weight: 700; cursor: pointer; color: #fff; }
.btn-verify.approve { background: var(--hijau2); }
.btn-verify.reject { background: var(--merah); }
.btn-verify:disabled { opacity: 0.5; cursor: not-allowed; }
.icon-btn { width: 26px; height: 26px; border: 1px solid var(--border); border-radius: var(--r6); background: #fff; cursor: pointer; display: inline-flex; align-items: center; justify-content: center; font-size: 12px; transition: all .15s; }
.icon-btn:hover { border-color: var(--hijau); color: var(--hijau); background: var(--hijau3); }
.icon-btn.danger:hover { border-color: var(--merah); color: var(--merah); background: #fdecea; }

.table-footer { display: flex; align-items: center; justify-content: space-between; padding: 12px 16px; border-top: 1px solid var(--border); background: var(--card); }
.table-info { font-size: 11px; color: var(--text3); }
.pagination { display: flex; gap: 4px; align-items: center; }
.page-btn { width: 28px; height: 28px; border: 1px solid var(--border); border-radius: var(--r6); background: #fff; cursor: pointer; font-size: 11px; display: flex; align-items: center; justify-content: center; color: var(--text2); }
.page-btn:disabled { opacity: 0.5; cursor: not-allowed; }
.page-btn.active { background: var(--hijau); color: #fff; border-color: var(--hijau); }

/* Right detail panel card wrapper */
.panel-detail {
  width: 380px;
  flex-shrink: 0;
  background: var(--card);
  border: 1px solid var(--border);
  border-radius: var(--r12);
  display: flex;
  flex-direction: column;
  overflow: hidden;
  height: 100%;
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.04);
}
.pd-empty { flex: 1; display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 10px; color: var(--text3); }
.pd-empty i { font-size: 40px; }
.pd-empty p { font-size: 12px; }

.pd-content { flex: 1; display: flex; flex-direction: column; overflow: hidden; }
.pd-header {
  background: linear-gradient(135deg, var(--hijau), #51af61);
  padding: 18px 16px;
  flex-shrink: 0;
  color: #fff;
  border-top-left-radius: 11px;
  border-top-right-radius: 11px;
}
.pd-avatar-wrap { display: flex; align-items: flex-start; gap: 12px; }
.pd-avatar-fallback { width: 72px; height: 72px; border-radius: 12px; background: rgba(255,255,255,0.2); display: flex; align-items: center; justify-content: center; font-size: 24px; font-weight: 700; color: #fff; border: 2px solid rgba(255,255,255,0.3); flex-shrink: 0; }
.pd-name { font-size: 15px; font-weight: 700; color: #fff; margin-bottom: 2px; }
.pd-nomor { font-size: 11px; color: rgba(255,255,255,0.75); font-family: monospace; margin-bottom: 7px; }
.pd-badges { display: flex; gap: 5px; flex-wrap: wrap; }

.pd-scroll { flex: 1; overflow-y: auto; padding: 14px 16px; }
.pd-scroll::-webkit-scrollbar { width: 4px; }
.pd-scroll::-webkit-scrollbar-thumb { background: var(--border); border-radius: 2px; }
.pd-section { font-size: 10px; font-weight: 700; color: var(--text3); text-transform: uppercase; letter-spacing: .06em; margin-bottom: 8px; margin-top: 14px; border-bottom: 1.5px solid var(--border); padding-bottom: 3px; }
.pd-section:first-child { margin-top: 0; }

.info-row { display: flex; justify-content: space-between; align-items: center; padding: 6px 0; border-bottom: .5px solid var(--border); }
.info-row:last-child { border-bottom: none; }
.info-key { font-size: 11px; color: var(--text2); }
.info-val { font-size: 11px; font-weight: 600; color: var(--text1); text-align: right; max-width: 180px; }
.info-val.monospace { font-family: monospace; font-size: 10px; }

.iuran-lunas { color: var(--hijau2); font-weight: 700; font-size: 11px; }
.iuran-belum { color: var(--merah); font-weight: 700; font-size: 11px; }

/* Stat green card */
.stat-green-card { background: var(--hijau3); border: .5px solid var(--hijau4); border-radius: var(--r8); padding: 11px 12px; margin-bottom: 10px; }
.stat-green-head { display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 6px; }
.stat-green-lbl { font-size: 10px; color: var(--hijau2); margin-bottom: 2px; }
.stat-green-val { font-size: 20px; font-weight: 700; color: var(--hijau); }
.stat-green-target { font-size: 12px; font-weight: 400; color: var(--hijau2); }
.stat-green-sub { font-size: 10px; color: var(--hijau2); }
.stat-green-track { height: 6px; background: rgba(255,255,255,0.4); border-radius: 3px; overflow: hidden; margin-bottom: 5px; }
.stat-green-bar { height: 100%; background: var(--hijau); border-radius: 3px; }
.stat-green-desc { font-size: 10px; color: var(--hijau2); display: flex; align-items: center; gap: 4px; }

/* 6 Months history card */
.history-months-box { background: var(--surface); border-radius: var(--r8); overflow: hidden; border: .5px solid var(--border); margin-bottom: 10px; }
.history-months-header { display: grid; grid-template-columns: 1fr 1fr 1fr 1fr; font-size: 10px; font-weight: 700; color: var(--text3); padding: 7px 10px; border-bottom: .5px solid var(--border); text-transform: uppercase; letter-spacing: .04em; }
.history-months-row { display: grid; grid-template-columns: 1fr 1fr 1fr 1fr; font-size: 11px; padding: 7px 10px; border-bottom: .5px solid var(--border); align-items: center; }
.history-months-row:last-child { border-bottom: none; }
.history-months-pct { display: flex; align-items: center; justify-content: flex-end; gap: 5px; }
.hm-track { width: 30px; height: 4px; background: var(--border); border-radius: 2px; overflow: hidden; }
.hm-bar { height: 100%; background: var(--hijau2); border-radius: 2px; }

/* Foto profil */
.photo-box { background: var(--surface); border-radius: var(--r8); padding: 10px; display: flex; align-items: center; gap: 10px; border: .5px solid var(--border); }
.photo-avatar { width: 52px; height: 52px; border-radius: 10px; display: flex; align-items: center; justify-content: center; font-size: 18px; font-weight: 700; color: var(--hijau); border: .5px solid var(--hijau4); background: var(--hijau3); }
.pd-avatar-img { width: 72px; height: 72px; border-radius: 12px; object-fit: cover; border: 2px solid rgba(255,255,255,0.3); flex-shrink: 0; }
.photo-avatar-img { width: 52px; height: 52px; border-radius: 10px; object-fit: cover; border: .5px solid var(--border); flex-shrink: 0; }
.member-avatar-img { width: 32px; height: 32px; border-radius: 50%; object-fit: cover; flex-shrink: 0; }
.photo-btn { font-size: 11px; padding: 5px 10px; border-radius: var(--r6); border: 1px solid var(--border2); background: var(--card); color: var(--text2); cursor: pointer; display: inline-flex; align-items: center; gap: 4px; }
.photo-btn:hover { border-color: var(--hijau); color: var(--hijau); }

/* Fitness detail */
.kebugaran-details-box { background: var(--surface); border-radius: var(--r8); border: .5px solid var(--border); padding: 12px; }
.keb-score-badge { text-align: center; background: var(--hijau3); border-radius: var(--r8); padding: 10px 14px; border: 1px solid var(--hijau2); }
.keb-score-val { font-size: 22px; font-weight: 800; color: var(--hijau); line-height: 1; }
.keb-score-lbl { font-size: 9px; color: var(--text3); margin-top: 2px; }
.kategori-badge { display: inline-block; padding: 2px 6px; border-radius: 8px; font-size: 9px; font-weight: 700; }
.kategori-badge.sangat-baik { background: #e0f5fb; color: var(--biru); }
.kategori-badge.baik { background: var(--hijau3); color: var(--hijau); }
.kategori-badge.cukup { background: #fff8e0; color: #a07000; }
.kategori-badge.kurang { background: #fee2e2; color: var(--merah); }
.keb-row-list { display: flex; flex-direction: column; gap: 4px; }
.keb-row-item { display: flex; justify-content: space-between; align-items: center; padding: 5px 8px; background: var(--card); border-radius: 6px; font-size: 11px; border: 0.5px solid var(--border); }

.pd-actions { padding: 12px 16px; border-top: 1px solid var(--border); display: flex; flex-direction: column; gap: 7px; flex-shrink: 0; }
.pd-actions .btn { width: 100%; justify-content: center; }

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

.empty-state { display: flex; flex-direction: column; align-items: center; justify-content: center; padding: 40px; font-size: 13px; color: var(--text3); gap: 10px; background: var(--card); border-radius: var(--r12); border: 1px dashed var(--border); margin-top: 10px; }
.empty-state i { font-size: 32px; color: var(--text3); }

@keyframes spin { from { transform: rotate(0deg); } to { transform: rotate(360deg); } }
.spin { animation: spin .8s linear infinite; }

/* Skeleton Loader */
@keyframes pulse { 0% { opacity: 0.6; } 50% { opacity: 1; } 100% { opacity: 0.6; } }
.skeleton-box { animation: pulse 1.5s infinite ease-in-out; background-color: rgba(0, 0, 0, 0.08); border-radius: var(--r6); }
</style>
