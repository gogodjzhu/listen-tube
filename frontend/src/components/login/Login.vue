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
        .then((ret) => {
          if (ret) {
            this.updateCurrentUser()
            this.showLoginModal = false
            this.username = ''
            this.password = ''
          } else {
            alert('Login failed: ' + ret)
          }
        })
        .catch(error => {
          console.log('Internal error: ' + error)
        })
    },
    handleLogout () {
      Auth.logout()
        .then(() => {
          this.updateCurrentUser()
        })
        .catch(error => {
          console.log('Internal error: ' + error)
        })
    },
    // exange token for user info
    updateCurrentUser () {
      Auth.currentUser()
        .then(user => {
          this.currentUser = user
        })
        .catch(error => {
          console.log('Internal error: ' + error)
          this.currentUser = null
        })
    },
    getCurrentUser () {
      return this.currentUser
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
