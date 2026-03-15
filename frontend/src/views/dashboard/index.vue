<template>
  <div class="dashboard">
    <div class="page-header">
      <div class="page-title">
        <h1>{{ $t('dashboard.title') }}</h1>
        <p>{{ $t('dashboard.subtitle') }}</p>
      </div>
      <div class="page-actions">
        <el-date-picker
          v-model="dateRange"
          type="daterange"
          :range-separator="$t('common.search')"
          :start-placeholder="$t('contract.startDate')"
          :end-placeholder="$t('contract.endDate')"
          size="default"
        />
      </div>
    </div>

    <div class="stats-grid">
      <div v-for="(item, index) in statsCards" :key="item.key" class="stat-card" :style="{ animationDelay: `${index * 0.1}s` }">
        <div class="stat-card-content">
          <div class="stat-icon" :style="{ background: item.gradient }">
            <el-icon :size="24"><component :is="item.icon" /></el-icon>
          </div>
          <div class="stat-info">
            <span class="stat-label">{{ item.label }}</span>
            <span class="stat-value">{{ formatNumber(overview[item.key] || 0) }}</span>
          </div>
        </div>
        <div class="stat-trend" v-if="item.trend">
          <el-icon :class="item.trend > 0 ? 'trend-up' : 'trend-down'">
            <component :is="item.trend > 0 ? trendIcons.Top : trendIcons.Bottom" />
          </el-icon>
          <span :class="item.trend > 0 ? 'trend-up' : 'trend-down'">
            {{ Math.abs(item.trend) }}%
          </span>
        </div>
      </div>
    </div>

    <div class="charts-row">
      <div class="chart-card chart-large">
        <div class="chart-header">
          <h3>{{ $t('dashboard.customerGrowth') }}</h3>
          <div class="chart-legend">
            <span class="legend-item"><i class="legend-dot primary"></i>{{ $t('dashboard.newCustomers') }}</span>
            <span class="legend-item"><i class="legend-dot secondary"></i>{{ $t('dashboard.totalCustomersCount') }}</span>
          </div>
        </div>
        <div ref="customerChartRef" class="chart-container"></div>
      </div>
      <div class="chart-card chart-small">
        <div class="chart-header">
          <h3>{{ $t('dashboard.contractDistribution') }}</h3>
        </div>
        <div ref="contractChartRef" class="chart-container"></div>
      </div>
    </div>

    <div class="data-row">
      <div class="data-card">
        <div class="card-header">
          <h3>{{ $t('dashboard.renewalAlerts') }}</h3>
          <el-button type="primary" link>{{ $t('dashboard.viewAll') }}</el-button>
        </div>
        <div class="alert-list">
          <div v-for="item in renewalAlerts" :key="item.id" class="alert-item">
            <div class="alert-icon" :class="getRiskClass(item.days_until_expiry)">
              <el-icon><Warning /></el-icon>
            </div>
            <div class="alert-content">
              <div class="alert-title">{{ item.contract_no }} - {{ item.customer_name }}</div>
              <div class="alert-meta">
                <span class="alert-date">{{ $t('dashboard.daysRemaining') }} {{ item.days_until_expiry }} {{ $t('dashboard.days') }}</span>
                <el-tag :type="getRiskType(item.days_until_expiry)" size="small">
                  {{ getRiskLabel(item.days_until_expiry) }}
                </el-tag>
              </div>
            </div>
          </div>
          <div v-if="renewalAlerts.length === 0" class="empty-state">
            <el-icon :size="32"><CircleCheck /></el-icon>
            <p>{{ $t('dashboard.noRenewalRisk') }}</p>
          </div>
        </div>
      </div>

      <div class="data-card">
        <div class="card-header">
          <h3>{{ $t('dashboard.topCustomers') }}</h3>
          <el-button type="primary" link>{{ $t('dashboard.viewAll') }}</el-button>
        </div>
        <div class="top-customer-list">
          <div v-for="(item, index) in topCustomers" :key="item.id" class="customer-item">
            <div class="customer-rank">{{ index + 1 }}</div>
            <div class="customer-info">
              <div class="customer-name">{{ item.customer_name }}</div>
              <div class="customer-meta">
                <span>{{ item.total_contracts }} {{ $t('dashboard.contractsCount') }}</span>
                <span>¥{{ formatNumber(item.total_value) }}</span>
              </div>
            </div>
            <div class="customer-score">
              <el-progress 
                type="circle" 
                :percentage="item.health_score || 0" 
                :width="48"
                :stroke-width="4"
                :color="getHealthColor(item.health_score)"
              />
            </div>
          </div>
          <div v-if="topCustomers.length === 0" class="empty-state">
            <el-icon :size="32"><User /></el-icon>
            <p>{{ $t('dashboard.noCustomerData') }}</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import * as echarts from 'echarts'
import { useI18n } from 'vue-i18n'
import { dashboardApi } from '../../api'
import { User, Document, ChatDotRound, Warning, CircleCheck, Top, Bottom } from '@element-plus/icons-vue'

const { t } = useI18n()
const customerChartRef = ref()
const contractChartRef = ref()
const overview = ref<any>({})
const renewalAlerts = ref<any[]>([])
const topCustomers = ref<any[]>([])
const dateRange = ref<[Date, Date] | null>(null)

const trendIcons = { Top, Bottom }

const statsCards = computed(() => [
  { 
    key: 'total_customers', 
    label: t('dashboard.totalCustomers'), 
    icon: User, 
    gradient: 'linear-gradient(135deg, #06b6d4 0%, #0891b2 100%)',
    trend: 12
  },
  { 
    key: 'active_contracts', 
    label: t('dashboard.activeContracts'), 
    icon: Document, 
    gradient: 'linear-gradient(135deg, #10b981 0%, #059669 100%)',
    trend: 8
  },
  { 
    key: 'monthly_interactions', 
    label: t('dashboard.monthlyInteractions'), 
    icon: ChatDotRound, 
    gradient: 'linear-gradient(135deg, #f59e0b 0%, #d97706 100%)',
    trend: -3
  },
  { 
    key: 'expiring_contracts', 
    label: t('dashboard.expiringContracts'), 
    icon: Warning, 
    gradient: 'linear-gradient(135deg, #ef4444 0%, #dc2626 100%)',
    trend: null
  }
])

const formatNumber = (num: number) => {
  if (num >= 10000) {
    return (num / 10000).toFixed(1) + '万'
  }
  return num.toLocaleString()
}

const getRiskClass = (days: number) => {
  if (days < 7) return 'risk-high'
  if (days < 30) return 'risk-medium'
  return 'risk-low'
}

const getRiskType = (days: number) => {
  if (days < 7) return 'danger'
  if (days < 30) return 'warning'
  return 'success'
}

const getRiskLabel = (days: number) => {
  if (days < 7) return t('dashboard.urgent')
  if (days < 30) return t('dashboard.attention')
  return t('dashboard.normal')
}

const getHealthColor = (score: number) => {
  if (score >= 70) return '#10b981'
  if (score >= 40) return '#f59e0b'
  return '#ef4444'
}

const loadOverview = async () => {
  try {
    const res = await dashboardApi.getOverview()
    overview.value = res.data.data
  } catch (error) {
    console.error('Failed to load overview:', error)
  }
}

const loadRenewalAlerts = async () => {
  try {
    const res = await dashboardApi.getRenewalAlerts()
    renewalAlerts.value = res.data.data || []
  } catch (error) {
    console.error('Failed to load renewal alerts:', error)
  }
}

const loadTopCustomers = async () => {
  try {
    const res = await dashboardApi.getTopCustomers()
    topCustomers.value = res.data.data || []
  } catch (error) {
    console.error('Failed to load top customers:', error)
  }
}

const initCustomerChart = async () => {
  try {
    const res = await dashboardApi.getCustomerGrowth({})
    const chart = echarts.init(customerChartRef.value)
    const data = res.data.data || []
    
    chart.setOption({
      tooltip: {
        trigger: 'axis',
        backgroundColor: 'rgba(17, 24, 39, 0.9)',
        borderColor: 'rgba(255, 255, 255, 0.1)',
        textStyle: { color: '#f1f5f9' }
      },
      grid: {
        left: '3%',
        right: '4%',
        bottom: '3%',
        top: '10%',
        containLabel: true
      },
      xAxis: {
        type: 'category',
        data: data.map((d: any) => d.month),
        axisLine: { lineStyle: { color: 'rgba(255, 255, 255, 0.1)' } },
        axisLabel: { color: '#94a3b8' }
      },
      yAxis: {
        type: 'value',
        axisLine: { show: false },
        splitLine: { lineStyle: { color: 'rgba(255, 255, 255, 0.05)' } },
        axisLabel: { color: '#94a3b8' }
      },
      series: [
        {
          name: t('dashboard.newCustomers'),
          type: 'bar',
          data: data.map((d: any) => d.new),
          itemStyle: {
            borderRadius: [4, 4, 0, 0],
            color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
              { offset: 0, color: '#06b6d4' },
              { offset: 1, color: '#0891b2' }
            ])
          }
        },
        {
          name: t('dashboard.totalCustomersCount'),
          type: 'line',
          smooth: true,
          data: data.map((d: any) => d.total),
          lineStyle: { color: '#8b5cf6', width: 2 },
          areaStyle: {
            color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
              { offset: 0, color: 'rgba(139, 92, 246, 0.3)' },
              { offset: 1, color: 'rgba(139, 92, 246, 0)' }
            ])
          },
          symbol: 'circle',
          symbolSize: 6,
          itemStyle: { color: '#8b5cf6' }
        }
      ]
    })
  } catch (error) {
    console.error('Failed to init customer chart:', error)
  }
}

const initContractChart = async () => {
  try {
    const res = await dashboardApi.getContractDistribution()
    const chart = echarts.init(contractChartRef.value)
    const data = res.data.data || []
    
    chart.setOption({
      tooltip: {
        trigger: 'item',
        backgroundColor: 'rgba(17, 24, 39, 0.9)',
        borderColor: 'rgba(255, 255, 255, 0.1)',
        textStyle: { color: '#f1f5f9' }
      },
      series: [{
        type: 'pie',
        radius: ['50%', '75%'],
        center: ['50%', '50%'],
        avoidLabelOverlap: false,
        itemStyle: {
          borderRadius: 8,
          borderColor: 'rgba(17, 24, 39, 1)',
          borderWidth: 2
        },
        label: {
          show: true,
          position: 'outside',
          color: '#94a3b8',
          formatter: '{b}\n{d}%'
        },
        labelLine: {
          lineStyle: { color: 'rgba(255, 255, 255, 0.2)' }
        },
        data: data.map((d: any, i: number) => ({
          name: d.status,
          value: d.amount,
          itemStyle: {
            color: ['#06b6d4', '#8b5cf6', '#10b981', '#f59e0b'][i % 4]
          }
        }))
      }]
    })
  } catch (error) {
    console.error('Failed to init contract chart:', error)
  }
}

onMounted(() => {
  loadOverview()
  loadRenewalAlerts()
  loadTopCustomers()
  setTimeout(() => {
    initCustomerChart()
    initContractChart()
  }, 100)
})
</script>

<style scoped>
.dashboard {
  max-width: 1400px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 32px;
}

.page-title h1 {
  font-family: var(--font-display);
  font-size: 28px;
  font-weight: 700;
  color: var(--color-text-primary);
  margin: 0 0 8px 0;
}

.page-title p {
  font-size: 14px;
  color: var(--color-text-muted);
  margin: 0;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 20px;
  margin-bottom: 24px;
}

.stat-card {
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  padding: 24px;
  animation: fadeIn 0.5s ease-out forwards;
  opacity: 0;
}

.stat-card-content {
  display: flex;
  align-items: center;
  gap: 16px;
}

.stat-icon {
  width: 56px;
  height: 56px;
  border-radius: var(--radius-sm);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.stat-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.stat-label {
  font-size: 13px;
  color: var(--color-text-muted);
}

.stat-value {
  font-family: var(--font-display);
  font-size: 28px;
  font-weight: 700;
  color: var(--color-text-primary);
}

.stat-trend {
  display: flex;
  align-items: center;
  gap: 4px;
  margin-top: 12px;
  padding-top: 12px;
  border-top: 1px solid var(--color-border);
  font-size: 13px;
}

.trend-up {
  color: var(--color-success);
}

.trend-down {
  color: var(--color-danger);
}

.charts-row {
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: 20px;
  margin-bottom: 24px;
}

.chart-card {
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  padding: 24px;
}

.chart-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.chart-header h3 {
  font-family: var(--font-display);
  font-size: 16px;
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0;
}

.chart-legend {
  display: flex;
  gap: 16px;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: var(--color-text-muted);
}

.legend-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
}

.legend-dot.primary {
  background: var(--color-accent);
}

.legend-dot.secondary {
  background: var(--color-accent-secondary);
}

.chart-container {
  height: 280px;
}

.data-row {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 20px;
}

.data-card {
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  padding: 24px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.card-header h3 {
  font-family: var(--font-display);
  font-size: 16px;
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0;
}

.alert-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.alert-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background: var(--color-bg-tertiary);
  border-radius: var(--radius-sm);
  transition: all var(--transition-fast);
}

.alert-item:hover {
  background: rgba(255, 255, 255, 0.05);
}

.alert-icon {
  width: 40px;
  height: 40px;
  border-radius: var(--radius-sm);
  display: flex;
  align-items: center;
  justify-content: center;
}

.alert-icon.risk-high {
  background: rgba(239, 68, 68, 0.2);
  color: var(--color-danger);
}

.alert-icon.risk-medium {
  background: rgba(245, 158, 11, 0.2);
  color: var(--color-warning);
}

.alert-icon.risk-low {
  background: rgba(16, 185, 129, 0.2);
  color: var(--color-success);
}

.alert-content {
  flex: 1;
  min-width: 0;
}

.alert-title {
  font-size: 14px;
  font-weight: 500;
  color: var(--color-text-primary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.alert-meta {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 4px;
}

.alert-date {
  font-size: 12px;
  color: var(--color-text-muted);
}

.top-customer-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.customer-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background: var(--color-bg-tertiary);
  border-radius: var(--radius-sm);
  transition: all var(--transition-fast);
}

.customer-item:hover {
  background: rgba(255, 255, 255, 0.05);
}

.customer-rank {
  width: 28px;
  height: 28px;
  background: linear-gradient(135deg, var(--color-accent) 0%, var(--color-accent-secondary) 100%);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 700;
  color: white;
}

.customer-info {
  flex: 1;
  min-width: 0;
}

.customer-name {
  font-size: 14px;
  font-weight: 500;
  color: var(--color-text-primary);
}

.customer-meta {
  display: flex;
  gap: 12px;
  margin-top: 4px;
  font-size: 12px;
  color: var(--color-text-muted);
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
  color: var(--color-text-muted);
}

.empty-state p {
  margin-top: 12px;
  font-size: 14px;
}

@media (max-width: 1200px) {
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }
  
  .charts-row {
    grid-template-columns: 1fr;
  }
  
  .data-row {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .stats-grid {
    grid-template-columns: 1fr;
  }
  
  .page-header {
    flex-direction: column;
    gap: 16px;
  }
}
</style>
