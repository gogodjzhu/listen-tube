<template>
  <div>
    <!-- header -->
    <div class="header bg-light py-2 px-3 w-100 d-flex justify-content-between align-items-center">
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

    <!-- stream -->
    <div class="content">
      <Stream />
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
import Stream from '@/components/stream/Stream.vue'

export default {
  name: 'Home',
  components: {
    Login,
    Sidebar,
    Stream
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
.header {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 1000;
}

.search {
  position: relative;
}

.content {
  margin-top: 60px; /* Adjust based on header height */
}
</style>
