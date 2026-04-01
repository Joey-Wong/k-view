export interface ImageSize {
  width: number
  height: number
}

/**
 * Promise 结果包装：返回 [error, data] 元组，避免 try/catch 嵌套
 */
export const to = <T>(promise: Promise<T>): Promise<[Error | null, T | null]> =>
  promise
    .then((data): [null, T] => [null, data])
    .catch((err: Error): [Error, null] => [err, null])

/**
 * 获取图片的原始宽高
 */
export const getImgWH = (img: string): Promise<ImageSize> => {
  const image = new Image()
  image.crossOrigin = ''
  image.src = img
  return new Promise((resolve) => {
    image.onload = function () {
      const { width, height } = image
      resolve({ width, height })
    }
  })
}

/**
 * Fisher-Yates 洗牌算法，返回新数组
 */
export const shuffleArray = <T>(array: T[]): T[] => {
  const newArray = [...array]
  for (let i = newArray.length - 1; i > 0; i--) {
    const j = Math.floor(Math.random() * (i + 1))
    ;[newArray[i], newArray[j]] = [newArray[j], newArray[i]]
  }
  return newArray
}
