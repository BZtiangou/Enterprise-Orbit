<template>
  <div class="customer-list">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>客户列表</span>
          <el-button type="primary" @click="showDialog = true">新增客户</el-button>
        </div>
      </template>
      
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="关键词">
          <el-input v-model="searchForm.search" placeholder="客户名称/编码" clearable />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="全部" clearable>
            <el-option label="活跃" value="active" />
            <el-option label="暂停" value="inactive" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadCustomers">搜索</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="customers" v-loading="loading" style="width: 100%">
        <el-table-column prop="customer_code" label="客户编码" width="120" />
        <el-table-column prop="customer_name" label="客户名称" />
        <el-table-column prop="industry" label="行业" width="120" />
        <el-table-column prop="region" label="区域" width="100" />
        <el-table-column prop="overall_score" label="健康度" width="150">
          <template #default="{ row }">
            <el-progress :percentage="row.overall_score || 0" :color="row.overall_score > 70 ? '#67c23a' : row.overall_score > 40 ? '#e6a23c' : '#f56c6c'" />
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'info'">{{ row.status === 'active' ? '活跃' : '暂停' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200">
          <template #default="{ row }">
            <el-button type="primary" link @click="viewDetail(row.id)">详情</el-button>
            <el-button type="primary" link @click="editCustomer(row)">编辑</el-button>
            <el-button type="danger" link @click="deleteCustomer(row.id)">删除</el-button>
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
    </el-card>

    <el-dialog v-model="showDialog" :title="editingCustomer ? '编辑客户' : '新增客户'" width="500px">
      <el-form :model="customerForm" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="客户编码" prop="customer_code">
          <el-input v-model="customerForm.customer_code" :disabled="!!editingCustomer" />
        </el-form-item>
        <el-form-item label="客户名称" prop="customer_name">
          <el-input v-model="customerForm.customer_name" />
        </el-form-item>
        <el-form-item label="行业">
          <el-input v-model="customerForm.industry" />
        </el-form-item>
        <el-form-item label="规模">
          <el-select v-model="customerForm.scale" placeholder="请选择">
            <el-option label="大型企业" value="large" />
            <el-option label="中型企业" value="medium" />
            <el-option label="小型企业" value="small" />
          </el-select>
        </el-form-item>
        <el-form-item label="区域">
          <el-input v-model="customerForm.region" />
        </el-form-item>
        <el-form-item label="地址">
          <el-input v-model="customerForm.address" type="textarea" />
        </el-form-item>
        <el-form-item label="网站">
          <el-input v-model="customerForm.website" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showDialog = false">取消</el-button>
        <el-button type="primary" @click="saveCustomer" :loading="saving">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { customerApi } from '../../api'

const router = useRouter()
const loading = ref(false)
const saving = ref(false)
const showDialog = ref(false)
const customers = ref<any[]>([])
const editingCustomer = ref<any>(null)
const formRef = ref()

const searchForm = reactive({
  search: '',
  status: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const customerForm = reactive({
  customer_code: '',
  customer_name: '',
  industry: '',
  scale: '',
  region: '',
  address: '',
  website: ''
})

const rules = {
  customer_code: [{ required: true, message: '请输入客户编码', trigger: 'blur' }],
  customer_name: [{ required: true, message: '请输入客户名称', trigger: 'blur' }]
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
    console.error('Failed to load customers:', error)
  } finally {
    loading.value = false
  }
}

const viewDetail = (id: number) => {
  router.push(`/customers/${id}`)
}

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
      ElMessage.success('更新成功')
    } else {
      await customerApi.create(customerForm)
      ElMessage.success('创建成功')
    }
    showDialog.value = false
    loadCustomers()
  } catch (error: any) {
    ElMessage.error(error.response?.data?.error || '操作失败')
  } finally {
    saving.value = false
  }
}

const deleteCustomer = async (id: number) => {
  await ElMessageBox.confirm('确定要删除该客户吗？', '提示', { type: 'warning' })
  try {
    await customerApi.delete(id)
    ElMessage.success('删除成功')
    loadCustomers()
  } catch (error: any) {
    ElMessage.error(error.response?.data?.error || '删除失败')
  }
}

onMounted(() => {
  loadCustomers()
})
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.search-form {
  margin-bottom: 20px;
}
</style>
