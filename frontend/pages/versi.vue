<template>
  <div class="page-container">
    <div class="page-header">
      <div class="header-left">
        <h1 class="page-title"><i class="ti ti-versions"></i> Kelola Versi Aplikasi</h1>
        <p class="page-subtitle">Atur versi aplikasi mobile dan update mandatory.</p>
      </div>
      <div class="header-actions">
        <button class="btn btn-primary" @click="openAddModal">
          <i class="ti ti-plus"></i> Tambah Versi Baru
        </button>
      </div>
    </div>

    <div class="card p-0">
      <div class="table-responsive">
        <table class="table" style="width: 100%;">
          <thead>
            <tr>
              <th>Versi (Nama)</th>
              <th>Build Number</th>
              <th>Status Update</th>
              <th>Tanggal Rilis</th>
              <th>Catatan Rilis</th>
              <th width="120" style="text-align:right;">Aksi</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="loading">
              <td colspan="6" class="text-center py-4 text-muted">Memuat data versi...</td>
            </tr>
            <tr v-else-if="versions.length === 0">
              <td colspan="6" class="text-center py-4 text-muted">Belum ada rilis versi aplikasi.</td>
            </tr>
            <tr v-for="v in versions" :key="v.id">
              <td style="font-weight: 500;">
                <div class="d-flex align-items-center">
                  <span class="badge" style="background:var(--hijau-transparan);color:var(--hijau);margin-right:8px;">v{{ v.version_name }}</span>
                </div>
              </td>
              <td>Code: {{ v.build_number }}</td>
              <td>
                <span class="badge" :class="v.is_mandatory ? 'badge-danger' : 'badge-soft-success'">
                  {{ v.is_mandatory ? 'Wajib Update' : 'Opsional' }}
                </span>
              </td>
              <td style="color:var(--text3);font-size:13px;">{{ new Date(v.created_at).toLocaleDateString('id-ID') }}</td>
              <td style="color:var(--text3);font-size:13px;max-width:200px;white-space:nowrap;overflow:hidden;text-overflow:ellipsis;" :title="v.release_notes">
                {{ v.release_notes }}
              </td>
              <td style="text-align:right;">
                <button class="btn-icon text-danger" @click="deleteVersion(v.id)" title="Hapus">
                  <i class="ti ti-trash"></i>
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Modal Form -->
    <div v-if="showModal" class="modal-backdrop" @click.self="showModal = false">
      <div class="modal-content" style="max-width:500px;">
        <div class="modal-header">
          <h3 class="modal-title">Tambah Versi Rilis Baru</h3>
          <button class="modal-close" @click="showModal = false"><i class="ti ti-x"></i></button>
        </div>
        <div class="modal-body">
          <form @submit.prevent="saveVersion" class="modal-form">
            <div class="form-group">
              <label>Versi Name (Contoh: 1.0.1)</label>
              <input type="text" v-model="form.version_name" class="form-control" required placeholder="1.0.1" />
            </div>
            <div class="form-group">
              <label>Build Number / Version Code (Angka)</label>
              <input type="number" v-model="form.build_number" class="form-control" required placeholder="2" />
              <small style="color:var(--text3);margin-top:4px;display:block;">Harus lebih besar dari versi sebelumnya agar terdeteksi mobile app.</small>
            </div>
            <div class="form-group">
              <label>Catatan Rilis (Release Notes)</label>
              <textarea v-model="form.release_notes" class="form-control" rows="3" placeholder="Fitur baru, perbaikan bug..."></textarea>
            </div>
            <div class="form-group checkbox-group" style="margin-top: 16px;">
              <label class="custom-checkbox">
                <input type="checkbox" v-model="form.is_mandatory">
                <span class="checkmark"></span>
                <span class="label-text">Wajib Update (Mandatory)</span>
              </label>
              <small style="color:var(--text3);margin-top:4px;display:block;margin-left:28px;">Centang jika update ini bersifat kritis dan aplikasi lama tidak boleh digunakan.</small>
            </div>

            <div class="form-actions mt-4">
              <button type="button" class="btn btn-outline" @click="showModal = false">Batal</button>
              <button type="submit" class="btn btn-primary" :disabled="submitting">
                <i class="ti ti-device-floppy"></i> {{ submitting ? 'Menyimpan...' : 'Rilis Versi' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api } from '~/utils/api'

definePageMeta({
  layout: 'default',
  middleware: ['auth']
})

const versions = ref<any[]>([])
const loading = ref(true)
const showModal = ref(false)
const submitting = ref(false)

const form = ref({
  version_name: '',
  build_number: 1,
  release_notes: '',
  is_mandatory: false
})

const fetchVersions = async () => {
  loading.value = true
  try {
    const res = await api.get('/versions')
    versions.value = res || []
  } catch (error) {
    console.error('Failed to fetch versions', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchVersions()
})

const openAddModal = () => {
  // Auto increment based on latest
  let nextBuild = 1
  let nextName = '1.0.0'
  if (versions.value.length > 0) {
    const latest = versions.value[0]
    nextBuild = latest.build_number + 1
    
    // Simple auto increment for version name (e.g. 1.0.1 -> 1.0.2)
    const parts = latest.version_name.split('.')
    if (parts.length === 3) {
      parts[2] = String(parseInt(parts[2]) + 1)
      nextName = parts.join('.')
    }
  }

  form.value = {
    version_name: nextName,
    build_number: nextBuild,
    release_notes: '',
    is_mandatory: false
  }
  showModal.value = true
}

const saveVersion = async () => {
  submitting.value = true
  try {
    await api.post('/versions', {
      version_name: form.value.version_name,
      build_number: Number(form.value.build_number),
      release_notes: form.value.release_notes,
      is_mandatory: form.value.is_mandatory
    })
    alert('Versi baru berhasil dirilis!')
    showModal.value = false
    await fetchVersions()
  } catch (error: any) {
    alert(error.response?.data?.message || 'Gagal menyimpan versi baru')
  } finally {
    submitting.value = false
  }
}

const deleteVersion = async (id: string) => {
  if (confirm('Apakah Anda yakin ingin menghapus versi ini? (Hati-hati, ini bisa mempengaruhi pengecekan update pada mobile app)')) {
    try {
      await api.delete(`/versions/${id}`)
      await fetchVersions()
    } catch (error: any) {
      alert(error.response?.data?.message || 'Gagal menghapus versi')
    }
  }
}
</script>

<style scoped>
/* Sama dengan CSS halaman admin lain */
.badge-soft-success {
  background: var(--hijau-transparan);
  color: var(--hijau);
}
.badge-danger {
  background: #fef2f2;
  color: #ef4444;
}
.custom-checkbox {
  display: flex;
  align-items: center;
  cursor: pointer;
  position: relative;
}
.custom-checkbox input {
  position: absolute;
  opacity: 0;
  cursor: pointer;
  height: 0;
  width: 0;
}
.checkmark {
  height: 20px;
  width: 20px;
  background-color: #fff;
  border: 1.5px solid var(--border);
  border-radius: 4px;
  margin-right: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
}
.custom-checkbox input:checked ~ .checkmark {
  background-color: var(--hijau);
  border-color: var(--hijau);
}
.custom-checkbox input:checked ~ .checkmark:after {
  content: "✓";
  color: white;
  font-size: 14px;
}
</style>
