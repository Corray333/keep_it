<!-- VkLoginButton.vue -->
<template>
  <div ref="vkContainer" class="vk-login-container"></div>
</template>

<script lang="ts" setup>
import { onMounted, ref } from 'vue'

const emit = defineEmits(['auth'])

const vkContainer = ref<HTMLElement | null>(null)

interface vkAuthData{
  id_token: string;
}

const vkidOnSuccess = (data: vkAuthData) => {
  console.log('Login success:', data)
  emit('auth', data.id_token)
}

const vkidOnError = (error: any) => {
  console.error('Login error:', error)
  // Здесь обработайте ошибку авторизации
}

const initVkLogin = () => {
  if ('VKIDSDK' in window && vkContainer.value) {
    const VKID = window.VKIDSDK

    VKID.Config.init({
      app: 53341208,
      redirectUrl: 'https://keep-it.xyz/auth/vk',
      responseMode: VKID.ConfigResponseMode.Callback,
      source: VKID.ConfigSource.LOWCODE,
      scope: '' // Добавьте необходимые права доступа
    })

    const oneTap = new VKID.OneTap()

    oneTap.render({
      container: vkContainer.value,
      showAlternativeLogin: true,
      styles: {
        borderRadius: 50,
        height: 32
      }
    })
      .on(VKID.WidgetEvents.ERROR, vkidOnError)
      .on(VKID.OneTapInternalEvents.LOGIN_SUCCESS, (payload: any) => {
        const { code, device_id: deviceId } = payload
        
        VKID.Auth.exchangeCode(code, deviceId)
          .then(vkidOnSuccess)
          .catch(vkidOnError)
      })
  }
}

onMounted(() => {
  // Загружаем SDK
  const script = document.createElement('script')
  script.src = 'https://unpkg.com/@vkid/sdk@<3.0.0/dist-sdk/umd/index.js'
  script.async = true
  script.onload = initVkLogin
  document.body.appendChild(script)
})

</script>

<style scoped>
.vk-login-container {
  display: inline-block;
}
</style>