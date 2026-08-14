<template>
  <div class="page-container" style="padding: 24px;">
    <div class="page-header" style="margin-bottom: 24px;">
      <h1 style="font-size: 24px; font-weight: 800; color: var(--text1); margin-bottom: 8px;">Profil Saya</h1>
      <p style="color: var(--text3); font-size: 14px;">Kelola informasi akun dan preferensi Anda.</p>
    </div>

    <div class="profile-layout" style="display: flex; gap: 24px; align-items: flex-start; flex-wrap: wrap;">
      <!-- Profil Card -->
      <div class="card" style="background: var(--card); border: 1px solid var(--border); border-radius: var(--r12); padding: 24px; flex: 1; min-width: 300px; text-align: center;">
        <div class="avatar-large" style="width: 100px; height: 100px; border-radius: 50%; background: var(--hijauSoft); color: var(--hijau); font-size: 32px; font-weight: 800; display: flex; align-items: center; justify-content: center; margin: 0 auto 16px;">
          {{ userInitials }}
        </div>
        <h2 style="font-size: 18px; font-weight: 700; color: var(--text1); margin-bottom: 4px;">{{ user.nama_lengkap || 'Sri Astuti' }}</h2>
        <div style="font-size: 13px; color: var(--text2); margin-bottom: 16px;">{{ user.email || 'admin@satrianusantara.id' }}</div>
        
        <span class="role-badge" style="display: inline-block; padding: 4px 12px; background: var(--hijau3); color: var(--hijau); border-radius: 20px; font-size: 12px; font-weight: 700; text-transform: uppercase;">
          {{ user.role_name || 'Admin Pusat' }}
        </span>

        <div style="margin-top: 24px; padding-top: 24px; border-top: 1px dashed var(--border); text-align: left;">
          <div style="display: flex; justify-content: space-between; margin-bottom: 12px; font-size: 13px;">
            <span style="color: var(--text3);">No. HP</span>
            <span style="font-weight: 600; color: var(--text1);">{{ user.phone || '-' }}</span>
          </div>
          <div style="display: flex; justify-content: space-between; margin-bottom: 12px; font-size: 13px;">
            <span style="color: var(--text3);">Cabang</span>
            <span style="font-weight: 600; color: var(--text1);">{{ user.cabang_nama || 'Pusat' }}</span>
          </div>
          <div style="display: flex; justify-content: space-between; font-size: 13px;">
            <span style="color: var(--text3);">Status</span>
            <span style="font-weight: 600; color: var(--hijau);">Aktif</span>
          </div>
        </div>
      </div>

      <!-- Settings Card -->
      <div class="card" style="background: var(--card); border: 1px solid var(--border); border-radius: var(--r12); padding: 24px; flex: 2; min-width: 300px;">
        <h3 style="font-size: 16px; font-weight: 700; color: var(--text1); margin-bottom: 16px; display: flex; align-items: center; gap: 8px;">
          <i class="ti ti-settings" style="color: var(--text3);"></i> Pengaturan Akun
        </h3>

        <div class="form-group" style="margin-bottom: 16px;">
          <label style="display: block; font-size: 12px; font-weight: 600; color: var(--text2); margin-bottom: 6px;">Nama Lengkap</label>
          <input type="text" :value="user.nama_lengkap || 'Sri Astuti'" class="form-input" disabled />
        </div>

        <div class="form-group" style="margin-bottom: 16px;">
          <label style="display: block; font-size: 12px; font-weight: 600; color: var(--text2); margin-bottom: 6px;">Email</label>
          <input type="email" :value="user.email || 'admin@satrianusantara.id'" class="form-input" disabled />
        </div>

        <div style="display: flex; justify-content: flex-end; margin-top: 24px;">
          <button class="btn btn-primary" style="background: var(--hijau); color: white; padding: 10px 20px; border: none; border-radius: 8px; font-weight: 600; cursor: pointer;">
            Simpan Perubahan
          </button>
        </div>

        <hr style="border: 0; border-top: 1px solid var(--border); margin: 32px 0;" />

        <h3 style="font-size: 16px; font-weight: 700; color: var(--text1); margin-bottom: 16px; display: flex; align-items: center; gap: 8px;">
          <i class="ti ti-lock" style="color: var(--text3);"></i> Keamanan
        </h3>

        <div class="form-group" style="margin-bottom: 16px;">
          <label style="display: block; font-size: 12px; font-weight: 600; color: var(--text2); margin-bottom: 6px;">Password Lama</label>
          <input type="password" placeholder="••••••••" class="form-input" />
        </div>
        <div class="form-group" style="margin-bottom: 16px;">
          <label style="display: block; font-size: 12px; font-weight: 600; color: var(--text2); margin-bottom: 6px;">Password Baru</label>
          <input type="password" placeholder="••••••••" class="form-input" />
        </div>

        <div style="display: flex; justify-content: flex-end; margin-top: 24px;">
          <button class="btn btn-outline" style="background: transparent; color: var(--text1); border: 1px solid var(--border); padding: 10px 20px; border-radius: 8px; font-weight: 600; cursor: pointer;">
            Ganti Password
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useAuthStore } from '~/stores/auth'

definePageMeta({
  title: 'Profil Saya'
})

const authStore = useAuthStore()

const user = computed(() => authStore.user || {})

const userInitials = computed(() => {
  if (!user.value.nama_lengkap) return 'SA'
  const parts = user.value.nama_lengkap.split(' ')
  if (parts.length > 1) {
    return (parts[0][0] + parts[1][0]).toUpperCase()
  }
  return parts[0].substring(0, 2).toUpperCase()
})
</script>

<style scoped>
.form-input {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid var(--border);
  border-radius: var(--r8);
  font-size: 13px;
  background: var(--surface);
  color: var(--text1);
}
.form-input:disabled {
  background: var(--bg);
  color: var(--text3);
  cursor: not-allowed;
}
</style>
