<template>
  <div class="contract-list">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>合同列表</span>
          <el-button type="primary" @click="showDialog = true">新增合同</el-button>
        </div>
      </template>
      
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="关键词">
          <el-input v-model="searchForm.search" placeholder="合同编号/名称" clearable />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="全部" clearable>
            <el-option label="草稿" value="draft" />
            <el-option label="待审批" value="pending" />
            <el-option label="生效中" value="active" />
            <el-option label="已过期" value="expired" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadContracts">搜索</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="contracts" v-loading="loading" style="width: 100%">
        <el-table-column prop="contract_no" label="合同编号" width="120" />
        <el-table-column prop="contract_name" label="合同名称" />
        <el-table-column prop="customer_name" label="客户" width="150" />
        <el-table-column prop="amount" label="金额" width="120">
          <template #default="{ row }">
            ¥{{ row.amount?.toLocaleString() }}
          </template>
        </el-table-column>
        <el-table-column prop="start_date" label="开始日期" width="110" />
        <el-table-column prop="end_date" label="结束日期" width="110" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="statusColors[row.status]">{{ statusLabels[row.status] }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="250">
          <template #default="{ row }">
            <el-button type="primary" link @click="viewDetail(row.id)">详情</el-button>
            <el-button v-if="row.status === 'draft'" type="success" link @click="submitApproval(row.id)">提交审批</el-button>
            <el-button v-if="row.status === 'pending'" type="warning" link @click="showApprovalDialog(row)">审批</el-button>
            <el-button type="danger" link @click="deleteContract(row.id)">删除</el-button>
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
    </el-card>

    <el-dialog v-model="showDialog" title="新增合同" width="600px">
      <el-form :model="contractForm" :rules="rules" ref="formRef" label-width="100px">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="合同编号" prop="contract_no">
              <el-input v-model="contractForm.contract_no" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="合同名称" prop="contract_name">
              <el-input v-model="contractForm.contract_name" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="客户" prop="customer_id">
              <el-select v-model="contractForm.customer_id" placeholder="选择客户" style="width: 100%">
                <el-option v-for="c in customers" :key="c.id" :label="c.customer_name" :value="c.id" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="合同类型">
              <el-select v-model="contractForm.contract_type" placeholder="选择类型" style="width: 100%">
                <el-option label="销售合同" value="sales" />
                <el-option label="服务合同" value="service" />
                <el-option label="采购合同" value="purchase" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="金额">
              <el-input-number v-model="contractForm.amount" :min="0" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="币种">
              <el-select v-model="contractForm.currency" style="width: 100%">
                <el-option label="人民币" value="CNY" />
                <el-option label="美元" value="USD" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="开始日期">
              <el-date-picker v-model="contractForm.start_date" type="date" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="结束日期">
              <el-date-picker v-model="contractForm.end_date" type="date" style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <el-button @click="showDialog = false">取消</el-button>
        <el-button type="primary" @click="saveContract" :loading="saving">保存</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="approvalDialogVisible" title="合同审批" width="400px">
      <el-form :model="approvalForm" label-width="80px">
        <el-form-item label="审批结果">
          <el-radio-group v-model="approvalForm.status">
            <el-radio label="approved">通过</el-radio>
            <el-radio label="rejected">拒绝</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="审批意见">
          <el-input v-model="approvalForm.comment" type="textarea" :rows="3" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="approvalDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="approveContract" :loading="approving">提交</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { contractApi, customerApi } from '../../api'

const router = useRouter()
const loading = ref(false)
const saving = ref(false)
const approving = ref(false)
const showDialog = ref(false)
const approvalDialogVisible = ref(false)
const contracts = ref<any[]>([])
const customers = ref<any[]>([])
const currentContract = ref<any>(null)
const formRef = ref()

const statusColors: any = { draft: 'info', pending: 'warning', active: 'success', expired: 'danger' }
const statusLabels: any = { draft: '草稿', pending: '待审批', active: '生效中', expired: '已过期' }

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
  end_date: ''
})

const approvalForm = reactive({ status: 'approved', comment: '' })

const rules = {
  contract_no: [{ required: true, message: '请输入合同编号', trigger: 'blur' }],
  contract_name: [{ required: true, message: '请输入合同名称', trigger: 'blur' }],
  customer_id: [{ required: true, message: '请选择客户', trigger: 'change' }]
}

const loadContracts = async () => {
  loading.value = true
  try {
    const res = await contractApi.getList({ page: pagination.page, page_size: pagination.pageSize, ...searchForm })
    contracts.value = res.data.data || []
    pagination.total = res.data.pagination?.total || 0
  } catch (error) {
    console.error('Failed to load contracts:', error)
  } finally {
    loading.value = false
  }
}

const loadCustomers = async () => {
  try {
    const res = await customerApi.getList({ page: 1, page_size: 100 })
    customers.value = res.data.data || []
  } catch (error) {
    console.error('Failed to load customers:', error)
  }
}

const viewDetail = (id: number) => router.push(`/contracts/${id}`)

const saveContract = async () => {
  await formRef.value.validate()
  saving.value = true
  try {
    await contractApi.create(contractForm)
    ElMessage.success('创建成功')
    showDialog.value = false
    loadContracts()
  } catch (error: any) {
    ElMessage.error(error.response?.data?.error || '操作失败')
  } finally {
    saving.value = false
  }
}

const submitApproval = async (id: number) => {
  try {
    await contractApi.submit(id)
    ElMessage.success('已提交审批')
    loadContracts()
  } catch (error: any) {
    ElMessage.error(error.response?.data?.error || '操作失败')
  }
}

const showApprovalDialog = (contract: any) => {
  currentContract.value = contract
  approvalForm.status = 'approved'
  approvalForm.comment = ''
  approvalDialogVisible.value = true
}

const approveContract = async () => {
  if (!currentContract.value) return
  approving.value = true
  try {
    await contractApi.approve(currentContract.value.id, approvalForm)
    ElMessage.success('审批完成')
    approvalDialogVisible.value = false
    loadContracts()
  } catch (error: any) {
    ElMessage.error(error.response?.data?.error || '操作失败')
  } finally {
    approving.value = false
  }
}

const deleteContract = async (id: number) => {
  await ElMessageBox.confirm('确定要删除该合同吗？', '提示', { type: 'warning' })
  try {
    await contractApi.delete(id)
    ElMessage.success('删除成功')
    loadContracts()
  } catch (error: any) {
    ElMessage.error(error.response?.data?.error || '删除失败')
  }
}

onMounted(() => {
  loadContracts()
  loadCustomers()
})
</script>

<style scoped>
.card-header { display: flex; justify-content: space-between; align-items: center; }
.search-form { margin-bottom: 20px; }
</style>
