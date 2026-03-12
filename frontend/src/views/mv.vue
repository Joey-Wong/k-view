<template>
  <div class="container" @mousemove="handleMouseMove" @click="handleClick">
    <div v-if="hasVideo" class="video-container">
      <div class="top-bar" :class="{ hidden: !showControls }">
        <div class="video-info">{{ videoInfoText }}</div>
      </div>
      <video 
        ref="videoPlayer" 
        preload="metadata"
        @timeupdate="updateProgress"
        @loadedmetadata="updateProgress"
        @play="handlePlay"
        @pause="handlePause"
        @ended="handleEnded"
        @error="handleError"
      ></video>
      <div class="bottom-bar" :class="{ hidden: !showControls }">
        <div class="video-controls">
          <button class="control-btn play-btn" @click.stop="togglePlay">{{ playBtnText }}</button>
          <button class="control-btn" @click.stop="switchVideo('prev')">⏮️</button>
          <button class="control-btn" @click.stop="switchVideo('next')">⏭️</button>
          <button class="control-btn delete-btn" @click.stop="deleteVideo">🗑️</button>
          <div class="auto-play-switch">
            <label class="switch">
              <input type="checkbox" v-model="autoPlayEnabled">
              <span class="slider"></span>
            </label>
          </div>
          <div class="progress-container" ref="progressContainer" @click.stop="handleProgressClick" @mousedown.stop="handleProgressMouseDown">
            <div class="progress-bar" :style="{ width: progressPercent + '%' }"></div>
          </div>
          <div class="time-display">{{ timeDisplayText }}</div>
          <div class="volume-container">
            <span class="volume-icon">🔊</span>
            <input type="range" class="volume-slider" v-model="volume" min="0" max="1" step="0.1" @input.stop="handleVolumeChange">
          </div>
          <button class="fullscreen-btn" @click.stop="toggleFullscreen">⛶</button>
        </div>
      </div>
    </div>
    <div v-else class="empty-tip">
      <p>暂无视频文件 📂</p>
      <p style="margin-top: 1rem; font-size: 0.9rem;">请将视频文件放入以下目录：</p>
      <p style="margin-top: 0.5rem; font-family: monospace; color: #4299e1;">videos/</p>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, nextTick } from 'vue'

// 模板引用
const videoPlayer = ref(null)
const progressContainer = ref(null)

// 状态管理
const currentVideo = ref('')
const videoList = ref([])
const autoPlayEnabled = ref(false)
const showControls = ref(true)
const volume = ref(1)
const currentTime = ref(0)
const duration = ref(0)
const isDragging = ref(false)

let hideTimeout = null

// 计算属性
const hasVideo = computed(() => currentVideo.value && currentVideo.value.length > 0)

const videoInfoText = computed(() => {
  if (!hasVideo.value) return ''
  const currentIndex = videoList.value.indexOf(currentVideo.value)
  const total = videoList.value.length
  const displayIndex = currentIndex + 1
  return `${displayIndex}/${total} ${currentVideo.value}`
})

const playBtnText = computed(() => {
  return videoPlayer.value && !videoPlayer.value.paused ? '⏸️' : '▶️'
})

const progressPercent = computed(() => {
  if (!duration.value) return 0
  return (currentTime.value / duration.value) * 100
})

const timeDisplayText = computed(() => {
  return `${formatTime(currentTime.value)} / ${formatTime(duration.value)}`
})

// 格式化时间
function formatTime(seconds) {
  const mins = Math.floor(seconds / 60)
  const secs = Math.floor(seconds % 60)
  return `${mins.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`
}

// 更新进度条和时间显示
function updateProgress() {
  if (videoPlayer.value) {
    currentTime.value = videoPlayer.value.currentTime
    duration.value = videoPlayer.value.duration
  }
}

// 显示控制栏
function showControlsBar() {
  showControls.value = true
  if (hideTimeout) {
    clearTimeout(hideTimeout)
  }
  hideTimeout = setTimeout(() => {
    showControls.value = false
  }, 3000)
}

// 鼠标移动事件
function handleMouseMove() {
  showControlsBar()
}

// 点击事件
function handleClick() {
  showControlsBar()
}

// 播放/暂停
function togglePlay() {
  if (videoPlayer.value) {
    if (videoPlayer.value.paused) {
      videoPlayer.value.play()
    } else {
      videoPlayer.value.pause()
    }
  }
}

// 播放事件
function handlePlay() {
  showControlsBar()
}

// 暂停事件
function handlePause() {
  showControlsBar()
}

// 视频播放结束
function handleEnded() {
  if (autoPlayEnabled.value) {
    switchVideo('next', false)
  }
}

// 视频加载错误处理
function handleError(e) {
  alert(`视频加载失败：${e.target.error.message}`)
}

// 点击进度条跳转
function handleProgressClick(e) {
  if (videoPlayer.value && progressContainer.value) {
    const rect = progressContainer.value.getBoundingClientRect()
    const pos = (e.clientX - rect.left) / rect.width
    videoPlayer.value.currentTime = pos * videoPlayer.value.duration
  }
}

// 进度条拖动开始
function handleProgressMouseDown(e) {
  isDragging.value = true
  handleProgressClick(e)
}

// 鼠标移动事件（用于进度条拖动）
function handleDocumentMouseMove(e) {
  if (isDragging.value && videoPlayer.value && progressContainer.value) {
    const rect = progressContainer.value.getBoundingClientRect()
    const pos = Math.max(0, Math.min(1, (e.clientX - rect.left) / rect.width))
    videoPlayer.value.currentTime = pos * videoPlayer.value.duration
  }
}

// 鼠标松开事件
function handleDocumentMouseUp() {
  isDragging.value = false
}

// 音量控制
function handleVolumeChange() {
  if (videoPlayer.value) {
    videoPlayer.value.volume = volume.value
  }
}

// 全屏切换
function toggleFullscreen() {
  if (document.fullscreenElement) {
    document.exitFullscreen()
  } else {
    document.documentElement.requestFullscreen()
  }
}

// 切换视频
async function switchVideo(direction, showUI = true) {
  if (!currentVideo.value) return
  try {
    const response = await fetch(`/mv/switch.json?current=${encodeURIComponent(currentVideo.value)}&direction=${direction}`)
    const data = await response.json()
    if (data.success && data.currentVideo) {
      currentVideo.value = data.currentVideo
      videoList.value = data.videoList || videoList.value
      if (videoPlayer.value) {
        videoPlayer.value.src = '/videos/' + currentVideo.value
        videoPlayer.value.play()
      }
      if (showUI) {
        showControlsBar()
      }
    }
  } catch (error) {
    alert(`切换视频失败：${error.message}`)
  }
}

// 删除视频（直接删除，无需确认）
async function deleteVideo() {
  if (!currentVideo.value) return

  try {
    const response = await fetch(`/mv/del.json?fileName=${encodeURIComponent(currentVideo.value)}`)
    const data = await response.json()
    if (data.success) {
      if (data.hasVideo) {
        // 从视频列表中移除当前视频
        const currentIndex = videoList.value.indexOf(currentVideo.value)
        videoList.value.splice(currentIndex, 1)
        
        // 更新当前视频
        currentVideo.value = data.currentVideo
        if (videoPlayer.value) {
          videoPlayer.value.src = '/videos/' + currentVideo.value
          videoPlayer.value.play()
        }
        showControlsBar()
      } else {
        // 没有视频了，重新加载页面
        window.location.reload()
      }
    } else {
      alert(`删除失败：${data.msg}`)
    }
  } catch (error) {
    alert(`删除失败：${error.message}`)
  }
}

// 页面加载完成后初始化
onMounted(async () => {
  // 添加全局事件监听
  document.addEventListener('mousemove', handleDocumentMouseMove)
  document.addEventListener('mouseup', handleDocumentMouseUp)
  
  // 初始化音量
  if (videoPlayer.value) {
    videoPlayer.value.volume = volume.value
  }
  
  // 等待 DOM 更新完成
  await nextTick()
  
  // 初始化视频源
  try {
    const response = await fetch('/mv/list.json')
    const data = await response.json()
    console.log('API 响应:', data)
    if (data.success && data.videoList && data.videoList.length > 0) {
      videoList.value = data.videoList
      // 使用 API 返回的 currentVideo，如果没有则使用第一个视频
      currentVideo.value = data.currentVideo || data.videoList[0]
      console.log('当前视频:', currentVideo.value)
      console.log('videoPlayer.value:', videoPlayer.value)
      
      // 再次等待确保 videoPlayer 已经准备好
      await nextTick()
      
      if (videoPlayer.value) {
        const videoSrc = '/videos/' + currentVideo.value
        console.log('视频源:', videoSrc)
        videoPlayer.value.src = videoSrc
        // 加载视频
        videoPlayer.value.load()
        console.log('视频已加载，duration:', videoPlayer.value.duration)
        // 尝试自动播放，如果失败则等待用户手动点击
        videoPlayer.value.play().catch(err => {
          console.log('自动播放失败，等待用户手动播放:', err.message)
        })
      }
    }
  } catch (error) {
    console.error('初始化视频失败：', error)
  }
  
  // 初始显示控制栏
  showControlsBar()
})

// 组件卸载时清理事件监听
onUnmounted(() => {
  document.removeEventListener('mousemove', handleDocumentMouseMove)
  document.removeEventListener('mouseup', handleDocumentMouseUp)
  if (hideTimeout) {
    clearTimeout(hideTimeout)
  }
})
</script>

<style scoped>
.container {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  position: relative;
}
.video-container {
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
  background: #000;
}
video {
  width: 100%;
  height: 100%;
  object-fit: contain;
}
.top-bar {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  padding: 20px;
  background: linear-gradient(to bottom, rgba(0,0,0,0.7), transparent);
  transition: opacity 0.3s ease;
  z-index: 10;
}
.video-info {
  font-size: 18px;
  color: #fff;
  text-shadow: 1px 1px 2px rgba(0,0,0,0.5);
}
.bottom-bar {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  padding: 20px;
  background: linear-gradient(to top, rgba(0,0,0,0.7), transparent);
  transition: opacity 0.3s ease;
  z-index: 10;
}
.btn-group {
  display: flex;
  gap: 15px;
  justify-content: center;
  align-items: center;
}
button {
  padding: 12px 24px;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 16px;
  font-weight: 500;
  transition: all 0.2s ease;
  color: white;
}
.prev-btn {
  background: #4299e1;
}
.next-btn {
  background: #38b2ac;
}
.delete-btn {
  background: #e53e3e;
}
button:hover {
  opacity: 0.9;
  transform: translateY(-2px);
}
button:active {
  transform: translateY(0);
}
.auto-play-switch {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #fff;
  font-size: 14px;
}
.switch {
  position: relative;
  display: inline-block;
  width: 44px;
  height: 24px;
}
.switch input {
  opacity: 0;
  width: 0;
  height: 0;
}
.slider {
  position: absolute;
  cursor: pointer;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: #555;
  transition: 0.3s;
  border-radius: 24px;
}
.slider:before {
  position: absolute;
  content: "";
  height: 18px;
  width: 18px;
  left: 3px;
  bottom: 3px;
  background-color: white;
  transition: 0.3s;
  border-radius: 50%;
}
input:checked + .slider {
  background-color: #38b2ac;
}
input:checked + .slider:before {
  transform: translateX(20px);
}
.empty-tip {
  font-size: 1.2rem;
  color: #6c757d;
  text-align: center;
  padding: 6rem 2rem;
  background: #1a1a1a;
  border-radius: 8px;
}
.hidden {
  opacity: 0 !important;
  pointer-events: none;
}
.video-controls {
  display: flex;
  align-items: center;
  gap: 12px;
  width: 100%;
}
.progress-container {
  flex: 1;
  height: 6px;
  background: rgba(255,255,255,0.3);
  border-radius: 3px;
  cursor: pointer;
  position: relative;
}
.progress-bar {
  height: 100%;
  background: #38b2ac;
  border-radius: 3px;
  position: relative;
}
.progress-bar:hover {
  background: #2c7a7b;
}
.progress-container:hover .progress-bar {
  background: #38b2ac;
}
.time-display {
  font-size: 14px;
  color: #fff;
  min-width: 100px;
  text-align: center;
}
.control-btn {
  background: rgba(255,255,255,0.2);
  border: none;
  color: white;
  padding: 2px 6px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 18px;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  justify-content: center;
}
.control-btn:hover {
  background: rgba(255,255,255,0.3);
  transform: scale(1.05);
}
.play-btn {
  font-size: 20px;
}
.volume-container {
  display: flex;
  align-items: center;
  gap: 8px;
}
.volume-icon {
  font-size: 18px;
}
.volume-slider {
  width: 100px;
  height: 4px;
  -webkit-appearance: none;
  background: rgba(255,255,255,0.3);
  border-radius: 2px;
  cursor: pointer;
}
.volume-slider::-webkit-slider-thumb {
  -webkit-appearance: none;
  width: 12px;
  height: 12px;
  background: #fff;
  border-radius: 50%;
  cursor: pointer;
}
.fullscreen-btn {
  background: rgba(255,255,255,0.2);
  border: none;
  color: white;
  padding: 2px 6px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 18px;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  justify-content: center;
}
.fullscreen-btn:hover {
  background: rgba(255,255,255,0.3);
  transform: scale(1.05);
}
</style>