<template>
  <div id="page-user">
    <!-- Header -->
    <div class="page-header">
      <div>
        <h1 class="page-title">Manajemen User</h1>
        <div class="page-subtitle">Kelola akun pengurus dan akses kontrol sistem</div>
      </div>
      <div class="header-actions">
        <div class="search-box-header">
          <i class="ti ti-search"></i>
          <input v-model="search" type="text" placeholder="Cari nama atau email..." />
        </div>
        <button class="btn btn-primary" @click="openAddModal"><i class="ti ti-user-plus"></i> Tambah User</button>
      </div>
    </div>

    <!-- Tabs -->
    <div class="tab-bar">
      <button v-for="t in tabs" :key="t.value" :class="['tab-btn', { active: activeTab === t.value }]" @click="activeTab = t.value">
        {{ t.label }} <span class="tab-count">{{ getTabCount(t.value) }}</span>
      </button>
    </div>

    <!-- User Grid -->
    <div v-if="filteredUsers.length === 0" class="empty-state">
      <i class="ti ti-user-x"></i>
      <p>Tidak ada user pengurus ditemukan.</p>
    </div>

    <div v-else class="user-grid">
      <div v-for="u in filteredUsers" :key="u.id" class="user-card" :class="{ disabled: u.status === 'nonactive' }">
        <div class="user-card-top">
          <div class="user-avatar" :style="{ background: u.bg }">{{ u.av }}</div>
          <div class="user-info">
            <div class="user-name">{{ u.name }}</div>
            <div class="user-role">{{ u.role }}</div>
            <div class="user-scope"><i class="ti ti-map-pin"></i> {{ u.scope }}</div>
          </div>
          <span :class="['akses-badge', u.akses]">{{ u.role }}</span>
        </div>
        <div class="user-detail">
          <div class="ud-item"><i class="ti ti-briefcase"></i> {{ u.jabatan }}</div>
          <div class="ud-item"><i class="ti ti-mail"></i> {{ u.email }}</div>
          <div class="ud-item"><i class="ti ti-phone"></i> {{ u.hp }}</div>
        </div>
        <div class="user-akses-list">
          <div v-for="a in aksesMap[u.akses]" :key="a[0]" :class="['akses-item', a[3]]">
            <i :class="'ti '+a[1]"></i> {{ a[0] }} <span>{{ a[2] }}</span>
          </div>
        </div>
        <div class="user-card-footer">
          <button class="btn-sm-outline" @click="openEditModal(u)"><i class="ti ti-edit"></i> Edit</button>
          <button v-if="u.status === 'active'" class="btn-sm-outline danger" @click="toggleUserStatus(u)">
            <i class="ti ti-lock"></i> Nonaktifkan
          </button>
          <button v-else class="btn-sm-outline green" @click="toggleUserStatus(u)">
            <i class="ti ti-lock-open"></i> Aktifkan
          </button>
        </div>
      </div>
    </div>

    <!-- MODAL: ADD / EDIT USER -->
    <div v-if="showModal" class="modal-overlay" @click.self="showModal = false">
      <div class="modal-card">
        <div class="modal-header">
          <div class="modal-title">{{ editingId ? 'Edit Akun User' : 'Tambah Akun User Baru' }}</div>
          <button class="modal-close" @click="showModal = false"><i class="ti ti-x"></i></button>
        </div>
        <form @submit.prevent="saveUser" class="modal-form">
          <div class="form-group">
            <label class="form-label">Nama Lengkap</label>
            <input v-model="form.name" type="text" class="form-input" placeholder="Masukkan nama lengkap..." required />
          </div>
          <div class="form-row">
            <div class="form-group half">
              <label class="form-label">Role Akses</label>
              <select v-model="form.akses" class="form-select" required>
                <option value="penuh">Admin Pusat</option>
                <option value="cabang">Pengurus Cabang</option>
                <option value="unit">PIC Unit</option>
                <option value="pelatih">Pelatih</option>
              </select>
            </div>
            <div class="form-group half">
              <label class="form-label">Jabatan Struktural</label>
              <input v-model="form.jabatan" type="text" class="form-input" placeholder="Sekretaris, PIC Unit, dll" required />
            </div>
          </div>
          <div class="form-group">
            <label class="form-label">Email Instansi / Pribadi</label>
            <input v-model="form.email" type="email" class="form-input" placeholder="nama@satria-nusantara.org" required />
          </div>
          <div class="form-group">
            <label class="form-label">No. HP / WhatsApp</label>
            <input v-model="form.hp" type="text" class="form-input" placeholder="Contoh: 081234567890" required />
          </div>
          <div class="form-group">
            <label class="form-label">Wilayah Kerja (Scope)</label>
            <input v-model="form.scope" type="text" class="form-input" placeholder="Contoh: Cabang · Kota Yogyakarta" required />
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-outline" style="width:auto" @click="showModal = false">Batal</button>
            <button type="submit" class="btn btn-primary" style="width:auto">{{ editingId ? 'Simpan Perubahan' : 'Buat Akun' }}</button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({ title: 'Manajemen User' })

const search = ref('')
const activeTab = ref('penuh')
const showModal = ref(false)
const editingId = ref<number | null>(null)

const tabs = [
  { label: 'Admin Pusat', value: 'penuh' },
  { label: 'Pengurus Cabang', value: 'cabang' },
  { label: 'PIC Unit', value: 'unit' },
  { label: 'Pelatih', value: 'pelatih' }
]

const aksesMap: Record<string, any[][]> = {
  penuh: [['Dashboard','ti-layout-dashboard','Akses penuh','akses-on'],['Kelola Cabang','ti-map-pin','Akses penuh','akses-on'],['Kelola Anggota','ti-users','Akses penuh','akses-on'],['BLBA','ti-wallet','Akses penuh','akses-on'],['Tes Kebugaran','ti-heartbeat','Akses penuh','akses-on'],['Monitoring Olah Nafas','ti-wind','Akses penuh','akses-on']],
  cabang: [['Dashboard','ti-layout-dashboard','Cabangnya saja','akses-on'],['Kelola Anggota','ti-users','Cabangnya saja','akses-on'],['BLBA','ti-wallet','Cabangnya saja','akses-on'],['Pengaturan','ti-settings','Tidak ada akses','akses-off'],['Tes Kebugaran','ti-heartbeat','Cabangnya saja','akses-on']],
  unit: [['Anggota Unit','ti-users','Unitnya saja','akses-on'],['BLBA Unit','ti-wallet','Unitnya saja','akses-on'],['Laporan Unit','ti-chart-bar','Unitnya saja','akses-on'],['Kelola Cabang','ti-map-pin','Tidak ada akses','akses-off'],['Pengaturan','ti-settings','Tidak ada akses','akses-off']],
  pelatih: [['Jadwal','ti-calendar','Generate QR','akses-on'],['Absensi','ti-qrcode','Input via mobile','akses-on'],['Tes Kebugaran','ti-heartbeat','Input via mobile','akses-on'],['BLBA','ti-wallet','Tidak ada akses','akses-off'],['Pengaturan','ti-settings','Tidak ada akses','akses-off']],
}

const users = ref([
  { id: 1, av: 'SA', bg: '#1a5c2a', name: 'Sri Astuti', role: 'Admin Pusat', jabatan: 'Administrator', email: 'sri.astuti@sn-pusat.org', hp: '0812-3456-7890', scope: 'Pusat · Seluruh Indonesia', akses: 'penuh', status: 'active' },
  { id: 2, av: 'BW', bg: '#5b21b6', name: 'Bambang Wiyono', role: 'Admin Pusat', jabatan: 'Administrator', email: 'bambang@sn-pusat.org', hp: '0813-1111-2222', scope: 'Pusat · Seluruh Indonesia', akses: 'penuh', status: 'active' },
  { id: 3, av: 'HW', bg: '#0c5478', name: 'Hadiwiyono W.', role: 'Pengurus Cabang', jabatan: 'Ketua Cabang', email: 'hadi@sn-yogya.org', hp: '0812-3456-7891', scope: 'Cabang · Kota Yogyakarta', akses: 'cabang', status: 'active' },
  { id: 4, av: 'SR', bg: '#0e7aad', name: 'Siti Rahayu', role: 'Pengurus Cabang', jabatan: 'Sekretaris', email: 'siti@sn-yogya.org', hp: '0814-5678-9012', scope: 'Cabang · Kota Yogyakarta', akses: 'cabang', status: 'active' },
  { id: 5, av: 'BH', bg: '#7a6000', name: 'Budi Hartono', role: 'PIC Unit', jabatan: 'PIC Unit', email: 'budi.h@sn-yogya.org', hp: '0812-0001-0002', scope: 'Unit · Unit Malioboro', akses: 'unit', status: 'active' },
  { id: 6, av: 'WS', bg: '#5b21b6', name: 'Wisnu Saputro', role: 'Pelatih', jabatan: 'Pelatih Daerah', email: 'wisnu@sn-yogya.org', hp: '0821-1234-5678', scope: 'Cabang · Kota Yogyakarta', akses: 'pelatih', status: 'active' },
])

const form = ref({
  name: '',
  akses: 'penuh',
  jabatan: '',
  email: '',
  hp: '',
  scope: ''
})

const getTabCount = (tabVal: string) => {
  return users.value.filter(u => u.akses === tabVal).length
}

const filteredUsers = computed(() => {
  return users.value.filter(u => {
    const matchesTab = u.akses === activeTab.value
    const matchesSearch = search.value === '' || 
      u.name.toLowerCase().includes(search.value.toLowerCase()) || 
      u.email.toLowerCase().includes(search.value.toLowerCase())
    return matchesTab && matchesSearch
  })
})

const openAddModal = () => {
  editingId.value = null
  form.value = {
    name: '',
    akses: activeTab.value,
    jabatan: '',
    email: '',
    hp: '',
    scope: ''
  }
  showModal.value = true
}

const openEditModal = (u: any) => {
  editingId.value = u.id
  form.value = {
    name: u.name,
    akses: u.akses,
    jabatan: u.jabatan,
    email: u.email,
    hp: u.hp,
    scope: u.scope
  }
  showModal.value = true
}

const toggleUserStatus = (u: any) => {
  const isActivating = u.status === 'nonactive'
  if (confirm(`Apakah Anda yakin ingin ${isActivating ? 'mengaktifkan' : 'menonaktifkan'} akun ${u.name}?`)) {
    u.status = isActivating ? 'active' : 'nonactive'
  }
}

const getInitials = (name: string) => {
  if (!name) return 'SA'
  const parts = name.split(' ')
  if (parts.length > 1) return (parts[0][0] + parts[1][0]).toUpperCase()
  return parts[0].substring(0, 2).toUpperCase()
}

const getAvatarBg = (name: string) => {
  const colors = ['#1a5c2a', '#5b21b6', '#0c5478', '#8c1c1c', '#7a6000', '#0e7aad']
  let hash = 0
  for (let i = 0; i < name.length; i++) hash = name.charCodeAt(i) + ((hash << 5) - hash)
  return colors[Math.abs(hash % colors.length)]
}

const saveUser = () => {
  const roleLabel: Record<string, string> = {
    penuh: 'Admin Pusat',
    cabang: 'Pengurus Cabang',
    unit: 'PIC Unit',
    pelatih: 'Pelatih'
  }

  if (editingId.value) {
    const idx = users.value.findIndex(u => u.id === editingId.value)
    if (idx !== -1) {
      users.value[idx] = {
        ...users.value[idx],
        name: form.value.name,
        akses: form.value.akses,
        role: roleLabel[form.value.akses],
        jabatan: form.value.jabatan,
        email: form.value.email,
        hp: form.value.hp,
        scope: form.value.scope,
        av: getInitials(form.value.name),
        bg: getAvatarBg(form.value.name)
      }
      alert('Informasi akun user berhasil diperbarui!')
    }
  } else {
    users.value.push({
      id: Date.now(),
      name: form.value.name,
      akses: form.value.akses,
      role: roleLabel[form.value.akses],
      jabatan: form.value.jabatan,
      email: form.value.email,
      hp: form.value.hp,
      scope: form.value.scope,
      av: getInitials(form.value.name),
      bg: getAvatarBg(form.value.name),
      status: 'active'
    })
    alert('Akun user baru berhasil dibuat!')
  }
  showModal.value = false
}
</script>

<style scoped>
.page-header { display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 20px; flex-wrap: wrap; gap: 12px; }
.page-title { font-size: 20px; font-weight: 800; margin: 0; }
.page-subtitle { font-size: 12px; color: var(--text3); margin-top: 2px; }

.header-actions { display: flex; gap: 10px; align-items: center; }
.search-box-header { display: flex; align-items: center; gap: 8px; background: var(--card); border: 1px solid var(--border); border-radius: var(--r8); padding: 7px 12px; width: 220px; }
.search-box-header i { color: var(--text3); font-size: 14px; }
.search-box-header input { border: none; outline: none; background: none; font-size: 12px; flex: 1; color: var(--text1); }

.btn { display: inline-flex; align-items: center; gap: 6px; padding: 8px 14px; border-radius: var(--r8); font-size: 12px; font-weight: 600; cursor: pointer; border: none; }
.btn-primary { background: var(--hijau); color: #fff; }
.btn-primary:hover { background: var(--hijau2); }
.btn-outline { background: #fff; color: var(--text2); border: 1px solid var(--border); }
.btn-outline:hover { background: var(--surface); }

.tab-bar { display: flex; gap: 4px; margin-bottom: 18px; border-bottom: 1px solid var(--border); }
.tab-btn { display: flex; align-items: center; gap: 6px; padding: 8px 16px; font-size: 12px; font-weight: 600; background: none; border: none; cursor: pointer; color: var(--text2); border-bottom: 2px solid transparent; margin-bottom: -1px; }
.tab-btn.active { color: var(--hijau); border-bottom-color: var(--hijau); }
.tab-count { background: var(--hijau3); color: var(--hijau); border-radius: 10px; padding: 1px 7px; font-size: 10px; }

.user-grid { display: grid; grid-template-columns: repeat(2, 1fr); gap: 14px; }
.user-card { background: var(--card); border: 1px solid var(--border); border-radius: var(--r12); padding: 16px; transition: all 0.2s; }
.user-card.disabled { opacity: 0.65; border-style: dashed; }
.user-card-top { display: flex; align-items: flex-start; gap: 12px; margin-bottom: 12px; }
.user-avatar { width: 40px; height: 40px; border-radius: 50%; display: flex; align-items: center; justify-content: center; font-size: 13px; font-weight: 800; color: #fff; flex-shrink: 0; box-shadow: 0 2px 6px rgba(0,0,0,0.1); }
.user-info { flex: 1; min-width: 0; }
.user-name { font-size: 13px; font-weight: 700; }
.user-role { font-size: 11px; color: var(--text3); }
.user-scope { font-size: 11px; color: var(--text2); display: flex; align-items: center; gap: 4px; margin-top: 2px; }
.user-scope i { color: var(--hijau); }

.akses-badge { display: inline-block; padding: 3px 10px; border-radius: 10px; font-size: 10px; font-weight: 700; white-space: nowrap; }
.akses-badge.penuh { background: #f0e0fb; color: #6b21a8; }
.akses-badge.cabang { background: var(--hijau3); color: var(--hijau); }
.akses-badge.unit { background: #fff8e0; color: #a07000; }
.akses-badge.pelatih { background: #e0f5fb; color: var(--biru); }

.user-detail { display: flex; flex-direction: column; gap: 4px; margin-bottom: 12px; padding: 10px; background: var(--surface); border-radius: var(--r8); border: 0.5px solid var(--border); }
.ud-item { font-size: 11px; color: var(--text2); display: flex; align-items: center; gap: 6px; }
.ud-item i { color: var(--text3); font-size: 12px; width: 14px; }

.user-akses-list { display: flex; flex-direction: column; gap: 3px; margin-bottom: 12px; }
.akses-item { display: flex; align-items: center; gap: 8px; padding: 5px 8px; border-radius: 6px; font-size: 11px; border: 0.5px solid transparent; }
.akses-item.akses-on { background: var(--hijau3); color: var(--hijau); border-color: rgba(46,125,50,0.1); }
.akses-item.akses-off { background: var(--surface); color: var(--text3); }
.akses-item i { font-size: 12px; width: 14px; }
.akses-item span { margin-left: auto; font-size: 10px; font-weight: 600; opacity: .8; }

.user-card-footer { display: flex; gap: 8px; border-top: 1px solid var(--border); padding-top: 12px; }
.btn-sm-outline { display: inline-flex; align-items: center; gap: 4px; padding: 6px 12px; border: 1px solid var(--border); border-radius: var(--r6); font-size: 11px; background: #fff; cursor: pointer; color: var(--text2); transition: all .15s; }
.btn-sm-outline:hover { border-color: var(--hijau); color: var(--hijau); }
.btn-sm-outline.danger { border-color: #fca5a5; color: var(--merah); }
.btn-sm-outline.danger:hover { background: #fdecea; }
.btn-sm-outline.green { border-color: var(--hijau4); color: var(--hijau); }
.btn-sm-outline.green:hover { background: var(--hijau3); }

/* Modals */
.modal-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.5); z-index: 1000; display: flex; align-items: center; justify-content: center; padding: 20px; }
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
.empty-state i { font-size: 32px; }
</style>
