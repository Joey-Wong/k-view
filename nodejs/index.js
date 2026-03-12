const Koa = require("koa");
const KoaStatic = require("koa-static");
const BodyParser = require("koa-bodyparser");
const fs = require("fs");
const { getImgs } = require("./utils/getImgs");
const { shuffleArray } = require("./utils");
const EXE_DIR = process.cwd().replace(/\\/g, "/");
const c = require("child_process");
const heicConvertMiddleware = require("./utils/heip2png");

const DEFAULT_CFG = {
  port: 666,
  dir: EXE_DIR,
  exts: [
    ".jpg",
    ".jpeg",
    ".png",
    ".gif",
    ".webp",
    ".JPG",
    ".JPEG",
    ".PNG",
    ".GIF",
    ".WEBP",
  ],
  disorder: false,
  delImg: false,
};

const start = async () => {
  console.log(`[程序所在目录]:`, EXE_DIR);
  let CFG = null;
  let CUSTOMER_CFG = null;
  try {
    CUSTOMER_CFG = JSON.parse(fs.readFileSync(`${EXE_DIR}/aikantu.json`));
  } catch (error) {
    CUSTOMER_CFG = {};
  }
  CFG = {
    ...DEFAULT_CFG,
    ...CUSTOMER_CFG,
  };
  console.log(`[配置信息]:\r\n`, CFG);
  console.log(`[图片目录]: ${CFG.dir}`);
  // 创建Koa实例
  const app = new Koa();
  const PORT = CFG.port || 3000;

  let DATA = null;
  try {
    DATA = JSON.parse(fs.readFileSync("./data.json"));
  } catch (error) {
    DATA = await getImgs(CFG);
  }
  // 托管静态资源
  app.use(
    heicConvertMiddleware({
      dir: CFG.dir, // 自定义文件目录
      quality: 1, // 自定义转换质量
    })
  );
  app.use(KoaStatic(`${__dirname}/public`));
  app.use(BodyParser());
  app.use(KoaStatic(CFG.dir));

  // 图片信息API端点
  app.use(async (ctx, next) => {
    if (ctx.path === "/getImgsList.json") {
      ctx.body = CFG.disorder ? shuffleArray(DATA) : DATA;
    }
    await next();
  });

  // 删除图片API端点
  app.use(async (ctx, next) => {
    if (ctx.path === "/delImg.json") {
      if (!CFG.delImg) {
        ctx.body = {
          code: "1",
          msg: `删除失败,配置文件未开启删除功能`,
        };
        return false;
      }
      const { path } = ctx.request.body;
      try {
        fs.unlinkSync(path);
        ctx.body = {
          code: "0",
          msg: "删除成功",
        };
      } catch (err) {
        ctx.body = {
          code: "1",
          msg: `删除失败,${err}`,
        };
      }
    }
    await next();
  });

  // 启动服务器
  app.listen(PORT, () => {
    console.log(`[服务器运行在]:http://localhost:${PORT}`);
    console.log(`[图片目录]: ${CFG.dir}`);
    // 使用默认浏览器打开
    c.exec(`start http://localhost:${PORT}`);
  });
};

start();

// esbuild index.js --bundle --minify --outfile=dist/dist.js --platform=node
// pkg -t win dist.js -o dist/dist.exe
