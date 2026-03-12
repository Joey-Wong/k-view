export const to = (promise) => promise.then((data) => [null, data]).catch((err) => [err]);
export const getImgWH = (img) => {
  const image = new Image();
  image.crossOrigin = "";
  image.src = img;
  return new Promise((resolve, reject) => {
    image.onload = function () {
      const { width, height } = image;
      resolve({ width, height });
    };
  });
};

export const shuffleArray = (array) => {
  const newArray = [...array];
  for (let i = newArray.length - 1; i > 0; i--) {
    const j = Math.floor(Math.random() * (i + 1));
    [newArray[i], newArray[j]] = [newArray[j], newArray[i]];
  }
  return newArray;
};
