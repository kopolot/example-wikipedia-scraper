<template>
    <div class="container-fluid py-4">
        <div class="row justify-content-center">
            <div class="col-12">
                <div class="mb-4 text-center"> 
                    <ElementsButtonBack/>
                    <h1 class="mb-3 text-center">{{ $t('panel.user.changePassword') }}</h1>
                </div>
            </div>
            <div class="col-12">
                <div class="card shadow-sm mb-4 p-4">
                    <div class="card-body d-flex justify-content-center">
                        <form class="col-5" @submit.prevent="handleChangePassword">
                            <div class="mb-3">
                                <label for="currentPassword" class="form-label">{{ $t('panel.user.currentPassword') }}</label>
                                <input id="currentPassword" v-model="currentPassword" type="password" class="form-control" required>
                            </div>
                            <div class="mb-3">
                                <label for="newPassword" class="form-label">{{ $t('panel.user.newPassword') }}</label>
                                <input id="newPassword" v-model="newPassword" type="password" class="form-control" required>
                            </div>
                            <div class="mb-3">
                                <label for="confirmPassword" class="form-label">{{ $t('panel.user.confirmNewPassword') }}</label>
                                <input id="confirmPassword" v-model="confirmPassword" type="password" class="form-control" required>
                            </div>
                            <button type="submit" class="btn btn-primary w-100">{{ $t('panel.user.savePassword') }}</button>
                            <div v-if="error" class="alert alert-danger mt-3">{{ error }}</div>
                            <div v-if="success" class="alert alert-success mt-3">{{ success }}</div>
                        </form>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script lang="ts" setup>
definePageMeta({
    middleware: 'auth',
    layout: 'panel'
})

const currentPassword = ref('')
const newPassword = ref('')
const confirmPassword = ref('')
const error = ref('')
const success = ref('')

const handleChangePassword = async () => {
    error.value = ''
    success.value = ''
    if (newPassword.value !== confirmPassword.value) {
        error.value = $t('panel.user.passwordsDoNotMatch')
        return
    }
    const result = await useUser().changePassword(currentPassword.value, newPassword.value)
    if (!result.success) {
        error.value = result.error || $t('panel.user.changePasswordError')
        return
    }
    success.value = $t('panel.user.changePasswordSuccess')
    setTimeout(() => {
        useAuth().logout()
    }, 1500)
}
</script>

<style lang="scss" scoped>
</style>