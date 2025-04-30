export const getSlideStyle = (slide: any) => {
  const backgroundStyle = slide.background;
  let background = '';

  if (backgroundStyle.type === 'image') {
    background = `url(${backgroundStyle.image.src}) no-repeat center center/cover`;
  } else if (backgroundStyle.type === 'gradient') {
    const gradient = backgroundStyle.gradient;
    const colors = gradient.colors.map((color: any) => `${color.color} ${color.pos}%`).join(', ');
    background = `linear-gradient(${gradient.rotate}deg, ${colors})`;
  } else if (backgroundStyle.type === 'solid') {
    background = backgroundStyle.color;
  }

  return {
    backgroundColor: background,
  };
};

export const getElementStyle = (element: any) => {
  return {
    width: `${element.width}px`,
    height: `${element.height}px`,
    left: `${element.left}px`,
    top: `${element.top}px`,
    transform: `rotate(${element.rotate}deg)`,
  };
};

export const getTextStyle = (element: any):any => {
  const shadow = element.shadow ? `${element.shadow.h}px ${element.shadow.v}px ${element.shadow.blur}px ${element.shadow.color}` : 'none';
  return {
    fontFamily: element.defaultFontName,
    color: element.defaultColor,
    lineHeight: element.lineHeight,
    textAlign: element.vertical ? 'center' : 'left',
    textShadow: shadow,
    outline: `${element.outline.width}px ${element.outline.style} ${element.outline.color}`,
    fill: element.fill,
  };
};

export const getImageStyle = (element: any) => {
  return {
    width: `${element.width}px`,
    height: `${element.height}px`,
    transform: `rotate(${element.rotate}deg)`,
    transformOrigin: 'center',
    filter: element.flipH ? 'scaleX(-1)' : '' + element.flipV ? 'scaleY(-1)' : '',
  };
};

export const getAudioStyle = (element: any) => {
  return {
    width: `${element.width}px`,
    height: `${element.height}px`,
  };
};

export const getVideoStyle = (element: any) => {
  return {
    width: `${element.width}px`,
    height: `${element.height}px`,
  };
};

export const getShapeStyle = (element: any) => {
  const shadow = element.shadow ? `${element.shadow.h}px ${element.shadow.v}px ${element.shadow.blur}px ${element.shadow.color}` : 'none';
  return {
    width: `${element.width}px`,
    height: `${element.height}px`,
    transform: `rotate(${element.rotate}deg)`,
    transformOrigin: 'center',
    filter: element.flipH ? 'scaleX(-1)' : '' + element.flipV ? 'scaleY(-1)' : '',
    boxShadow: shadow,
  };
};

export const getTableStyle = (element: any):any => {
  return {
    width: `${element.width}px`,
    height: `${element.height}px`,
    borderCollapse: 'collapse',
  };
};

export const getCellStyle = (element:any, cell: any) => {
  return {
    textAlign: cell.style.align,
    fontSize: cell.style.fontsize,
    fontFamily: cell.style.fontname,
    color: cell.style.color,
    fontWeight: cell.style.bold ? 'bold' : 'normal',
    backgroundColor: cell.style.backcolor,
    border: `${element.outline.width}px ${element.outline.style} ${element.outline.color}`,
  };
};