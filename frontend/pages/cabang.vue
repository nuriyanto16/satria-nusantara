<template>
  <div id="page-cabang" class="cabang-split-layout">
    <!-- ===== PANEL KIRI: DAFTAR CABANG ===== -->
    <div class="panel-left">
      <div class="pl-head">
        <div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:10px;">
          <div class="pl-title">Daftar Cabang</div>
          <button class="btn btn-primary btn-xs" @click="openCreateModal"><i class="ti ti-plus"></i> Tambah</button>
        </div>
        
        <!-- Filter Provinsi -->
        <div class="form-group-filter">
          <select v-model="filterProvinsi" class="filter-select">
            <option value="">Semua Provinsi</option>
            <option v-for="p in provinsiList" :key="p.id" :value="p.nama">{{ p.nama }}</option>
          </select>
        </div>

        <!-- Search Box -->
        <div class="search-box">
          <i class="ti ti-search"></i>
          <input v-model="search" placeholder="Cari nama cabang..." />
        </div>
      </div>

      <div class="pl-list">
        <div v-if="loadingCabangList" class="loading-state-list">
          <i class="ti ti-loader-2 spin"></i> Memuat cabang...
        </div>
        <div v-else-if="groupedCabang.length === 0" class="empty-state-list">
          Tidak ada cabang ditemukan.
        </div>
        <div v-else v-for="group in groupedCabang" :key="group.provinsi">
          <div class="prov-group">{{ group.provinsi }}</div>
          <div 
            v-for="c in group.items" 
            :key="c.id" 
            :class="['cabang-row', { active: selectedCabang && selectedCabang.id === c.id }]" 
            @click="selectCabang(c)"
          >
            <div class="cr-dot" :style="{ background: c.status === 'aktif' ? 'var(--hijau2)' : 'var(--text3)' }"></div>
            <div class="cr-info">
              <div class="cr-name">{{ c.nama }}</div>
              <div class="cr-meta">{{ c.jumlah_unit || 0 }} unit · {{ c.jumlah_anggota || 0 }} anggota</div>
            </div>
            <div class="cr-badge">{{ c.jumlah_anggota || 0 }}</div>
          </div>
        </div>
      </div>
    </div>

    <!-- ===== PANEL KANAN: DETAIL & TAB-TAB ===== -->
    <div class="panel-right" v-if="selectedCabang">
      <!-- Topbar Detail -->
      <div class="pr-topbar">
        <div class="pr-breadcrumb">
          <span class="pr-bc-link" @click="selectedCabang = null">Kelola Cabang</span>
          <span class="pr-bc-sep"><i class="ti ti-chevron-right"></i></span>
          <span class="pr-bc-cur">{{ selectedCabang.nama }}</span>
        </div>
        <div style="display:flex;gap:7px;flex-shrink:0;">
          <button class="btn btn-ghost" @click="openEditModal(selectedCabang)"><i class="ti ti-edit"></i> Edit info cabang</button>
          <button class="btn btn-primary" @click="openAddUnitModal"><i class="ti ti-plus"></i> Tambah Unit</button>
        </div>
      </div>

      <!-- Tab Bar -->
      <div class="tabs-bar">
        <button :class="['tab', { active: activeTab === 'profil' }]" @click="switchTab('profil')">Profil Cabang</button>
        <button :class="['tab', { active: activeTab === 'unit' }]" @click="switchTab('unit')">
          Unit Latihan <span class="tab-count">{{ unitList.length }}</span>
        </button>
        <button :class="['tab', { active: activeTab === 'pelatih' }]" @click="switchTab('pelatih')">
          Pelatih <span class="tab-count">{{ trainerList.length }}</span>
        </button>
        <button :class="['tab', { active: activeTab === 'statistik' }]" @click="switchTab('statistik')">
          <i class="ti ti-chart-bar"></i> Statistik
        </button>
      </div>

      <div class="pr-body">
        <!-- 1. TAB PROFIL -->
        <div v-if="activeTab === 'profil'" class="tab-content-pane">
          <!-- Info Cabang -->
          <div class="card">
            <div class="card-head">
              <div class="card-title"><i class="ti ti-building"></i> Informasi Cabang</div>
              <button class="btn btn-ghost btn-sm" @click="openEditModal(selectedCabang)"><i class="ti ti-edit"></i> Edit</button>
            </div>
            <div class="card-body">
              <div class="info-grid">
                <div class="info-item"><div class="info-label">Nama cabang</div><div class="info-val">{{ selectedCabang.nama }}</div></div>
                <div class="info-item"><div class="info-label">Provinsi</div><div class="info-val">{{ getProvinceName(selectedCabang.provinsi_id) }}</div></div>
                <div class="info-item"><div class="info-label">Kota/Kabupaten</div><div class="info-val">{{ selectedCabang.kota_kab }}</div></div>
                <div class="info-item">
                  <div class="info-label">Status</div>
                  <div class="info-val">
                    <span :class="['bdg', selectedCabang.status === 'aktif' ? 'bdg-g' : 'bdg-r']">
                      <i :class="selectedCabang.status === 'aktif' ? 'ti ti-circle-check' : 'ti ti-circle-x'"></i>
                      {{ selectedCabang.status === 'aktif' ? 'Aktif' : 'Nonaktif' }}
                    </span>
                  </div>
                </div>
                <div class="info-item" style="grid-column:span 2;">
                  <div class="info-label">Alamat kantor</div>
                  <div class="info-val">{{ selectedCabang.alamat || 'Belum diatur' }}</div>
                </div>
                <div class="info-item"><div class="info-label">No. telepon</div><div class="info-val">{{ selectedCabang.telepon || '0274-512345' }}</div></div>
                <div class="info-item"><div class="info-label">Email cabang</div><div class="info-val">{{ selectedCabang.email || 'ygy@satrianusantara.org' }}</div></div>
                <div class="info-item"><div class="info-label">Berdiri sejak</div><div class="info-val">{{ selectedCabang.berdiri_sejak || '12 Maret 1998' }}</div></div>
                <div class="info-item"><div class="info-label">Kode cabang</div><div class="info-val monospace">{{ selectedCabang.kode }}</div></div>
              </div>
            </div>
          </div>

          <!-- Struktur Kepengurusan -->
          <div class="card">
            <div class="card-head">
              <div class="card-title"><i class="ti ti-users"></i> Struktur Kepengurusan</div>
              <button class="btn btn-primary btn-sm" @click="openAddPengurusModal"><i class="ti ti-plus"></i> Tambah Pengurus</button>
            </div>
            <div class="card-body" style="padding:0;">
              <table class="pg-table">
                <thead>
                  <tr>
                    <th>Nama</th>
                    <th>Jabatan</th>
                    <th>Tingkatan</th>
                    <th>Kontak</th>
                    <th>Aksi</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="p in staffList" :key="p.id">
                    <td>
                      <div class="pg-avatar-wrap">
                        <div class="pg-avatar">{{ getInitials(p.nama) }}</div>
                        <div>
                          <div class="pg-name">{{ p.nama }}</div>
                          <div class="pg-nomor">{{ p.nomor }}</div>
                        </div>
                      </div>
                    </td>
                    <td><span class="bdg bdg-p">{{ p.jabatan }}</span></td>
                    <td><span class="bdg bdg-y"><i class="ti ti-award"></i> {{ p.tingkatan }}</span></td>
                    <td class="pg-contact">{{ p.kontak }}</td>
                    <td>
                      <div class="action-btns" style="gap: 8px;">
                        <button class="action-btn-styled btn-edit-styled" @click="openEditPengurusModal(p)" title="Edit Pengurus"><i class="ti ti-edit"></i> Edit</button>
                        <button class="action-btn-styled btn-delete-styled" @click="deletePengurus(p.id)" title="Hapus"><i class="ti ti-trash"></i> Hapus</button>
                      </div>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>

        <!-- 2. TAB LIST UNIT LATIHAN -->
        <div v-else-if="activeTab === 'unit'" class="tab-content-pane">
          <div v-if="selectedUnit" class="unit-detail-pane">
            <div style="display:flex;align-items:center;gap:10px;margin-bottom:16px;">
              <button class="btn btn-ghost btn-sm" @click="selectedUnit = null"><i class="ti ti-arrow-left"></i> Kembali ke daftar unit</button>
              <div style="font-size:14px;font-weight:700;">Detail: {{ selectedUnit.nama }}</div>
            </div>
            
            <div class="stat-mini-grid">
              <div class="stat-mini"><div class="stat-mini-val">{{ selectedUnit.jumlah_anggota || 124 }}</div><div class="stat-mini-label">Anggota aktif</div></div>
              <div class="stat-mini"><div class="stat-mini-val">2x</div><div class="stat-mini-label">Latihan/minggu</div></div>
              <div class="stat-mini"><div class="stat-mini-val">89%</div><div class="stat-mini-label">Kehadiran rata-rata</div></div>
              <div class="stat-mini"><div class="stat-mini-val" style="color:var(--hijau);">100%</div><div class="stat-mini-label">BLBA bulan ini</div></div>
            </div>

            <div class="card">
              <div class="card-head">
                <div class="card-title"><i class="ti ti-info-circle"></i> Info Unit</div>
                <button class="btn btn-ghost btn-sm" @click="openEditUnitModal(selectedUnit)"><i class="ti ti-edit"></i> Edit</button>
              </div>
              <div class="card-body">
                <div class="info-grid">
                  <div class="info-item"><div class="info-label">Nama unit</div><div class="info-val">{{ selectedUnit.nama }}</div></div>
                  <div class="info-item"><div class="info-label">Status</div><div class="info-val"><span class="bdg bdg-g">Aktif</span></div></div>
                  <div class="info-item" style="grid-column:span 2;"><div class="info-label">Lokasi latihan</div><div class="info-val">{{ selectedUnit.lokasi || 'Lapangan Malioboro, Jl. Malioboro, Yogyakarta' }}</div></div>
                  <div class="info-item"><div class="info-label">Jadwal tetap</div><div class="info-val">{{ selectedUnit.jadwal || 'Sabtu & Selasa, 06.00 WIB' }}</div></div>
                  <div class="info-item"><div class="info-label">Berdiri sejak</div><div class="info-val">5 Agustus 2002</div></div>
                </div>
              </div>
            </div>

            <!-- PIC/Coordinator -->
            <div class="card">
              <div class="card-head">
                <div class="card-title"><i class="ti ti-user-check"></i> PIC / Koordinator Unit</div>
                <button class="btn btn-ghost btn-sm" @click="openGantiPICModal(selectedUnit)"><i class="ti ti-refresh"></i> Ganti PIC</button>
              </div>
              <div class="card-body">
                <div class="pic-highlight-box">
                  <div class="pic-big-av">{{ getInitials(selectedUnit.pic || 'Wisnu Saputro') }}</div>
                  <div style="flex:1;">
                    <div style="font-size:13px;font-weight:700;">{{ selectedUnit.pic || 'Wisnu Saputro' }}</div>
                    <div style="font-size:11px;color:var(--text2);font-family:monospace;">YO-YGY-00234</div>
                    <div style="margin-top:4px;display:flex;gap:5px;">
                      <span class="bdg bdg-y"><i class="ti ti-award"></i> PH Jurus 7</span>
                      <span class="bdg bdg-p">Koordinator Unit</span>
                    </div>
                  </div>
                  <div class="pic-contacts">
                    0857-1234-5678<br>
                    wisnu@satrianusantara.org
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div v-else class="card">
            <div class="card-head">
              <div class="card-title"><i class="ti ti-home"></i> Daftar Unit Latihan</div>
              <button class="btn btn-primary btn-sm" @click="openAddUnitModal"><i class="ti ti-plus"></i> Tambah Unit</button>
            </div>
            <div class="card-body" style="padding:0;">
              <table class="unit-table">
                <thead>
                  <tr>
                    <th>Unit</th>
                    <th>Jadwal Tetap</th>
                    <th>PIC / Koordinator</th>
                    <th>Anggota</th>
                    <th>Status</th>
                    <th>Aksi</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="u in paginatedUnits" :key="u.id" @click="selectedUnit = u" style="cursor:pointer;">
                    <td>
                      <div class="unit-name-cell">
                        <div class="unit-icon"><i class="ti ti-home"></i></div>
                        <div>
                          <div class="unit-name">{{ u.nama }}</div>
                          <div class="unit-loc"><i class="ti ti-map-pin"></i> {{ u.lokasi || 'Kota Yogyakarta' }}</div>
                        </div>
                      </div>
                    </td>
                    <td>
                      <div class="jadwal-chips">
                        <span class="jadwal-chip"><i class="ti ti-calendar"></i> {{ u.jadwal || 'Minggu 07.00 WIB' }}</span>
                      </div>
                    </td>
                    <td>
                      <div class="pic-cell">
                        <div class="pic-av">{{ getInitials(u.pic || 'WS') }}</div>
                        <div>
                          <div style="font-size:12px;font-weight:600;">{{ u.pic || 'Wisnu Saputro' }}</div>
                          <div style="font-size:10px;color:var(--text3);">YO-YGY-00234</div>
                        </div>
                      </div>
                    </td>
                    <td><strong>{{ u.jumlah_anggota || 0 }}</strong></td>
                    <td><span class="bdg bdg-g">Aktif</span></td>
                    <td @click.stopPropagation>
                      <div class="action-btns" style="gap: 8px;">
                        <button class="action-btn-styled btn-view-styled" @click="selectedUnit = u" title="Lihat Detail"><i class="ti ti-eye"></i> Detail</button>
                        <button class="action-btn-styled btn-edit-styled" @click="openEditUnitModal(u)" title="Edit Unit"><i class="ti ti-edit"></i> Edit</button>
                        <button class="action-btn-styled btn-delete-styled" @click="deleteUnit(u.id)" title="Hapus Unit"><i class="ti ti-trash"></i> Hapus</button>
                      </div>
                    </td>
                  </tr>
                </tbody>
              </table>
              <Pagination 
                v-model:currentPage="pageUnit" 
                v-model:itemsPerPage="limitUnit" 
                :totalItems="unitList.length" 
              />
            </div>
          </div>
        </div>

        <!-- 3. TAB PELATIH -->
        <div v-else-if="activeTab === 'pelatih'" class="tab-content-pane">
          <div class="card">
            <div class="card-head">
              <div class="card-title"><i class="ti ti-user-star"></i> Tim Pelatih Cabang</div>
              <button class="btn btn-primary btn-sm" @click="openAddTrainerModal"><i class="ti ti-plus"></i> Tambah Pelatih</button>
            </div>
            <div class="card-body" style="padding:0;">
              <table class="pg-table">
                <thead>
                  <tr>
                    <th>Pelatih</th>
                    <th>Jenis</th>
                    <th>Tingkatan</th>
                    <th>Kategori Jarak</th>
                    <th>Transport</th>
                    <th>Status</th>
                    <th>Aksi</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="t in paginatedTrainers" :key="t.id">
                    <td>
                      <div class="pg-avatar-wrap">
                        <div class="pg-avatar" style="background:var(--hijau)">{{ getInitials(t.nama) }}</div>
                        <div>
                          <div class="pg-name">{{ t.nama }}</div>
                          <div class="pg-nomor">{{ t.nomor }}</div>
                        </div>
                      </div>
                    </td>
                    <td><span :class="['bdg', t.jenis === 'Pelatih Pusat' ? 'bdg-r' : 'bdg-b']">{{ t.jenis || 'Pelatih Daerah' }}</span></td>
                    <td><span class="bdg bdg-y"><i class="ti ti-award"></i> {{ t.tingkatan }}</span></td>
                    <td>{{ t.jarak || 'Jarak Dekat' }}</td>
                    <td style="font-weight:700;color:var(--text1)">Rp {{ formatRupiah(t.transport || 30000) }}</td>
                    <td><span :class="['bdg', t.status === 'Aktif' ? 'bdg-g' : 'bdg-r']">{{ t.status || 'Aktif' }}</span></td>
                    <td>
                      <div class="action-btns" style="gap: 8px;">
                        <button class="action-btn-styled btn-edit-styled" @click="openEditTrainerModal(t)" title="Edit Pelatih"><i class="ti ti-edit"></i> Edit</button>
                        <button class="action-btn-styled btn-delete-styled" @click="deleteTrainer(t.id)" title="Hapus"><i class="ti ti-trash"></i> Hapus</button>
                      </div>
                    </td>
                  </tr>
                </tbody>
              </table>
              <Pagination 
                v-model:currentPage="pageTrainer" 
                v-model:itemsPerPage="limitTrainer" 
                :totalItems="trainerList.length" 
              />
            </div>
          </div>
        </div>

        <!-- 4. TAB STATISTIK -->
        <div v-else-if="activeTab === 'statistik'" class="tab-content-pane">
          <!-- Toolbar & Options -->
          <div style="display:flex;align-items:center;justify-content:space-between;gap:10px;margin-bottom:8px;">
            <div style="display:flex;gap:6px;">
              <select v-model="statsPeriod" class="filter-select" style="width:160px;height:34px;">
                <option value="6">6 Bulan Terakhir</option>
                <option value="12">12 Bulan Terakhir</option>
                <option value="24">24 Bulan Terakhir</option>
              </select>
              <select v-model="statsMetric" class="filter-select" style="width:200px;height:34px;">
                <option value="kehadiran">Rata-rata Kehadiran (%)</option>
                <option value="iuran">Pencapaian BLBA (%)</option>
                <option value="anggota">Jumlah Anggota Aktif</option>
              </select>
            </div>
            <div style="display:flex;gap:6px;">
              <button class="export-btn" @click="alert('Pdf berhasil diekspor!')"><i class="ti ti-file-type-pdf"></i> PDF</button>
              <button class="export-btn" @click="alert('Excel berhasil diekspor!')"><i class="ti ti-file-spreadsheet"></i> Excel</button>
            </div>
          </div>

          <!-- Summary Mini Cards -->
          <div class="stat-mini-grid" style="margin-bottom:12px;">
            <div class="stat-mini">
              <div v-if="loadingTrends" class="skeleton-block" style="width:50px;height:22px;margin-bottom:4px;"></div>
              <div v-else class="stat-mini-val" style="color:var(--hijau);">{{ summaryAvgKehadiran }}%</div>
              <div class="stat-mini-label">Kehadiran Rata-rata</div>
            </div>
            <div class="stat-mini">
              <div v-if="loadingTrends" class="skeleton-block" style="width:50px;height:22px;margin-bottom:4px;"></div>
              <div v-else class="stat-mini-val" style="color:var(--hijau);">{{ summaryBlba }}%</div>
              <div class="stat-mini-label">Pencapaian BLBA</div>
            </div>
            <div class="stat-mini">
              <div v-if="loadingTrends" class="skeleton-block" style="width:50px;height:22px;margin-bottom:4px;"></div>
              <div v-else class="stat-mini-val" style="color:var(--biru);">{{ summaryAnggota }}</div>
              <div class="stat-mini-label">Anggota Aktif</div>
            </div>
            <div class="stat-mini">
              <div v-if="loadingTrends" class="skeleton-block" style="width:60px;height:22px;margin-bottom:4px;"></div>
              <div v-else class="stat-mini-val" style="color:var(--kuning);">{{ summaryKas }}</div>
              <div class="stat-mini-label">Kas Unit / Cabang</div>
            </div>
          </div>

          <!-- Trends Diagram Card -->
          <div class="card">
            <div class="card-head">
              <div class="card-title"><i class="ti ti-chart-line"></i> Tren {{ statsMetricLabel }} — {{ selectedCabang.nama }}</div>
            </div>
            <div class="card-body" style="padding: 16px; position: relative;">
              <div v-if="loadingTrends" style="height:180px;display:flex;align-items:center;justify-content:center;color:var(--text3);font-size:12px;">
                <i class="ti ti-loader-2 spin" style="font-size:24px;color:var(--hijau);margin-right:8px;"></i> Memuat data tren {{ statsMetricLabel.toLowerCase() }}...
              </div>
              <!-- Beautiful dynamic SVG Chart -->
              <svg v-else viewBox="0 0 540 180" width="100%" height="180" style="display:block;overflow:visible;">
                <defs>
                  <linearGradient id="trendGrad" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="0%" stop-color="var(--hijau)" stop-opacity="0.35" />
                    <stop offset="100%" stop-color="var(--hijau)" stop-opacity="0.0" />
                  </linearGradient>
                </defs>

                <!-- Grid lines -->
                <line x1="30" y1="30" x2="510" y2="30" stroke="#f1f3f5" stroke-dasharray="3,3" />
                <line x1="30" y1="68" x2="510" y2="68" stroke="#f1f3f5" stroke-dasharray="3,3" />
                <line x1="30" y1="106" x2="510" y2="106" stroke="#f1f3f5" stroke-dasharray="3,3" />
                <line x1="30" y1="144" x2="510" y2="144" stroke="#e2e8f0" />

                <!-- Area fill -->
                <path v-if="chartSvgData.areaPath" :d="chartSvgData.areaPath" fill="url(#trendGrad)" />

                <!-- Line stroke -->
                <path v-if="chartSvgData.path" :d="chartSvgData.path" fill="none" stroke="var(--hijau)" stroke-width="3" stroke-linecap="round" stroke-linejoin="round" />

                <!-- Data circles & labels -->
                <g v-for="(c, idx) in chartSvgData.coords" :key="idx">
                  <circle
                    :cx="c.x" :cy="c.y" :r="hoveredPoint === c ? 7 : (idx === chartSvgData.coords.length - 1 ? 5 : 4)"
                    :fill="hoveredPoint === c || idx === chartSvgData.coords.length - 1 ? 'var(--hijau)' : '#fff'"
                    :stroke="'var(--hijau)'"
                    stroke-width="2.5"
                    style="cursor:pointer;transition:all .2s ease;"
                    @mouseenter="hoveredPoint = c"
                    @mouseleave="hoveredPoint = null"
                  />
                  <!-- X Axis Month Label -->
                  <text 
                    v-if="chartSvgData.coords.length <= 12 || idx % Math.ceil(chartSvgData.coords.length / 12) === 0 || idx === chartSvgData.coords.length - 1"
                    :x="c.x" y="164" font-size="10" fill="#9ca3af" text-anchor="middle" font-weight="500"
                  >{{ c.month }}</text>
                </g>

                <!-- Hover Tooltip Callout -->
                <g v-if="hoveredPoint">
                  <rect
                    :x="Math.min(450, Math.max(10, hoveredPoint.x - 45))"
                    :y="Math.max(5, hoveredPoint.y - 32)"
                    width="90" height="24" rx="5" fill="#1e293b" opacity="0.9"
                  />
                  <text
                    :x="Math.min(450, Math.max(10, hoveredPoint.x - 45)) + 45"
                    :y="Math.max(5, hoveredPoint.y - 32) + 16"
                    font-size="11" font-weight="bold" fill="#ffffff" text-anchor="middle"
                  >
                    {{ hoveredPoint.val }}{{ statsMetric !== 'anggota' ? '%' : '' }}
                  </text>
                </g>

                <!-- Value Indicator Text at top right -->
                <text x="510" y="20" font-size="11" font-weight="bold" fill="var(--hijau)" text-anchor="end">
                  Terakhir: {{ statsMetricValue }}
                </text>
              </svg>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="panel-right empty" v-else>
      <i class="ti ti-map-pin" style="font-size: 48px; color: var(--text3);"></i>
      <p style="font-size: 13px; color: var(--text2); margin-top: 8px;">Silakan pilih salah satu cabang di sebelah kiri untuk melihat detail informasi cabang.</p>
    </div>

    <!-- MODAL ADD/EDIT CABANG -->
    <div v-if="showModal" class="modal-overlay" @click.self="showModal = false">
      <div class="modal-card">
        <div class="modal-header">
          <h2 class="modal-title">{{ isEdit ? 'Edit Info Cabang' : 'Tambah Cabang Baru' }}</h2>
          <button class="modal-close" @click="showModal = false"><i class="ti ti-x"></i></button>
        </div>
        <form @submit.prevent="saveCabang" class="modal-form">
          <div class="form-group">
            <label class="form-label">Kode Cabang</label>
            <input v-model="form.kode" type="text" class="form-input" placeholder="Contoh: YGY" required :disabled="isEdit" />
          </div>
          <div class="form-group">
            <label class="form-label">Nama Cabang</label>
            <input v-model="form.nama" type="text" class="form-input" placeholder="Contoh: Kota Yogyakarta" required />
          </div>
          <div class="form-group">
            <label class="form-label">Provinsi</label>
            <select v-model="form.provinsi_id" class="form-select" required>
              <option value="" disabled>Pilih Provinsi</option>
              <option v-for="p in provinsiList" :key="p.id" :value="p.id">{{ p.nama }}</option>
            </select>
          </div>
          <div class="form-group">
            <label class="form-label">Kota/Kabupaten</label>
            <input v-model="form.kota_kab" type="text" class="form-input" placeholder="Contoh: Kota Yogyakarta" required />
          </div>
          <div class="form-group">
            <label class="form-label">Alamat Sekretariat</label>
            <textarea v-model="form.alamat" class="form-input textarea" placeholder="Alamat lengkap kantor sekretariat..."></textarea>
          </div>
          <div class="form-group">
            <label class="form-label">Pilih Titik Lokasi Peta</label>
            <div id="map-cabang" style="height: 200px; width: 100%; border-radius: 8px; border: 1.5px solid var(--border2); z-index: 1;"></div>
          </div>
          <div class="form-row" style="display:flex; gap:10px;">
            <div class="form-group" style="flex:1;">
              <label class="form-label">Latitude Peta</label>
              <input v-model="form.lat" type="text" class="form-input" placeholder="Contoh: -7.7956" />
            </div>
            <div class="form-group" style="flex:1;">
              <label class="form-label">Longitude Peta</label>
              <input v-model="form.lng" type="text" class="form-input" placeholder="Contoh: 110.3695" />
            </div>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-outline" style="width:auto" @click="showModal = false">Batal</button>
            <button type="submit" class="btn btn-primary" style="width:auto" :disabled="submitting">Simpan</button>
          </div>
        </form>
      </div>
    </div>

    <!-- MODAL ADD/EDIT UNIT -->
    <div v-if="showUnitModal" class="modal-overlay" @click.self="showUnitModal = false">
      <div class="modal-card">
        <div class="modal-header">
          <h2 class="modal-title">{{ editingUnitId ? 'Edit Unit Latihan' : 'Tambah Unit Latihan Baru' }}</h2>
          <button class="modal-close" @click="showUnitModal = false"><i class="ti ti-x"></i></button>
        </div>
        <form @submit.prevent="saveUnit" class="modal-form">
          <div class="form-group">
            <label class="form-label">Nama Unit Latihan</label>
            <input v-model="unitForm.nama" type="text" class="form-input" placeholder="Contoh: Unit Malioboro" required />
          </div>
          <div class="form-group">
            <label class="form-label">Lokasi / Tempat Latihan</label>
            <input v-model="unitForm.lokasi" type="text" class="form-input" placeholder="Contoh: Lapangan depan Vredeburg" required />
          </div>
          <div class="form-group">
            <label class="form-label">Jadwal Rutin</label>
            <input v-model="unitForm.jadwal" type="text" class="form-input" placeholder="Contoh: Sabtu & Selasa, 06.00 WIB" required />
          </div>
          <div class="form-group">
            <label class="form-label">PIC / Koordinator Unit</label>
            <input v-model="unitForm.pic" type="text" class="form-input" placeholder="Nama Koordinator Unit" required />
          </div>
          <div class="form-group">
            <label class="form-label">Pilih Titik Lokasi Peta</label>
            <div id="map-unit" style="height: 200px; width: 100%; border-radius: 8px; border: 1.5px solid var(--border2); z-index: 1;"></div>
          </div>
          <div class="form-row" style="display:flex; gap:10px;">
            <div class="form-group" style="flex:1;">
              <label class="form-label">Latitude Peta</label>
              <input v-model="unitForm.lat" type="text" class="form-input" placeholder="Contoh: -7.7956" />
            </div>
            <div class="form-group" style="flex:1;">
              <label class="form-label">Longitude Peta</label>
              <input v-model="unitForm.lng" type="text" class="form-input" placeholder="Contoh: 110.3695" />
            </div>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-outline" style="width:auto" @click="showUnitModal = false">Batal</button>
            <button type="submit" class="btn btn-primary" style="width:auto">Simpan Unit</button>
          </div>
        </form>
      </div>
    </div>

    <!-- MODAL ADD/EDIT PENGURUS -->
    <div v-if="showPengurusModal" class="modal-overlay" @click.self="showPengurusModal = false">
      <div class="modal-card">
        <div class="modal-header">
          <h2 class="modal-title">{{ editingPengurusId ? 'Edit Struktur Pengurus' : 'Tambah Struktur Pengurus' }}</h2>
          <button class="modal-close" @click="showPengurusModal = false"><i class="ti ti-x"></i></button>
        </div>
        <form @submit.prevent="savePengurus" class="modal-form">
          <div class="form-group">
            <label class="form-label">Nama Pengurus</label>
            <input v-model="pengurusForm.nama" type="text" class="form-input" required />
          </div>
          <div class="form-group">
            <label class="form-label">Jabatan</label>
            <input v-model="pengurusForm.jabatan" type="text" class="form-input" placeholder="Contoh: Sekretaris Cabang" required />
          </div>
          <div class="form-group">
            <label class="form-label">Tingkatan Sabuk / Jurus</label>
            <input v-model="pengurusForm.tingkatan" type="text" class="form-input" placeholder="Contoh: PK Jurus 5" required />
          </div>
          <div class="form-group">
            <label class="form-label">No. Telepon / WA</label>
            <input v-model="pengurusForm.kontak" type="text" class="form-input" required />
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-outline" style="width:auto" @click="showPengurusModal = false">Batal</button>
            <button type="submit" class="btn btn-primary" style="width:auto">Simpan Pengurus</button>
          </div>
        </form>
      </div>
    </div>

    <!-- MODAL ADD/EDIT PELATIH -->
    <div v-if="showTrainerModal" class="modal-overlay" @click.self="showTrainerModal = false">
      <div class="modal-card">
        <div class="modal-header">
          <h2 class="modal-title">{{ editingTrainerId ? 'Edit Data Pelatih' : 'Tambah Pelatih Cabang' }}</h2>
          <button class="modal-close" @click="showTrainerModal = false"><i class="ti ti-x"></i></button>
        </div>
        <form @submit.prevent="saveTrainer" class="modal-form">
          <div class="form-group" v-if="!editingTrainerId">
            <label class="form-label">Nama / Anggota</label>
            <input v-model="trainerForm.nama" type="text" class="form-input" placeholder="Nama Pelatih baru" required />
          </div>
          <div class="form-group">
            <label class="form-label">Jenis Pelatih</label>
            <select v-model="trainerForm.jenis" class="form-select">
              <option value="Pelatih Daerah">Pelatih Daerah</option>
              <option value="Pelatih Pusat">Pelatih Pusat</option>
            </select>
          </div>
          <div class="form-group">
            <label class="form-label">Kategori Jarak / Transport</label>
            <select v-model="trainerForm.jarak" class="form-select" @change="onJarakChange">
              <option value="Sukarela">Sukarela — Rp 0</option>
              <option value="Jarak Dekat">Jarak Dekat — Rp 30.000</option>
              <option value="Jarak Sedang">Jarak Sedang — Rp 50.000</option>
              <option value="Jarak Jauh">Jarak Jauh — Rp 75.000</option>
            </select>
          </div>
          <div class="form-group" v-if="!editingTrainerId">
            <label class="form-label">Tanggal Mulai Bertugas</label>
            <input v-model="trainerForm.tanggal_mulai" type="date" class="form-input" />
          </div>
          <div class="form-group" v-if="!editingTrainerId">
            <label class="form-label">No. SK / Surat Tugas</label>
            <input v-model="trainerForm.sk" type="text" class="form-input" placeholder="No. SK Tugas (opsional)" />
          </div>
          <div class="form-group" v-if="editingTrainerId">
            <label class="form-label">Status Keaktifan</label>
            <select v-model="trainerForm.status" class="form-select">
              <option value="Aktif">Aktif</option>
              <option value="Nonaktif">Nonaktif</option>
            </select>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-outline" style="width:auto" @click="showTrainerModal = false">Batal</button>
            <button type="submit" class="btn btn-primary" style="width:auto">Simpan Pelatih</button>
          </div>
        </form>
      </div>
    </div>

    <!-- MODAL GANTI PIC UNIT -->
    <div v-if="showGantiPICModal" class="modal-overlay" @click.self="showGantiPICModal = false">
      <div class="modal-card">
        <div class="modal-header">
          <h2 class="modal-title">Ganti PIC Koordinator Unit</h2>
          <button class="modal-close" @click="showGantiPICModal = false"><i class="ti ti-x"></i></button>
        </div>
        <form @submit.prevent="savePIC" class="modal-form">
          <div style="background:var(--surface); border:1px solid var(--border); padding:10px; border-radius:var(--r8); margin-bottom:8px;">
            <div style="font-size:10px; color:var(--text3); font-weight:700;">PIC SAAT INI</div>
            <div style="font-size:12px; font-weight:700; margin-top:2px;">{{ selectedUnitPic }}</div>
          </div>
          <div class="form-group">
            <label class="form-label">Nama PIC Koordinator Baru</label>
            <input v-model="picFormName" type="text" class="form-input" placeholder="Nama Koordinator Baru" required />
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-outline" style="width:auto" @click="showGantiPICModal = false">Batal</button>
            <button type="submit" class="btn btn-primary" style="width:auto">Ganti Koordinator</button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({ title: 'Kelola Cabang' })

const api = useApi()

const search = ref('')
const filterProvinsi = ref('')
const showModal = ref(false)
const showUnitModal = ref(false)
const showPengurusModal = ref(false)
const showTrainerModal = ref(false)
const showGantiPICModal = ref(false)
const isEdit = ref(false)
const selectedCabang = ref<any>(null)
const selectedUnit = ref<any>(null)
const activeTab = ref('profil')
const loadingCabangList = ref(false)
const submitting = ref(false)

const statsPeriod = ref('12')
const statsMetric = ref('kehadiran')
const loadingTrends = ref(false)
const trendsResponse = ref<any>(null)
const hoveredPoint = ref<any>(null)

const statsMetricLabel = computed(() => {
  if (statsMetric.value === 'kehadiran') return 'Rata-rata Kehadiran (%)'
  if (statsMetric.value === 'iuran') return 'Pencapaian BLBA (%)'
  return 'Jumlah Anggota Aktif'
})

const trendPoints = computed(() => trendsResponse.value?.points || [])
const summaryAvgKehadiran = computed(() => trendsResponse.value?.avg_kehadiran_pct ?? 0)
const summaryBlba = computed(() => trendsResponse.value?.blba_pct ?? 0)
const summaryAnggota = computed(() => trendsResponse.value?.total_anggota ?? 0)
const summaryKas = computed(() => {
  if (trendsResponse.value?.kas_unit) {
    return 'Rp ' + (trendsResponse.value.kas_unit / 1000000).toFixed(1) + 'Jt'
  }
  return 'Rp 0'
})

const statsMetricValue = computed(() => {
  if (!trendsResponse.value || !trendPoints.value.length) {
    if (statsMetric.value === 'kehadiran') return '0%'
    if (statsMetric.value === 'iuran') return '0%'
    return '0 Anggota'
  }
  const last = trendPoints.value[trendPoints.value.length - 1]
  if (statsMetric.value === 'kehadiran') return `${last.kehadiran_pct}%`
  if (statsMetric.value === 'iuran') return `${last.iuran_pct}%`
  return `${last.anggota_count} Anggota`
})

const chartSvgData = computed(() => {
  const points = trendPoints.value
  if (!points || points.length === 0) return { path: '', areaPath: '', coords: [] }

  const width = 540
  const height = 180
  const padLeft = 40
  const padRight = 30
  const padTop = 30
  const padBottom = 40

  const chartW = width - padLeft - padRight
  const chartH = height - padTop - padBottom

  const vals = points.map((p: any) => {
    if (statsMetric.value === 'kehadiran') return p.kehadiran_pct
    if (statsMetric.value === 'iuran') return p.iuran_pct
    return p.anggota_count
  })

  let minVal = Math.min(...vals)
  let maxVal = Math.max(...vals)
  if (statsMetric.value !== 'anggota') {
    minVal = 0
    maxVal = 100
  } else {
    if (minVal === maxVal) {
      minVal = Math.max(0, minVal - 10)
      maxVal = maxVal + 10
    }
  }

  const coords = points.map((p: any, idx: number) => {
    const val = vals[idx]
    const x = padLeft + (idx / Math.max(1, points.length - 1)) * chartW
    const y = padTop + (1 - (val - minVal) / Math.max(1, maxVal - minVal)) * chartH
    return {
      x: Math.round(x * 10) / 10,
      y: Math.round(y * 10) / 10,
      val,
      month: p.month,
      fullMonth: p.full_month
    }
  })

  const pathStr = coords.map((c: any, i: number) => `${i === 0 ? 'M' : 'L'} ${c.x} ${c.y}`).join(' ')
  const lastCoord = coords[coords.length - 1]
  const firstCoord = coords[0]
  const areaPathStr = `${pathStr} L ${lastCoord.x} ${height - padBottom} L ${firstCoord.x} ${height - padBottom} Z`

  return { path: pathStr, areaPath: areaPathStr, coords }
})

const fetchTrends = async () => {
  if (!selectedCabang.value) return
  loadingTrends.value = true
  try {
    const res = await api.get(`/organization/cabang/${selectedCabang.value.id}/trends?period=${statsPeriod.value}`)
    trendsResponse.value = res
  } catch (e) {
    console.error('Failed to fetch trends', e)
  } finally {
    loadingTrends.value = false
  }
}

watch(statsPeriod, () => {
  fetchTrends()
})


const cabangData = ref<any[]>([])
const totalCabang = ref(0)
const provinsiList = ref<any[]>([])

// Form states
const form = ref({
  id: '',
  kode: '',
  nama: '',
  provinsi_id: '',
  kota_kab: '',
  alamat: '',
  telepon: '',
  email: '',
  berdiri_sejak: '',
  lat: '',
  lng: ''
})

const unitForm = ref({
  id: '',
  nama: '',
  lokasi: '',
  jadwal: '',
  pic: '',
  lat: '',
  lng: ''
})
const editingUnitId = ref<string | null>(null)

const pengurusForm = ref({
  id: '',
  nama: '',
  jabatan: '',
  tingkatan: '',
  kontak: '',
  nomor: ''
})
const editingPengurusId = ref<string | null>(null)

// Trainer modal states
const trainerForm = ref({
  id: '',
  nama: '',
  nomor: '',
  jenis: 'Pelatih Daerah',
  tingkatan: 'PH Jurus 6',
  jarak: 'Jarak Dekat',
  transport: 30000,
  tanggal_mulai: '',
  sk: '',
  status: 'Aktif'
})
const editingTrainerId = ref<string | null>(null)

// PIC modal states
const selectedUnitPic = ref('')
const picFormName = ref('')

// Sub lists
const unitList = ref<any[]>([])
const staffList = ref<any[]>([])
const trainerList = ref<any[]>([])

// Unit Pagination states
const pageUnit = ref(1)
const limitUnit = ref(10)
const paginatedUnits = computed(() => {
  return unitList.value.slice((pageUnit.value - 1) * limitUnit.value, pageUnit.value * limitUnit.value)
})

// Trainer Pagination states
const pageTrainer = ref(1)
const limitTrainer = ref(10)
const paginatedTrainers = computed(() => {
  return trainerList.value.slice((pageTrainer.value - 1) * limitTrainer.value, pageTrainer.value * limitTrainer.value)
})

const fetchProvinsi = async () => {
  provinsiList.value = [
    { id: 1, nama: 'DI Yogyakarta' },
    { id: 2, nama: 'DKI Jakarta' },
    { id: 3, nama: 'Jawa Barat' },
    { id: 4, nama: 'Jawa Tengah' },
    { id: 5, nama: 'Jawa Timur' },
    { id: 6, nama: 'Banten' },
    { id: 7, nama: 'Sulawesi Selatan' }
  ]
}

const getProvinceName = (provId: any) => {
  const p = provinsiList.value.find(x => x.id === Number(provId))
  return p ? p.nama : 'DI Yogyakarta'
}

const fetchCabang = async () => {
  loadingCabangList.value = true
  try {
    const data = await api.get(`/organization/cabang?search=${search.value}`)
    cabangData.value = data.data || []
    totalCabang.value = data.total || 0
    if (cabangData.value.length > 0 && !selectedCabang.value) {
      selectCabang(cabangData.value[0])
    }
  } catch (e) {
    console.error('Failed to fetch cabang', e)
  } finally {
    loadingCabangList.value = false
  }
}

const selectCabang = async (c: any) => {
  selectedCabang.value = c
  selectedUnit.value = null
  activeTab.value = 'profil'
  await fetchUnitsForCabang(c.id)
  await fetchTrainersForCabang(c.id)
  await fetchTrends()
  await fetchPengurusForCabang(c.id)
}

const fetchPengurusForCabang = async (cabangId: string) => {
  try {
    const res = await api.get(`/organization/cabang/${cabangId}/pengurus`)
    staffList.value = res.data || []
  } catch (error) {
    console.error('Gagal memuat daftar pengurus:', error)
    staffList.value = []
  }
}

const fetchUnitsForCabang = async (cabangId: string) => {
  try {
    const res = await api.get(`/organization/cabang/${cabangId}/unit`)
    unitList.value = res || []
  } catch (e) {
    console.error(e)
    unitList.value = [
      { id: '1', nama: 'Unit Malioboro', lokasi: 'Jl. Malioboro, Yogyakarta', jadwal: 'Sabtu & Selasa, 06.00 WIB', pic: 'Wisnu Saputro', jumlah_anggota: 124, lat: '-7.7926', lng: '110.3658' },
      { id: '2', nama: 'Unit Kotagede', lokasi: 'Kotagede, Yogyakarta', jadwal: 'Minggu 07.00 WIB', pic: 'Pratiwi Riyadi', jumlah_anggota: 98, lat: '-7.8284', lng: '110.4008' },
      { id: '3', nama: 'Unit Gondokusuman', lokasi: 'Gondokusuman, Yogyakarta', jadwal: 'Rabu & Sabtu, 06.30 WIB', pic: 'Sunaryo', jumlah_anggota: 87, lat: '-7.7865', lng: '110.3831' }
    ]
  }
}

const fetchTrainersForCabang = async (cabangId: string) => {
  try {
    const res = await api.get(`/organization/pelatih?cabang_id=${cabangId}`)
    trainerList.value = (res || []).map((p: any) => {
      let transport = 30000
      if (p.kategori_transport === 'Jarak Sedang') transport = 50000
      else if (p.kategori_transport === 'Jarak Jauh') transport = 75000
      else if (p.kategori_transport === 'Sukarela') transport = 0

      return {
        id: p.id,
        nama: p.nama_lengkap,
        nomor: p.nomor_anggota || '-',
        jenis: p.jenis === 'pusat' ? 'Pelatih Pusat' : 'Pelatih Daerah',
        tingkatan: p.tingkatan || 'Dasar',
        jarak: p.kategori_transport || 'Jarak Dekat',
        transport: transport,
        status: p.is_active ? 'Aktif' : 'Nonaktif'
      }
    })
  } catch (e) {
    console.error('Failed to fetch trainers:', e)
    trainerList.value = []
  }
}

const groupedCabang = computed(() => {
  const provMap: Record<string, any[]> = {}
  
  cabangData.value.forEach(c => {
    const provName = getProvinceName(c.provinsi_id)
    if (filterProvinsi.value && filterProvinsi.value !== provName) return
    if (!provMap[provName]) {
      provMap[provName] = []
    }
    provMap[provName].push(c)
  })

  return Object.keys(provMap).map(k => ({
    provinsi: k,
    items: provMap[k]
  }))
})

const switchTab = (tab: string) => {
  activeTab.value = tab
  selectedUnit.value = null
}

const openCreateModal = () => {
  isEdit.value = false
  form.value = {
    id: '',
    kode: '',
    nama: '',
    provinsi_id: '',
    kota_kab: '',
    alamat: '',
    telepon: '',
    email: '',
    berdiri_sejak: '',
    lat: '',
    lng: ''
  }
  showModal.value = true
  initMap('map-cabang', form)
}

const openEditModal = (c: any) => {
  isEdit.value = true
  form.value = {
    id: c.id,
    kode: c.kode,
    nama: c.nama,
    provinsi_id: c.provinsi_id || '',
    kota_kab: c.kota_kab,
    alamat: c.alamat || '',
    telepon: c.telepon || '0274-512345',
    email: c.email || 'ygy@satrianusantara.org',
    berdiri_sejak: c.berdiri_sejak || '12 Maret 1998',
    lat: c.lat || '',
    lng: c.lng || ''
  }
  showModal.value = true
  initMap('map-cabang', form)
}

const saveCabang = async () => {
  submitting.value = true
  try {
    if (isEdit.value) {
      const payload = {
        nama: form.value.nama,
        provinsi_id: parseInt(form.value.provinsi_id),
        kota_kab: form.value.kota_kab,
        alamat: form.value.alamat
      }
      await api.put(`/organization/cabang/${form.value.id}`, payload)
      alert('Informasi cabang berhasil diperbarui!')
    } else {
      const payload = {
        kode: form.value.kode,
        nama: form.value.nama,
        provinsi_id: parseInt(form.value.provinsi_id),
        kota_kab: form.value.kota_kab,
        alamat: form.value.alamat
      }
      await api.post('/organization/cabang', payload)
      alert('Cabang baru berhasil dibuat!')
    }
    showModal.value = false
    await fetchCabang()
  } catch (e: any) {
    alert(e.message || 'Gagal menyimpan cabang')
  } finally {
    submitting.value = false
  }
}

// Unit Actions
const openAddUnitModal = () => {
  editingUnitId.value = null
  unitForm.value = { id: '', nama: '', lokasi: '', jadwal: '', pic: '', lat: '-7.7956', lng: '110.3695' }
  showUnitModal.value = true
  initMap('map-unit', unitForm)
}

const openEditUnitModal = (u: any) => {
  editingUnitId.value = u.id
  unitForm.value = {
    id: u.id,
    nama: u.nama,
    lokasi: u.lokasi || '',
    jadwal: u.jadwal || '',
    pic: u.pic || '',
    lat: u.lat || '-7.7956',
    lng: u.lng || '110.3695'
  }
  showUnitModal.value = true
  initMap('map-unit', unitForm)
}

const saveUnit = async () => {
  if (!selectedCabang.value) return
  try {
    const payload = {
      cabang_id: selectedCabang.value.id,
      nama: unitForm.value.nama,
      lokasi_nama: unitForm.value.lokasi,
      lokasi_alamat: '',
      pic_user_id: '',
      lat: String(unitForm.value.lat || '-7.7956'),
      lng: String(unitForm.value.lng || '110.3695')
    }

    if (editingUnitId.value) {
      await api.put(`/organization/unit/${editingUnitId.value}`, payload)
      alert('Unit latihan berhasil diperbarui!')
    } else {
      await api.post(`/organization/unit`, payload)
      alert('Unit latihan baru berhasil ditambahkan!')
    }
    // Refresh unit list
    const res = await api.get(`/organization/cabang/${selectedCabang.value.id}/unit`)
    unitList.value = res || []
    
    showUnitModal.value = false
    selectedUnit.value = null
  } catch (e: any) {
    alert(e.response?.data?.message || 'Gagal menyimpan unit')
  }
}

const deleteUnit = async (id: string) => {
  if (confirm('Apakah Anda yakin ingin menghapus unit latihan ini?')) {
    try {
      await api.delete(`/organization/unit/${id}`)
      const res = await api.get(`/organization/cabang/${selectedCabang.value.id}/unit`)
      unitList.value = res || []
    } catch (e: any) {
      alert(e.response?.data?.message || 'Gagal menghapus unit')
    }
  }
}

// Staff Actions
const openAddPengurusModal = () => {
  editingPengurusId.value = null
  pengurusForm.value = { id: '', nama: '', jabatan: '', tingkatan: '', kontak: '', nomor: '' }
  showPengurusModal.value = true
}

const openEditPengurusModal = (p: any) => {
  editingPengurusId.value = p.id
  pengurusForm.value = {
    id: p.id,
    nama: p.nama,
    jabatan: p.jabatan,
    tingkatan: p.tingkatan,
    kontak: p.kontak,
    nomor: p.nomor || ''
  }
  showPengurusModal.value = true
}

const savePengurus = async () => {
  if (!selectedCabang.value) return
  submitting.value = true
  try {
    if (editingPengurusId.value) {
      await api.put(`/organization/pengurus/${editingPengurusId.value}`, pengurusForm.value)
      alert('Informasi pengurus berhasil diperbarui!')
    } else {
      await api.post(`/organization/cabang/${selectedCabang.value.id}/pengurus`, pengurusForm.value)
      alert('Pengurus baru berhasil ditambahkan!')
    }
    await fetchPengurusForCabang(selectedCabang.value.id)
    showPengurusModal.value = false
  } catch (error: any) {
    alert(error.response?.data?.message || 'Gagal menyimpan data pengurus')
  } finally {
    submitting.value = false
  }
}

const deletePengurus = async (id: string) => {
  if (confirm('Apakah Anda yakin ingin menghapus pengurus ini dari kepengurusan cabang?')) {
    try {
      await api.delete(`/organization/pengurus/${id}`)
      if (selectedCabang.value) {
        await fetchPengurusForCabang(selectedCabang.value.id)
      }
    } catch (error: any) {
      alert(error.response?.data?.message || 'Gagal menghapus pengurus')
    }
  }
}

// Trainer Actions
const openAddTrainerModal = () => {
  editingTrainerId.value = null
  trainerForm.value = {
    id: '',
    nama: '',
    nomor: '',
    jenis: 'Pelatih Daerah',
    tingkatan: 'PH Jurus 6',
    jarak: 'Jarak Dekat',
    transport: 30000,
    tanggal_mulai: new Date().toISOString().substring(0, 10),
    sk: '',
    status: 'Aktif'
  }
  showTrainerModal.value = true
}

const openEditTrainerModal = (t: any) => {
  editingTrainerId.value = t.id
  trainerForm.value = {
    id: t.id,
    nama: t.nama,
    nomor: t.nomor,
    jenis: t.jenis || 'Pelatih Daerah',
    tingkatan: t.tingkatan || 'PH Jurus 6',
    jarak: t.jarak || 'Jarak Dekat',
    transport: t.transport || 30000,
    tanggal_mulai: '',
    sk: '',
    status: t.status || 'Aktif'
  }
  showTrainerModal.value = true
}

const onJarakChange = () => {
  if (trainerForm.value.jarak === 'Sukarela') trainerForm.value.transport = 0
  else if (trainerForm.value.jarak === 'Jarak Dekat') trainerForm.value.transport = 30000
  else if (trainerForm.value.jarak === 'Jarak Sedang') trainerForm.value.transport = 50000
  else if (trainerForm.value.jarak === 'Jarak Jauh') trainerForm.value.transport = 75000
}

const saveTrainer = async () => {
  if (!selectedCabang.value) return
  try {
    const payload = {
      cabang_id: selectedCabang.value.id,
      nama: trainerForm.value.nama,
      nomor: trainerForm.value.nomor,
      jenis: trainerForm.value.jenis,
      tingkatan: trainerForm.value.tingkatan,
      jarak: trainerForm.value.jarak,
      transport: trainerForm.value.transport,
      tanggal_mulai: trainerForm.value.tanggal_mulai,
      sk: trainerForm.value.sk,
      status: trainerForm.value.status
    }

    if (editingTrainerId.value) {
      await api.put(`/organization/pelatih/${editingTrainerId.value}`, payload)
      alert('Data pelatih berhasil diperbarui!')
    } else {
      await api.post(`/organization/pelatih`, payload)
      alert('Pelatih baru berhasil ditambahkan!')
    }
    const res = await api.get(`/organization/pelatih?cabang_id=${selectedCabang.value.id}`)
    trainerList.value = res || []
    
    showTrainerModal.value = false
  } catch (e: any) {
    alert(e.response?.data?.message || 'Gagal menyimpan pelatih')
  }
}

const deleteTrainer = async (id: string) => {
  if (confirm('Apakah Anda yakin ingin menghapus pelatih ini?')) {
    try {
      await api.delete(`/organization/pelatih/${id}`)
      const res = await api.get(`/organization/pelatih?cabang_id=${selectedCabang.value.id}`)
      trainerList.value = res || []
    } catch (e: any) {
      alert(e.response?.data?.message || 'Gagal menghapus pelatih')
    }
  }
}

// PIC Koordinator Actions
const openGantiPICModal = (unit: any) => {
  selectedUnitPic.value = unit.pic || 'Belum diatur'
  picFormName.value = ''
  showGantiPICModal.value = true
}

const savePIC = () => {
  if (selectedUnit.value) {
    selectedUnit.value.pic = picFormName.value
    // Also update in unit list
    const idx = unitList.value.findIndex(x => x.id === selectedUnit.value.id)
    if (idx !== -1) {
      unitList.value[idx].pic = picFormName.value
    }
    showGantiPICModal.value = false
    alert('Koordinator unit latihan berhasil diganti!')
  }
}

// Helpers
const getInitials = (name: string) => {
  if (!name) return 'SA'
  const parts = name.split(' ')
  if (parts.length > 1) return (parts[0][0] + parts[1][0]).toUpperCase()
  return parts[0].substring(0, 2).toUpperCase()
}

const formatRupiah = (val: number) => {
  return new Intl.NumberFormat('id-ID').format(val)
}

let mapInstance: any = null
let markerInstance: any = null

const initMap = async (mapId: string, formRef: any) => {
  if (!process.client) return
  const L = (await import('leaflet')).default

  nextTick(() => {
    const el = document.getElementById(mapId)
    if (!el) return
    if (mapInstance) {
      mapInstance.remove()
      mapInstance = null
    }

    const initialLat = parseFloat(formRef.value.lat) || -7.7956
    const initialLng = parseFloat(formRef.value.lng) || 110.3695

    mapInstance = L.map(mapId).setView([initialLat, initialLng], 13)
    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png').addTo(mapInstance)

    delete (L.Icon.Default.prototype as any)._getIconUrl
    L.Icon.Default.mergeOptions({
      iconRetinaUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.9.4/images/marker-icon-2x.png',
      iconUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.9.4/images/marker-icon.png',
      shadowUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.9.4/images/marker-shadow.png',
    })

    markerInstance = L.marker([initialLat, initialLng], { draggable: true }).addTo(mapInstance)

    const updateForm = (pos: any) => {
      formRef.value.lat = pos.lat.toFixed(6)
      formRef.value.lng = pos.lng.toFixed(6)
    }

    markerInstance.on('dragend', (e: any) => updateForm(e.target.getLatLng()))
    mapInstance.on('click', (e: any) => {
      markerInstance.setLatLng(e.latlng)
      updateForm(e.latlng)
    })
    
    // Fix resize issue in modals using ResizeObserver
    const ro = new ResizeObserver(() => {
      if (mapInstance) mapInstance.invalidateSize()
    })
    ro.observe(el)

    // Fallback timeouts just in case
    setTimeout(() => { if (mapInstance) mapInstance.invalidateSize() }, 200)
    setTimeout(() => { if (mapInstance) mapInstance.invalidateSize() }, 600)
  })
}

onMounted(() => {
  fetchCabang()
  fetchProvinsi()
})
</script>

<style scoped>
@import "leaflet/dist/leaflet.css";

.cabang-split-layout {
  display: flex;
  height: calc(100vh - 64px - 44px);
  background: var(--bg);
  overflow: hidden;
  margin: -16px;
}

/* Left list area */
.panel-left {
  width: 300px;
  background: var(--card);
  border-right: 1px solid var(--border);
  display: flex;
  flex-direction: column;
  height: 100%;
  flex-shrink: 0;
}
.pl-head { padding: 16px 20px; border-bottom: 1px solid var(--border); }
.pl-title { font-size: 14px; font-weight: 800; }
.btn-xs { padding: 4px 8px; font-size: 11px; }

.form-group-filter { margin-bottom: 8px; margin-top: 8px; }
.filter-select { width: 100%; padding: 6px 10px; border: 1px solid var(--border2); border-radius: var(--r8); font-size: 12px; outline: none; background: #fff; }

.search-box { display: flex; align-items: center; gap: 8px; background: var(--surface); border: 1px solid var(--border); border-radius: var(--r8); padding: 6px 12px; }
.search-box i { color: var(--text3); }
.search-box input { border: none; outline: none; background: none; font-size: 11px; flex: 1; color: var(--text1); }

.pl-list { flex: 1; overflow-y: auto; padding: 10px 0; }
.prov-group { font-size: 10px; font-weight: 800; color: var(--text3); text-transform: uppercase; padding: 8px 20px 4px 20px; letter-spacing: .06em; }

.cabang-row { display: flex; align-items: center; gap: 10px; padding: 10px 20px; cursor: pointer; border-left: 3px solid transparent; transition: all .15s; }
.cabang-row:hover { background: var(--surface); }
.cabang-row.active { background: var(--hijau3); border-left-color: var(--hijau); }
.cr-dot { width: 8px; height: 8px; border-radius: 50%; flex-shrink: 0; }
.cr-info { flex: 1; min-width: 0; }
.cr-name { font-size: 12px; font-weight: 700; color: var(--text1); white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.cr-meta { font-size: 10px; color: var(--text3); margin-top: 2px; }
.cr-badge { font-size: 9px; font-weight: 700; background: var(--surface); color: var(--text2); padding: 2px 6px; border-radius: 10px; flex-shrink: 0; }
.cabang-row.active .cr-badge { background: var(--card); color: var(--hijau); }

/* Right panel area */
.panel-right { flex: 1; display: flex; flex-direction: column; height: 100%; min-width: 0; background: var(--surface); }
.panel-right.empty { align-items: center; justify-content: center; color: var(--text3); gap: 10px; padding: 40px; }

.pr-topbar { display: flex; justify-content: space-between; align-items: center; padding: 14px 20px; border-bottom: 1px solid var(--border); background: var(--card); flex-shrink: 0; }
.pr-breadcrumb { display: flex; align-items: center; gap: 6px; font-size: 12px; }
.pr-bc-link { color: var(--hijau); cursor: pointer; font-weight: 600; }
.pr-bc-sep { color: var(--text3); display: flex; align-items: center; }
.pr-bc-cur { color: var(--text1); font-weight: 700; }

.tabs-bar { display: flex; gap: 4px; border-bottom: 1px solid var(--border); background: var(--card); padding: 0 20px; flex-shrink: 0; }
.tab { padding: 12px 16px; font-size: 12px; font-weight: 600; color: var(--text2); border: none; background: none; cursor: pointer; border-bottom: 2px solid transparent; display: flex; align-items: center; gap: 4px; }
.tab.active { color: var(--hijau); border-bottom-color: var(--hijau); }
.tab-count { background: var(--hijau3); color: var(--hijau); border-radius: 10px; padding: 1px 5px; font-size: 9px; font-weight: 700; }

.pr-body { flex: 1; overflow-y: auto; padding: 20px; }
.tab-content-pane { display: flex; flex-direction: column; gap: 16px; }

.card { background: var(--card); border: 1px solid var(--border); border-radius: var(--r12); overflow: hidden; }
.card-head { display: flex; justify-content: space-between; align-items: center; padding: 12px 16px; border-bottom: 1px solid var(--border); }
.card-title { font-size: 12px; font-weight: 700; display: flex; align-items: center; gap: 6px; }
.card-title i { color: var(--hijau); font-size: 14px; }

.info-grid { display: grid; grid-template-columns: repeat(2, 1fr); gap: 14px; padding: 16px; }
.info-item { display: flex; flex-direction: column; gap: 4px; }
.info-label { font-size: 10px; color: var(--text3); text-transform: uppercase; font-weight: 600; letter-spacing: 0.5px; }
.info-val { font-size: 12px; font-weight: 600; color: var(--text1); }
.info-val.monospace { font-family: monospace; }

.bdg { display: inline-flex; align-items: center; gap: 4px; padding: 2px 8px; border-radius: 10px; font-size: 10px; font-weight: 700; }
.bdg-g { background: var(--hijau3); color: var(--hijau); }
.bdg-r { background: #fde8e8; color: var(--merah); }
.bdg-p { background: #f0e0fb; color: #6b21a8; }
.bdg-b { background: #e0f5fb; color: var(--biru); }
.bdg-y { background: #fff8e0; color: #a07000; }

/* Table styling inside cards */
.pg-table, .unit-table { width: 100%; border-collapse: collapse; }
.pg-table th, .unit-table th { padding: 10px 16px; font-size: 10px; font-weight: 700; color: var(--text3); text-transform: uppercase; background: var(--surface); border-bottom: 1px solid var(--border); text-align: left; }
.pg-table td, .unit-table td { padding: 12px 16px; font-size: 12px; border-bottom: 1px solid var(--border); vertical-align: middle; }
.pg-table tr:last-child td, .unit-table tr:last-child td { border-bottom: none; }
.unit-table tr:hover td { background: var(--surface); }

.pg-avatar-wrap { display: flex; align-items: center; gap: 8px; }
.pg-avatar { width: 32px; height: 32px; border-radius: 50%; background: var(--surface); display: flex; align-items: center; justify-content: center; font-size: 11px; font-weight: 700; color: var(--text2); border: 1px solid var(--border); }
.pg-name { font-size: 12px; font-weight: 600; }
.pg-nomor { font-size: 10px; color: var(--text3); font-family: monospace; }
.pg-contact { font-size: 11px; color: var(--text3); }

.action-btns { display: flex; gap: 4px; }
.icon-btn-sm { width: 24px; height: 24px; border: 1px solid var(--border); border-radius: var(--r6); background: #fff; cursor: pointer; display: inline-flex; align-items: center; justify-content: center; font-size: 11px; color: var(--text2); transition: all 0.15s; }
.icon-btn-sm:hover { border-color: var(--hijau); color: var(--hijau); }
.icon-btn-sm.danger:hover { border-color: #fca5a5; color: var(--merah); background: #fdecea; }

/* Unit row styles */
.unit-name-cell { display: flex; align-items: center; gap: 8px; }
.unit-icon { width: 28px; height: 28px; border-radius: var(--r6); background: var(--hijau3); color: var(--hijau); display: flex; align-items: center; justify-content: center; font-size: 14px; }
.unit-name { font-size: 12px; font-weight: 600; }
.unit-loc { font-size: 10px; color: var(--text3); display: flex; align-items: center; gap: 3px; }
.jadwal-chips { display: flex; flex-direction: column; gap: 3px; }
.jadwal-chip { display: inline-flex; align-items: center; gap: 4px; padding: 2px 6px; border-radius: 4px; font-size: 9px; background: var(--surface); border: 0.5px solid var(--border); color: var(--text2); font-weight: 600; }
.pic-cell { display: flex; align-items: center; gap: 6px; }
.pic-av { width: 24px; height: 24px; border-radius: 50%; background: var(--surface); display: flex; align-items: center; justify-content: center; font-size: 9px; font-weight: 700; color: var(--text2); border: 0.5px solid var(--border); }

/* Unit detail sub-pane styling */
.unit-detail-pane { display: flex; flex-direction: column; gap: 16px; }
.stat-mini-grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 12px; }
.stat-mini { background: var(--card); border: 1px solid var(--border); border-radius: var(--r12); padding: 14px; text-align: center; }
.stat-mini-val { font-size: 20px; font-weight: 800; line-height: 1; }
.stat-mini-label { font-size: 10px; color: var(--text3); margin-top: 4px; text-transform: uppercase; font-weight: 600; }

.pic-highlight-box { display: flex; align-items: center; gap: 12px; padding: 12px; background: var(--hijau3); border-radius: var(--r8); border: 0.5px solid var(--hijau4); }
.pic-big-av { width: 40px; height: 40px; border-radius: 50%; background: var(--hijau); display: flex; align-items: center; justify-content: center; font-size: 14px; font-weight: 700; color: #fff; }
.pic-contacts { text-align: right; font-size: 11px; color: var(--text2); line-height: 1.4; }

.export-btn { display: inline-flex; align-items: center; gap: 5px; border: 1px solid var(--border); padding: 6px 12px; border-radius: var(--r8); font-size: 11px; font-weight: 600; cursor: pointer; background: #fff; color: var(--text2); }
.export-btn:hover { border-color: var(--hijau); color: var(--hijau); }

/* Modals */
.modal-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.5); z-index: 1000; display: flex; align-items: center; justify-content: center; padding: 20px; }
.modal-card { background: var(--card); border-radius: var(--r12); width: 100%; max-width: 480px; padding: 20px; box-shadow: var(--shadow-lg); }
.modal-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px; }
.modal-title { font-size: 14px; font-weight: 700; }
.modal-close { background: none; border: none; font-size: 16px; cursor: pointer; color: var(--text3); }
.modal-form { display: flex; flex-direction: column; gap: 12px; }
.form-group { display: flex; flex-direction: column; gap: 4px; }
.form-label { font-size: 11px; font-weight: 600; color: var(--text2); }
.form-input, .form-select { padding: 8px 12px; border: 1.5px solid var(--border2); border-radius: var(--r8); font-size: 12px; outline: none; background: #fff; width: 100%; }
.form-input:focus, .form-select:focus { border-color: var(--hijau); }
.textarea { min-height: 70px; resize: vertical; font-family: inherit; }
.modal-footer { display: flex; justify-content: flex-end; gap: 8px; margin-top: 12px; }

.loading-state-list, .empty-state-list { text-align: center; padding: 20px; font-size: 12px; color: var(--text3); }
.spin { animation: spin .8s linear infinite; }
@keyframes spin { from { transform: rotate(0deg); } to { transform: rotate(360deg); } }
</style>

<style scoped>
.action-btn-styled {
  display: inline-flex; align-items: center; gap: 4px; padding: 4px 10px; border-radius: 6px; font-size: 11px; font-weight: 600; cursor: pointer; border: none; transition: all 0.2s;
}
.btn-edit-styled { background: #eff6ff; color: #3b82f6; }
.btn-edit-styled:hover { background: #dbeafe; }
.btn-delete-styled { background: #fef2f2; color: #ef4444; }
.btn-delete-styled:hover { background: #fee2e2; }
.btn-view-styled { background: #f0fdf4; color: #22c55e; }
.btn-view-styled:hover { background: #dcfce7; }
</style>
