<script setup lang="ts">
import { t} from "@/i18n/index";
import { useModelStore } from "@/stores/model";
const props = defineProps({
    model:{
        type:String,
        default:""
    }
})
const model = props.model;
const modelStore = useModelStore();
const modelInfo = modelStore.getModelInfo(model);
</script>
<template>
    <el-scrollbar>
    <div class="model-info">
    
        <el-row justify="space-around">
            <el-col :span="10" class="tc"><el-text>{{ t('model.modelNames') }}</el-text></el-col>
            <el-col :span="14">
                <el-tag type="success">{{model}}</el-tag>
            </el-col>
        </el-row>
        <el-row justify="space-around">
            <el-col :span="10" class="tc"><el-text>{{ t('model.modelSize') }}</el-text></el-col>
            <el-col :span="14">
                <el-tag type="info">{{modelInfo.info.size}}</el-tag>
            </el-col>
        </el-row>
        <el-row justify="space-around">
            <el-col :span="10" class="tc"><el-text>{{ t('model.modelEngine') }}</el-text></el-col>
            <el-col :span="14">
                <el-tag type="primary">{{modelInfo.engine}}</el-tag>
            </el-col>
        </el-row>
        <el-row justify="space-around" v-if="modelInfo.action">
            <el-col :span="10" class="tc"><el-text>{{ t('model.applicableScope') }}</el-text></el-col>
            <el-col :span="14">
                <el-tag type="primary" v-for="item in modelInfo.action" style="margin-right: 5px;">{{t('model.' + item)}}</el-tag>
            </el-col>
        </el-row>
        <el-row justify="space-around" v-if="modelInfo.info.context_length">
            <el-col :span="10" class="tc"><el-text>{{ t('model.contextLengths') }}</el-text></el-col>
            <el-col :span="14">
                <el-tag type="primary">{{modelInfo.info.context_length}}</el-tag>
            </el-col>
        </el-row>
        <el-row justify="space-around" v-if="modelInfo.info.pb">
            <el-col :span="10" class="tc"><el-text>{{ t('model.parameterSize') }}</el-text></el-col>
            <el-col :span="14">
                <el-tag type="primary">{{modelInfo.info.pb}}</el-tag>
            </el-col>
        </el-row>
        <el-row justify="space-around" v-if="modelInfo.info.cpu">
            <el-col :span="10" class="tc"><el-text>{{ t('model.requiredCPU') }}</el-text></el-col>
            <el-col :span="14">
                <el-tag type="danger">{{modelInfo.info.cpu}}</el-tag>
            </el-col>
        </el-row>
        <el-row justify="space-around" v-if="modelInfo.info.gpu">
            <el-col :span="10" class="tc"><el-text>{{ t('model.requiredGPU') }}</el-text></el-col>
            <el-col :span="14">
                <el-tag type="warning">{{modelInfo.info.gpu}}</el-tag>
            </el-col>
        </el-row>
        <el-row justify="space-around" v-if="modelInfo.info.template">
            <el-col :span="10" class="tc"><el-text>{{ t('model.modelTemplate') }}</el-text></el-col>
            <el-col :span="14">
                <el-input type="textarea" :row="3" v-model="modelInfo.info.template"></el-input>
            </el-col>
        </el-row>
        <el-row justify="space-around" v-if="modelInfo.info.parameters">
            <el-col :span="10" class="tc"><el-text>{{ t('model.modelParameters') }}</el-text></el-col>
            <el-col :span="14">
                <el-input type="textarea" :row="3" v-model="modelInfo.info.parameters"></el-input>
            </el-col>
        </el-row>
    
  </div>
</el-scrollbar>
</template>
<style>
.model-info{
    padding:20px;
}
.el-row{
    margin:10px;
    border-bottom: 1px solid #eee;
    padding-bottom: 10px;
}
.tc{
    text-align:left
}
</style>
