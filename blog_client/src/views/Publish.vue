<template>
    <div class="topbar">
        <n-button @click="goback" strong quaternary round style="position: absolute; left: 50px; top: 7px; font-size: 24px;" color="#7B3DE0">
            <n-icon>
                <return-up-back />
            </n-icon>
        </n-button>
        <text style="position:absolute; left: 200px; line-height: 50px; color: #383838">标题</text>
        <n-input v-model:value="addArticle.title" round placeholder="请输入标题" style="position:absolute; left: 265px; top: 8px; width: 1000px; background-color: #F3F0F9;" />
        <n-avatar round size="medium" :src=user.avatarUrl style="position: absolute; right: 190px; top: 8px; "/>
        <div style="position: absolute; right: 50px; top: 8px">
            <n-button round color="#7B3DE0" @click="showModalModal">
                <template #icon>
                    <n-icon>
                        <send />
                    </n-icon>
                </template>
                发布
            </n-button>
        </div>     
    </div>
    <div class="tabs">
        <n-card>
            <rich-text-editor v-model:modelValue="addArticle.content"></rich-text-editor>
        </n-card>
    </div>

    <n-modal v-model:show="showModal">
        <div style="width: 400px; height: 450px; background: white;">

            <n-card title="封面" :bordered="false" >
                <div v-if="!newHeadImage" style="width: 300px; height: 150px; margin: 0 auto;">
                    <n-upload
                    multiple
                    directory-dnd
                    :max="1"
                    @before-upload="beforeUpload"
                    :custom-request="customRequest"
                    >
                        <n-upload-dragger>
                            <div style="margin-bottom: 12px">
                                <n-icon size="48" :depth="3">
                                    <archive-icon />
                                </n-icon>
                            </div>
                            <n-text style="font-size: 16px">
                                点击或者拖动图片到此处
                            </n-text>
                        </n-upload-dragger>
                    </n-upload>
                </div>
                <div v-else style="width: 230px; margin: 0 auto;">
                    <n-image height="150" width="300" :src=serverUrl+addArticle.headImage />    
                    <n-button @click="deleteImage" circle style="position: absolute; left: 298px; top: 50px;" color="#383838">
                        <template #icon>
                            <n-icon><close /></n-icon>
                        </template>
                    </n-button>                
                </div>
            </n-card>

            <n-card title="分类" :bordered="false">
                <div style="width:300px; margin: 0 auto;">
                    <n-select v-model:value="addArticle.categoryId" :options="categoryOptions" placeholder="请选择分类"/>
                </div>
            </n-card>
            <div style="position: absolute; right: 100px; bottom: 30px;">
                <n-button type="default" @click="closeSubmitModal">
                    取消
                </n-button>
            </div>
            <div style="position: absolute; right: 30px; bottom: 30px;">
                <n-button type="primary" @click="submit">
                    确认
                </n-button>
            </div>                
        </div>
    </n-modal>

</template>

<script setup>
import {ref,reactive,inject, onMounted} from 'vue'
import RichTextEditor from '../components/RichTextEditor.vue'
import { ArchiveOutline as ArchiveIcon } from "@vicons/ionicons5"
import { Send } from "@vicons/ionicons5"
import { ReturnUpBack } from "@vicons/ionicons5"
import { Close } from "@vicons/ionicons5"

import {useRouter, useRoute} from 'vue-router'
const router = useRouter()
const route = useRoute()

const serverUrl = inject("serverUrl")
const axios = inject("axios")
const message = inject("message")

const login = ref(false)
const user = reactive({
    avatarUrl: "",
    id: 0
})
const categoryOptions = ref([])
const addArticle = reactive({
    id: 0,
    categoryId: 0,
    title:"",
    content:"",
    headImage:"",
})
const showModal = ref(false)
const newHeadImage = ref(false)

onMounted(() => {
    loadAvatar()
    loadCategories()
})

const loadAvatar= async() => {
    let res = await axios.get("/user")    
    console.log(res)
    if (res.data.code == 200) {
        user.avatarUrl = serverUrl + res.data.data.avatar
        user.id = res.data.data.id
        login.value = true
    }
} 

const loadCategories = async() =>{
    let res = await axios.get("/category")
    console.log(res)
    categoryOptions.value = res.data.data.categories.map((item)=>{
      return {
        label:item.name,
        value:item.id
      }
    })
}

const showModalModal = () => {
    showModal.value = true
}

const closeSubmitModal = () => {
    showModal.value = false  
}

const beforeUpload = async(data) => {
    if (data.file.file?.type !== "image/png") {
        message.error("只能上传png格式的图片")
        return false;
    }
    return true;
}

const customRequest = async({file}) => {
    const formData = new FormData()
    formData.append('file', file.file)
    let res = await axios.post("/upload", formData)
    console.log(res)
    addArticle.headImage = res.data.data.filePath
    newHeadImage.value = true
}

const deleteImage = () => {
    addArticle.headImage = ""
    newHeadImage.value = false
}

const submit = async() => {
    let res = await axios.post("/article", {
        category_id: addArticle.categoryId,
        title: addArticle.title,
        content: addArticle.content,
        head_image: addArticle.headImage
    })
    console.log(res)   
    if (res.data.code == 200) {
        message.success(res.data.msg) 
        goback()  
    } else {
        message.error(res.data.msg)
    }
}

const goback= () => {
    router.go(-1)    
}


</script>

<style lang="scss" scoped>
.topbar {
    position: sticky;
    top: 0;
    height: 50px;
    background: white;
    box-shadow: 0px 1px 5px #D3D4D8;
}
.tabs {
    position: absolute;
    top: 75px;
    left: 0;
    right: 0;
    margin: auto;
    width: 1000px;
    height: auto;
    background: white;  
    box-shadow: 0px 1px 3px #D3D4D8; 
    border-radius: 5px; 
}
</style>