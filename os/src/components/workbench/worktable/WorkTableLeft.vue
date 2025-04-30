    <template>
        <el-menu class="el-menu-vertical-demo" :default-openeds="opens" :default-active="store.currentTab">
            <el-sub-menu index="my">
                <template #title>
                    <el-icon>
                        <UserFilled />
                    </el-icon>
                    <span>我的</span>
                </template>
                <el-menu-item index="0" @click="store.setCurrentTab('0')">
                    <el-icon>
                        <Compass />
                    </el-icon>
                    <span>我的任务</span>
                </el-menu-item>
                <el-menu-item index="1" @click="store.setCurrentTab('1')">
                    <el-icon>
                        <Open />
                    </el-icon>
                    <span>我的申请</span>
                </el-menu-item>
                <el-menu-item index="2" @click="store.setCurrentTab('2')">
                    <el-icon>
                        <Tickets />
                    </el-icon>
                    <span>任务列表</span>
                </el-menu-item>
                <el-menu-item index="3" @click="store.setCurrentTab('3')">
                    <el-icon>
                        <Memo />
                    </el-icon>
                    <span>申请列表</span>
                </el-menu-item>
                <el-menu-item index="4" @click="store.setCurrentTab('4')">
                    <el-icon>
                        <Files />
                    </el-icon>
                    <span>历史抄送</span>
                </el-menu-item>
            </el-sub-menu>
            <el-sub-menu index="approval">
                <template #title>
                    <el-icon><icon-menu /></el-icon>
                    <span>审批</span>
                </template>
                <el-menu-item index="5" @click="store.setCurrentTab('5')">
                    <el-icon>
                        <Reading />
                    </el-icon>
                    <span>我的审批</span>
                </el-menu-item>
                <el-menu-item index="6" @click="store.setCurrentTab('6')">
                    <el-icon>
                        <Postcard />
                    </el-icon>
                    <span>我的抄送</span>
                </el-menu-item>
            </el-sub-menu>
        </el-menu>
    </template>

<script lang="ts" setup>
import {
    Menu as IconMenu,
} from '@element-plus/icons-vue'
import { useWorkTableStore } from '@/stores/worktable';
const opens = ['my', 'approval']
const store = useWorkTableStore()
</script>



<style scoped lang="scss">
.el-menu {
    width: 100%;
    border-left: 1px solid var(--el-menu-border-color);
    border-right: 0;
}
</style>