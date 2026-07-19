<template>
  <div id="page-konten">
    <!-- Header -->
    <div class="page-header">
      <div>
        <h1 class="page-title">Konten & Artikel</h1>
        <div class="page-subtitle">Kelola pengumuman, artikel, dan info kegiatan untuk anggota</div>
      </div>
      <button class="btn btn-primary" @click="openAddModal"><i class="ti ti-plus"></i> Buat Konten</button>
    </div>

    <!-- Filters -->
    <div class="filter-bar">
      <div class="search-box">
        <i class="ti ti-search"></i>
        <input v-model="search" placeholder="Cari judul konten..." />
      </div>
      <select v-model="filterJenis" class="filter-select">
        <option value="">Semua Jenis</option>
        <option>Pengumuman</option>
        <option>Artikel</option>
        <option>Info Kegiatan</option>
      </select>
      <select v-model="filterScope" class="filter-select">
        <option value="">Semua Scope</option>
        <option>Nasional</option>
        <option>Provinsi</option>
        <option>Cabang</option>
      </select>
      <select v-model="filterStatus" class="filter-select">
        <option value="">Semua Status</option>
        <option value="published">Published</option>
        <option value="draft">Draft</option>
      </select>
    </div>

    <!-- Konten Grid -->
    <div v-if="filteredKonten.length === 0" class="empty-state">
      <i class="ti ti-notes-off"></i>
      <p>Tidak ada konten atau artikel ditemukan.</p>
    </div>

    <div v-else class="konten-grid">
      <div v-for="k in filteredKonten" :key="k.id" class="konten-card">
        <div v-if="k.image" class="konten-card-image">
          <img :src="k.image" alt="Banner Konten" />
        </div>
        <div class="konten-card-header">
          <span :class="['jenis-badge', getJenisClass(k.jenis)]">{{ k.jenis }}</span>
          <span :class="['status-badge', k.status === 'published' ? 'pub' : 'draft']">
            {{ k.status === 'published' ? 'Published' : 'Draft' }}
          </span>
        </div>
        <div class="konten-judul">{{ k.judul }}</div>
        <div class="konten-meta">
          <span><i class="ti ti-map-pin"></i> {{ k.scope }}</span>
          <span><i class="ti ti-calendar"></i> {{ formatDate(k.tanggal) }}</span>
          <span><i class="ti ti-user"></i> {{ k.author }}</span>
        </div>
        <div class="konten-actions">
          <button class="btn-sm-outline" @click="openPreview(k)"><i class="ti ti-eye"></i> Preview</button>
          <button class="btn-sm-outline" @click="openEditModal(k)"><i class="ti ti-edit"></i> Edit</button>
          <button v-if="k.status === 'draft'" class="btn-sm-outline green" @click="publishContent(k)">
            <i class="ti ti-send"></i> Publish
          </button>
          <button class="btn-sm-outline danger" @click="deleteContent(k.id)"><i class="ti ti-trash"></i> Hapus</button>
        </div>
      </div>
    </div>

    <!-- MODAL: CREATE / EDIT KONTEN -->
    <div v-if="showModal" class="modal-overlay" @click.self="showModal = false">
      <div class="modal-card">
        <div class="modal-header">
          <div class="modal-title">{{ editingId ? 'Edit Konten & Artikel' : 'Buat Konten Baru' }}</div>
          <button class="modal-close" @click="showModal = false"><i class="ti ti-x"></i></button>
        </div>
        <form @submit.prevent="saveContent" class="modal-form">
          <div class="form-group">
            <label class="form-label">Judul Konten</label>
            <input v-model="form.judul" type="text" class="form-input" placeholder="Masukkan judul menarik..." required />
          </div>
          <div class="form-row">
            <div class="form-group half">
              <label class="form-label">Jenis Konten</label>
              <select v-model="form.jenis" class="form-select" required>
                <option>Pengumuman</option>
                <option>Artikel</option>
                <option>Info Kegiatan</option>
              </select>
            </div>
            <div class="form-group half">
              <label class="form-label">Scope Akses</label>
              <select v-model="form.scope" class="form-select" required>
                <option>Nasional</option>
                <option>Provinsi DIY</option>
                <option>Provinsi DKI Jakarta</option>
                <option>Cabang Yogyakarta</option>
                <option>Cabang Bandung</option>
              </select>
            </div>
          </div>
          <div class="form-row">
            <div class="form-group half">
              <label class="form-label">Tanggal Rilis</label>
              <input v-model="form.tanggal" type="date" class="form-input" required />
            </div>
            <div class="form-group half">
              <label class="form-label">Status Awal</label>
              <select v-model="form.status" class="form-select" required>
                <option value="draft">Draft (Simpan saja)</option>
                <option value="published">Published (Siap Tayang)</option>
              </select>
            </div>
          </div>
          <div class="form-group">
            <label class="form-label">Gambar Banner</label>
            <div class="image-upload-wrapper">
              <div v-if="form.image" class="image-preview-container">
                <img :src="form.image" class="image-preview" />
                <button type="button" class="btn-remove-image" @click="form.image = ''">
                  <i class="ti ti-trash"></i> Hapus Foto
                </button>
              </div>
              <div v-else class="upload-placeholder" @click="triggerFileInput">
                <i class="ti ti-camera"></i>
                <span>Upload Foto Banner (PNG/JPG)</span>
              </div>
              <input
                ref="fileInput"
                type="file"
                accept="image/*"
                style="display: none"
                @change="handleImageUpload"
              />
            </div>
          </div>
          <div class="form-group">
            <label class="form-label">Isi Konten / Artikel</label>
            <textarea v-model="form.isi" class="form-input text-area" rows="6" placeholder="Tuliskan detail berita, artikel, atau pengumuman di sini..." required></textarea>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-outline" style="width:auto" @click="showModal = false">Batal</button>
            <button type="submit" class="btn btn-primary" style="width:auto">{{ editingId ? 'Simpan Perubahan' : 'Terbitkan Konten' }}</button>
          </div>
        </form>
      </div>
    </div>

    <!-- MODAL: PREVIEW KONTEN -->
    <div v-if="showPreviewModal && previewItem" class="modal-overlay" @click.self="showPreviewModal = false">
      <div class="modal-card large">
        <div class="modal-header">
          <div class="modal-title">Pratinjau Tampilan Artikel</div>
          <button class="modal-close" @click="showPreviewModal = false"><i class="ti ti-x"></i></button>
        </div>
        <div class="modal-body" style="padding: 20px;">
          <div v-if="previewItem.image" class="preview-banner-container">
            <img :src="previewItem.image" alt="Preview Banner" class="preview-banner" />
          </div>
          <div class="preview-badge-row">
            <span :class="['jenis-badge', getJenisClass(previewItem.jenis)]">{{ previewItem.jenis }}</span>
            <span :class="['status-badge', previewItem.status === 'published' ? 'pub' : 'draft']">
              {{ previewItem.status === 'published' ? 'Published' : 'Draft' }}
            </span>
          </div>
          <h1 class="preview-judul">{{ previewItem.judul }}</h1>
          <div class="preview-meta">
            <span><i class="ti ti-user"></i> {{ previewItem.author }}</span>
            <span><i class="ti ti-calendar"></i> {{ formatDate(previewItem.tanggal) }}</span>
            <span><i class="ti ti-map-pin"></i> Scope: {{ previewItem.scope }}</span>
          </div>
          <div class="preview-body-content">
            <p v-for="(p, idx) in splitParagraphs(previewItem.isi)" :key="idx">{{ p }}</p>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn-primary" style="width:auto" @click="showPreviewModal = false">Tutup</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({ title: 'Konten & Artikel' })

const search = ref('')
const filterJenis = ref('')
const filterScope = ref('')
const filterStatus = ref('')
const showModal = ref(false)
const showPreviewModal = ref(false)
const editingId = ref<number | null>(null)
const previewItem = ref<any>(null)
const fileInput = ref<HTMLInputElement | null>(null)

const kontenData = ref([
  { id: 1, jenis: 'Pengumuman', judul: 'Jadwal Latgab + EKT Jurus Provinsi DIY — Juli 2026', scope: 'Provinsi DIY', tanggal: '2026-07-15', author: 'Sri Astuti', status: 'published', isi: 'Diberitahukan kepada seluruh anggota Cabang Yogyakarta, Sleman, dan Bantul bahwa Latihan Gabungan dan Evaluasi Kenaikan Tingkat (EKT) daerah Yogyakarta akan diselenggarakan pada tanggal 26 Juli 2026 bertempat di Alun-Alun Kidul Yogyakarta.\n\nHarap setiap unit mengumpulkan daftar peserta selambat-lambatnya 19 Juli 2026 melalui aplikasi admin.', image: 'https://images.unsplash.com/photo-1544033527-b192daee1f5b?w=400' },
  { id: 2, jenis: 'Info Kegiatan', judul: 'Pendaftaran Pelatnas 2026 Dibuka — Kuota Terbatas!', scope: 'Nasional', tanggal: '2026-07-10', author: 'Sri Astuti', status: 'published', isi: 'Pusat mengumumkan pembukaan pendaftaran Latihan Nasional (Pelatnas) Angkatan 2026 yang akan diselenggarakan di Kaliurang, Yogyakarta.\n\nKegiatan ini ditujukan untuk tingkatan minimal Calon Pelatih (PH Jurus 3 ke atas). Pendaftaran akan ditutup otomatis apabila kuota 150 peserta telah terpenuhi.', image: 'https://images.unsplash.com/photo-1555597673-b21d5c935865?w=400' },
  { id: 3, jenis: 'Artikel', judul: 'Tips Meningkatkan Performa Latihan Seni Pernapasan', scope: 'Nasional', tanggal: '2026-07-08', author: 'Bambang Wiyono', status: 'draft', isi: 'Melatih pernapasan memerlukan konsentrasi penuh dan keharmonisan antara gerakan fisik serta pengaturan udara.\n\nBeberapa tips utama: 1) Lakukan latihan rutin minimal 3 kali seminggu, 2) Selalu perhatikan ketepatan jurus sebelum mempercepat gerakan, 3) Jaga asupan air putih dan nutrisi seimbang untuk mendukung stamina tubuh.', image: 'https://images.unsplash.com/photo-1517649763962-0c623066013b?w=400' },
  { id: 4, jenis: 'Pengumuman', judul: 'Perubahan Jadwal Latihan Unit Malioboro — Agu 2026', scope: 'Cabang Yogyakarta', tanggal: '2026-07-05', author: 'Hadiwiyono W.', status: 'published', isi: 'Diberitahukan kepada anggota Unit Malioboro bahwa mulai tanggal 1 Agustus 2026, sesi latihan rutin hari Sabtu malam dipindahkan ke hari Minggu pagi pukul 06.00 WIB bertempat di pelataran depan benteng Vredeburg.\n\nHal ini dikarenakan adanya pengerjaan revitalisasi kawasan pedestrian Malioboro.', image: 'https://images.unsplash.com/photo-1517838277536-f5f99be501cd?w=400' },
])

const form = ref({
  judul: '',
  jenis: 'Pengumuman',
  scope: 'Nasional',
  tanggal: '2026-07-18',
  status: 'draft',
  isi: '',
  image: ''
})

const getJenisClass = (jenis: string) => {
  if (jenis === 'Pengumuman') return 'pengumuman'
  if (jenis === 'Info Kegiatan') return 'kegiatan'
  return 'artikel'
}

const filteredKonten = computed(() => {
  return kontenData.value.filter(k => {
    const matchesSearch = search.value === '' || k.judul.toLowerCase().includes(search.value.toLowerCase())
    const matchesJenis = filterJenis.value === '' || k.jenis === filterJenis.value
    const matchesScope = filterScope.value === '' || k.scope.toLowerCase().includes(filterScope.value.toLowerCase())
    const matchesStatus = filterStatus.value === '' || k.status === filterStatus.value
    return matchesSearch && matchesJenis && matchesScope && matchesStatus
  })
})

const triggerFileInput = () => {
  fileInput.value?.click()
}

const handleImageUpload = (event: Event) => {
  const target = event.target as HTMLInputElement
  if (target.files && target.files[0]) {
    const file = target.files[0]
    const reader = new FileReader()
    reader.onload = (e) => {
      form.value.image = e.target?.result as string
    }
    reader.readAsDataURL(file)
  }
}

const openAddModal = () => {
  editingId.value = null
  form.value = {
    judul: '',
    jenis: 'Pengumuman',
    scope: 'Nasional',
    tanggal: new Date().toISOString().substring(0, 10),
    status: 'draft',
    isi: '',
    image: ''
  }
  showModal.value = true
}

const openEditModal = (k: any) => {
  editingId.value = k.id
  form.value = {
    judul: k.judul,
    jenis: k.jenis,
    scope: k.scope,
    tanggal: k.tanggal,
    status: k.status,
    isi: k.isi,
    image: k.image || ''
  }
  showModal.value = true
}

const openPreview = (k: any) => {
  previewItem.value = k
  showPreviewModal.value = true
}

const publishContent = (k: any) => {
  k.status = 'published'
  alert(`Konten "${k.judul}" berhasil diterbitkan!`)
}

const deleteContent = (id: number) => {
  if (confirm('Apakah Anda yakin ingin menghapus konten ini?')) {
    kontenData.value = kontenData.value.filter(k => k.id !== id)
  }
}

const saveContent = () => {
  if (editingId.value) {
    const idx = kontenData.value.findIndex(k => k.id === editingId.value)
    if (idx !== -1) {
      kontenData.value[idx] = {
        ...kontenData.value[idx],
        judul: form.value.judul,
        jenis: form.value.jenis,
        scope: form.value.scope,
        tanggal: form.value.tanggal,
        status: form.value.status,
        isi: form.value.isi,
        image: form.value.image
      }
      alert('Konten berhasil diperbarui!')
    }
  } else {
    kontenData.value.unshift({
      id: Date.now(),
      jenis: form.value.jenis,
      judul: form.value.judul,
      scope: form.value.scope,
      tanggal: form.value.tanggal,
      author: 'Sri Astuti',
      status: form.value.status,
      isi: form.value.isi,
      image: form.value.image
    })
    alert('Konten baru berhasil dibuat!')
  }
  showModal.value = false
}

const formatDate = (dateStr: string) => {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  return d.toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' })
}

const splitParagraphs = (text: string) => {
  if (!text) return []
  return text.split('\n\n')
}
</script>

<style scoped>
.page-header { display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 20px; flex-wrap: wrap; gap: 12px; }
.page-title { font-size: 20px; font-weight: 800; margin: 0; }
.page-subtitle { font-size: 12px; color: var(--text3); margin-top: 2px; }

.btn { display: inline-flex; align-items: center; gap: 6px; padding: 8px 14px; border-radius: var(--r8); font-size: 12px; font-weight: 600; cursor: pointer; border: none; }
.btn-primary { background: var(--hijau); color: #fff; }
.btn-primary:hover { background: var(--hijau2); }
.btn-outline { background: #fff; color: var(--text2); border: 1px solid var(--border); }
.btn-outline:hover { background: var(--surface); }

.filter-bar { display: flex; gap: 8px; margin-bottom: 16px; flex-wrap: wrap; }
.search-box { display: flex; align-items: center; gap: 8px; flex: 1; min-width: 200px; background: var(--card); border: 1px solid var(--border); border-radius: var(--r8); padding: 7px 12px; }
.search-box i { color: var(--text3); }
.search-box input { border: none; outline: none; background: none; font-size: 12px; flex: 1; color: var(--text1); }
.filter-select { padding: 7px 10px; border: 1px solid var(--border); border-radius: var(--r8); font-size: 12px; background: var(--card); color: var(--text2); cursor: pointer; }

.konten-grid { display: grid; grid-template-columns: repeat(2, 1fr); gap: 14px; }
.konten-card { background: var(--card); border: 1px solid var(--border); border-radius: var(--r12); padding: 16px; display: flex; flex-direction: column; }
.konten-card-image { height: 130px; width: calc(100% + 32px); margin-left: -16px; margin-top: -16px; margin-bottom: 12px; overflow: hidden; border-radius: var(--r12) var(--r12) 0 0; }
.konten-card-image img { width: 100%; height: 100%; object-fit: cover; }

.konten-card-header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 10px; }

.jenis-badge { display: inline-block; padding: 3px 9px; border-radius: 10px; font-size: 10px; font-weight: 700; text-transform: uppercase; }
.jenis-badge.pengumuman { background: #fff8e0; color: #a07000; }
.jenis-badge.kegiatan { background: #e0f5fb; color: var(--biru); }
.jenis-badge.artikel { background: var(--hijau3); color: var(--hijau); }

.status-badge { display: inline-block; padding: 3px 9px; border-radius: 10px; font-size: 10px; font-weight: 700; text-transform: capitalize; }
.status-badge.pub { background: var(--hijau3); color: var(--hijau); }
.status-badge.draft { background: var(--surface); color: var(--text3); border: 0.5px dashed var(--border); }

.konten-judul { font-size: 14px; font-weight: 700; margin-bottom: 10px; line-height: 1.4; color: var(--text1); }
.konten-meta { display: flex; flex-wrap: wrap; gap: 10px; font-size: 11px; color: var(--text3); margin-bottom: 14px; }
.konten-meta span { display: flex; align-items: center; gap: 4px; }

.konten-actions { display: flex; gap: 6px; margin-top: auto; }
.btn-sm-outline { display: inline-flex; align-items: center; gap: 4px; padding: 6px 10px; border: 1px solid var(--border); border-radius: var(--r6); font-size: 11px; background: #fff; cursor: pointer; color: var(--text2); transition: all 0.15s; }
.btn-sm-outline:hover { border-color: var(--hijau); color: var(--hijau); }
.btn-sm-outline.green { border-color: var(--hijau4); color: var(--hijau); background: var(--hijau3); }
.btn-sm-outline.danger:hover { border-color: #fca5a5; color: var(--merah); background: #fdecea; }

/* Modals */
.modal-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.5); z-index: 1000; display: flex; align-items: center; justify-content: center; padding: 20px; }
.modal-card { background: var(--card); border-radius: var(--r12); width: 100%; max-width: 520px; padding: 20px; box-shadow: var(--shadow-lg); }
.modal-card.large { max-width: 680px; }
.modal-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px; }
.modal-title { font-size: 16px; font-weight: 700; }
.modal-close { background: none; border: none; font-size: 16px; cursor: pointer; color: var(--text3); }
.modal-form { display: flex; flex-direction: column; gap: 12px; }
.form-group { display: flex; flex-direction: column; gap: 4px; }
.form-row { display: flex; gap: 10px; }
.form-group.half { flex: 1; }
.form-label { font-size: 11px; font-weight: 600; color: var(--text2); }
.form-input, .form-select { padding: 8px 12px; border: 1.5px solid var(--border2); border-radius: var(--r8); font-size: 12px; outline: none; background: #fff; width: 100%; }
.form-input:focus, .form-select:focus { border-color: var(--hijau); }
.form-input.text-area { resize: vertical; font-family: inherit; }
.modal-footer { display: flex; justify-content: flex-end; gap: 8px; margin-top: 12px; }

/* Image upload */
.image-upload-wrapper { border: 1.5px dashed var(--border2); border-radius: var(--r8); padding: 12px; background: var(--surface); display: flex; justify-content: center; align-items: center; min-height: 100px; cursor: pointer; transition: all 0.15s; }
.image-upload-wrapper:hover { border-color: var(--hijau); }
.upload-placeholder { display: flex; flex-direction: column; align-items: center; gap: 6px; color: var(--text3); font-size: 11px; }
.upload-placeholder i { font-size: 20px; }
.image-preview-container { position: relative; width: 100%; height: 150px; border-radius: var(--r6); overflow: hidden; }
.image-preview { width: 100%; height: 100%; object-fit: cover; }
.btn-remove-image { position: absolute; bottom: 8px; right: 8px; background: rgba(220, 38, 38, 0.9); color: #fff; border: none; padding: 4px 8px; border-radius: var(--r4); font-size: 10px; font-weight: 700; cursor: pointer; display: flex; align-items: center; gap: 4px; }

/* Preview Modal content */
.preview-banner-container { height: 220px; width: 100%; overflow: hidden; border-radius: var(--r8); margin-bottom: 16px; }
.preview-banner { width: 100%; height: 100%; object-fit: cover; }
.preview-badge-row { display: flex; gap: 6px; margin-bottom: 10px; }
.preview-judul { font-size: 18px; font-weight: 800; line-height: 1.35; margin-bottom: 10px; color: var(--text1); }
.preview-meta { display: flex; gap: 14px; font-size: 11px; color: var(--text3); border-bottom: 1px solid var(--border); padding-bottom: 12px; margin-bottom: 14px; }
.preview-meta span { display: flex; align-items: center; gap: 4px; }
.preview-body-content { font-size: 13px; line-height: 1.6; color: var(--text2); }
.preview-body-content p { margin-bottom: 12px; }

.empty-state { display: flex; flex-direction: column; align-items: center; justify-content: center; padding: 40px; font-size: 13px; color: var(--text3); gap: 10px; background: var(--card); border-radius: var(--r12); border: 1px dashed var(--border); margin-top: 10px; }
.empty-state i { font-size: 32px; }
</style>
