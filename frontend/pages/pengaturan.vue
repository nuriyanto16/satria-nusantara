<template>
  <div id="page-pengaturan">
    <div class="page-header">
      <div><h1 class="page-title">Pengaturan Sistem</h1><div class="page-subtitle">Konfigurasi parameter sistem — hanya dapat diakses Admin Pusat</div></div>
      <button class="btn btn-primary" @click="simpanPengaturan">
        <i :class="saving ? 'ti ti-loader-2 spin' : 'ti ti-device-floppy'"></i> {{ saving ? 'Menyimpan...' : 'Simpan Perubahan' }}
      </button>
    </div>

    <!-- Banner Sukses -->
    <transition name="fade">
      <div v-if="saveSuccess" class="save-success-banner">
        <i class="ti ti-circle-check"></i> Pengaturan dan tema sistem berhasil disimpan!
      </div>
    </transition>

    <div class="setting-layout">
      <!-- Sidebar kategori -->
      <div class="setting-sidebar">
        <div v-for="k in kategoriList" :key="k.value" :class="['setting-menu',{active:activeKat===k.value}]" @click="activeKat=k.value">
          <i :class="'ti '+k.icon"></i> {{ k.label }}
        </div>
      </div>
      <!-- Content -->
      <div class="setting-content">
        <!-- Tema Selector -->
        <div v-if="activeKat === 'tema'" class="setting-section">
          <label class="setting-label">Tema Tampilan Web Admin</label>
          <div class="setting-desc">Pilih mode tampilan sesuai dengan kenyamanan mata Anda</div>
          <div class="theme-options-grid">
            <div :class="['theme-option-card', { active: currentTheme === 'light' }]" @click="setTheme('light')">
              <div class="toc-icon light"><i class="ti ti-sun"></i></div>
              <div>
                <div class="toc-title">Mode Terang (Light Mode)</div>
                <div class="toc-desc">Tampilan standar bersih dengan latar belakang terang</div>
              </div>
            </div>
            <div :class="['theme-option-card', { active: currentTheme === 'dark' }]" @click="setTheme('dark')">
              <div class="toc-icon dark"><i class="ti ti-moon"></i></div>
              <div>
                <div class="toc-title">Mode Gelap (Dark Mode)</div>
                <div class="toc-desc">Nuansa hijau & deep dark navy yang nyaman di malam hari</div>
              </div>
            </div>
          </div>
        </div>

        <!-- Dynamic Settings -->
        <template v-else>
          <div class="setting-section" v-for="s in currentSettings" :key="s.kunci">
            <label class="setting-label">{{ s.label }}</label>
            <div class="setting-desc">{{ s.desc }}</div>
            <input :type="s.type||'text'" class="setting-input" v-model="s.nilai" />
          </div>
        </template>
      </div>
    </div>
  </div>
</template>
<script setup>
definePageMeta({ title: 'Pengaturan' })
const activeKat = ref('tema')
const currentTheme = ref('light')
const saving = ref(false)
const saveSuccess = ref(false)

const setTheme = (mode) => {
  currentTheme.value = mode
  localStorage.setItem('theme', mode)
  if (mode === 'dark') {
    document.documentElement.classList.add('dark')
    document.body.classList.add('dark')
  } else {
    document.documentElement.classList.remove('dark')
    document.body.classList.remove('dark')
  }
}

const simpanPengaturan = () => {
  saving.value = true
  setTheme(currentTheme.value)
  setTimeout(() => {
    saving.value = false
    saveSuccess.value = true
    setTimeout(() => {
      saveSuccess.value = false
    }, 3500)
  }, 300)
}

onMounted(() => {
  const saved = localStorage.getItem('theme')
  if (saved === 'dark' || (!saved && window.matchMedia('(prefers-color-scheme: dark)').matches)) {
    currentTheme.value = 'dark'
  } else {
    currentTheme.value = 'light'
  }
})

const kategoriList = [
  {label:'Tema Tampilan',value:'tema',icon:'ti-palette'},
  {label:'BLBA',value:'blba',icon:'ti-wallet'},
  {label:'Transport Pelatih',value:'transport',icon:'ti-car'},
  {label:'Syarat EKT',value:'ekt',icon:'ti-award'},
  {label:'Tes Kebugaran',value:'kebugaran',icon:'ti-heartbeat'},
  {label:'Profil Organisasi',value:'profil',icon:'ti-building'},
  {label:'Notifikasi',value:'notifikasi',icon:'ti-bell'},
  {label:'Aplikasi Mobile',value:'mobile',icon:'ti-device-mobile'},
]
const allSettings = reactive({
  blba:[
    {kunci:'blba_nominal',label:'Nominal BLBA per Bulan',desc:'Besaran iuran bulanan wajib anggota aktif (Rupiah)',nilai:'50000',type:'number'},
    {kunci:'blba_grace_period',label:'Grace Period Pembayaran (hari)',desc:'Toleransi hari pembayaran setelah tanggal jatuh tempo',nilai:'10',type:'number'},
  ],
  transport:[
    {kunci:'transport_jarak_dekat',label:'Nominal Transport Jarak Dekat (Rp/sesi)',desc:'Untuk pelatih yang tinggal di wilayah yang sama',nilai:'30000',type:'number'},
    {kunci:'transport_jarak_sedang',label:'Nominal Transport Jarak Sedang (Rp/sesi)',desc:'Untuk pelatih yang tinggal di kecamatan berbeda',nilai:'50000',type:'number'},
    {kunci:'transport_jarak_jauh',label:'Nominal Transport Jarak Jauh (Rp/sesi)',desc:'Untuk pelatih yang tinggal di kota/kabupaten berbeda',nilai:'75000',type:'number'},
  ],
  ekt:[
    {kunci:'ekt_syarat_hadir_pra_dasar',label:'Min. Kehadiran — Pra Dasar',desc:'Jumlah minimum sesi hadir untuk EKT dari Pra Dasar',nilai:'10',type:'number'},
    {kunci:'ekt_syarat_hadir_dasar',label:'Min. Kehadiran — Dasar',desc:'Jumlah minimum sesi hadir untuk EKT dari Dasar',nilai:'8',type:'number'},
    {kunci:'ekt_syarat_hadir_gabungan',label:'Min. Kehadiran — Gabungan',desc:'Jumlah minimum sesi hadir untuk EKT dari Gabungan',nilai:'12',type:'number'},
  ],
  kebugaran:[
    {kunci:'tes_kebugaran_periode_per_tahun',label:'Jumlah Periode Tes per Tahun',desc:'Berapa kali tes kebugaran diadakan setiap tahun',nilai:'3',type:'number'},
  ],
  profil:[
    {kunci:'nama_organisasi',label:'Nama Resmi Organisasi',desc:'Nama lengkap organisasi yang tampil di sistem dan dokumen',nilai:'Satria Nusantara'},
    {kunci:'kota_pusat',label:'Kota Pusat',desc:'Kota lokasi kantor pusat',nilai:'Kota Yogyakarta'},
  ],
  notifikasi:[{kunci:'qr_berlaku_jam',label:'Durasi Berlaku QR Absensi (jam)',desc:'Setelah durasi ini, QR yang digenerate pelatih akan expired',nilai:'4',type:'number'}],
  mobile:[{kunci:'qr_berlaku_jam',label:'Durasi Berlaku QR Absensi (jam)',desc:'Masa berlaku QR Code per sesi yang digenerate pelatih',nilai:'4',type:'number'}],
})
const currentSettings = computed(() => allSettings[activeKat.value] || [])
</script>
<style scoped>
.page-header{display:flex;justify-content:space-between;align-items:flex-start;margin-bottom:24px}
.page-title{font-size:20px;font-weight:800;margin:0}.page-subtitle{font-size:12px;color:var(--text3);margin-top:2px}
.btn{display:inline-flex;align-items:center;gap:6px;padding:8px 14px;border-radius:var(--r8);font-size:12px;font-weight:600;cursor:pointer;border:none}
.btn-primary{background:var(--hijau);color:#fff}
.setting-layout{display:grid;grid-template-columns:200px 1fr;gap:16px}
.setting-sidebar{background:var(--card);border:1px solid var(--border);border-radius:var(--r12);padding:8px;display:flex;flex-direction:column;gap:2px;height:fit-content}
.setting-menu{display:flex;align-items:center;gap:8px;padding:9px 12px;border-radius:var(--r8);font-size:12px;font-weight:600;cursor:pointer;color:var(--text2);transition:all .15s}
.setting-menu.active{background:var(--hijau3);color:var(--hijau)}.setting-menu:hover{background:var(--surface)}
.setting-content{background:var(--card);border:1px solid var(--border);border-radius:var(--r12);padding:20px;display:flex;flex-direction:column;gap:16px}
.setting-section{border-bottom:1px solid var(--border);padding-bottom:16px}
.setting-section:last-child{border-bottom:none;padding-bottom:0}
.setting-label{font-size:13px;font-weight:700;display:block;margin-bottom:4px}
.setting-desc{font-size:11px;color:var(--text3);margin-bottom:8px}
.setting-input{width:100%;max-width:300px;padding:8px 10px;border:1px solid var(--border2);border-radius:var(--r8);font-size:13px;outline:none}
.setting-input:focus{border-color:var(--hijau)}

.theme-options-grid { display: grid; grid-template-columns: repeat(2, 1fr); gap: 14px; margin-top: 10px; }
@media (max-width: 600px) { .theme-options-grid { grid-template-columns: 1fr; } }
.theme-option-card { display: flex; align-items: flex-start; gap: 12px; padding: 14px; border: 1.5px solid var(--border); border-radius: var(--r12); background: var(--surface); cursor: pointer; transition: all .15s; }
.theme-option-card:hover { border-color: var(--hijau); }
.theme-option-card.active { border-color: var(--hijau); background: var(--hijau3); }
.toc-icon { width: 36px; height: 36px; border-radius: 50%; display: flex; align-items: center; justify-content: center; font-size: 18px; flex-shrink: 0; }
.toc-icon.light { background: #fff8e0; color: var(--kuning); }
.toc-icon.dark { background: #152219; color: #5bb8d4; }
.toc-title { font-size: 13px; font-weight: 700; color: var(--text1); margin-bottom: 2px; }
.toc-desc { font-size: 11px; color: var(--text3); line-height: 1.3; }

.save-success-banner {
  background: var(--hijau3);
  color: var(--hijau);
  border: 1px solid var(--border);
  border-radius: var(--r8);
  padding: 10px 14px;
  font-size: 13px;
  font-weight: 600;
  margin-bottom: 16px;
  display: flex;
  align-items: center;
  gap: 8px;
}
.spin { animation: spin 1s linear infinite; }
@keyframes spin { 100% { transform: rotate(360deg); } }
.fade-enter-active, .fade-leave-active { transition: opacity .3s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
</style>
