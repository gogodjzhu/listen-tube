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
import Auth from '@/components/utils/Auth'

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
      Auth.login(this.username, this.password)
        .then(() => {
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
      Auth.logout()
      this.updateCurrentUser()
    },
    getCurrentUser () {
      return this.currentUser
    },
    // exange token for user info
    updateCurrentUser () {
      if (localStorage.getItem('token') == null) {
        this.currentUser = null
        return
      }
      Auth.currentUser()
        .then(response => {
          this.currentUser = {
            username: response.UserName,
            userCredit: response.UserCredit
          }
        })
        .catch(error => {
          alert('Get current user failed: ' + error.message)
          this.currentUser = null
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
