const { getFilesSortByMTime } = require('./index');
const path = require('path');
const { imageSizeFromFile } = require('image-size/fromFile')
const getImgs = async (cfg) => {
    const files = getFilesSortByMTime(cfg.dir, true, cfg.exts);
    const dirLen = cfg.dir.length;


    const imgs = files.filter((file) => {
        const ext = path.extname(file);
        return cfg.exts.includes(ext);
    });
    const len = imgs.length;
    const data = [];

    let i = 0;
    while (i < len) {
        const img = imgs[i];
        const dimensions = await imageSizeFromFile(img);
        const item = {
            pic: encodeURI(img.slice(dirLen)), // 可访问的URL路径
            picW: dimensions.width,
            picH: dimensions.height,
            path: img,
        };
        console.log(img);
        data.push(item);
        i++;
    }

    return data;
}

module.exports = { getImgs };