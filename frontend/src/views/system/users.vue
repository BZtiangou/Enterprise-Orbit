<template>
  <div class="user-management">
    <div class="page-header">
      <div class="page-title">
        <h1>{{ $t('user.title') }}</h1>
      </div>
      <el-button type="primary" @click="openCreateDialog">
        <el-icon><Plus /></el-icon>
        {{ $t('common.create') }}
      </el-button>
    </div>

    <div class="table-card glass">
      <el-table :data="users" v-loading="loading" :element-loading-text="$t('common.loading')" element-loading-background="rgba(10, 15, 28, 0.7)" style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="username" :label="$t('user.username')" />
        <el-table-column prop="email" :label="$t('user.email')" />
        <el-table-column prop="role" :label="$t('user.role')" width="120">
          <template #default="{ row }">
            <el-tag :type="row.role === 'admin' ? 'danger' : 'primary'">
              {{ row.role === 'admin' ? $t('user.admin') : $t('user.user') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" :label="$t('common.status')" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'danger'">
              {{ row.status === 'active' ? $t('common.active') : $t('common.inactive') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" :label="$t('user.createdAt')" width="180" />
        <el-table-column :label="$t('common.action')" width="250" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="editUser(row)">{{ $t('common.edit') }}</el-button>
            <el-button type="warning" link @click="resetPassword(row)">{{ $t('user.resetPassword') }}</el-button>
            <el-button type="danger" link @click="deleteUser(row.id)">{{ $t('common.delete') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <el-dialog v-model="showDialog" :title="editingUser ? $t('common.edit') : $t('common.create')" width="500px">
      <el-form :model="userForm" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item :label="$t('user.username')" prop="username">
          <el-input v-model="userForm.username" :disabled="!!editingUser" />
        </el-form-item>
        <el-form-item :label="$t('user.email')" prop="email">
          <el-input v-model="userForm.email" />
        </el-form-item>
        <el-form-item v-if="!editingUser" :label="$t('user.password')" prop="password">
          <el-input v-model="userForm.password" type="password" show-password />
        </el-form-item>
        <el-form-item :label="$t('user.role')">
          <el-select v-model="userForm.role" style="width: 100%">
            <el-option :label="$t('user.admin')" value="admin" />
            <el-option :label="$t('user.user')" value="user" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('common.status')">
          <el-select v-model="userForm.status" style="width: 100%">
            <el-option :label="$t('common.active')" value="active" />
            <el-option :label="$t('common.inactive')" value="inactive" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showDialog = false">{{ $t('common.cancel') }}</el-button>
        <el-button type="primary" @click="saveUser" :loading="saving">{{ $t('common.save') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { mockUsers } from '../../mock/data'

const { t } = useI18n()
const loading = ref(false)
const saving = ref(false)
const showDialog = ref(false)
const users = ref<any[]>([])
const editingUser = ref<any>(null)
const formRef = ref()

const userForm = reactive({
  username: '',
  email: '',
  password: '',
  role: 'user',
  status: 'active'
})

const rules = computed(() => ({
  username: [{ required: true, message: t('common.required'), trigger: 'blur' }],
  email: [
    { required: true, message: t('common.required'), trigger: 'blur' },
    { type: 'email', message: t('user.emailFormat'), trigger: 'blur' }
  ],
  password: [{ required: true, message: t('common.required'), trigger: 'blur' }]
}))

const loadUsers = async () => {
  loading.value = true
  try {
    users.value = mockUsers
  } finally {
    loading.value = false
  }
}

const openCreateDialog = () => {
  editingUser.value = null
  Object.assign(userForm, { username: '', email: '', password: '', role: 'user', status: 'active' })
  showDialog.value = true
}

const editUser = (user: any) => {
  editingUser.value = user
  Object.assign(userForm, { ...user, password: '' })
  showDialog.value = true
}

const saveUser = async () => {
  await formRef.value.validate()
  saving.value = true
  try {
    ElMessage.success(t('common.success'))
    showDialog.value = false
    loadUsers()
  } catch (error: any) {
    ElMessage.error(error.response?.data?.error || t('common.error'))
  } finally {
    saving.value = false
  }
}

const resetPassword = async (user: any) => {
  await ElMessageBox.confirm(t('user.resetPasswordConfirm', { username: user.username }), t('user.resetPassword'), { type: 'warning' })
  ElMessage.success(t('user.resetPasswordSuccess'))
}

const deleteUser = async (_id: number) => {
  await ElMessageBox.confirm(t('common.confirm') + '?', t('common.delete'), { type: 'warning' })
  ElMessage.success(t('common.success'))
  loadUsers()
}

onMounted(() => loadUsers())
</script>

<style scoped>
.user-management { max-width: 1400px; }
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 24px; }
.page-title h1 { font-family: var(--font-display); font-size: 24px; font-weight: 700; background: linear-gradient(135deg, var(--color-accent) 0%, var(--color-accent-secondary) 100%); -webkit-background-clip: text; -webkit-text-fill-color: transparent; background-clip: text; margin: 0; }
.table-card { padding: 20px; border-radius: var(--radius-md); position: relative; overflow: hidden; }
.table-card::before { content: ''; position: absolute; top: 0; left: 0; right: 0; height: 2px; background: linear-gradient(90deg, var(--color-accent), var(--color-accent-secondary)); }
</style>
