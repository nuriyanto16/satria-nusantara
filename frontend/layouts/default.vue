<template>
  <div class="layout" :class="{ 'sidebar-open': isSidebarOpen }" @click="closeSidebarOnOverlay">
    <!-- SIDEBAR -->
    <div class="sidebar">
      <div class="sb-header">
        <div class="sb-logo">
          <div class="sb-icon">SN</div>
          <div>
            <div class="sb-title">Satria Nusantara</div>
            <div class="sb-sub">LSP Web Admin</div>
          </div>
        </div>
        <div class="sb-badge">Pusat · Seluruh Indonesia</div>
      </div>

      <div class="sb-section">Utama</div>
      <NuxtLink to="/" class="sb-item"><i class="ti ti-layout-dashboard"></i> Dashboard</NuxtLink>

      <div class="sb-section">Organisasi</div>
      <NuxtLink to="/cabang"   class="sb-item"><i class="ti ti-map-pin"></i> Kelola Cabang</NuxtLink>
      <NuxtLink to="/anggota"  class="sb-item"><i class="ti ti-users"></i> Kelola Anggota</NuxtLink>
      <NuxtLink to="/jadwal"   class="sb-item"><i class="ti ti-calendar"></i> Jadwal Latihan</NuxtLink>
      <NuxtLink to="/absensi"  class="sb-item"><i class="ti ti-qrcode"></i> Absensi & Kehadiran</NuxtLink>
      <NuxtLink to="/latgab"   class="sb-item"><i class="ti ti-award"></i> Latgab / EKT / Pelatnas</NuxtLink>
      <NuxtLink to="/kebugaran" class="sb-item"><i class="ti ti-heartbeat"></i> Tes Kebugaran</NuxtLink>
      <NuxtLink to="/nafas"     class="sb-item"><i class="ti ti-wind"></i> Monitoring Olah Nafas</NuxtLink>

      <div class="sb-section">Keuangan</div>
      <NuxtLink to="/iuran"   class="sb-item"><i class="ti ti-wallet"></i> BLBA</NuxtLink>
      <NuxtLink to="/laporan" class="sb-item"><i class="ti ti-chart-bar"></i> Laporan</NuxtLink>

      <div class="sb-section">Konten</div>
      <NuxtLink to="/konten" class="sb-item"><i class="ti ti-news"></i> Konten & Artikel</NuxtLink>

      <div class="sb-section">Sistem</div>
      <NuxtLink to="/user"       class="sb-item"><i class="ti ti-user-cog"></i> Manajemen User</NuxtLink>
      <NuxtLink to="/pengaturan" class="sb-item"><i class="ti ti-settings"></i> Pengaturan</NuxtLink>
      <NuxtLink to="/dokumentasi" class="sb-item"><i class="ti ti-file-text"></i> Dokumentasi</NuxtLink>
    </div>
    <!-- TOPBAR + MAIN -->
    <div style="flex:1;overflow:hidden;position:relative;display:flex;flex-direction:column">
      <div class="topbar">
        <div class="tb-left">
          <button class="tb-menu-toggle-btn" @click.stop="toggleSidebar">
            <i class="ti ti-menu-2"></i>
          </button>
          <div class="tb-breadcrumb">
            <span class="tb-bc-home" title="Beranda" @click="navigateTo('/')"><i class="ti ti-home" style="font-size:15px"></i></span>
            <span class="tb-bc-sep"><i class="ti ti-chevron-right" style="font-size:11px"></i></span>
            <span class="tb-bc-cur">{{ currentTitle }}</span>
          </div>

          <!-- Scope indicator pill -->
          <div class="tb-scope-badge">
            <span class="scope-dot"></span>
            <span>Pusat · Indonesia</span>
          </div>

          <!-- Premium search bar -->
          <div class="tb-search-bar" title="Cari peserta / anggota (Cmd+K)" @click.stop>
            <i class="ti ti-search"></i>
            <input
              ref="headerInputRef"
              v-model="headerSearchQuery"
              type="text"
              placeholder="Cari peserta, anggota, NIK..."
              @input="onHeaderSearchInput"
              @focus="onHeaderSearchFocus"
              @keydown.esc="closeSearchDropdown"
              @keydown.enter="goToAnggotaPage"
            />
            <span v-if="!headerSearchQuery" class="search-kbd">⌘K</span>
            <button v-else class="search-clear-btn" @click.stop="clearHeaderSearch">
              <i class="ti ti-x"></i>
            </button>

            <!-- Live Search Results Dropdown -->
            <div v-if="showSearchModal" class="search-dropdown" @click.stop>
              <div class="search-dropdown-header">
                <span>Hasil Pencarian Peserta ({{ searchResults.length }})</span>
                <span v-if="searchLoading" class="search-spinner"><i class="ti ti-loader-2 spin"></i> Memuat...</span>
              </div>
              
              <div class="search-dropdown-body">
                <div v-if="searchLoading && searchResults.length === 0" class="search-loading-state">
                  <i class="ti ti-loader-2 spin"></i> Mencari data peserta...
                </div>
                
                <div v-else-if="!searchLoading && searchResults.length === 0" class="search-empty-state">
                  <i class="ti ti-user-x"></i> Tidak ada peserta ditemukan untuk "{{ headerSearchQuery }}"
                </div>
                
                <div
                  v-for="item in searchResults"
                  :key="item.id"
                  class="search-result-item"
                  @click="selectSearchResult(item)"
                >
                  <div class="sri-avatar-wrap">
                    <img v-if="item.foto_url" :src="item.foto_url" class="sri-avatar-img" />
                    <div v-else class="sri-avatar-fallback">
                      {{ getSearchInitials(item.nama_lengkap) }}
                    </div>
                  </div>
                  <div class="sri-info">
                    <div class="sri-name-row">
                      <span class="sri-name">{{ item.nama_lengkap }}</span>
                      <span :class="['sri-status-badge', item.status]">{{ item.status || 'aktif' }}</span>
                    </div>
                    <div class="sri-meta">
                      <span><i class="ti ti-id"></i> {{ item.nomor_anggota || 'PENDING' }}</span>
                      <span><i class="ti ti-building"></i> {{ item.cabang_nama || 'Pusat' }}</span>
                      <span><i class="ti ti-map-pin"></i> {{ item.unit_nama || 'Unit' }}</span>
                      <span class="sri-tingkat">{{ item.tingkatan || 'Dasar' }}</span>
                    </div>
                  </div>
                  <i class="ti ti-chevron-right sri-arrow"></i>
                </div>
              </div>

              <div class="search-dropdown-footer">
                <span class="sdf-tip">Tekan <kbd>ESC</kbd> untuk menutup · <kbd>Enter</kbd> cari semua</span>
                <button class="sdf-view-all" @click="goToAnggotaPage">Lihat Semua di Kelola Anggota →</button>
              </div>
            </div>
          </div>
        </div>
        <div class="tb-right">
          <div class="tb-notif-btn" @click.stop="toggleTheme" :title="isDarkMode ? 'Beralih ke Mode Terang' : 'Beralih ke Mode Gelap'">
            <i :class="isDarkMode ? 'ti ti-sun' : 'ti ti-moon'"></i>
          </div>
          <div class="tb-notif-btn" @click.stop="showNotificationsPopover = !showNotificationsPopover" title="Notifikasi terbaru">
            <i class="ti ti-bell"></i>
            <span v-if="unreadCount > 0" class="tb-notif-dot"></span>
 
            <!-- Popover Dropdown -->
            <div v-if="showNotificationsPopover" class="notif-dropdown" @click.stop>
              <div class="notif-header">
                <span>Notifikasi ({{ unreadCount }} Baru)</span>
                <button class="notif-clear-btn" @click="markAllAsRead">Tandai semua dibaca</button>
              </div>
              <div class="notif-body">
                <div v-if="notifications.length === 0" class="notif-empty">Tidak ada notifikasi baru</div>
                <div v-for="n in notifications" :key="n.id" :class="['notif-item', { unread: !n.isRead }]" @click="openNotifDetail(n)">
                  <div class="notif-item-header">
                    <span class="notif-type-badge">{{ n.type }}</span>
                    <span class="notif-time">{{ n.time }}</span>
                  </div>
                  <div class="notif-item-title">{{ n.title }}</div>
                  <div class="notif-item-preview">{{ n.message }}</div>
                </div>
              </div>
            </div>
          </div>
          <div class="tb-user-info" @click="handleLogout" title="Klik untuk keluar dari sistem">
            <div class="tb-avatar">{{ userInitials }}</div>
            <div class="user-meta-details">
              <div class="tb-uname">{{ authStore.user?.nama_lengkap || 'Sri Astuti' }}</div>
              <div class="tb-urole">{{ authStore.user?.role_name || 'Admin Pusat' }}</div>
            </div>
            <i class="ti ti-logout" style="font-size:15px;color:var(--merah);margin-left:4px"></i>
          </div>
        </div>
      </div>
 
      <div class="page-container">
        <slot />
      </div>
    </div>

    <!-- NOTIFIKASI DETAIL MODAL -->
    <div v-if="selectedNotification" class="modal-overlay" @click.self="selectedNotification = null">
      <div class="modal-card">
        <div class="modal-header">
          <h2 class="modal-title">Detail Notifikasi</h2>
          <button class="modal-close" @click="selectedNotification = null"><i class="ti ti-x"></i></button>
        </div>
        <div class="modal-detail-content">
          <div class="nd-row"><strong>Pengirim:</strong> LSP Satria Nusantara Pusat</div>
          <div class="nd-row"><strong>Kategori:</strong> <span class="notif-type-badge">{{ selectedNotification.type }}</span></div>
          <div class="nd-row"><strong>Tanggal:</strong> {{ selectedNotification.time }}</div>
          <div class="notif-detail-body">
            <h4 class="nd-body-title">{{ selectedNotification.title }}</h4>
            <p class="nd-body-text">{{ selectedNotification.message }}</p>
          </div>
        </div>
        <div class="modal-footer" style="margin-top: 16px;">
          <button class="btn btn-outline" @click="selectedNotification = null">Tutup</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useAuthStore } from '~/stores/auth'

const route = useRoute()
const authStore = useAuthStore()

const isSidebarOpen = ref(false)
const isDarkMode = ref(false)

const applyTheme = (isDark) => {
  if (isDark) {
    document.documentElement.classList.add('dark')
    document.body.classList.add('dark')
    const layout = document.querySelector('.layout')
    if (layout) layout.classList.add('dark')
    localStorage.setItem('theme', 'dark')
  } else {
    document.documentElement.classList.remove('dark')
    document.body.classList.remove('dark')
    const layout = document.querySelector('.layout')
    if (layout) layout.classList.remove('dark')
    localStorage.setItem('theme', 'light')
  }
}

const toggleTheme = () => {
  isDarkMode.value = !isDarkMode.value
  applyTheme(isDarkMode.value)
}

const api = useApi()
const headerSearchQuery = ref('')
const showSearchModal = ref(false)
const searchLoading = ref(false)
const searchResults = ref([])
const headerInputRef = ref(null)
let searchTimer = null

const getSearchInitials = (name) => {
  if (!name) return 'SN'
  const parts = name.trim().split(' ')
  if (parts.length > 1) return (parts[0][0] + parts[1][0]).toUpperCase()
  return name.substring(0, 2).toUpperCase()
}

const onHeaderSearchFocus = () => {
  if (headerSearchQuery.value.trim().length > 0) {
    showSearchModal.value = true
  }
}

const onHeaderSearchInput = () => {
  clearTimeout(searchTimer)
  const query = headerSearchQuery.value.trim()
  if (query.length === 0) {
    searchResults.value = []
    showSearchModal.value = false
    return
  }

  showSearchModal.value = true
  searchLoading.value = true

  searchTimer = setTimeout(async () => {
    try {
      const res = await api.get(`/organization/anggota?limit=8&search=${encodeURIComponent(query)}`)
      searchResults.value = res.data || []
    } catch (e) {
      console.error('Header search error:', e)
      searchResults.value = []
    } finally {
      searchLoading.value = false
    }
  }, 250)
}

const clearHeaderSearch = () => {
  headerSearchQuery.value = ''
  searchResults.value = []
  showSearchModal.value = false
}

const closeSearchDropdown = () => {
  showSearchModal.value = false
}

const selectSearchResult = (item) => {
  showSearchModal.value = false
  navigateTo(`/anggota?search=${encodeURIComponent(item.nama_lengkap)}`)
}

const goToAnggotaPage = () => {
  const q = headerSearchQuery.value.trim()
  showSearchModal.value = false
  if (q) {
    navigateTo(`/anggota?search=${encodeURIComponent(q)}`)
  } else {
    navigateTo('/anggota')
  }
}

const showNotificationsPopover = ref(false)
const selectedNotification = ref(null)

const notifications = ref([
  { id: 1, type: 'Pendaftaran', time: '10 menit yang lalu', title: 'Pendaftaran Anggota Baru', message: 'Ahmad Santoso mendaftar ke Cabang Kota Yogyakarta (Unit Malioboro). Butuh verifikasi berkas dan KTP.', isRead: false },
  { id: 2, type: 'Jadwal', time: '1 jam yang lalu', title: 'Jadwal Sesi Baru', message: 'Pelatih Budi Susanto membuat sesi Latihan Rutin baru untuk Unit Rungkut, Cabang Surabaya.', isRead: false },
  { id: 3, type: 'Keuangan', time: 'Kemarin', title: 'Pembayaran BLBA Masuk', message: 'Cabang Kota Bandung menyetor rekap BLBA bulan Juli senilai Rp 12.500.000 via Midtrans.', isRead: true },
])

const unreadCount = computed(() => notifications.value.filter(n => !n.isRead).length)

const toggleSidebar = () => {
  isSidebarOpen.value = !isSidebarOpen.value
}

const closeSidebarOnOverlay = () => {
  if (isSidebarOpen.value) {
    isSidebarOpen.value = false
  }
}

const markAllAsRead = () => {
  notifications.value.forEach(n => n.isRead = true)
}

const openNotifDetail = (n) => {
  n.isRead = true
  selectedNotification.value = n
  showNotificationsPopover.value = false
}

onMounted(() => {
  const savedTheme = localStorage.getItem('theme')
  if (savedTheme === 'dark' || (!savedTheme && window.matchMedia('(prefers-color-scheme: dark)').matches)) {
    isDarkMode.value = true
    applyTheme(true)
  } else {
    isDarkMode.value = false
    applyTheme(false)
  }

  authStore.rehydrate()
  if (!authStore.isAuthenticated) {
    navigateTo('/login')
  }

  // Close notifications & search popovers on click outside
  document.addEventListener('click', () => {
    showNotificationsPopover.value = false
    showSearchModal.value = false
  })

  // Keyboard shortcut Cmd+K or Ctrl+K
  document.addEventListener('keydown', (e) => {
    if ((e.metaKey || e.ctrlKey) && e.key.toLowerCase() === 'k') {
      e.preventDefault()
      headerInputRef.value?.focus()
      if (headerSearchQuery.value.trim().length > 0) {
        showSearchModal.value = true
      }
    } else if (e.key === 'Escape') {
      showSearchModal.value = false
    }
  })
})

const userInitials = computed(() => {
  if (!authStore.user?.nama_lengkap) return 'SA'
  const parts = authStore.user.nama_lengkap.split(' ')
  if (parts.length > 1) {
    return (parts[0][0] + parts[1][0]).toUpperCase()
  }
  return parts[0].substring(0, 2).toUpperCase()
})

const handleLogout = () => {
  if (confirm('Apakah Anda yakin ingin keluar?')) {
    authStore.logout()
  }
}

const titleMap = {
  '/': 'Dashboard',
  '/cabang': 'Kelola Cabang',
  '/anggota': 'Kelola Anggota',
  '/jadwal': 'Jadwal Latihan',
  '/absensi': 'Absensi & Kehadiran',
  '/latgab': 'Latgab / EKT / Pelatnas',
  '/kebugaran': 'Tes Kebugaran',
  '/nafas': 'Monitoring Olah Nafas',
  '/iuran': 'BLBA',
  '/laporan': 'Laporan',
  '/konten': 'Konten & Artikel',
  '/user': 'Manajemen User',
  '/pengaturan': 'Pengaturan',
  '/dokumentasi': 'Dokumentasi Sistem',
}

const currentTitle = computed(() => titleMap[route.path] || 'Dashboard')

// Close sidebar on navigation change
watch(() => route.path, () => {
  isSidebarOpen.value = false
})
</script>

<style scoped>
.layout {
  display: flex;
  height: 100vh;
  overflow: hidden;
}

.topbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 64px;
  background: rgba(255,255,255,0.85);
  backdrop-filter: blur(16px);
  border-bottom: 1px solid var(--border);
  padding: 0 24px;
  position: relative;
  z-index: 100;
}

.btn { display: inline-flex; align-items: center; gap: 6px; padding: 8px 14px; border-radius: var(--r8); font-size: 12px; font-weight: 600; cursor: pointer; border: none; }
.btn-outline { background: #fff; border: 1px solid var(--border); color: var(--text2); }

/* Popover dropdown style */
.notif-dropdown {
  position: absolute;
  top: 42px; right: 0;
  width: 320px;
  background: var(--card);
  border: 1px solid var(--border);
  border-radius: var(--r12);
  box-shadow: var(--shadow-lg);
  z-index: 1002;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}
.notif-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 14px;
  border-bottom: 1px solid var(--border);
  font-size: 12px;
  font-weight: 700;
}
.notif-clear-btn {
  background: none;
  border: none;
  color: var(--hijau);
  font-weight: 600;
  cursor: pointer;
}
.notif-body {
  max-height: 280px;
  overflow-y: auto;
}
.notif-empty {
  padding: 20px;
  text-align: center;
  font-size: 12px;
  color: var(--text3);
}
.notif-item {
  padding: 12px 14px;
  border-bottom: 1px solid var(--border);
  cursor: pointer;
  transition: background .15s;
  display: flex;
  flex-direction: column;
  gap: 4px;
}
.notif-item:hover {
  background: var(--surface);
}
.notif-item.unread {
  background: rgba(26, 92, 42, 0.04);
}
.notif-item.unread:hover {
  background: rgba(26, 92, 42, 0.07);
}
.notif-item-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 10px;
}
.notif-type-badge {
  background: var(--hijau3);
  color: var(--hijau);
  padding: 2px 6px;
  border-radius: 4px;
  font-weight: 700;
}
.notif-time {
  color: var(--text3);
}
.notif-item-title {
  font-size: 12px;
  font-weight: 700;
  color: var(--text1);
}
.notif-item-preview {
  font-size: 11px;
  color: var(--text2);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

/* Modal detail content */
.modal-overlay { position: fixed; top: 0; left: 0; right: 0; bottom: 0; background: rgba(0,0,0,0.5); z-index: 1000; display: flex; align-items: center; justify-content: center; }
.modal-card { background: var(--card); border-radius: var(--r12); width: 100%; max-width: 440px; padding: 20px; box-shadow: var(--shadow-lg); }
.modal-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px; }
.modal-title { font-size: 16px; font-weight: 700; }
.modal-close { background: none; border: none; font-size: 16px; cursor: pointer; color: var(--text3); }
.modal-detail-content { display: flex; flex-direction: column; gap: 10px; font-size: 13px; }
.nd-row { border-bottom: 1px solid var(--border); padding-bottom: 8px; }
.notif-detail-body { margin-top: 12px; padding: 12px; background: var(--surface); border-radius: var(--r8); border: 1px solid var(--border) }
.nd-body-title { margin-bottom: 6px; font-weight: 700; font-size: 14px; }
.nd-body-text { font-size: 12px; line-height: 1.5; color: var(--text2); }

/* Header Search Bar & Live Dropdown */
.tb-search-bar {
  position: relative;
  display: flex;
  align-items: center;
  gap: 8px;
  background: var(--surface);
  border: 1.5px solid var(--border);
  border-radius: var(--r8);
  padding: 6px 12px;
  width: 240px;
  transition: all 0.2s;
}
.tb-search-bar:focus-within {
  border-color: var(--hijau);
  background: #fff;
  box-shadow: 0 0 0 3px var(--hijauSoft);
  width: 300px;
}
.tb-search-bar i { color: var(--text3); font-size: 14px; }
.tb-search-bar input {
  border: none;
  outline: none;
  background: none;
  font-size: 12px;
  flex: 1;
  color: var(--text1);
  width: 100%;
}
.search-kbd {
  font-size: 10px;
  background: var(--border2);
  color: var(--text3);
  padding: 2px 5px;
  border-radius: 4px;
  font-weight: 600;
}
.search-clear-btn {
  background: none;
  border: none;
  cursor: pointer;
  padding: 0;
  color: var(--text3);
  display: flex;
  align-items: center;
}
.search-clear-btn:hover { color: var(--merah); }

/* Dropdown Popup */
.search-dropdown {
  position: absolute;
  top: calc(100% + 8px);
  left: 0;
  width: 420px;
  background: var(--card);
  border: 1px solid var(--border);
  border-radius: var(--r12);
  box-shadow: var(--shadow-lg);
  z-index: 1005;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  max-height: 480px;
}

.search-dropdown-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 14px;
  background: var(--surface);
  border-bottom: 1px solid var(--border);
  font-size: 11px;
  font-weight: 700;
  color: var(--text2);
}
.search-spinner { color: var(--hijau); display: flex; align-items: center; gap: 4px; }
.spin { animation: spin 1s linear infinite; }
@keyframes spin { 100% { transform: rotate(360deg); } }

.search-dropdown-body {
  overflow-y: auto;
  max-height: 340px;
}
.search-loading-state, .search-empty-state {
  padding: 24px;
  text-align: center;
  font-size: 12px;
  color: var(--text3);
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}
.search-empty-state i { font-size: 28px; color: var(--text3); }

.search-result-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 14px;
  border-bottom: 1px solid var(--border);
  cursor: pointer;
  transition: background 0.15s;
}
.search-result-item:hover { background: var(--surface); }
.search-result-item:last-child { border-bottom: none; }

.sri-avatar-wrap { width: 36px; height: 36px; border-radius: 50%; overflow: hidden; flex-shrink: 0; }
.sri-avatar-img { width: 100%; height: 100%; object-fit: cover; }
.sri-avatar-fallback {
  width: 100%; height: 100%;
  background: var(--hijauSoft);
  color: var(--hijau);
  font-size: 12px;
  font-weight: 700;
  display: flex;
  align-items: center;
  justify-content: center;
}

.sri-info { flex: 1; min-width: 0; }
.sri-name-row { display: flex; align-items: center; justify-content: space-between; gap: 6px; }
.sri-name { font-size: 13px; font-weight: 700; color: var(--text1); white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.sri-status-badge {
  font-size: 9px;
  font-weight: 700;
  padding: 2px 6px;
  border-radius: 8px;
  text-transform: uppercase;
}
.sri-status-badge.aktif { background: var(--hijauSoft); color: var(--hijau); }
.sri-status-badge.pending { background: #fff8e0; color: #a07000; }
.sri-status-badge.nonaktif { background: #fdecea; color: var(--merah); }

.sri-meta { display: flex; align-items: center; gap: 10px; font-size: 10px; color: var(--text3); margin-top: 3px; flex-wrap: wrap; }
.sri-meta span { display: flex; align-items: center; gap: 3px; }
.sri-tingkat { background: var(--border2); color: var(--text2); padding: 1px 5px; border-radius: 4px; font-weight: 600; }
.sri-arrow { color: var(--text3); font-size: 14px; }

.search-dropdown-footer {
  padding: 10px 14px;
  background: var(--surface);
  border-top: 1px solid var(--border);
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 10px;
}
.sdf-tip { color: var(--text3); }
.sdf-tip kbd { background: #fff; border: 1px solid var(--border); padding: 1px 4px; border-radius: 3px; }
.sdf-view-all { background: none; border: none; color: var(--hijau); font-weight: 700; cursor: pointer; font-size: 10px; }
.sdf-view-all:hover { text-decoration: underline; }

</style>
