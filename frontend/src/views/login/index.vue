<template>
  <div class="auth-container">
    <div class="auth-background">
      <div class="grid-overlay"></div>
      <div class="gradient-orb orb-1"></div>
      <div class="gradient-orb orb-2"></div>
      <div class="gradient-orb orb-3"></div>
    </div>
    
    <div class="auth-content">
      <div class="auth-card glass" :class="{ 'flipped': isRegister }">
        <div class="auth-form-container">
          <div class="auth-header">
            <div class="logo-container animate-float">
              <div class="logo-icon">
                <svg viewBox="0 0 40 40" fill="none" xmlns="http://www.w3.org/2000/svg">
                  <circle cx="20" cy="20" r="18" stroke="url(#gradient)" stroke-width="2"/>
                  <circle cx="20" cy="20" r="8" fill="url(#gradient)"/>
                  <circle cx="20" cy="8" r="3" fill="#06b6d4"/>
                  <circle cx="32" cy="20" r="3" fill="#8b5cf6"/>
                  <circle cx="20" cy="32" r="3" fill="#06b6d4"/>
                  <circle cx="8" cy="20" r="3" fill="#8b5cf6"/>
                  <defs>
                    <linearGradient id="gradient" x1="0" y1="0" x2="40" y2="40">
                      <stop offset="0%" stop-color="#06b6d4"/>
                      <stop offset="100%" stop-color="#8b5cf6"/>
                    </linearGradient>
                  </defs>
                </svg>
              </div>
            </div>
            <h1 class="auth-title">Enterprise Orbit</h1>
            <p class="auth-subtitle">B2B 关系管理系统</p>
          </div>

          <transition name="fade" mode="out-in">
            <div v-if="!isRegister" key="login" class="auth-form">
              <el-form :model="loginForm" :rules="loginRules" ref="loginFormRef">
                <el-form-item prop="username">
                  <el-input
                    v-model="loginForm.username"
                    placeholder="用户名"
                    size="large"
                    :prefix-icon="User"
                  />
                </el-form-item>
                <el-form-item prop="password">
                  <el-input
                    v-model="loginForm.password"
                    type="password"
                    placeholder="密码"
                    size="large"
                    :prefix-icon="Lock"
                    show-password
                    @keyup.enter="handleLogin"
                  />
                </el-form-item>
                <el-form-item>
                  <el-button
                    type="primary"
                    size="large"
                    :loading="loading"
                    @click="handleLogin"
                    class="auth-button"
                  >
                    <span v-if="!loading">登录系统</span>
                    <span v-else>正在验证...</span>
                  </el-button>
                </el-form-item>
              </el-form>
              
              <div class="auth-footer">
                <span class="auth-footer-text">还没有账号？</span>
                <a class="auth-link" @click="toggleMode">立即注册</a>
              </div>
            </div>

            <div v-else key="register" class="auth-form">
              <el-form :model="registerForm" :rules="registerRules" ref="registerFormRef">
                <el-form-item prop="username">
                  <el-input
                    v-model="registerForm.username"
                    placeholder="用户名"
                    size="large"
                    :prefix-icon="User"
                  />
                </el-form-item>
                <el-form-item prop="email">
                  <el-input
                    v-model="registerForm.email"
                    placeholder="邮箱地址"
                    size="large"
                    :prefix-icon="Message"
                  />
                </el-form-item>
                <el-form-item prop="password">
                  <el-input
                    v-model="registerForm.password"
                    type="password"
                    placeholder="设置密码（至少6位）"
                    size="large"
                    :prefix-icon="Lock"
                    show-password
                  />
                </el-form-item>
                <el-form-item prop="confirmPassword">
                  <el-input
                    v-model="registerForm.confirmPassword"
                    type="password"
                    placeholder="确认密码"
                    size="large"
                    :prefix-icon="Lock"
                    show-password
                    @keyup.enter="handleRegister"
                  />
                </el-form-item>
                <el-form-item>
                  <el-button
                    type="primary"
                    size="large"
                    :loading="loading"
                    @click="handleRegister"
                    class="auth-button"
                  >
                    <span v-if="!loading">创建账号</span>
                    <span v-else>正在创建...</span>
                  </el-button>
                </el-form-item>
              </el-form>
              
              <div class="auth-footer">
                <span class="auth-footer-text">已有账号？</span>
                <a class="auth-link" @click="toggleMode">返回登录</a>
              </div>
            </div>
          </transition>
        </div>
      </div>

      <div class="auth-info">
        <div class="info-content">
          <h2 class="info-title">
            <span class="gradient-text">智能化的</span>
            <br />
            企业关系管理
          </h2>
          <p class="info-description">
            基于社会交换理论与动态能力理论，构建信任、互惠、承诺三维健康度评估模型，
            助力企业实现客户关系的精细化运营与价值最大化。
          </p>
          <div class="features-grid">
            <div class="feature-item">
              <div class="feature-icon">
                <el-icon :size="24"><User /></el-icon>
              </div>
              <div class="feature-text">
                <h4>客户全景档案</h4>
                <p>360°客户视图</p>
              </div>
            </div>
            <div class="feature-item">
              <div class="feature-icon">
                <el-icon :size="24"><Document /></el-icon>
              </div>
              <div class="feature-text">
                <h4>智能合同管理</h4>
                <p>全生命周期追踪</p>
              </div>
            </div>
            <div class="feature-item">
              <div class="feature-icon">
                <el-icon :size="24"><ChatDotRound /></el-icon>
              </div>
              <div class="feature-text">
                <h4>互动日志分析</h4>
                <p>知识智能提取</p>
              </div>
            </div>
            <div class="feature-item">
              <div class="feature-icon">
                <el-icon :size="24"><Odometer /></el-icon>
              </div>
              <div class="feature-text">
                <h4>数据驱动决策</h4>
                <p>可视化仪表盘</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { User, Lock, Message, Document, ChatDotRound, Odometer } from '@element-plus/icons-vue'
import { useUserStore } from '../../stores/user'
import { authApi } from '../../api'

const router = useRouter()
const userStore = useUserStore()
const loginFormRef = ref()
const registerFormRef = ref()
const loading = ref(false)
const isRegister = ref(false)

const loginForm = reactive({
  username: '',
  password: ''
})

const registerForm = reactive({
  username: '',
  email: '',
  password: '',
  confirmPassword: ''
})

const validateConfirmPassword = (_rule: any, value: string, callback: any) => {
  if (value !== registerForm.password) {
    callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

const loginRules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
}

const registerRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度在 3 到 20 个字符', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱地址', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于 6 个字符', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请再次输入密码', trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' }
  ]
}

const toggleMode = () => {
  isRegister.value = !isRegister.value
}

const handleLogin = async () => {
  await loginFormRef.value.validate()
  loading.value = true
  try {
    await userStore.login(loginForm.username, loginForm.password)
    ElMessage.success('登录成功，欢迎回来！')
    router.push('/dashboard')
  } catch (error: any) {
    ElMessage.error(error.response?.data?.error || '登录失败，请检查用户名和密码')
  } finally {
    loading.value = false
  }
}

const handleRegister = async () => {
  await registerFormRef.value.validate()
  loading.value = true
  try {
    await authApi.register({
      username: registerForm.username,
      email: registerForm.email,
      password: registerForm.password
    })
    ElMessage.success('注册成功！请登录')
    isRegister.value = false
    loginForm.username = registerForm.username
  } catch (error: any) {
    ElMessage.error(error.response?.data?.error || '注册失败，请稍后重试')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.auth-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
}

.auth-background {
  position: absolute;
  inset: 0;
  background: var(--color-bg-primary);
  overflow: hidden;
}

.grid-overlay {
  position: absolute;
  inset: 0;
  background-image: 
    linear-gradient(rgba(6, 182, 212, 0.03) 1px, transparent 1px),
    linear-gradient(90deg, rgba(6, 182, 212, 0.03) 1px, transparent 1px);
  background-size: 50px 50px;
  animation: gridMove 20s linear infinite;
}

@keyframes gridMove {
  0% { transform: translate(0, 0); }
  100% { transform: translate(50px, 50px); }
}

.gradient-orb {
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
  opacity: 0.5;
  animation: orbFloat 10s ease-in-out infinite;
}

.orb-1 {
  width: 400px;
  height: 400px;
  background: radial-gradient(circle, var(--color-accent) 0%, transparent 70%);
  top: -100px;
  right: -100px;
  animation-delay: 0s;
}

.orb-2 {
  width: 300px;
  height: 300px;
  background: radial-gradient(circle, var(--color-accent-secondary) 0%, transparent 70%);
  bottom: -50px;
  left: -50px;
  animation-delay: -3s;
}

.orb-3 {
  width: 200px;
  height: 200px;
  background: radial-gradient(circle, #06b6d4 0%, transparent 70%);
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  animation-delay: -5s;
}

@keyframes orbFloat {
  0%, 100% { transform: translate(0, 0) scale(1); }
  33% { transform: translate(30px, -30px) scale(1.1); }
  66% { transform: translate(-20px, 20px) scale(0.9); }
}

.auth-content {
  position: relative;
  z-index: 1;
  display: flex;
  gap: 60px;
  align-items: center;
  padding: 40px;
}

.auth-card {
  width: 420px;
  padding: 48px 40px;
  border-radius: var(--radius-lg);
  transition: all 0.5s cubic-bezier(0.4, 0, 0.2, 1);
}

.auth-card.flipped {
  transform: rotateY(0deg);
}

.auth-header {
  text-align: center;
  margin-bottom: 40px;
}

.logo-container {
  width: 80px;
  height: 80px;
  margin: 0 auto 24px;
}

.logo-icon {
  width: 100%;
  height: 100%;
}

.logo-icon svg {
  width: 100%;
  height: 100%;
}

.auth-title {
  font-family: var(--font-display);
  font-size: 28px;
  font-weight: 700;
  color: var(--color-text-primary);
  margin-bottom: 8px;
  letter-spacing: -0.5px;
}

.auth-subtitle {
  font-size: 14px;
  color: var(--color-text-muted);
  font-weight: 400;
}

.auth-form {
  animation: fadeIn 0.3s ease-out;
}

.auth-button {
  width: 100%;
  height: 48px;
  font-size: 16px;
  font-weight: 600;
  letter-spacing: 0.5px;
}

.auth-footer {
  text-align: center;
  margin-top: 24px;
  padding-top: 24px;
  border-top: 1px solid var(--color-border);
}

.auth-footer-text {
  color: var(--color-text-muted);
  font-size: 14px;
}

.auth-link {
  color: var(--color-accent);
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  margin-left: 8px;
  transition: all var(--transition-fast);
}

.auth-link:hover {
  color: var(--color-accent-secondary);
  text-decoration: underline;
}

.auth-info {
  max-width: 480px;
}

.info-content {
  animation: slideIn 0.6s ease-out;
}

.info-title {
  font-family: var(--font-display);
  font-size: 48px;
  font-weight: 700;
  line-height: 1.2;
  color: var(--color-text-primary);
  margin-bottom: 20px;
}

.info-description {
  font-size: 16px;
  line-height: 1.7;
  color: var(--color-text-secondary);
  margin-bottom: 40px;
}

.features-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 20px;
}

.feature-item {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px;
  background: rgba(255, 255, 255, 0.02);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  transition: all var(--transition-normal);
}

.feature-item:hover {
  background: rgba(255, 255, 255, 0.05);
  border-color: var(--color-border-hover);
  transform: translateY(-2px);
}

.feature-icon {
  width: 48px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, rgba(6, 182, 212, 0.2) 0%, rgba(139, 92, 246, 0.2) 100%);
  border-radius: var(--radius-sm);
  color: var(--color-accent);
}

.feature-text h4 {
  font-family: var(--font-display);
  font-size: 14px;
  font-weight: 600;
  color: var(--color-text-primary);
  margin-bottom: 4px;
}

.feature-text p {
  font-size: 12px;
  color: var(--color-text-muted);
}

.fade-enter-active,
.fade-leave-active {
  transition: all 0.3s ease;
}

.fade-enter-from {
  opacity: 0;
  transform: translateX(20px);
}

.fade-leave-to {
  opacity: 0;
  transform: translateX(-20px);
}

@media (max-width: 1024px) {
  .auth-content {
    flex-direction: column;
    gap: 40px;
  }
  
  .auth-info {
    display: none;
  }
  
  .auth-card {
    width: 100%;
    max-width: 420px;
  }
}

@media (max-width: 480px) {
  .auth-card {
    padding: 32px 24px;
  }
  
  .auth-title {
    font-size: 24px;
  }
}
</style>
