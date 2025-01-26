import { defineStore } from 'pinia'
import type { User } from '@/components/utils/Auth'

interface UserInfo {
  currentUser: User | null
}

const useUserInfo = defineStore('UserInfo', {
  state: (): UserInfo => ({
    currentUser: null,
  }),
  getters: {
    name: (state) => state.currentUser?.UserName,
  },
  actions: {
    updateCurrentUser (user: User | null) {
      this.currentUser = user
    },
  },
})

export default useUserInfo