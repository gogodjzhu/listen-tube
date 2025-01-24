import Request from './Request.ts'
import axios from 'axios'

export interface User {
  UserName: string
  UserCredit: string
}

class Auth {
  /**
   * Register a new user with the given username and password
   * @param {string} username The username of the user
   * @param {string} password The password of the user
   * @returns {Promise<boolean>}
   */
  static async register (username: string, password: string) {
    try {
      return await Request.post('/auth/register', { username, password })
        .then(response => {
          return true
        })
    } catch (error) {
      throw new Error('Registration failed: ' + error)
    }
  }

  /**
   * Login the user with the given username and password
   * @param {string} username The username of the user
   * @param {string} password The password of the user
   * @returns {Promise<boolean>}
   */
  static async login (username: string, password: string) {
    try {
      return await Request.post('/auth/login', { username, password })
        .then(resp => {
          if (resp.status !== 200) {
            return false
          }
          // get the token from the response and store it in the local storage
          localStorage.setItem('token', resp.data.token)
          return true
        })
    } catch (error) {
      throw new Error('Login failed: ' + error)
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
      throw new Error('Logout failed: ' + error)
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

      return Request.get<User>('/auth/current_user')
        .then(response => {
          console.log('Response:', response)
          return response
          // return User(response.username, response.credit)
        })
    } catch (error) {
      // if (error.response && error.response.status === 401) {
      //   localStorage.removeItem('token')
      //   return null
      // }
      throw new Error('Failed to get current user: ' + error)
    }
  }
}

export default Auth
