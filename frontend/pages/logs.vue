<template>
  <div id="page-logs" style="padding: 20px; flex: 1; overflow-y: auto;">
    <div style="margin-bottom: 24px;">
      <h1 style="font-size: 24px; font-weight: 800; color: var(--text1); margin-bottom: 8px;">Sistem Log & Audit</h1>
      <p style="font-size: 13px; color: var(--text3);">Pantau seluruh aktivitas aplikasi dan riwayat koneksi ke Gateway Eksternal</p>
    </div>

    <!-- Tab Navigation -->
    <div class="tabs-bar" style="margin-bottom: 20px;">
      <button :class="['tab', { active: activeTab === 'audit' }]" @click="switchTab('audit')">
        <i class="ti ti-user-scan"></i> Audit Trail (Akses)
      </button>
      <button :class="['tab', { active: activeTab === 'payment' }]" @click="switchTab('payment')">
        <i class="ti ti-receipt"></i> Log Payment Gateway
      </button>
      <button :class="['tab', { active: activeTab === 'error' }]" @click="switchTab('error')">
        <i class="ti ti-alert-triangle"></i> Error Aplikasi
      </button>
    </div>

    <div class="card">
      <div v-if="loading" class="card-body" style="padding: 40px; text-align: center; color: var(--text3);">
        <i class="ti ti-loader-2 spin" style="font-size: 32px; color: var(--hijau); margin-bottom: 12px; display: inline-block;"></i>
        <div>Memuat log...</div>
      </div>

      <!-- AUDIT TAB -->
      <div v-else-if="activeTab === 'audit'" class="card-body" style="padding: 0; overflow-x: auto;">
        <table class="pg-table" style="width: 100%;">
          <thead>
            <tr>
              <th>Waktu</th>
              <th>Pengguna</th>
              <th>Aksi</th>
              <th>Entitas / IP</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="log in auditLogs" :key="log.id">
              <td style="font-size: 13px; color: var(--text2);">{{ formatDate(log.createdAt) }}</td>
              <td>
                <div style="font-weight: 600; font-size: 13px;">{{ log.userName || 'Sistem' }}</div>
                <div style="font-size: 11px; color: var(--text3); font-family: monospace;">{{ log.userId ? log.userId.substring(0, 8) + '...' : '-' }}</div>
              </td>
              <td>
                <span class="bdg bdg-b">{{ log.action }}</span>
              </td>
              <td>
                <div style="font-size: 13px;">{{ log.entity || '-' }} ({{ log.entityId ? log.entityId.substring(0, 8) + '...' : '-' }})</div>
                <div style="font-size: 11px; color: var(--text3);"><i class="ti ti-network"></i> {{ log.ipAddress || 'Unknown' }}</div>
              </td>
            </tr>
            <tr v-if="auditLogs.length === 0">
              <td colspan="4" style="padding: 40px; text-align: center; color: var(--text3);">
                <i class="ti ti-inbox" style="font-size: 32px; color: var(--border2); margin-bottom: 12px; display: inline-block;"></i>
                <div>Belum ada catatan aktivitas</div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- PAYMENT TAB -->
      <div v-else-if="activeTab === 'payment'" class="card-body" style="padding: 0; overflow-x: auto;">
        <table class="pg-table" style="width: 100%;">
          <thead>
            <tr>
              <th>Waktu</th>
              <th>Gateway / TxID</th>
              <th>Status HTTP</th>
              <th style="text-align: center;">Payload JSON</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="log in paymentLogs" :key="log.id">
              <td style="font-size: 13px; color: var(--text2);">{{ formatDate(log.createdAt) }}</td>
              <td>
                <div style="font-weight: 700; font-size: 13px;">{{ log.provider }}</div>
                <div style="font-size: 11px; color: var(--text3); font-family: monospace;" :title="log.endpoint">{{ log.endpoint.substring(0, 30) }}...</div>
                <div style="font-size: 11px; color: var(--text3); font-family: monospace; margin-top: 4px;">Tx: {{ log.transactionId ? log.transactionId.substring(0, 8) + '...' : '-' }}</div>
              </td>
              <td>
                <span :class="['bdg', log.statusCode >= 200 && log.statusCode < 300 ? 'bdg-g' : 'bdg-r']">
                  {{ log.statusCode || 'N/A' }}
                </span>
                <div v-if="log.errorMessage" style="font-size: 11px; color: var(--merah); margin-top: 6px; max-width: 200px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis;" :title="log.errorMessage">
                  {{ log.errorMessage }}
                </div>
              </td>
              <td style="text-align: center;">
                <button class="btn btn-ghost btn-sm" @click="viewJson(log)">
                  <i class="ti ti-brackets-contain"></i> Lihat Data
                </button>
              </td>
            </tr>
            <tr v-if="paymentLogs.length === 0">
              <td colspan="4" style="padding: 40px; text-align: center; color: var(--text3);">
                <i class="ti ti-inbox" style="font-size: 32px; color: var(--border2); margin-bottom: 12px; display: inline-block;"></i>
                <div>Belum ada log dari Gateway Pembayaran</div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- ERROR TAB -->
      <div v-else-if="activeTab === 'error'" class="card-body" style="padding: 0; overflow-x: auto;">
        <table class="pg-table" style="width: 100%;">
          <thead>
            <tr>
              <th>Waktu</th>
              <th>Service / Type</th>
              <th>Error Message</th>
              <th style="text-align: center;">Stack Trace</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="log in appErrors" :key="log.id">
              <td style="font-size: 13px; color: var(--text2);">{{ formatDate(log.createdAt) }}</td>
              <td>
                <div style="font-weight: 600; font-size: 13px;">{{ log.service || '-' }}</div>
                <div style="font-size: 11px; color: var(--text3); font-family: monospace; margin-top: 4px;">{{ log.errorType || '-' }}</div>
              </td>
              <td>
                <div style="font-size: 13px; color: var(--merah); font-weight: 500; display: -webkit-box; -webkit-line-clamp: 2; -webkit-box-orient: vertical; overflow: hidden;" :title="log.message">{{ log.message }}</div>
              </td>
              <td style="text-align: center;">
                <button class="btn btn-ghost btn-sm" style="color: var(--merah);" @click="viewStackTrace(log)">
                  <i class="ti ti-bug"></i> Detail Trace
                </button>
              </td>
            </tr>
            <tr v-if="appErrors.length === 0">
              <td colspan="4" style="padding: 40px; text-align: center; color: var(--text3);">
                <i class="ti ti-shield-check" style="font-size: 32px; color: var(--hijauSoft); margin-bottom: 12px; display: inline-block;"></i>
                <div>Sistem berjalan mulus tanpa error terdeteksi</div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Pagination Controls -->
      <div class="card-footer" style="padding: 16px; border-top: 1px solid var(--border); display: flex; justify-content: space-between; align-items: center;" v-if="!loading">
        <div style="font-size: 12px; color: var(--text3);">Halaman {{ page }}</div>
        <div style="display: flex; gap: 8px;">
          <button class="btn btn-outline btn-sm" style="width: auto;" @click="changePage(-1)" :disabled="page <= 1">Prev</button>
          <button class="btn btn-outline btn-sm" style="width: auto;" @click="changePage(1)">Next</button>
        </div>
      </div>
    </div>

    <!-- JSON Viewer Modal -->
    <div v-if="showJsonModal" class="modal-overlay" @click.self="showJsonModal = false">
      <div class="modal-card" style="max-width: 800px; width: 90%;">
        <div class="modal-header">
          <h2 class="modal-title">Detail Payload & Respons</h2>
          <button class="modal-close" @click="showJsonModal = false"><i class="ti ti-x"></i></button>
        </div>
        <div class="modal-body" style="padding: 20px; display: flex; gap: 20px; max-height: 60vh; overflow-y: auto;">
          <div style="flex: 1;">
            <div style="font-size: 12px; font-weight: 700; margin-bottom: 8px; color: var(--text1);"><i class="ti ti-arrow-right" style="color: var(--hijau);"></i> Request Payload (Dikirim)</div>
            <pre style="background: #1e293b; color: #cbd5e1; padding: 16px; border-radius: 8px; font-size: 11px; font-family: monospace; white-space: pre-wrap; height: 350px; overflow-y: auto;">{{ JSON.stringify(selectedLog?.requestPayload || {}, null, 2) }}</pre>
          </div>
          <div style="flex: 1;">
            <div style="font-size: 12px; font-weight: 700; margin-bottom: 8px; color: var(--text1);"><i class="ti ti-arrow-left" style="color: var(--biru);"></i> Response Payload (Diterima)</div>
            <pre style="background: #1e293b; color: #cbd5e1; padding: 16px; border-radius: 8px; font-size: 11px; font-family: monospace; white-space: pre-wrap; height: 350px; overflow-y: auto;">{{ JSON.stringify(selectedLog?.responsePayload || {}, null, 2) }}</pre>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn-primary" style="width: auto;" @click="showJsonModal = false">Tutup</button>
        </div>
      </div>
    </div>

    <!-- Stack Trace Modal -->
    <div v-if="showTraceModal" class="modal-overlay" @click.self="showTraceModal = false">
      <div class="modal-card" style="max-width: 800px; width: 90%;">
        <div class="modal-header" style="background: #fef2f2; border-bottom: 1px solid #fee2e2;">
          <h2 class="modal-title" style="color: var(--merah);"><i class="ti ti-alert-octagon"></i> Detail Error App</h2>
          <button class="modal-close" @click="showTraceModal = false"><i class="ti ti-x"></i></button>
        </div>
        <div class="modal-body" style="padding: 20px; background: #0f172a; max-height: 60vh; overflow-y: auto;">
          <div style="background: rgba(220, 38, 38, 0.15); border: 1px solid rgba(220, 38, 38, 0.3); padding: 12px; border-radius: 6px; color: #fca5a5; font-size: 12px; font-family: monospace; margin-bottom: 16px;">
            {{ selectedLog?.message }}
          </div>
          <div style="font-size: 10px; color: #94a3b8; text-transform: uppercase; letter-spacing: 1px; margin-bottom: 8px;">Stack Trace</div>
          <pre style="color: #cbd5e1; font-size: 11px; font-family: monospace; white-space: pre-wrap;">{{ selectedLog?.stackTrace || 'No stack trace available' }}</pre>
          
          <div v-if="selectedLog?.context" style="font-size: 10px; color: #94a3b8; text-transform: uppercase; letter-spacing: 1px; margin-top: 24px; margin-bottom: 8px;">Contextual Data</div>
          <pre v-if="selectedLog?.context" style="background: #1e293b; color: #93c5fd; padding: 12px; border-radius: 6px; font-size: 11px; font-family: monospace; white-space: pre-wrap;">{{ JSON.stringify(selectedLog.context, null, 2) }}</pre>
        </div>
        <div class="modal-footer">
          <button class="btn btn-outline" style="width: auto;" @click="showTraceModal = false">Tutup</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

definePageMeta({ title: 'Log & Audit' })

const activeTab = ref('audit') // 'audit' | 'payment' | 'error'
const loading = ref(true)
const page = ref(1)

const auditLogs = ref<any[]>([])
const paymentLogs = ref<any[]>([])
const appErrors = ref<any[]>([])

const showJsonModal = ref(false)
const showTraceModal = ref(false)
const selectedLog = ref<any>(null)

// API Config
const api = useApi()

const fetchLogs = async () => {
  loading.value = true
  try {
    const endpoint = `/admin/logs/${activeTab.value === 'error' ? 'errors' : activeTab.value === 'audit' ? 'audit' : 'payments'}?page=${page.value}`
    const res = await api.get(endpoint)
    
    if (activeTab.value === 'audit') {
      auditLogs.value = res.data || []
    } else if (activeTab.value === 'payment') {
      paymentLogs.value = res.data || []
    } else {
      appErrors.value = res.data || []
    }
  } catch (e) {
    console.error('Failed fetching logs', e)
    // If backend is not ready, mock empty array
    if (activeTab.value === 'audit') auditLogs.value = []
    if (activeTab.value === 'payment') paymentLogs.value = []
    if (activeTab.value === 'error') appErrors.value = []
  } finally {
    loading.value = false
  }
}

const switchTab = (tab: string) => {
  if (activeTab.value !== tab) {
    activeTab.value = tab
    page.value = 1
    fetchLogs()
  }
}

const changePage = (delta: number) => {
  const newPage = page.value + delta
  if (newPage >= 1) {
    page.value = newPage
    fetchLogs()
  }
}

const formatDate = (isoString: string) => {
  if (!isoString) return '-'
  const d = new Date(isoString)
  return d.toLocaleString('id-ID', { 
    year: 'numeric', month: 'short', day: '2-digit',
    hour: '2-digit', minute: '2-digit', second: '2-digit'
  })
}

const viewJson = (log: any) => {
  selectedLog.value = log
  showJsonModal.value = true
}

const viewStackTrace = (log: any) => {
  selectedLog.value = log
  showTraceModal.value = true
}

onMounted(() => {
  fetchLogs()
})
</script>

<style scoped>
.tabs-bar { display: flex; gap: 4px; border-bottom: 1px solid var(--border); background: var(--card); padding: 0 20px; flex-shrink: 0; }
.tab { padding: 12px 16px; font-size: 12px; font-weight: 600; color: var(--text2); border: none; background: none; cursor: pointer; border-bottom: 2px solid transparent; display: flex; align-items: center; gap: 6px; }
.tab.active { color: var(--hijau); border-bottom-color: var(--hijau); }
.tab-count { background: var(--hijau3); color: var(--hijau); border-radius: 10px; padding: 1px 5px; font-size: 9px; font-weight: 700; }
</style>
