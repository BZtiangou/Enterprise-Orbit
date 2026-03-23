<template>
  <div class="settings-page">
    <div class="page-header">
      <div class="page-title">
        <h1>{{ $t('settings.title') }}</h1>
      </div>
    </div>

    <el-row :gutter="24">
      <el-col :span="6">
        <div class="settings-nav glass">
          <div
            v-for="item in navItems"
            :key="item.key"
            class="nav-item"
            :class="{ active: activeSection === item.key }"
            @click="activeSection = item.key"
          >
            <el-icon><component :is="item.icon" /></el-icon>
            <span>{{ item.label }}</span>
          </div>
        </div>
      </el-col>

      <el-col :span="18">
        <div class="settings-content glass">
          <div v-if="activeSection === 'general'" class="section">
            <h3>{{ $t('settings.general') }}</h3>
            <el-form label-width="140px">
              <el-form-item :label="$t('settings.language')">
                <el-select v-model="settings.language" @change="changeLanguage">
                  <el-option label="中文" value="zh" />
                  <el-option label="English" value="en" />
                </el-select>
              </el-form-item>
              <el-form-item :label="$t('settings.theme')">
                <el-radio-group v-model="settings.theme">
                  <el-radio label="dark">{{ $t('settings.darkMode') }}</el-radio>
                  <el-radio label="light">{{ $t('settings.lightMode') }}</el-radio>
                  <el-radio label="auto">{{ $t('settings.autoMode') }}</el-radio>
                </el-radio-group>
              </el-form-item>
            </el-form>
          </div>

          <div v-if="activeSection === 'notification'" class="section">
            <h3>{{ $t('settings.notification') }}</h3>
            <el-form label-width="180px">
              <el-form-item :label="$t('settings.emailNotification')">
                <el-switch v-model="settings.emailNotification" />
              </el-form-item>
              <el-form-item :label="$t('settings.contractReminder')">
                <el-switch v-model="settings.contractReminder" />
                <div class="form-tip">{{ $t('settings.contractReminderTip') }}</div>
              </el-form-item>
              <el-form-item :label="$t('settings.interactionReminder')">
                <el-switch v-model="settings.interactionReminder" />
                <div class="form-tip">{{ $t('settings.interactionReminderTip') }}</div>
              </el-form-item>
            </el-form>
          </div>

          <div v-if="activeSection === 'security'" class="section">
            <h3>{{ $t('settings.security') }}</h3>
            <el-form label-width="180px">
              <el-form-item :label="$t('settings.sessionTimeout')">
                <el-select v-model="settings.sessionTimeout" style="width: 200px">
                  <el-option :label="$t('settings.minutes30')" :value="30" />
                  <el-option :label="$t('settings.hour1')" :value="60" />
                  <el-option :label="$t('settings.hour2')" :value="120" />
                  <el-option :label="$t('settings.hour4')" :value="240" />
                </el-select>
              </el-form-item>
              <el-form-item :label="$t('settings.twoFactorAuth')">
                <el-switch v-model="settings.twoFactorAuth" />
                <div class="form-tip">{{ $t('settings.twoFactorAuthTip') }}</div>
              </el-form-item>
              <el-form-item :label="$t('settings.changePassword')">
                <el-button type="primary" @click="goToProfile">{{ $t('settings.changePassword') }}</el-button>
              </el-form-item>
            </el-form>
          </div>

          <div class="section-footer">
            <el-button type="primary" @click="saveSettings" :loading="saving">{{ $t('settings.saveSettings') }}</el-button>
          </div>
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { Setting, Bell, Lock } from '@element-plus/icons-vue'
import { setLocale } from '@/locales'

const { t } = useI18n()
const router = useRouter()
const activeSection = ref('general')
const saving = ref(false)

const navItems = computed(() => [
  { key: 'general', label: t('settings.general'), icon: Setting },
  { key: 'notification', label: t('settings.notification'), icon: Bell },
  { key: 'security', label: t('settings.security'), icon: Lock }
])

const settings = reactive({
  language: localStorage.getItem('locale') || 'zh',
  theme: 'dark',
  emailNotification: true,
  contractReminder: true,
  interactionReminder: true,
  sessionTimeout: 60,
  twoFactorAuth: false
})

const changeLanguage = (lang: string) => {
  setLocale(lang as 'zh' | 'en')
}

const saveSettings = async () => {
  saving.value = true
  try {
    localStorage.setItem('settings', JSON.stringify(settings))
    ElMessage.success(t('settings.settingsSaved'))
  } finally {
    saving.value = false
  }
}

const goToProfile = () => {
  router.push('/profile')
}
</script>

<style scoped>
.settings-page {
  max-width: 1200px;
}

.page-header {
  margin-bottom: 24px;
}

.page-title h1 {
  font-family: var(--font-display);
  font-size: 24px;
  font-weight: 700;
  background: linear-gradient(135deg, var(--color-accent) 0%, var(--color-accent-secondary) 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin: 0;
}

.settings-nav {
  padding: 16px;
  border-radius: var(--radius-md);
  position: relative;
  overflow: hidden;
}

.settings-nav::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 2px;
  background: linear-gradient(90deg, var(--color-accent), var(--color-accent-secondary));
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 16px;
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: all 0.2s ease;
  color: var(--color-text-secondary);
  margin-bottom: 4px;
  border: 1px solid transparent;
}

.nav-item:hover {
  background: rgba(255, 255, 255, 0.05);
  color: var(--color-accent);
  border-color: rgba(6, 182, 212, 0.2);
}

.nav-item.active {
  background: linear-gradient(135deg, rgba(6, 182, 212, 0.2) 0%, rgba(139, 92, 246, 0.2) 100%);
  color: var(--color-accent);
  border-color: rgba(6, 182, 212, 0.3);
}

.settings-content {
  padding: 24px;
  border-radius: var(--radius-md);
  min-height: 400px;
  position: relative;
  overflow: hidden;
}

.settings-content::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 2px;
  background: linear-gradient(90deg, var(--color-accent), var(--color-accent-secondary));
}

.section h3 {
  font-family: var(--font-display);
  font-size: 18px;
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0 0 24px 0;
  padding-bottom: 16px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
}

.form-tip {
  font-size: 12px;
  color: var(--color-text-secondary);
  margin-top: 4px;
}

.section-footer {
  margin-top: 32px;
  padding-top: 24px;
  border-top: 1px solid rgba(255, 255, 255, 0.05);
}
</style>
