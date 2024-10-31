  <script setup>
  import { ref } from "vue";
  import { useChatStore } from "@/stores/chat";
  const store = useChatStore()
  const form = ref({
      nickname: '',
  })

  const onSubmit = () => {
      console.log('submit!')
  }

  </script>
  <template>
      <el-form :model="form" label-width="160px" style="padding: 30px;">
          <el-form-item label="头像">
              <el-avatar shape="square" :size="80" :src="store.userInfo.avatar"/>
          </el-form-item>
          <el-form-item label="昵称">
              <el-input v-model="store.userInfo.nickname" />
          </el-form-item>
          <el-form-item label="手机号">
              <el-input v-model="store.userInfo.phone" />
          </el-form-item>
          <el-form-item label="工号">
              <el-input v-model="store.userInfo.job_number" />
          </el-form-item>
          <el-form-item label="自我介绍">
              <el-input v-model="store.userInfo.desc" />
          </el-form-item>
          <el-form-item>
              <el-button type="primary" @click="onSubmit">保存</el-button>
          </el-form-item>
      </el-form>
  </template>
