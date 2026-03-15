<template>
  <div class="interaction-list">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>互动日志</span>
          <el-button type="primary" @click="showDialog = true">记录互动</el-button>
        </div>
      </template>
      
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="客户">
          <el-select v-model="searchForm.customer_id" placeholder="全部" clearable style="width: 200px">
            <el-option v-for="c in customers" :key="c.id" :label="c.customer_name" :value="c.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="类型">
          <el-select v-model="searchForm.type" placeholder="全部" clearable>
            <el-option label="会议" value="meeting" />
            <el-option label="电话" value="phone" />
            <el-option label="邮件" value="email" />
            <el-option label="拜访" value="visit" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadInteractions">搜索</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="interactions" v-loading="loading" style="width: 100%">
        <el-table-column prop="interaction_date" label="日期" width="110" />
        <el-table-column prop="customer_name" label="客户" width="150" />
        <el-table-column prop="interaction_type" label="类型" width="80">
          <template #default="{ row }">
            <el-tag size="small">{{ typeLabels[row.interaction_type] || row.interaction_type }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="topic" label="主题" />
        <el-table-column prop="outcome" label="结果" width="100">
          <template #default="{ row }">
            <el-tag v-if="row.outcome" :type="row.outcome === 'positive' ? 'success' : row.outcome === 'negative' ? 'danger' : 'info'" size="small">
              {{ row.outcome }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="importance_level" label="重要性" width="100">
          <template #default="{ row }">
            <el-rate v-model="row.importance_level" disabled />
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150">
          <template #default="{ row }">
            <el-button type="primary" link @click="editInteraction(row)">编辑</el-button>
            <el-button type="danger" link @click="deleteInteraction(row.id)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :total="pagination.total"
        :page-sizes="[10, 20, 50]"
        layout="total, sizes, prev, pager, next"
        @size-change="loadInteractions"
        @current-change="loadInteractions"
        style="margin-top: 20px; justify-content: flex-end"
      />
    </el-card>

    <el-dialog v-model="showDialog" :title="editingInteraction ? '编辑互动' : '记录互动'" width="600px">
      <el-form :model="interactionForm" :rules="rules" ref="formRef" label-width="100px">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="客户" prop="customer_id">
              <el-select v-model="interactionForm.customer_id" placeholder="选择客户" style="width: 100%">
                <el-option v-for="c in customers" :key="c.id" :label="c.customer_name" :value="c.id" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="互动类型" prop="interaction_type">
              <el-select v-model="interactionForm.interaction_type" style="width: 100%">
                <el-option label="会议" value="meeting" />
                <el-option label="电话" value="phone" />
                <el-option label="邮件" value="email" />
                <el-option label="拜访" value="visit" />
                <el-option label="其他" value="other" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="互动日期">
              <el-date-picker v-model="interactionForm.interaction_date" type="datetime" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="重要性">
              <el-rate v-model="interactionForm.importance_level" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="主题" prop="topic">
          <el-input v-model="interactionForm.topic" />
        </el-form-item>
        <el-form-item label="内容摘要">
          <el-input v-model="interactionForm.content_summary" type="textarea" :rows="3" />
        </el-form-item>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="结果">
              <el-select v-model="interactionForm.outcome" style="width: 100%">
                <el-option label="积极" value="positive" />
                <el-option label="中性" value="neutral" />
                <el-option label="消极" value="negative" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="下一步行动">
              <el-input v-model="interactionForm.next_action" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <el-button @click="showDialog = false">取消</el-button>
        <el-button type="primary" @click="saveInteraction" :loading="saving">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { interactionApi, customerApi } from '../../api'

const loading = ref(false)
const saving = ref(false)
const showDialog = ref(false)
const interactions = ref<any[]>([])
const customers = ref<any[]>([])
const editingInteraction = ref<any>(null)
const formRef = ref()

const typeLabels: any = { meeting: '会议', phone: '电话', email: '邮件', visit: '拜访', other: '其他' }

const searchForm = reactive({ customer_id: '', type: '' })
const pagination = reactive({ page: 1, pageSize: 10, total: 0 })

const interactionForm = reactive({
  customer_id: null as number | null,
  interaction_type: 'meeting',
  interaction_date: '',
  topic: '',
  content_summary: '',
  outcome: '',
  next_action: '',
  importance_level: 3
})

const rules = {
  customer_id: [{ required: true, message: '请选择客户', trigger: 'change' }],
  interaction_type: [{ required: true, message: '请选择类型', trigger: 'change' }],
  topic: [{ required: true, message: '请输入主题', trigger: 'blur' }]
}

const loadInteractions = async () => {
  loading.value = true
  try {
    const res = await interactionApi.getList({ page: pagination.page, page_size: pagination.pageSize, ...searchForm })
    interactions.value = res.data.data || []
    pagination.total = res.data.pagination?.total || 0
  } catch (error) {
    console.error('Failed to load interactions:', error)
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

const editInteraction = (interaction: any) => {
  editingInteraction.value = interaction
  Object.assign(interactionForm, interaction)
  showDialog.value = true
}

const saveInteraction = async () => {
  await formRef.value.validate()
  saving.value = true
  try {
    if (editingInteraction.value) {
      await interactionApi.update(editingInteraction.value.id, interactionForm)
      ElMessage.success('更新成功')
    } else {
      await interactionApi.create(interactionForm)
      ElMessage.success('创建成功')
    }
    showDialog.value = false
    loadInteractions()
  } catch (error: any) {
    ElMessage.error(error.response?.data?.error || '操作失败')
  } finally {
    saving.value = false
  }
}

const deleteInteraction = async (id: number) => {
  await ElMessageBox.confirm('确定要删除该互动记录吗？', '提示', { type: 'warning' })
  try {
    await interactionApi.delete(id)
    ElMessage.success('删除成功')
    loadInteractions()
  } catch (error: any) {
    ElMessage.error(error.response?.data?.error || '删除失败')
  }
}

onMounted(() => {
  loadInteractions()
  loadCustomers()
})
</script>

<style scoped>
.card-header { display: flex; justify-content: space-between; align-items: center; }
.search-form { margin-bottom: 20px; }
</style>
