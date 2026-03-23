<template>
  <div class="customer-list">
    <div class="page-header">
      <div class="page-title">
        <h1>{{ $t('customer.title') }}</h1>
      </div>
      <el-button type="primary" @click="openCreateDialog">
        <el-icon><Plus /></el-icon>
        {{ $t('common.create') }}
      </el-button>
    </div>

    <div class="search-card glass">
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item>
          <el-input v-model="searchForm.search" :placeholder="$t('customer.customerName') + '/' + $t('customer.customerCode')" clearable>
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item>
          <el-select v-model="searchForm.status" :placeholder="$t('common.status')" clearable>
            <el-option :label="$t('common.active')" value="active" />
            <el-option :label="$t('common.inactive')" value="inactive" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadCustomers">{{ $t('common.search') }}</el-button>
          <el-button @click="resetSearch">{{ $t('common.reset') }}</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="table-card glass">
      <el-table :data="customers" v-loading="loading" :element-loading-text="$t('common.loading')" element-loading-background="rgba(10, 15, 28, 0.7)" style="width: 100%">
        <el-table-column prop="customer_code" :label="$t('customer.customerCode')" width="140" />
        <el-table-column prop="customer_name" :label="$t('customer.customerName')" />
        <el-table-column prop="industry" :label="$t('customer.industry')" width="120" />
        <el-table-column prop="region" :label="$t('customer.region')" width="100" />
        <el-table-column prop="overall_score" :label="$t('customer.healthScore')" width="150">
          <template #default="{ row }">
            <el-progress :percentage="row.overall_score || 0" :color="getHealthColor(row.overall_score)" />
          </template>
        </el-table-column>
        <el-table-column prop="status" :label="$t('common.status')" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'info'">
              {{ row.status === 'active' ? $t('common.active') : $t('common.inactive') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('common.action')" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="viewDetail(row.id)">{{ $t('common.detail') }}</el-button>
            <el-button type="primary" link @click="editCustomer(row)">{{ $t('common.edit') }}</el-button>
            <el-button type="danger" link @click="deleteCustomer(row.id)">{{ $t('common.delete') }}</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :total="pagination.total"
        :page-sizes="[10, 20, 50]"
        layout="total, sizes, prev, pager, next"
        @size-change="loadCustomers"
        @current-change="loadCustomers"
        style="margin-top: 20px; justify-content: flex-end"
      />
    </div>

    <el-dialog v-model="showDialog" :title="editingCustomer ? $t('common.edit') : $t('common.create')" width="550px">
      <el-form :model="customerForm" :rules="rules" ref="formRef" label-width="120px">
        <el-form-item :label="$t('customer.customerCode')" prop="customer_code">
          <el-input v-model="customerForm.customer_code" :disabled="!!editingCustomer" :placeholder="$t('customer.customerCode')" />
        </el-form-item>
        <el-form-item :label="$t('customer.customerName')" prop="customer_name">
          <el-input v-model="customerForm.customer_name" :placeholder="$t('customer.customerName')" />
        </el-form-item>
        <el-form-item :label="$t('customer.industry')">
          <el-input v-model="customerForm.industry" :placeholder="$t('customer.industry')" />
        </el-form-item>
        <el-form-item :label="$t('customer.scale')">
          <el-select v-model="customerForm.scale" :placeholder="$t('customer.scale')" style="width: 100%">
            <el-option :label="$t('customer.large')" value="large" />
            <el-option :label="$t('customer.medium')" value="medium" />
            <el-option :label="$t('customer.small')" value="small" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('customer.region')">
          <el-input v-model="customerForm.region" :placeholder="$t('customer.region')" />
        </el-form-item>
        <el-form-item :label="$t('customer.address')">
          <el-input v-model="customerForm.address" type="textarea" :placeholder="$t('customer.address')" />
        </el-form-item>
        <el-form-item :label="$t('customer.website')">
          <el-input v-model="customerForm.website" :placeholder="$t('customer.website')" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showDialog = false">{{ $t('common.cancel') }}</el-button>
        <el-button type="primary" @click="saveCustomer" :loading="saving">{{ $t('common.save') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Search } from '@element-plus/icons-vue'
import { customerApi } from '../../api'
import { mockCustomers } from '../../mock/data'

const { t } = useI18n()
const router = useRouter()
const loading = ref(false)
const saving = ref(false)
const showDialog = ref(false)
const customers = ref<any[]>([])
const editingCustomer = ref<any>(null)
const formRef = ref()

const searchForm = reactive({ search: '', status: '' })
const pagination = reactive({ page: 1, pageSize: 10, total: 0 })

const customerForm = reactive({
  customer_code: '',
  customer_name: '',
  industry: '',
  scale: '',
  region: '',
  address: '',
  website: ''
})

const rules = computed(() => ({
  customer_code: [{ required: true, message: t('common.required'), trigger: 'blur' }],
  customer_name: [{ required: true, message: t('common.required'), trigger: 'blur' }]
}))

const getHealthColor = (score: number) => {
  if (score > 70) return '#10b981'
  if (score > 40) return '#f59e0b'
  return '#ef4444'
}

const loadCustomers = async () => {
  loading.value = true
  try {
    const res = await customerApi.getList({
      page: pagination.page,
      page_size: pagination.pageSize,
      ...searchForm
    })
    customers.value = res.data.data || []
    pagination.total = res.data.pagination?.total || 0
  } catch (error) {
    customers.value = mockCustomers
    pagination.total = mockCustomers.length
  } finally {
    loading.value = false
  }
}

const resetSearch = () => {
  searchForm.search = ''
  searchForm.status = ''
  loadCustomers()
}

const openCreateDialog = () => {
  editingCustomer.value = null
  Object.assign(customerForm, { customer_code: '', customer_name: '', industry: '', scale: '', region: '', address: '', website: '' })
  showDialog.value = true
}

const viewDetail = (id: number) => router.push(`/customers/${id}`)

const editCustomer = (customer: any) => {
  editingCustomer.value = customer
  Object.assign(customerForm, customer)
  showDialog.value = true
}

const saveCustomer = async () => {
  await formRef.value.validate()
  saving.value = true
  try {
    if (editingCustomer.value) {
      await customerApi.update(editingCustomer.value.id, customerForm)
      ElMessage.success(t('common.success'))
    } else {
      await customerApi.create(customerForm)
      ElMessage.success(t('common.success'))
    }
    showDialog.value = false
    loadCustomers()
  } catch (error: any) {
    ElMessage.error(error.response?.data?.error || t('common.error'))
  } finally {
    saving.value = false
  }
}

const deleteCustomer = async (id: number) => {
  await ElMessageBox.confirm(t('common.confirm') + '?', t('common.delete'), { type: 'warning' })
  try {
    await customerApi.delete(id)
    ElMessage.success(t('common.success'))
    loadCustomers()
  } catch (error: any) {
    ElMessage.error(error.response?.data?.error || t('common.error'))
  }
}

onMounted(() => loadCustomers())
</script>

<style scoped>
.customer-list {
  max-width: 1400px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.page-title h1 {
  font-family: var(--font-display);
  font-size: 24px;
  font-weight: 700;
  background: linear-gradient(135deg, var(--color-accent) 0%, var(--color-accent-secondary) 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin: 0;
}

.search-card {
  padding: 20px;
  margin-bottom: 20px;
  border-radius: var(--radius-md);
  position: relative;
  overflow: hidden;
}

.search-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 2px;
  background: linear-gradient(90deg, var(--color-accent), var(--color-accent-secondary));
}

.table-card {
  padding: 20px;
  border-radius: var(--radius-md);
  position: relative;
  overflow: hidden;
}

.table-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 2px;
  background: linear-gradient(90deg, var(--color-accent), var(--color-accent-secondary));
}

.search-form {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}
</style>
