export default {
  mounted(el: HTMLElement, binding: any) {
    el.clickOutsideEvent = function (event: Event) {
      if (!(el == event.target || el.contains(event.target))) {
        binding.value()
      }
    }
    document.body.addEventListener('click', el.clickOutsideEvent)
  },
  unmounted(el: HTMLElement) {
    document.body.removeEventListener('click', el.clickOutsideEvent)
  },
}
