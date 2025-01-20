import Request from './Request'

// ...existing code...

class Auth {
  static async register (username, password) {
    try {
      const response = await Request.post('/auth/register', { username, password })
      return response
    } catch (error) {
      throw new Error('Registration failed: ' + error.message)
    }
  }

  static async login (username, password) {
    try {
      const response = await Request.post('/auth/login', { username, password })
      localStorage.setItem('token', response.token)
      return response
    } catch (error) {
      throw new Error('Login failed: ' + error.message)
    }
  }

  static async logout () {
    try {
      localStorage.removeItem('token')
    } catch (error) {
      throw new Error('Logout failed: ' + error.message)
    }
  }

  static async currentUser () {
    try {
      const response = await Request.get('/auth/current_user')
      return response
    } catch (error) {
      throw new Error('Failed to get current user: ' + error.message)
    }
  }
}

export default Auth
