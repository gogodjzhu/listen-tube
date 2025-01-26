<template>
  <div>
    <!-- Login components -->
    <button v-if="userInfoStore.currentUser == null" class="btn btn-primary" @click="showLoginModal = true">Login & Sign
      In</button>
    <span v-if="userInfoStore.currentUser != null">
      <span>Hello, {{ userInfoStore.name }}</span>
      <span></span>
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

<script lang="ts">
import authAPI from '../utils/Auth'
import useUserInfoStore from '../../stores/UserInfo'

export default {
  name: 'Login',
  data () {
    return {
      showLoginModal: false,
      username: '',
      password: ''
    }
  },
  setup () {
    const userInfoStore = useUserInfoStore()
    return {
      userInfoStore: userInfoStore
    }
  },
  methods: {
    handleLogin () {
      authAPI.login(this.username, this.password)
        .then((ret) => {
          if (ret) {
            this.handleUpdateCurrentUser()
            this.showLoginModal = false
            this.username = ''
            this.password = ''
          }
        })
        .catch(error => {
          console.log('Internal error: ' + error)
        })
    },
    handleLogout () {
      authAPI.logout()
        .then(() => {
          this.handleUpdateCurrentUser()
        })
        .catch(error => {
          console.log('Internal error: ' + error)
        })
    },
    handleUpdateCurrentUser () {
      authAPI.currentUser()
        .then(user => {
          this.userInfoStore.updateCurrentUser(user)
        })
        .catch(error => {
          console.log('Internal error: ' + error)
          this.userInfoStore.$reset()
        })
    },
  },
  mounted () {
    this.handleUpdateCurrentUser()
  }
}
</script>

<style scoped>
.modal {
  display: block;
  background: rgba(0, 0, 0, 0.5);
}
</style>
