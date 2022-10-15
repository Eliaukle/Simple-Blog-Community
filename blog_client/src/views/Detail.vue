<template>
    <div style="position: fixed; top: 0; height: 50px; width: 100%; z-index: 999; "><TopBar></TopBar></div>

    <div class="tabs">
        <n-card>
            <n-h1>{{articleInfo.title}}</n-h1>
            <div style="height: 75px; background-color: #FCFAF7;">
                <n-avatar @click="toOtherUser" round size="medium" :src=userUrl style="position: relative; left: 20px; top: 20px; cursor: pointer;"/>
                <text style="position: relative; left: 36px; color: #808080;">发布时间：{{articleInfo.createdAt}} </text>
                <div style="position: relative; left: 70px; color: #808080;">
                    文章分类：
                    <n-tag type="warning">{{categoryName}}</n-tag>
                </div>
                <n-button v-if="self" @click="toUpdate" ghost style="bottom: 45px; left: 805px;" color="#7B3DE0">修改</n-button>
                <n-button v-if="self" @click="toDelete" ghost style="bottom: 45px; left: 815px;" color="#7B3DE0">删除</n-button>
            </div>
            <n-divider />
            <n-button v-if=!collected text @click="newCollect" style="position: absolute; right: 40px; top: 25px; cursor: pointer;" round ghost color="#FFA876">
                <template #icon>
                    <n-icon>
                        <star-outline />
                    </n-icon>
                </template>
                收藏
            </n-button>
            <n-button v-else text @click="unCollect" style="position: absolute; right: 40px; top: 25px; cursor: pointer;" round ghost color="#FFA876">
                <template #icon>
                    <n-icon>
                        <star />
                    </n-icon>
                </template>
                已收藏
            </n-button>
            <div class="article-content">
                <div v-html="articleInfo.content"></div>
            </div>
        </n-card>
    </div>
</template>

<script setup>
import {ref,reactive,inject, onMounted} from 'vue'
import TopBar from '../components/TopBar.vue'
import { Star } from "@vicons/ionicons5"
import { StarOutline } from "@vicons/ionicons5"

import {useRouter, useRoute} from 'vue-router'
const router = useRouter()
const route = useRoute()

const serverUrl = inject("serverUrl")
const axios = inject("axios")
const message = inject("message")
const dialog = inject("dialog")

const articleInfo = ref({})
const categoryName = ref("")
const user = ref({})
const userUrl = ref("")
const collected = ref(false)
const index = ref(0)
const self = ref(false)

onMounted(() => {
    loadArticle()
})

const loadArticle = async() => {
    let res1 = await axios.get("article/" + route.query.id)
    console.log(res1)
    if (res1.data.code == 200) {
        articleInfo.value = res1.data.data.article 
        let res2 = await axios.get("category/" + res1.data.data.article.category_id) 
        console.log(res2)
        if (res2.data.code == 200) {
            categoryName.value = res2.data.data.categoryName
        }
        let res3 = await axios.get("user/briefInfo/" + res1.data.data.article.user_id)
        console.log(res3)
        if (res3.data.code == 200) {
            user.value = res3.data.data
            userUrl.value = serverUrl + user.value.avatar
            if (user.value.id == user.value.loginId) {
                self.value = true
            }
        }    
        let res4 = await axios.get("collects/" + route.query.id) 
        console.log(res4)
        if (res4.data.code == 200) {
            collected.value = res4.data.data.collected
            index.value = res4.data.data.index
        }  
    }
}

const newCollect = async() => {
    let res = await axios.put("collects/new/" + route.query.id)
    console.log(res)  
    if (res.data.code == 200) {
        message.warning("已收藏", {showIcon: false})  
        loadArticle()  
    }
}

const unCollect = async() => {
    let res = await axios.delete("collects/" + index.value)
    console.log(res)  
    if (res.data.code == 200) {
        message.warning("取消收藏", {showIcon: false})  
        loadArticle()  
    }
}

const toOtherUser = () => {
    if (user.value.id == user.value.loginId) {
        router.push({
            path: "/myself",
            query: {
                id: user.value.id
            }
        })   
    } else {
        router.push({
            path: "/others",
            query: {
                id: user.value.id
            }
        })
    }
}

const toUpdate = () => {
    router.push({
        path: "/update",
        query: {
            id: articleInfo.value.id
        }
    })
}

const toDelete = async (blog) => {
    dialog.warning({
      title: '警告',
      content: '是否要删除',
      positiveText: '确定',
      negativeText: '取消',
      onPositiveClick: async () => {
            let res = await axios.delete("article/" + articleInfo.value.id)
            if(res.data.code == 200){
                message.info(res.data.msg)
                goback()
            }else{
                message.error(res.data.msg)
            }  
        },
        onNegativeClick: () => {}
    })    
}

const goback= () => {
    router.go(-1)    
}

</script>

<style lang="scss" scoped>
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
.article-content img{
    max-width: 100% !important;
}
</style>