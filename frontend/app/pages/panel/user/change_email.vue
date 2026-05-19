<template>
    <div class="container-fluid py-4">
        <div class="row justify-content-center">
            <div class="col-12">
                <div class="mb-4 text-center"> 
                    <ElementsButtonBack/>
                    <h1 class="mb-3 text-center">{{ $t('panel.user.changeEmail') }}</h1>
                </div>
            </div>
            <div class="col-12">
                <div class="card shadow-sm mb-4 p-4">
                    <div class="card-body d-flex justify-content-center">
                        <form class="col-5" @submit.prevent="handleChangeEmail">
                            <div class="mb-3">
                                <label for="newEmail" class="form-label">{{ $t('panel.user.newEmail') }}</label>
                                <input id="newEmail" v-model="newEmail" type="email" class="form-control" required>
                            </div>
                            <div class="mb-3">
                                <label for="password" class="form-label">{{ $t('panel.user.passwordConfirm') }}</label>
                                <input id="password" v-model="password" type="password" class="form-control" required>
                            </div>
                            <button type="submit" class="btn btn-primary w-100">{{ $t('panel.user.saveEmail') }}</button>
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

const newEmail = ref('')
const password = ref('')
const error = ref('')
const success = ref('')

const handleChangeEmail = async () => {
    error.value = ''
    success.value = ''
    const result = await useUser().changeEmail(newEmail.value, password.value)
    if (!result.success) {
        error.value = result.error || $t('panel.user.changeEmailError')
        return
    }
    success.value = $t('panel.user.changeEmailSuccess')
    setTimeout(() => {
        useAuth().logout()
    }, 1500)
}
</script>

<style lang="scss" scoped>
</style>
