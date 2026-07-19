<template>
  <div id="page-laporan" class="laporan-split-layout">
    <!-- ===== PANEL KIRI: MENU LAPORAN & SETTING ===== -->
    <div class="panel-left">
      <div class="pl-head">
        <div class="pl-title">Jenis Laporan</div>
        <div class="pl-sub">Pilih & atur kriteria laporan</div>
      </div>
      <div class="pl-scroll">
        <!-- Organisasi Group -->
        <div class="lap-group">
          <div class="lap-group-title">Organisasi</div>
          <button 
            :class="['lap-item', { active: activeReportType === 'keanggotaan' }]" 
            @click="activeReportType = 'keanggotaan'"
          >
            <div class="lap-icon"><i class="ti ti-users"></i></div>
            <div class="lap-details">
              <div class="lap-name">Keanggotaan</div>
              <div class="lap-desc">Jumlah, tingkatan, pertumbuhan</div>
            </div>
          </button>
          
          <button 
            :class="['lap-item', { active: activeReportType === 'kehadiran' }]" 
            @click="activeReportType = 'kehadiran'"
          >
            <div class="lap-icon"><i class="ti ti-qrcode"></i></div>
            <div class="lap-details">
              <div class="lap-name">Kehadiran</div>
              <div class="lap-desc">Rekap latihan per unit/cabang</div>
            </div>
          </button>
        </div>

        <!-- Keuangan Group -->
        <div class="lap-group">
          <div class="lap-group-title">Keuangan & Event</div>
          <button 
            :class="['lap-item', { active: activeReportType === 'iuran' }]" 
            @click="activeReportType = 'iuran'"
          >
            <div class="lap-icon"><i class="ti ti-wallet"></i></div>
            <div class="lap-details">
              <div class="lap-name">Keuangan BLBA</div>
              <div class="lap-desc">Terkumpul, tunggakan, target</div>
            </div>
          </button>
          
          <button 
            :class="['lap-item', { active: activeReportType === 'transport' }]" 
            @click="activeReportType = 'transport'"
          >
            <div class="lap-icon"><i class="ti ti-motorbike"></i></div>
            <div class="lap-details">
              <div class="lap-name">Transport Pelatih</div>
              <div class="lap-desc">Rekap biaya per pelatih</div>
            </div>
          </button>

          <button 
            :class="['lap-item', { active: activeReportType === 'ekt' }]" 
            @click="activeReportType = 'ekt'"
          >
            <div class="lap-icon"><i class="ti ti-award"></i></div>
            <div class="lap-details">
              <div class="lap-name">EKT & Latgab</div>
              <div class="lap-desc">Peserta & pendapatan event</div>
            </div>
          </button>
        </div>

        <!-- Period Section -->
        <div class="period-section">
          <span class="period-label">Periode</span>
          <div class="period-btns">
            <button 
              v-for="p in periods" 
              :key="p" 
              :class="['period-btn', { active: selectedPeriod === p }]" 
              @click="setPeriod(p)"
            >
              {{ p }}
            </button>
          </div>
          <div class="lp-form-row">
            <div>
              <label class="date-input-label">Dari</label>
              <input v-model="startDate" class="lp-form-input" type="date" />
            </div>
            <div>
              <label class="date-input-label">Sampai</label>
              <input v-model="endDate" class="lp-form-input" type="date" />
            </div>
          </div>
        </div>

        <!-- Lingkup Section -->
        <div class="scope-section">
          <span class="period-label">Lingkup Wilayah</span>
          <select v-model="selectedScope" class="lp-form-select">
            <option value="nasional">Seluruh Indonesia</option>
            <option value="cabang">Per Cabang</option>
          </select>
          
          <div v-if="selectedScope === 'cabang'" class="branch-select-group">
            <select v-model="selectedCabangName" class="lp-form-select" style="margin-top: 8px;">
              <option value="">Pilih Cabang</option>
              <option v-for="c in cabangList" :key="c.id" :value="c.nama">{{ c.nama }}</option>
            </select>
          </div>
        </div>

        <!-- Action Buttons -->
        <div class="action-buttons-group">
          <button class="export-btn primary" @click="generateReport" :disabled="generating">
            <i :class="['ti', generating ? 'ti-loader-2 spin' : 'ti-refresh']"></i>
            {{ generating ? 'Generating...' : 'Generate Laporan' }}
          </button>
          <button class="export-btn" @click="exportReport('PDF')"><i class="ti ti-file-type-pdf"></i> Export PDF</button>
          <button class="export-btn" @click="exportReport('Excel')"><i class="ti ti-file-spreadsheet"></i> Export Excel</button>
        </div>
      </div>
    </div>

    <!-- ===== PANEL KANAN: PREVIEW LAPORAN ===== -->
    <div class="panel-right">
      <div class="pr-toolbar">
        <div>
          <div class="pr-title">{{ getReportTitle(activeReportType) }}</div>
          <div class="pr-meta">{{ formattedPeriod }} · {{ formattedScope }}</div>
        </div>
        <div class="export-row">
          <button class="export-btn" @click="exportReport('PDF')"><i class="ti ti-file-type-pdf"></i> PDF</button>
          <button class="export-btn" @click="exportReport('Excel')"><i class="ti ti-file-spreadsheet"></i> Excel</button>
        </div>
      </div>

      <div class="pr-scroll" v-if="!generating">
        <!-- 1. KEANGGOTAAN -->
        <div v-if="activeReportType === 'keanggotaan'">
          <div class="sum-grid">
            <div class="sum-card green"><div class="sc-val">12.847</div><div class="sc-label">Total Anggota Aktif</div><div class="sc-sub">↑ 342 dari kuartal lalu</div></div>
            <div class="sum-card blue"><div class="sc-val">286</div><div class="sc-label">Anggota Baru</div><div class="sc-sub">Kuartal ini</div></div>
            <div class="sum-card yellow"><div class="sc-val">590</div><div class="sc-label">Anggota Nonaktif</div><div class="sc-sub">↓ 12 dari kuartal lalu</div></div>
            <div class="sum-card green"><div class="sc-val">148</div><div class="sc-label">Cabang Aktif</div><div class="sc-sub">Di 30 provinsi</div></div>
          </div>
          <div class="card">
            <div class="card-head"><div class="card-title"><i class="ti ti-map-pin"></i> Anggota per Cabang (Top 5)</div></div>
            <table class="data-table">
              <thead>
                <tr>
                  <th>#</th>
                  <th>Cabang</th>
                  <th>Provinsi</th>
                  <th>Aktif</th>
                  <th>Baru (kuartal)</th>
                  <th>Nonaktif</th>
                  <th>Growth</th>
                </tr>
              </thead>
              <tbody>
                <tr><td class="muted-td">1</td><td style="font-weight:600;">Kota Bandung</td><td>Jawa Barat</td><td class="bold-green">1.856</td><td>42</td><td class="text-red">67</td><td><span class="bdg bdg-g">↑ 2.3%</span></td></tr>
                <tr><td class="muted-td">2</td><td style="font-weight:600;">Kota Surabaya</td><td>Jawa Timur</td><td class="bold-green">1.524</td><td>38</td><td class="text-red">54</td><td><span class="bdg bdg-g">↑ 2.5%</span></td></tr>
                <tr><td class="muted-td">3</td><td style="font-weight:600;">Kota Semarang</td><td>Jawa Tengah</td><td class="bold-green">1.642</td><td>35</td><td class="text-red">48</td><td><span class="bdg bdg-g">↑ 2.1%</span></td></tr>
                <tr><td class="muted-td">4</td><td style="font-weight:600;">Kota Jakarta Central</td><td>DKI Jakarta</td><td class="bold-green">1.248</td><td>29</td><td class="text-red">42</td><td><span class="bdg bdg-g">↑ 2.3%</span></td></tr>
                <tr><td class="muted-td">5</td><td style="font-weight:600;">Kota Yogyakarta</td><td>DI Yogyakarta</td><td class="bold-green">892</td><td>18</td><td class="text-red">32</td><td><span class="bdg bdg-g">↑ 2.0%</span></td></tr>
              </tbody>
            </table>
          </div>
          <div class="card">
            <div class="card-head"><div class="card-title"><i class="ti ti-award"></i> Distribusi Tingkatan</div></div>
            <table class="data-table">
              <thead>
                <tr>
                  <th>Tingkatan</th>
                  <th>Jumlah</th>
                  <th>Proporsi</th>
                  <th>Visual</th>
                </tr>
              </thead>
              <tbody>
                <tr><td style="font-weight:600;">Pra Dasar</td><td style="font-weight:700;">4.820</td><td>37.5%</td><td><div class="pct-bar"><div class="pb-wrap"><div class="pb-fill pb-g" style="width:38%;"></div></div></div></td></tr>
                <tr><td style="font-weight:600;">Dasar</td><td style="font-weight:700;">3.612</td><td>28.1%</td><td><div class="pct-bar"><div class="pb-wrap"><div class="pb-fill pb-g" style="width:28%;"></div></div></div></td></tr>
                <tr><td style="font-weight:600;">PH</td><td style="font-weight:700;">2.104</td><td>16.4%</td><td><div class="pct-bar"><div class="pb-wrap"><div class="pb-fill pb-g" style="width:16%;"></div></div></div></td></tr>
                <tr><td style="font-weight:600;">Gabungan</td><td style="font-weight:700;">1.248</td><td>9.7%</td><td><div class="pct-bar"><div class="pb-wrap"><div class="pb-fill pb-y" style="width:10%;"></div></div></div></td></tr>
                <tr><td style="font-weight:600;">PK</td><td style="font-weight:700;">642</td><td>5.0%</td><td><div class="pct-bar"><div class="pb-wrap"><div class="pb-fill pb-y" style="width:5%;"></div></div></div></td></tr>
                <tr><td style="font-weight:600;">GPK</td><td style="font-weight:700;">284</td><td>2.2%</td><td><div class="pct-bar"><div class="pb-wrap"><div class="pb-fill pb-r" style="width:2%;"></div></div></div></td></tr>
              </tbody>
              <tfoot>
                <tr><td>Total Nasional</td><td>12.847</td><td>100%</td><td></td></tr>
              </tfoot>
            </table>
          </div>
        </div>

        <!-- 2. KEHADIRAN -->
        <div v-else-if="activeReportType === 'kehadiran'">
          <div class="sum-grid">
            <div class="sum-card green"><div class="sc-val">4.284</div><div class="sc-label">Total Sesi Latihan</div><div class="sc-sub">Kuartal ini</div></div>
            <div class="sum-card blue"><div class="sc-val">78%</div><div class="sc-label">Rata-rata Kehadiran</div><div class="sc-sub">Nasional</div></div>
            <div class="sum-card yellow"><div class="sc-val">1.842</div><div class="sc-label">Anggota Kurang Aktif</div><div class="sc-sub">&lt;50% Kehadiran</div></div>
            <div class="sum-card green"><div class="sc-val">9.214</div><div class="sc-label">Memenuhi Syarat EKT</div><div class="sc-sub">≥10 Sesi</div></div>
          </div>
          <div class="card">
            <div class="card-head"><div class="card-title"><i class="ti ti-chart-bar"></i> Persentase Kehadiran per Cabang</div></div>
            <table class="data-table">
              <thead>
                <tr>
                  <th>Cabang</th>
                  <th>Total Sesi</th>
                  <th>Kehadiran</th>
                  <th>Syarat EKT ✓</th>
                </tr>
              </thead>
              <tbody>
                <tr><td style="font-weight:600;">Kota Bandung</td><td>312</td><td><div class="progress-bar-cell"><div class="pb-wrap"><div class="pb-fill pb-g" style="width:82%;"></div></div><span class="bold-green">82%</span></div></td><td style="font-weight:700;color:var(--hijau);">1.524</td></tr>
                <tr><td style="font-weight:600;">Kota Surabaya</td><td>276</td><td><div class="progress-bar-cell"><div class="pb-wrap"><div class="pb-fill pb-g" style="width:79%;"></div></div><span class="bold-green">79%</span></div></td><td style="font-weight:700;color:var(--hijau);">1.187</td></tr>
                <tr><td style="font-weight:600;">Kota Semarang</td><td>294</td><td><div class="progress-bar-cell"><div class="pb-wrap"><div class="pb-fill pb-g" style="width:75%;"></div></div><span class="bold-green">75%</span></div></td><td style="font-weight:700;color:var(--hijau);">1.201</td></tr>
                <tr><td style="font-weight:600;">Kota Yogyakarta</td><td>198</td><td><div class="progress-bar-cell"><div class="pb-wrap"><div class="pb-fill pb-g" style="width:78%;"></div></div><span class="bold-green">78%</span></div></td><td style="font-weight:700;color:var(--hijau);">687</td></tr>
                <tr><td style="font-weight:600;">Kota Makassar</td><td>162</td><td><div class="progress-bar-cell"><div class="pb-wrap"><div class="pb-fill pb-y" style="width:61%;"></div></div><span class="bold-yellow">61%</span></div></td><td style="font-weight:700;color:#9a7000;">234</td></tr>
              </tbody>
            </table>
          </div>
        </div>

        <!-- 3. KEUANGAN BLBA -->
        <div v-else-if="activeReportType === 'iuran'">
          <div class="sum-grid">
            <div class="sum-card green"><div class="sc-val">Rp 1,54 M</div><div class="sc-label">Total Terkumpul</div><div class="sc-sub">Kuartal ini</div></div>
            <div class="sum-card blue"><div class="sc-val">84%</div><div class="sc-label">Pencapaian Target</div><div class="sc-sub">Nasional</div></div>
            <div class="sum-card red"><div class="sc-val">Rp 295 Jt</div><div class="sc-label">Tunggakan</div><div class="sc-sub">Belum Terbayar</div></div>
            <div class="sum-card yellow"><div class="sc-val">1.842</div><div class="sc-label">Anggota Menunggak</div><div class="sc-sub">&gt;1 Bulan</div></div>
          </div>
          <div class="card">
            <div class="card-head"><div class="card-title"><i class="ti ti-wallet"></i> Rekap Biaya Bulanan per Cabang</div></div>
            <table class="data-table">
              <thead>
                <tr>
                  <th>Cabang</th>
                  <th>Anggota Aktif</th>
                  <th>Target Dana</th>
                  <th>Terkumpul</th>
                  <th>Pencapaian</th>
                </tr>
              </thead>
              <tbody>
                <tr><td style="font-weight:600;">Kota Bandung</td><td>1.856</td><td>Rp 222,7 Jt</td><td class="bold-green">Rp 196,4 Jt</td><td><div class="progress-bar-cell"><div class="pb-wrap"><div class="pb-fill pb-g" style="width:88%;"></div></div><span class="bold-green">88%</span></div></td></tr>
                <tr><td style="font-weight:600;">Kota Surabaya</td><td>1.524</td><td>Rp 182,9 Jt</td><td class="bold-green">Rp 152,4 Jt</td><td><div class="progress-bar-cell"><div class="pb-wrap"><div class="pb-fill pb-g" style="width:83%;"></div></div><span class="bold-green">83%</span></div></td></tr>
                <tr><td style="font-weight:600;">Kota Semarang</td><td>1.642</td><td>Rp 197,0 Jt</td><td class="bold-green">Rp 157,6 Jt</td><td><div class="progress-bar-cell"><div class="pb-wrap"><div class="pb-fill pb-g" style="width:80%;"></div></div><span class="bold-green">80%</span></div></td></tr>
                <tr><td style="font-weight:600;">Kota Makassar</td><td>387</td><td>Rp 46,4 Jt</td><td class="bold-yellow">Rp 27,8 Jt</td><td><div class="progress-bar-cell"><div class="pb-wrap"><div class="pb-fill pb-y" style="width:60%;"></div></div><span class="bold-yellow">60%</span></div></td></tr>
              </tbody>
              <tfoot>
                <tr><td colspan="2">Total Nasional</td><td>Rp 1,84 M</td><td style="color:var(--hijau);">Rp 1,54 M</td><td style="color:var(--hijau);">84%</td></tr>
              </tfoot>
            </table>
          </div>
        </div>

        <!-- 4. TRANSPORT PELATIH -->
        <div v-else-if="activeReportType === 'transport'">
          <div class="sum-grid">
            <div class="sum-card green"><div class="sc-val">Rp 68,4 Jt</div><div class="sc-label">Total Transport</div><div class="sc-sub">Kuartal ini</div></div>
            <div class="sum-card blue"><div class="sc-val">2.280</div><div class="sc-label">Total Sesi Latihan</div><div class="sc-sub">Seluruh Cabang</div></div>
            <div class="sum-card green"><div class="sc-val">342</div><div class="sc-label">Pelatih Aktif</div><div class="sc-sub">Daerah & Pusat</div></div>
            <div class="sum-card yellow"><div class="sc-val">Rp 30.000</div><div class="sc-label">Rata-rata per Sesi</div><div class="sc-sub">Default nominal</div></div>
          </div>
          <div class="card">
            <div class="card-head"><div class="card-title"><i class="ti ti-motorbike"></i> Biaya Transport per Cabang</div></div>
            <table class="data-table">
              <thead>
                <tr>
                  <th>Cabang</th>
                  <th>Jumlah Pelatih</th>
                  <th>Total Sesi</th>
                  <th>Total Biaya</th>
                </tr>
              </thead>
              <tbody>
                <tr><td style="font-weight:600;">Kota Bandung</td><td>48</td><td>384</td><td class="bold-green">Rp 11.880.000</td></tr>
                <tr><td style="font-weight:600;">Kota Surabaya</td><td>38</td><td>312</td><td class="bold-green">Rp 9.720.000</td></tr>
                <tr><td style="font-weight:600;">Kota Semarang</td><td>42</td><td>336</td><td class="bold-green">Rp 10.320.000</td></tr>
                <tr><td style="font-weight:600;">Kota Yogyakarta</td><td>18</td><td>198</td><td class="bold-green">Rp 5.940.000</td></tr>
              </tbody>
              <tfoot>
                <tr><td colspan="2">Total Nasional</td><td>2.280 sesi</td><td style="color:var(--hijau)">Rp 68.400.000</td></tr>
              </tfoot>
            </table>
          </div>
        </div>

        <!-- 5. EKT & LATGAB -->
        <div v-else-if="activeReportType === 'ekt'">
          <div class="sum-grid">
            <div class="sum-card green"><div class="sc-val">24</div><div class="sc-label">Total Event</div><div class="sc-sub">14 EKT · 10 Latgab</div></div>
            <div class="sum-card blue"><div class="sc-val">2.847</div><div class="sc-label">Total Peserta</div><div class="sc-sub">Semua event</div></div>
            <div class="sum-card green"><div class="sc-val">Rp 94,5 Jt</div><div class="sc-label">Pendapatan EKT</div><div class="sc-sub">Event Berbayar</div></div>
            <div class="sum-card yellow"><div class="sc-val">Rp 12,3 Jt</div><div class="sc-label">Pendapatan Latgab</div><div class="sc-sub">Event Berbayar</div></div>
          </div>
          <div class="card">
            <div class="card-head"><div class="card-title"><i class="ti ti-award"></i> Daftar Event Terselenggara</div></div>
            <table class="data-table">
              <thead>
                <tr>
                  <th>Event</th>
                  <th>Tanggal</th>
                  <th>Cabang</th>
                  <th>Peserta</th>
                  <th>Pendapatan</th>
                </tr>
              </thead>
              <tbody>
                <tr><td style="font-weight:600;">EKT Cabang Bandung<br><span class="bdg bdg-r" style="font-size:9px;">EKT</span></td><td style="font-size:11px;">15 Jun 2026</td><td style="font-size:11px;">Kota Bandung</td><td style="font-weight:700;">42</td><td class="bold-green">Rp 3.150.000</td></tr>
                <tr><td style="font-weight:600;">Latgab Se-Jawa Barat<br><span class="bdg bdg-b" style="font-size:9px;">Latgab</span></td><td style="font-size:11px;">22 Jun 2026</td><td style="font-size:11px;">Kota Bandung</td><td style="font-weight:700;">187</td><td class="muted-text">Gratis / Kas</td></tr>
                <tr><td style="font-weight:600;">EKT Cabang Surabaya<br><span class="bdg bdg-r" style="font-size:9px;">EKT</span></td><td style="font-size:11px;">29 Jun 2026</td><td style="font-size:11px;">Kota Surabaya</td><td style="font-weight:700;">38</td><td class="bold-green">Rp 2.850.000</td></tr>
                <tr><td style="font-weight:600;">EKT Cabang Yogyakarta<br><span class="bdg bdg-r" style="font-size:9px;">EKT</span></td><td style="font-size:11px;">1 Jun 2026</td><td style="font-size:11px;">Kota Yogyakarta</td><td style="font-weight:700;">45</td><td class="bold-green">Rp 3.375.000</td></tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>

      <!-- Loading State -->
      <div v-else class="generating-state">
        <i class="ti ti-loader-2 spin"></i>
        <div style="font-size:14px; font-weight:700; margin-top:12px;">Menghitung Data Laporan...</div>
        <div style="font-size:11px; color:var(--text3); margin-top:4px;">Mengambil ringkasan data dari database organisasi.</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({ title: 'Laporan' })

const api = useApi()

const activeReportType = ref('keanggotaan')
const selectedPeriod = ref('Kuartal')
const startDate = ref('2026-04-01')
const endDate = ref('2026-06-30')
const selectedScope = ref('nasional')
const selectedCabangName = ref('')
const generating = ref(false)

const periods = ['Bulan ini', 'Kuartal', 'Semester', 'Tahunan']
const cabangList = ref<any[]>([])

const fetchCabang = async () => {
  try {
    const data = await api.get('/organization/cabang?limit=100')
    cabangList.value = data.data || []
    if (cabangList.value.length > 0) {
      selectedCabangName.value = cabangList.value[0].nama
    }
  } catch (e) {
    console.error(e)
  }
}

const getReportTitle = (type: string) => {
  if (type === 'keanggotaan') return 'Laporan Keanggotaan'
  if (type === 'kehadiran') return 'Laporan Kehadiran Latihan'
  if (type === 'iuran') return 'Laporan Keuangan BLBA'
  if (type === 'transport') return 'Laporan Transport Pelatih'
  return 'Laporan EKT & Latgab'
}

const formattedPeriod = computed(() => {
  const dStart = new Date(startDate.value)
  const dEnd = new Date(endDate.value)
  
  const mStart = dStart.toLocaleDateString('id-ID', { month: 'short', year: 'numeric' })
  const mEnd = dEnd.toLocaleDateString('id-ID', { month: 'short', year: 'numeric' })
  
  return `${mStart} – ${mEnd}`
})

const formattedScope = computed(() => {
  if (selectedScope.value === 'nasional') return 'Seluruh Indonesia'
  return selectedCabangName.value ? `Cabang ${selectedCabangName.value}` : 'Semua Cabang'
})

const setPeriod = (period: string) => {
  selectedPeriod.value = period
  const today = new Date(2026, 5, 30) // Base on June 2026

  if (period === 'Bulan ini') {
    startDate.value = '2026-06-01'
    endDate.value = '2026-06-30'
  } else if (period === 'Kuartal') {
    startDate.value = '2026-04-01'
    endDate.value = '2026-06-30'
  } else if (period === 'Semester') {
    startDate.value = '2026-01-01'
    endDate.value = '2026-06-30'
  } else {
    startDate.value = '2026-01-01'
    endDate.value = '2026-12-31'
  }
}

const generateReport = () => {
  generating.value = true
  setTimeout(() => {
    generating.value = false
    alert('Laporan berhasil diperbarui berdasarkan kriteria penyaringan baru!')
  }, 600)
}

const exportReport = (format: string) => {
  alert(`Laporan ${getReportTitle(activeReportType.value)} berhasil diekspor ke format ${format}!`)
}

onMounted(() => {
  fetchCabang()
})
</script>

<style scoped>
.laporan-split-layout {
  display: flex;
  height: calc(100vh - 64px - 44px);
  background: var(--bg);
  overflow: hidden;
  margin: -16px;
}

/* Left panel menu */
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
.pl-title { font-size: 15px; font-weight: 800; }
.pl-sub { font-size: 11px; color: var(--text3); margin-top: 2px; }
.pl-scroll { flex: 1; overflow-y: auto; padding: 20px; }

.lap-group { margin-bottom: 18px; }
.lap-group-title { font-size: 10px; font-weight: 800; color: var(--text3); text-transform: uppercase; margin-bottom: 8px; letter-spacing: .06em; }
.lap-item { display: flex; align-items: center; gap: 10px; padding: 8px 12px; border-radius: var(--r8); width: 100%; background: none; border: none; cursor: pointer; text-align: left; transition: all .15s; margin-bottom: 4px; }
.lap-item:hover { background: var(--surface); }
.lap-item.active { background: var(--hijau3); }
.lap-icon { width: 30px; height: 30px; border-radius: var(--r6); background: var(--surface); color: var(--text2); display: flex; align-items: center; justify-content: center; font-size: 14px; flex-shrink: 0; }
.lap-item.active .lap-icon { background: var(--card); color: var(--hijau); }
.lap-details { min-width: 0; flex: 1; }
.lap-name { font-size: 12px; font-weight: 700; color: var(--text1); }
.lap-desc { font-size: 10px; color: var(--text3); white-space: nowrap; overflow: hidden; text-overflow: ellipsis; margin-top: 1px; }

/* Period filters style */
.period-section, .scope-section { margin-bottom: 18px; padding-top: 14px; border-top: 1px dashed var(--border); }
.period-label { font-size: 10px; font-weight: 800; color: var(--text3); text-transform: uppercase; display: block; margin-bottom: 8px; letter-spacing: .06em; }

.period-btns { display: grid; grid-template-columns: repeat(4, 1fr); gap: 4px; margin-bottom: 12px; }
.period-btn { padding: 6px 0; border: 1px solid var(--border2); border-radius: var(--r6); background: #fff; font-size: 10px; font-weight: 600; cursor: pointer; color: var(--text2); text-align: center; }
.period-btn.active { background: var(--hijau); color: #fff; border-color: var(--hijau); }

.lp-form-row { display: grid; grid-template-columns: 1fr 1fr; gap: 8px; }
.date-input-label { font-size: 9px; color: var(--text3); font-weight: 600; margin-bottom: 2px; display: block; }
.lp-form-input, .lp-form-select { width: 100%; padding: 6px 10px; border: 1.5px solid var(--border2); border-radius: var(--r8); font-size: 11px; outline: none; background: #fff; color: var(--text1); }

.action-buttons-group { display: flex; flex-direction: column; gap: 6px; margin-top: 18px; padding-top: 14px; border-top: 1px dashed var(--border); }
.export-btn { border: 1px solid var(--border); padding: 7px 12px; border-radius: var(--r8); font-size: 11px; font-weight: 600; cursor: pointer; background: #fff; display: inline-flex; align-items: center; justify-content: center; gap: 6px; color: var(--text2); }
.export-btn:hover { color: var(--hijau); border-color: var(--hijau); }
.export-btn.primary { background: var(--hijau); color: #fff; border-color: var(--hijau); }
.export-btn.primary:hover { background: var(--hijau2); }

/* Right panel preview style */
.panel-right { flex: 1; display: flex; flex-direction: column; height: 100%; min-width: 0; background: var(--surface); }
.pr-toolbar { display: flex; justify-content: space-between; align-items: center; padding: 14px 20px; border-bottom: 1px solid var(--border); background: var(--card); flex-shrink: 0; }
.pr-title { font-size: 14px; font-weight: 800; }
.pr-meta { font-size: 11px; color: var(--text3); margin-top: 2px; }
.export-row { display: flex; gap: 6px; }

.pr-scroll { flex: 1; overflow-y: auto; padding: 20px; display: flex; flex-direction: column; gap: 16px; }

.sum-grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 12px; }
.sum-card { background: var(--card); border: 1px solid var(--border); border-radius: var(--r12); padding: 14px; position: relative; overflow: hidden; }
.sum-card::before { content: ''; position: absolute; top: 0; left: 0; right: 0; height: 3px; }
.sum-card.green::before { background: var(--hijau2); }
.sum-card.blue::before { background: var(--biru); }
.sum-card.yellow::before { background: var(--kuning); }
.sum-card.red::before { background: var(--merah); }
.sc-val { font-size: 20px; font-weight: 800; line-height: 1; }
.sc-label { font-size: 10px; color: var(--text3); margin-top: 6px; font-weight: 700; text-transform: uppercase; }
.sc-sub { font-size: 9px; color: var(--text3); margin-top: 2px; }

.card { background: var(--card); border: 1px solid var(--border); border-radius: var(--r12); overflow: hidden; }
.card-head { display: flex; justify-content: space-between; align-items: center; padding: 12px 16px; border-bottom: 1px solid var(--border); }
.card-title { font-size: 12px; font-weight: 700; display: flex; align-items: center; gap: 6px; }
.card-title i { color: var(--hijau); font-size: 14px; }

/* Table styles */
.data-table { width: 100%; border-collapse: collapse; }
.data-table th { padding: 10px 16px; font-size: 10px; font-weight: 700; color: var(--text3); text-transform: uppercase; background: var(--surface); border-bottom: 1px solid var(--border); text-align: left; }
.data-table td { padding: 10px 16px; font-size: 12px; border-bottom: 1px solid var(--border); vertical-align: middle; }
.data-table tfoot td { font-weight: 700; background: var(--surface); padding: 10px 16px; border-top: 1.5px solid var(--border); }

.muted-td { color: var(--text3); }
.bold-green { font-weight: 700; color: var(--hijau); }
.bold-yellow { font-weight: 700; color: #a07000; }
.text-red { color: var(--merah); }
.muted-text { color: var(--text3); font-size: 11px; }

.bdg { display: inline-flex; align-items: center; gap: 4px; padding: 2px 8px; border-radius: 10px; font-size: 10px; font-weight: 700; }
.bdg-g { background: var(--hijau3); color: var(--hijau); }
.bdg-r { background: #fde8e8; color: var(--merah); }
.bdg-b { background: #e0f5fb; color: var(--biru); }
.bdg-y { background: #fff8e0; color: #a07000; }

/* Visual Pct Bars */
.pct-bar, .progress-bar-cell { display: flex; align-items: center; gap: 8px; }
.pb-wrap { width: 60px; height: 5px; background: var(--border); border-radius: 3px; overflow: hidden; }
.pb-fill { height: 100%; border-radius: 3px; }
.pb-fill.pb-g { background: var(--hijau2); }
.pb-fill.pb-y { background: var(--kuning); }
.pb-fill.pb-r { background: var(--merah); }

/* Generating spinner overlay */
.generating-state { flex: 1; display: flex; flex-direction: column; align-items: center; justify-content: center; color: var(--text2); }
.spin { animation: spin .8s linear infinite; }
@keyframes spin { from { transform: rotate(0deg); } to { transform: rotate(360deg); } }
</style>
