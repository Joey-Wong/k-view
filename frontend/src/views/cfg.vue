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
          <button @click="selectVideoDir" class="e-button small-btn">浏览</button>
        </div>
      </div>
            <!-- 允许删除复选框 -->
      <div class="e-row">
        <label for="allowDel" class="e-label">深度扫描:</label>
        <input 
          type="checkbox" 
          v-model="config.isDeep" 
          id="isDeep" 
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
          type="checkbox" 
          v-model="config.isAllowDel" 
          id="allowDel" 
          class="e-checkbox"
        />
      </div>
      
      <!-- 按钮组 -->
      <div class="e-row button-group">
        <button @click="startServer" :disabled="isServerRunning" class="e-button">启动</button>
        <button @click="stopServer" :disabled="!isServerRunning" class="e-button">停止</button>
      </div>
      
      <!-- 状态信息 -->
      <div v-if="statusMsg" class="e-status">
        {{ statusMsg }}
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { SetConfig, GetConfig, StartServer, StopServer, SelectDirectory, OpenBrowser } from '../../wailsjs/go/main/App'

const config = ref({
  videoDir: '',
  port: 3000,
  allowedExtsStr: '.mp4,.avi,.mov,.mkv,.flv,.wmv',
  imageExtsStr: '.jpg,.jpeg,.png,.gif,.webp,.bmp',
  isAllowDel: false,
  isDeep: true
})

const isServerRunning = ref(false)
const statusMsg = ref('')

// 初始化配置
onMounted(async () => {
  try {
    const currentConfig = await GetConfig()
    config.value.videoDir = currentConfig.videoDir
    config.value.port = currentConfig.port
    config.value.allowedExtsStr = currentConfig.allowedExts.join(',')
    config.value.imageExtsStr = currentConfig.imageExts.join(',')
    config.value.isAllowDel = currentConfig.isAllowDel
    config.value.isDeep = currentConfig.isDeep
  } catch (error) {
    console.error('获取配置失败:', error)
  }
})

// 选择视频目录
const selectVideoDir = async () => {
  try {
    const directory = await SelectDirectory()
    if (directory) {
      config.value.videoDir = directory
    }
  } catch (error) {
    console.error('选择目录失败:', error)
    statusMsg.value = '选择目录失败: ' + error.message
  }
}

// 启动服务
const startServer = async () => {
  try {
    const allowedExts = config.value.allowedExtsStr.split(',').map(ext => ext.trim()).filter(ext => ext)
    const imageExts = config.value.imageExtsStr.split(',').map(ext => ext.trim()).filter(ext => ext)
    await SetConfig(
      config.value.videoDir,
      config.value.port,
      allowedExts,
      imageExts,
      config.value.isAllowDel,
      config.value.isDeep
    )
    const port = await StartServer()
    isServerRunning.value = true
    statusMsg.value = `启动成功: http://localhost:${port}/mv`
    // 自动打开浏览器
    await OpenBrowser()
    setTimeout(() => {
      statusMsg.value = ''
    }, 2000)
  } catch (error) {
    console.error('启动服务失败:', error)
    statusMsg.value = '启动服务失败: ' + error.message
  }
}

// 停止服务
const stopServer = async () => {
  try {
    await StopServer()
    isServerRunning.value = false
    statusMsg.value = '服务停止成功'
    setTimeout(() => {
      statusMsg.value = ''
    }, 2000)
  } catch (error) {
    console.error('停止服务失败:', error)
    statusMsg.value = '停止服务失败: ' + error.message
  }
}
</script>

<style scoped>
/* 核心容器 - 严格适配400x300 */
.e-lang-container {
  width: 100%;
  height: 100vh;
  background-color: #E8EDF2; /* 易语言经典背景色 */
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
  border: 1px solid #99B4D1; /* 易语言经典边框色 */
  border-radius: 3px; /* 易语言圆角 */
  background-color: #FFFFFF;
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
  border: 1px solid #2C69A1; /* 易语言按钮边框 */
  border-radius: 3px; /* 易语言圆角 */
  background: linear-gradient(to bottom, #F0F7FF 0%, #D6E8FF 100%); /* 易语言按钮渐变 */
  color: #1E4E79;
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
  background: #E0E0E0;
  color: #888888;
  border-color: #CCCCCC;
  cursor: not-allowed;
}

/* 按钮hover/点击效果 */
.e-button:not(:disabled):hover {
  background: linear-gradient(to bottom, #D6E8FF 0%, #B9D1EE 100%);
}

.e-button:not(:disabled):active {
  background: linear-gradient(to bottom, #B9D1EE 0%, #D6E8FF 100%);
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
  border: 1px solid #99B4D1;
  border-radius: 2px;
  background-color: #FFFFFF;
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
  color: #1E4E79; /* 易语言提示色 */
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
  border-color: #4A89DC; /* 易语言焦点色 */
}
</style>