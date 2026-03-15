<template>
  <div class="auth-container">
    <div class="auth-background">
      <div class="grid-overlay"></div>
      <div class="gradient-orb orb-1"></div>
      <div class="gradient-orb orb-2"></div>
      <div class="gradient-orb orb-3"></div>
    </div>
    
    <div class="lang-switch-wrapper">
      <LangSwitch />
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
            <p class="auth-subtitle">{{ $t('auth.loginTitle') }}</p>
          </div>

          <transition name="fade" mode="out-in">
            <div v-if="!isRegister" key="login" class="auth-form">
              <el-form :model="loginForm" :rules="loginRules" ref="loginFormRef">
                <el-form-item prop="username">
                  <el-input
                    v-model="loginForm.username"
                    :placeholder="$t('auth.username')"
                    size="large"
                    :prefix-icon="User"
                  />
                </el-form-item>
                <el-form-item prop="password">
                  <el-input
                    v-model="loginForm.password"
                    type="password"
                    :placeholder="$t('auth.password')"
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
                    <span v-if="!loading">{{ $t('auth.login') }}</span>
                    <span v-else>{{ $t('auth.verifying') }}</span>
                  </el-button>
                </el-form-item>
              </el-form>
              
              <div class="auth-footer">
                <span class="auth-footer-text">{{ $t('auth.noAccount') }}</span>
                <a class="auth-link" @click="toggleMode">{{ $t('auth.registerNow') }}</a>
              </div>
            </div>

            <div v-else key="register" class="auth-form">
              <el-form :model="registerForm" :rules="registerRules" ref="registerFormRef">
                <el-form-item prop="username">
                  <el-input
                    v-model="registerForm.username"
                    :placeholder="$t('auth.username')"
                    size="large"
                    :prefix-icon="User"
                  />
                </el-form-item>
                <el-form-item prop="email">
                  <el-input
                    v-model="registerForm.email"
                    :placeholder="$t('auth.email')"
                    size="large"
                    :prefix-icon="Message"
                  />
                </el-form-item>
                <el-form-item prop="password">
                  <el-input
                    v-model="registerForm.password"
                    type="password"
                    :placeholder="$t('auth.password')"
                    size="large"
                    :prefix-icon="Lock"
                    show-password
                  />
                </el-form-item>
                <el-form-item prop="confirmPassword">
                  <el-input
                    v-model="registerForm.confirmPassword"
                    type="password"
                    :placeholder="$t('auth.confirmPassword')"
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
                    <span v-if="!loading">{{ $t('auth.register') }}</span>
                    <span v-else>{{ $t('auth.creating') }}</span>
                  </el-button>
                </el-form-item>
              </el-form>
              
              <div class="auth-footer">
                <span class="auth-footer-text">{{ $t('auth.hasAccount') }}</span>
                <a class="auth-link" @click="toggleMode">{{ $t('auth.backToLogin') }}</a>
              </div>
            </div>
          </transition>
        </div>
      </div>

      <div class="auth-info">
        <div class="info-content">
          <h2 class="info-title">
            <span class="gradient-text">{{ $t('dashboard.title') }}</span>
          </h2>
          <p class="info-description">
            {{ $t('dashboard.subtitle') }}
          </p>
          <div class="features-grid">
            <div class="feature-item">
              <div class="feature-icon">
                <el-icon :size="24"><User /></el-icon>
              </div>
              <div class="feature-text">
                <h4>{{ $t('nav.customers') }}</h4>
                <p>360° {{ $t('customer.basicInfo') }}</p>
              </div>
            </div>
            <div class="feature-item">
              <div class="feature-icon">
                <el-icon :size="24"><Document /></el-icon>
              </div>
              <div class="feature-text">
                <h4>{{ $t('nav.contracts') }}</h4>
                <p>{{ $t('contract.performanceRecords') }}</p>
              </div>
            </div>
            <div class="feature-item">
              <div class="feature-icon">
                <el-icon :size="24"><ChatDotRound /></el-icon>
              </div>
              <div class="feature-text">
                <h4>{{ $t('nav.interactions') }}</h4>
                <p>{{ $t('interaction.contentSummary') }}</p>
              </div>
            </div>
            <div class="feature-item">
              <div class="feature-icon">
                <el-icon :size="24"><Odometer /></el-icon>
              </div>
              <div class="feature-text">
                <h4>{{ $t('nav.dashboard') }}</h4>
                <p>{{ $t('dashboard.customerGrowth') }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useI18n } from 'vue-i18n'
import { User, Lock, Message, Document, ChatDotRound, Odometer } from '@element-plus/icons-vue'
import { useUserStore } from '../../stores/user'
import { authApi } from '../../api'
import LangSwitch from '../../components/LangSwitch.vue'

const { t } = useI18n()
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
    callback(new Error(t('auth.passwordMismatch')))
  } else {
    callback()
  }
}

const loginRules = computed(() => ({
  username: [{ required: true, message: t('auth.usernameRequired'), trigger: 'blur' }],
  password: [{ required: true, message: t('auth.passwordRequired'), trigger: 'blur' }]
}))

const registerRules = computed(() => ({
  username: [
    { required: true, message: t('auth.usernameRequired'), trigger: 'blur' },
    { min: 3, max: 20, message: t('auth.usernameLength'), trigger: 'blur' }
  ],
  email: [
    { required: true, message: t('auth.emailRequired'), trigger: 'blur' },
    { type: 'email', message: t('auth.emailFormat'), trigger: 'blur' }
  ],
  password: [
    { required: true, message: t('auth.passwordRequired'), trigger: 'blur' },
    { min: 6, message: t('auth.passwordLength'), trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: t('auth.confirmPasswordRequired'), trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' }
  ]
}))

const toggleMode = () => {
  isRegister.value = !isRegister.value
}

const handleLogin = async () => {
  await loginFormRef.value.validate()
  loading.value = true
  try {
    await userStore.login(loginForm.username, loginForm.password)
    ElMessage.success(t('auth.loginSuccess'))
    router.push('/dashboard')
  } catch (error: any) {
    ElMessage.error(error.response?.data?.error || t('auth.loginFailed'))
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
    ElMessage.success(t('auth.registerSuccess'))
    isRegister.value = false
    loginForm.username = registerForm.username
  } catch (error: any) {
    ElMessage.error(error.response?.data?.error || t('auth.registerFailed'))
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

.lang-switch-wrapper {
  position: absolute;
  top: 24px;
  right: 24px;
  z-index: 10;
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
