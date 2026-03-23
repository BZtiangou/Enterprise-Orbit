<template>
  <div class="customer-detail">
    <div class="page-header">
      <el-page-header @back="$router.back()" :title="$t('common.view')">
        <template #content>
          <span class="customer-name">{{ customer.customer_name || $t('common.loading') }}</span>
        </template>
      </el-page-header>
    </div>

    <el-tabs v-model="activeTab" class="detail-tabs">
      <el-tab-pane :label="$t('customer.basicInfo')" name="basic">
        <el-row :gutter="24">
          <el-col :span="16">
            <div class="info-card">
              <h3>{{ $t('customer.basicInfo') }}</h3>
              <el-descriptions :column="2" border>
                <el-descriptions-item :label="$t('customer.customerCode')">{{ customer.customer_code }}</el-descriptions-item>
                <el-descriptions-item :label="$t('customer.customerName')">{{ customer.customer_name }}</el-descriptions-item>
                <el-descriptions-item :label="$t('customer.industry')">{{ customer.industry }}</el-descriptions-item>
                <el-descriptions-item :label="$t('customer.scale')">{{ scaleLabel(customer.scale) }}</el-descriptions-item>
                <el-descriptions-item :label="$t('customer.region')">{{ customer.region }}</el-descriptions-item>
                <el-descriptions-item :label="$t('common.status')">
                  <el-tag :type="customer.status === 'active' ? 'success' : 'info'">
                    {{ customer.status === 'active' ? $t('common.active') : $t('common.inactive') }}
                  </el-tag>
                </el-descriptions-item>
                <el-descriptions-item :label="$t('customer.address')" :span="2">{{ customer.address }}</el-descriptions-item>
                <el-descriptions-item :label="$t('customer.website')" :span="2">
                  <a :href="customer.website" target="_blank" class="website-link">{{ customer.website }}</a>
                </el-descriptions-item>
              </el-descriptions>
            </div>
          </el-col>

          <el-col :span="8">
            <div class="info-card health-card">
              <div class="card-header">
                <h3>{{ $t('customer.relationshipHealth') }}</h3>
                <el-button type="primary" size="small" @click="calculateHealth">{{ $t('customer.recalculate') }}</el-button>
              </div>
              <div class="health-score">
                <el-progress type="dashboard" :percentage="Math.round(health.overall_score || 0)" :color="healthColor" :width="120" />
                <div class="score-label">{{ $t('customer.overallScore') }}</div>
              </div>
              <div class="health-breakdown">
                <div class="health-item">
                  <span class="label">{{ $t('customer.trustScore') }}</span>
                  <el-progress :percentage="Math.round(health.trust_score || 0)" :stroke-width="8" />
                </div>
                <div class="health-item">
                  <span class="label">{{ $t('customer.commitmentScore') }}</span>
                  <el-progress :percentage="Math.round(health.commitment_score || 0)" :stroke-width="8" />
                </div>
                <div class="health-item">
                  <span class="label">{{ $t('customer.reciprocityScore') }}</span>
                  <el-progress :percentage="Math.round(health.reciprocity_score || 0)" :stroke-width="8" />
                </div>
              </div>
              <div class="risk-level">
                <span class="label">{{ $t('customer.riskLevel') }}</span>
                <el-tag :type="riskTagType" size="large">{{ riskLabel }}</el-tag>
              </div>
            </div>
          </el-col>
        </el-row>
      </el-tab-pane>

      <el-tab-pane :label="$t('customer.organizationStructure', '组织架构')" name="org">
        <div class="info-card">
          <div class="card-header">
            <h3>{{ $t('customer.organizationStructure', '组织架构') }}</h3>
            <el-button type="primary" size="small" @click="showOrgDialog = true">{{ $t('common.create') }}</el-button>
          </div>
          <div class="org-tree-container">
            <el-tree
              v-if="orgTree.length > 0"
              :data="orgTree"
              :props="{ label: 'dept_name', children: 'children' }"
              default-expand-all
              node-key="id"
              :expand-on-click-node="false"
            >
              <template #default="{ node, data }">
                <div class="org-node">
                  <span class="dept-name">{{ node.label }}</span>
                  <span class="dept-function">{{ data.dept_function }}</span>
                </div>
              </template>
            </el-tree>
            <el-empty v-else :description="$t('common.noData', '暂无数据')" />
          </div>
        </div>
      </el-tab-pane>

      <el-tab-pane :label="$t('customer.keyDecisionMakers', '关键决策人')" name="contacts">
        <div class="info-card">
          <div class="card-header">
            <h3>{{ $t('customer.keyDecisionMakers', '关键决策人') }}</h3>
            <el-button type="primary" size="small" @click="showContactDialog = true">{{ $t('common.create') }}</el-button>
          </div>
          <el-table :data="contacts" v-loading="contactsLoading" :element-loading-text="$t('common.loading')" element-loading-background="rgba(10, 15, 28, 0.7)">
            <el-table-column prop="name" :label="$t('customer.contactName', '姓名')" width="120" />
            <el-table-column prop="position" :label="$t('profile.position')" width="140" />
            <el-table-column prop="department" :label="$t('profile.department')" width="120" />
            <el-table-column prop="phone" :label="$t('profile.phone')" width="130" />
            <el-table-column prop="email" :label="$t('auth.email')" />
            <el-table-column :label="$t('customer.isKeyDecisionMaker', '关键决策人')" width="100">
              <template #default="{ row }">
                <el-tag :type="row.is_key_decision_maker ? 'danger' : 'info'" size="small">
                  {{ row.is_key_decision_maker ? $t('common.yes', '是') : $t('common.no', '否') }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column :label="$t('customer.influenceLevel', '影响力')" width="150">
              <template #default="{ row }">
                <el-rate v-model="row.influence_level" disabled :max="5" />
              </template>
            </el-table-column>
            <el-table-column :label="$t('common.action')" width="120">
              <template #default="{ row }">
                <el-button type="primary" link @click="editContact(row)">{{ $t('common.edit') }}</el-button>
                <el-button type="danger" link @click="deleteContact(row.id)">{{ $t('common.delete') }}</el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>
      </el-tab-pane>

      <el-tab-pane :label="$t('customer.cooperationHistory', '合作历史')" name="history">
        <div class="info-card">
          <div class="card-header">
            <h3>{{ $t('customer.cooperationHistory', '合作历史') }}</h3>
            <el-button type="primary" size="small" @click="showHistoryDialog = true">{{ $t('common.create') }}</el-button>
          </div>
          <el-timeline v-if="cooperationHistory.length > 0">
            <el-timeline-item
              v-for="item in cooperationHistory"
              :key="item.id"
              :timestamp="item.start_date"
              placement="top"
              :type="item.status === 'completed' ? 'success' : 'primary'"
            >
              <div class="history-card">
                <div class="history-header">
                  <h4>{{ item.project_name }}</h4>
                  <el-tag :type="historyStatusType(item.status)" size="small">
                    {{ historyStatusLabel(item.status) }}
                  </el-tag>
                </div>
                <div class="history-content">
                  <div class="history-info">
                    <span class="info-item">
                      <span class="label">{{ $t('contract.amount') }}:</span>
                      <span class="value">¥{{ item.contract_amount?.toLocaleString() }}</span>
                    </span>
                    <span class="info-item">
                      <span class="label">{{ $t('contract.startDate') }}:</span>
                      <span class="value">{{ item.start_date }}</span>
                    </span>
                    <span class="info-item" v-if="item.end_date">
                      <span class="label">{{ $t('contract.endDate') }}:</span>
                      <span class="value">{{ item.end_date }}</span>
                    </span>
                  </div>
                  <div class="history-desc">{{ item.description }}</div>
                  <div class="history-satisfaction" v-if="item.satisfaction_score">
                    <span class="label">{{ $t('customer.satisfactionScore', '满意度') }}:</span>
                    <el-rate v-model="item.satisfaction_score" disabled :max="5" />
                  </div>
                </div>
              </div>
            </el-timeline-item>
          </el-timeline>
          <el-empty v-else :description="$t('common.noData', '暂无数据')" />
        </div>
      </el-tab-pane>

      <el-tab-pane :label="$t('customer.interactionRecords')" name="interactions">
        <div class="info-card">
          <div class="card-header">
            <h3>{{ $t('customer.interactionRecords') }}</h3>
            <el-button type="primary" size="small" @click="showInteractionDialog = true">{{ $t('customer.addInteraction') }}</el-button>
          </div>
          <el-timeline v-if="interactions.length > 0">
            <el-timeline-item
              v-for="item in interactions"
              :key="item.id"
              :timestamp="item.interaction_date"
              placement="top"
              :type="outcomeType(item.outcome)"
            >
              <div class="interaction-card">
                <div class="interaction-header">
                  <h4>{{ item.topic }}</h4>
                  <div class="interaction-tags">
                    <el-tag size="small">{{ interactionTypeLabel(item.interaction_type) }}</el-tag>
                    <el-tag v-if="item.outcome" size="small" :type="outcomeTagType(item.outcome)" style="margin-left: 8px">
                      {{ outcomeLabel(item.outcome) }}
                    </el-tag>
                  </div>
                </div>
                <p class="interaction-summary">{{ item.summary }}</p>
                <div class="interaction-footer" v-if="item.next_action">
                  <span class="next-action">{{ $t('interaction.nextAction') }}: {{ item.next_action }}</span>
                </div>
              </div>
            </el-timeline-item>
          </el-timeline>
          <el-empty v-else :description="$t('common.noData', '暂无数据')" />
        </div>
      </el-tab-pane>

      <el-tab-pane :label="$t('customer.relatedContracts')" name="contracts">
        <div class="info-card">
          <div class="card-header">
            <h3>{{ $t('customer.relatedContracts') }}</h3>
          </div>
          <el-table :data="contracts" v-loading="contractsLoading" :element-loading-text="$t('common.loading')" element-loading-background="rgba(10, 15, 28, 0.7)">
            <el-table-column prop="contract_no" :label="$t('contract.contractNo')" width="140" />
            <el-table-column prop="contract_name" :label="$t('contract.contractName')" />
            <el-table-column prop="amount" :label="$t('contract.amount')" width="120">
              <template #default="{ row }">¥{{ row.amount?.toLocaleString() }}</template>
            </el-table-column>
            <el-table-column prop="start_date" :label="$t('contract.startDate')" width="110" />
            <el-table-column prop="end_date" :label="$t('contract.endDate')" width="110" />
            <el-table-column prop="status" :label="$t('common.status')" width="100">
              <template #default="{ row }">
                <el-tag :type="contractStatusType(row.status)">{{ contractStatusLabel(row.status) }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column :label="$t('common.action')" width="100">
              <template #default="{ row }">
                <el-button type="primary" link @click="$router.push(`/contracts/${row.id}`)">{{ $t('common.detail') }}</el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>
      </el-tab-pane>
    </el-tabs>

    <el-dialog v-model="showContactDialog" :title="editingContact ? $t('common.edit') : $t('common.create')" width="550px">
      <el-form :model="contactForm" :rules="contactRules" ref="contactFormRef" label-width="120px">
        <el-form-item :label="$t('customer.contactName', '姓名')" prop="name">
          <el-input v-model="contactForm.name" />
        </el-form-item>
        <el-form-item :label="$t('profile.position')" prop="position">
          <el-input v-model="contactForm.position" />
        </el-form-item>
        <el-form-item :label="$t('profile.department')">
          <el-input v-model="contactForm.department" />
        </el-form-item>
        <el-form-item :label="$t('profile.phone')">
          <el-input v-model="contactForm.phone" />
        </el-form-item>
        <el-form-item :label="$t('auth.email')">
          <el-input v-model="contactForm.email" />
        </el-form-item>
        <el-form-item :label="$t('customer.isKeyDecisionMaker', '关键决策人')">
          <el-switch v-model="contactForm.is_key_decision_maker" />
        </el-form-item>
        <el-form-item :label="$t('customer.influenceLevel', '影响力')">
          <el-rate v-model="contactForm.influence_level" :max="5" />
        </el-form-item>
        <el-form-item :label="$t('customer.notes', '备注')">
          <el-input v-model="contactForm.notes" type="textarea" :rows="3" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showContactDialog = false">{{ $t('common.cancel') }}</el-button>
        <el-button type="primary" @click="saveContact" :loading="saving">{{ $t('common.save') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { customerApi, interactionApi, contractApi } from '../../api'
import { mockCustomers, mockContacts, mockOrgStructures, mockCooperationHistory, mockInteractions, mockContracts } from '../../mock/data'

const { t } = useI18n()
const route = useRoute()
const customerId = Number(route.params.id)

const activeTab = ref('basic')
const customer = ref<any>({})
const health = ref<any>({})
const contacts = ref<any[]>([])
const orgTree = ref<any[]>([])
const cooperationHistory = ref<any[]>([])
const interactions = ref<any[]>([])
const contracts = ref<any[]>([])
const contactsLoading = ref(false)
const contractsLoading = ref(false)
const saving = ref(false)
const showContactDialog = ref(false)
const showOrgDialog = ref(false)
const showHistoryDialog = ref(false)
const showInteractionDialog = ref(false)
const editingContact = ref<any>(null)
const contactFormRef = ref()

const contactForm = reactive({
  name: '',
  position: '',
  department: '',
  phone: '',
  email: '',
  is_key_decision_maker: false,
  influence_level: 3,
  notes: ''
})

const contactRules = computed(() => ({
  name: [{ required: true, message: t('common.required'), trigger: 'blur' }],
  position: [{ required: true, message: t('common.required'), trigger: 'blur' }]
}))

const healthColor = computed(() => {
  const score = health.value.overall_score || 0
  return score > 70 ? '#10b981' : score > 40 ? '#f59e0b' : '#ef4444'
})

const riskTagType = computed(() => {
  const level = health.value.risk_level
  return level === 'low' ? 'success' : level === 'medium' ? 'warning' : 'danger'
})

const riskLabel = computed(() => {
  const level = health.value.risk_level
  return level === 'low' ? t('customer.lowRisk') : level === 'medium' ? t('customer.mediumRisk') : t('customer.highRisk')
})

const scaleLabel = (scale: string) => {
  const labels: any = { large: t('customer.large'), medium: t('customer.medium'), small: t('customer.small') }
  return labels[scale] || scale
}

const historyStatusType = (status: string) => {
  const types: any = { completed: 'success', in_progress: 'primary', pending: 'info' }
  return types[status] || 'info'
}

const historyStatusLabel = (status: string) => {
  const labels: any = { completed: t('contract.completed'), in_progress: t('contract.inProgress'), pending: t('contract.draft') }
  return labels[status] || status
}

const outcomeType = (outcome: string) => {
  const types: any = { positive: 'success', neutral: 'primary', negative: 'danger' }
  return types[outcome] || 'primary'
}

const outcomeTagType = (outcome: string) => {
  const types: any = { positive: 'success', neutral: 'info', negative: 'danger' }
  return types[outcome] || 'info'
}

const outcomeLabel = (outcome: string) => {
  const labels: any = { positive: t('interaction.positive'), neutral: t('interaction.neutral'), negative: t('interaction.negative') }
  return labels[outcome] || outcome
}

const interactionTypeLabel = (type: string) => {
  const labels: any = { meeting: t('interaction.meeting'), phone: t('interaction.phone'), email: t('interaction.email'), visit: t('interaction.visit'), other: t('interaction.other') }
  return labels[type] || type
}

const contractStatusType = (status: string) => {
  const types: any = { active: 'success', pending: 'warning', draft: 'info', expired: 'danger' }
  return types[status] || 'info'
}

const contractStatusLabel = (status: string) => {
  const labels: any = { active: t('contract.active'), pending: t('contract.pending'), draft: t('contract.draft'), expired: t('contract.expired') }
  return labels[status] || status
}

const loadCustomer = async () => {
  try {
    const res = await customerApi.getDetail(customerId)
    customer.value = res.data.data || mockCustomers.find(c => c.id === customerId) || {}
  } catch (error) {
    customer.value = mockCustomers.find(c => c.id === customerId) || {}
  }
}

const loadHealth = async () => {
  try {
    const res = await customerApi.getHealth(customerId)
    health.value = res.data.data || {}
  } catch (error) {
    const mockCustomer = mockCustomers.find(c => c.id === customerId)
    if (mockCustomer) {
      health.value = {
        overall_score: mockCustomer.overall_score || 75,
        trust_score: mockCustomer.trust_score || 75,
        commitment_score: mockCustomer.commitment_score || 75,
        reciprocity_score: mockCustomer.reciprocity_score || 75,
        risk_level: mockCustomer.overall_score > 70 ? 'low' : mockCustomer.overall_score > 50 ? 'medium' : 'high'
      }
    } else {
      health.value = { overall_score: 75, trust_score: 80, commitment_score: 70, reciprocity_score: 75, risk_level: 'low' }
    }
  }
}

const calculateHealth = async () => {
  try {
    const res = await customerApi.calculateHealth(customerId)
    health.value = res.data.data
    ElMessage.success(t('common.success'))
  } catch (error) {
    health.value = { overall_score: Math.floor(Math.random() * 30 + 60), trust_score: Math.floor(Math.random() * 30 + 60), commitment_score: Math.floor(Math.random() * 30 + 60), reciprocity_score: Math.floor(Math.random() * 30 + 60), risk_level: ['low', 'medium', 'high'][Math.floor(Math.random() * 3)] }
    ElMessage.success(t('common.success'))
  }
}

const loadContacts = async () => {
  contactsLoading.value = true
  try {
    const res = await customerApi.getContacts(customerId)
    contacts.value = res.data.data || []
  } catch (error) {
    contacts.value = mockContacts.filter(c => c.customer_id === customerId)
  } finally {
    contactsLoading.value = false
  }
}

const loadOrgStructure = async () => {
  try {
    const res = await customerApi.getOrgStructure(customerId)
    orgTree.value = res.data.data || []
  } catch (error) {
    orgTree.value = mockOrgStructures.filter(o => o.customer_id === customerId)
  }
}

const loadCooperationHistory = async () => {
  try {
    const res = await customerApi.getCooperationHistory(customerId)
    cooperationHistory.value = res.data.data || []
  } catch (error) {
    cooperationHistory.value = mockCooperationHistory.filter(h => h.customer_id === customerId)
  }
}

const loadInteractions = async () => {
  try {
    const res = await interactionApi.getList({ customer_id: customerId })
    interactions.value = res.data.data || []
  } catch (error) {
    interactions.value = mockInteractions.filter(i => i.customer_id === customerId)
  }
}

const loadContracts = async () => {
  contractsLoading.value = true
  try {
    const res = await contractApi.getList({ customer_id: customerId })
    contracts.value = res.data.data || []
  } catch (error) {
    contracts.value = mockContracts.filter(c => c.customer_id === customerId)
  } finally {
    contractsLoading.value = false
  }
}

const editContact = (contact: any) => {
  editingContact.value = contact
  Object.assign(contactForm, contact)
  showContactDialog.value = true
}

const saveContact = async () => {
  await contactFormRef.value.validate()
  saving.value = true
  try {
    if (editingContact.value) {
      const index = contacts.value.findIndex(c => c.id === editingContact.value.id)
      if (index > -1) contacts.value[index] = { ...editingContact.value, ...contactForm }
    } else {
      contacts.value.push({ id: Date.now(), customer_id: customerId, ...contactForm, created_at: new Date().toISOString() })
    }
    ElMessage.success(t('common.success'))
    showContactDialog.value = false
    resetContactForm()
  } catch (error) {
    ElMessage.error(t('common.error'))
  } finally {
    saving.value = false
  }
}

const deleteContact = async (id: number) => {
  await ElMessageBox.confirm(t('common.confirm') + '?', t('common.delete'), { type: 'warning' })
  contacts.value = contacts.value.filter(c => c.id !== id)
  ElMessage.success(t('common.success'))
}

const resetContactForm = () => {
  editingContact.value = null
  Object.assign(contactForm, { name: '', position: '', department: '', phone: '', email: '', is_key_decision_maker: false, influence_level: 3, notes: '' })
}

onMounted(() => {
  loadCustomer()
  loadHealth()
  loadContacts()
  loadOrgStructure()
  loadCooperationHistory()
  loadInteractions()
  loadContracts()
})
</script>

<style scoped>
.customer-detail { max-width: 1400px; }
.page-header { margin-bottom: 24px; }
.customer-name { font-family: var(--font-display); font-size: 20px; font-weight: 600; background: linear-gradient(135deg, var(--color-accent) 0%, var(--color-accent-secondary) 100%); -webkit-background-clip: text; -webkit-text-fill-color: transparent; background-clip: text; }
.detail-tabs { margin-top: 20px; }
.info-card { padding: 24px; border-radius: var(--radius-md); margin-bottom: 20px; position: relative; overflow: hidden; background: transparent; border: 1px solid var(--color-border); }
.info-card::before { content: ''; position: absolute; top: 0; left: 0; right: 0; height: 3px; background: linear-gradient(90deg, var(--color-accent), var(--color-accent-secondary)); }
.info-card h3 { font-family: var(--font-display); font-size: 18px; font-weight: 600; color: var(--color-text-primary); margin: 0 0 20px 0; }
.card-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; }
.card-header h3 { margin: 0; }
.website-link { color: var(--color-accent); text-decoration: none; transition: all var(--transition-fast); }
.website-link:hover { text-decoration: underline; text-shadow: 0 0 10px var(--color-accent-glow); }
.health-card { text-align: center; }
.health-card::before { display: none; }
.health-score { padding: 20px 0; }
.score-label { margin-top: 10px; font-size: 14px; color: var(--color-text-secondary); }
.health-breakdown { margin-top: 20px; text-align: left; }
.health-item { margin-bottom: 16px; padding: 12px; background: transparent; border-radius: var(--radius-sm); transition: all var(--transition-fast); border: 1px solid var(--color-border); }
.health-item:hover { background: rgba(255, 255, 255, 0.02); border-color: var(--color-accent); }
.health-item .label { display: block; margin-bottom: 8px; font-size: 13px; color: var(--color-text-secondary); }
.risk-level { margin-top: 20px; padding-top: 20px; border-top: 1px solid var(--color-border); display: flex; justify-content: space-between; align-items: center; }
.risk-level .label { font-size: 14px; color: var(--color-text-secondary); }
.org-tree-container { min-height: 200px; padding: 16px; background: transparent; border-radius: var(--radius-sm); }
.org-node { display: flex; align-items: center; gap: 12px; padding: 8px 12px; border-radius: var(--radius-sm); transition: all var(--transition-fast); }
.org-node:hover { background: rgba(6, 182, 212, 0.1); }
.dept-name { font-weight: 500; color: var(--color-text-primary); }
.dept-function { font-size: 12px; color: var(--color-accent); background: rgba(6, 182, 212, 0.15); padding: 2px 8px; border-radius: 4px; }
.history-card { background: transparent; padding: 20px; border-radius: var(--radius-md); border: 1px solid var(--color-border); transition: all var(--transition-normal); }
.history-card:hover { border-color: var(--color-accent); box-shadow: 0 4px 20px rgba(6, 182, 212, 0.1); }
.history-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 12px; }
.history-header h4 { margin: 0; font-size: 16px; color: var(--color-text-primary); font-family: var(--font-display); }
.history-info { display: flex; flex-wrap: wrap; gap: 24px; margin-bottom: 12px; }
.info-item { font-size: 13px; }
.info-item .label { color: var(--color-text-secondary); margin-right: 8px; }
.info-item .value { color: var(--color-accent); font-weight: 500; }
.history-desc { color: var(--color-text-secondary); font-size: 14px; line-height: 1.6; }
.history-satisfaction { margin-top: 12px; display: flex; align-items: center; gap: 8px; }
.history-satisfaction .label { font-size: 13px; color: var(--color-text-secondary); }
.interaction-card { background: transparent; padding: 20px; border-radius: var(--radius-md); border: 1px solid var(--color-border); transition: all var(--transition-normal); }
.interaction-card:hover { border-color: var(--color-accent); box-shadow: 0 4px 20px rgba(6, 182, 212, 0.1); }
.interaction-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 12px; }
.interaction-header h4 { margin: 0; font-size: 16px; color: var(--color-text-primary); font-family: var(--font-display); }
.interaction-tags { display: flex; }
.interaction-summary { color: var(--color-text-secondary); font-size: 14px; line-height: 1.6; margin: 0; }
.interaction-footer { margin-top: 12px; padding-top: 12px; border-top: 1px solid var(--color-border); }
.next-action { font-size: 13px; color: var(--color-accent); }
</style>
