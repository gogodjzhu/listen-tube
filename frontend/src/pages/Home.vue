<template>
  <div class="bg-light py-2 px-3 w-100">

    <!-- header -->
    <div class="d-flex justify-content-between align-items-center">
      <div class="d-flex w-50">
        <!-- sidebar toggle -->
        <button class="btn" data-bs-toggle="offcanvas" data-bs-target="#offcanvasExample">
          <i class="fa fa-bars fa-lg"></i>
        </button>
        <!-- search -->
        <div class="search w-100">
          <input type="text" class="form-control" placeholder="Search..." v-model="form.search" @input="search">
        </div>
      </div>
      <!-- user info -->
      <Login ref="Login" />
    </div>

    <!-- sidebar offcanvas -->
    <div class="offcanvas offcanvas-start" tabindex="-1" id="offcanvasExample" aria-labelledby="offcanvasExampleLabel">
      <div class="offcanvas-body">
        <Sidebar />
      </div>
    </div>
  </div>
</template>

<script>
import Login from '@/components/login/Login.vue'
import Sidebar from '@/components/sidebar/Sidebar.vue'

export default {
  name: 'Home',
  components: {
    Login,
    Sidebar
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
