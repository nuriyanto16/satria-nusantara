<template>
  <div id="page-iuran" class="iuran-split-layout">
    <!-- ===== PANEL KIRI: FILTER, TARGET, SUMMARY & REKENING ===== -->
    <div class="panel-left">
      <div class="pl-head">
        <div class="pl-title">Filter & Status</div>
        <div class="pl-sub">Pilih cabang & unit latihan</div>
      </div>
      <div class="pl-scroll">
        <!-- Pilih Cabang -->
        <div class="form-group-filter">
          <label class="filter-label">Cabang</label>
          <select v-model="selectedCabangId" class="filter-select" @change="onCabangChange">
            <option v-for="c in cabangList" :key="c.id" :value="c.id">{{ c.nama }}</option>
          </select>
        </div>

        <!-- Pilih Unit -->
        <div class="form-group-filter">
          <label class="filter-label">Unit Latihan</label>
          <select v-model="selectedUnitId" class="filter-select" @change="onUnitChange">
            <option v-for="u in unitList" :key="u.id" :value="u.id">{{ u.nama }}</option>
          </select>
        </div>

        <!-- Periode Bulan -->
        <div class="form-group-filter">
          <label class="filter-label">Periode</label>
          <div class="month-nav">
            <button class="month-btn" @click="changeMonth(-1)"><i class="ti ti-chevron-left"></i></button>
            <div class="month-label">{{ currentMonthLabel }}</div>
            <button class="month-btn" @click="changeMonth(1)"><i class="ti ti-chevron-right"></i></button>
          </div>
        </div>

        <!-- Target Box -->
        <div class="target-box">
          <div class="tb-label">Terkumpul bulan ini</div>
          <div class="tb-val">Rp {{ formatRupiah(totalTerkumpul) }}</div>
          <div class="tb-sub">{{ lunasCount }} dari {{ totalAnggotaUnit }} anggota · {{ currentMonthLabel }}</div>
          <div class="tb-bar-wrap">
            <div class="tb-bar" :style="{ width: percentCollected + '%', background: getPercentColor(percentCollected) }"></div>
          </div>
          <div style="display:flex;justify-content:space-between;margin-top:6px;">
            <div class="tb-pct" :style="{ color: getPercentColor(percentCollected) }">{{ percentCollected }}%</div>
            <div class="tb-target">Target: Rp {{ formatRupiah(targetBulanan) }}</div>
          </div>
        </div>

        <!-- Summary -->
        <div class="summary-section">
          <div class="sum-title">Ringkasan · {{ activeUnitName }}</div>
          <div class="summary-details">
            <div class="sum-row"><span class="sum-key">Anggota aktif</span><span class="sum-val">{{ totalAnggotaUnit }}</span></div>
            <div class="sum-row"><span class="sum-key">Sudah bayar</span><span class="sum-val green">{{ lunasCount }}</span></div>
            <div class="sum-row"><span class="sum-key">Belum bayar</span><span class="sum-val red">{{ belumCount }}</span></div>
            <div class="sum-row"><span class="sum-key">Nominal BLBA</span><span class="sum-val">Rp {{ formatRupiah(nominalBLBA) }}</span></div>
            <div class="sum-row"><span class="sum-key">Total terkumpul</span><span class="sum-val green">Rp {{ formatRupiah(totalTerkumpul) }}</span></div>
            <div class="sum-row"><span class="sum-key">Kurang dari target</span><span class="sum-val red">Rp {{ formatRupiah(kurangDariTarget) }}</span></div>
          </div>
        </div>

        <!-- Rekening Box -->
        <div class="rekening-box">
          <div class="rek-title">
            Rekening BLBA Unit
            <span class="rek-edit" @click="showRekeningModal = true"><i class="ti ti-edit"></i> Edit</span>
          </div>
          <div class="rek-row"><span class="rek-key">Bank</span><span class="rek-val">{{ rekening.bank }}</span></div>
          <div class="rek-row"><span class="rek-key">No. rekening</span><span class="rek-val monospace">{{ rekening.nomor }}</span></div>
          <div class="rek-row"><span class="rek-key">Atas nama</span><span class="rek-val">{{ rekening.atasNama }}</span></div>
          <div class="rek-row"><span class="rek-key">Diset oleh</span><span class="rek-val" style="color:var(--text3);">PIC Unit</span></div>
        </div>
      </div>
    </div>

    <!-- ===== PANEL KANAN: TABS & DATA LIST ===== -->
    <div class="panel-right">
      <div class="tabs-header">
        <button :class="['tab', { active: activeTab === 'anggota' }]" @click="activeTab = 'anggota'">
          <i class="ti ti-users"></i> Per Anggota
        </button>
        <button :class="['tab', { active: activeTab === 'unit' }]" @click="activeTab = 'unit'">
          <i class="ti ti-home"></i> Rekap Per Unit
        </button>
        <button :class="['tab', { active: activeTab === 'transaksi' }]" @click="activeTab = 'transaksi'" style="position:relative">
          <i class="ti ti-arrows-exchange"></i> Transaksi Terbaru
          <span v-if="pendingCount > 0" class="tab-badge">{{ pendingCount }}</span>
        </button>
      </div>

      <!-- TAB 1: PER ANGGOTA -->
      <div v-if="activeTab === 'anggota'" class="tab-container">
        <div class="toolbar">
          <div style="display:flex;gap:12px;align-items:center;">
            <select v-model="filterStatus" class="filter-select-sm">
              <option value="">Semua anggota</option>
              <option value="lunas">Sudah bayar</option>
              <option value="belum">Belum bayar</option>
            </select>
            <div class="toolbar-info">{{ activeUnitName }} · {{ currentMonthLabel }}</div>
          </div>
          <div style="display:flex;gap:7px;">
            <button class="export-btn" @click="exportData('PDF')"><i class="ti ti-file-type-pdf"></i> PDF</button>
            <button class="export-btn" @click="exportData('Excel')"><i class="ti ti-file-spreadsheet"></i> Excel</button>
          </div>
        </div>

        <div class="list-scroll">
          <div v-if="loading" style="padding: 20px; text-align: center;">
            <i class="ti ti-loader-2 spin" style="font-size:24px; color:var(--hijau)"></i>
            <div style="font-size:12px; color:var(--text3); margin-top:8px;">Memuat data iuran...</div>
          </div>
          <table v-else class="data-table">
            <thead>
              <tr>
                <th>Anggota</th>
                <th>Unit</th>
                <th>Status {{ currentMonthShort }}</th>
                <th>Tanggal Bayar</th>
                <th>Metode</th>
                <th>BLBA Berjalan</th>
                <th>Aksi</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="item in filteredAnggotaBLBA" :key="item.id" class="table-row-clickable" @click="viewHistory(item)">
                <td>
                  <div class="anggota-cell">
                    <div class="av" :style="{ background: getAvatarBg(item.nama) }">{{ getInitials(item.nama) }}</div>
                    <div>
                      <div class="member-name">{{ item.nama }}</div>
                      <div class="member-number">{{ item.nomor }}</div>
                    </div>
                  </div>
                </td>
                <td style="font-size:11px;">{{ item.unit }}</td>
                <td>
                  <span :class="['status-badge', item.status === 'lunas' ? 'lunas' : 'belum']">
                    <i :class="item.status === 'lunas' ? 'ti ti-circle-check' : 'ti ti-clock'"></i>
                    {{ item.status === 'lunas' ? 'Lunas' : 'Belum' }}
                  </span>
                </td>
                <td style="font-size:11px;color:var(--text2);">{{ item.tanggalBayar || '—' }}</td>
                <td>
                  <span v-if="item.status === 'lunas'" class="metode-badge">{{ item.metode }}</span>
                  <span v-else class="muted-text">—</span>
                </td>
                <td class="history-cell">{{ item.blbaBerjalan }} bln</td>
                <td @click.stopPropagation>
                  <div class="action-btns">
                    <button v-if="item.status === 'belum'" class="action-btn" title="Tandai Lunas" @click="setLunas(item)">
                      <i class="ti ti-check"></i>
                    </button>
                    <button class="action-btn" title="Riwayat" @click="viewHistory(item)">
                      <i class="ti ti-history"></i>
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <div class="list-footer">
          <div class="lf-info">
            Menampilkan {{ filteredAnggotaBLBA.length }} anggota ·
            <strong style="color:var(--hijau);">{{ lunasCount }} lunas</strong> ·
            <strong style="color:var(--merah);">{{ belumCount }} belum</strong>
          </div>
          <div style="display:flex;gap:7px;">
            <button class="export-btn" @click="exportData('PDF')"><i class="ti ti-file-type-pdf"></i> PDF</button>
            <button class="export-btn" @click="exportData('Excel')"><i class="ti ti-file-spreadsheet"></i> Excel</button>
          </div>
        </div>
      </div>

      <!-- TAB 2: REKAP PER UNIT -->
      <div v-else-if="activeTab === 'unit'" class="tab-container">
        <div class="toolbar">
          <div class="toolbar-info">Rekap BLBA Per Unit · {{ activeCabangName }} · {{ currentMonthLabel }}</div>
          <div style="display:flex;gap:7px;">
            <button class="export-btn" @click="exportData('PDF')"><i class="ti ti-file-type-pdf"></i> PDF</button>
            <button class="export-btn" @click="exportData('Excel')"><i class="ti ti-file-spreadsheet"></i> Excel</button>
          </div>
        </div>

        <div class="list-scroll">
          <table class="data-table">
            <thead>
              <tr>
                <th>Unit</th>
                <th>Aktif</th>
                <th>Lunas</th>
                <th>Belum</th>
                <th>Terkumpul</th>
                <th>Target</th>
                <th>Pencapaian</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="unit in rekapUnitList" :key="unit.name">
                <td>
                  <div style="font-weight:600;">{{ unit.name }}</div>
                  <div style="font-size:10px;color:var(--text3);">PIC: {{ unit.pic }}</div>
                </td>
                <td style="font-weight:700;">{{ unit.aktif }}</td>
                <td style="font-weight:700;color:var(--hijau);">{{ unit.lunas }}</td>
                <td style="font-weight:700;color:var(--merah);">{{ unit.belum }}</td>
                <td style="font-weight:700;">Rp {{ formatRupiah(unit.terkumpul) }}</td>
                <td style="color:var(--text2);">Rp {{ formatRupiah(unit.target) }}</td>
                <td>
                  <div class="unit-bar-cell">
                    <div class="ub-wrap"><div class="ub-fill" :style="{ width: unit.pencapaian + '%', background: getPercentColor(unit.pencapaian) }"></div></div>
                    <span style="font-size:11px;font-weight:700;" :style="{ color: getPercentColor(unit.pencapaian) }">{{ unit.pencapaian }}%</span>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <div class="branch-summary-bar">
          <div style="display:flex;gap:24px;align-items:center;">
            <div>
              <div style="font-size:10px;color:var(--text3);margin-bottom:2px;">Total Terkumpul Cabang</div>
              <div style="font-size:20px;font-weight:700;color:var(--hijau);">Rp {{ formatRupiah(rekapCabangTotal) }}</div>
            </div>
            <div style="border-left:1px solid var(--border);padding-left:24px;">
              <div style="font-size:10px;color:var(--text3);margin-bottom:2px;">Target Cabang</div>
              <div style="font-size:20px;font-weight:700;color:var(--text2);">Rp {{ formatRupiah(rekapCabangTarget) }}</div>
            </div>
            <div style="border-left:1px solid var(--border);padding-left:24px;">
              <div style="font-size:10px;color:var(--text3);margin-bottom:2px;">Pencapaian</div>
              <div style="font-size:20px;font-weight:700;color:var(--hijau);">{{ rekapCabangPct }}%</div>
            </div>
          </div>
          <div style="display:flex;gap:7px;">
            <button class="export-btn" @click="exportData('PDF')"><i class="ti ti-file-type-pdf"></i> PDF</button>
            <button class="export-btn" @click="exportData('Excel')"><i class="ti ti-file-spreadsheet"></i> Excel</button>
          </div>
        </div>
      </div>

      <!-- TAB 3: TRANSAKSI TERBARU -->
      <div v-else-if="activeTab === 'transaksi'" class="tab-container">
        <div class="toolbar">
          <div style="display:flex;gap:12px;align-items:center;">
            <select v-model="filterTrxStatus" class="filter-select-sm">
              <option value="">Semua transaksi</option>
              <option value="lunas">Berhasil</option>
              <option value="pending">Menunggu konfirmasi</option>
              <option value="ditolak">Ditolak</option>
            </select>
            <div class="toolbar-info">{{ filteredTrxList.length }} transaksi</div>
          </div>
          <div style="display:flex;gap:7px;align-items:center;">
            <div v-if="pendingCount > 0" class="pending-alert-chip">
              <i class="ti ti-clock-exclamation"></i>
              {{ pendingCount }} menunggu konfirmasi
            </div>
            <button class="export-btn" @click="exportData('PDF')"><i class="ti ti-file-type-pdf"></i> PDF</button>
          </div>
        </div>

        <div class="list-scroll">
          <div v-if="filteredTrxList.length === 0" class="empty-state">
            <i class="ti ti-inbox" style="font-size:36px;color:var(--text3);"></i>
            <div style="font-size:13px;color:var(--text3);margin-top:8px;">Belum ada transaksi</div>
          </div>
          <table v-else class="data-table">
            <thead>
              <tr>
                <th>Anggota</th>
                <th>Periode</th>
                <th>Metode</th>
                <th>Nominal</th>
                <th>Waktu</th>
                <th>Status</th>
                <th>Aksi</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="trx in filteredTrxList" :key="trx.id" :class="['table-row-clickable', trx.status === 'pending' ? 'row-pending' : '']">
                <td>
                  <div class="anggota-cell">
                    <div class="av" :style="{ background: getAvatarBg(trx.nama) }">{{ getInitials(trx.nama) }}</div>
                    <div>
                      <div class="member-name">{{ trx.nama }}</div>
                      <div class="member-number">{{ trx.nomor }}</div>
                    </div>
                  </div>
                </td>
                <td style="font-size:11px;color:var(--text2);">BLBA {{ trx.bulan }}</td>
                <td>
                  <div style="display:flex;flex-direction:column;gap:3px;">
                    <span class="metode-badge" :class="trx.metode.startsWith('Transfer') ? 'metode-bank' : ''">
                      <i :class="trx.metode.startsWith('Transfer') ? 'ti ti-building-bank' : 'ti ti-wallet'"></i>
                      {{ trx.metode }}
                    </span>
                    <span v-if="trx.buktiUrl" style="font-size:10px;color:var(--biru);cursor:pointer;" @click.stop="viewBukti(trx)"><i class="ti ti-paperclip"></i> Lihat bukti</span>
                  </div>
                </td>
                <td style="font-weight:700;color:var(--text1);">Rp {{ formatRupiah(trx.nominal) }}</td>
                <td style="font-size:11px;color:var(--text3);">
                  <div>{{ trx.waktu }}</div>
                  <div style="font-size:10px;color:var(--text3);">{{ trx.relativeTime }}</div>
                </td>
                <td>
                  <span :class="['status-badge', trx.status === 'lunas' ? 'lunas' : trx.status === 'pending' ? 'pending' : 'belum']">
                    <i :class="trx.status === 'lunas' ? 'ti ti-circle-check' : trx.status === 'pending' ? 'ti ti-clock' : 'ti ti-circle-x'"></i>
                    {{ trx.status === 'lunas' ? 'Berhasil' : trx.status === 'pending' ? 'Menunggu' : 'Ditolak' }}
                  </span>
                </td>
                <td @click.stop>
                  <div v-if="trx.status === 'pending'" class="action-btns">
                    <button class="action-btn green" title="Konfirmasi Lunas" @click="konfirmasiTrx(trx)">
                      <i class="ti ti-check"></i> Konfirmasi
                    </button>
                    <button class="action-btn red" title="Tolak" @click="tolakTrx(trx)">
                      <i class="ti ti-x"></i>
                    </button>
                  </div>
                  <span v-else style="font-size:11px;color:var(--text3);">—</span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <div class="list-footer">
          <div class="lf-info">
            {{ filteredTrxList.length }} transaksi ·
            <strong style="color:var(--hijau);">{{ trxBerhasil }} berhasil</strong> ·
            <strong style="color:var(--kuning);">{{ pendingCount }} pending</strong> ·
            <strong style="color:var(--merah);">{{ trxDitolak }} ditolak</strong>
          </div>
          <div style="font-size:11px;color:var(--text3);">
            Total: <strong style="color:var(--hijau);">Rp {{ formatRupiah(totalTrxBerhasil) }}</strong>
          </div>
        </div>
      </div>
    </div>

    <!-- MODAL: LIHAT BUKTI TRANSFER -->
    <div v-if="showBuktiModal && selectedTrx" class="modal-overlay" @click.self="showBuktiModal = false">
      <div class="modal-card mini">
        <div class="modal-header">
          <div class="modal-title">Bukti Transfer — {{ selectedTrx.nama }}</div>
          <button class="modal-close" @click="showBuktiModal = false"><i class="ti ti-x"></i></button>
        </div>
        <div class="modal-body">
          <div style="display:flex;flex-direction:column;gap:12px;">
            <div class="rekening-current-box">
              <div class="rcb-title">Detail Transaksi</div>
              <div class="rcb-desc">{{ selectedTrx.nama }} · {{ selectedTrx.nomor }}</div>
            </div>
            <div class="sum-row"><span class="sum-key">Periode</span><span class="sum-val">BLBA {{ selectedTrx.bulan }}</span></div>
            <div class="sum-row"><span class="sum-key">Metode</span><span class="sum-val">{{ selectedTrx.metode }}</span></div>
            <div class="sum-row"><span class="sum-key">Nominal</span><span class="sum-val green">Rp {{ formatRupiah(selectedTrx.nominal) }}</span></div>
            <div class="sum-row"><span class="sum-key">Waktu</span><span class="sum-val">{{ selectedTrx.waktu }}</span></div>
            <div style="background:var(--bg2);border-radius:10px;padding:12px;text-align:center;border:1px dashed var(--border);">
              <i class="ti ti-photo" style="font-size:48px;color:var(--text3);"></i>
              <div style="font-size:12px;color:var(--text3);margin-top:6px;">{{ selectedTrx.buktiUrl || 'bukti_transfer.jpg' }}</div>
              <div style="font-size:10px;color:var(--text3);margin-top:4px;">Preview tidak tersedia di mode demo</div>
            </div>
          </div>
        </div>
        <div class="modal-footer">
          <button v-if="selectedTrx.status === 'pending'" class="btn btn-outline red" style="width:auto" @click="tolakTrx(selectedTrx); showBuktiModal = false">
            <i class="ti ti-x"></i> Tolak
          </button>
          <button v-if="selectedTrx.status === 'pending'" class="btn btn-primary" style="width:auto" @click="konfirmasiTrx(selectedTrx); showBuktiModal = false">
            <i class="ti ti-check"></i> Konfirmasi Lunas
          </button>
          <button v-else class="btn btn-outline" style="width:auto" @click="showBuktiModal = false">Tutup</button>
        </div>
      </div>
    </div>

    <!-- MODAL: EDIT REKENING BLBA -->
    <div v-if="showRekeningModal" class="modal-overlay" @click.self="showRekeningModal = false">
      <div class="modal-card mini">
        <div class="modal-header">
          <div class="modal-title">Edit Rekening BLBA — {{ activeUnitName }}</div>
          <button class="modal-close" @click="showRekeningModal = false"><i class="ti ti-x"></i></button>
        </div>
        <div class="modal-body">
          <div class="alert-info-box">
            <i class="ti ti-info-circle"></i>
            <span>Rekening ini digunakan untuk pembayaran BLBA {{ activeUnitName }}. Perubahan berlaku mulai bulan depan.</span>
          </div>
          <div class="rekening-current-box">
            <div class="rcb-title">Rekening saat ini</div>
            <div class="rcb-desc"><strong>{{ rekening.bank }}</strong> · {{ rekening.nomor }} · {{ rekening.atasNama }}</div>
          </div>
          <div class="form-group">
            <label class="form-label">Bank</label>
            <select v-model="editRekening.bank" class="form-select">
              <option>BCA</option>
              <option>BRI</option>
              <option>BNI</option>
              <option>Mandiri</option>
              <option>BSI</option>
              <option>DANA</option>
              <option>OVO</option>
              <option>GoPay</option>
              <option>ShopeePay</option>
            </select>
          </div>
          <div class="form-group">
            <label class="form-label">Nomor Rekening / Akun</label>
            <input v-model="editRekening.nomor" type="text" class="form-input" placeholder="Masukkan nomor rekening..." />
          </div>
          <div class="form-group">
            <label class="form-label">Atas Nama</label>
            <input v-model="editRekening.atasNama" type="text" class="form-input" placeholder="Nama pemilik rekening..." />
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn-outline" style="width: auto" @click="showRekeningModal = false">Batal</button>
          <button class="btn btn-primary" style="width: auto" @click="saveRekening"><i class="ti ti-check"></i> Simpan</button>
        </div>
      </div>
    </div>

    <!-- MODAL: RIWAYAT PEMBAYARAN BLBA -->
    <div v-if="showRiwayatModal && selectedAnggota" class="modal-overlay" @click.self="showRiwayatModal = false">
      <div class="modal-card">
        <div class="modal-header">
          <div class="modal-title">Riwayat Pembayaran BLBA</div>
          <button class="modal-close" @click="showRiwayatModal = false"><i class="ti ti-x"></i></button>
        </div>
        <div class="modal-body">
          <div class="riwayat-header">
            <div class="av-large" :style="{ background: getAvatarBg(selectedAnggota.nama) }">
              {{ getInitials(selectedAnggota.nama) }}
            </div>
            <div>
              <div class="rh-name">{{ selectedAnggota.nama }}</div>
              <div class="rh-sub">{{ selectedAnggota.nomor }} · {{ selectedAnggota.unit }}</div>
              <div style="margin-top:6px;display:flex;gap:6px;">
                <span class="badge border-green">{{ selectedAnggota.blbaBerjalan }} bulan lunas</span>
                <span class="badge border-green">100% konsisten</span>
              </div>
            </div>
          </div>
          <div class="table-card" style="margin-top: 16px;">
            <table class="history-table">
              <thead>
                <tr>
                  <th>Bulan</th>
                  <th>Status</th>
                  <th>Tanggal Bayar</th>
                  <th>Nominal</th>
                  <th>Metode</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="h in selectedAnggota.riwayatList" :key="h.bulan">
                  <td style="font-weight: 700;">{{ h.bulan }}</td>
                  <td>
                    <span :class="['status-badge', h.status === 'Lunas' ? 'lunas' : 'belum']">
                      {{ h.status }}
                    </span>
                  </td>
                  <td style="font-size:11px;color:var(--text2);">{{ h.tanggal || '—' }}</td>
                  <td style="font-weight:600;">Rp {{ formatRupiah(nominalBLBA) }}</td>
                  <td>
                    <span v-if="h.status === 'Lunas'" class="metode-badge">{{ h.metode }}</span>
                    <span v-else class="muted-text">—</span>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn-outline" style="width: auto" @click="exportHistory"><i class="ti ti-download"></i> Export PDF</button>
          <button class="btn btn-primary" style="width: auto" @click="showRiwayatModal = false">Tutup</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({ title: 'BLBA' })

const api = useApi()

const activeTab = ref('anggota')
const filterStatus = ref('')
const filterTrxStatus = ref('')
const loading = ref(false)

// Month details
const currentMonthOffset = ref(0)
const baseDate = new Date(2026, 6, 1) // Base on July 2026

const currentMonthLabel = computed(() => {
  const d = new Date(baseDate.getFullYear(), baseDate.getMonth() + currentMonthOffset.value, 1)
  return d.toLocaleDateString('id-ID', { month: 'long', year: 'numeric' })
})

const currentMonthShort = computed(() => {
  const d = new Date(baseDate.getFullYear(), baseDate.getMonth() + currentMonthOffset.value, 1)
  return d.toLocaleDateString('id-ID', { month: 'short', year: 'numeric' })
})

const changeMonth = (val: number) => {
  currentMonthOffset.value += val
  generateIuranData()
}

// Selection & Dropdowns
const selectedCabangId = ref<string>('')
const selectedUnitId = ref<string>('')
const activeUnitName = ref('Unit Malioboro')
const activeCabangName = ref('Kota Yogyakarta')

const cabangList = ref<any[]>([])
const unitList = ref<any[]>([])

// Modals
const showRekeningModal = ref(false)
const showRiwayatModal = ref(false)
const selectedAnggota = ref<any>(null)

// Rekening Details
const rekening = ref({
  bank: 'BCA',
  nomor: '0881234567',
  atasNama: 'Unit Malioboro SN'
})

const editRekening = ref({
  bank: 'BCA',
  nomor: '0881234567',
  atasNama: 'Unit Malioboro SN'
})

const saveRekening = () => {
  rekening.value = { ...editRekening.value }
  showRekeningModal.value = false
  alert('Rekening pembayaran BLBA berhasil diperbarui!')
}

// Master rates
const nominalBLBA = ref(40000)

// Dynamic BLBA Member List
const anggotaBLBAList = ref<any[]>([])

// Transaksi list
const showBuktiModal = ref(false)
const selectedTrx = ref<any>(null)

const viewBukti = (trx: any) => {
  selectedTrx.value = trx
  showBuktiModal.value = true
}

const transaksiList = ref<any[]>([
  { id: 't1', nama: 'Budi Santoso', nomor: 'YO-YGY-00142', bulan: 'Juli 2026', metode: 'DANA', nominal: 40000, waktu: '22 Jul 2026, 08:14', relativeTime: 'Hari ini', status: 'lunas', buktiUrl: null },
  { id: 't2', nama: 'Sari Rahmawati', nomor: 'YO-YGY-00098', bulan: 'Juli 2026', metode: 'Transfer Bank - BCA', nominal: 40000, waktu: '22 Jul 2026, 07:52', relativeTime: 'Hari ini', status: 'pending', buktiUrl: 'bukti_transfer_1753145520123.jpg' },
  { id: 't3', nama: 'Hendra Kusuma', nomor: 'YO-YGY-00067', bulan: 'Juli 2026', metode: 'GoPay', nominal: 40000, waktu: '21 Jul 2026, 19:30', relativeTime: 'Kemarin', status: 'lunas', buktiUrl: null },
  { id: 't4', nama: 'Agus Prasetyo', nomor: 'YO-YGY-00201', bulan: 'Juli 2026', metode: 'Transfer Bank - BRI', nominal: 40000, waktu: '21 Jul 2026, 15:10', relativeTime: 'Kemarin', status: 'pending', buktiUrl: 'bukti_transfer_1753098610234.jpg' },
  { id: 't5', nama: 'Dwi Wahyuni', nomor: 'YO-YGY-00167', bulan: 'Juli 2026', metode: 'OVO', nominal: 40000, waktu: '20 Jul 2026, 10:05', relativeTime: '2 hari lalu', status: 'lunas', buktiUrl: null },
  { id: 't6', nama: 'Farid Nugroho', nomor: 'YO-YGY-00312', bulan: 'Juni 2026', metode: 'ShopeePay', nominal: 40000, waktu: '19 Jul 2026, 09:00', relativeTime: '3 hari lalu', status: 'lunas', buktiUrl: null },
  { id: 't7', nama: 'Nurul Rahayu', nomor: 'YO-YGY-00089', bulan: 'Juli 2026', metode: 'Transfer Bank - Mandiri', nominal: 40000, waktu: '18 Jul 2026, 14:22', relativeTime: '4 hari lalu', status: 'ditolak', buktiUrl: 'bukti_transfer_1752925320456.jpg' },
])

const filteredTrxList = computed(() => {
  if (!filterTrxStatus.value) return transaksiList.value
  return transaksiList.value.filter(t => t.status === filterTrxStatus.value)
})

const pendingCount = computed(() => transaksiList.value.filter(t => t.status === 'pending').length)
const trxBerhasil = computed(() => transaksiList.value.filter(t => t.status === 'lunas').length)
const trxDitolak = computed(() => transaksiList.value.filter(t => t.status === 'ditolak').length)
const totalTrxBerhasil = computed(() => transaksiList.value.filter(t => t.status === 'lunas').reduce((s, t) => s + t.nominal, 0))

const konfirmasiTrx = (trx: any) => {
  trx.status = 'lunas'
  trx.relativeTime = 'Baru dikonfirmasi'
  // Also update anggotaBLBAList if match found
  const member = anggotaBLBAList.value.find(a => a.nomor === trx.nomor)
  if (member) { member.status = 'lunas'; member.tanggalBayar = trx.waktu.split(',')[0]; member.metode = trx.metode }
  alert(`✅ Pembayaran ${trx.nama} untuk BLBA ${trx.bulan} berhasil dikonfirmasi!`)
}

const tolakTrx = (trx: any) => {
  if (!confirm(`Tolak pembayaran transfer dari ${trx.nama}?`)) return
  trx.status = 'ditolak'
  trx.relativeTime = 'Baru ditolak'
  alert(`❌ Pembayaran ${trx.nama} ditolak. Anggota akan diberitahu.`)
}

const fetchCabang = async () => {
  try {
    const res = await api.get('/organization/cabang?limit=100')
    cabangList.value = res.data || []
  } catch (e) {
    console.error(e)
    cabangList.value = []
  }

  if (cabangList.value.length === 0) {
    cabangList.value = [
      { id: '1', nama: 'Kota Yogyakarta' },
      { id: '2', nama: 'Kota Bandung' },
      { id: '3', nama: 'Kota Semarang' },
      { id: '4', nama: 'Kota Surabaya' }
    ]
  }

  selectedCabangId.value = cabangList.value[0].id
  activeCabangName.value = cabangList.value[0].nama
  await fetchUnits()
}

const fetchUnits = async () => {
  if (!selectedCabangId.value) return
  try {
    const res = await api.get(`/organization/cabang/${selectedCabangId.value}/unit`)
    unitList.value = res || []
  } catch (e) {
    console.error(e)
    unitList.value = []
  }

  if (unitList.value.length === 0) {
    unitList.value = [
      { id: '1', nama: 'Unit Malioboro', cabang_id: selectedCabangId.value },
      { id: '2', nama: 'Unit Kotagede', cabang_id: selectedCabangId.value },
      { id: '3', nama: 'Unit Gondokusuman', cabang_id: selectedCabangId.value },
      { id: '4', nama: 'Unit Mantrijeron', cabang_id: selectedCabangId.value }
    ]
  }

  selectedUnitId.value = unitList.value[0].id
  activeUnitName.value = unitList.value[0].nama
  await generateIuranData()
}

const onCabangChange = async () => {
  const match = cabangList.value.find(c => c.id == selectedCabangId.value)
  if (match) activeCabangName.value = match.nama
  await fetchUnits()
}

const onUnitChange = async () => {
  const match = unitList.value.find(u => u.id == selectedUnitId.value)
  if (match) activeUnitName.value = match.nama
  await generateIuranData()
}

const generateIuranData = async () => {
  loading.value = true
  try {
    let members: any[] = []
    try {
      const res = await api.get(`/organization/anggota?limit=100&unit_id=${selectedUnitId.value}&status=aktif`)
      members = res.data || []
    } catch (e) {
      console.warn('API lookup failed, falling back to mockup', e)
    }

    if (members.length === 0) {
      if (activeUnitName.value.includes('Kotagede')) {
        members = [
          { id: '8', nama_lengkap: 'Pratiwi Riyadi', nomor_anggota: 'YO-YGY-00287', unit_nama: activeUnitName.value },
          { id: '9', nama_lengkap: 'Rian Wijaya', nomor_anggota: 'YO-YGY-00288', unit_nama: activeUnitName.value },
          { id: '10', nama_lengkap: 'Siti Aminah', nomor_anggota: 'YO-YGY-00289', unit_nama: activeUnitName.value }
        ]
      } else if (activeUnitName.value.includes('Gondokusuman')) {
        members = [
          { id: '11', nama_lengkap: 'Sunaryo', nomor_anggota: 'YO-YGY-00312', unit_nama: activeUnitName.value },
          { id: '12', nama_lengkap: 'Edi Santoso', nomor_anggota: 'YO-YGY-00313', unit_nama: activeUnitName.value }
        ]
      } else if (activeUnitName.value.includes('Mantrijeron')) {
        members = [
          { id: '13', nama_lengkap: 'Haryono N.', nomor_anggota: 'YO-YGY-00398', unit_nama: activeUnitName.value },
          { id: '14', nama_lengkap: 'Lilis Karlina', nomor_anggota: 'YO-YGY-00399', unit_nama: activeUnitName.value }
        ]
      } else {
        members = [
          { id: '1', nama_lengkap: 'Budi Santoso', nomor_anggota: 'YO-YGY-00142', unit_nama: activeUnitName.value },
          { id: '2', nama_lengkap: 'Sari Rahmawati', nomor_anggota: 'YO-YGY-00098', unit_nama: activeUnitName.value },
          { id: '3', nama_lengkap: 'Hendra Kusuma', nomor_anggota: 'YO-YGY-00067', unit_nama: activeUnitName.value },
          { id: '4', nama_lengkap: 'Agus Prasetyo', nomor_anggota: 'YO-YGY-00201', unit_nama: activeUnitName.value },
          { id: '5', nama_lengkap: 'Dwi Wahyuni', nomor_anggota: 'YO-YGY-00167', unit_nama: activeUnitName.value },
          { id: '6', nama_lengkap: 'Farid Nugroho', nomor_anggota: 'YO-YGY-00312', unit_nama: activeUnitName.value },
          { id: '7', nama_lengkap: 'Nurul Rahayu', nomor_anggota: 'YO-YGY-00089', unit_nama: activeUnitName.value }
        ]
      }
    }

    const todayStr = new Date().toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' })

    anggotaBLBAList.value = members.map((m, idx) => {
      const isLunas = idx % 3 !== 2
      return {
        id: m.id,
        nama: m.nama_lengkap,
        nomor: m.nomor_anggota || `YO-YGY-${100 + idx}`,
        unit: m.unit_nama || activeUnitName.value,
        status: isLunas ? 'lunas' : 'belum',
        tanggalBayar: isLunas ? `${idx + 1} ${currentMonthShort.value}` : null,
        metode: isLunas ? ['DANA', 'GoPay', 'ShopeePay', 'OVO'][idx % 4] : '—',
        blbaBerjalan: isLunas ? 15 : 14,
        riwayatList: [
          { bulan: `Jun ${baseDate.getFullYear()}`, status: 'Lunas', tanggal: `3 Jun ${baseDate.getFullYear()}`, metode: 'DANA' },
          { bulan: `Mei ${baseDate.getFullYear()}`, status: 'Lunas', tanggal: `1 Mei ${baseDate.getFullYear()}`, metode: 'DANA' },
          { bulan: `Apr ${baseDate.getFullYear()}`, status: 'Lunas', tanggal: `2 Apr ${baseDate.getFullYear()}`, metode: 'GoPay' },
          { bulan: `Mar ${baseDate.getFullYear()}`, status: 'Lunas', tanggal: `5 Mar ${baseDate.getFullYear()}`, metode: 'GoPay' },
          { bulan: `Feb ${baseDate.getFullYear()}`, status: 'Lunas', tanggal: `3 Feb ${baseDate.getFullYear()}`, metode: 'OVO' }
        ]
      }
    })
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

const filteredAnggotaBLBA = computed(() => {
  if (!filterStatus.value) return anggotaBLBAList.value
  return anggotaBLBAList.value.filter(a => a.status === filterStatus.value)
})

// Summary numbers
const totalAnggotaUnit = computed(() => anggotaBLBAList.value.length)
const lunasCount = computed(() => anggotaBLBAList.value.filter(a => a.status === 'lunas').length)
const belumCount = computed(() => totalAnggotaUnit.value - lunasCount.value)

const totalTerkumpul = computed(() => lunasCount.value * nominalBLBA.value)
const targetBulanan = computed(() => totalAnggotaUnit.value * nominalBLBA.value)
const kurangDariTarget = computed(() => targetBulanan.value - totalTerkumpul.value)

const percentCollected = computed(() => {
  if (totalAnggotaUnit.value === 0) return 0
  return Math.round((lunasCount.value / totalAnggotaUnit.value) * 100)
})

// Set member status to Lunas
const setLunas = (item: any) => {
  item.status = 'lunas'
  const today = new Date()
  item.tanggalBayar = today.toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' })
  item.metode = 'QRIS'
  item.blbaBerjalan += 1
  alert(`Pembayaran BLBA untuk ${item.nama} telah diverifikasi Lunas!`)
}

// Modal View History
const viewHistory = (item: any) => {
  selectedAnggota.value = item
  showRiwayatModal.value = true
}

// Tab 2 Rekap per Unit
const rekapUnitList = computed(() => {
  return [
    { name: 'Unit Malioboro', pic: 'Budi Hartono', aktif: totalAnggotaUnit.value, lunas: lunasCount.value, belum: belumCount.value, terkumpul: totalTerkumpul.value, target: targetBulanan.value, pencapaian: percentCollected.value },
    { name: 'Unit Kotagede', pic: 'Sari Wulandari', aktif: 87, lunas: 87, belum: 0, terkumpul: 3480000, target: 3480000, pencapaian: 100 },
    { name: 'Unit Gondokusuman', pic: 'Hendra K.', aktif: 64, lunas: 51, belum: 13, terkumpul: 2040000, target: 2560000, pencapaian: 80 },
    { name: 'Unit Mantrijeron', pic: 'Dewi N.', aktif: 52, lunas: 31, belum: 21, terkumpul: 1240000, target: 2080000, pencapaian: 60 }
  ]
})

const rekapCabangTotal = computed(() => {
  return rekapUnitList.value.reduce((acc, u) => acc + u.terkumpul, 0)
})
const rekapCabangTarget = computed(() => {
  return rekapUnitList.value.reduce((acc, u) => acc + u.target, 0)
})
const rekapCabangPct = computed(() => {
  if (rekapCabangTarget.value === 0) return 0
  return Math.round((rekapCabangTotal.value / rekapCabangTarget.value) * 100)
})

// Visual formatting helpers
const formatRupiah = (val: number) => {
  return val.toString().replace(/\B(?=(\d{3})+(?!\d))/g, '.')
}

const getPercentColor = (pct: number) => {
  if (pct >= 85) return 'var(--hijau)'
  if (pct >= 60) return 'var(--kuning)'
  return 'var(--merah)'
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

const getInitials = (name: string) => {
  if (!name) return 'AN'
  const parts = name.split(' ')
  if (parts.length > 1) {
    return (parts[0][0] + parts[1][0]).toUpperCase()
  }
  return parts[0].substring(0, 2).toUpperCase()
}

const exportData = (type: string) => {
  alert(`Data pembayaran BLBA berhasil diekspor ke format ${type}!`)
}

const exportHistory = () => {
  alert('Riwayat pembayaran berhasil diekspor!')
}

onMounted(() => {
  fetchCabang()
})
</script>

<style scoped>
.iuran-split-layout {
  display: flex;
  height: calc(100vh - 64px - 44px);
  background: var(--bg);
  overflow: hidden;
  margin: -16px;
}

/* Left panel style */
.panel-left {
  width: 320px;
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

.form-group-filter { margin-bottom: 16px; }
.filter-label { font-size: 11px; font-weight: 600; color: var(--text2); margin-bottom: 5px; display: block; text-transform: uppercase; letter-spacing: 0.5px; }
.filter-select { width: 100%; padding: 7px 10px; border: 1px solid var(--border2); border-radius: var(--r8); font-size: 12px; outline: none; background: #fff; color: var(--text1); }

/* Month navigate button style */
.month-nav { display: flex; align-items: center; border: 1px solid var(--border2); border-radius: var(--r8); background: #fff; overflow: hidden; }
.month-btn { width: 32px; height: 32px; border: none; background: none; cursor: pointer; color: var(--text2); display: flex; align-items: center; justify-content: center; font-size: 14px; }
.month-btn:hover { background: var(--surface); color: var(--hijau); }
.month-label { flex: 1; text-align: center; font-size: 12px; font-weight: 700; color: var(--text1); }

/* Target Box styling */
.target-box { background: var(--surface); border: 1px solid var(--border); border-radius: var(--r8); padding: 12px; margin-bottom: 18px; }
.tb-label { font-size: 10px; color: var(--text3); text-transform: uppercase; }
.tb-val { font-size: 20px; font-weight: 800; margin-top: 4px; line-height: 1; }
.tb-sub { font-size: 10px; color: var(--text3); margin-top: 4px; }
.tb-bar-wrap { height: 6px; background: var(--border); border-radius: 3px; overflow: hidden; margin-top: 8px; }
.tb-bar { height: 100%; border-radius: 3px; }
.tb-pct { font-size: 12px; font-weight: 800; }
.tb-target { font-size: 10px; color: var(--text3); text-align: right; }

/* Summary section details */
.summary-section { margin-bottom: 18px; }
.sum-title { font-size: 11px; font-weight: 700; color: var(--text2); text-transform: uppercase; letter-spacing: 0.5px; margin-bottom: 8px; }
.summary-details { background: var(--surface); border: 0.5px solid var(--border); border-radius: var(--r8); padding: 4px 10px; }
.sum-row { display: flex; justify-content: space-between; padding: 6px 0; border-bottom: 0.5px dashed var(--border); font-size: 11px; }
.sum-row:last-child { border-bottom: none; }
.sum-key { color: var(--text2); }
.sum-val { font-weight: 700; }
.sum-val.green { color: var(--hijau); }
.sum-val.red { color: var(--merah); }

/* Rekening Box */
.rekening-box { background: var(--surface); border: 0.5px solid var(--border); border-radius: var(--r8); padding: 10px 12px; }
.rek-title { display: flex; justify-content: space-between; align-items: center; font-size: 11px; font-weight: 700; color: var(--text2); text-transform: uppercase; letter-spacing: 0.5px; margin-bottom: 8px; }
.rek-edit { color: var(--hijau); font-size: 10px; cursor: pointer; text-transform: none; font-weight: normal; }
.rek-row { display: flex; justify-content: space-between; font-size: 11px; padding: 3px 0; }
.rek-key { color: var(--text3); }
.rek-val { font-weight: 600; }

/* Right panel tabs & content */
.panel-right { flex: 1; display: flex; flex-direction: column; height: 100%; min-width: 0; background: var(--surface); }
.tabs-header { display: flex; border-bottom: 1px solid var(--border); background: var(--card); padding: 0 16px; flex-shrink: 0; }
.tab { display: flex; align-items: center; gap: 6px; padding: 12px 16px; font-size: 12px; font-weight: 600; background: none; border: none; cursor: pointer; color: var(--text2); border-bottom: 2px solid transparent; }
.tab.active { color: var(--hijau); border-bottom-color: var(--hijau); }

.tab-container { flex: 1; display: flex; flex-direction: column; overflow: hidden; height: 100%; }
.toolbar { display: flex; justify-content: space-between; align-items: center; padding: 10px 16px; border-bottom: 1px solid var(--border); background: var(--card); flex-shrink: 0; }
.filter-select-sm { padding: 5px 8px; border: 1px solid var(--border); border-radius: var(--r6); font-size: 11px; background: var(--surface); color: var(--text2); cursor: pointer; }
.toolbar-info { font-size: 11px; color: var(--text3); }
.export-btn { border: 1px solid var(--border); padding: 5px 10px; border-radius: 6px; font-size: 11px; font-weight: 600; cursor: pointer; background: #fff; display: inline-flex; align-items: center; gap: 4px; color: var(--text2); }
.export-btn:hover { color: var(--hijau); border-color: var(--hijau); }

.list-scroll { flex: 1; overflow-y: auto; background: var(--card); }
.data-table { width: 100%; border-collapse: collapse; }
.data-table th { padding: 10px 14px; font-size: 10px; font-weight: 700; color: var(--text3); text-transform: uppercase; text-align: left; background: var(--surface); border-bottom: 1px solid var(--border); }
.data-table td { padding: 12px 14px; font-size: 12px; border-bottom: 1px solid var(--border); vertical-align: middle; }
.table-row-clickable { cursor: pointer; transition: background .15s; }
.table-row-clickable:hover td { background: var(--surface); }

.anggota-cell { display: flex; align-items: center; gap: 8px; }
.av { width: 28px; height: 28px; border-radius: 50%; display: flex; align-items: center; justify-content: center; font-size: 10px; font-weight: 700; color: #fff; flex-shrink: 0; }
.member-name { font-size: 12px; font-weight: 600; }
.member-number { font-size: 10px; color: var(--text3); font-family: monospace; }

.status-badge { display: inline-flex; align-items: center; gap: 4px; padding: 2px 8px; border-radius: 10px; font-size: 10px; font-weight: 700; }
.status-badge.lunas { background: var(--hijau3); color: var(--hijau); }
.status-badge.belum { background: #fde8e8; color: var(--merah); }

.metode-badge { background: var(--hijau3); color: var(--hijau); border-radius: 4px; padding: 2px 6px; font-size: 10px; font-weight: 700; }
.history-cell { font-weight: 700; color: var(--hijau); font-size: 11px; }

.action-btns { display: flex; gap: 4px; }
.action-btn { width: 26px; height: 26px; border: 1px solid var(--border); border-radius: var(--r6); background: #fff; cursor: pointer; display: inline-flex; align-items: center; justify-content: center; font-size: 12px; color: var(--text2); }
.action-btn:hover { border-color: var(--hijau); color: var(--hijau); background: var(--hijau3); }

.list-footer { display: flex; justify-content: space-between; align-items: center; padding: 10px 16px; border-top: 1px solid var(--border); background: var(--card); flex-shrink: 0; }
.lf-info { font-size: 11px; color: var(--text3); }

/* Rekap per unit style */
.unit-bar-cell { display: flex; align-items: center; gap: 8px; }
.ub-wrap { width: 60px; height: 5px; background: var(--border); border-radius: 3px; overflow: hidden; }
.ub-fill { height: 100%; border-radius: 3px; }

.branch-summary-bar { border-top: 1px solid var(--border); background: var(--card); padding: 12px 16px; display: flex; align-items: center; justify-content: space-between; flex-shrink: 0; }

/* Modals */
.modal-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.5); z-index: 1000; display: flex; align-items: center; justify-content: center; padding: 20px; }
.modal-card { background: var(--card); border-radius: var(--r12); width: 100%; max-width: 600px; max-height: 90vh; overflow-y: auto; box-shadow: var(--shadow-lg); }
.modal-card.mini { max-width: 440px; }
.modal-header { padding: 16px 20px; border-bottom: 1px solid var(--border); display: flex; justify-content: space-between; align-items: center; }
.modal-title { font-size: 14px; font-weight: 700; }
.modal-close { background: none; border: none; font-size: 16px; cursor: pointer; color: var(--text3); }
.modal-footer { padding: 12px 20px; border-top: 1px solid var(--border); display: flex; justify-content: flex-end; gap: 8px; }

.alert-info-box { background: #fff8e0; border: 0.5px solid var(--kuning); border-radius: var(--r8); padding: 9px 12px; font-size: 11px; color: #7a6000; display: flex; align-items: center; gap: 8px; margin-bottom: 14px; }
.alert-info-box i { font-size: 14px; }

.rekening-current-box { background: var(--surface); border-radius: var(--r8); padding: 10px 12px; margin-bottom: 14px; border: 0.5px solid var(--border); }
.rcb-title { font-size: 10px; color: var(--text3); text-transform: uppercase; font-weight: 700; margin-bottom: 4px; }
.rcb-desc { font-size: 12px; }

.form-group { margin-bottom: 12px; }
.form-label { font-size: 11px; font-weight: 600; color: var(--text2); margin-bottom: 4px; display: block; }
.form-input, .form-select { width: 100%; padding: 8px 10px; border: 1.5px solid var(--border2); border-radius: var(--r8); font-size: 12px; outline: none; background: #fff; }
.form-input:focus, .form-select:focus { border-color: var(--hijau); }

/* Riwayat Modal content style */
.riwayat-header { display: flex; align-items: center; gap: 12px; border-bottom: 1px solid var(--border); padding-bottom: 14px; }
.av-large { width: 48px; height: 48px; border-radius: 50%; display: flex; align-items: center; justify-content: center; font-size: 16px; font-weight: 700; color: #fff; }
.rh-name { font-weight: 800; font-size: 14px; }
.rh-sub { font-size: 11px; color: var(--text3); margin-top: 2px; }
.badge { display: inline-block; padding: 2px 8px; border-radius: 10px; font-size: 9px; font-weight: 700; }
.badge.border-green { background: var(--hijau3); color: var(--hijau); border: 0.5px solid var(--hijau4); }

.history-table { width: 100%; border-collapse: collapse; font-size: 11px; }
.history-table th { padding: 6px 10px; background: var(--surface); color: var(--text3); font-weight: 700; text-align: left; }
.history-table td { padding: 8px 10px; border-bottom: 1px solid var(--border); }
.history-table tr:last-child td { border-bottom: none; }

.spin { animation: spin .8s linear infinite; }
@keyframes spin { from { transform: rotate(0deg); } to { transform: rotate(360deg); } }
</style>
