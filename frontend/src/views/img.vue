<template>
  <div ref="view" class="wrap" :style="{ padding: `0 ${SPACE_HORIZONTAL}px` }">
    <img class="setting-img" src="@/assets/imgs/setting.png" @click="openSet()" />
    <div v-for="(v, i) in LoopResultRes" :key="`${i}`" :style="{ width: `${COL_WIDTH}px`, marginRight: i < LoopResultRes.length - 1 ? `${SPACE_HORIZONTAL}px` : '0' }">
      <div
        class="content-wrap"
        v-for="(item, index) in v"
        :key="`${i}-${index}-${item.path}`"
        :style="{ height: `${item.viewH}px`, marginBottom: `${SPACE_VERTICAL}px` }"
        v-show="DelList.indexOf(item.path) === -1"
      >
        <img style="width: 100%" v-lazy="item.pic" />
        <img v-if="oldSet.delImg" class="del-img" src="@/assets/imgs/del.png" @click="delImgHandle(item, i, index)" />
      </div>
    </div>
    <n-modal
      :style="{
        maxWidth: '500px',
        backgroundColor: '#fff',
      }"
      v-model:show="showModal"
    >
      <n-card header-class="card-header" size="small" title="设置" :bordered="false">
        <n-form ref="formRef" label-placement="left" label-width="auto" require-mark-placement="right-hanging">
          <n-form-item label="列数">
            <n-input v-model:value="oldSet.rowCount" :placeholder="`请输入列数`" />
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
          <n-button @click="cancle()">取消</n-button>
          <n-button @click="save()" type="success">保存</n-button>
        </n-space>
      </n-card>
    </n-modal>
  </div>
</template>

<script>
import axios from "axios";
import { to, shuffleArray } from "@/utils";
import { NForm, NFormItem, NInput, NModal, NCard, NSpace, NButton, useMessage, NSwitch } from "naive-ui";
const API_PREFIX = "img";

const getMinIndex = (arr) => {
  let min = arr[0];
  //声明了个变量 保存下标值
  let index = 0;
  for (var i = 0; i < arr.length; i++) {
    if (min > arr[i]) {
      min = arr[i];
      index = i;
    }
  }
  return index;
};
// 图片默认最小宽高
const SPACE_HORIZONTAL = 4; // 元素块横向间距
const SPACE_VERTICAL = 4; // 元素块纵向间距
export default {
  name: "Search",
  setup() {
    window.$message = useMessage();
  },
  components: {
    NForm,
    NFormItem,
    NInput,
    NCard,
    NModal,
    NSpace,
    NButton,
    NSwitch,
  },
  data() {
    return {
      isShowPay: true,
      COL_WIDTH: 0,
      SPACE_HORIZONTAL: SPACE_HORIZONTAL,
      SPACE_VERTICAL: SPACE_VERTICAL,
      rowCount: 0, // 每行个数
      randomImg: true,
      delImg: false,
      timerID: null,
      List: [], // 列表展示数据 等于所有数据 [[],[]] // 每个数组一行
      LoopResult: [], // 所有数据 []
      LoopResultRes: [],
      DelList: [], // 已删除的图片
      ColHeightList: [], // 每列的已有高度
      showModal: false, // 展示设置弹窗
      oldSet: {
        rowCount: "0",
        SPACE_VERTICAL: `${SPACE_VERTICAL}`,
        SPACE_HORIZONTAL: `${SPACE_HORIZONTAL}`,
        randomImg: false,
        delImg: false,
      },
      key: 0,
      ViewWidth: 0,
    };
  },
  async mounted() {
    const isShowPay = localStorage.getItem("IsShowPay");
    if (isShowPay) {
      this.isShowPay = false;
    }
    const DelList = sessionStorage.getItem("DelList");
    if (DelList) {
      this.DelList = JSON.parse(DelList);
    }
    this.ViewWidth = this.$refs.view.clientWidth;

    const cfg = localStorage.getItem("cfg");
    if (cfg) {
      this.oldSet = JSON.parse(cfg);
      this.rowCount = Number(this.oldSet.rowCount);
      this.SPACE_VERTICAL = Number(this.oldSet.SPACE_VERTICAL);
      this.SPACE_HORIZONTAL = Number(this.oldSet.SPACE_HORIZONTAL);
      this.randomImg = this.oldSet.randomImg;
      this.delImg = this.oldSet.delImg;
    }

    await this.search();
    this.getViewWidth();
    window.addEventListener("resize", this.resizePage);
  },
  methods: {
    openSet() {
      this.oldSet = {
        rowCount: `${this.rowCount}`,
        SPACE_VERTICAL: `${this.SPACE_VERTICAL}`,
        SPACE_HORIZONTAL: `${this.SPACE_HORIZONTAL}`,
        randomImg: this.randomImg,
        delImg: this.delImg,
      };
      this.showModal = true;
    },
    cancle() {
      this.showModal = false;
    },
    async save() {
      this.showModal = false;
      this.rowCount = Number(this.oldSet.rowCount);
      this.SPACE_VERTICAL = Number(this.oldSet.SPACE_VERTICAL);
      this.SPACE_HORIZONTAL = Number(this.oldSet.SPACE_HORIZONTAL);
      // 触发乱序
      if (!this.randomImg && this.oldSet.randomImg) {
        this.key++;
        this.LoopResult = shuffleArray(this.LoopResult);
      }
      this.randomImg = this.oldSet.randomImg;
      this.delImg = this.oldSet.delImg;
      localStorage.setItem("cfg", JSON.stringify(this.oldSet));
      this.setElsWidth();
    },
    resizePage() {
      // 防抖
      clearTimeout(this.timerID);
      this.timerID = setTimeout(() => {
        const width = this.$refs.view.clientWidth;
        if (width !== this.ViewWidth) {
          this.ViewWidth = width;
          this.getViewWidth();
        }
      }, 100);
    },
    async search() {
      console.log(1);
      const searchUrl = `${API_PREFIX}/getImgsList.json?t=${Date.now()}`;
      console.log(2);
      const [err, res] = await to(axios.get(searchUrl));
      console.log(err)
      console.log(3);
      console.log(res);
      if (err) {
        window.$message.error(err);
        return false;
      }
      const data = res.data.imageList;
      console.log(4);
      console.log(data);
      if (!data || data.length === 0) {
        window.$message.error("暂无数据");
        return false;
      }
      let resData = data;
      if (this.randomImg) {
        resData = shuffleArray(data);
      }
      // const { name, type, pic, panLink, des } = res.data[0];
      this.LoopResult = resData;
    },
    setElsWidth() {
      // 0-768-1200-
      // 使用getBoundingClientRect获取实际可用宽度（不包括滚动条）
      const rect = this.$refs.view.getBoundingClientRect();
      const width = rect.width;
      // 减去左右padding
      const contentWidth = width - 2 * this.SPACE_HORIZONTAL;
      // 计算每列宽度：(内容宽度 - (列数-1)*列间距) / 列数
      this.COL_WIDTH = parseInt((contentWidth - (this.rowCount - 1) * this.SPACE_HORIZONTAL) / this.rowCount);
      this.splitData();
    },
    getViewWidth() {
      // 0-768-1200-
      if (this.rowCount) {
        this.setElsWidth();
        return false;
      }
      const width = this.$refs.view.clientWidth;
      if (width <= 800) {
        this.rowCount = 2;
        this.setElsWidth();
      } else {
        this.rowCount = parseInt(width / 800) * 2;
        this.setElsWidth();
      }
    },
    splitData() {
      this.ColHeightList = Array(this.rowCount).fill(this.SPACE_VERTICAL);
      this.LoopResultRes = Array(this.rowCount)
        .fill("")
        .map(() => []); // 存N列，确保每个子数组都是独立实例
      if (!this.LoopResult.length) {
        return false;
      }
      this.LoopResult.forEach((item) => {
        const mixIndex = getMinIndex(this.ColHeightList);
        const height = parseInt((this.COL_WIDTH * item.picH) / item.picW);

        this.ColHeightList[mixIndex] += height + this.SPACE_VERTICAL;
        this.LoopResultRes[mixIndex].push({
          viewH: height,
          pic: item.pic,
          path: item.path,
        });
      });
    },
    async delImgHandle({ path, viewH }, i, index) {
      const [err, res] = await to(
        axios({
          method: "post",
          url: `${API_PREFIX}/delImg.json`,
          data: { path },
        })
      );
      if (err) {
        window.$message.error(err);
        return false;
      }
      const { code, msg } = res.data;
      if (code !== "0") {
        window.$message.error(msg);
        return false;
      }
      window.$message.success(`图片已成功从磁盘删除`);
      setTimeout(() => {
        this.DelList.push(path);
        sessionStorage.setItem("DelList", JSON.stringify(this.DelList));
        this.ColHeightList[i] -= viewH - this.SPACE_VERTICAL;
      }, 300);
    },
  },
};
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
/* 整体页面容器样式 */
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
    // border-radius: 12px;
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
