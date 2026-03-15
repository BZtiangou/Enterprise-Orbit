<template>
  <div class="main-layout">
    <aside class="sidebar">
      <div class="sidebar-header">
        <div class="logo">
          <div class="logo-icon">
            <svg viewBox="0 0 40 40" fill="none" xmlns="http://www.w3.org/2000/svg">
              <circle cx="20" cy="20" r="18" stroke="url(#gradient2)" stroke-width="2"/>
              <circle cx="20" cy="20" r="8" fill="url(#gradient2)"/>
              <circle cx="20" cy="8" r="3" fill="#06b6d4"/>
              <circle cx="32" cy="20" r="3" fill="#8b5cf6"/>
              <circle cx="20" cy="32" r="3" fill="#06b6d4"/>
              <circle cx="8" cy="20" r="3" fill="#8b5cf6"/>
              <defs>
                <linearGradient id="gradient2" x1="0" y1="0" x2="40" y2="40">
                  <stop offset="0%" stop-color="#06b6d4"/>
                  <stop offset="100%" stop-color="#8b5cf6"/>
                </linearGradient>
              </defs>
            </svg>
          </div>
          <div class="logo-text">
            <h1>Enterprise</h1>
            <span>Orbit</span>
          </div>
        </div>
      </div>
      
      <nav class="sidebar-nav">
        <div class="nav-section">
          <span class="nav-section-title">主菜单</span>
          <ul class="nav-list">
            <li>
              <router-link to="/dashboard" class="nav-item" :class="{ active: activeMenu === '/dashboard' }">
                <el-icon class="nav-icon"><Odometer /></el-icon>
                <span class="nav-text">仪表盘</span>
              </router-link>
            </li>
            <li>
              <router-link to="/customers" class="nav-item" :class="{ active: activeMenu.startsWith('/customers') }">
                <el-icon class="nav-icon"><User /></el-icon>
                <span class="nav-text">客户管理</span>
              </router-link>
            </li>
            <li>
              <router-link to="/contracts" class="nav-item" :class="{ active: activeMenu.startsWith('/contracts') }">
                <el-icon class="nav-icon"><Document /></el-icon>
                <span class="nav-text">合同管理</span>
              </router-link>
            </li>
            <li>
              <router-link to="/interactions" class="nav-item" :class="{ active: activeMenu === '/interactions' }">
                <el-icon class="nav-icon"><ChatDotRound /></el-icon>
                <span class="nav-text">互动日志</span>
              </router-link>
            </li>
          </ul>
        </div>
        
        <div class="nav-section">
          <span class="nav-section-title">系统设置</span>
          <ul class="nav-list">
            <li>
              <router-link to="/users" class="nav-item" :class="{ active: activeMenu === '/users' }">
                <el-icon class="nav-icon"><UserFilled /></el-icon>
                <span class="nav-text">用户管理</span>
              </router-link>
            </li>
          </ul>
        </div>
      </nav>
      
      <div class="sidebar-footer">
        <div class="user-info">
          <div class="user-avatar">
            <el-icon :size="20"><UserFilled /></el-icon>
          </div>
          <div class="user-details">
            <span class="user-name">管理员</span>
            <span class="user-role">系统管理员</span>
          </div>
          <el-dropdown trigger="click" @command="handleCommand">
            <el-icon class="user-menu-icon"><ArrowDown /></el-icon>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="profile">个人中心</el-dropdown-item>
                <el-dropdown-item command="settings">系统设置</el-dropdown-item>
                <el-dropdown-item divided command="logout">退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </div>
    </aside>
    
    <div class="main-container">
      <header class="header">
        <div class="header-left">
          <el-breadcrumb separator="/">
            <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
            <el-breadcrumb-item v-if="currentTitle !== '仪表盘'">{{ currentTitle }}</el-breadcrumb-item>
          </el-breadcrumb>
        </div>
        <div class="header-right">
          <div class="header-actions">
            <el-badge :value="notifications" :max="99" class="notification-badge">
              <el-button :icon="Bell" circle />
            </el-badge>
          </div>
        </div>
      </header>
      
      <main class="main-content">
        <router-view v-slot="{ Component }">
          <transition name="page" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '../stores/user'
import { Odometer, User, Document, ChatDotRound, UserFilled, ArrowDown, Bell } from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const notifications = ref(3)

const activeMenu = computed(() => route.path)
const currentTitle = computed(() => route.meta.title as string || '仪表盘')

const handleCommand = (command: string) => {
  if (command === 'logout') {
    userStore.logout()
    router.push('/login')
  }
}
</script>

<style scoped>
.main-layout {
  display: flex;
  min-height: 100vh;
  background: var(--color-bg-primary);
}

.sidebar {
  width: 260px;
  background: var(--color-bg-secondary);
  border-right: 1px solid var(--color-border);
  display: flex;
  flex-direction: column;
  position: fixed;
  top: 0;
  left: 0;
  height: 100vh;
  z-index: 100;
}

.sidebar-header {
  padding: 24px 20px;
  border-bottom: 1px solid var(--color-border);
}

.logo {
  display: flex;
  align-items: center;
  gap: 12px;
}

.logo-icon {
  width: 40px;
  height: 40px;
}

.logo-icon svg {
  width: 100%;
  height: 100%;
}

.logo-text h1 {
  font-family: var(--font-display);
  font-size: 18px;
  font-weight: 700;
  color: var(--color-text-primary);
  margin: 0;
  letter-spacing: -0.5px;
}

.logo-text span {
  font-family: var(--font-display);
  font-size: 12px;
  font-weight: 500;
  color: var(--color-accent);
  text-transform: uppercase;
  letter-spacing: 1px;
}

.sidebar-nav {
  flex: 1;
  padding: 16px 12px;
  overflow-y: auto;
}

.nav-section {
  margin-bottom: 24px;
}

.nav-section-title {
  display: block;
  font-size: 11px;
  font-weight: 600;
  color: var(--color-text-muted);
  text-transform: uppercase;
  letter-spacing: 1px;
  padding: 0 12px;
  margin-bottom: 8px;
}

.nav-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  border-radius: var(--radius-sm);
  color: var(--color-text-secondary);
  text-decoration: none;
  transition: all var(--transition-fast);
  margin-bottom: 4px;
}

.nav-item:hover {
  background: rgba(255, 255, 255, 0.05);
  color: var(--color-text-primary);
}

.nav-item.active {
  background: linear-gradient(135deg, rgba(6, 182, 212, 0.15) 0%, rgba(139, 92, 246, 0.15) 100%);
  color: var(--color-accent);
}

.nav-item.active .nav-icon {
  color: var(--color-accent);
}

.nav-icon {
  font-size: 20px;
  transition: color var(--transition-fast);
}

.nav-text {
  font-size: 14px;
  font-weight: 500;
}

.sidebar-footer {
  padding: 16px;
  border-top: 1px solid var(--color-border);
}

.user-info {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background: var(--color-bg-tertiary);
  border-radius: var(--radius-sm);
}

.user-avatar {
  width: 40px;
  height: 40px;
  background: linear-gradient(135deg, var(--color-accent) 0%, var(--color-accent-secondary) 100%);
  border-radius: var(--radius-sm);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.user-details {
  flex: 1;
  min-width: 0;
}

.user-name {
  display: block;
  font-size: 14px;
  font-weight: 600;
  color: var(--color-text-primary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.user-role {
  display: block;
  font-size: 12px;
  color: var(--color-text-muted);
}

.user-menu-icon {
  color: var(--color-text-muted);
  cursor: pointer;
  transition: color var(--transition-fast);
}

.user-menu-icon:hover {
  color: var(--color-accent);
}

.main-container {
  flex: 1;
  margin-left: 260px;
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}

.header {
  height: 64px;
  background: var(--color-bg-secondary);
  border-bottom: 1px solid var(--color-border);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  position: sticky;
  top: 0;
  z-index: 50;
}

.header-left {
  display: flex;
  align-items: center;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.header-actions {
  display: flex;
  align-items: center;
}

.notification-badge :deep(.el-button) {
  background: transparent;
  border: 1px solid var(--color-border);
  color: var(--color-text-secondary);
}

.notification-badge :deep(.el-button:hover) {
  border-color: var(--color-accent);
  color: var(--color-accent);
}

.main-content {
  flex: 1;
  padding: 24px;
  background: var(--color-bg-primary);
}

.page-enter-active,
.page-leave-active {
  transition: all 0.3s ease-out;
}

.page-enter-from {
  opacity: 0;
  transform: translateY(10px);
}

.page-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}
</style>
