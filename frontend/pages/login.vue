<template>
  <div class="login-page">
    <!-- Visual Branding Left Panel -->
    <div class="login-bg">
      <div class="bg-gradient-shapes">
        <div class="shape shape-1"></div>
        <div class="shape shape-2"></div>
      </div>
      <div class="login-brand">
        <div class="brand-logo">
          <div class="brand-icon">SN</div>
        </div>
        <div class="brand-name">Satria Nusantara</div>
        <div class="brand-sub">Sistem Informasi Manajemen Anggota Nasional</div>
        <div class="brand-features">
          <div class="bf-item"><i class="ti ti-circle-check"></i> Kelola seluruh Cabang &amp; Unit Latihan</div>
          <div class="bf-item"><i class="ti ti-circle-check"></i> Monitoring Iuran BLBA &amp; Kehadiran secara real-time</div>
          <div class="bf-item"><i class="ti ti-circle-check"></i> Pembuatan Absensi QR Code Cepat &amp; Aman</div>
          <div class="bf-item"><i class="ti ti-circle-check"></i> Laporan Statistik Perkembangan Nasional</div>
        </div>
      </div>
    </div>

    <!-- Login Panel Right Side -->
    <div class="login-panel">
      <div class="login-card">
        <div class="login-header">
          <h2 class="login-title">Selamat Datang</h2>
          <p class="login-sub">Masuk dengan kredensial akun pengurus Anda</p>
        </div>

        <form @submit.prevent="handleLogin" class="login-form">
          <div class="form-group">
            <label class="form-label">Alamat Email</label>
            <div class="input-wrapper">
              <i class="ti ti-mail input-icon"></i>
              <input v-model="email" type="email" class="form-input" placeholder="admin@satria-nusantara.org" required />
            </div>
          </div>

          <div class="form-group">
            <label class="form-label">Kata Sandi</label>
            <div class="input-wrapper">
              <i class="ti ti-lock input-icon"></i>
              <input v-model="password" :type="showPw ? 'text' : 'password'" class="form-input" placeholder="••••••••" required />
              <button type="button" class="pw-toggle" @click="showPw = !showPw">
                <i :class="showPw ? 'ti ti-eye-off' : 'ti ti-eye'"></i>
              </button>
            </div>
          </div>

          <div class="login-extra">
            <label class="remember-label">
              <input type="checkbox" v-model="rememberMe" class="custom-checkbox" /> Ingat Saya
            </label>
            <a href="#" class="forgot-link">Lupa kata sandi?</a>
          </div>

          <button type="submit" class="btn-login" :disabled="loading">
            <i v-if="loading" class="ti ti-loader-2 spin"></i>
            <i v-else class="ti ti-login"></i>
            {{ loading ? 'Sedang Masuk...' : 'Masuk ke Dashboard' }}
          </button>

          <div v-if="errorMsg" class="error-msg">
            <i class="ti ti-alert-circle"></i> {{ errorMsg }}
          </div>
        </form>

        <div class="login-divider"><span>Lembaga Seni Pernapasan</span></div>

        <div class="login-footer">
          <div class="login-version">Versi 8.0 (LSP-SIMA)</div>
          <div class="login-org">© 2026 Yayasan Satria Nusantara</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from '~/stores/auth'

definePageMeta({
  title: 'Login — Satria Nusantara Web Admin',
  layout: false  // No sidebar for login
})

const email = ref('')
const password = ref('')
const showPw = ref(false)
const rememberMe = ref(false)
const loading = ref(false)
const errorMsg = ref('')

const authStore = useAuthStore()
const api = useApi()

const handleLogin = async () => {
  loading.value = true
  errorMsg.value = ''
  try {
    const data = await api.post('/auth/login', {
      email: email.value,
      password: password.value
    })
    authStore.setAuth(data.token, data.user)
    navigateTo('/')
  } catch (e: any) {
    errorMsg.value = e.message || 'Email atau password salah. Silakan coba lagi.'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-page {
  display: flex;
  height: 100vh;
  width: 100vw;
  overflow: hidden;
  font-family: 'Inter', system-ui, -apple-system, sans-serif;
  background: var(--bg);
}

/* Background panel left */
.login-bg {
  flex: 1.2;
  background: linear-gradient(135deg, #103c1b 0%, #1a5c2a 60%, #2e8b3a 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 60px;
  position: relative;
  overflow: hidden;
}

.bg-gradient-shapes .shape {
  position: absolute;
  border-radius: 50%;
  background: radial-gradient(circle, rgba(255,255,255,0.08) 0%, rgba(255,255,255,0) 70%);
}
.shape-1 {
  width: 600px; height: 600px;
  top: -150px; left: -150px;
}
.shape-2 {
  width: 500px; height: 500px;
  bottom: -100px; right: -100px;
}

.login-brand {
  z-index: 2;
  max-width: 440px;
  color: #ffffff;
}

.brand-logo {
  margin-bottom: 24px;
}

.brand-icon {
  width: 64px; height: 64px;
  background: rgba(255, 255, 255, 0.12);
  border: 1.5px solid rgba(255, 255, 255, 0.25);
  border-radius: 16px;
  display: flex; align-items: center; justify-content: center;
  font-size: 24px; font-weight: 900; color: #fff;
  backdrop-filter: blur(8px);
}

.brand-name {
  font-size: 32px;
  font-weight: 800;
  letter-spacing: -0.5px;
  margin-bottom: 8px;
}

.brand-sub {
  font-size: 15px;
  color: rgba(255, 255, 255, 0.7);
  margin-bottom: 48px;
  font-weight: 500;
}

.brand-features {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.bf-item {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 14px;
  color: rgba(255, 255, 255, 0.85);
  font-weight: 500;
}

.bf-item i {
  color: #4ade80;
  font-size: 18px;
}

/* Right input panel */
.login-panel {
  flex: 0.9;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--card);
  padding: 40px;
  box-shadow: -4px 0 24px rgba(0, 0, 0, 0.03);
}

.login-card {
  width: 100%;
  max-width: 380px;
}

.login-header {
  margin-bottom: 32px;
}

.login-title {
  font-size: 26px;
  font-weight: 800;
  color: var(--text1);
  letter-spacing: -0.5px;
  margin-bottom: 6px;
}

.login-sub {
  font-size: 13px;
  color: var(--text3);
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: 18px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.form-label {
  font-size: 12px;
  font-weight: 600;
  color: var(--text2);
}

.input-wrapper {
  position: relative;
  display: flex;
  align-items: center;
}

.input-icon {
  position: absolute;
  left: 14px;
  font-size: 16px;
  color: var(--text3);
}

.form-input {
  width: 100%;
  padding: 12px 14px 12px 42px;
  border: 1.5px solid var(--border2);
  border-radius: var(--r8);
  font-size: 13px;
  color: var(--text1);
  background: #ffffff;
  outline: none;
  transition: all 0.15s ease-in-out;
}

.form-input:focus {
  border-color: var(--hijau);
  box-shadow: 0 0 0 3.5px rgba(26, 92, 42, 0.08);
}

.pw-toggle {
  position: absolute;
  right: 14px;
  background: none;
  border: none;
  cursor: pointer;
  font-size: 16px;
  color: var(--text3);
  display: flex;
  align-items: center;
  padding: 4px;
}

.login-extra {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 12px;
}

.remember-label {
  display: flex;
  align-items: center;
  gap: 8px;
  color: var(--text2);
  cursor: pointer;
}

.custom-checkbox {
  accent-color: var(--hijau);
}

.forgot-link {
  color: var(--hijau);
  text-decoration: none;
  font-weight: 600;
}
.forgot-link:hover {
  text-decoration: underline;
}

.btn-login {
  width: 100%;
  padding: 13px;
  background: linear-gradient(135deg, var(--hijau2) 0%, var(--hijau) 100%);
  color: #ffffff;
  border: none;
  border-radius: var(--r8);
  font-size: 14px;
  font-weight: 700;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  transition: all 0.2s ease;
  margin-top: 6px;
  box-shadow: 0 4px 12px rgba(26, 92, 42, 0.15);
}

.btn-login:hover {
  transform: translateY(-1px);
  box-shadow: 0 6px 20px rgba(26, 92, 42, 0.25);
}

.btn-login:disabled {
  opacity: 0.7;
  cursor: not-allowed;
  transform: none;
}

.error-msg {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px;
  background: #fde8e8;
  border: 1px solid #fca5a5;
  border-radius: var(--r8);
  font-size: 12px;
  color: var(--merah);
  font-weight: 500;
}

.login-divider {
  margin: 28px 0 20px;
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 10px;
  color: var(--text3);
  text-transform: uppercase;
  letter-spacing: 0.05em;
  font-weight: 600;
}

.login-divider::before, .login-divider::after {
  content: '';
  flex: 1;
  height: 1px;
  background: var(--border);
}

.login-footer {
  display: flex;
  justify-content: space-between;
  font-size: 11px;
  color: var(--text3);
}

/* Animations */
@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}
.spin {
  animation: spin 0.8s linear infinite;
}

/* Responsive Styles */
@media (max-width: 992px) {
  .login-bg {
    display: none;
  }
  .login-panel {
    flex: 1;
    background: var(--bg);
  }
  .login-card {
    background: #ffffff;
    padding: 32px;
    border-radius: var(--r12);
    border: 1px solid var(--border);
    box-shadow: var(--shadow-md);
  }
}
</style>
