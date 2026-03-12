const fs = require('fs');

const to = (promise) => promise.then((data) => [null, data]).catch((err) => [err]);
const getFileInfo = (filePath) => {
  // 找到最后一个斜杠的索引
  const lastSlashIndex = filePath.lastIndexOf("/");
  // 获取所在路径
  const directory = filePath.slice(0, lastSlashIndex);
  // 获取文件名
  const fullFileName = filePath.slice(lastSlashIndex + 1);
  // 找到文件名中最后一个点的索引
  const lastDotIndex = fullFileName.lastIndexOf(".");
  // 获取文件后缀，如果没有点则为空字符串
  const fileExtension = lastDotIndex !== -1 ? fullFileName.slice(lastDotIndex + 1) : "";
  // 获取不包含后缀的文件名
  const fileNameWithoutExtension = lastDotIndex !== -1 ? fullFileName.slice(0, lastDotIndex) : fullFileName;
  return {
    ext: fileExtension,
    name: fileNameWithoutExtension,
    dir: directory,
  };
};

const getFiles = (dir, isDeep) => {
  const files = fs.readdirSync(dir);
  const filesPath = [];
  for (let filaName of files) {
    const fPath = `${dir}/${filaName}`;
    const stats = fs.statSync(fPath);
    if (stats.isDirectory()) {
      isDeep && filesPath.push(...getFiles(fPath, isDeep));
    } else {
      // 存储文件路径和修改时间
      filesPath.push({ path: fPath, mtime: stats.mtime });
    }
  }

  // 提取排序后的路径数组
  return filesPath;
};

const getFilesSortByMTime = (dir, isDeep) => {
  const arr = getFiles(dir, isDeep);
  arr.sort((a, b) => b.mtime - a.mtime);
  // 提取排序后的路径数组
  return arr.map(file => file.path);
};

const shuffleArray = (array) => {
  const newArray = [...array];
  for (let i = newArray.length - 1; i > 0; i--) {
    const j = Math.floor(Math.random() * (i + 1));
    [newArray[i], newArray[j]] = [newArray[j], newArray[i]];
  }
  return newArray;
};

module.exports = { getFilesSortByMTime, to, getFileInfo, shuffleArray };
