<template>
  <div ref="viewRef" class="wrap" :style="{ padding: `0 ${SPACE_HORIZONTAL}px` }">
    <img class="setting-img" src="@/assets/imgs/setting.png" @click="openSet()" />
    <div
      v-for="(col, i) in loopResultRes"
      :key="`col-${i}`"
      :style="{
        width: `${colWidth}px`,
        marginRight: i < loopResultRes.length - 1 ? `${SPACE_HORIZONTAL}px` : '0',
      }"
    >
      <div
        v-for="(item, index) in col"
        :key="`${i}-${index}-${item.path}`"
        class="content-wrap"
        :style="{ height: `${item.viewH}px`, marginBottom: `${SPACE_VERTICAL}px` }"
        v-show="!delList.includes(item.path)"
      >
        <img style="width: 100%" v-lazy="item.pic" />
        <img
          v-if="oldSet.delImg"
          class="del-img"
          src="@/assets/imgs/del.png"
          @click="delImgHandle(item, i)"
        />
      </div>
    </div>

    <n-modal
      v-model:show="showModal"
      :style="{ maxWidth: '500px', backgroundColor: '#fff' }"
    >
      <n-card header-class="card-header" size="small" title="设置" :bordered="false">
        <n-form ref="formRef" label-placement="left" label-width="auto" require-mark-placement="right-hanging">
          <n-form-item label="列数">
            <n-input v-model:value="oldSet.rowCount" placeholder="请输入列数" />
          </n-form-item>
          <n-form-item label="横间距">
            <n-input v-model:value="oldSet.SPACE_HORIZONTAL" placeholder="请输入横间距" />
          </n-form-item>
          <n-form-item label="纵间距">
            <n-input v-model:value="oldSet.SPACE_VERTICAL" placeholder="请输入纵间距" />
          </n-form-item>
          <n-form-item label="随机图片">
            <n-switch v-model:value="oldSet.randomImg" />
          </n-form-item>
          <n-form-item label="删除功能">
            <n-switch v-model:value="oldSet.delImg" />
          </n-form-item>
          <n-form-item label="图片总数">
            <n-input disabled :value="String(allCount)" placeholder="图片总数" />
          </n-form-item>
          <n-form-item v-show="isShowPay" label="打赏一下">
            <n-space>
              <div class="pay-wrap">
                <div class="pay-img">
                  <img src="@/assets/imgs/wechatpay.png" />
                  <img src="@/assets/imgs/zhifupay.png" />
                </div>
                <div class="pay-des">
                  <span>微信</span>
                  <span>支付宝</span>
                </div>
              </div>
            </n-space>
          </n-form-item>
        </n-form>
        <n-space justify="end">
          <n-button @click="cancel()">取消</n-button>
          <n-button type="success" @click="save()">保存</n-button>
        </n-space>
      </n-card>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onUnmounted } from 'vue'
import axios from 'axios'
import {
  NForm,
  NFormItem,
  NInput,
  NModal,
  NCard,
  NSpace,
  NButton,
  NSwitch,
  useMessage,
} from 'naive-ui'
import { to, shuffleArray } from '@/utils'
import { GetRandomAPI, GetImageList, DeleteImage, type ImageItem } from '@/services/index'

// ── 类型定义 ──────────────────────────────────────────
interface ViewItem {
  viewH: number
  pic: string
  path: string
}

interface Settings {
  rowCount: string
  SPACE_VERTICAL: string
  SPACE_HORIZONTAL: string
  randomImg: boolean
  delImg: boolean
}

// ── 常量 ──────────────────────────────────────────────
const DEFAULT_SPACE_HORIZONTAL = 4
const DEFAULT_SPACE_VERTICAL = 4

// ── 工具函数 ──────────────────────────────────────────
const getMinIndex = (arr: number[]): number => {
  let min = arr[0]
  let index = 0
  for (let i = 0; i < arr.length; i++) {
    if (min > arr[i]) {
      min = arr[i]
      index = i
    }
  }
  return index
}

// ── 消息 ──────────────────────────────────────────────
window.$message = useMessage()

// ── 模板引用 ──────────────────────────────────────────
const viewRef = ref<HTMLDivElement | null>(null)

// ── 状态 ──────────────────────────────────────────────
const randomAPI = ref<string>('')
const isShowPay = ref<boolean>(true)
const colWidth = ref<number>(0)
const SPACE_HORIZONTAL = ref<number>(DEFAULT_SPACE_HORIZONTAL)
const SPACE_VERTICAL = ref<number>(DEFAULT_SPACE_VERTICAL)
const rowCount = ref<number>(0)
const randomImg = ref<boolean>(true)
const delImg = ref<boolean>(false)
const timerID = ref<ReturnType<typeof setTimeout> | null>(null)
const loopResult = ref<ImageItem[]>([])
const loopResultRes = ref<ViewItem[][]>([])
const delList = ref<string[]>([])
const colHeightList = ref<number[]>([])
const showModal = ref<boolean>(false)
const viewWidth = ref<number>(0)
const allCount = ref<number>(0)

const oldSet = reactive<Settings>({
  rowCount: '0',
  SPACE_VERTICAL: `${DEFAULT_SPACE_VERTICAL}`,
  SPACE_HORIZONTAL: `${DEFAULT_SPACE_HORIZONTAL}`,
  randomImg: false,
  delImg: false,
})

// ── 方法 ──────────────────────────────────────────────
function openSet(): void {
  oldSet.rowCount = `${rowCount.value}`
  oldSet.SPACE_VERTICAL = `${SPACE_VERTICAL.value}`
  oldSet.SPACE_HORIZONTAL = `${SPACE_HORIZONTAL.value}`
  oldSet.randomImg = randomImg.value
  oldSet.delImg = delImg.value
  showModal.value = true
}

function cancel(): void {
  showModal.value = false
}

async function save(): Promise<void> {
  showModal.value = false
  rowCount.value = Number(oldSet.rowCount)
  SPACE_VERTICAL.value = Number(oldSet.SPACE_VERTICAL)
  SPACE_HORIZONTAL.value = Number(oldSet.SPACE_HORIZONTAL)
  // 触发乱序
  if (!randomImg.value && oldSet.randomImg) {
    loopResult.value = shuffleArray(loopResult.value)
  }
  randomImg.value = oldSet.randomImg
  delImg.value = oldSet.delImg
  localStorage.setItem('cfg', JSON.stringify({ ...oldSet }))
  setElsWidth()
}

function resizePage(): void {
  if (timerID.value) clearTimeout(timerID.value)
  timerID.value = setTimeout(() => {
    const width = viewRef.value?.clientWidth ?? 0
    if (width !== viewWidth.value) {
      viewWidth.value = width
      getViewWidth()
    }
  }, 100)
}

async function search(): Promise<void> {
  const searchUrl =
    GetImageList.replace('@GetRandomAPI', randomAPI.value) + `?t=${Date.now()}`
  const [err, res] = await to(axios.get<{ imageList: ImageItem[] }>(searchUrl))
  if (err || !res) {
    window.$message.error(String(err))
    return
  }
  const data = res.data.imageList
  if (!data || data.length === 0) {
    window.$message.error('暂无数据')
    return
  }
  loopResult.value = randomImg.value ? shuffleArray(data) : data
  allCount.value = data.length
}

function setElsWidth(): void {
  if (!viewRef.value) return
  const rect = viewRef.value.getBoundingClientRect()
  const width = rect.width
  const contentWidth = width - 2 * SPACE_HORIZONTAL.value
  colWidth.value = Math.floor(
    (contentWidth - (rowCount.value - 1) * SPACE_HORIZONTAL.value) / rowCount.value,
  )
  splitData()
}

function getViewWidth(): void {
  if (rowCount.value) {
    setElsWidth()
    return
  }
  const width = viewRef.value?.clientWidth ?? 0
  rowCount.value = width <= 800 ? 2 : Math.floor(width / 800) * 2
  setElsWidth()
}

function splitData(): void {
  colHeightList.value = Array<number>(rowCount.value).fill(SPACE_VERTICAL.value)
  loopResultRes.value = Array.from({ length: rowCount.value }, () => [] as ViewItem[])
  if (!loopResult.value.length) return

  loopResult.value.forEach((item) => {
    const minIdx = getMinIndex(colHeightList.value)
    const height = Math.floor((colWidth.value * item.picH) / item.picW)
    colHeightList.value[minIdx] += height + SPACE_VERTICAL.value
    loopResultRes.value[minIdx].push({
      viewH: height,
      pic: item.pic,
      path: item.path,
    })
  })
}

async function delImgHandle(item: ViewItem, colIndex: number): Promise<void> {
  const deleteUrl = DeleteImage.replace('@GetRandomAPI', randomAPI.value)
  const [err, res] = await to(
    axios.get<{ success: boolean; msg: string }>(deleteUrl, { params: { path: item.path } }),
  )
  if (err || !res) {
    window.$message.error(String(err))
    return
  }
  const { success, msg } = res.data
  if (!success) {
    window.$message.error(msg)
    return
  }
  window.$message.success('图片已成功从磁盘删除')
  setTimeout(() => {
    delList.value.push(item.path)
    sessionStorage.setItem('DelList', JSON.stringify(delList.value))
    colHeightList.value[colIndex] -= item.viewH - SPACE_VERTICAL.value
    allCount.value--
  }, 300)
}

// ── 生命周期 ──────────────────────────────────────────
onMounted(async () => {
  // 打赏状态
  if (localStorage.getItem('IsShowPay')) {
    isShowPay.value = false
  }
  // 已删除列表
  const savedDelList = sessionStorage.getItem('DelList')
  if (savedDelList) {
    delList.value = JSON.parse(savedDelList) as string[]
  }
  viewWidth.value = viewRef.value?.clientWidth ?? 0

  // 读取本地配置
  const cfg = localStorage.getItem('cfg')
  if (cfg) {
    const parsed = JSON.parse(cfg) as Settings
    Object.assign(oldSet, parsed)
    rowCount.value = Number(parsed.rowCount)
    SPACE_VERTICAL.value = Number(parsed.SPACE_VERTICAL)
    SPACE_HORIZONTAL.value = Number(parsed.SPACE_HORIZONTAL)
    randomImg.value = parsed.randomImg
    delImg.value = parsed.delImg
  }

  // 获取 randomAPI
  try {
    const response = await fetch(GetRandomAPI)
    const data: { randomAPI: string } = await response.json()
    randomAPI.value = data.randomAPI
  } catch (error) {
    console.error('获取randomAPI失败:', error)
    window.$message.error('无法连接到服务器，请确保服务已启动')
    return
  }

  await search()
  getViewWidth()
  window.addEventListener('resize', resizePage)
})

onUnmounted(() => {
  window.removeEventListener('resize', resizePage)
  if (timerID.value) clearTimeout(timerID.value)
})
</script>

<style>
:root {
  --search-input-bg-color: #fff;
  --search-input-font-color: #000;
  --search-btn-bg-color: #01847f;
  --search-btn-font-color: #fff;
  --bg-color: #263238;
  --tag-bg-color: #01847f;
  --tag-font-color: #fff;
  --block-box-shadow: 2px 2px 10px rgba(128, 128, 128, 0.2);
  --block-bg-color: #263238;
  --block-title-font-color: #fffffff2;
  --block-title-des-color: #fff9;
}
.light {
  --search-input-bg-color: #fff;
  --search-input-font-color: #666;
  --search-btn-bg-color: #01847f;
  --search-btn-font-color: #fff;
  --bg-color: #f3f3f3;
  --tag-bg-color: #fcfcfc;
  --tag-font-color: #999;
  --block-box-shadow: 2px 2px 10px rgba(128, 128, 128, 0.6);
  --block-bg-color: #fff;
  --block-title-font-color: #372929;
  --block-title-des-color: #716969;
}
</style>

<style lang="less" scoped>
.wrap {
  background-color: var(--bg-color);
  width: 100%;
  box-sizing: border-box;
  display: flex;
  justify-content: space-between;
  .setting-img {
    position: fixed;
    top: 10px;
    right: 10px;
    width: 30px;
    z-index: 999;
  }
  .content-wrap {
    background-color: var(--block-bg-color);
    box-shadow: var(--block-box-shadow);
    position: relative;
    .del-img {
      display: none;
      position: absolute;
      top: 10px;
      left: 10px;
      width: 20px;
    }
    &:hover {
      .del-img {
        display: block;
        cursor: pointer;
      }
    }
  }
}

.pay-wrap {
  .pay-img {
    display: flex;
    align-items: center;
    justify-content: space-around;
    font-size: 0;
    img {
      width: 45%;
    }
  }
  .pay-des {
    display: flex;
    align-items: center;
    font-size: 0;
    span {
      width: 50%;
      text-align: center;
      display: inline-block;
      font-size: 14px;
      color: #666;
    }
  }
}
</style>
