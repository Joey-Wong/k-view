const fs = require("fs-extra");
const path = require("path");
const convert = require("heic-convert");

// 封装 HEIC 转换中间件
const heicConvertMiddleware = (options = {}) => {
  // 默认配置
  const config = {
    dir: "./", // 默认文件存储目录
    quality: 1, // 默认转换质量
    ...options,
  };

  // 检测文件是否为 HEIC 格式（兼容所有常见 HEIC 变体，更可靠）
  async function isHeicFile(filePath) {
    try {
      // 读取文件开头 16 字节（覆盖所有 HEIC 变体的魔术字节范围）
      const buffer = await fs.readFile(filePath, { length: 16 });

      // 情况1：传统 HEIC 魔术字节（00 00 00 18 + ftypheic）
      const heicMagic1 = Buffer.from([
        0x00, 0x00, 0x00, 0x18, 0x66, 0x74, 0x79, 0x70, 0x68, 0x65, 0x69, 0x63,
      ]);
      // 情况2：小红书/苹果变体 HEIC 魔术字节（00 00 00 20 + ftypheic）
      const heicMagic2 = Buffer.from([
        0x00, 0x00, 0x00, 0x20, 0x66, 0x74, 0x79, 0x70, 0x68, 0x65, 0x69, 0x63,
      ]);

      // 情况3：通用匹配（只要包含 ftypheic 标识，忽略开头的容器大小）
      const heicCoreMagic = Buffer.from([
        0x66, 0x74, 0x79, 0x70, 0x68, 0x65, 0x69, 0x63,
      ]); // ftypheic
      const hasCoreMagic =
        buffer.length >= 12 && buffer.slice(4, 12).compare(heicCoreMagic) === 0;

      // 满足任意一种情况，都判定为 HEIC 文件
      const isMagic1Match =
        buffer.length >= 12 && buffer.compare(heicMagic1, 0, 12, 0, 12) === 0;
      const isMagic2Match =
        buffer.length >= 12 && buffer.compare(heicMagic2, 0, 12, 0, 12) === 0;

      return isMagic1Match || isMagic2Match || hasCoreMagic;
    } catch (err) {
      // 捕获文件不存在、权限不足、读取失败等异常，均返回 false
      return false;
    }
  }

  // 中间件核心逻辑
  return async (ctx, next) => {
    // 获取多级文件路径参数（如 b/a.png）
    const queryPath = ctx.path;
    const filePath = decodeURIComponent(queryPath.slice(1));
    // 如果没有文件路径参数，直接执行后续中间件
    if (!filePath) {
      await next();
      return false;
    }

    // 拼接完整的物理文件路径（处理路径分隔符兼容）
    const filePathFull = `${config.dir}/${filePath}`;
    // 检查文件是否存在
    if (!(await fs.pathExists(filePathFull))) {
      await next();
      return false;
    }
    console.log("[FullPath]");
    console.log(filePathFull);
    // 判断是否为 HEIC 格式
    if (await isHeicFile(filePathFull)) {
      console.log("[是 HEIC]");
      const heicBuffer = await fs.readFile(filePathFull);
      // 转换 HEIC 为 JPEG
      const jpgBuffer = await convert({
        buffer: heicBuffer,
        format: "JPEG",
        quality: config.quality,
      });
      // 设置响应头，返回 JPG 图片
      ctx.set("Content-Type", "image/jpeg");
      ctx.body = jpgBuffer;
    } else {
      console.log("[不是 HEIC]");
      // 非 HEIC 格式，直接执行后续中间件
      await next();
    }
  };
};

module.exports = heicConvertMiddleware;
