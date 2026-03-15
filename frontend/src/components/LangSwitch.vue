<template>
  <div class="lang-switch" :title="$t('lang.switch')">
    <div class="lang-switch-track" @click="toggleLang">
      <div class="lang-switch-thumb" :class="{ 'thumb-en': currentLang === 'en' }">
        <span class="lang-icon">{{ currentLang === 'zh' ? '中' : 'EN' }}</span>
      </div>
      <div class="lang-labels">
        <span class="lang-label" :class="{ active: currentLang === 'zh' }">中</span>
        <span class="lang-label" :class="{ active: currentLang === 'en' }">EN</span>
      </div>
    </div>
    <div class="lang-glow"></div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { setLocale } from '@/locales'

const { locale } = useI18n()

const currentLang = computed(() => locale.value)

const toggleLang = () => {
  const newLang = currentLang.value === 'zh' ? 'en' : 'zh'
  setLocale(newLang)
}
</script>

<style scoped>
.lang-switch {
  position: relative;
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
}

.lang-switch-track {
  position: relative;
  width: 72px;
  height: 32px;
  background: linear-gradient(135deg, rgba(6, 182, 212, 0.1) 0%, rgba(139, 92, 246, 0.1) 100%);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 16px;
  display: flex;
  align-items: center;
  padding: 3px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  overflow: hidden;
}

.lang-switch-track::before {
  content: '';
  position: absolute;
  inset: 0;
  background: linear-gradient(135deg, rgba(6, 182, 212, 0.05) 0%, rgba(139, 92, 246, 0.05) 100%);
  opacity: 0;
  transition: opacity 0.3s ease;
}

.lang-switch-track:hover::before {
  opacity: 1;
}

.lang-switch-track:hover {
  border-color: rgba(6, 182, 212, 0.3);
  box-shadow: 0 0 20px rgba(6, 182, 212, 0.15);
}

.lang-switch-thumb {
  position: absolute;
  left: 3px;
  width: 28px;
  height: 26px;
  background: linear-gradient(135deg, #06b6d4 0%, #8b5cf6 100%);
  border-radius: 13px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  z-index: 2;
  box-shadow: 0 2px 8px rgba(6, 182, 212, 0.4);
}

.lang-switch-thumb.thumb-en {
  left: calc(100% - 31px);
}

.lang-icon {
  font-size: 12px;
  font-weight: 700;
  color: white;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.2);
  letter-spacing: 0.5px;
}

.lang-labels {
  position: relative;
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 10px;
  z-index: 1;
}

.lang-label {
  font-size: 11px;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.4);
  transition: all 0.3s ease;
  user-select: none;
}

.lang-label.active {
  opacity: 0;
}

.lang-glow {
  position: absolute;
  width: 40px;
  height: 40px;
  background: radial-gradient(circle, rgba(6, 182, 212, 0.3) 0%, transparent 70%);
  border-radius: 50%;
  pointer-events: none;
  opacity: 0;
  transition: opacity 0.3s ease;
  left: 50%;
  top: 50%;
  transform: translate(-50%, -50%);
}

.lang-switch:hover .lang-glow {
  opacity: 1;
}
</style>
