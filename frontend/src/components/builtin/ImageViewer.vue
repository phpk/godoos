<template>
  <div class="viewer">
    <div class="loading" v-if="!content">
      Loading
      <WinLoading></WinLoading>
    </div>
    <img class="viewer-img" v-else :src="content" />
  </div>
</template>
<script setup lang="ts">
import { inject, onMounted, ref } from "vue";
import { OsFileSystem, System, BrowserWindow } from "@/system";

const browserWindow: BrowserWindow | undefined = inject("browserWindow");
const sys = inject<System>("system")!;
const content = ref("");

onMounted(() => {
  const fileContent = browserWindow?.config.content;

  if (sys.fs instanceof OsFileSystem) {
    sys.fs.checkVolumePath(fileContent);
  }
  if (typeof fileContent === "string" && fileContent?.startsWith("http")) {
    // http://admin:admin@example.com:5244/dav/a.jpg to admin:admin
    const url = new URL(fileContent);
    const cred = url.password ? `${url.username}:${url.password}` : "";
    if (cred) {
      const resource = url.href.replace(`${url.username}:${url.password}@`, "");
      const headers = new Headers({
        Authorization: `Basic ${btoa(cred)}`,
      });
      fetch(resource, {
        headers: headers,
      }).then((res) => {
        res.blob().then((blob) => {
          const reader = new FileReader();
          reader.readAsDataURL(blob);
          reader.onload = function (e) {
            content.value = e.target?.result as string;
          };
        });
      });
      return;
    } else {
      content.value = fileContent;
    }
  } else {
    //console.log(browserWindow?.config)
    const ext = browserWindow?.config.path.split(".").pop();
    //let data = JSON.parse(fileContent)
    if (typeof fileContent === "string") {
      content.value = `data:image/${ext};base64,` + fileContent;
    }
    if (fileContent instanceof ArrayBuffer) {
      //console.log(fileContent)
      const blob = new Blob([fileContent]);
      const reader = new FileReader();
      reader.onload = function (e: any) {
        console.log(e.target.result);
        content.value = `data:image/${ext};base64,` + e.target.result.split(",")[1];
        //console.log(base64String);
      };
      reader.readAsDataURL(blob);
    }

    //content.value = data.content;
    return;
  }
});
</script>
<style scoped>
.viewer {
  width: 100%;
  height: 100%;
  background-color: #000;
}
.loading {
  width: 100%;
  height: 100%;
}
.viewer-img {
  width: 100%;
  height: 100%;
  object-fit: contain;
}
</style>
