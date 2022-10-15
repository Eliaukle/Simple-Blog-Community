<template>
    <div>
        <div style="position: fixed; top: 0; height: 50px; width: 100%; z-index: 999; "><TopBar></TopBar></div>

        <div class="card">
            <div style="position: absolute; left: 40px; bottom: 20px">
                <n-avatar round :size="120" :src=user.avatarUrl :bordered=true />
            </div>
            <div style="position: absolute; top: 25px;left: 200px; font-size: 20px;">{{user.name}}</div>
            <div style="position: absolute; top: 70px;left: 200px;">
                <text style="font-weight:bold; font-size: 20px;">{{user.number}}</text>
                <text style="margin-left: 5px; font-size: 14px;">文章</text>
                <text style="font-weight:bold; font-size: 20px; margin-left: 20px;">{{user.fans}}</text>
                <text style="margin-left: 5px; font-size: 14px;">粉丝</text>
            </div>
            <n-dropdown trigger="hover" :options="options" @select="handleSelect">
                <n-button style="position: absolute; right: 40px; top: 25px;" round ghost color="#7B3DE0">修改资料</n-button>
            </n-dropdown>
        </div>

        <n-modal v-model:show="showAvatarModal">
            <div style="width: 600px; height: 320px; background: white;">
                <n-card title="修改头像" :bordered="false">
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
                </n-card>
                <div style="position: absolute; right: 90px; bottom: 20px;">
                    <n-button type="default" @click="closeAvatarModal">
                        取消
                    </n-button>
                </div>
                <div style="position: absolute; right: 20px; bottom: 20px;">
                    <n-button v-if="newAvatar" @click="modifyAvatar" type="primary">
                        确认
                    </n-button>
                    <n-button v-else type="primary" disabled>
                        确认
                    </n-button>
                </div>
            </div>
        </n-modal>

        <n-modal v-model:show="showNameModal">
            <div style="width: 440px; height: 185px; background: white;">
                <n-card title="修改用户名" :bordered="false">
                    <div style="width:350px;">
                        <n-input v-model:value="newName" size="large" round type="text" placeholder="请输入用户名" />
                    </div>
                </n-card>
                <div style="position: absolute; right: 90px; bottom: 20px;">
                    <n-button type="default" @click="closeNameModal">
                        取消
                    </n-button>
                </div>
                <div style="position: absolute; right: 20px; bottom: 20px;">
                    <n-button type="primary" @click="modifyName">
                        确认
                    </n-button>
                </div>                
            </div>
        </n-modal>

        <div class="tabs">
            <n-card>
                <n-tabs type="line" >
                <n-tab-pane  name="articles" tab="我的文章">
                    <div v-for="(article,index) in articles" style="margin-bottom:15px">
                        <n-card v-if="article.head_image" @click="toDetail(article)" style="cursor: pointer;" hoverable >
                            <n-image height="135" width="200" :src=serverUrl+article.head_image style="float: left" />
                            <div style="position: absolute; left: 240px; width: 690px;">
                                <text  style="font-weight:bold; font-size: 20px;">{{article.title}}</text>
                                <p >{{article.content+"..."}}</p>
                                <div style="position: absolute; margin-top: 10px;">发布时间：{{article.created_at}}</div>
                            </div>
                        </n-card>
                        <n-card v-else @click="toDetail(article)" style="cursor: pointer;" hoverable >
                            <div style="height: 140px; ">
                                <text style="font-weight:bold; font-size: 20px;">{{article.title}}</text>
                                <p >{{article.content+"..."}}</p>
                                <div style="position: absolute; margin-top: 10px;">发布时间：{{article.created_at}}</div>                           
                            </div>
                        </n-card>
                    </div>
                </n-tab-pane>
                <n-tab-pane  name="collects" tab="我的收藏">
                    <div v-for="(col,index) in collects" style="margin-bottom:15px">
                        <n-card v-if="col.head_image" @click="toDetail(col)" style="cursor: pointer;" hoverable >
                            <n-image height="135" width="200" :src=serverUrl+col.head_image style="float: left" />
                            <div style="position: absolute; left: 240px; width: 690px;">
                                <text style="font-weight:bold; font-size: 20px;">{{col.title}}</text>
                                <p>{{col.content+"..."}}</p>
                                <div style="position: absolute; margin-top: 10px;">发布时间：{{col.created_at}}</div>
                            </div>
                        </n-card>
                        <n-card v-else style="cursor: pointer;" hoverable >
                            <div style="height: 140px; ">
                                <text @click="toDetail(col)" style="font-weight:bold; font-size: 20px;">{{col.title}}</text>
                                <p @click="toDetail(col)" >{{col.content+"..."}}</p>
                                <div style="position: absolute; margin-top: 10px;">发布时间：{{col.created_at}}</div>                        
                            </div>
                        </n-card>
                    </div>
                </n-tab-pane>
                <n-tab-pane name="following" tab="我的关注">
                    <div v-for="(fol,index) in following" style="margin-bottom:15px">
                        <n-card>
                            <n-avatar @click="toOtherUser(fol)" round size="large" :src=serverUrl+fol.avatar style="float: left; cursor: pointer;" />
                            <text style="position: absolute; left: 90px; top: 25px; font-size: 20px;">{{fol.userName}}</text>
                        </n-card>
                    </div>
                </n-tab-pane>
                </n-tabs>
            </n-card>
        </div>
    </div>
</template>

<script setup>
import {ref,reactive,inject, onMounted} from 'vue'
import TopBar from '../components/TopBar.vue'
import { ArchiveOutline as ArchiveIcon} from "@vicons/ionicons5"

import {useRouter, useRoute} from 'vue-router'
const router = useRouter()
const route = useRoute()

const serverUrl = inject("serverUrl")
const axios = inject("axios")
const message = inject("message")
const options = reactive([
    {label: "修改头像", key: "avatar"},
    {label: "修改用户名", key: "name"},
])
const user = reactive({
    self: false,
    avatarUrl: "",
    name: "",
    number: 0,
    fans: 0,
    id: 0
})
const newUrl = ref("")
const newAvatar = ref(false)
const newName = ref("")
const showAvatarModal = ref(false)
const showNameModal = ref(false)
const articles = ref([])
const collects = ref([])
const following = ref([])

onMounted(() => {
    loadDetailedInfo()
})

const loadDetailedInfo= async() => {
    let res = await axios.get("user/detailedInfo/" + route.query.id)        
    console.log(res)
    if (res.data.code == 200) {
        user.self = res.data.data.self
        user.avatarUrl = serverUrl + res.data.data.avatar
        user.name = res.data.data.name
        user.number = res.data.data.articles.length
        user.fans = res.data.data.fans
        user.id = res.data.data.id
        articles.value = res.data.data.articles
        collects.value = res.data.data.collects
        following.value = res.data.data.following
        newName.value =  user.name
    }
} 

const handleSelect = (key) => {
    if (String(key) == "avatar") {
        showAvatarModal.value = true
    } 
    if (String(key) == "name") {
        showNameModal.value = true
    } 
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
    newUrl.value = res.data.data.filePath
    newAvatar.value = true
}

const modifyAvatar = async() => {
    let res = await axios.put("user/avatar/" + route.query.id,
    {
        avatar: newUrl.value,
    })     
    console.log(res) 
    if (res.data.code == 200) {
        message.success(res.data.msg)    
        showAvatarModal.value = false  
        loadDetailedInfo()
    }  else {
        message.error(res.data.msg)  
    }
}

const modifyName = async() => {
    let res = await axios.put("user/name/" + route.query.id,
    {
        userName: newName.value,
    })     
    console.log(res) 
    if (res.data.code == 200) {
        message.success(res.data.msg)    
        showNameModal.value = false  
        loadDetailedInfo()
    }  else {
        message.error(res.data.msg)  
    }
}

const closeAvatarModal = () => {
    showAvatarModal.value = false  
}

const closeNameModal = () => {
    showNameModal.value = false  
}

const toOtherUser = (fol) => {
    router.push({
        path: "/others",
        query: {
            id: fol.id
        }
    }) 
}

const toDetail = (article) => {
    router.push({
        path: "/detail",
        query: {
            id: article.id
        }
    }) 
}

</script>

<style lang="scss" scoped>
.card {
    position: absolute;
    top: 100px;
    left: 0;
    right: 0;
    margin: auto;
    width: 1000px;
    height: 130px;
    background: white;  
    box-shadow: 0px 1px 3px #D3D4D8; 
    border-radius: 5px;
}
.tabs {
    position: absolute;
    top: 250px;
    left: 0;
    right: 0;
    margin: auto;
    width: 1000px;
    height: auto;
    background: white;  
    box-shadow: 0px 1px 3px #D3D4D8; 
    border-radius: 5px; 
}
// .head_image {
//     float: left;
//     width: 20%;
// }
.cardInfo {
    float: right;
    width: 80%;
}
</style>