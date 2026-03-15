<template>
  <div class="customer-detail">
    <el-page-header @back="$router.back()" title="返回">
      <template #content>
        <span>{{ customer.customer_name }}</span>
      </template>
    </el-page-header>

    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="16">
        <el-card>
          <template #header>
            <span>基本信息</span>
          </template>
          <el-descriptions :column="2" border>
            <el-descriptions-item label="客户编码">{{ customer.customer_code }}</el-descriptions-item>
            <el-descriptions-item label="客户名称">{{ customer.customer_name }}</el-descriptions-item>
            <el-descriptions-item label="行业">{{ customer.industry }}</el-descriptions-item>
            <el-descriptions-item label="规模">{{ customer.scale }}</el-descriptions-item>
            <el-descriptions-item label="区域">{{ customer.region }}</el-descriptions-item>
            <el-descriptions-item label="状态">
              <el-tag :type="customer.status === 'active' ? 'success' : 'info'">
                {{ customer.status === 'active' ? '活跃' : '暂停' }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="地址" :span="2">{{ customer.address }}</el-descriptions-item>
            <el-descriptions-item label="网站" :span="2">
              <a :href="customer.website" target="_blank">{{ customer.website }}</a>
            </el-descriptions-item>
          </el-descriptions>
        </el-card>

        <el-card style="margin-top: 20px">
          <template #header>
            <div style="display: flex; justify-content: space-between; align-items: center">
              <span>互动记录</span>
              <el-button type="primary" size="small" @click="showInteractionDialog = true">添加互动</el-button>
            </div>
          </template>
          <el-timeline>
            <el-timeline-item v-for="item in interactions" :key="item.id" :timestamp="item.interaction_date" placement="top">
              <el-card>
                <h4>{{ item.topic }}</h4>
                <p>{{ item.content_summary }}</p>
                <div style="margin-top: 10px">
                  <el-tag size="small">{{ item.interaction_type }}</el-tag>
                  <el-tag v-if="item.outcome" size="small" type="success" style="margin-left: 8px">{{ item.outcome }}</el-tag>
                </div>
              </el-card>
            </el-timeline-item>
          </el-timeline>
        </el-card>
      </el-col>

      <el-col :span="8">
        <el-card>
          <template #header>
            <div style="display: flex; justify-content: space-between; align-items: center">
              <span>关系健康度</span>
              <el-button type="primary" size="small" @click="calculateHealth">重新计算</el-button>
            </div>
          </template>
          <div class="health-score">
            <el-progress type="dashboard" :percentage="health.overall_score || 0" :color="healthColor" />
            <div class="score-label">综合评分</div>
          </div>
          <el-descriptions :column="1" border style="margin-top: 20px">
            <el-descriptions-item label="信任度">{{ health.trust_score?.toFixed(1) }}</el-descriptions-item>
            <el-descriptions-item label="承诺度">{{ health.commitment_score?.toFixed(1) }}</el-descriptions-item>
            <el-descriptions-item label="互惠度">{{ health.reciprocity_score?.toFixed(1) }}</el-descriptions-item>
            <el-descriptions-item label="风险等级">
              <el-tag :type="health.risk_level === 'low' ? 'success' : health.risk_level === 'medium' ? 'warning' : 'danger'">
                {{ health.risk_level === 'low' ? '低风险' : health.risk_level === 'medium' ? '中风险' : '高风险' }}
              </el-tag>
            </el-descriptions-item>
          </el-descriptions>
        </el-card>

        <el-card style="margin-top: 20px">
          <template #header>
            <span>相关合同</span>
          </template>
          <el-table :data="contracts" style="width: 100%">
            <el-table-column prop="contract_no" label="合同编号" width="100" />
            <el-table-column prop="contract_name" label="名称" />
            <el-table-column prop="status" label="状态" width="80">
              <template #default="{ row }">
                <el-tag :type="row.status === 'active' ? 'success' : 'info'" size="small">{{ row.status }}</el-tag>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { customerApi, interactionApi, contractApi } from '../../api'

const route = useRoute()
const customerId = Number(route.params.id)

const customer = ref<any>({})
const health = ref<any>({})
const interactions = ref<any[]>([])
const contracts = ref<any[]>([])
const showInteractionDialog = ref(false)

const healthColor = computed(() => {
  const score = health.value.overall_score || 0
  return score > 70 ? '#67c23a' : score > 40 ? '#e6a23c' : '#f56c6c'
})

const loadCustomer = async () => {
  try {
    const res = await customerApi.getDetail(customerId)
    customer.value = res.data.data
  } catch (error) {
    console.error('Failed to load customer:', error)
  }
}

const loadHealth = async () => {
  try {
    const res = await customerApi.getHealth(customerId)
    health.value = res.data.data || {}
  } catch (error) {
    console.error('Failed to load health:', error)
  }
}

const calculateHealth = async () => {
  try {
    const res = await customerApi.calculateHealth(customerId)
    health.value = res.data.data
    ElMessage.success('健康度计算完成')
  } catch (error) {
    console.error('Failed to calculate health:', error)
  }
}

const loadInteractions = async () => {
  try {
    const res = await interactionApi.getList({ customer_id: customerId })
    interactions.value = res.data.data || []
  } catch (error) {
    console.error('Failed to load interactions:', error)
  }
}

const loadContracts = async () => {
  try {
    const res = await contractApi.getList({ customer_id: customerId })
    contracts.value = res.data.data || []
  } catch (error) {
    console.error('Failed to load contracts:', error)
  }
}

onMounted(() => {
  loadCustomer()
  loadHealth()
  loadInteractions()
  loadContracts()
})
</script>

<style scoped>
.health-score {
  text-align: center;
  padding: 20px 0;
}

.score-label {
  margin-top: 10px;
  font-size: 16px;
  color: #666;
}
</style>
