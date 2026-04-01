// ========================
// API 响应类型定义
// ========================

export interface RandomAPIResponse {
  randomAPI: string
}

export interface VideoListResponse {
  hasVideo: boolean
  videoList: string[]
  currentVideo: string
}

export interface SwitchVideoResponse {
  hasVideo: boolean
  currentVideo: string
  videoList: string[]
}

export interface DeleteVideoResponse {
  success: boolean
  msg: string
  hasVideo: boolean
  currentVideo: string
}

export interface ImageItem {
  pic: string
  path: string
  picW: number
  picH: number
}

export interface ImageListResponse {
  imageList: ImageItem[]
}

export interface DeleteImageResponse {
  success: boolean
  msg: string
}

// ========================
// API 路径常量
// ========================

/** [获取randomAPI] 参数：无 */
export const GetRandomAPI = `/GetRandomAPI.json`

/** [获取视频列表] 参数：无 */
export const GetVideoList = `/@GetRandomAPI/list.json`

/** [切换视频] 参数：current-当前视频文件名, direction-切换方向(prev/next) */
export const SwitchVideo = `/@GetRandomAPI/switch.json`

/** [删除视频] 参数：fileName-要删除的视频文件名 */
export const DeleteVideo = `/@GetRandomAPI/del.json`

/** [获取图片列表] 参数：无 */
export const GetImageList = `/@GetRandomAPI/getImgsList.json`

/** [删除图片] 参数：path-要删除的图片路径 */
export const DeleteImage = `/@GetRandomAPI/imagedel.json`

/** [视频文件路径] 使用方式：`${VideoPath}/${filename}` */
export const VideoPath = `/@GetRandomAPI/videos`

/** [图片文件路径] 使用方式：`${ImagePath}/${filename}` */
export const ImagePath = `/@GetRandomAPI/images`
