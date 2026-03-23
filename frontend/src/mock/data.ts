export const mockCustomers = [
  {
    id: 1,
    customer_code: 'CUS120260115001',
    customer_name: '华为技术有限公司',
    industry: '通信设备',
    scale: 'large',
    region: '深圳',
    address: '广东省深圳市龙岗区坂田华为基地',
    website: 'https://www.huawei.com',
    status: 'active',
    trust_score: 85,
    commitment_score: 78,
    reciprocity_score: 82,
    overall_score: 82,
    created_by: 1,
    created_at: '2024-01-15T08:00:00Z',
    updated_at: '2024-03-10T10:30:00Z'
  },
  {
    id: 2,
    customer_code: 'CUS220260115002',
    customer_name: '腾讯科技有限公司',
    industry: '互联网',
    scale: 'large',
    region: '深圳',
    address: '广东省深圳市南山区科技园腾讯大厦',
    website: 'https://www.tencent.com',
    status: 'active',
    trust_score: 92,
    commitment_score: 88,
    reciprocity_score: 90,
    overall_score: 90,
    created_by: 1,
    created_at: '2024-02-20T09:00:00Z',
    updated_at: '2024-03-12T14:20:00Z'
  },
  {
    id: 3,
    customer_code: 'CUS320260115003',
    customer_name: '阿里巴巴集团',
    industry: '电子商务',
    scale: 'large',
    region: '杭州',
    address: '浙江省杭州市余杭区阿里巴巴西溪园区',
    website: 'https://www.alibaba.com',
    status: 'active',
    trust_score: 75,
    commitment_score: 70,
    reciprocity_score: 72,
    overall_score: 72,
    created_by: 1,
    created_at: '2024-01-10T10:00:00Z',
    updated_at: '2024-02-28T16:45:00Z'
  },
  {
    id: 4,
    customer_code: 'CUS420260115004',
    customer_name: '字节跳动科技有限公司',
    industry: '互联网',
    scale: 'large',
    region: '北京',
    address: '北京市海淀区北三环西路43号中航广场',
    website: 'https://www.bytedance.com',
    status: 'active',
    trust_score: 88,
    commitment_score: 85,
    reciprocity_score: 86,
    overall_score: 86,
    created_by: 2,
    created_at: '2024-03-01T11:00:00Z',
    updated_at: '2024-03-15T09:30:00Z'
  },
  {
    id: 5,
    customer_code: 'CUS520260115005',
    customer_name: '小米科技有限责任公司',
    industry: '消费电子',
    scale: 'large',
    region: '北京',
    address: '北京市海淀区清河中街68号华润五彩城',
    website: 'https://www.mi.com',
    status: 'active',
    trust_score: 68,
    commitment_score: 65,
    reciprocity_score: 66,
    overall_score: 66,
    created_by: 1,
    created_at: '2024-02-05T14:00:00Z',
    updated_at: '2024-03-08T11:15:00Z'
  }
]

export const mockContacts = [
  {
    id: 1,
    customer_id: 1,
    name: '张伟',
    position: '采购总监',
    department: '采购部',
    phone: '13800138001',
    email: 'zhangwei@huawei.com',
    wechat: 'zhangwei_hw',
    is_key_decision_maker: true,
    influence_level: 5,
    notes: '负责公司整体采购决策，对价格敏感，重视长期合作关系',
    created_at: '2024-01-20T08:00:00Z',
    updated_at: '2024-03-10T10:00:00Z'
  },
  {
    id: 2,
    customer_id: 1,
    name: '李娜',
    position: '技术总监',
    department: '研发部',
    phone: '13800138002',
    email: 'lina@huawei.com',
    wechat: 'lina_hw',
    is_key_decision_maker: true,
    influence_level: 4,
    notes: '技术选型关键人物，关注产品技术指标和稳定性',
    created_at: '2024-01-25T09:00:00Z',
    updated_at: '2024-03-05T14:00:00Z'
  },
  {
    id: 3,
    customer_id: 2,
    name: '王强',
    position: 'VP',
    department: '战略合作部',
    phone: '13800138003',
    email: 'wangqiang@tencent.com',
    wechat: 'wangqiang_tx',
    is_key_decision_maker: true,
    influence_level: 5,
    notes: '公司高层，直接参与重大合作决策',
    created_at: '2024-02-22T10:00:00Z',
    updated_at: '2024-03-12T15:00:00Z'
  },
  {
    id: 4,
    customer_id: 2,
    name: '陈明',
    position: '项目经理',
    department: '产品部',
    phone: '13800138004',
    email: 'chenming@tencent.com',
    wechat: 'chenming_tx',
    is_key_decision_maker: false,
    influence_level: 3,
    notes: '项目执行负责人，负责日常对接',
    created_at: '2024-02-25T11:00:00Z',
    updated_at: '2024-03-10T09:00:00Z'
  }
]

export const mockOrgStructures = [
  {
    id: 1,
    customer_id: 1,
    parent_id: null,
    dept_name: '华为技术有限公司',
    dept_function: '总部',
    order_num: 1,
    children: [
      {
        id: 2,
        customer_id: 1,
        parent_id: 1,
        dept_name: '运营商BG',
        dept_function: '运营商业务',
        order_num: 1,
        children: [
          { id: 5, customer_id: 1, parent_id: 2, dept_name: '采购部', dept_function: '采购管理', order_num: 1 },
          { id: 6, customer_id: 1, parent_id: 2, dept_name: '技术部', dept_function: '技术研发', order_num: 2 }
        ]
      },
      {
        id: 3,
        customer_id: 1,
        parent_id: 1,
        dept_name: '企业BG',
        dept_function: '企业业务',
        order_num: 2,
        children: [
          { id: 7, customer_id: 1, parent_id: 3, dept_name: '销售部', dept_function: '销售管理', order_num: 1 },
          { id: 8, customer_id: 1, parent_id: 3, dept_name: '服务部', dept_function: '客户服务', order_num: 2 }
        ]
      },
      {
        id: 4,
        customer_id: 1,
        parent_id: 1,
        dept_name: '消费者BG',
        dept_function: '消费者业务',
        order_num: 3,
        children: [
          { id: 9, customer_id: 1, parent_id: 4, dept_name: '产品部', dept_function: '产品管理', order_num: 1 },
          { id: 10, customer_id: 1, parent_id: 4, dept_name: '市场部', dept_function: '市场推广', order_num: 2 }
        ]
      }
    ]
  }
]

export const mockCooperationHistory = [
  {
    id: 1,
    customer_id: 1,
    project_name: '5G基站设备采购项目',
    contract_amount: 5000000,
    start_date: '2023-01-15',
    end_date: '2023-12-31',
    status: 'completed',
    satisfaction_score: 4.5,
    description: '为华为提供5G基站核心设备，项目按时交付，客户满意度高',
    created_at: '2023-01-10T08:00:00Z'
  },
  {
    id: 2,
    customer_id: 1,
    project_name: '云计算服务合作',
    contract_amount: 2000000,
    start_date: '2023-06-01',
    end_date: '2024-05-31',
    status: 'in_progress',
    satisfaction_score: 4.2,
    description: '提供云计算基础设施服务，目前合作顺利',
    created_at: '2023-05-25T10:00:00Z'
  },
  {
    id: 3,
    customer_id: 2,
    project_name: '游戏服务器采购',
    contract_amount: 3000000,
    start_date: '2023-03-01',
    end_date: '2024-02-28',
    status: 'completed',
    satisfaction_score: 4.8,
    description: '为腾讯游戏提供高性能服务器设备',
    created_at: '2023-02-20T09:00:00Z'
  },
  {
    id: 4,
    customer_id: 2,
    project_name: 'AI算力平台建设',
    contract_amount: 8000000,
    start_date: '2024-01-01',
    end_date: '2024-12-31',
    status: 'in_progress',
    satisfaction_score: 4.0,
    description: '协助腾讯建设AI算力基础设施平台',
    created_at: '2023-12-15T14:00:00Z'
  }
]

export const mockContracts = [
  {
    id: 1,
    contract_no: 'HT2024001',
    contract_name: '华为5G设备采购合同',
    customer_id: 1,
    contract_type: 'sales',
    amount: 5000000,
    currency: 'CNY',
    start_date: '2024-01-01',
    end_date: '2024-12-31',
    status: 'active',
    approval_status: 'approved',
    created_by: 1,
    sign_date: '2024-01-15',
    effective_date: '2024-01-01',
    created_at: '2024-01-10T08:00:00Z',
    updated_at: '2024-01-15T10:00:00Z'
  },
  {
    id: 2,
    contract_no: 'HT2024002',
    contract_name: '腾讯AI平台服务合同',
    customer_id: 2,
    contract_type: 'service',
    amount: 8000000,
    currency: 'CNY',
    start_date: '2024-01-01',
    end_date: '2024-12-31',
    status: 'active',
    approval_status: 'approved',
    created_by: 1,
    sign_date: '2024-01-20',
    effective_date: '2024-01-01',
    created_at: '2024-01-05T09:00:00Z',
    updated_at: '2024-01-20T11:00:00Z'
  },
  {
    id: 3,
    contract_no: 'HT2024003',
    contract_name: '阿里云服务采购合同',
    customer_id: 3,
    contract_type: 'purchase',
    amount: 2000000,
    currency: 'CNY',
    start_date: '2024-03-01',
    end_date: '2025-02-28',
    status: 'pending',
    approval_status: 'pending',
    created_by: 2,
    created_at: '2024-02-28T10:00:00Z',
    updated_at: '2024-02-28T10:00:00Z'
  },
  {
    id: 4,
    contract_no: 'HT2024004',
    contract_name: '字节跳动服务器采购框架协议',
    customer_id: 4,
    contract_type: 'sales',
    amount: 12000000,
    currency: 'CNY',
    start_date: '2024-02-01',
    end_date: '2025-01-31',
    status: 'active',
    approval_status: 'approved',
    created_by: 1,
    sign_date: '2024-02-10',
    effective_date: '2024-02-01',
    created_at: '2024-01-25T14:00:00Z',
    updated_at: '2024-02-10T16:00:00Z'
  },
  {
    id: 5,
    contract_no: 'HT2024005',
    contract_name: '小米智能硬件合作协议',
    customer_id: 5,
    contract_type: 'service',
    amount: 3500000,
    currency: 'CNY',
    start_date: '2023-06-01',
    end_date: '2024-05-31',
    status: 'expiring',
    approval_status: 'approved',
    created_by: 1,
    sign_date: '2023-05-20',
    effective_date: '2023-06-01',
    created_at: '2023-05-15T08:00:00Z',
    updated_at: '2023-05-20T10:00:00Z'
  }
]

export const mockContractTemplates = [
  {
    id: 1,
    template_name: '标准销售合同模板',
    template_type: 'sales',
    content: '本合同由甲方（买方）和乙方（卖方）于签订日期签署...',
    variables: '["甲方名称", "乙方名称", "合同金额", "签订日期", "生效日期", "终止日期"]',
    status: 'active',
    created_by: 1,
    created_at: '2024-01-01T08:00:00Z',
    updated_at: '2024-01-01T08:00:00Z'
  },
  {
    id: 2,
    template_name: '服务协议模板',
    template_type: 'service',
    content: '本服务协议由服务提供方和服务接受方于签订日期签署...',
    variables: '["服务内容", "服务期限", "服务费用", "付款方式"]',
    status: 'active',
    created_by: 1,
    created_at: '2024-01-01T08:00:00Z',
    updated_at: '2024-01-01T08:00:00Z'
  },
  {
    id: 3,
    template_name: '采购合同模板',
    template_type: 'purchase',
    content: '本采购合同由采购方和供应方于签订日期签署...',
    variables: '["采购物品", "数量", "单价", "总价", "交货日期"]',
    status: 'active',
    created_by: 1,
    created_at: '2024-01-01T08:00:00Z',
    updated_at: '2024-01-01T08:00:00Z'
  },
  {
    id: 4,
    template_name: '保密协议(NDA)模板',
    template_type: 'nda',
    content: '本保密协议由披露方和接收方于签订日期签署...',
    variables: '["保密信息范围", "保密期限", "违约责任"]',
    status: 'active',
    created_by: 1,
    created_at: '2024-01-01T08:00:00Z',
    updated_at: '2024-01-01T08:00:00Z'
  }
]

export const mockInteractions = [
  {
    id: 1,
    customer_id: 1,
    interaction_date: '2024-03-15',
    interaction_type: 'meeting',
    topic: 'Q2合作规划会议',
    summary: '与华为采购总监张伟讨论Q2合作计划，确认采购需求和技术支持要求',
    outcome: 'positive',
    importance: 5,
    next_action: '准备详细的技术方案和报价',
    created_by: 1,
    created_at: '2024-03-15T16:00:00Z'
  },
  {
    id: 2,
    customer_id: 1,
    interaction_date: '2024-03-10',
    interaction_type: 'phone',
    topic: '项目进度确认',
    summary: '电话确认云计算服务项目进度，客户对当前进展表示满意',
    outcome: 'positive',
    importance: 3,
    next_action: '发送项目进度报告',
    created_by: 1,
    created_at: '2024-03-10T14:00:00Z'
  },
  {
    id: 3,
    customer_id: 2,
    interaction_date: '2024-03-14',
    interaction_type: 'visit',
    topic: '高层互访',
    summary: '公司CEO带队拜访腾讯总部，与王强VP深入交流战略合作方向',
    outcome: 'positive',
    importance: 5,
    next_action: '准备战略合作框架协议',
    created_by: 1,
    created_at: '2024-03-14T18:00:00Z'
  },
  {
    id: 4,
    customer_id: 2,
    interaction_date: '2024-03-08',
    interaction_type: 'email',
    topic: '技术方案讨论',
    summary: '邮件往来讨论AI平台技术方案细节',
    outcome: 'neutral',
    importance: 4,
    next_action: '安排技术评审会议',
    created_by: 2,
    created_at: '2024-03-08T10:00:00Z'
  },
  {
    id: 5,
    customer_id: 3,
    interaction_date: '2024-03-05',
    interaction_type: 'meeting',
    topic: '续约谈判',
    summary: '与阿里团队讨论云服务合同续约事宜',
    outcome: 'negative',
    importance: 5,
    next_action: '准备更有竞争力的报价方案',
    created_by: 1,
    created_at: '2024-03-05T15:00:00Z'
  },
  {
    id: 6,
    customer_id: 4,
    interaction_date: '2024-03-12',
    interaction_type: 'phone',
    topic: '服务器交付确认',
    summary: '确认服务器交付时间和安装安排',
    outcome: 'positive',
    importance: 4,
    next_action: '安排物流和安装团队',
    created_by: 2,
    created_at: '2024-03-12T11:00:00Z'
  }
]

export const mockDashboardData = {
  overview: {
    total_customers: 156,
    active_contracts: 42,
    monthly_interactions: 328,
    expiring_contracts: 8
  },
  customerGrowth: [
    { month: '2024-01', new_customers: 12, lost_customers: 2, total: 145 },
    { month: '2024-02', new_customers: 8, lost_customers: 1, total: 152 },
    { month: '2024-03', new_customers: 15, lost_customers: 3, total: 164 },
    { month: '2024-04', new_customers: 10, lost_customers: 2, total: 172 },
    { month: '2024-05', new_customers: 18, lost_customers: 1, total: 189 },
    { month: '2024-06', new_customers: 12, lost_customers: 4, total: 197 }
  ],
  contractDistribution: [
    { range: '0-50万', count: 25, total_amount: 8500000 },
    { range: '50-100万', count: 12, total_amount: 8900000 },
    { range: '100-500万', count: 8, total_amount: 22000000 },
    { range: '500万以上', count: 5, total_amount: 45000000 }
  ],
  interactionHeatmap: [
    { date: '2024-03-01', count: 12 },
    { date: '2024-03-02', count: 8 },
    { date: '2024-03-03', count: 0 },
    { date: '2024-03-04', count: 15 },
    { date: '2024-03-05', count: 22 },
    { date: '2024-03-06', count: 18 },
    { date: '2024-03-07', count: 10 },
    { date: '2024-03-08', count: 14 },
    { date: '2024-03-09', count: 5 },
    { date: '2024-03-10', count: 0 },
    { date: '2024-03-11', count: 20 },
    { date: '2024-03-12', count: 25 },
    { date: '2024-03-13', count: 18 },
    { date: '2024-03-14', count: 16 },
    { date: '2024-03-15', count: 22 }
  ],
  renewalAlerts: [
    {
      contract_id: 5,
      contract_name: '小米智能硬件合作协议',
      customer_name: '小米科技有限责任公司',
      end_date: '2024-05-31',
      days_remaining: 77,
      renewal_probability: 0.65,
      predicted_amount: 3800000,
      risk_level: 'medium',
      last_interaction: '2024-02-20'
    },
    {
      contract_id: 6,
      contract_name: '京东物流服务合同',
      customer_name: '京东集团',
      end_date: '2024-04-30',
      days_remaining: 46,
      renewal_probability: 0.45,
      predicted_amount: 2500000,
      risk_level: 'high',
      last_interaction: '2024-01-15'
    },
    {
      contract_id: 7,
      contract_name: '美团云服务协议',
      customer_name: '美团',
      end_date: '2024-06-30',
      days_remaining: 107,
      renewal_probability: 0.85,
      predicted_amount: 4200000,
      risk_level: 'low',
      last_interaction: '2024-03-10'
    }
  ],
  topCustomers: [
    { id: 4, name: '字节跳动科技有限公司', total_amount: 12000000, health_score: 86, ltv: 48000000 },
    { id: 2, name: '腾讯科技有限公司', total_amount: 11000000, health_score: 90, ltv: 55000000 },
    { id: 1, name: '华为技术有限公司', total_amount: 7000000, health_score: 82, ltv: 35000000 },
    { id: 5, name: '小米科技有限责任公司', total_amount: 3500000, health_score: 66, ltv: 14000000 },
    { id: 3, name: '阿里巴巴集团', total_amount: 2000000, health_score: 72, ltv: 10000000 }
  ]
}

export const mockUsers = [
  {
    id: 1,
    username: 'admin',
    email: 'admin@enterprise-orbit.com',
    role: 'admin',
    status: 'active',
    created_at: '2024-01-01T00:00:00Z'
  },
  {
    id: 2,
    username: 'sales_manager',
    email: 'sales@enterprise-orbit.com',
    role: 'user',
    status: 'active',
    created_at: '2024-01-05T08:00:00Z'
  },
  {
    id: 3,
    username: 'contract_manager',
    email: 'contract@enterprise-orbit.com',
    role: 'user',
    status: 'active',
    created_at: '2024-01-10T09:00:00Z'
  }
]

export const mockApprovalFlow = [
  { level: 1, role: '部门经理', amount_threshold: 100000, description: '10万以下部门经理审批' },
  { level: 2, role: '法务审核', amount_threshold: 500000, description: '10-50万需法务审核' },
  { level: 3, role: 'CEO', amount_threshold: 999999999, description: '50万以上需CEO审批' }
]

export const mockPerformanceRecords = [
  {
    id: 1,
    contract_id: 1,
    milestone_name: '首付款支付',
    milestone_desc: '合同签订后支付30%首付款',
    planned_date: '2024-01-20',
    actual_date: '2024-01-18',
    completion_status: 'completed',
    amount_planned: 1500000,
    amount_paid: 1500000,
    notes: '客户提前支付，信用良好',
    created_at: '2024-01-15T08:00:00Z'
  },
  {
    id: 2,
    contract_id: 1,
    milestone_name: '设备交付',
    milestone_desc: '完成首批设备交付安装',
    planned_date: '2024-03-15',
    actual_date: null,
    completion_status: 'in_progress',
    amount_planned: 0,
    amount_paid: 0,
    notes: '正在准备发货',
    created_at: '2024-01-15T08:00:00Z'
  },
  {
    id: 3,
    contract_id: 1,
    milestone_name: '验收确认',
    milestone_desc: '客户验收并签署验收报告',
    planned_date: '2024-04-01',
    actual_date: null,
    completion_status: 'pending',
    amount_planned: 0,
    amount_paid: 0,
    notes: '',
    created_at: '2024-01-15T08:00:00Z'
  },
  {
    id: 4,
    contract_id: 1,
    milestone_name: '尾款支付',
    milestone_desc: '验收后支付70%尾款',
    planned_date: '2024-04-15',
    actual_date: null,
    completion_status: 'pending',
    amount_planned: 3500000,
    amount_paid: 0,
    notes: '',
    created_at: '2024-01-15T08:00:00Z'
  }
]
