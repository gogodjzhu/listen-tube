import type { ApiResponse } from './Types'
import API from './ApiService'

export interface User {
  UserName: string
  UserCredit: string
}

class AuthAPI {
  public async register (username: string, password: string): Promise<boolean | null> {
    try {
      return API.POST<null>('/auth/register', { username, password })
        .then(success => {
          return true
        })
    } catch (error) {
      throw new Error('Registration failed: ' + error)
    }
  }

  public async login (username: string, password: string): Promise<boolean> {
    try {
      return API.POST<string>('/auth/login', { username, password })
        .then(token => {
          localStorage.setItem('token', token)
          return true
        })
    } catch (error) {
      throw new Error('Login failed: ' + error)
    }
  }

  public async logout (): Promise<boolean> {
    try {
      return API.POST<string>('/auth/logout', {})
        .then(token => {
          return true
        })
    } catch (error) {
      throw new Error('Logout failed: ' + error)
    } finally {
      localStorage.removeItem('token')
    }
  }

  public async currentUser (): Promise<User | null> {
    try {
      if (!localStorage.getItem('token')) {
        return new Promise((resolve, reject) => {
          resolve(null)
        })
      }
      return API.GET<User>('/auth/current_user')
        .then(user => {
          return user
        })
    } catch (error) {
      throw new Error('Failed to get current user: ' + error)
    }
  }
}

const auth = new AuthAPI()
export default auth
