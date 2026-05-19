export default defineNuxtRouteMiddleware((to, from) => {
    //  to działa
    const auth = useAuth()
    if (!auth.isAuthenticated()) {
        return navigateTo('/')
    }
})