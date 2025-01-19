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
      <div class="offcanvas-header">
        <h5 class="offcanvas-title" id="offcanvasExampleLabel">Offcanvas</h5>
        <button type="button" class="btn-close" data-bs-dismiss="offcanvas" aria-label="Close"></button>
      </div>
      <div class="offcanvas-body">
        <div>
          Some text as placeholder. In real life you can have the elements you have chosen. Like, text, images, lists,
          etc.
        </div>
        <div class="dropdown mt-3">
          <button class="btn btn-secondary dropdown-toggle" type="button" data-bs-toggle="dropdown">
            Dropdown button
          </button>
          <ul class="dropdown-menu">
            <li><a class="dropdown-item" href="#">Action</a></li>
            <li><a class="dropdown-item" href="#">Another action</a></li>
            <li><a class="dropdown-item" href="#">Something else here</a></li>
          </ul>
        </div>
      </div>
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
