

<template>
	<div class="container py-5">
		<div class="row justify-content-center">
			<div class="col-md-6 col-lg-5">
				<div class="card shadow-sm">
					<div class="card-body">
						<ElementsButtonBack to="/" />
						<form @submit.prevent="onSubmit">
							<h2 class="mb-4 text-center">{{ $t('forgot.title') }}</h2>
							<div class="mb-3">
								<label for="email" class="form-label">{{ $t('forgot.email') }}</label>
								<input id="email" v-model="email" type="email" class="form-control" required>
							</div>
							<div v-if="error" class="alert alert-danger text-center py-2">{{ error }}</div>
							<div v-if="success" class="alert alert-success text-center py-2">{{ $t('forgot.success') }}</div>
							<button type="submit" class="btn btn-primary w-100" :disabled="loading">
								{{ loading ? $t('forgot.submited') + '...' : $t('forgot.submit') }}
							</button>
						</form>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
const t = useI18n().t,
	loading = ref(false),
	email = ref(''),
	error = ref<string | null>(null),
	success = ref(false),
	onSubmit = async () => {
		try {
			loading.value = true;
			const { success: successResponse, errors} = await useAuth().forgotPassword(email.value);
			loading.value = false;
			if (successResponse) {
				error.value = null;
				success.value = true;
				setTimeout(() => {
					navigateTo('/');
				}, 3000);
			} else {
				error.value = errors ? errors.join(', ') : t('forgot.error');
			}
		}catch (e) {
			console.error('Forgot Password exception:', e);
			loading.value = false;
			error.value = t('forgot.error');
		}
	};
</script>
