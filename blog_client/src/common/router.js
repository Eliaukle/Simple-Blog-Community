import { createRouter, createWebHashHistory } from "vue-router";

let routes = [
    { path: "/login", component: () => import("../views/Login.vue") },
    { path: "/register", component: () => import("../views/Register.vue") },
    { path: "/", component: () => import("../views/MainFrame.vue") },
    { path: "/publish", component: () => import("../views/Publish.vue") },
    { path:"/myself", component: () => import("../views/Myself.vue") },
    { path:"/others", component: () => import("../views/Others.vue") },
    { path:"/detail", component: () => import("../views/Detail.vue") },
    { path:"/update", component: () => import("../views/Update.vue") },
]

const router = createRouter({
    history: createWebHashHistory(),
    routes,
});

export { router, routes }   // routes暂时用不上