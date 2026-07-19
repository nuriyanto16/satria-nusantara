<template>
  <div id="page-dokumentasi" class="doc-split-layout">
    <!-- ===== PANEL KIRI: DAFTAR DOKUMEN & SUB-MENU ===== -->
    <div class="panel-left">
      <div class="pl-head">
        <div class="pl-title">Dokumentasi</div>
        <div class="pl-sub">Panduan arsitektur & referensi teknis</div>
      </div>
      <div class="pl-scroll">
        <!-- Main Navigation Groups -->
        <div class="doc-menu-group">
          <div class="dmg-title">Panduan Umum</div>
          <button :class="['dmg-item', { active: activeTab === 'stack' }]" @click="switchTab('stack')">
            <div class="dmg-icon"><i class="ti ti-layers-intersect"></i></div>
            <div class="dmg-info">
              <div class="dmg-name">Tech Stack</div>
              <div class="dmg-desc">Spesifikasi core teknologi</div>
            </div>
          </button>
          
          <button :class="['dmg-item', { active: activeTab === 'guide' }]" @click="switchTab('guide')">
            <div class="dmg-icon"><i class="ti ti-terminal-2"></i></div>
            <div class="dmg-info">
              <div class="dmg-name">Setup Development</div>
              <div class="dmg-desc">Instalasi lokal & panduan Docker</div>
            </div>
          </button>
        </div>

        <div class="doc-menu-group">
          <div class="dmg-title">Arsitektur Data</div>
          <button :class="['dmg-item', { active: activeTab === 'db' }]" @click="switchTab('db')">
            <div class="dmg-icon"><i class="ti ti-database"></i></div>
            <div class="dmg-info">
              <div class="dmg-name">Skema Database</div>
              <div class="dmg-desc">Entitas relasional & tipe data</div>
            </div>
          </button>

          <button :class="['dmg-item', { active: activeTab === 'api' }]" @click="switchTab('api')">
            <div class="dmg-icon"><i class="ti ti-api"></i></div>
            <div class="dmg-info">
              <div class="dmg-name">Referensi API</div>
              <div class="dmg-desc">Endpoint, payload, & respon</div>
            </div>
          </button>
        </div>

        <!-- System Architecture Card -->
        <div class="arch-card">
          <div class="ac-title"><i class="ti ti-git-fork"></i> Alur Data Sistem</div>
          <div class="ac-desc">
            API Gateway bertindak sebagai reverse proxy ke microservices (Auth, Org, Training) dengan Postgres sebagai DB utama.
          </div>
        </div>
      </div>
    </div>

    <!-- ===== PANEL KANAN: DOKUMEN CONTENT ===== -->
    <div class="panel-right">
      <div class="pr-toolbar">
        <div>
          <div class="pr-title">{{ getTabTitle(activeTab) }}</div>
          <div class="pr-meta">Developer Documentation · v1.0.0-Stable</div>
        </div>
        <button class="btn btn-outline btn-sm" @click="copyDocLink"><i class="ti ti-link"></i> Salin Link</button>
      </div>

      <div class="pr-scroll">
        <!-- 1. TAB TECH STACK -->
        <div v-if="activeTab === 'stack'" class="doc-pane animate-fade">
          <p class="pane-intro">
            Aplikasi administrasi Satria Nusantara dibangun menggunakan arsitektur microservices terintegrasi yang andal, scalable, dan modern. Berikut adalah rincian stack teknologi yang digunakan:
          </p>

          <div class="tech-grid">
            <div class="tech-card">
              <div class="tc-header">
                <div class="tc-icon nuxt"><i class="ti ti-brand-vue"></i></div>
                <div>
                  <div class="tc-title">Nuxt.js 3</div>
                  <div class="tc-subtitle">Vue.js SSR Framework</div>
                </div>
              </div>
              <p class="tc-desc">Digunakan untuk membangun Web Admin Console. Menggunakan sistem SSR (Server-Side Rendering) untuk rendering yang instan, SEO friendly, dan struktur folder berbasis file-routing.</p>
              <div class="tc-meta">Versi: 3.11+ · Engine: Vue 3 + Vite</div>
            </div>

            <div class="tech-card">
              <div class="tc-header">
                <div class="tc-icon go"><i class="ti ti-code"></i></div>
                <div>
                  <div class="tc-title">Go (Golang)</div>
                  <div class="tc-subtitle">Backend Microservices</div>
                </div>
              </div>
              <p class="tc-desc">Seluruh backend API dibangun menggunakan bahasa Go untuk memastikan eksekusi performa tinggi, konkurensi goroutine yang efisien, serta konsumsi memori minimum.</p>
              <div class="tc-meta">Versi: 1.21+ · Router: chi-router</div>
            </div>

            <div class="tech-card">
              <div class="tc-header">
                <div class="tc-icon pg"><i class="ti ti-database"></i></div>
                <div>
                  <div class="tc-title">PostgreSQL 15</div>
                  <div class="tc-subtitle">Relational Database</div>
                </div>
              </div>
              <p class="tc-desc">Penyimpanan data relasional terpusat untuk menyimpan profil anggota, histori iuran BLBA, rekap absensi, dan data pelatihan secara konsisten dengan transaksi ACID.</p>
              <div class="tc-meta">Database Engine: PostgreSQL</div>
            </div>

            <div class="tech-card">
              <div class="tc-header">
                <div class="tc-icon redis"><i class="ti ti-bolt"></i></div>
                <div>
                  <div class="tc-title">Redis Cache</div>
                  <div class="tc-subtitle">Session & Temporary Cache</div>
                </div>
              </div>
              <p class="tc-desc">Berfungsi mempercepat verifikasi sesi token JWT, menyimpan sementara data QR absensi, serta caching data laporan eksekutif agar mempercepat respon backend.</p>
              <div class="tc-meta">Versi: 7-Alpine · In-Memory DB</div>
            </div>

            <div class="tech-card">
              <div class="tc-header">
                <div class="tc-icon docker"><i class="ti ti-brand-docker"></i></div>
                <div>
                  <div class="tc-title">Docker Containers</div>
                  <div class="tc-subtitle">DevOps & Infrastructure</div>
                </div>
              </div>
              <p class="tc-desc">Memastikan keselarasan lingkungan pengembangan dari lokal hingga staging. Terbungkus rapi menggunakan Docker Compose untuk instalasi sekali klik.</p>
              <div class="tc-meta">Compose Schema: v3.8</div>
            </div>

            <div class="tech-card">
              <div class="tc-header">
                <div class="tc-icon flutter"><i class="ti ti-device-mobile"></i></div>
                <div>
                  <div class="tc-title">Flutter SDK</div>
                  <div class="tc-subtitle">Mobile App Companion</div>
                </div>
              </div>
              <p class="tc-desc">Aplikasi mobile companion untuk pelatih (pemindai absensi QR & input tes kebugaran) dan anggota (kartu digital & pembayaran iuran BLBA).</p>
              <div class="tc-meta">Versi: 3.19+ · Dart Language</div>
            </div>
          </div>
        </div>

        <!-- 2. TAB SETUP DEVELOPMENT -->
        <div v-else-if="activeTab === 'guide'" class="doc-pane animate-fade">
          <div class="guide-tabs">
            <button :class="['guide-tab-btn', { active: activeGuideTab === 'docker' }]" @click="activeGuideTab = 'docker'">Dengan Docker (Rekomendasi)</button>
            <button :class="['guide-tab-btn', { active: activeGuideTab === 'local' }]" @click="activeGuideTab = 'local'">Lokal Manual (Tanpa Docker)</button>
          </div>

          <div v-if="activeGuideTab === 'docker'" class="guide-content">
            <h3 class="guide-section-title">1. Persiapan Awal</h3>
            <p class="guide-text">Pastikan Anda telah menginstal **Docker Desktop** dan **Docker Compose** di komputer Anda.</p>

            <h3 class="guide-section-title">2. Jalankan Cluster Database & Services</h3>
            <p class="guide-text">Jalankan perintah berikut di direktori root project untuk mengompilasi dan mengaktifkan semua service kontainer di latar belakang:</p>
            <div class="code-box">
              <pre><code>docker-compose up -d --build</code></pre>
              <button class="btn-copy" @click="copyCode('docker-compose up -d --build')"><i class="ti ti-copy"></i></button>
            </div>

            <h3 class="guide-section-title">3. Verifikasi Status Container</h3>
            <div class="code-box">
              <pre><code>docker-compose ps</code></pre>
              <button class="btn-copy" @click="copyCode('docker-compose ps')"><i class="ti ti-copy"></i></button>
            </div>

            <h3 class="guide-section-title">4. Alamat Akses Layanan</h3>
            <table class="data-table">
              <thead>
                <tr><th>Layanan</th><th>URL Akses</th><th>Keterangan</th></tr>
              </thead>
              <tbody>
                <tr><td>Web Admin Console</td><td><a href="http://localhost:3000" target="_blank" class="link-url">http://localhost:3000</a></td><td>Aplikasi Web Utama</td></tr>
                <tr><td>API Gateway (Reverse Proxy)</td><td><a href="http://localhost:8080" target="_blank" class="link-url">http://localhost:8080</a></td><td>Semua API Routing</td></tr>
                <tr><td>pgAdmin DB GUI</td><td><a href="http://localhost:5050" target="_blank" class="link-url">http://localhost:5050</a></td><td>Email: admin@satria-nusantara.org</td></tr>
              </tbody>
            </table>
          </div>

          <div v-else class="guide-content">
            <h3 class="guide-section-title">1. Setup Database PostgreSQL lokal</h3>
            <p class="guide-text">Eksekusi skrip inisialisasi tabel SQL dari file <code class="code-tag">infra/sql/init.sql</code> ke database Postgres Anda.</p>

            <h3 class="guide-section-title">2. Menjalankan Backend API (Golang)</h3>
            <div class="code-box">
              <pre><code>cd backend
go mod tidy
go run main.go</code></pre>
              <button class="btn-copy" @click="copyCode('cd backend && go run main.go')"><i class="ti ti-copy"></i></button>
            </div>

            <h3 class="guide-section-title">3. Menjalankan Frontend Web (Nuxt 3)</h3>
            <div class="code-box">
              <pre><code>cd frontend
npm install
npm run dev</code></pre>
              <button class="btn-copy" @click="copyCode('cd frontend && npm run dev')"><i class="ti ti-copy"></i></button>
            </div>
          </div>
        </div>

        <!-- 3. TAB DATABASE SCHEMA -->
        <div v-else-if="activeTab === 'db'" class="doc-pane db-layout animate-fade">
          <!-- Left Table list tree -->
          <div class="db-tree">
            <div class="db-tree-title">Daftar Tabel</div>
            <button 
              v-for="(tInfo, tName) in dbTables" 
              :key="tName"
              :class="['db-tree-item', { active: selectedTable === tName }]"
              @click="selectedTable = tName"
            >
              <i class="ti ti-table"></i> {{ tName }}
            </button>
          </div>

          <!-- Right Table detail columns -->
          <div class="db-details" v-if="selectedTable && dbTables[selectedTable]">
            <div class="db-table-head">
              <div class="db-table-name"><i class="ti ti-table-alias"></i> {{ selectedTable }}</div>
              <p class="db-table-desc">{{ dbTables[selectedTable].desc }}</p>
            </div>
            
            <table class="data-table">
              <thead>
                <tr>
                  <th>Kolom</th>
                  <th>Tipe Data</th>
                  <th>Atribut</th>
                  <th>Deskripsi</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="col in dbTables[selectedTable].cols" :key="col.name">
                  <td class="col-name">{{ col.name }}</td>
                  <td class="col-type">{{ col.type }}</td>
                  <td>
                    <span v-if="col.pk" class="col-attr pk">Primary Key</span>
                    <span v-if="col.fk" class="col-attr fk">Foreign Key</span>
                    <span v-if="!col.nullable" class="col-attr nn">Not Null</span>
                  </td>
                  <td style="font-size: 11px;">{{ col.desc }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <!-- 4. TAB REFERENSI API -->
        <div v-else-if="activeTab === 'api'" class="doc-pane api-layout animate-fade">
          <div class="api-tree">
            <div class="db-tree-title">Kategori API</div>
            <div v-for="cat in apiEndpoints" :key="cat.category" class="api-cat-group">
              <div class="api-cat-title">{{ cat.category }}</div>
              <button 
                v-for="ep in cat.items" 
                :key="ep.path"
                :class="['api-tree-item', { active: selectedEndpoint && selectedEndpoint.path === ep.path }]"
                @click="selectedEndpoint = ep"
              >
                <span :class="['method-badge', ep.method]">{{ ep.method }}</span>
                <span class="ep-url-text">{{ ep.path }}</span>
              </button>
            </div>
          </div>

          <div class="api-details" v-if="selectedEndpoint">
            <div class="api-detail-head">
              <div style="display:flex;align-items:center;gap:10px;">
                <span :class="['method-badge large', selectedEndpoint.method]">{{ selectedEndpoint.method }}</span>
                <div class="api-detail-path">{{ selectedEndpoint.path }}</div>
              </div>
              <p class="api-detail-desc">{{ selectedEndpoint.desc }}</p>
            </div>

            <!-- Parameters table -->
            <div class="api-section-title">Query Parameters / URL Param</div>
            <table class="data-table" v-if="selectedEndpoint.params && selectedEndpoint.params.length > 0">
              <thead>
                <tr><th>Parameter</th><th>Tipe</th><th>Wajib</th><th>Deskripsi</th></tr>
              </thead>
              <tbody>
                <tr v-for="p in selectedEndpoint.params" :key="p.name">
                  <td style="font-weight: 700;">{{ p.name }}</td>
                  <td class="col-type">{{ p.type }}</td>
                  <td><span :class="['col-attr', p.required ? 'nn' : '']">{{ p.required ? 'Ya' : 'Tidak' }}</span></td>
                  <td style="font-size: 11px;">{{ p.desc }}</td>
                </tr>
              </tbody>
            </table>
            <div v-else class="api-no-param">Tidak memerlukan parameter khusus.</div>

            <!-- Payload json body -->
            <div v-if="selectedEndpoint.payload" style="margin-top:16px;">
              <div class="api-section-title">Request Body (JSON)</div>
              <div class="code-box-api">
                <pre><code>{{ JSON.stringify(selectedEndpoint.payload, null, 2) }}</code></pre>
              </div>
            </div>

            <!-- Response json body -->
            <div style="margin-top:16px;">
              <div class="api-section-title">Response Schema (Success 200)</div>
              <div class="code-box-api">
                <pre><code>{{ JSON.stringify(selectedEndpoint.response, null, 2) }}</code></pre>
              </div>
            </div>

            <!-- cURL snippet -->
            <div style="margin-top:16px;">
              <div class="api-section-title">cURL CLI Request</div>
              <div class="code-box">
                <pre><code>{{ generateCurl(selectedEndpoint) }}</code></pre>
                <button class="btn-copy" @click="copyCode(generateCurl(selectedEndpoint))"><i class="ti ti-copy"></i></button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({ title: 'Dokumentasi Sistem' })

const activeTab = ref('stack')
const activeGuideTab = ref('docker')
const selectedTable = ref('users')
const selectedEndpoint = ref<any>(null)

const switchTab = (tab: string) => {
  activeTab.value = tab
}

const getTabTitle = (tab: string) => {
  if (tab === 'stack') return 'Tech Stack & Spesifikasi'
  if (tab === 'guide') return 'Setup & Panduan Development'
  if (tab === 'db') return 'Arsitektur & Skema Database'
  return 'Referensi API & Kontrak Service'
}

// Mock Database schema
const dbTables: Record<string, any> = {
  users: {
    desc: 'Menyimpan kredensial login akun pengurus, pelatih, administrator pusat, dan tingkat hak akses ke sistem web admin.',
    cols: [
      { name: 'id', type: 'BIGSERIAL', pk: true, nullable: false, desc: 'Unique auto-incrementing record ID' },
      { name: 'name', type: 'VARCHAR(255)', nullable: false, desc: 'Nama lengkap pengurus' },
      { name: 'email', type: 'VARCHAR(255)', nullable: false, desc: 'Alamat email instansi/pribadi (digunakan untuk login)' },
      { name: 'password_hash', type: 'VARCHAR(255)', nullable: false, desc: 'Argon2id/Bcrypt hash password' },
      { name: 'role', type: 'VARCHAR(50)', nullable: false, desc: 'Hak akses tingkat otorisasi (admin, pengurus_cabang, pic_unit, pelatih)' },
      { name: 'cabang_id', type: 'BIGINT', fk: true, nullable: true, desc: 'Membatasi lingkup akses kerja pengurus ke cabang tertentu' }
    ]
  },
  cabang: {
    desc: 'Daftar pengurus cabang resmi Satria Nusantara di tingkat Kota/Kabupaten di Indonesia.',
    cols: [
      { name: 'id', type: 'BIGSERIAL', pk: true, nullable: false, desc: 'ID internal cabang' },
      { name: 'kode', type: 'VARCHAR(10)', nullable: false, desc: 'Kode cabang 3 huruf unik (contoh: YGY, BDG)' },
      { name: 'nama', type: 'VARCHAR(255)', nullable: false, desc: 'Nama resmi cabang administratif' },
      { name: 'provinsi_id', type: 'INT', nullable: false, desc: 'ID referensi provinsi cakupan kerja' },
      { name: 'kota_kab', type: 'VARCHAR(255)', nullable: false, desc: 'Wilayah Kota/Kabupaten sekretariat cabang' }
    ]
  },
  unit: {
    desc: 'Daftar unit latihan formal di bawah naungan kepengurusan cabang tingkat lokal.',
    cols: [
      { name: 'id', type: 'BIGSERIAL', pk: true, nullable: false, desc: 'ID unit latihan' },
      { name: 'cabang_id', type: 'BIGINT', fk: true, nullable: false, desc: 'Relasi ke cabang naungan' },
      { name: 'nama', type: 'VARCHAR(255)', nullable: false, desc: 'Nama unit latihan (contoh: Unit Malioboro)' },
      { name: 'lokasi', type: 'TEXT', nullable: false, desc: 'Alamat atau deskripsi tempat latihan rutin' },
      { name: 'jadwal_rutin', type: 'VARCHAR(255)', nullable: false, desc: 'Deskripsi jadwal hari & jam latihan tetap' }
    ]
  },
  anggota: {
    desc: 'Menyimpan profil lengkap anggota Satria Nusantara beserta tingkatan sabuk dan status keaktifannya.',
    cols: [
      { name: 'id', type: 'BIGSERIAL', pk: true, nullable: false, desc: 'ID anggota' },
      { name: 'unit_id', type: 'BIGINT', fk: true, nullable: false, desc: 'Relasi ke unit latihan tempat berlatih' },
      { name: 'nomor_anggota', type: 'VARCHAR(50)', nullable: false, desc: 'Kode anggota nasional (contoh: YO-YGY-00142)' },
      { name: 'nama_lengkap', type: 'VARCHAR(255)', nullable: false, desc: 'Nama lengkap anggota sesuai KTP' },
      { name: 'jenis_kelamin', type: 'CHAR(1)', nullable: false, desc: 'Kode L (Laki-laki) atau P (Perempuan)' },
      { name: 'tingkatan', type: 'VARCHAR(50)', nullable: false, desc: 'Peringkat sabuk saat ini (Pra Dasar, Dasar, PH, PK)' }
    ]
  }
}

// Mock API Endpoints
const apiEndpoints = [
  {
    category: 'Authentication API',
    items: [
      {
        method: 'POST',
        path: '/api/v1/auth/login',
        desc: 'Melakukan login otentikasi menggunakan email dan kata sandi untuk memperoleh JWT token akses.',
        params: [],
        payload: { email: 'admin@satria-nusantara.org', password: 'password123' },
        response: { status: 'success', token: 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...', user: { id: 1, name: 'Sri Astuti', role: 'admin' } }
      },
      {
        method: 'GET',
        path: '/api/v1/auth/me',
        desc: 'Memperoleh info profil pengurus yang sedang login berdasarkan token JWT di Header.',
        params: [
          { name: 'Authorization', type: 'Header', required: true, desc: 'Bearer JWT_TOKEN_STRING' }
        ],
        payload: null,
        response: { status: 'success', user: { id: 1, name: 'Sri Astuti', email: 'sri.astuti@sn-pusat.org', role: 'admin' } }
      }
    ]
  },
  {
    category: 'Organization API',
    items: [
      {
        method: 'GET',
        path: '/api/v1/organization/cabang',
        desc: 'Mengambil daftar seluruh cabang terdaftar dengan opsi pencarian kata kunci.',
        params: [
          { name: 'search', type: 'Query', required: false, desc: 'Cari berdasarkan nama/kota cabang' }
        ],
        payload: null,
        response: { total: 1, data: [{ id: 1, kode: 'YGY', nama: 'Kota Yogyakarta', kota_kab: 'Yogyakarta', status: 'aktif' }] }
      },
      {
        method: 'POST',
        path: '/api/v1/organization/cabang',
        desc: 'Membuat entitas cabang baru di lingkup administratif nasional.',
        params: [],
        payload: { kode: 'BDG', nama: 'Kota Bandung', provinsi_id: 3, kota_kab: 'Bandung', alamat: 'Jl. Dago No 15' },
        response: { status: 'success', message: 'Cabang baru berhasil dibuat', id: 2 }
      }
    ]
  }
]

// Auto select first endpoint on load
onMounted(() => {
  if (apiEndpoints.length > 0 && apiEndpoints[0].items.length > 0) {
    selectedEndpoint.value = apiEndpoints[0].items[0]
  }
})

// Helpers
const generateCurl = (ep: any) => {
  let curl = `curl -X ${ep.method} http://localhost:8080${ep.path}`
  if (ep.method === 'POST') {
    curl += ` \\\n  -H "Content-Type: application/json"`
    if (ep.payload) {
      curl += ` \\\n  -d '${JSON.stringify(ep.payload)}'`
    }
  }
  return curl
}

const copyCode = (text: string) => {
  navigator.clipboard.writeText(text)
  alert('Kode perintah/URL disalin ke clipboard!')
}

const copyDocLink = () => {
  navigator.clipboard.writeText(window.location.href)
  alert('Link halaman dokumentasi disalin!')
}
</script>

<style scoped>
.doc-split-layout {
  display: flex;
  height: calc(100vh - 64px - 44px);
  background: var(--bg);
  overflow: hidden;
  margin: -16px;
}

/* Left Sidebar Menu */
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
.pl-scroll { flex: 1; overflow-y: auto; padding: 16px; }

.doc-menu-group { margin-bottom: 18px; }
.dmg-title { font-size: 10px; font-weight: 800; color: var(--text3); text-transform: uppercase; margin-bottom: 8px; letter-spacing: .06em; }
.dmg-item { display: flex; align-items: center; gap: 10px; padding: 8px 12px; border-radius: var(--r8); width: 100%; background: none; border: none; cursor: pointer; text-align: left; transition: all .15s; margin-bottom: 4px; }
.dmg-item:hover { background: var(--surface); }
.dmg-item.active { background: var(--hijau3); }
.dmg-icon { width: 30px; height: 30px; border-radius: var(--r6); background: var(--surface); color: var(--text2); display: flex; align-items: center; justify-content: center; font-size: 14px; flex-shrink: 0; }
.dmg-item.active .dmg-icon { background: var(--card); color: var(--hijau); }
.dmg-info { min-width: 0; flex: 1; }
.dmg-name { font-size: 12px; font-weight: 700; color: var(--text1); }
.dmg-desc { font-size: 10px; color: var(--text3); white-space: nowrap; overflow: hidden; text-overflow: ellipsis; margin-top: 1px; }

.arch-card { margin-top: 18px; padding: 12px; background: var(--surface); border: 0.5px solid var(--border); border-radius: var(--r8); }
.ac-title { font-size: 11px; font-weight: 700; color: var(--hijau); display: flex; align-items: center; gap: 4px; margin-bottom: 6px; }
.ac-desc { font-size: 10px; color: var(--text3); line-height: 1.4; }

/* Right Panel Viewport */
.panel-right { flex: 1; display: flex; flex-direction: column; height: 100%; min-width: 0; background: var(--surface); }
.pr-toolbar { display: flex; justify-content: space-between; align-items: center; padding: 14px 20px; border-bottom: 1px solid var(--border); background: var(--card); flex-shrink: 0; }
.pr-title { font-size: 14px; font-weight: 800; }
.pr-meta { font-size: 11px; color: var(--text3); margin-top: 2px; }
.btn { display: inline-flex; align-items: center; gap: 6px; padding: 8px 14px; border-radius: var(--r8); font-size: 12px; font-weight: 600; cursor: pointer; border: none; transition: all .15s; }
.btn-outline { background: #fff; color: var(--text2); border: 1px solid var(--border); }
.btn-outline:hover { background: var(--surface); }

.pr-scroll { flex: 1; overflow-y: auto; padding: 20px; }
.doc-pane { display: flex; flex-direction: column; gap: 16px; }
.pane-intro { font-size: 13px; color: var(--text2); line-height: 1.6; }

/* Tech Grid */
.tech-grid { display: grid; grid-template-columns: repeat(2, 1fr); gap: 14px; }
.tech-card { background: var(--card); border: 1px solid var(--border); border-radius: var(--r12); padding: 16px; display: flex; flex-direction: column; }
.tc-header { display: flex; align-items: center; gap: 12px; margin-bottom: 12px; }
.tc-icon { width: 36px; height: 36px; border-radius: var(--r8); display: flex; align-items: center; justify-content: center; font-size: 18px; color: #fff; }
.tc-icon.nuxt { background: #00dc82; }
.tc-icon.go { background: #00add8; }
.tc-icon.pg { background: #336791; }
.tc-icon.redis { background: #d82c20; }
.tc-icon.docker { background: #2496ed; }
.tc-icon.flutter { background: #02569b; }
.tc-title { font-size: 13px; font-weight: 700; color: var(--text1); }
.tc-subtitle { font-size: 10px; color: var(--text3); }
.tc-desc { font-size: 11px; color: var(--text2); line-height: 1.5; margin-bottom: 14px; }
.tc-meta { font-size: 10px; color: var(--text3); border-top: 1px solid var(--border); padding-top: 8px; margin-top: auto; }

/* Guide Tab / Setup styles */
.guide-tabs { display: flex; gap: 6px; border-bottom: 1px solid var(--border); padding-bottom: 8px; margin-bottom: 10px; }
.guide-tab-btn { padding: 6px 12px; font-size: 11px; font-weight: 600; border: none; background: none; cursor: pointer; color: var(--text2); border-bottom: 2px solid transparent; }
.guide-tab-btn.active { color: var(--hijau); border-bottom-color: var(--hijau); }
.guide-content { display: flex; flex-direction: column; gap: 10px; }
.guide-section-title { font-size: 12px; font-weight: 700; margin-top: 10px; }
.guide-text { font-size: 12px; color: var(--text2); line-height: 1.5; }

.code-box { display: flex; background: var(--surface); border: 1px solid var(--border); border-radius: var(--r8); padding: 10px 14px; align-items: center; justify-content: space-between; position: relative; }
.code-box pre { font-family: monospace; font-size: 11px; color: var(--text1); margin: 0; white-space: pre-wrap; word-break: break-all; }
.btn-copy { background: none; border: none; cursor: pointer; color: var(--text3); font-size: 14px; padding: 4px; border-radius: 4px; }
.btn-copy:hover { color: var(--hijau); }
.code-tag { background: var(--surface); border: 0.5px solid var(--border); padding: 1px 4px; border-radius: 4px; font-family: monospace; font-size: 11px; }

/* Database tree/layout styles */
.db-layout { display: flex; gap: 16px; flex-direction: row; height: 100%; align-items: flex-start; }
.db-tree { width: 180px; flex-shrink: 0; background: var(--card); border: 1px solid var(--border); border-radius: var(--r12); padding: 12px; }
.db-tree-title { font-size: 10px; font-weight: 800; color: var(--text3); text-transform: uppercase; margin-bottom: 8px; letter-spacing: 0.5px; }
.db-tree-item { display: flex; align-items: center; gap: 6px; padding: 6px 8px; font-size: 11px; font-weight: 600; width: 100%; border: none; background: none; cursor: pointer; text-align: left; border-radius: 6px; color: var(--text2); margin-bottom: 2px; }
.db-tree-item:hover { background: var(--surface); }
.db-tree-item.active { background: var(--hijau3); color: var(--hijau); }

.db-details { flex: 1; min-width: 0; background: var(--card); border: 1px solid var(--border); border-radius: var(--r12); padding: 16px; }
.db-table-head { border-bottom: 1px solid var(--border); padding-bottom: 12px; margin-bottom: 12px; }
.db-table-name { font-size: 14px; font-weight: 800; display: flex; align-items: center; gap: 6px; color: var(--text1); }
.db-table-name i { color: var(--hijau); }
.db-table-desc { font-size: 11px; color: var(--text3); margin-top: 4px; line-height: 1.4; }

.data-table { width: 100%; border-collapse: collapse; }
.data-table th { padding: 8px 12px; background: var(--surface); color: var(--text3); font-size: 10px; font-weight: 700; text-transform: uppercase; text-align: left; border-bottom: 1px solid var(--border); }
.data-table td { padding: 10px 12px; border-bottom: 1px solid var(--border); font-size: 12px; vertical-align: middle; }
.data-table tfoot td { font-weight: 700; background: var(--surface); border-top: 1.5px solid var(--border); }
.link-url { color: var(--hijau); font-weight: 600; text-decoration: none; }
.link-url:hover { text-decoration: underline; }

.col-name { font-weight: 600; font-family: monospace; font-size: 11px; }
.col-type { font-family: monospace; color: #b45309; font-size: 11px; }
.col-attr { font-size: 9px; font-weight: 700; display: inline-block; padding: 1px 4px; border-radius: 4px; margin-right: 4px; }
.col-attr.pk { background: #fef3c7; color: #d97706; }
.col-attr.fk { background: #dbeafe; color: #2563eb; }
.col-attr.nn { background: #fee2e2; color: #dc2626; }

/* API Layout styles */
.api-layout { display: flex; gap: 16px; flex-direction: row; height: 100%; align-items: flex-start; }
.api-tree { width: 220px; flex-shrink: 0; background: var(--card); border: 1px solid var(--border); border-radius: var(--r12); padding: 12px; max-height: 100%; overflow-y: auto; }
.api-cat-group { margin-bottom: 12px; }
.api-cat-title { font-size: 9px; font-weight: 800; color: var(--text3); text-transform: uppercase; margin-bottom: 6px; letter-spacing: 0.5px; border-bottom: 0.5px solid var(--border); padding-bottom: 2px; }
.api-tree-item { display: flex; align-items: center; gap: 6px; padding: 5px 6px; font-size: 10px; width: 100%; border: none; background: none; cursor: pointer; text-align: left; border-radius: 6px; color: var(--text2); margin-bottom: 2px; }
.api-tree-item:hover { background: var(--surface); }
.api-tree-item.active { background: var(--hijau3); color: var(--hijau); }
.ep-url-text { font-family: monospace; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; flex: 1; }

.method-badge { font-size: 8px; font-weight: 800; padding: 2px 4px; border-radius: 4px; color: #fff; width: 36px; text-align: center; flex-shrink: 0; }
.method-badge.GET { background: var(--hijau2); }
.method-badge.POST { background: var(--biru); }
.method-badge.large { font-size: 10px; padding: 4px 8px; width: 48px; }

.api-details { flex: 1; min-width: 0; background: var(--card); border: 1px solid var(--border); border-radius: var(--r12); padding: 16px; }
.api-detail-head { border-bottom: 1px solid var(--border); padding-bottom: 12px; margin-bottom: 14px; }
.api-detail-path { font-size: 15px; font-family: monospace; font-weight: 700; color: var(--text1); }
.api-detail-desc { font-size: 12px; color: var(--text2); margin-top: 4px; line-height: 1.4; }
.api-section-title { font-size: 10px; font-weight: 800; color: var(--text3); text-transform: uppercase; margin-bottom: 6px; letter-spacing: 0.5px; }
.api-no-param { font-size: 11px; color: var(--text3); padding: 6px 0; font-style: italic; }

.code-box-api { background: var(--surface); border: 1px solid var(--border); border-radius: var(--r8); padding: 10px; max-height: 160px; overflow-y: auto; }
.code-box-api pre { margin: 0; font-family: monospace; font-size: 11px; color: var(--text2); }

.animate-fade { animation: fadeIn .3s ease-in-out; }
@keyframes fadeIn { from { opacity: 0; transform: translateY(4px); } to { opacity: 1; transform: translateY(0); } }
</style>
