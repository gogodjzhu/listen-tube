import { defineStore } from 'pinia'
import type { Content } from '@/components/utils/Content'
import type { Subscription } from '@/components/utils/Subscribe'

interface BuzzInfo {
  subscriptions: Subscription[]
  contents: Content[]
}

const useBuzzInfo = defineStore('BuzzInfo', {
  state: (): BuzzInfo => ({
    subscriptions: [],
    contents: [],
  }),
  getters: {
  },
  actions: {
    clearSubscriptions () {
      this.subscriptions = []
    },
    updateSubscriptions (subscriptions: Subscription[]) {
      this.subscriptions = subscriptions
    },
    appendSubscriptions (subscriptions: Subscription[]) {
      this.subscriptions.push(...subscriptions)
    },
    clearContents () {
      this.contents = []
    },
    updateContents (contents: Content[]) {
      this.contents = contents
    },
    appendContents (content: Content[]) {
      this.contents.push(...content)
    },
  },
})

export default useBuzzInfo