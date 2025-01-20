import Request from './Request'

class User {
  constructor (userName, userCredit) {
    this.userName = userName
    this.userCredit = userCredit
  }
}

class Auth {
  /**
   * Register a new user with the given username and password
   * @param {string} username The username of the user
   * @param {string} password The password of the user
   * @returns {Promise<boolean>}
   */
  static async register (username, password) {
    try {
      return await Request.post('/auth/register', { username, password })
        .then(response => {
          if (response.status !== 200) {
            return false
          }
          return true
        })
    } catch (error) {
      throw new Error('Registration failed: ' + error.message)
    }
  }

  /**
   * Login the user with the given username and password
   * @param {string} username The username of the user
   * @param {string} password The password of the user
   * @returns {Promise<boolean>}
   */
  static async login (username, password) {
    try {
      return await Request.post('/auth/login', { username, password })
        .then(response => {
          if (response.code !== 200) {
            return false
          }
          localStorage.setItem('token', response.token)
          return true
        })
    } catch (error) {
      if (error.response && error.response.status === 401) {
        return false
      }
      throw new Error('Login failed: ' + error.message)
    }
  }

  /**
   * Logout the user by removing the token from the local storage
   * @returns {Promise<null>}
   */
  static async logout () {
    try {
      return new Promise((resolve, reject) => {
        localStorage.removeItem('token')
        resolve(null)
      })
    } catch (error) {
      throw new Error('Logout failed: ' + error.message)
    }
  }

  /**
   * Get the current user from the server with the token
   * @returns {Promise<User>}
   */
  static async currentUser () {
    try {
      if (!localStorage.getItem('token')) {
        return new Promise((resolve, reject) => {
          resolve(null)
        })
      }
      return Request.get('/auth/current_user')
        .then(response => {
          return new User(response.username, response.credit)
        })
    } catch (error) {
      if (error.response && error.response.status === 401) {
        localStorage.removeItem('token')
        return null
      }
      throw new Error('Failed to get current user: ' + error.message)
    }
  }
}

export default Auth
