import { createI18n } from 'vue-i18n'

const i18n = createI18n({
  locale: 'en', // set locale
  fallbackLocale: 'en', // set fallback locale
  messages: {}, // initially empty
})

export async function loadLocaleMessages(locale: string, page: string) {
    const messages = await import(`./locales/${locale}/${page}.json`)
    i18n.global.setLocaleMessage(locale, {
        ...i18n.global.getLocaleMessage(locale),
        ...messages.default,
    })
    i18n.global.locale = locale
}

export default i18n
