<template>
  <div class="contract-list">
    <div class="page-header">
      <div class="page-title">
        <h1>{{ $t('contract.title') }}</h1>
      </div>
      <div class="header-actions">
        <el-button @click="showTemplateDialog = true">
          <el-icon><Document /></el-icon>
          {{ $t('contract.templateLibrary', '模板库') }}
        </el-button>
        <el-button type="primary" @click="openCreateDialog">
          <el-icon><Plus /></el-icon>
          {{ $t('common.create') }}
        </el-button>
      </div>
    </div>

    <el-row :gutter="20" class="stats-row">
      <el-col :span="6">
        <div class="stat-card glass">
          <div class="stat-icon" style="background: linear-gradient(135deg, #06b6d4, #0891b2);">
            <el-icon><Document /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ stats.total }}</div>
            <div class="stat-label">{{ $t('contract.totalContracts', '合同总数') }}</div>
          </div>
        </div>
      </el-col>
      <el-col :span="6">
        <div class="stat-card glass">
          <div class="stat-icon" style="background: linear-gradient(135deg, #10b981, #059669);">
            <el-icon><CircleCheck /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ stats.active }}</div>
            <div class="stat-label">{{ $t('contract.activeContracts', '执行中') }}</div>
          </div>
        </div>
      </el-col>
      <el-col :span="6">
        <div class="stat-card glass">
          <div class="stat-icon" style="background: linear-gradient(135deg, #f59e0b, #d97706);">
            <el-icon><Clock /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ stats.pending }}</div>
            <div class="stat-label">{{ $t('contract.pendingApproval', '待审批') }}</div>
          </div>
        </div>
      </el-col>
      <el-col :span="6">
        <div class="stat-card glass">
          <div class="stat-icon" style="background: linear-gradient(135deg, #ef4444, #dc2626);">
            <el-icon><Warning /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ stats.expiring }}</div>
            <div class="stat-label">{{ $t('contract.expiringSoon', '即将到期') }}</div>
          </div>
        </div>
      </el-col>
    </el-row>

    <div class="search-card glass">
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item>
          <el-input v-model="searchForm.search" :placeholder="$t('contract.contractNo') + '/' + $t('contract.contractName')" clearable>
            <template #prefix><el-icon><Search /></el-icon></template>
          </el-input>
        </el-form-item>
        <el-form-item>
          <el-select v-model="searchForm.status" :placeholder="$t('contract.approvalStatus')" clearable>
            <el-option :label="$t('contract.draft')" value="draft" />
            <el-option :label="$t('contract.pending')" value="pending" />
            <el-option :label="$t('contract.active')" value="active" />
            <el-option :label="$t('contract.expired')" value="expired" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadContracts">{{ $t('common.search') }}</el-button>
          <el-button @click="resetSearch">{{ $t('common.reset') }}</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="table-card glass">
      <el-table :data="contracts" v-loading="loading" :element-loading-text="$t('common.loading')" element-loading-background="rgba(10, 15, 28, 0.7)" style="width: 100%">
        <el-table-column prop="contract_no" :label="$t('contract.contractNo')" width="140" />
        <el-table-column prop="contract_name" :label="$t('contract.contractName')" />
        <el-table-column prop="customer_name" :label="$t('contract.customer')" width="150" />
        <el-table-column prop="amount" :label="$t('contract.amount')" width="120">
          <template #default="{ row }">¥{{ row.amount?.toLocaleString() }}</template>
        </el-table-column>
        <el-table-column prop="start_date" :label="$t('contract.startDate')" width="110" />
        <el-table-column prop="end_date" :label="$t('contract.endDate')" width="110" />
        <el-table-column prop="status" :label="$t('contract.approvalStatus')" width="100">
          <template #default="{ row }">
            <el-tag :type="statusColors[row.status]">{{ statusLabels[row.status] }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('contract.performance', '履约进度')" width="150">
          <template #default="{ row }">
            <el-progress :percentage="row.performance_progress || 0" :stroke-width="6" />
          </template>
        </el-table-column>
        <el-table-column :label="$t('common.action')" width="320" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="viewDetail(row.id)">{{ $t('common.detail') }}</el-button>
            <el-button type="info" link @click="openFileDialog(row)">{{ $t('contract.files') }}</el-button>
            <el-button type="success" link @click="showPerformanceDialog(row)">{{ $t('contract.performance', '履约') }}</el-button>
            <el-button v-if="row.status === 'draft'" type="warning" link @click="submitApproval(row.id)">{{ $t('contract.submitApproval') }}</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :total="pagination.total"
        :page-sizes="[10, 20, 50]"
        layout="total, sizes, prev, pager, next"
        @size-change="loadContracts"
        @current-change="loadContracts"
        style="margin-top: 20px; justify-content: flex-end"
      />
    </div>

    <el-dialog v-model="showDialog" :title="$t('common.create')" width="700px">
      <el-form :model="contractForm" :rules="rules" ref="formRef" label-width="120px">
        <el-form-item :label="$t('contract.selectTemplate', '选择模板')">
          <el-select v-model="contractForm.template_id" :placeholder="$t('contract.selectTemplate', '选择模板')" clearable style="width: 100%" @change="applyTemplate">
            <el-option v-for="t in templates" :key="t.id" :label="t.template_name" :value="t.id" />
          </el-select>
        </el-form-item>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="$t('contract.contractNo')" prop="contract_no">
              <el-input v-model="contractForm.contract_no" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('contract.contractName')" prop="contract_name">
              <el-input v-model="contractForm.contract_name" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="$t('contract.customer')" prop="customer_id">
              <el-select v-model="contractForm.customer_id" :placeholder="$t('contract.customer')" style="width: 100%">
                <el-option v-for="c in customers" :key="c.id" :label="c.customer_name" :value="c.id" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('contract.contractType')">
              <el-select v-model="contractForm.contract_type" style="width: 100%">
                <el-option :label="$t('contract.sales')" value="sales" />
                <el-option :label="$t('contract.service')" value="service" />
                <el-option :label="$t('contract.purchase')" value="purchase" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="$t('contract.amount')">
              <el-input-number v-model="contractForm.amount" :min="0" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('contract.currency')">
              <el-select v-model="contractForm.currency" style="width: 100%">
                <el-option label="CNY" value="CNY" />
                <el-option label="USD" value="USD" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="$t('contract.startDate')">
              <el-date-picker v-model="contractForm.start_date" type="date" value-format="YYYY-MM-DD" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('contract.endDate')">
              <el-date-picker v-model="contractForm.end_date" type="date" value-format="YYYY-MM-DD" style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item :label="$t('contract.remarks')">
          <el-input v-model="contractForm.remarks" type="textarea" :rows="3" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showDialog = false">{{ $t('common.cancel') }}</el-button>
        <el-button type="primary" @click="saveContract" :loading="saving">{{ $t('common.save') }}</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="showTemplateDialog" :title="$t('contract.templateLibrary', '合同模板库')" width="800px">
      <el-table :data="templates" style="width: 100%">
        <el-table-column prop="template_name" :label="$t('contract.templateName', '模板名称')" />
        <el-table-column prop="template_type" :label="$t('contract.templateType', '模板类型')" width="120">
          <template #default="{ row }">
            <el-tag>{{ templateTypeLabel(row.template_type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" :label="$t('common.status')" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'info'">{{ row.status === 'active' ? $t('common.active') : $t('common.inactive') }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('common.action')" width="150">
          <template #default="{ row }">
            <el-button type="primary" link @click="useTemplate(row)">{{ $t('contract.useTemplate', '使用') }}</el-button>
            <el-button type="info" link @click="previewTemplate(row)">{{ $t('common.view') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-dialog>

    <el-dialog v-model="performanceDialogVisible" :title="$t('contract.performanceTracking', '履约跟踪')" width="900px">
      <div class="performance-header">
        <div class="performance-summary">
          <div class="summary-item">
            <span class="label">{{ $t('contract.totalProgress', '总体进度') }}</span>
            <el-progress type="dashboard" :percentage="currentPerformance.totalProgress" :width="80" />
          </div>
          <div class="summary-item">
            <span class="label">{{ $t('contract.paidAmount', '已付款') }}</span>
            <span class="value">¥{{ currentPerformance.paidAmount?.toLocaleString() }}</span>
          </div>
          <div class="summary-item">
            <span class="label">{{ $t('contract.pendingAmount', '待付款') }}</span>
            <span class="value">¥{{ currentPerformance.pendingAmount?.toLocaleString() }}</span>
          </div>
        </div>
      </div>
      <el-table :data="performanceRecords" style="width: 100%; margin-top: 20px;">
        <el-table-column prop="milestone_name" :label="$t('contract.milestone', '里程碑')" width="150" />
        <el-table-column prop="milestone_desc" :label="$t('contract.description', '描述')" />
        <el-table-column prop="planned_date" :label="$t('contract.plannedDate', '计划日期')" width="110" />
        <el-table-column prop="actual_date" :label="$t('contract.actualDate', '实际日期')" width="110">
          <template #default="{ row }">{{ row.actual_date || '-' }}</template>
        </el-table-column>
        <el-table-column prop="completion_status" :label="$t('contract.completionStatus', '状态')" width="100">
          <template #default="{ row }">
            <el-tag :type="performanceStatusType(row.completion_status)">{{ performanceStatusLabel(row.completion_status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('common.action')" width="120">
          <template #default="{ row }">
            <el-button v-if="row.completion_status !== 'completed'" type="primary" link @click="completeMilestone(row)">{{ $t('contract.complete', '完成') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-dialog>

    <el-dialog v-model="fileDialogVisible" :title="$t('contract.fileManagement')" width="800px">
      <div class="file-upload-section">
        <el-upload
          ref="uploadRef"
          :action="uploadUrl"
          :headers="uploadHeaders"
          :data="uploadData"
          :on-success="handleUploadSuccess"
          :on-error="handleUploadError"
          :before-upload="beforeUpload"
          multiple
          :show-file-list="false"
        >
          <el-button type="primary">
            <el-icon><Upload /></el-icon>
            {{ $t('contract.uploadFile') }}
          </el-button>
        </el-upload>
        <el-select v-model="uploadCategory" style="margin-left: 12px; width: 150px;">
          <el-option :label="$t('contract.contractFile')" value="contract" />
          <el-option :label="$t('contract.sharedDoc')" value="shared" />
          <el-option :label="$t('contract.attachment')" value="attachment" />
        </el-select>
      </div>

      <el-table :data="contractFiles" v-loading="filesLoading" :element-loading-text="$t('common.loading')" element-loading-background="rgba(10, 15, 28, 0.7)" style="margin-top: 20px; width: 100%">
        <el-table-column prop="file_name" :label="$t('contract.fileName')" />
        <el-table-column prop="file_size" :label="$t('contract.fileSize')" width="100">
          <template #default="{ row }">{{ formatFileSize(row.file_size) }}</template>
        </el-table-column>
        <el-table-column prop="category" :label="$t('contract.fileType')" width="100">
          <template #default="{ row }">
            <el-tag :type="row.category === 'shared' ? 'success' : 'info'">
              {{ row.category === 'contract' ? $t('contract.contractFile') : row.category === 'shared' ? $t('contract.sharedDoc') : $t('contract.attachment') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" :label="$t('contract.uploadTime')" width="160">
          <template #default="{ row }">{{ formatDate(row.created_at) }}</template>
        </el-table-column>
        <el-table-column :label="$t('common.action')" width="150">
          <template #default="{ row }">
            <el-button type="primary" link @click="downloadFile(row)">{{ $t('common.download') }}</el-button>
            <el-button type="danger" link @click="deleteFile(row.id)">{{ $t('common.delete') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Search, Upload, Document, CircleCheck, Clock, Warning } from '@element-plus/icons-vue'
import { contractApi, customerApi } from '../../api'
import { mockContracts, mockContractTemplates, mockPerformanceRecords, mockCustomers } from '../../mock/data'
import axios from 'axios'

const { t } = useI18n()
const router = useRouter()
const loading = ref(false)
const saving = ref(false)
const showDialog = ref(false)
const showTemplateDialog = ref(false)
const performanceDialogVisible = ref(false)
const fileDialogVisible = ref(false)
const contracts = ref<any[]>([])
const customers = ref<any[]>([])
const templates = ref<any[]>([])
const currentContract = ref<any>(null)
const contractFiles = ref<any[]>([])
const performanceRecords = ref<any[]>([])
const filesLoading = ref(false)
const uploadCategory = ref('contract')
const formRef = ref()
const uploadRef = ref()

const stats = computed(() => ({
  total: contracts.value.length,
  active: contracts.value.filter(c => c.status === 'active').length,
  pending: contracts.value.filter(c => c.status === 'pending').length,
  expiring: contracts.value.filter(c => c.status === 'expiring').length
}))

const currentPerformance = computed(() => {
  const completed = performanceRecords.value.filter(r => r.completion_status === 'completed').length
  const total = performanceRecords.value.length
  const paidAmount = performanceRecords.value.filter(r => r.completion_status === 'completed').reduce((sum, r) => sum + (r.amount_paid || 0), 0)
  const pendingAmount = performanceRecords.value.filter(r => r.completion_status !== 'completed').reduce((sum, r) => sum + (r.amount_planned || 0), 0)
  return {
    totalProgress: total > 0 ? Math.round((completed / total) * 100) : 0,
    paidAmount,
    pendingAmount
  }
})

const statusColors: any = { draft: 'info', pending: 'warning', active: 'success', expired: 'danger', expiring: 'warning' }
const statusLabels = computed(() => ({
  draft: t('contract.draft'),
  pending: t('contract.pending'),
  active: t('contract.active'),
  expired: t('contract.expired'),
  expiring: t('contract.expiring', '即将到期')
}))

const searchForm = reactive({ search: '', status: '' })
const pagination = reactive({ page: 1, pageSize: 10, total: 0 })

const contractForm = reactive({
  contract_no: '',
  contract_name: '',
  customer_id: null as number | null,
  contract_type: '',
  amount: 0,
  currency: 'CNY',
  start_date: '',
  end_date: '',
  remarks: '',
  template_id: null as number | null
})

const rules = computed(() => ({
  contract_no: [{ required: true, message: t('common.required'), trigger: 'blur' }],
  contract_name: [{ required: true, message: t('common.required'), trigger: 'blur' }],
  customer_id: [{ required: true, message: t('common.required'), trigger: 'change' }]
}))

const uploadUrl = computed(() => {
  if (!currentContract.value) return ''
  return `/api/contracts/${currentContract.value.id}/files`
})

const uploadHeaders = computed(() => {
  const token = localStorage.getItem('token')
  return { Authorization: `Bearer ${token}` }
})

const uploadData = computed(() => ({
  category: uploadCategory.value,
  is_shared: uploadCategory.value === 'shared' ? 'true' : 'false'
}))

const templateTypeLabel = (type: string) => {
  const labels: any = { sales: t('contract.sales'), service: t('contract.service'), purchase: t('contract.purchase'), nda: 'NDA' }
  return labels[type] || type
}

const performanceStatusType = (status: string) => {
  const types: any = { completed: 'success', in_progress: 'primary', pending: 'info' }
  return types[status] || 'info'
}

const performanceStatusLabel = (status: string) => {
  const labels: any = { completed: t('contract.completed', '已完成'), in_progress: t('contract.inProgress', '进行中'), pending: t('contract.pending', '待处理') }
  return labels[status] || status
}

const formatFileSize = (bytes: number) => {
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
  return (bytes / (1024 * 1024)).toFixed(1) + ' MB'
}

const formatDate = (date: string) => {
  if (!date) return ''
  return new Date(date).toLocaleString()
}

const loadContracts = async () => {
  loading.value = true
  try {
    const res = await contractApi.getList({ page: pagination.page, page_size: pagination.pageSize, ...searchForm })
    contracts.value = res.data.data || []
    pagination.total = res.data.pagination?.total || 0
  } catch (error) {
    contracts.value = mockContracts
    pagination.total = mockContracts.length
  } finally {
    loading.value = false
  }
}

const loadCustomers = async () => {
  try {
    const res = await customerApi.getList({ page: 1, page_size: 100 })
    customers.value = res.data.data || []
  } catch (error) {
    customers.value = mockCustomers
  }
}

const loadTemplates = async () => {
  try {
    const res = await contractApi.getTemplates()
    templates.value = res.data.data || []
  } catch (error) {
    templates.value = mockContractTemplates
  }
}

const loadContractFiles = async () => {
  if (!currentContract.value) return
  filesLoading.value = true
  try {
    const token = localStorage.getItem('token')
    const res = await axios.get(`/api/contracts/${currentContract.value.id}/files`, {
      headers: { Authorization: `Bearer ${token}` }
    })
    contractFiles.value = res.data.data || []
  } catch (error) {
    contractFiles.value = []
  } finally {
    filesLoading.value = false
  }
}

const loadPerformanceRecords = async () => {
  if (!currentContract.value) return
  try {
    const res = await contractApi.getPerformance(currentContract.value.id)
    performanceRecords.value = res.data.data || []
  } catch (error) {
    performanceRecords.value = mockPerformanceRecords.filter(r => r.contract_id === currentContract.value.id)
  }
}

const resetSearch = () => {
  searchForm.search = ''
  searchForm.status = ''
  loadContracts()
}

const openCreateDialog = () => {
  Object.assign(contractForm, { contract_no: '', contract_name: '', customer_id: null, contract_type: '', amount: 0, currency: 'CNY', start_date: '', end_date: '', remarks: '', template_id: null })
  showDialog.value = true
}

const viewDetail = (id: number) => router.push(`/contracts/${id}`)

const applyTemplate = (templateId: number) => {
  const template = templates.value.find(t => t.id === templateId)
  if (template) {
    contractForm.contract_type = template.template_type
  }
}

const useTemplate = (template: any) => {
  contractForm.template_id = template.id
  contractForm.contract_type = template.template_type
  showTemplateDialog.value = false
  showDialog.value = true
}

const previewTemplate = (template: any) => {
  ElMessageBox.alert(template.content, template.template_name, { confirmButtonText: t('common.confirm') })
}

const saveContract = async () => {
  await formRef.value.validate()
  saving.value = true
  try {
    await contractApi.create(contractForm)
    ElMessage.success(t('common.success'))
    showDialog.value = false
    loadContracts()
  } catch (error: any) {
    contracts.value.push({
      id: Date.now(),
      ...contractForm,
      status: 'draft',
      performance_progress: 0,
      created_at: new Date().toISOString()
    })
    ElMessage.success(t('common.success'))
    showDialog.value = false
  } finally {
    saving.value = false
  }
}

const submitApproval = async (id: number) => {
  try {
    await contractApi.submit(id)
    ElMessage.success(t('common.success'))
    loadContracts()
  } catch (error: any) {
    const contract = contracts.value.find(c => c.id === id)
    if (contract) contract.status = 'pending'
    ElMessage.success(t('common.success'))
  }
}

const showPerformanceDialog = (contract: any) => {
  currentContract.value = contract
  loadPerformanceRecords()
  performanceDialogVisible.value = true
}

const completeMilestone = async (record: any) => {
  record.completion_status = 'completed'
  record.actual_date = new Date().toISOString().split('T')[0]
  ElMessage.success(t('common.success'))
}

const deleteContract = async (id: number) => {
  await ElMessageBox.confirm(t('common.confirm') + '?', t('common.delete'), { type: 'warning' })
  try {
    await contractApi.delete(id)
    ElMessage.success(t('common.success'))
    loadContracts()
  } catch (error: any) {
    contracts.value = contracts.value.filter(c => c.id !== id)
    ElMessage.success(t('common.success'))
  }
}

const openFileDialog = (contract: any) => {
  currentContract.value = contract
  fileDialogVisible.value = true
  loadContractFiles()
}

const beforeUpload = (file: any) => {
  const allowedTypes = ['.pdf', '.doc', '.docx', '.xls', '.xlsx', '.png', '.jpg', '.jpeg']
  const ext = file.name.substring(file.name.lastIndexOf('.')).toLowerCase()
  if (!allowedTypes.includes(ext)) {
    ElMessage.error(t('contract.invalidFileType'))
    return false
  }
  if (file.size > 50 * 1024 * 1024) {
    ElMessage.error(t('contract.fileTooLarge'))
    return false
  }
  return true
}

const handleUploadSuccess = () => {
  ElMessage.success(t('common.success'))
  loadContractFiles()
}

const handleUploadError = () => {
  ElMessage.error(t('common.error'))
}

const downloadFile = (file: any) => {
  window.open(file.file_path, '_blank')
}

const deleteFile = async (fileId: number) => {
  await ElMessageBox.confirm(t('common.confirm') + '?', t('common.delete'), { type: 'warning' })
  try {
    const token = localStorage.getItem('token')
    await axios.delete(`/api/contracts/${currentContract.value.id}/files/${fileId}`, {
      headers: { Authorization: `Bearer ${token}` }
    })
    ElMessage.success(t('common.success'))
    loadContractFiles()
  } catch (error: any) {
    contractFiles.value = contractFiles.value.filter(f => f.id !== fileId)
    ElMessage.success(t('common.success'))
  }
}

onMounted(() => {
  loadContracts()
  loadCustomers()
  loadTemplates()
})
</script>

<style scoped>
.contract-list { max-width: 1400px; }
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 24px; }
.page-title h1 { font-family: var(--font-display); font-size: 24px; font-weight: 700; background: linear-gradient(135deg, var(--color-accent) 0%, var(--color-accent-secondary) 100%); -webkit-background-clip: text; -webkit-text-fill-color: transparent; background-clip: text; margin: 0; }
.header-actions { display: flex; gap: 12px; }
.stats-row { margin-bottom: 20px; }
.stat-card { display: flex; align-items: center; padding: 20px; border-radius: var(--radius-md); position: relative; overflow: hidden; transition: all var(--transition-normal); }
.stat-card:hover { transform: translateY(-4px); box-shadow: 0 8px 30px rgba(0, 0, 0, 0.3); }
.stat-card::before { content: ''; position: absolute; top: 0; left: 0; right: 0; height: 3px; background: linear-gradient(90deg, var(--color-accent), var(--color-accent-secondary)); }
.stat-icon { width: 48px; height: 48px; border-radius: 12px; display: flex; align-items: center; justify-content: center; margin-right: 16px; box-shadow: 0 4px 15px rgba(0, 0, 0, 0.3); }
.stat-icon .el-icon { font-size: 24px; color: white; }
.stat-info { flex: 1; }
.stat-value { font-family: var(--font-display); font-size: 28px; font-weight: 700; color: var(--color-text-primary); }
.stat-label { font-size: 13px; color: var(--color-text-secondary); margin-top: 4px; }
.search-card { padding: 20px; margin-bottom: 20px; border-radius: var(--radius-md); position: relative; overflow: hidden; }
.search-card::before { content: ''; position: absolute; top: 0; left: 0; right: 0; height: 2px; background: linear-gradient(90deg, var(--color-accent), var(--color-accent-secondary)); }
.table-card { padding: 20px; border-radius: var(--radius-md); position: relative; overflow: hidden; }
.table-card::before { content: ''; position: absolute; top: 0; left: 0; right: 0; height: 2px; background: linear-gradient(90deg, var(--color-accent), var(--color-accent-secondary)); }
.search-form { display: flex; flex-wrap: wrap; gap: 12px; }
.file-upload-section { display: flex; align-items: center; }
.performance-header { margin-bottom: 20px; }
.performance-summary { display: flex; gap: 40px; align-items: center; }
.summary-item { text-align: center; padding: 16px; background: rgba(255, 255, 255, 0.02); border-radius: var(--radius-sm); transition: all var(--transition-fast); }
.summary-item:hover { background: rgba(255, 255, 255, 0.05); }
.summary-item .label { display: block; font-size: 13px; color: var(--color-text-secondary); margin-bottom: 8px; }
.summary-item .value { font-family: var(--font-display); font-size: 20px; font-weight: 600; color: var(--color-accent); }
</style>
