<template>
    <div>
        <div style="position: fixed; top: 0; height: 50px; width: 100%; z-index: 999; ">
            <div class="card">
                <n-popselect @update:value="searchByCategory" v-model:value="selectedCategory" :options="categoryOptions" trigger="click">
                    <n-button text style="position:absolute; left: 50px; top: 22px; font-size: 18px;">{{categoryName}}</n-button>
                </n-popselect>

                <n-input v-model:value="pageInfo.keyword" round placeholder="请输入关键字" style="position:absolute; left: 125px; top: 15px; width: 1000px; background-color: #F3F0F9;" />

                <n-button @click="loadArticles(0)" round color="#7B3DE0" style="position:absolute; left: 1150px; top: 15px;">
                    <template #icon>
                        <n-icon>
                            <search />
                        </n-icon>
                    </template>
                搜索
                </n-button>
            </div>            
        </div>

        <div style="position: fixed; top: 0; height: 50px; width: 100%; z-index: 999; "><TopBar></TopBar></div>

        <div class="tabs">
            <n-card>
                <div v-for="(article,index) in articleList" style="margin-bottom:15px">
                    <n-card v-if="article.head_image" @click="toDetail(article)" style="cursor: pointer;" hoverable >
                        <n-image width="200" :src=serverUrl+article.head_image style="float: left" />
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
                <n-pagination @update:page="loadArticles" v-model:page="pageInfo.pageNum" :page-count="pageInfo.pageCount" />
            </n-card>
        </div>
    </div>

</template>

<script setup>
import TopBar from '../components/TopBar.vue'
import {ref,reactive,inject,onMounted,computed} from 'vue'
import {Search} from '@vicons/ionicons5'

import {useRouter, useRoute} from 'vue-router'
const router = useRouter()
const route = useRoute()

const serverUrl = inject("serverUrl")
const axios = inject("axios")
const message = inject("message")

const selectedCategory = ref(0)
const categoryOptions = ref([])
const articleList = ref([])
const pageInfo = reactive({
  pageNum:1,
  pageSize:5,
  pageCount:0,
  count:0,
  keyword:"",
  categoryId:0
})

onMounted(()=>{
    loadArticles()
    loadCategories()
})

const loadArticles = async(pageNum = 0) =>{
    if (pageNum != 0){
        pageInfo.pageNum = pageNum;
    }
    let res = await axios.post(`/article/list?keyword=${pageInfo.keyword}&pageNum=${pageInfo.pageNum}&pageSize=${pageInfo.pageSize}&categoryId=${pageInfo.categoryId}`)
    console.log(res)
    if (res.data.code == 200) {
        articleList.value = res.data.data.article
    }
    pageInfo.count = res.data.data.count;
    pageInfo.pageCount = parseInt(pageInfo.count / pageInfo.pageSize) + (pageInfo.count % pageInfo.pageSize > 0 ? 1 : 0)
    console.log(pageInfo.pageNum, pageInfo.pageCount, pageInfo.count)
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

const categoryName = computed(() => {
    let selectedOption = categoryOptions.value.find((option) => {return option.value == selectedCategory.value})
    console.log(selectedOption)
    return selectedOption ? selectedOption.label : ""
})

const searchByCategory = (categoryId) => {
    pageInfo.categoryId = categoryId
    pageInfo.pageNum = 1
    loadArticles()
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
    top: 50px;
    left: 0;
    right: 0;
    margin: auto;
    height: 60px;
    background: white;  
    box-shadow: 0px 1px 3px #D3D4D8; 
    border-radius: 5px;
}
.tabs {
    position: absolute;
    top: 150px;
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