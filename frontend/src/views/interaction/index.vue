<template>
  <div class="interaction-list">
    <div class="page-header">
      <div class="page-title">
        <h1>{{ $t('interaction.title') }}</h1>
      </div>
      <el-button type="primary" @click="openCreateDialog">
        <el-icon><Plus /></el-icon>
        {{ $t('interaction.record') }}
      </el-button>
    </div>

    <div class="search-card glass">
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item>
          <el-select v-model="searchForm.customer_id" :placeholder="$t('interaction.customer')" clearable style="width: 200px">
            <el-option v-for="c in customers" :key="c.id" :label="c.customer_name" :value="c.id" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-select v-model="searchForm.type" :placeholder="$t('interaction.type')" clearable>
            <el-option :label="$t('interaction.meeting')" value="meeting" />
            <el-option :label="$t('interaction.phone')" value="phone" />
            <el-option :label="$t('interaction.email')" value="email" />
            <el-option :label="$t('interaction.visit')" value="visit" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadInteractions">{{ $t('common.search') }}</el-button>
          <el-button @click="resetSearch">{{ $t('common.reset') }}</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="table-card glass">
      <el-table :data="interactions" v-loading="loading" :element-loading-text="$t('common.loading')" element-loading-background="rgba(10, 15, 28, 0.7)" style="width: 100%">
        <el-table-column prop="interaction_date" :label="$t('interaction.date')" width="120" />
        <el-table-column prop="customer_name" :label="$t('interaction.customer')" width="150" />
        <el-table-column prop="interaction_type" :label="$t('interaction.type')" width="100">
          <template #default="{ row }">
            <el-tag size="small">{{ typeLabels[row.interaction_type] || row.interaction_type }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="topic" :label="$t('interaction.topic')" />
        <el-table-column prop="outcome" :label="$t('interaction.outcome')" width="100">
          <template #default="{ row }">
            <el-tag v-if="row.outcome" :type="row.outcome === 'positive' ? 'success' : row.outcome === 'negative' ? 'danger' : 'info'" size="small">
              {{ outcomeLabels[row.outcome] || row.outcome }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="importance_level" :label="$t('interaction.importance')" width="140">
          <template #default="{ row }">
            <el-rate v-model="row.importance_level" disabled />
          </template>
        </el-table-column>
        <el-table-column :label="$t('common.action')" width="150" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="editInteraction(row)">{{ $t('common.edit') }}</el-button>
            <el-button type="danger" link @click="deleteInteraction(row.id)">{{ $t('common.delete') }}</el-button>
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
    </div>

    <el-dialog v-model="showDialog" :title="editingInteraction ? $t('common.edit') : $t('interaction.record')" width="650px">
      <el-form :model="interactionForm" :rules="rules" ref="formRef" label-width="120px">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="$t('interaction.customer')" prop="customer_id">
              <el-select v-model="interactionForm.customer_id" style="width: 100%">
                <el-option v-for="c in customers" :key="c.id" :label="c.customer_name" :value="c.id" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('interaction.type')" prop="interaction_type">
              <el-select v-model="interactionForm.interaction_type" style="width: 100%">
                <el-option :label="$t('interaction.meeting')" value="meeting" />
                <el-option :label="$t('interaction.phone')" value="phone" />
                <el-option :label="$t('interaction.email')" value="email" />
                <el-option :label="$t('interaction.visit')" value="visit" />
                <el-option :label="$t('interaction.other')" value="other" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="$t('interaction.date')">
              <el-date-picker v-model="interactionForm.interaction_date" type="datetime" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('interaction.importance')">
              <el-rate v-model="interactionForm.importance_level" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item :label="$t('interaction.topic')" prop="topic">
          <el-input v-model="interactionForm.topic" />
        </el-form-item>
        <el-form-item :label="$t('interaction.summary')">
          <el-input v-model="interactionForm.content_summary" type="textarea" :rows="3" />
        </el-form-item>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="$t('interaction.outcome')">
              <el-select v-model="interactionForm.outcome" style="width: 100%">
                <el-option :label="$t('interaction.positive')" value="positive" />
                <el-option :label="$t('interaction.neutral')" value="neutral" />
                <el-option :label="$t('interaction.negative')" value="negative" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('interaction.nextAction')">
              <el-input v-model="interactionForm.next_action" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <el-button @click="showDialog = false">{{ $t('common.cancel') }}</el-button>
        <el-button type="primary" @click="saveInteraction" :loading="saving">{{ $t('common.save') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { interactionApi, customerApi } from '../../api'
import { mockInteractions, mockCustomers } from '../../mock/data'

const { t } = useI18n()
const loading = ref(false)
const saving = ref(false)
const showDialog = ref(false)
const interactions = ref<any[]>([])
const customers = ref<any[]>([])
const editingInteraction = ref<any>(null)
const formRef = ref()

const typeLabels = computed(() => ({
  meeting: t('interaction.meeting'),
  phone: t('interaction.phone'),
  email: t('interaction.email'),
  visit: t('interaction.visit'),
  other: t('interaction.other')
}))

const outcomeLabels = computed(() => ({
  positive: t('interaction.positive'),
  neutral: t('interaction.neutral'),
  negative: t('interaction.negative')
}))

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

const rules = computed(() => ({
  customer_id: [{ required: true, message: t('common.required'), trigger: 'change' }],
  interaction_type: [{ required: true, message: t('common.required'), trigger: 'change' }],
  topic: [{ required: true, message: t('common.required'), trigger: 'blur' }]
}))

const loadInteractions = async () => {
  loading.value = true
  try {
    const res = await interactionApi.getList({ page: pagination.page, page_size: pagination.pageSize, ...searchForm })
    interactions.value = res.data.data || []
    pagination.total = res.data.pagination?.total || 0
  } catch (error) {
    interactions.value = mockInteractions
    pagination.total = mockInteractions.length
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

const resetSearch = () => {
  searchForm.customer_id = ''
  searchForm.type = ''
  loadInteractions()
}

const openCreateDialog = () => {
  editingInteraction.value = null
  Object.assign(interactionForm, { customer_id: null, interaction_type: 'meeting', interaction_date: '', topic: '', content_summary: '', outcome: '', next_action: '', importance_level: 3 })
  showDialog.value = true
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
      ElMessage.success(t('common.success'))
    } else {
      await interactionApi.create(interactionForm)
      ElMessage.success(t('common.success'))
    }
    showDialog.value = false
    loadInteractions()
  } catch (error: any) {
    ElMessage.error(error.response?.data?.error || t('common.error'))
  } finally {
    saving.value = false
  }
}

const deleteInteraction = async (id: number) => {
  await ElMessageBox.confirm(t('common.confirm') + '?', t('common.delete'), { type: 'warning' })
  try {
    await interactionApi.delete(id)
    ElMessage.success(t('common.success'))
    loadInteractions()
  } catch (error: any) {
    ElMessage.error(error.response?.data?.error || t('common.error'))
  }
}

onMounted(() => {
  loadInteractions()
  loadCustomers()
})
</script>

<style scoped>
.interaction-list { max-width: 1400px; }
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 24px; }
.page-title h1 { font-family: var(--font-display); font-size: 24px; font-weight: 700; background: linear-gradient(135deg, var(--color-accent) 0%, var(--color-accent-secondary) 100%); -webkit-background-clip: text; -webkit-text-fill-color: transparent; background-clip: text; margin: 0; }
.search-card { padding: 20px; margin-bottom: 20px; border-radius: var(--radius-md); position: relative; overflow: hidden; }
.search-card::before { content: ''; position: absolute; top: 0; left: 0; right: 0; height: 2px; background: linear-gradient(90deg, var(--color-accent), var(--color-accent-secondary)); }
.table-card { padding: 20px; border-radius: var(--radius-md); position: relative; overflow: hidden; }
.table-card::before { content: ''; position: absolute; top: 0; left: 0; right: 0; height: 2px; background: linear-gradient(90deg, var(--color-accent), var(--color-accent-secondary)); }
.search-form { display: flex; flex-wrap: wrap; gap: 12px; }
</style>
