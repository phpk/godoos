<template>
  <div class="view-ppt" ref="wrapperRef">
    <div v-html="htmlContent" class="content-container"></div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { parse } from 'pptxtojson';

const wrapperRef = ref(null);
const htmlContent = ref('');

const props = defineProps<{
  content: any;
}>();

const getStyleString = (element: any): string => {
  return `
    position: absolute;
    left: ${element.left}px;
    top: ${element.top}px;
    width: ${element.width}px;
    height: ${element.height}px;
    transform: rotate(${element.rotate}deg);
  `;
};

const renderText = (element: any): string => {
  return `<div style="${getStyleString(element)}; font-size: ${element.fontSize}pt; color: ${element.fillColor}; text-align: center;">${element.content}</div>`;
};

const renderImage = (element: any): string => {
  return `<img src="${element.src}" style="${getStyleString(element)}" />`;
};

const renderShape = (element: any): string => {
  return `<div style="${getStyleString(element)}; background-color: ${element.fillColor}; border: ${element.borderWidth}px ${element.borderType} ${element.borderColor};"></div>`;
};

const renderTableCell = (cell: any): string => {
  return `<td style="width: ${cell.width}px; height: ${cell.height}px; background-color: ${cell.fillColor}; color: ${cell.fontColor}; text-align: ${cell.textAlign || 'center'};">${cell.content}</td>`;
};

const renderTable = (element: any): string => {
  let tableHtml = `<table style="${getStyleString(element)}">`;
  element.data.forEach((row: any) => {
    tableHtml += `<tr>`;
    row.forEach((cell:any) => {
      tableHtml += renderTableCell(cell);
    });
    tableHtml += `</tr>`;
  });
  tableHtml += `</table>`;
  return tableHtml;
};

const renderAudio = (element: any): string => {
  return `<audio controls style="${getStyleString(element)}">
    <source src="${element.blob}" type="audio/${element.type}">
    您的浏览器不支持 audio 元素。
  </audio>`;
};

const renderVideo = (element: any): string => {
  return `<video controls style="${getStyleString(element)}">
    <source src="${element.blob}" type="video/${element.type}">
    您的浏览器不支持 video 元素。
  </video>`;
};

const renderGroup = (group: any): string => {
  return group.elements.map(renderElement).join('');
};

const renderElement = (element: any): string => {
  switch (element.type) {
    case 'text': return renderText(element);
    case 'image': return renderImage(element);
    case 'shape': return renderShape(element);
    case 'table': return renderTable(element);
    case 'audio': return renderAudio(element);
    case 'video': return renderVideo(element);
    case 'group': return renderGroup(element);
    default: return '';
  }
};

const importPPTX = async (content: any) => {
  try {
    const json = await parse(content);
    const elements = json.slides.flatMap(slide => slide.elements);
    const htmlString = elements.map(renderElement).join('');
    htmlContent.value = htmlString;
  } catch (error) {
    console.error('无法正确读取 / 解析PPTX文件:', error);
  }
};

onMounted(async () => {
  await importPPTX(props.content);
});
</script>

<style scoped>
.view-ppt {
  position: relative;
  margin: auto;
  text-align: center;
  overflow: auto;
  width: 100%;
  height: 100vh; /* 或者您想要的任何高度 */
}

.content-container {
  position: relative;
  width: 100%;
  height: 100%;
}
</style>