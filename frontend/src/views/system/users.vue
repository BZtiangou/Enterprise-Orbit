<template>
  <div class="user-management">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>用户管理</span>
          <el-button type="primary" @click="showDialog = true">新增用户</el-button>
        </div>
      </template>
      
      <el-table :data="users" v-loading="loading" style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="username" label="用户名" />
        <el-table-column prop="email" label="邮箱" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'danger'">
              {{ row.status === 'active' ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="200">
          <template #default="{ row }">
            <el-button type="primary" link @click="editUser(row)">编辑</el-button>
            <el-button type="warning" link @click="resetPassword(row)">重置密码</el-button>
            <el-button type="danger" link @click="deleteUser(row.id)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="showDialog" :title="editingUser ? '编辑用户' : '新增用户'" width="500px">
      <el-form :model="userForm" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="userForm.username" :disabled="!!editingUser" />
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="userForm.email" />
        </el-form-item>
        <el-form-item v-if="!editingUser" label="密码" prop="password">
          <el-input v-model="userForm.password" type="password" show-password />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="userForm.status" style="width: 100%">
            <el-option label="启用" value="active" />
            <el-option label="禁用" value="inactive" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showDialog = false">取消</el-button>
        <el-button type="primary" @click="saveUser" :loading="saving">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'

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
  status: 'active'
})

const rules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  email: [{ required: true, message: '请输入邮箱', trigger: 'blur' }, { type: 'email', message: '邮箱格式不正确', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
}

const loadUsers = async () => {
  loading.value = true
  try {
    users.value = [
      { id: 1, username: 'admin', email: 'admin@example.com', status: 'active', created_at: '2024-01-01' },
      { id: 2, username: 'user1', email: 'user1@example.com', status: 'active', created_at: '2024-01-02' }
    ]
  } finally {
    loading.value = false
  }
}

const editUser = (user: any) => {
  editingUser.value = user
  Object.assign(userForm, user)
  showDialog.value = true
}

const saveUser = async () => {
  await formRef.value.validate()
  saving.value = true
  try {
    ElMessage.success(editingUser.value ? '更新成功' : '创建成功')
    showDialog.value = false
  } catch (error: any) {
    ElMessage.error(error.response?.data?.error || '操作失败')
  } finally {
    saving.value = false
  }
}

const resetPassword = async (user: any) => {
  await ElMessageBox.confirm(`确定要重置用户 ${user.username} 的密码吗？`, '提示', { type: 'warning' })
  ElMessage.success('密码已重置为默认密码')
}

const deleteUser = async (_id: number) => {
  await ElMessageBox.confirm('确定要删除该用户吗？', '提示', { type: 'warning' })
  ElMessage.success('删除成功')
  loadUsers()
}

onMounted(() => {
  loadUsers()
})
</script>

<style scoped>
.card-header { display: flex; justify-content: space-between; align-items: center; }
</style>
