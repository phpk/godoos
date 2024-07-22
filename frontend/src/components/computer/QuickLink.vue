<template>
  <div class="contaion">
    <div class="showName">{{ i18n("favorite") }}</div>

    <div
      class="file-item mode-list"
      v-for="item in quickLinks"
      :key="item.path"
      @click="quick(item)"
    >
      <div class="file-item_img">
        <!-- <img :src="item.img" alt="" /> -->
        <svg class="icon" aria-hidden="true">
          <use :xlink:href="'#icon-' + item.img"></use>
        </svg>
      </div>
      <div class="file-item_title">
        {{ item.name }}
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
import { i18n } from "@/i18n";
const props = defineProps({
  onOpen: {
    type: Function,
    default: () => {
      //
    },
  },
});
const quickLinks = [
  {
    img: "word",
    name: i18n("document"),
    path: "/C/Users/Doc",
  },
  {
    img: "excel",
    name: i18n("excel"),
    path: "/C/Users/Execl",
  },
  {
    img: "pptexe",
    name: i18n("ppt"),
    path: "/C/Users/PPT",
  },
  {
    img: "markdown",
    name: i18n("markdown"),
    path: "/C/Users/Markdown",
  },
  {
    img: "mindexe",
    name: i18n("mindmap"),
    path: "/C/Users/Mind",
  },
  {
    img: "baiban",
    name: i18n("whiteBoard"),
    path: "/C/Users/Baiban",
  },
  {
    img: "kanban",
    name: i18n("board"),
    path: "/C/Users/Kanban",
  },
  {
    img: "gallery",
    name: i18n("gallery"),
    path: "/C/Users/Photo",
  },
  {
    img: "music",
    name: i18n("music"),
    path: "/C/Users/Music",
  },
];
function quick(item: any) {
  props.onOpen(item.path);
}
</script>
<style lang="scss" scoped>
.contaion {
  width: 92%;
  padding: 4%;
}
.file-item-title {
  font-size: 11px;
  text-indent: 5px;
}
.file-item {
  position: relative;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  width: 100%;
  color: var(--icon-title-color);

  font-size: var(--ui-font-size);
  border: 1px solid transparent;
  margin: 2px;
  user-select: none;
}
.file-item:hover {
  background-color: #b1f1ff4c;
}

.mode-list {
  display: flex;
  flex-direction: row;
  justify-content: flex-start;
  height: calc(0.8 * var(--menulist-item-height));
  width: 100%;
  overflow: hidden;

  .file-item_img {
    width: var(--menulist-item-height);
    height: calc(0.6 * var(--menulist-item-height));
    // margin: 0px auto;
    user-select: none;
    flex-shrink: 0;
    svg {
      width: 100%;
      height: 100%;
      object-fit: contain;
    }
  }
  .file-item_title {
    height: min-content;

    word-break: break-all;
    overflow: hidden;
    text-overflow: ellipsis;
    word-break: break-all;
    display: -webkit-box;
    -webkit-box-orient: vertical;
    -webkit-line-clamp: 1;
  }
}
.showName {
  font-size: 11px;
  text-indent: 5px;
}
</style>
