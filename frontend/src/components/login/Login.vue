<template>
  <div>
    <!-- Login components -->
    <button v-if="currentUser == null" class="btn btn-primary" @click="showLoginModal = true">Login & Sign In</button>
    <span v-if="currentUser != null">
      <span>Hello, {{ currentUser.username }}</span>
      <button class="btn btn-secondary" @click="handleLogout">Logout</button>
    </span>

    <!-- Login Modal -->
    <div v-if="showLoginModal" class="modal" tabindex="-1" role="dialog">
      <div class="modal-dialog" role="document">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">Login</h5>
          </div>
          <div class="modal-body">
            <input v-model="username" type="text" class="form-control mb-2" placeholder="Username">
            <input v-model="password" type="password" class="form-control" placeholder="Password">
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-primary" @click="handleLogin">Login</button>
            <button type="button" class="btn btn-secondary" @click="showLoginModal = false">Close</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import request from '@/components/utils/request'

export default {
  name: 'Login',
  data () {
    return {
      showLoginModal: false,
      currentUser: null,
      username: '',
      password: ''
    }
  },
  methods: {
    handleLogin () {
      return request({
        url: '/auth/login',
        method: 'post',
        data: {
          username: this.username,
          password: this.password
        }
      })
        .then(response => {
          localStorage.setItem('token', response.token)
          this.updateCurrentUser()
        })
        .catch(error => {
          alert('Login failed: ' + error.message)
        })
        .finally(() => {
          this.showLoginModal = false
          this.username = ''
          this.password = ''
        })
    },
    handleLogout () {
      this.currentUser = null
      localStorage.removeItem('token')
    },
    getCurrentUser () {
      return this.currentUser
    },
    updateCurrentUser () {
      if (localStorage.getItem('token') == null) {
        return
      }
      return request({
        url: '/auth/current_user',
        method: 'get'
      })
        .then(response => {
          this.currentUser = {
            username: response.UserName,
            userCredit: response.UserCredit
          }
          return this.currentUser
        })
        .catch(error => {
          alert('Get current user failed: ' + error.message)
        })
    }
  },
  mounted () {
    this.updateCurrentUser()
  }
}
</script>

<style scoped>
.modal {
  display: block;
  background: rgba(0, 0, 0, 0.5);
}
</style>
