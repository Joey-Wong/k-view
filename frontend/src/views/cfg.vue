<template>
  <div class="e-lang-container">
    <!-- 易语言风格主面板 -->
    <div class="e-panel">
      <!-- 视频目录行 -->
      <div class="e-row">
        <label class="e-label">目录:</label>
        <div class="e-input-group">
          <input
            type="text"
            v-model="config.videoDir"
            placeholder="请选择目录"
            readonly
            class="e-edit"
          />
          <button class="e-button small-btn" @click="selectVideoDir">浏览</button>
        </div>
      </div>
      <!-- 深度扫描复选框 -->
      <div class="e-row">
        <label for="isDeep" class="e-label">深度扫描:</label>
        <input
          id="isDeep"
          type="checkbox"
          v-model="config.isDeep"
          class="e-checkbox"
        />
      </div>
      <!-- 服务端口行 -->
      <div class="e-row">
        <label class="e-label">服务端口:</label>
        <input
          type="number"
          v-model.number="config.port"
          min="1"
          max="65535"
          class="e-edit port-edit"
        />
      </div>
      <!-- 视频格式行 -->
      <div class="e-row">
        <label class="e-label">视频格式:</label>
        <input
          type="text"
          v-model="config.allowedExtsStr"
          placeholder=".mp4,.avi,.mov"
          class="e-edit"
        />
      </div>
      <!-- 图片格式行 -->
      <div class="e-row">
        <label class="e-label">图片格式:</label>
        <input
          type="text"
          v-model="config.imageExtsStr"
          placeholder=".jpg,.png,.gif"
          class="e-edit"
        />
      </div>
      <!-- 允许删除复选框 -->
      <div class="e-row">
        <label for="allowDel" class="e-label">允许删除:</label>
        <input
          id="allowDel"
          type="checkbox"
          v-model="config.isAllowDel"
          class="e-checkbox"
        />
      </div>
      <!-- 按钮组 -->
      <div class="e-row button-group">
        <button :disabled="isServerRunning" class="e-button" @click="startServer">启动</button>
        <button :disabled="!isServerRunning" class="e-button" @click="stopServer">停止</button>
      </div>
      <!-- 状态信息 -->
      <div v-if="statusMsg" class="e-status" v-html="statusMsg" />
    </div>
    <span class="version">v1.0.0</span>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { SetConfig, GetConfig, StartServer, StopServer, SelectDirectory } from '../../wailsjs/go/main/App'

interface Config {
  videoDir: string
  port: number
  allowedExtsStr: string
  imageExtsStr: string
  isAllowDel: boolean
  isDeep: boolean
}

const config = reactive<Config>({
  videoDir: '',
  port: 3000,
  allowedExtsStr: '.mp4,.avi,.mov,.mkv,.flv,.wmv',
  imageExtsStr: '.jpg,.jpeg,.png,.gif,.webp,.bmp',
  isAllowDel: false,
  isDeep: true,
})

const isServerRunning = ref<boolean>(false)
const statusMsg = ref<string>('')

// 初始化配置
onMounted(async () => {
  try {
    const currentConfig = await GetConfig()
    config.videoDir = currentConfig.videoDir
    config.port = currentConfig.port
    config.allowedExtsStr = currentConfig.allowedExts.join(',')
    config.imageExtsStr = currentConfig.imageExts.join(',')
    config.isAllowDel = currentConfig.isAllowDel
    config.isDeep = currentConfig.isDeep
  } catch (error) {
    console.error('获取配置失败:', error)
  }
})

// 选择视频目录
const selectVideoDir = async (): Promise<void> => {
  try {
    const directory = await SelectDirectory()
    if (directory) {
      config.videoDir = directory
    }
  } catch (error) {
    console.error('选择目录失败:', error)
    statusMsg.value = '选择目录失败: ' + (error as Error).message
  }
}

// 启动服务
const startServer = async (): Promise<void> => {
  try {
    const allowedExts = config.allowedExtsStr
      .split(',')
      .map((ext) => ext.trim())
      .filter((ext) => ext)
    const imageExts = config.imageExtsStr
      .split(',')
      .map((ext) => ext.trim())
      .filter((ext) => ext)
    await SetConfig(
      config.videoDir,
      config.port,
      allowedExts,
      imageExts,
      config.isAllowDel,
      config.isDeep,
    )
    const result = await StartServer()
    isServerRunning.value = true
    statusMsg.value = `[图片]${result.imageCount} [视频]${result.videoCount}`
  } catch (error) {
    console.error('启动服务失败:', error)
    statusMsg.value = '启动服务失败: ' + (error as Error).message
  }
}

// 停止服务
const stopServer = async (): Promise<void> => {
  try {
    await StopServer()
    isServerRunning.value = false
    statusMsg.value = '服务停止成功'
    setTimeout(() => {
      statusMsg.value = ''
    }, 2000)
  } catch (error) {
    console.error('停止服务失败:', error)
    statusMsg.value = '停止服务失败: ' + (error as Error).message
  }
}
</script>

<style scoped>
/* 核心容器 - 严格适配400x300 */
.e-lang-container {
  width: 100%;
  height: 100vh;
  background-color: #e8edf2; /* 易语言经典背景色 */
  font-family: "微软雅黑", "SimHei", sans-serif;
  font-size: 12px;
  overflow: hidden; /* 禁止滚动条 */
  box-sizing: border-box;
  padding: 10px 20px;
}

/* 易语言面板 */
.e-panel {
  width: 100%;
  height: 100%;
  box-sizing: border-box;
}

/* 行容器 - 紧凑布局 */
.e-row {
  display: flex;
  align-items: center;
  margin-bottom: 8px; /* 紧凑间距 */
  height: 26px; /* 统一行高 */
  box-sizing: border-box;
}

/* 易语言标签样式 */
.e-label {
  text-align: right;
  margin-right: 6px;
  color: #333333;
  font-weight: 500;
}

/* 输入框组 */
.e-input-group {
  display: flex;
  flex: 1;
  gap: 4px; /* 紧凑间距 */
}

/* 易语言编辑框样式 */
.e-edit {
  flex: 1;
  height: 24px; /* 紧凑高度 */
  padding: 0 4px;
  border: 1px solid #99b4d1; /* 易语言经典边框色 */
  border-radius: 3px; /* 易语言圆角 */
  background-color: #ffffff;
  font-size: 12px;
  box-sizing: border-box;
  outline: none;
}

/* 端口输入框 - 更紧凑 */
.port-edit {
  width: 70px;
  flex: none;
}

/* 易语言按钮样式 */
.e-button {
  height: 24px; /* 紧凑高度 */
  min-width: 55px; /* 紧凑宽度 */
  border: 1px solid #2c69a1; /* 易语言按钮边框 */
  border-radius: 3px; /* 易语言圆角 */
  background: linear-gradient(to bottom, #f0f7ff 0%, #d6e8ff 100%); /* 易语言按钮渐变 */
  color: #1e4e79;
  font-size: 12px;
  cursor: pointer;
  padding: 0 6px;
  box-sizing: border-box;
}

.e-button + .e-button {
  margin-left: 8px;
}

/* 小尺寸按钮（浏览） */
.small-btn {
  min-width: 45px;
}

/* 按钮禁用状态 */
.e-button:disabled {
  background: #e0e0e0;
  color: #888888;
  border-color: #cccccc;
  cursor: not-allowed;
}

/* 按钮hover/点击效果 */
.e-button:not(:disabled):hover {
  background: linear-gradient(to bottom, #d6e8ff 0%, #b9d1ee 100%);
}

.e-button:not(:disabled):active {
  background: linear-gradient(to bottom, #b9d1ee 0%, #d6e8ff 100%);
}

/* 复选框行 */
.checkbox-row {
  padding-left: 76px; /* 对齐标签位置 */
  height: 24px;
}

/* 易语言复选框样式 */
.e-checkbox {
  width: 14px;
  height: 14px;
  margin: 0 4px 0 0;
  border: 1px solid #99b4d1;
  border-radius: 2px;
  background-color: #ffffff;
  cursor: pointer;
}

/* 复选框标签 */
.e-checkbox-label {
  margin-left: 3px;
  cursor: pointer;
  color: #333333;
}

/* 按钮组 */
.button-group {
  align-items: center;
  justify-content: center;
  display: flex;
}

/* 状态信息 */
.e-status {
  margin-top: 4px;
  color: #1e4e79; /* 易语言提示色 */
  height: 18px;
  line-height: 18px;
  font-size: 11px; /* 更小字体 */
}

/* 全局重置 - 防止默认样式干扰 */
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

input:focus {
  border-color: #4a89dc; /* 易语言焦点色 */
}

.version {
  font-size: 12px;
  color: #666;
  font-weight: 500;
  line-height: 18px;
  position: fixed;
  bottom: 6px;
  right: 6px;
}
</style>
