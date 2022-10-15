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
            <n-button v-if=!followed @click="newFollow" style="position: absolute; right: 40px; top: 25px; cursor: pointer;" round ghost color="#ED4557">
                <template #icon>
                    <n-icon>
                        <heart-outline />
                    </n-icon>
                </template>
                关注
            </n-button>
            <n-button v-else @click="unFollow" style="position: absolute; right: 40px; top: 25px; cursor: pointer;" round ghost color="#ED4557">
                <template #icon>
                    <n-icon>
                        <heart />
                    </n-icon>
                </template>
                已关注
            </n-button>
        </div>

        <div class="tabs">
            <n-card>
                <n-tabs type="line" >
                <n-tab-pane  name="articles" tab="TA的文章">
                    <div v-for="(article,index) in articles" style="margin-bottom:15px">
                        <n-card v-if="article.head_image" @click="toDetail(article)" style="cursor: pointer;" hoverable >
                            <n-image height="135" width="200" :src=serverUrl+article.head_image style="float: left" />
                            <div style="position: absolute; left: 240px; width: 690px;">
                                <text style="font-weight:bold; font-size: 20px;">{{article.title}}</text>
                                <p>{{article.content+"..."}}</p>
                                <div style="position: absolute; margin-top: 10px;">发布时间：{{article.created_at}}</div>
                            </div>
                        </n-card>
                        <n-card v-else style="cursor: pointer;" hoverable >
                            <div style="height: 140px; ">
                                <text @click="toDetail(article)" style="font-weight:bold; font-size: 20px;">{{article.title}}</text>
                                <p @click="toDetail(article)" >{{article.content+"..."}}</p>
                                <div style="position: absolute; margin-top: 10px;">发布时间：{{article.created_at}}</div>                        
                            </div>
                        </n-card>
                    </div>
                </n-tab-pane>
                <n-tab-pane  name="collects" tab="TA的收藏">
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
                                <p @click="toDetail(col)" >{{article.content+"..."}}</p>
                                <div style="position: absolute; margin-top: 10px;">发布时间：{{col.created_at}}</div>                        
                            </div>
                        </n-card>
                    </div>
                </n-tab-pane>
                <n-tab-pane name="following" tab="TA的关注">
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
import {HeartOutline} from '@vicons/ionicons5'
import {Heart} from '@vicons/ionicons5'

import {useRouter, useRoute} from 'vue-router'
const router = useRouter()
const route = useRoute()

const serverUrl = inject("serverUrl")
const axios = inject("axios")
const message = inject("message")

const user = reactive({
    self: false,
    avatarUrl: "",
    name: "",
    number: 0,
    fans: 0,
    id: 0,
    loginId: 0,
})
const articles = ref([])
const collects = ref([])
const following = ref([])
const followed = ref(false)
const index = ref(0)

onMounted(() => {
    loadDetailedInfo()
})

const loadDetailedInfo = async() => {
    let res1 = await axios.get("user/detailedInfo/" + route.query.id)        
    console.log(res1)
    if (res1.data.code == 200) {
        user.self = res1.data.data.self
        user.avatarUrl = serverUrl + res1.data.data.avatar
        user.name = res1.data.data.name
        user.number = res1.data.data.articles.length
        user.fans = res1.data.data.fans
        user.id = res1.data.data.id
        user.loginId = res1.data.data.loginId
        articles.value = res1.data.data.articles
        collects.value = res1.data.data.collects
        following.value = res1.data.data.following
        let res2 = await axios.get("following/" + route.query.id) 
        console.log(res2)
        if (res2.data.code == 200) {
            followed.value = res2.data.data.followed
            index.value = res2.data.data.index
        }
    }
} 

const newFollow = async() => {
    let res1 = await axios.put("following/new/" + route.query.id)
    console.log(res1)  
    if (res1.data.code == 200) {
        message.warning("已关注", {showIcon: false})  
        loadDetailedInfo()  
    }
}

const unFollow = async() => {
    let res1 = await axios.delete("following/" + index.value)
    console.log(res1)  
    if (res1.data.code == 200) {
        message.warning("取消关注", {showIcon: false})  
        loadDetailedInfo()  
    }
}

const toOtherUser = (fol) => {
    console.log(fol.id, user.loginId)
    if (fol.id == user.loginId) {
        router.push({
            path: "/myself",
            query: {
                id: fol.id
            }
        })   
    } else {
        router.push({
            path: "/others",
            query: {
                id: fol.id
            }
        })
    }
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