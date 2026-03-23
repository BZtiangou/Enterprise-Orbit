<template>
  <div class="profile-page">
    <div class="page-header">
      <div class="page-title">
        <h1>{{ $t('profile.title') }}</h1>
      </div>
    </div>

    <el-row :gutter="24">
      <el-col :span="8">
        <div class="profile-card glass">
          <div class="avatar-section">
            <el-avatar :size="120" :src="userInfo.avatar || defaultAvatar" class="avatar">
              <el-icon :size="60"><User /></el-icon>
            </el-avatar>
            <el-upload
              :show-file-list="false"
              :before-upload="beforeAvatarUpload"
              :http-request="uploadAvatar"
              class="avatar-upload"
            >
              <el-button type="primary" size="small">{{ $t('profile.uploadAvatar') }}</el-button>
            </el-upload>
          </div>
          <div class="user-info">
            <h2>{{ userInfo.nickname || userInfo.username }}</h2>
            <p class="email">{{ userInfo.email }}</p>
            <div class="user-meta">
              <el-tag :type="userInfo.role === 'admin' ? 'danger' : 'primary'" size="small">
                {{ userInfo.role === 'admin' ? $t('user.admin') : $t('user.user') }}
              </el-tag>
              <el-tag :type="userInfo.status === 'active' ? 'success' : 'danger'" size="small">
                {{ userInfo.status === 'active' ? $t('common.active') : $t('common.inactive') }}
              </el-tag>
            </div>
          </div>
        </div>
      </el-col>

      <el-col :span="16">
        <div class="info-card glass">
          <el-tabs v-model="activeTab">
            <el-tab-pane :label="$t('profile.basicInfo')" name="basic">
              <el-form :model="profileForm" :rules="profileRules" ref="profileFormRef" label-width="120px">
                <el-row :gutter="20">
                  <el-col :span="12">
                    <el-form-item :label="$t('auth.username')">
                      <el-input v-model="userInfo.username" disabled />
                    </el-form-item>
                  </el-col>
                  <el-col :span="12">
                    <el-form-item :label="$t('auth.email')">
                      <el-input v-model="userInfo.email" disabled />
                    </el-form-item>
                  </el-col>
                </el-row>
                <el-row :gutter="20">
                  <el-col :span="12">
                    <el-form-item :label="$t('profile.nickname')" prop="nickname">
                      <el-input v-model="profileForm.nickname" />
                    </el-form-item>
                  </el-col>
                  <el-col :span="12">
                    <el-form-item :label="$t('profile.phone')" prop="phone">
                      <el-input v-model="profileForm.phone" />
                    </el-form-item>
                  </el-col>
                </el-row>
                <el-row :gutter="20">
                  <el-col :span="12">
                    <el-form-item :label="$t('profile.department')">
                      <el-input v-model="profileForm.department" />
                    </el-form-item>
                  </el-col>
                  <el-col :span="12">
                    <el-form-item :label="$t('profile.position')">
                      <el-input v-model="profileForm.position" />
                    </el-form-item>
                  </el-col>
                </el-row>
                <el-form-item>
                  <el-button type="primary" @click="updateProfile" :loading="saving">{{ $t('common.save') }}</el-button>
                </el-form-item>
              </el-form>
            </el-tab-pane>

            <el-tab-pane :label="$t('profile.changePassword')" name="password">
              <el-form :model="passwordForm" :rules="passwordRules" ref="passwordFormRef" label-width="120px">
                <el-form-item :label="$t('profile.currentPassword')" prop="current_password">
                  <el-input v-model="passwordForm.current_password" type="password" show-password />
                </el-form-item>
                <el-form-item :label="$t('profile.newPassword')" prop="new_password">
                  <el-input v-model="passwordForm.new_password" type="password" show-password />
                </el-form-item>
                <el-form-item :label="$t('profile.confirmNewPassword')" prop="confirm_password">
                  <el-input v-model="passwordForm.confirm_password" type="password" show-password />
                </el-form-item>
                <el-form-item>
                  <el-button type="primary" @click="changePassword" :loading="changingPassword">{{ $t('common.save') }}</el-button>
                </el-form-item>
              </el-form>
            </el-tab-pane>

            <el-tab-pane :label="$t('profile.accountInfo')" name="account">
              <div class="account-info">
                <div class="info-item">
                  <span class="label">{{ $t('profile.accountCreated') }}</span>
                  <span class="value">{{ formatDate(userInfo.created_at) }}</span>
                </div>
                <div class="info-item">
                  <span class="label">{{ $t('profile.lastLogin') }}</span>
                  <span class="value">{{ userInfo.last_login ? formatDate(userInfo.last_login) : '-' }}</span>
                </div>
              </div>
            </el-tab-pane>
          </el-tabs>
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { User } from '@element-plus/icons-vue'
import axios from 'axios'

const { t } = useI18n()
const activeTab = ref('basic')
const saving = ref(false)
const changingPassword = ref(false)
const profileFormRef = ref()
const passwordFormRef = ref()

const defaultAvatar = 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png'

const userInfo = ref<any>({
  username: '',
  email: '',
  nickname: '',
  phone: '',
  avatar: '',
  department: '',
  position: '',
  role: 'user',
  status: 'active',
  created_at: '',
  last_login: null
})

const profileForm = reactive({
  nickname: '',
  phone: '',
  department: '',
  position: ''
})

const passwordForm = reactive({
  current_password: '',
  new_password: '',
  confirm_password: ''
})

const profileRules = computed(() => ({}))

const validateConfirmPassword = (_rule: any, value: string, callback: any) => {
  if (value !== passwordForm.new_password) {
    callback(new Error(t('profile.passwordMismatch')))
  } else {
    callback()
  }
}

const passwordRules = computed(() => ({
  current_password: [{ required: true, message: t('common.required'), trigger: 'blur' }],
  new_password: [
    { required: true, message: t('common.required'), trigger: 'blur' },
    { min: 6, message: t('auth.passwordLength'), trigger: 'blur' }
  ],
  confirm_password: [
    { required: true, message: t('common.required'), trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' }
  ]
}))

const formatDate = (date: string | Date) => {
  if (!date) return '-'
  return new Date(date).toLocaleString()
}

const loadProfile = async () => {
  try {
    const token = localStorage.getItem('token')
    const res = await axios.get('/api/profile', {
      headers: { Authorization: `Bearer ${token}` }
    })
    userInfo.value = res.data
    Object.assign(profileForm, {
      nickname: res.data.nickname || '',
      phone: res.data.phone || '',
      department: res.data.department || '',
      position: res.data.position || ''
    })
  } catch (error) {
    console.error('Failed to load profile:', error)
  }
}

const updateProfile = async () => {
  await profileFormRef.value.validate()
  saving.value = true
  try {
    const token = localStorage.getItem('token')
    const res = await axios.put('/api/profile', profileForm, {
      headers: { Authorization: `Bearer ${token}` }
    })
    userInfo.value = res.data
    ElMessage.success(t('profile.updateSuccess'))
  } catch (error: any) {
    ElMessage.error(error.response?.data?.error || t('common.error'))
  } finally {
    saving.value = false
  }
}

const changePassword = async () => {
  await passwordFormRef.value.validate()
  changingPassword.value = true
  try {
    const token = localStorage.getItem('token')
    await axios.post('/api/profile/change-password', {
      current_password: passwordForm.current_password,
      new_password: passwordForm.new_password
    }, {
      headers: { Authorization: `Bearer ${token}` }
    })
    ElMessage.success(t('profile.passwordChanged'))
    passwordForm.current_password = ''
    passwordForm.new_password = ''
    passwordForm.confirm_password = ''
  } catch (error: any) {
    ElMessage.error(error.response?.data?.error || t('common.error'))
  } finally {
    changingPassword.value = false
  }
}

const beforeAvatarUpload = (file: File) => {
  const isImage = file.type.startsWith('image/')
  const isLt2M = file.size / 1024 / 1024 < 2
  if (!isImage) {
    ElMessage.error(t('profile.imageOnly'))
    return false
  }
  if (!isLt2M) {
    ElMessage.error(t('profile.imageSizeLimit'))
    return false
  }
  return true
}

const uploadAvatar = async (options: any) => {
  const formData = new FormData()
  formData.append('avatar', options.file)
  try {
    const token = localStorage.getItem('token')
    const res = await axios.post('/api/profile/avatar', formData, {
      headers: {
        Authorization: `Bearer ${token}`,
        'Content-Type': 'multipart/form-data'
      }
    })
    userInfo.value.avatar = res.data.url
    ElMessage.success(t('profile.updateSuccess'))
  } catch (error: any) {
    ElMessage.error(error.response?.data?.error || t('common.error'))
  }
}

onMounted(() => loadProfile())
</script>

<style scoped>
.profile-page {
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

.profile-card {
  padding: 30px;
  border-radius: var(--radius-md);
  text-align: center;
  position: relative;
  overflow: hidden;
}

.profile-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 3px;
  background: linear-gradient(90deg, var(--color-accent), var(--color-accent-secondary));
}

.avatar-section {
  margin-bottom: 20px;
}

.avatar {
  border: 4px solid rgba(6, 182, 212, 0.3);
  background: linear-gradient(135deg, #06b6d4 0%, #8b5cf6 100%);
  box-shadow: 0 4px 20px rgba(6, 182, 212, 0.3);
}

.avatar-upload {
  margin-top: 16px;
}

.user-info h2 {
  font-family: var(--font-display);
  font-size: 20px;
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0 0 8px 0;
}

.user-info .email {
  color: var(--color-text-secondary);
  margin: 0 0 16px 0;
}

.user-meta {
  display: flex;
  justify-content: center;
  gap: 8px;
}

.info-card {
  padding: 24px;
  border-radius: var(--radius-md);
  position: relative;
  overflow: hidden;
}

.info-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 2px;
  background: linear-gradient(90deg, var(--color-accent), var(--color-accent-secondary));
}

.account-info {
  padding: 20px 0;
}

.info-item {
  display: flex;
  justify-content: space-between;
  padding: 16px 0;
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
  transition: all var(--transition-fast);
}

.info-item:hover {
  background: rgba(255, 255, 255, 0.02);
  padding-left: 8px;
  padding-right: 8px;
  margin-left: -8px;
  margin-right: -8px;
  border-radius: var(--radius-sm);
}

.info-item:last-child {
  border-bottom: none;
}

.info-item .label {
  color: var(--color-text-secondary);
}

.info-item .value {
  color: var(--color-accent);
  font-weight: 500;
}
</style>
