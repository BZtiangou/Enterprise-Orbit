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
          value-format="YYYY-MM-DD"
        />
      </div>
    </div>

    <div class="stats-grid">
      <div v-for="(item, index) in statsCards" :key="item.key" class="stat-card glass" :style="{ animationDelay: `${index * 0.1}s` }">
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
      <div class="chart-card chart-large glass">
        <div class="chart-header">
          <h3>{{ $t('dashboard.customerGrowthTrend', '客户增长趋势') }}</h3>
          <div class="chart-actions">
            <el-radio-group v-model="growthPeriod" size="small">
              <el-radio-button label="month">{{ $t('dashboard.monthly', '月度') }}</el-radio-button>
              <el-radio-button label="quarter">{{ $t('dashboard.quarterly', '季度') }}</el-radio-button>
            </el-radio-group>
          </div>
        </div>
        <div ref="customerChartRef" class="chart-container"></div>
      </div>
      <div class="chart-card chart-small glass">
        <div class="chart-header">
          <h3>{{ $t('dashboard.contractValueDistribution', '合同价值分布') }}</h3>
        </div>
        <div ref="contractChartRef" class="chart-container"></div>
      </div>
    </div>

    <div class="charts-row">
      <div class="chart-card glass">
        <div class="chart-header">
          <h3>{{ $t('dashboard.interactionHeatmap', '互动热力图') }}</h3>
          <span class="chart-subtitle">{{ $t('dashboard.last30Days', '近30天互动频次') }}</span>
        </div>
        <div class="heatmap-container">
          <div class="heatmap-grid">
            <div
              v-for="(item, index) in interactionHeatmap"
              :key="index"
              class="heatmap-cell"
              :style="{ backgroundColor: getHeatmapColor(item.count) }"
              :title="`${item.date}: ${item.count} ${t('dashboard.interactions')}`"
            >
              <span v-if="item.count > 0" class="heatmap-count">{{ item.count }}</span>
            </div>
          </div>
          <div class="heatmap-legend">
            <span>{{ $t('dashboard.less', '少') }}</span>
            <div class="legend-scale">
              <span v-for="i in 5" :key="i" :style="{ backgroundColor: getHeatmapColor((i - 1) * 6) }"></span>
            </div>
            <span>{{ $t('dashboard.more', '多') }}</span>
          </div>
        </div>
      </div>

      <div class="chart-card glass">
        <div class="chart-header">
          <h3>{{ $t('dashboard.renewalRiskWarning', '续约风险预警') }}</h3>
          <el-button type="primary" link>{{ $t('dashboard.viewAll') }}</el-button>
        </div>
        <div class="renewal-list">
          <div v-for="item in renewalAlerts" :key="item.contract_id" class="renewal-item" :class="'risk-' + item.risk_level">
            <div class="renewal-left">
              <div class="renewal-customer">{{ item.customer_name }}</div>
              <div class="renewal-contract">{{ item.contract_name }}</div>
              <div class="renewal-meta">
                <span>{{ $t('dashboard.daysRemaining', '剩余') }}: {{ item.days_remaining }}{{ $t('dashboard.days', '天') }}</span>
                <span>{{ $t('dashboard.lastInteraction', '上次互动') }}: {{ item.last_interaction }}</span>
              </div>
            </div>
            <div class="renewal-right">
              <div class="renewal-probability">
                <el-progress type="dashboard" :percentage="Math.round(item.renewal_probability * 100)" :width="60" :color="getProbabilityColor(item.renewal_probability)" />
              </div>
              <div class="renewal-amount">¥{{ formatNumber(item.predicted_amount) }}</div>
            </div>
          </div>
          <div v-if="renewalAlerts.length === 0" class="empty-state">
            <el-icon :size="32"><CircleCheck /></el-icon>
            <p>{{ $t('dashboard.noRenewalRisk') }}</p>
          </div>
        </div>
      </div>
    </div>

    <div class="data-row">
      <div class="data-card glass">
        <div class="card-header">
          <h3>{{ $t('dashboard.topCustomersLTV', '高价值客户 (LTV预测)') }}</h3>
          <el-button type="primary" link>{{ $t('dashboard.viewAll') }}</el-button>
        </div>
        <div class="top-customer-list">
          <div v-for="(item, index) in topCustomers" :key="item.id" class="customer-item">
            <div class="customer-rank" :class="'rank-' + (index + 1)">{{ index + 1 }}</div>
            <div class="customer-info">
              <div class="customer-name">{{ item.name }}</div>
              <div class="customer-meta">
                <span>{{ $t('dashboard.currentValue', '当前价值') }}: ¥{{ formatNumber(item.total_amount) }}</span>
                <span>{{ $t('dashboard.healthScore', '健康度') }}: {{ item.health_score }}</span>
              </div>
            </div>
            <div class="customer-ltv">
              <div class="ltv-label">{{ $t('dashboard.ltv', 'LTV') }}</div>
              <div class="ltv-value">¥{{ formatNumber(item.ltv) }}</div>
            </div>
          </div>
          <div v-if="topCustomers.length === 0" class="empty-state">
            <el-icon :size="32"><User /></el-icon>
            <p>{{ $t('dashboard.noCustomerData') }}</p>
          </div>
        </div>
      </div>

      <div class="data-card glass">
        <div class="card-header">
          <h3>{{ $t('dashboard.teamPerformance', '团队绩效对比') }}</h3>
        </div>
        <div ref="teamChartRef" class="chart-container" style="height: 280px;"></div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue'
import * as echarts from 'echarts'
import { useI18n } from 'vue-i18n'
import { dashboardApi } from '../../api'
import { mockDashboardData } from '../../mock/data'
import { User, Document, ChatDotRound, Warning, CircleCheck, Top, Bottom, TrendCharts } from '@element-plus/icons-vue'

const { t } = useI18n()
const customerChartRef = ref()
const contractChartRef = ref()
const teamChartRef = ref()
const overview = ref<any>({})
const renewalAlerts = ref<any[]>([])
const topCustomers = ref<any[]>([])
const interactionHeatmap = ref<any[]>([])
const dateRange = ref<[string, string] | null>(null)
const growthPeriod = ref('month')

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
  if (num >= 100000000) return (num / 100000000).toFixed(1) + t('dashboard.billion')
  if (num >= 10000) return (num / 10000).toFixed(1) + t('dashboard.tenThousand')
  return num.toLocaleString()
}

const getHeatmapColor = (count: number) => {
  const colors = [
    'rgba(100, 116, 139, 0.1)',
    'rgba(6, 182, 212, 0.2)',
    'rgba(6, 182, 212, 0.4)',
    'rgba(6, 182, 212, 0.6)',
    'rgba(6, 182, 212, 0.8)',
    'rgba(6, 182, 212, 1)'
  ]
  const index = Math.min(Math.floor(count / 5), 5)
  return colors[index]
}

const getProbabilityColor = (probability: number) => {
  if (probability >= 0.7) return '#10b981'
  if (probability >= 0.4) return '#f59e0b'
  return '#ef4444'
}

const loadOverview = async () => {
  try {
    const res = await dashboardApi.getOverview()
    overview.value = res.data.data
  } catch (error) {
    overview.value = mockDashboardData.overview
  }
}

const loadRenewalAlerts = async () => {
  try {
    const res = await dashboardApi.getRenewalAlerts()
    renewalAlerts.value = res.data.data || []
  } catch (error) {
    renewalAlerts.value = mockDashboardData.renewalAlerts
  }
}

const loadTopCustomers = async () => {
  try {
    const res = await dashboardApi.getTopCustomers()
    topCustomers.value = res.data.data || []
  } catch (error) {
    topCustomers.value = mockDashboardData.topCustomers
  }
}

const loadInteractionHeatmap = async () => {
  try {
    const res = await dashboardApi.getInteractionHeatmap()
    interactionHeatmap.value = res.data.data || []
  } catch (error) {
    interactionHeatmap.value = mockDashboardData.interactionHeatmap
  }
}

const initCustomerChart = async () => {
  const chart = echarts.init(customerChartRef.value)
  let data: any[] = []
  
  try {
    const res = await dashboardApi.getCustomerGrowth({})
    data = res.data.data || []
  } catch (error) {
    data = mockDashboardData.customerGrowth
  }
  
  chart.setOption({
    tooltip: {
      trigger: 'axis',
      backgroundColor: 'rgba(17, 24, 39, 0.9)',
      borderColor: 'rgba(255, 255, 255, 0.1)',
      textStyle: { color: '#f1f5f9' }
    },
    legend: {
      data: [t('dashboard.newCustomers'), t('dashboard.lostCustomers', '流失客户'), t('dashboard.netGrowth', '净增长')],
      textStyle: { color: '#94a3b8' },
      top: 0
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      top: '15%',
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
        stack: 'total',
        data: data.map((d: any) => d.new_customers),
        itemStyle: {
          borderRadius: [0, 0, 0, 0],
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: '#10b981' },
            { offset: 1, color: '#059669' }
          ])
        }
      },
      {
        name: t('dashboard.lostCustomers', '流失客户'),
        type: 'bar',
        stack: 'total',
        data: data.map((d: any) => -d.lost_customers),
        itemStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: '#ef4444' },
            { offset: 1, color: '#dc2626' }
          ])
        }
      },
      {
        name: t('dashboard.netGrowth', '净增长'),
        type: 'line',
        data: data.map((d: any) => d.new_customers - d.lost_customers),
        lineStyle: { color: '#06b6d4', width: 3 },
        symbol: 'circle',
        symbolSize: 8,
        itemStyle: { color: '#06b6d4' }
      }
    ]
  })
}

const initContractChart = async () => {
  const chart = echarts.init(contractChartRef.value)
  let data: any[] = []
  
  try {
    const res = await dashboardApi.getContractDistribution()
    data = res.data.data || []
  } catch (error) {
    data = mockDashboardData.contractDistribution
  }
  
  chart.setOption({
    tooltip: {
      trigger: 'item',
      backgroundColor: 'rgba(17, 24, 39, 0.9)',
      borderColor: 'rgba(255, 255, 255, 0.1)',
      textStyle: { color: '#f1f5f9' },
      formatter: (params: any) => {
        return `${params.name}<br/>${t('contract.amount')}: ¥${formatNumber(params.value)}<br/>${t('dashboard.contractsCount')}: ${data[params.dataIndex]?.count || 0}`
      }
    },
    series: [{
      type: 'pie',
      radius: ['45%', '70%'],
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
        name: d.range,
        value: d.total_amount,
        itemStyle: {
          color: ['#06b6d4', '#8b5cf6', '#10b981', '#f59e0b'][i % 4]
        }
      }))
    }]
  })
}

const initTeamChart = async () => {
  const chart = echarts.init(teamChartRef.value)
  
  const teamData = [
    { name: '销售一组', healthAvg: 82, renewalRate: 78, interactionCount: 156 },
    { name: '销售二组', healthAvg: 75, renewalRate: 72, interactionCount: 132 },
    { name: '销售三组', healthAvg: 88, renewalRate: 85, interactionCount: 178 },
    { name: '销售四组', healthAvg: 70, renewalRate: 68, interactionCount: 98 }
  ]
  
  chart.setOption({
    tooltip: {
      trigger: 'axis',
      backgroundColor: 'rgba(17, 24, 39, 0.9)',
      borderColor: 'rgba(255, 255, 255, 0.1)',
      textStyle: { color: '#f1f5f9' }
    },
    legend: {
      data: [t('dashboard.healthAvg', '健康度均值'), t('dashboard.renewalRate', '续约率'), t('dashboard.interactionCount', '互动次数')],
      textStyle: { color: '#94a3b8' },
      top: 0
    },
    radar: {
      indicator: [
        { name: t('dashboard.healthAvg', '健康度'), max: 100 },
        { name: t('dashboard.renewalRate', '续约率'), max: 100 },
        { name: t('dashboard.interactionCount', '互动'), max: 200 }
      ],
      axisLine: { lineStyle: { color: 'rgba(255, 255, 255, 0.1)' } },
      splitLine: { lineStyle: { color: 'rgba(255, 255, 255, 0.05)' } },
      splitArea: { show: false }
    },
    series: [{
      type: 'radar',
      data: teamData.map((team, i) => ({
        name: team.name,
        value: [team.healthAvg, team.renewalRate, team.interactionCount],
        lineStyle: { color: ['#06b6d4', '#8b5cf6', '#10b981', '#f59e0b'][i] },
        areaStyle: { color: ['rgba(6, 182, 212, 0.2)', 'rgba(139, 92, 246, 0.2)', 'rgba(16, 185, 129, 0.2)', 'rgba(245, 158, 11, 0.2)'][i] }
      }))
    }]
  })
}

watch(growthPeriod, () => {
  initCustomerChart()
})

onMounted(() => {
  loadOverview()
  loadRenewalAlerts()
  loadTopCustomers()
  loadInteractionHeatmap()
  setTimeout(() => {
    initCustomerChart()
    initContractChart()
    initTeamChart()
  }, 100)
})
</script>

<style scoped>
.dashboard { max-width: 1400px; }
.page-header { display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 32px; }
.page-title h1 { font-family: var(--font-display); font-size: 28px; font-weight: 700; background: linear-gradient(135deg, var(--color-accent) 0%, var(--color-accent-secondary) 100%); -webkit-background-clip: text; -webkit-text-fill-color: transparent; background-clip: text; margin: 0 0 8px 0; }
.page-title p { font-size: 14px; color: var(--color-text-muted); margin: 0; }
.stats-grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 20px; margin-bottom: 24px; }
.stat-card { padding: 24px; border-radius: var(--radius-md); animation: fadeIn 0.5s ease-out forwards; opacity: 0; position: relative; overflow: hidden; transition: all var(--transition-normal); }
.stat-card:hover { transform: translateY(-4px); box-shadow: 0 8px 30px rgba(0, 0, 0, 0.3); }
.stat-card::before { content: ''; position: absolute; top: 0; left: 0; right: 0; height: 3px; background: linear-gradient(90deg, var(--color-accent), var(--color-accent-secondary)); }
.stat-card-content { display: flex; align-items: center; gap: 16px; }
.stat-icon { width: 56px; height: 56px; border-radius: var(--radius-sm); display: flex; align-items: center; justify-content: center; color: white; box-shadow: 0 4px 15px rgba(0, 0, 0, 0.3); }
.stat-info { display: flex; flex-direction: column; gap: 4px; }
.stat-label { font-size: 13px; color: var(--color-text-muted); }
.stat-value { font-family: var(--font-display); font-size: 28px; font-weight: 700; color: var(--color-text-primary); }
.stat-trend { display: flex; align-items: center; gap: 4px; margin-top: 12px; padding-top: 12px; border-top: 1px solid var(--color-border); font-size: 13px; }
.trend-up { color: var(--color-success); }
.trend-down { color: var(--color-danger); }
.charts-row { display: grid; grid-template-columns: 2fr 1fr; gap: 20px; margin-bottom: 24px; }
.chart-card { padding: 24px; border-radius: var(--radius-md); position: relative; overflow: hidden; }
.chart-card::before { content: ''; position: absolute; top: 0; left: 0; right: 0; height: 2px; background: linear-gradient(90deg, var(--color-accent), var(--color-accent-secondary)); }
.chart-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; }
.chart-header h3 { font-family: var(--font-display); font-size: 16px; font-weight: 600; color: var(--color-text-primary); margin: 0; }
.chart-subtitle { font-size: 12px; color: var(--color-text-muted); }
.chart-actions { display: flex; gap: 12px; }
.chart-container { height: 280px; }
.heatmap-container { padding: 10px 0; }
.heatmap-grid { display: grid; grid-template-columns: repeat(15, 1fr); gap: 4px; }
.heatmap-cell { aspect-ratio: 1; border-radius: 4px; display: flex; align-items: center; justify-content: center; font-size: 10px; color: var(--color-text-primary); cursor: pointer; transition: all 0.2s; }
.heatmap-cell:hover { transform: scale(1.2); box-shadow: 0 2px 8px rgba(0, 0, 0, 0.3); }
.heatmap-count { font-weight: 600; }
.heatmap-legend { display: flex; align-items: center; justify-content: flex-end; gap: 8px; margin-top: 12px; font-size: 12px; color: var(--color-text-muted); }
.legend-scale { display: flex; gap: 2px; }
.legend-scale span { width: 16px; height: 12px; border-radius: 2px; }
.renewal-list { display: flex; flex-direction: column; gap: 12px; max-height: 280px; overflow-y: auto; }
.renewal-item { display: flex; justify-content: space-between; padding: 16px; background: linear-gradient(135deg, rgba(255, 255, 255, 0.03) 0%, rgba(255, 255, 255, 0.01) 100%); border-radius: var(--radius-sm); border-left: 3px solid; transition: all var(--transition-fast); }
.renewal-item:hover { background: rgba(255, 255, 255, 0.05); }
.renewal-item.risk-high { border-left-color: var(--color-danger); }
.renewal-item.risk-medium { border-left-color: var(--color-warning); }
.renewal-item.risk-low { border-left-color: var(--color-success); }
.renewal-left { flex: 1; }
.renewal-customer { font-weight: 500; color: var(--color-text-primary); }
.renewal-contract { font-size: 12px; color: var(--color-text-secondary); margin-top: 4px; }
.renewal-meta { font-size: 11px; color: var(--color-text-muted); margin-top: 8px; display: flex; gap: 16px; }
.renewal-right { text-align: center; }
.renewal-amount { font-family: var(--font-display); font-size: 14px; font-weight: 600; color: var(--color-accent); margin-top: 8px; }
.data-row { display: grid; grid-template-columns: repeat(2, 1fr); gap: 20px; }
.data-card { padding: 24px; border-radius: var(--radius-md); position: relative; overflow: hidden; }
.data-card::before { content: ''; position: absolute; top: 0; left: 0; right: 0; height: 2px; background: linear-gradient(90deg, var(--color-accent), var(--color-accent-secondary)); }
.card-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; }
.card-header h3 { font-family: var(--font-display); font-size: 16px; font-weight: 600; color: var(--color-text-primary); margin: 0; }
.top-customer-list { display: flex; flex-direction: column; gap: 12px; }
.customer-item { display: flex; align-items: center; gap: 12px; padding: 16px; background: linear-gradient(135deg, rgba(255, 255, 255, 0.03) 0%, rgba(255, 255, 255, 0.01) 100%); border-radius: var(--radius-sm); transition: all var(--transition-fast); border: 1px solid transparent; }
.customer-item:hover { background: rgba(255, 255, 255, 0.05); border-color: var(--color-accent); }
.customer-rank { width: 28px; height: 28px; border-radius: 50%; display: flex; align-items: center; justify-content: center; font-size: 12px; font-weight: 700; color: white; background: var(--color-bg-secondary); }
.customer-rank.rank-1 { background: linear-gradient(135deg, #f59e0b, #d97706); box-shadow: 0 2px 10px rgba(245, 158, 11, 0.4); }
.customer-rank.rank-2 { background: linear-gradient(135deg, #94a3b8, #64748b); box-shadow: 0 2px 10px rgba(148, 163, 184, 0.4); }
.customer-rank.rank-3 { background: linear-gradient(135deg, #b45309, #92400e); box-shadow: 0 2px 10px rgba(180, 83, 9, 0.4); }
.customer-info { flex: 1; min-width: 0; }
.customer-name { font-size: 14px; font-weight: 500; color: var(--color-text-primary); }
.customer-meta { display: flex; gap: 12px; margin-top: 4px; font-size: 12px; color: var(--color-text-muted); }
.customer-ltv { text-align: right; }
.ltv-label { font-size: 11px; color: var(--color-text-muted); }
.ltv-value { font-family: var(--font-display); font-size: 16px; font-weight: 600; color: var(--color-accent); }
.empty-state { display: flex; flex-direction: column; align-items: center; justify-content: center; padding: 40px 20px; color: var(--color-text-muted); }
.empty-state p { margin-top: 12px; font-size: 14px; }
@keyframes fadeIn { from { opacity: 0; transform: translateY(10px); } to { opacity: 1; transform: translateY(0); } }
@media (max-width: 1200px) {
  .stats-grid { grid-template-columns: repeat(2, 1fr); }
  .charts-row { grid-template-columns: 1fr; }
  .data-row { grid-template-columns: 1fr; }
}
@media (max-width: 768px) {
  .stats-grid { grid-template-columns: 1fr; }
  .page-header { flex-direction: column; gap: 16px; }
}
</style>
