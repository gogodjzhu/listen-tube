<template>
  <div class="bg-light py-2 px-3 w-100">
    <div class="d-flex justify-content-between align-items-center">
      <!-- search -->
      <div class="search w-50">
        <input type="text" class="form-control" placeholder="Search..." v-model="form.search" @input="search">
      </div>
      <!-- user info -->
      <Login ref="Login" />
    </div>
  </div>
</template>

<script>
import Login from '@/components/login/Login.vue'

export default {
  name: 'Home',
  components: {
    Login
  },
  data () {
    return {
      form: {
        search: ''
      }
    }
  },
  methods: {
    search () {
      const currentUser = this.$refs.Login.getCurrentUser()
      console.log('Search:', this.form.search, 'User:', currentUser)
    }
  }
}
</script>

<style scoped>
.search {
  position: relative;
}
</style>
