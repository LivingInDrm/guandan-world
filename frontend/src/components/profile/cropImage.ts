export interface CropArea {
  x: number;
  y: number;
  width: number;
  height: number;
}

const OUTPUT_SIZE = 512;

export async function getCroppedImage(
  imageSrc: string,
  cropArea: CropArea,
  fileName = 'avatar.png'
): Promise<File> {
  const image = await createImage(imageSrc);
  const canvas = document.createElement('canvas');
  const ctx = canvas.getContext('2d');

  if (!ctx) {
    throw new Error('Failed to get canvas context');
  }

  canvas.width = OUTPUT_SIZE;
  canvas.height = OUTPUT_SIZE;

  ctx.drawImage(
    image,
    cropArea.x,
    cropArea.y,
    cropArea.width,
    cropArea.height,
    0,
    0,
    OUTPUT_SIZE,
    OUTPUT_SIZE
  );

  return new Promise((resolve, reject) => {
    canvas.toBlob(
      (blob) => {
        if (blob) {
          resolve(new File([blob], fileName, { type: 'image/png' }));
        } else {
          reject(new Error('Failed to create blob'));
        }
      },
      'image/png',
      1
    );
  });
}

function createImage(src: string): Promise<HTMLImageElement> {
  return new Promise((resolve, reject) => {
    const img = new Image();
    img.onload = () => resolve(img);
    img.onerror = (error) => reject(error);
    if (!src.startsWith('blob:')) {
      img.crossOrigin = 'anonymous';
    }
    img.src = src;
  });
}
