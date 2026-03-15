<template>
  <div class="contract-detail">
    <el-page-header @back="$router.back()" title="返回">
      <template #content>
        <span>{{ contract.contract_name }}</span>
      </template>
    </el-page-header>

    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="16">
        <el-card>
          <template #header>
            <span>合同信息</span>
          </template>
          <el-descriptions :column="2" border>
            <el-descriptions-item label="合同编号">{{ contract.contract_no }}</el-descriptions-item>
            <el-descriptions-item label="合同名称">{{ contract.contract_name }}</el-descriptions-item>
            <el-descriptions-item label="客户">{{ contract.customer_name }}</el-descriptions-item>
            <el-descriptions-item label="合同类型">{{ contract.contract_type }}</el-descriptions-item>
            <el-descriptions-item label="金额">¥{{ contract.amount?.toLocaleString() }}</el-descriptions-item>
            <el-descriptions-item label="币种">{{ contract.currency }}</el-descriptions-item>
            <el-descriptions-item label="开始日期">{{ contract.start_date }}</el-descriptions-item>
            <el-descriptions-item label="结束日期">{{ contract.end_date }}</el-descriptions-item>
            <el-descriptions-item label="状态">
              <el-tag :type="statusColors[contract.status]">{{ statusLabels[contract.status] }}</el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="审批状态">
              <el-tag :type="contract.approval_status === 'approved' ? 'success' : contract.approval_status === 'rejected' ? 'danger' : 'warning'">
                {{ contract.approval_status === 'approved' ? '已通过' : contract.approval_status === 'rejected' ? '已拒绝' : '待审批' }}
              </el-tag>
            </el-descriptions-item>
          </el-descriptions>
        </el-card>

        <el-card style="margin-top: 20px">
          <template #header>
            <div style="display: flex; justify-content: space-between; align-items: center">
              <span>履约记录</span>
              <el-button type="primary" size="small" @click="showPerformanceDialog = true">添加记录</el-button>
            </div>
          </template>
          <el-table :data="performances" style="width: 100%">
            <el-table-column prop="record_date" label="日期" width="120" />
            <el-table-column prop="milestone" label="里程碑" />
            <el-table-column prop="completion_rate" label="完成率" width="100">
              <template #default="{ row }">
                <el-progress :percentage="row.completion_rate" />
              </template>
            </el-table-column>
            <el-table-column prop="status" label="状态" width="80">
              <template #default="{ row }">
                <el-tag :type="row.status === 'completed' ? 'success' : row.status === 'delayed' ? 'danger' : 'warning'" size="small">
                  {{ row.status === 'completed' ? '完成' : row.status === 'delayed' ? '延期' : '进行中' }}
                </el-tag>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>

      <el-col :span="8">
        <el-card>
          <template #header>
            <span>审批流程</span>
          </template>
          <el-timeline>
            <el-timeline-item v-for="item in approvalFlows" :key="item.id" :timestamp="item.approval_time" placement="top">
              <el-card>
                <p><strong>审批人:</strong> {{ item.approver_name }}</p>
                <p><strong>结果:</strong> 
                  <el-tag :type="item.approval_status === 'approved' ? 'success' : 'danger'" size="small">
                    {{ item.approval_status === 'approved' ? '通过' : '拒绝' }}
                  </el-tag>
                </p>
                <p v-if="item.comment"><strong>意见:</strong> {{ item.comment }}</p>
              </el-card>
            </el-timeline-item>
          </el-timeline>
        </el-card>
      </el-col>
    </el-row>

    <el-dialog v-model="showPerformanceDialog" title="添加履约记录" width="400px">
      <el-form :model="performanceForm" label-width="80px">
        <el-form-item label="日期">
          <el-date-picker v-model="performanceForm.record_date" type="date" style="width: 100%" />
        </el-form-item>
        <el-form-item label="里程碑">
          <el-input v-model="performanceForm.milestone" />
        </el-form-item>
        <el-form-item label="完成率">
          <el-slider v-model="performanceForm.completion_rate" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="performanceForm.status" style="width: 100%">
            <el-option label="进行中" value="in_progress" />
            <el-option label="已完成" value="completed" />
            <el-option label="延期" value="delayed" />
          </el-select>
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="performanceForm.remarks" type="textarea" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showPerformanceDialog = false">取消</el-button>
        <el-button type="primary" @click="addPerformance" :loading="saving">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { contractApi } from '../../api'

const route = useRoute()
const contractId = Number(route.params.id)

const statusColors: any = { draft: 'info', pending: 'warning', active: 'success', expired: 'danger' }
const statusLabels: any = { draft: '草稿', pending: '待审批', active: '生效中', expired: '已过期' }

const contract = ref<any>({})
const performances = ref<any[]>([])
const approvalFlows = ref<any[]>([])
const showPerformanceDialog = ref(false)
const saving = ref(false)

const performanceForm = reactive({
  record_date: '',
  milestone: '',
  completion_rate: 50,
  status: 'in_progress',
  remarks: ''
})

const loadContract = async () => {
  try {
    const res = await contractApi.getDetail(contractId)
    contract.value = res.data.data
  } catch (error) {
    console.error('Failed to load contract:', error)
  }
}

const loadPerformances = async () => {
  try {
    const res = await contractApi.getPerformance(contractId)
    performances.value = res.data.data || []
  } catch (error) {
    console.error('Failed to load performances:', error)
  }
}

const addPerformance = async () => {
  saving.value = true
  try {
    await contractApi.createPerformance(contractId, performanceForm)
    ElMessage.success('添加成功')
    showPerformanceDialog.value = false
    loadPerformances()
  } catch (error: any) {
    ElMessage.error(error.response?.data?.error || '操作失败')
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  loadContract()
  loadPerformances()
})
</script>
