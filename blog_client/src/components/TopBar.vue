<template>
    <div class="container">
        <div class="topbar">
            <div class="bigtitle" @click="toMain">首页</div>
            <n-dropdown v-if="login" trigger="hover" :options="options" @select="handleSelect">
                <n-avatar @click="toHome" round size="medium" :src=user.avatarUrl style="position: absolute; right: 190px; top: 8px; cursor: pointer;"/>
            </n-dropdown>
            <div v-if="!login" class="smalltitle" @click="toLogin">登录/注册</div>
            <div style="position: absolute; right: 50px; top: 8px">
                <n-button round color="#7B3DE0" @click="toPublish">发布文章</n-button>
            </div>      
        </div>
    </div>
</template>

<script setup>
import {ref,reactive,inject, onMounted} from 'vue'

import {useRouter, useRoute} from 'vue-router'
const router = useRouter()
const route = useRoute()

const serverUrl = inject("serverUrl")
const axios = inject("axios")
const message = inject("message")

const options = reactive([{label: "退出登录", key: "login"}])
const login = ref(false)
const user = reactive({
    avatarUrl: "",
    id: 0
})

onMounted(() => {
    loadAvatar()
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

const toMain = () => {
    router.push("/")    
}

const toLogin = () => {
    router.push("/login")    
}

const toHome = () => {
    router.push({
        path: "/myself",
        query: {
            id: user.id
        }
    })    
}

const toPublish = () => {
    if (login.value == false) {
        message.warning("请先登录")
    } else {
        router.push("/publish")   
    }  
}

const handleSelect = (key) => {
    router.push("/" + String(key)) 
}

</script>

<style lang="scss" scoped>
.container {
    .topbar {
        position: sticky;
        top: 0;
        height: 50px;
        background: white;
        box-shadow: 0px 1px 5px #D3D4D8;
        .bigtitle {
            position: absolute;
            font-size: 20px;
            left: 50px;
            line-height: 50px;
            color: #7B3DE0;
            cursor: pointer;
        }
        .smalltitle {
            position: absolute;
            font-size: 16px;
            right: 175px;
            line-height: 50px;
            color: #7B3DE0;
            cursor: pointer;            
        }
    }
}
</style>