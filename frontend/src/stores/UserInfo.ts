import { defineStore } from 'pinia'
import type { User } from '@/components/utils/Auth'
import type { Subscription } from '@/components/utils/Subscribe'
import type { Content } from '@/components/utils/Content'

interface UserInfo {
  currentUser: User | null
  subscriptions: Subscription[]
  contents: Content[]
}

const useUserInfo = defineStore('UserInfo', {
  state: (): UserInfo => ({
    currentUser: null,
    subscriptions: [],
    contents: [],
  }),
  getters: {
    name: (state) => state.currentUser?.UserName,
  },
  actions: {
    updateCurrentUser (user: User | null) {
      this.currentUser = user
    },
    updateSubscriptions (subscriptions: Subscription[]) {
      this.subscriptions = subscriptions
    },
    appendContents (content: Content[]) {
      this.contents.push(...content)
    },
  },
})

export default useUserInfo