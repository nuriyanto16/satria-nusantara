<template>
  <div id="page-dashboard">
    <!-- Period Selector -->
    <div class="period-bar">
      <button
        v-for="p in periods" :key="p.value"
        :class="['period-btn', { active: activePeriod === p.value }]"
        @click="activePeriod = p.value"
      >{{ p.label }}</button>
      <span class="period-label">Data diperbarui hari ini pukul 07:00 WIB</span>
    </div>

    <!-- Stat Cards Grid -->
    <div class="stat-grid">
      <div class="stat-card green" @click="navigateTo('/anggota')">
        <div class="stat-icon" style="background:var(--hijau3);color:var(--hijau)"><i class="ti ti-users"></i></div>
        <div v-if="loadingStats" class="skeleton-block" style="width:80px;height:28px;margin-bottom:6px;"></div>
        <div v-else class="stat-val">{{ totalAnggota }}</div>
        <div class="stat-label">Total Anggota</div>
        <div v-if="loadingStats" class="skeleton-block" style="width:100px;height:12px;"></div>
        <div v-else class="stat-trend trend-up"><i class="ti ti-trending-up"></i> +{{ newAnggotaThisWeek }} minggu ini</div>
      </div>

      <div class="stat-card yellow" @click="navigateTo('/cabang')">
        <div class="stat-icon" style="background:#fff8e0;color:var(--kuning)"><i class="ti ti-map-pin"></i></div>
        <div v-if="loadingStats" class="skeleton-block" style="width:60px;height:28px;margin-bottom:6px;"></div>
        <div v-else class="stat-val">{{ totalCabang }}</div>
        <div class="stat-label">Cabang Aktif</div>
        <div v-if="loadingStats" class="skeleton-block" style="width:80px;height:12px;"></div>
        <div v-else class="stat-sub">Di {{ totalProvinsi }} Provinsi</div>
      </div>

      <div class="stat-card green" @click="navigateTo('/cabang')">
        <div class="stat-icon" style="background:var(--hijau3);color:var(--hijau)"><i class="ti ti-home"></i></div>
        <div v-if="loadingStats" class="skeleton-block" style="width:60px;height:28px;margin-bottom:6px;"></div>
        <div v-else class="stat-val">{{ totalUnit }}</div>
        <div class="stat-label">Unit Latihan</div>
        <div v-if="loadingStats" class="skeleton-block" style="width:100px;height:12px;"></div>
        <div v-else class="stat-trend trend-up"><i class="ti ti-trending-up"></i> +{{ newCabangThisMonth }} unit baru</div>
      </div>

      <div class="stat-card blue" @click="navigateTo('/iuran')">
        <div class="stat-icon" style="background:#e0f5fb;color:var(--biru)"><i class="ti ti-wallet"></i></div>
        <div v-if="loadingStats" class="skeleton-block" style="width:100px;height:28px;margin-bottom:6px;"></div>
        <div v-else class="stat-val">{{ formatRupiahJuta(financialTerkumpul) }}</div>
        <div class="stat-label">Keuangan BLBA</div>
        <div v-if="loadingStats" class="skeleton-block" style="width:120px;height:12px;"></div>
        <div v-else class="stat-sub">{{ financialPct }}% dari target {{ formatRupiahJuta(financialTarget) }}</div>
      </div>

      <div class="stat-card blue" @click="navigateTo('/absensi')">
        <div class="stat-icon" style="background:#e0f5fb;color:var(--biru)"><i class="ti ti-calendar-event"></i></div>
        <div v-if="loadingStats" class="skeleton-block" style="width:70px;height:28px;margin-bottom:6px;"></div>
        <div v-else class="stat-val">{{ attendanceRate }}%</div>
        <div class="stat-label">Kehadiran Sesi</div>
        <div v-if="loadingStats" class="skeleton-block" style="width:110px;height:12px;"></div>
        <div v-else class="stat-trend trend-up"><i class="ti ti-trending-up"></i> ▲ 4% vs bulan lalu</div>
      </div>

      <div class="stat-card teal" @click="navigateTo('/kebugaran')">
        <div class="stat-icon" style="background:#e0f7f5;color:#2abcaa"><i class="ti ti-heartbeat"></i></div>
        <div v-if="loadingStats" class="skeleton-block" style="width:60px;height:28px;margin-bottom:6px;"></div>
        <div v-else class="stat-val">{{ kebugaranBaikPct }}%</div>
        <div class="stat-label">Nilai Kebugaran Baik</div>
        <div v-if="loadingStats" class="skeleton-block" style="width:100px;height:12px;"></div>
        <div v-else class="stat-sub">{{ kebugaranTotalDites }} anggota dites</div>
      </div>

      <div class="stat-card green" @click="navigateTo('/nafas')">
        <div class="stat-icon" style="background:#e8f5e9;color:var(--hijau)"><i class="ti ti-wind"></i></div>
        <div class="stat-val">Modul Nafas</div>
        <div class="stat-label">Monitoring Olah Nafas</div>
        <div class="stat-sub">Pantau Riwayat Latihan →</div>
      </div>
    </div>

    <!-- MAIN GRID (MAPS LIBRARY) -->
    <div class="main-grid">
      <div class="card">
        <div class="card-header">
          <div class="card-title">
            <svg viewBox="0 0 24 24"><path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"/></svg>
            Peta Sebaran Anggota se-Indonesia (Leaflet Maps Library)
          </div>
          <div class="card-action">Fullscreen</div>
        </div>
        <div class="map-wrap" style="height: 480px; position: relative; background: #e8f4f8;">
          <div v-if="loadingMap" class="skeleton-panel" style="height:480px;">
            <div class="skeleton-pulse-icon"><i class="ti ti-map-2" style="font-size:32px;color:#c8dfc5;"></i></div>
            <div class="skeleton-pulse-text">Memuat peta sebaran...</div>
          </div>
          <ClientOnly v-else>
            <div id="leaflet-map" style="width: 100%; height: 100%; z-index: 1;"></div>
          </ClientOnly>
        </div>
      </div>
    </div>

    <!-- FITNESS RANKING & TRENDS -->
    <div class="fitness-dashboard-grid" style="margin-top: 16px; display: grid; grid-template-columns: 1fr 2fr; gap: 16px; margin-bottom: 16px;">
      <!-- Card kiri: Ranking list (per-panel skeleton) -->
      <div class="card">
        <div class="card-header">
          <div class="card-title">
            <svg viewBox="0 0 24 24" style="width:16px;height:16px;fill:#5b21b6;display:inline-block;vertical-align:middle;margin-right:6px;"><path d="M12 21.593c-5.63-5.539-11-10.297-11-14.402 0-3.791 3.068-5.191 5.281-5.191 1.312 0 4.151.501 5.719 4.457 1.59-3.968 4.464-4.447 5.726-4.447 2.54 0 5.274 1.621 5.274 5.181 0 4.069-5.136 8.625-11 14.402z"/></svg>
            Ranking Kebugaran Cabang
          </div>
          <div style="display:flex;align-items:center;gap:6px;">
            <select style="font-size:11px;padding:3px 6px;border:1px solid var(--border);border-radius:6px;background:var(--bg);color:var(--text2);outline:none;">
              <option>Periode 2 — 2026</option>
              <option>Periode 1 — 2026</option>
            </select>
          </div>
        </div>
        <div v-if="loadingRanking" class="card-body">
          <div class="panel-skeleton-header"></div>
          <div class="kehadiran-row" v-for="n in 5" :key="n" style="gap:10px;padding:10px 0;">
            <div class="skeleton-block" style="width:26px;height:26px;border-radius:50%;"></div>
            <div style="flex:1;display:flex;flex-direction:column;gap:5px;">
              <div class="skeleton-block" style="width:110px;height:12px;"></div>
              <div class="skeleton-block" style="width:70px;height:10px;"></div>
            </div>
            <div class="skeleton-block" style="width:80px;height:6px;border-radius:3px;"></div>
            <div class="skeleton-block" style="width:36px;height:14px;"></div>
          </div>
        </div>
        <div v-else class="card-body" style="padding:0;">
          <div style="padding:8px 14px 4px;display:flex;gap:6px;font-size:10px;font-weight:700;color:var(--text3);text-transform:uppercase;letter-spacing:.04em;border-bottom:.5px solid var(--border);">
            <div style="width:28px;">RNK</div><div style="flex:1;">Cabang</div><div style="width:100px;">Skor</div><div style="width:54px;text-align:center;">Kat.</div>
          </div>
          <div v-for="(r, index) in rankingKebugaran" :key="r.cabang" class="kehadiran-row" style="padding:10px 14px;display:flex;align-items:center;gap:6px;border-bottom:.5px solid var(--border);">
            <div style="width:28px;">
              <div :class="['rank-num', { top: index < 3 }]" style="width:20px;height:20px;border-radius:50%;display:flex;align-items:center;justify-content:center;font-size:10px;font-weight:800;">
                {{ index + 1 }}
              </div>
            </div>
            <div style="flex:1;">
              <div style="font-size:12px;font-weight:600;">{{ r.cabang }}</div>
              <div style="font-size:10px;color:var(--text3);">DI Yogyakarta</div>
            </div>
            <div style="width:100px;display:flex;align-items:center;gap:4px;">
              <div style="flex:1;height:5px;background:var(--border);border-radius:3px;">
                <div :style="{ width: r.skor + '%', background: getKebugaranColor(r.skor) }" style="height:100%;border-radius:3px;"></div>
              </div>
              <span style="font-size:11px;font-weight:700;" :style="{ color: getKebugaranColor(r.skor) }">{{ r.skor }}%</span>
            </div>
            <div style="width:54px;text-align:center;">
              <span :style="getBadgeStyle(r.kategori)" style="font-size:9px;padding:2px 6px;border-radius:6px;font-weight:700;">
                {{ r.kategori }}
              </span>
            </div>
          </div>
        </div>
      </div>

      <!-- Card kanan: Line chart tren kebugaran -->
      <div class="card">
        <div class="card-header">
          <div class="card-title">
            <svg viewBox="0 0 24 24" style="width:15px;height:15px;fill:none;stroke:#5b21b6;stroke-width:2;stroke-linecap:round;stroke-linejoin:round;display:inline-block;vertical-align:middle;margin-right:6px;"><polyline points="3,17 8,11 13,14 18,6 21,9"/></svg>
            Tren Kebugaran per Periode
          </div>
          <div style="font-size:10px;color:var(--text3);">Data tingkat nasional per periode tes</div>
        </div>
        <div v-if="loadingTren" class="card-body" style="padding:16px;">
          <div class="panel-skeleton-header" style="width:200px;"></div>
          <div style="display:flex;align-items:flex-end;gap:6px;height:200px;padding-top:12px;">
            <div v-for="n in 8" :key="n" class="skeleton-block" :style="{flex:1, height: (40 + n * 12) + 'px', borderRadius:'4px 4px 0 0'}"></div>
          </div>
        </div>
        <div v-else class="card-body" style="padding:14px 16px 10px;">
          <ClientOnly>
            <apexchart 
              type="area" 
              height="250" 
              :options="optionsKebugaran" 
              :series="seriesKebugaran"
            />
          </ClientOnly>
        </div>
      </div>
    </div>

    <!-- BOTTOM GRID -->
    <div class="bottom-grid">
      <!-- STATISTIK ANGGOTA -->
      <div class="card">
        <div class="card-header">
          <div class="card-title">
            <i class="ti ti-chart-pie" style="font-size:16px;color:var(--hijau)"></i>
            Statistik Anggota
          </div>
        </div>
        <div v-if="loadingStatistik" class="card-body">
          <div class="panel-skeleton-header" style="width:100px;margin-bottom:12px;"></div>
          <div style="display:flex;align-items:center;gap:16px;margin-bottom:20px;">
            <div class="skeleton-block" style="width:100px;height:100px;border-radius:50%;"></div>
            <div style="flex:1;display:flex;flex-direction:column;gap:8px;">
              <div class="skeleton-block" style="width:90px;height:12px;"></div>
              <div class="skeleton-block" style="width:70px;height:12px;"></div>
            </div>
          </div>
          <div class="panel-skeleton-header" style="width:110px;margin-bottom:12px;"></div>
          <div v-for="n in 4" :key="n" style="display:flex;align-items:center;gap:8px;margin-bottom:10px;">
            <div class="skeleton-block" style="width:65px;height:12px;"></div>
            <div class="skeleton-block" style="flex:1;height:10px;border-radius:5px;"></div>
            <div class="skeleton-block" style="width:32px;height:12px;"></div>
          </div>
        </div>
        <div v-else class="card-body">
          <div style="font-size:11px;font-weight:700;color:var(--text3);margin-bottom:8px;">Jenis Kelamin</div>
          <div style="margin-bottom:14px;">
            <ClientOnly>
              <apexchart 
                type="donut" 
                height="150" 
                :options="optionsGender" 
                :series="seriesGender"
              />
            </ClientOnly>
          </div>

          <div style="font-size:11px;font-weight:700;color:var(--text3);margin-bottom:8px;">Kelompok Usia</div>
          <ClientOnly>
            <apexchart 
              type="bar" 
              height="180" 
              :options="optionsUsia" 
              :series="seriesUsia"
            />
          </ClientOnly>
        </div>
      </div>

      <!-- REKAP BLBA -->
      <div class="card">
        <div class="card-header">
          <div class="card-title">
            <i class="ti ti-wallet" style="font-size:16px;color:var(--hijau)"></i>
            Rekap BLBA per Cabang
          </div>
        </div>
        <div v-if="loadingBlba" class="card-body">
          <div class="skeleton-block" style="width:100%;height:38px;border-radius:6px;margin-bottom:12px;"></div>
          <div v-for="n in 5" :key="n" style="display:flex;gap:8px;align-items:center;margin-bottom:10px;">
            <div class="skeleton-block" style="flex:2;height:12px;"></div>
            <div class="skeleton-block" style="flex:1;height:12px;"></div>
            <div class="skeleton-block" style="flex:1;height:12px;"></div>
            <div class="skeleton-block" style="width:36px;height:18px;border-radius:4px;"></div>
          </div>
        </div>
        <div v-else class="card-body" style="padding:0;">
          <div style="padding:10px 16px;background:#fff8e0;border-bottom:1px solid var(--border);display:flex;justify-content:space-between;align-items:center;">
            <div style="font-size:12px;color:#7a6000;"><strong>Terkumpul:</strong> {{ formatRupiahJuta(financialTerkumpul) }}</div>
            <div style="font-size:12px;color:#7a6000;"><strong>Target:</strong> {{ formatRupiahJuta(financialTarget) }}</div>
            <div style="font-weight:700;color:#e8a020;font-size:13px;">{{ financialPct }}%</div>
          </div>
          <div style="padding:10px 16px;">
            <table class="iuran-table">
              <thead><tr><th>Cabang</th><th>Terkumpul</th><th>Target</th><th>%</th></tr></thead>
              <tbody>
                <tr v-for="i in iuranList" :key="i.cabang">
                  <td><strong>{{ i.cabang }}</strong></td>
                  <td>{{ formatRupiahJuta(i.terkumpul) }}</td>
                  <td>{{ formatRupiahJuta(i.target) }}</td>
                  <td>
                    <span :class="['iuran-pct', i.class]">{{ i.pct }}%</span>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>

      <!-- KEHADIRAN -->
      <div class="card">
        <div class="card-header">
          <div class="card-title">
            <i class="ti ti-calendar" style="font-size:16px;color:var(--hijau)"></i>
            Kehadiran Latihan
          </div>
        </div>
        <div v-if="loadingKehadiran" class="card-body">
          <div class="panel-skeleton-header" style="width:160px;margin-bottom:14px;"></div>
          <div v-for="n in 4" :key="n" class="kehadiran-row" style="gap:10px;margin-bottom:10px;">
            <div class="skeleton-block" style="width:22px;height:22px;border-radius:50%;"></div>
            <div style="flex:1;display:flex;flex-direction:column;gap:5px;">
              <div class="skeleton-block" style="width:90px;height:11px;"></div>
              <div class="skeleton-block" style="width:60px;height:9px;"></div>
            </div>
            <div class="skeleton-block" style="width:80px;height:6px;border-radius:3px;"></div>
            <div class="skeleton-block" style="width:30px;height:12px;"></div>
          </div>
        </div>
        <div v-else class="card-body">
          <div style="font-size:11px;font-weight:700;color:var(--text3);margin-bottom:10px;">Top Unit Latihan Paling Rajin</div>
          <div class="kehadiran-row" v-for="(unit, index) in topUnits" :key="unit.nama">
            <div :class="['rank-num', { top: index < 3 }]">{{ index + 1 }}</div>
            <div class="kunit-name">
              <div style="font-size:12px;font-weight:600;">{{ unit.nama }}</div>
              <div style="font-size:10px;color:var(--text3);">{{ unit.cabang }}</div>
            </div>
            <div class="kunit-bar-wrap">
              <div class="kunit-bar" :style="{ width: unit.pct + '%', background: 'var(--hijau2)' }"></div>
            </div>
            <div class="kunit-pct">{{ unit.pct }}%</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, onMounted } from 'vue'

definePageMeta({ title: 'Dashboard' })

const api = useApi()

// ─── Per-panel loading states (independent async skeletons) ─────────────────
const loadingStats    = ref(true)
const loadingMap      = ref(true)
const loadingRanking  = ref(true)
const loadingTren     = ref(true)
const loadingStatistik = ref(true)
const loadingBlba     = ref(true)
const loadingKehadiran = ref(true)

const activePeriod = ref('bulan_ini')
const totalCabang = ref(0)
const totalAnggota = ref(0)
const totalUnit = ref(0)
const totalProvinsi = ref(5)
const newAnggotaThisWeek = ref(0)
const newCabangThisMonth = ref(0)
const attendanceRate = ref(87)

const periods = [
  { label: 'Bulan Ini', value: 'bulan_ini' },
  { label: 'Bulan Lalu', value: 'bulan_lalu' },
  { label: 'Tahun Ini', value: 'tahun_ini' },
]

// Geographic coordinates of provinces in Indonesia
const geoMap: Record<string, [number, number]> = {
  'Aceh': [-4.6951, 96.7494],
  'Sumatera Utara': [2.1121, 99.1342],
  'Sumatera Barat': [-0.7399, 100.8000],
  'Riau': [0.2933, 101.7068],
  'Jambi': [-1.6186, 102.7725],
  'Sumatera Selatan': [-3.3194, 104.9142],
  'Bengkulu': [-3.7928, 102.2608],
  'Lampung': [-4.5586, 105.4000],
  'Kepulauan Riau': [3.9456, 108.1428],
  'Bangka Belitung': [-2.7410, 106.4406],
  'DKI Jakarta': [-6.2088, 106.8456],
  'Jawa Barat': [-6.9175, 107.6191],
  'Banten': [-6.4058, 106.0640],
  'Jawa Tengah': [-7.0051, 110.4381],
  'DI Yogyakarta': [-7.7956, 110.3695],
  'Jawa Timur': [-7.2575, 112.7521],
  'Sulawesi Selatan': [-5.1477, 119.4327],
  'Bali': [-8.4095, 115.1889]
}

const statGender = ref({ laki_laki: 75, perempuan: 25 })
const statUsia = ref({ under_17: 18, usia_18_25: 32, usia_26_40: 35, over_40: 15 })

const iuranList = ref<any[]>([])
const topUnits = ref<any[]>([])
const rankingKebugaran = ref<any[]>([])
const trenKebugaran = ref<any>({ labels: [], datasets: [] })

const financialTerkumpul = ref(0)
const financialTarget = ref(0)
const financialPct = ref(0)

const kebugaranBaikPct = ref(72)
const kebugaranTotalDites = ref(0)

const totalGender = computed(() => statGender.value.laki_laki + statGender.value.perempuan || 1)
const lakiPct = computed(() => Math.round((statGender.value.laki_laki / totalGender.value) * 100))
const perempuanPct = computed(() => Math.round((statGender.value.perempuan / totalGender.value) * 100))

const totalUsia = computed(() => {
  const val = statUsia.value
  return val.under_17 + val.usia_18_25 + val.usia_26_40 + val.over_40 || 1
})
const under17Pct = computed(() => Math.round((statUsia.value.under_17 / totalUsia.value) * 100))
const usia18_25Pct = computed(() => Math.round((statUsia.value.usia_18_25 / totalUsia.value) * 100))
const usia26_40Pct = computed(() => Math.round((statUsia.value.usia_26_40 / totalUsia.value) * 100))
const over40Pct = computed(() => Math.round((statUsia.value.over_40 / totalUsia.value) * 100))

// ─── APEXCHARTS COMPUTED SERIES & OPTIONS ────────────────────────────────────

const seriesKebugaran = computed(() => {
  if (!trenKebugaran.value.datasets) return []
  return trenKebugaran.value.datasets.map((ds: any) => ({
    name: ds.label.replace('Kota ', ''),
    data: ds.data
  }))
})

const optionsKebugaran = computed(() => ({
  chart: {
    type: 'area',
    toolbar: { show: false },
    zoom: { enabled: false },
    fontFamily: 'Outfit, Inter, sans-serif'
  },
  stroke: { curve: 'smooth', width: 2.5 },
  fill: {
    type: 'gradient',
    gradient: {
      shadeIntensity: 1,
      opacityFrom: 0.25,
      opacityTo: 0.05,
      stops: [0, 90, 100]
    }
  },
  colors: ['#1a5c2a', '#3b82f6', '#f59e0b', '#8b5cf6', '#ef4444'],
  xaxis: {
    categories: trenKebugaran.value.labels || [],
    labels: { style: { colors: '#9ca3af', fontSize: '9px' } }
  },
  yaxis: {
    min: 40,
    max: 100,
    tickAmount: 6,
    labels: {
      formatter: (val: number) => `${Math.round(val)}%`,
      style: { colors: '#9ca3af', fontSize: '9px' }
    }
  },
  tooltip: {
    y: {
      formatter: (val: number) => `${val}% Lulus`
    }
  },
  legend: {
    position: 'bottom',
    fontSize: '11px',
    labels: { colors: 'var(--text2)' }
  }
}))

const seriesGender = computed(() => [
  statGender.value.laki_laki || 0,
  statGender.value.perempuan || 0
])

const optionsGender = computed(() => ({
  chart: {
    type: 'donut',
    fontFamily: 'Outfit, Inter, sans-serif'
  },
  labels: ['Laki-laki', 'Perempuan'],
  colors: ['#1a5c2a', '#5bb8d4'],
  legend: {
    show: true,
    position: 'right',
    fontSize: '11px',
    labels: { colors: 'var(--text2)' }
  },
  plotOptions: {
    pie: {
      donut: {
        size: '65%',
        labels: {
          show: true,
          total: {
            show: true,
            label: 'Total',
            formatter: () => `${totalGender.value}`
          }
        }
      }
    }
  },
  dataLabels: {
    enabled: true,
    formatter: (val: number) => `${Math.round(val)}%`
  }
}))

const seriesUsia = computed(() => [{
  name: 'Persentase',
  data: [under17Pct.value, usia18_25Pct.value, usia26_40Pct.value, over40Pct.value]
}])

const optionsUsia = computed(() => ({
  chart: {
    type: 'bar',
    toolbar: { show: false },
    fontFamily: 'Outfit, Inter, sans-serif'
  },
  plotOptions: {
    bar: {
      horizontal: true,
      barHeight: '55%',
      borderRadius: 4,
      distributed: true
    }
  },
  colors: ['#5bb8d4', '#2e8b3a', '#1a5c2a', '#e8c42a'],
  xaxis: {
    categories: ['<17 tahun', '18–25 thn', '26–40 thn', '>40 tahun'],
    labels: {
      formatter: (val: number) => `${val}%`,
      style: { colors: '#9ca3af', fontSize: '9px' }
    }
  },
  yaxis: {
    labels: {
      style: { colors: 'var(--text2)', fontSize: '11px' }
    }
  },
  dataLabels: {
    enabled: true,
    formatter: (val: number) => `${val}%`
  },
  legend: {
    show: false
  }
}))

const formatRupiahJuta = (val: number) => {
  if (val >= 1000000000) {
    return `Rp ${(val / 1000000000).toFixed(1).replace('.', ',')} M`
  }
  return `Rp ${(val / 1000000).toFixed(1).replace('.', ',')} jt`
}

const getKebugaranColor = (skor: number) => {
  if (skor >= 75) return '#1a5c2a'
  if (skor >= 55) return '#3b82f6'
  return '#ef4444'
}

const getBadgeStyle = (kat: string) => {
  if (kat === 'Baik') return 'border: 1px solid var(--hijau4); background: var(--hijau3); color: var(--hijau);'
  if (kat === 'Cukup') return 'border: 1px solid #e8c42a; background: #fff8e0; color: #7a6000;'
  return 'border: 1px solid var(--merah); background: #fee2e2; color: var(--merah);'
}

useHead({
  link: [
    { rel: 'stylesheet', href: 'https://unpkg.com/leaflet@1.9.4/dist/leaflet.css' }
  ],
  script: [
    { src: 'https://unpkg.com/leaflet@1.9.4/dist/leaflet.js' }
  ]
})

const initLeafletMap = (sebaranList: any[]) => {
  if (!process.client) return

  const checkLeaflet = setInterval(() => {
    if ((window as any).L) {
      clearInterval(checkLeaflet)
      const L = (window as any).L

      // Center around Indonesia
      const map = L.map('leaflet-map', {
        zoomControl: true,
        scrollWheelZoom: false
      }).setView([-2.5489, 118.0149], 5)

      L.tileLayer('https://{s}.basemaps.cartocdn.com/light_all/{z}/{x}/{y}{r}.png', {
        attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors &copy; <a href="https://carto.com/attributions">CARTO</a>',
        subdomains: 'abcd',
        maxZoom: 20
      }).addTo(map)

      sebaranList.forEach((s: any) => {
        const coords = geoMap[s.provinsi]
        if (coords) {
          const count = s.jumlah_anggota
          let radius = 8
          let fillColor = '#4aa855'

          if (count > 1000) {
            radius = 22
            fillColor = '#0d3d1a'
          } else if (count > 500) {
            radius = 16
            fillColor = '#1a5c2a'
          } else if (count > 200) {
            radius = 12
            fillColor = '#2e8b3a'
          } else if (count > 0) {
            radius = 10
            fillColor = '#2e8b3a'
          }

          const marker = L.circleMarker(coords, {
            radius,
            fillColor,
            color: '#fff',
            weight: 2,
            opacity: 1,
            fillOpacity: 0.85
          }).addTo(map)

          // Premium formatted Leaflet tooltip
          marker.bindPopup(`
            <div style="font-family:'Segoe UI',sans-serif; font-size:12px; line-height:1.4; min-width: 150px;">
              <strong style="font-size:13px; color:var(--hijau); display:block; border-bottom:1px solid #eee; padding-bottom:4px; margin-bottom:6px;">${s.provinsi}</strong>
              <div style="display:flex; justify-content:space-between; margin-bottom:2px;">
                <span style="color:#6b6b67;">Total Anggota:</span>
                <strong>${count}</strong>
              </div>
              <div style="display:flex; justify-content:space-between; margin-bottom:2px;">
                <span style="color:#6b6b67;">Cabang Utama:</span>
                <strong>${s.cabang_utama || '-'}</strong>
              </div>
              <div style="display:flex; justify-content:space-between;">
                <span style="color:#6b6b67;">Total Unit:</span>
                <strong>${s.jumlah_unit} unit</strong>
              </div>
            </div>
          `)
        }
      })
    }
  }, 100)
}

// ─── Stat Cards ──────────────────────────────────────────────────────────────
const fetchStats = async () => {
  try {
    loadingStats.value = true
    const res = await api.get('/organization/dashboard-stats')
    totalCabang.value = res.total_cabang || 0
    totalAnggota.value = res.total_anggota || 0
    totalUnit.value = res.total_unit || 0
    newAnggotaThisWeek.value = res.new_anggota_this_week || 0
    newCabangThisMonth.value = res.new_cabang_this_month || 0
    attendanceRate.value = res.attendance_rate || 87

    // Hand off shared data to other panels
    _sharedStats.value = res
  } catch (e) { console.error('fetchStats', e) }
  finally { loadingStats.value = false }
}

// Shared response cache so we only call the API once
const _sharedStats = ref<any>(null)

// ─── Map Panel ───────────────────────────────────────────────────────────────
const fetchMap = async () => {
  try {
    loadingMap.value = true
    const sebaranData = await api.get('/organization/sebaran')
    initLeafletMap(sebaranData)
  } catch (e) { console.error('fetchMap', e) }
  finally { loadingMap.value = false }
}

// ─── Ranking Kebugaran ───────────────────────────────────────────────────────
const fetchRanking = async () => {
  try {
    loadingRanking.value = true
    // Wait for shared stats if not yet loaded
    if (!_sharedStats.value) await fetchStats()
    const res = _sharedStats.value
    rankingKebugaran.value = res?.ranking_kebugaran || []
    if (rankingKebugaran.value.length > 0) {
      const goodCount = rankingKebugaran.value.filter((r: any) => r.kategori === 'Baik' || r.kategori === 'Cukup').length
      kebugaranBaikPct.value = Math.round((goodCount / rankingKebugaran.value.length) * 100)
      kebugaranTotalDites.value = rankingKebugaran.value.length
    } else {
      kebugaranBaikPct.value = 72
      kebugaranTotalDites.value = 247
    }
  } catch (e) { console.error('fetchRanking', e) }
  finally { loadingRanking.value = false }
}

// ─── Tren Kebugaran Chart ────────────────────────────────────────────────────
const fetchTren = async () => {
  try {
    loadingTren.value = true
    if (!_sharedStats.value) await fetchStats()
    trenKebugaran.value = _sharedStats.value?.tren_kebugaran || { labels: [], datasets: [] }
  } catch (e) { console.error('fetchTren', e) }
  finally { loadingTren.value = false }
}

// ─── Statistik Anggota Chart ─────────────────────────────────────────────────
const fetchStatistik = async () => {
  try {
    loadingStatistik.value = true
    if (!_sharedStats.value) await fetchStats()
    const res = _sharedStats.value
    statGender.value = res?.stat_gender || { laki_laki: 75, perempuan: 25 }
    statUsia.value = res?.stat_usia || { under_17: 18, usia_18_25: 32, usia_26_40: 35, over_40: 15 }
  } catch (e) { console.error('fetchStatistik', e) }
  finally { loadingStatistik.value = false }
}

// ─── Rekap BLBA ──────────────────────────────────────────────────────────────
const fetchBlba = async () => {
  try {
    loadingBlba.value = true
    if (!_sharedStats.value) await fetchStats()
    iuranList.value = _sharedStats.value?.rekap_blba || []
    let sumTerkumpul = 0, sumTarget = 0
    iuranList.value.forEach((item: any) => {
      sumTerkumpul += item.terkumpul
      sumTarget += item.target
    })
    financialTerkumpul.value = sumTerkumpul
    financialTarget.value = sumTarget
    financialPct.value = sumTarget > 0 ? Math.round((sumTerkumpul / sumTarget) * 100) : 0
  } catch (e) { console.error('fetchBlba', e) }
  finally { loadingBlba.value = false }
}

// ─── Kehadiran / Top Units ───────────────────────────────────────────────────
const fetchKehadiran = async () => {
  try {
    loadingKehadiran.value = true
    if (!_sharedStats.value) await fetchStats()
    topUnits.value = _sharedStats.value?.top_units || []
  } catch (e) { console.error('fetchKehadiran', e) }
  finally { loadingKehadiran.value = false }
}

// ─── Main orchestrator: fire all panels concurrently ─────────────────────────
const fetchDashboardStats = async () => {
  // Step 1: fetch the single backend stats endpoint first
  await fetchStats()
  // Step 2: fire all other panels in parallel (they use cached _sharedStats)
  Promise.all([
    fetchMap(),
    fetchRanking(),
    fetchTren(),
    fetchStatistik(),
    fetchBlba(),
    fetchKehadiran()
  ])
}

onMounted(() => {
  fetchDashboardStats()
})
</script>

<style scoped>
#page-dashboard { padding: 20px; flex: 1; overflow-y: auto; }

.period-bar { display: flex; align-items: center; gap: 8px; margin-bottom: 20px; flex-wrap: wrap; }
.period-btn { font-size: 12px; padding: 6px 14px; border-radius: 20px; border: 1px solid var(--border); background: var(--card); color: var(--text2); cursor: pointer; transition: all .15s; }
.period-btn.active { background: var(--hijau); color: #fff; border-color: var(--hijau); }
.period-label { font-size: 11px; color: var(--text3); margin-left: 4px; }

.stat-grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 12px; margin-bottom: 20px; }
@media (max-width: 1200px) { .stat-grid { grid-template-columns: repeat(2, 1fr); } }
@media (max-width: 600px) { .stat-grid { grid-template-columns: 1fr; } }

.stat-card { background: var(--card); border: 1px solid var(--border); border-radius: var(--r12); padding: 14px 16px; cursor: pointer; transition: all .2s; position: relative; overflow: hidden; }
.stat-card:hover { box-shadow: 0 4px 16px rgba(0,0,0,0.08); transform: translateY(-2px); }
.stat-card::before { content: ''; position: absolute; top: 0; left: 0; right: 0; height: 3px; }
.stat-card.green::before { background: var(--hijau2); }
.stat-card.yellow::before { background: var(--kuning); }
.stat-card.blue::before { background: var(--biru); }
.stat-card.red::before { background: var(--merah); }
.stat-card.teal::before { background: #2abcaa; }

.stat-icon { width: 36px; height: 36px; border-radius: var(--r8); display: flex; align-items: center; justify-content: center; margin-bottom: 10px; font-size: 16px; }
.stat-val { font-size: 20px; font-weight: 800; color: var(--text1); line-height: 1; }
.stat-label { font-size: 11px; color: var(--text3); margin-top: 4px; }
.stat-trend { font-size: 10px; margin-top: 6px; display: flex; align-items: center; gap: 3px; }
.trend-up { color: #2e8b3a; }
.trend-down { color: var(--merah); }
.stat-sub { font-size: 10px; color: var(--text3); margin-top: 6px; }

.main-grid { display: grid; grid-template-columns: 1fr; gap: 16px; margin-bottom: 16px; }
.bottom-grid { display: grid; grid-template-columns: 1fr 1fr 1fr; gap: 16px; }
@media (max-width: 1000px) { .bottom-grid { grid-template-columns: 1fr; } }

.card { background: var(--card); border: 1px solid var(--border); border-radius: var(--r12); overflow: hidden; display: flex; flex-direction: column; }
.card-header { padding: 14px 16px; border-bottom: 1px solid var(--border); display: flex; align-items: center; justify-content: space-between; background: var(--card); }
.card-title { font-size: 13px; font-weight: 700; display: flex; align-items: center; gap: 6px; }
.card-title svg { width: 15px; height: 15px; fill: var(--hijau); }
.card-action { font-size: 11px; color: var(--hijau2); cursor: pointer; }
.card-body { padding: 14px 16px; }

.map-wrap { position: relative; background: #e8f4f8; overflow: hidden; width: 100%; }

/* Donut Chart */
.donut-wrap { display: flex; align-items: center; gap: 14px; }
.donut-svg { flex-shrink: 0; }
.donut-legend { display: flex; flex-direction: column; gap: 6px; }
.donut-leg-row { display: flex; align-items: center; gap: 6px; font-size: 11px; color: var(--text2); }
.donut-dot { width: 8px; height: 8px; border-radius: 50%; }

/* Chart Bars */
.chart-bar-row { display: flex; align-items: center; gap: 8px; margin-bottom: 8px; }
.chart-label { font-size: 11px; color: var(--text2); width: 80px; flex-shrink: 0; text-align: right; }
.chart-bar-wrap { flex: 1; height: 18px; background: var(--bg); border-radius: 4px; overflow: hidden; position: relative; }
.chart-bar-fill { height: 100%; border-radius: 4px; transition: width .6s ease; }
.chart-val { font-size: 11px; font-weight: 600; color: var(--text1); width: 42px; text-align: right; }

/* Iuran Table */
.iuran-table { width: 100%; border-collapse: collapse; font-size: 11px; }
.iuran-table th { text-align: left; color: var(--text3); font-weight: 700; padding: 6px 8px; border-bottom: 1px solid var(--border); text-transform: uppercase; }
.iuran-table td { padding: 8px; border-bottom: 1px solid var(--border); }
.iuran-table tr:last-child td { border-bottom: none; }
.iuran-pct { font-weight: 700; }
.iuran-pct.high { color: var(--hijau2); }
.iuran-pct.mid { color: #e8a020; }
.iuran-pct.low { color: var(--merah); }

/* Kehadiran Rows */
.kehadiran-row { display: flex; align-items: center; gap: 8px; margin-bottom: 10px; }
.rank-num { width: 20px; height: 20px; border-radius: 50%; display: flex; align-items: center; justify-content: center; font-size: 10px; font-weight: 800; background: var(--bg); color: var(--text2); }
.rank-num.top { background: #f59e0b; color: #fff; }
.kunit-name { width: 100px; flex-shrink: 0; }
.kunit-bar-wrap { flex: 1; height: 6px; background: var(--bg); border-radius: 3px; overflow: hidden; }
.kunit-bar { height: 100%; border-radius: 3px; }
.kunit-pct { font-size: 11px; font-weight: 700; width: 32px; text-align: right; }

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

/* Premium shimmer skeleton block */
.skeleton-block {
  position: relative;
  overflow: hidden;
  background: linear-gradient(90deg, #f0f0f0 25%, #e0e0e0 50%, #f0f0f0 75%);
  background-size: 200% 100%;
  animation: skeleton-shimmer 1.4s infinite;
  border-radius: 6px;
  min-height: 10px;
}

@keyframes skeleton-shimmer {
  0% { background-position: 200% 0; }
  100% { background-position: -200% 0; }
}

/* Full-panel skeleton placeholder */
.skeleton-panel {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 14px;
  background: linear-gradient(135deg, #f5f9f5 0%, #edf4ed 100%);
  border-radius: 0 0 12px 12px;
  width: 100%;
}

.skeleton-pulse-icon {
  animation: pulse-opacity 1.8s ease-in-out infinite;
}

.skeleton-pulse-text {
  font-size: 12px;
  color: #9ca3af;
  font-weight: 500;
  animation: pulse-opacity 1.8s ease-in-out infinite;
  letter-spacing: 0.03em;
}

@keyframes pulse-opacity {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.35; }
}

/* Section header placeholder line */
.panel-skeleton-header {
  height: 13px;
  background: linear-gradient(90deg, #e8e8e8 25%, #d4d4d4 50%, #e8e8e8 75%);
  background-size: 200% 100%;
  animation: skeleton-shimmer 1.4s infinite;
  border-radius: 6px;
  width: 140px;
  margin-bottom: 14px;
}

</style>
