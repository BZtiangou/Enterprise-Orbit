<template>
  <div class="contract-detail">
    <el-page-header @back="$router.back()" :title="$t('common.view')">
      <template #content>
        <span>{{ contract.contract_name }}</span>
      </template>
    </el-page-header>

    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="16">
        <el-card>
          <template #header>
            <span>{{ $t('contract.contractInfo') }}</span>
          </template>
          <el-descriptions :column="2" border>
            <el-descriptions-item :label="$t('contract.contractNo')">{{ contract.contract_no }}</el-descriptions-item>
            <el-descriptions-item :label="$t('contract.contractName')">{{ contract.contract_name }}</el-descriptions-item>
            <el-descriptions-item :label="$t('contract.customer')">{{ contract.customer_name }}</el-descriptions-item>
            <el-descriptions-item :label="$t('contract.contractType')">{{ contract.contract_type }}</el-descriptions-item>
            <el-descriptions-item :label="$t('contract.amount')">¥{{ contract.amount?.toLocaleString() }}</el-descriptions-item>
            <el-descriptions-item :label="$t('contract.currency')">{{ contract.currency }}</el-descriptions-item>
            <el-descriptions-item :label="$t('contract.startDate')">{{ contract.start_date }}</el-descriptions-item>
            <el-descriptions-item :label="$t('contract.endDate')">{{ contract.end_date }}</el-descriptions-item>
            <el-descriptions-item :label="$t('common.status')">
              <el-tag :type="statusColors[contract.status]">{{ statusLabels[contract.status as keyof typeof statusLabels] }}</el-tag>
            </el-descriptions-item>
            <el-descriptions-item :label="$t('contract.approvalStatus')">
              <el-tag :type="contract.approval_status === 'approved' ? 'success' : contract.approval_status === 'rejected' ? 'danger' : 'warning'">
                {{ contract.approval_status === 'approved' ? $t('contract.approved') : contract.approval_status === 'rejected' ? $t('contract.rejected') : $t('contract.pending') }}
              </el-tag>
            </el-descriptions-item>
          </el-descriptions>
        </el-card>

        <el-card style="margin-top: 20px">
          <template #header>
            <div style="display: flex; justify-content: space-between; align-items: center">
              <span>{{ $t('contract.performanceRecords') }}</span>
              <el-button type="primary" size="small" @click="showPerformanceDialog = true">{{ $t('contract.addRecord') }}</el-button>
            </div>
          </template>
          <el-table :data="performances" style="width: 100%">
            <el-table-column prop="milestone_name" :label="$t('contract.milestone')" width="150" />
            <el-table-column prop="milestone_desc" :label="$t('contract.description')" />
            <el-table-column prop="planned_date" :label="$t('contract.plannedDate')" width="110" />
            <el-table-column prop="actual_date" :label="$t('contract.actualDate')" width="110">
              <template #default="{ row }">{{ row.actual_date || '-' }}</template>
            </el-table-column>
            <el-table-column prop="completion_status" :label="$t('contract.completionStatus')" width="100">
              <template #default="{ row }">
                <el-tag :type="row.completion_status === 'completed' ? 'success' : row.completion_status === 'delayed' ? 'danger' : 'warning'" size="small">
                  {{ row.completion_status === 'completed' ? $t('contract.completed') : row.completion_status === 'delayed' ? $t('contract.delayed') : $t('contract.inProgress') }}
                </el-tag>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>

      <el-col :span="8">
        <el-card>
          <template #header>
            <span>{{ $t('contract.approvalFlow') }}</span>
          </template>
          <el-timeline>
            <el-timeline-item v-for="item in approvalFlows" :key="item.id" :timestamp="item.approval_time" placement="top">
              <el-card>
                <p><strong>{{ $t('contract.approver') }}:</strong> {{ item.approver_name }}</p>
                <p><strong>{{ $t('contract.result') }}:</strong> 
                  <el-tag :type="item.approval_status === 'approved' ? 'success' : 'danger'" size="small">
                    {{ item.approval_status === 'approved' ? $t('contract.approved') : $t('contract.rejected') }}
                  </el-tag>
                </p>
                <p v-if="item.comment"><strong>{{ $t('contract.comment') }}:</strong> {{ item.comment }}</p>
              </el-card>
            </el-timeline-item>
          </el-timeline>
        </el-card>
      </el-col>
    </el-row>

    <el-dialog v-model="showPerformanceDialog" :title="$t('contract.addRecord')" width="500px">
      <el-form :model="performanceForm" label-width="80px">
        <el-form-item :label="$t('contract.milestone')">
          <el-input v-model="performanceForm.milestone_name" />
        </el-form-item>
        <el-form-item :label="$t('contract.description')">
          <el-input v-model="performanceForm.milestone_desc" type="textarea" :rows="2" />
        </el-form-item>
        <el-form-item :label="$t('contract.plannedDate')">
          <el-date-picker v-model="performanceForm.planned_date" type="date" style="width: 100%" />
        </el-form-item>
        <el-form-item :label="$t('contract.completionStatus')">
          <el-select v-model="performanceForm.completion_status" style="width: 100%">
            <el-option :label="$t('contract.inProgress')" value="in_progress" />
            <el-option :label="$t('contract.completed')" value="completed" />
            <el-option :label="$t('contract.delayed')" value="delayed" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('contract.notes')">
          <el-input v-model="performanceForm.notes" type="textarea" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showPerformanceDialog = false">{{ $t('common.cancel') }}</el-button>
        <el-button type="primary" @click="addPerformance" :loading="saving">{{ $t('common.save') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { contractApi } from '../../api'
import { mockContracts, mockPerformanceRecords } from '../../mock/data'

const { t } = useI18n()
const route = useRoute()
const contractId = Number(route.params.id)

const statusColors: any = { draft: 'info', pending: 'warning', active: 'success', expired: 'danger' }
const statusLabels = computed(() => ({
  draft: t('contract.draft'),
  pending: t('contract.pending'),
  active: t('contract.active'),
  expired: t('contract.expired')
}))

const contract = ref<any>({})
const performances = ref<any[]>([])
const approvalFlows = ref<any[]>([])
const showPerformanceDialog = ref(false)
const saving = ref(false)

const performanceForm = reactive({
  milestone_name: '',
  milestone_desc: '',
  planned_date: '',
  completion_status: 'in_progress',
  notes: ''
})

const loadContract = async () => {
  try {
    const res = await contractApi.getDetail(contractId)
    contract.value = res.data.data
  } catch (error) {
    contract.value = mockContracts.find(c => c.id === contractId) || {}
  }
}

const loadPerformances = async () => {
  try {
    const res = await contractApi.getPerformance(contractId)
    performances.value = res.data.data || []
  } catch (error) {
    performances.value = mockPerformanceRecords.filter(p => p.contract_id === contractId)
  }
}

const addPerformance = async () => {
  saving.value = true
  try {
    await contractApi.createPerformance(contractId, performanceForm)
    ElMessage.success(t('common.success'))
    showPerformanceDialog.value = false
    loadPerformances()
  } catch (error: any) {
    performances.value.push({
      id: Date.now(),
      contract_id: contractId,
      ...performanceForm,
      planned_date: performanceForm.planned_date || new Date().toISOString().split('T')[0]
    })
    ElMessage.success(t('common.success'))
    showPerformanceDialog.value = false
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  loadContract()
  loadPerformances()
})
</script>
