<template>
  <div>
    <h1>Sidebar</h1>
    <!-- a button to invoke updateSubscription -->
    <button @click="updateSubscription">Update Subscription</button>
  </div>
</template>

<script>
// import Request from '@/components/utils/Request'
import Auth from '@/components/utils/Auth'

export default {
  name: 'Sidebar',
  components: {
    Auth
  },
  data () {
    return {
      currentUser: null
    }
  },
  methods: {
    updateSubscription () {
      const currentUser = Auth.currentUser()
      console.log('Current user:', currentUser)
    }
  }
  // mounted () {
  //   this.updateSubscription()
  // }
}
</script>

<style scoped>
</style>
