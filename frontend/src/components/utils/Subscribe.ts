import API from './ApiService'

export interface Subscription {
  platform: string
  channel_name: string
  channel_thumbnail: string
  create_at: number
  update_at: number
}

class SubscribeAPI {
  public async subscribe (channelId: string): Promise<boolean> {
    try {
      return await API.POST<boolean>('/buzz/subscription/add', { channel_id: channelId })
        .then(response => {
          return true
        })
    } catch (error) {
      throw new Error('subscribe failed: ' + error)
    }
  }

  public async unsubscribe (channelId: string): Promise<boolean> {
    try {
      return await API.POST<boolean>('/buzz/subscription/delete', { channel_id: channelId })
        .then(response => {
          return true
        })
    } catch (error) {
      throw new Error('unsubscribe failed: ' + error)
    }
  }

  public async listSubscriptions (): Promise<Subscription[]> {
    try {
      return await API.GET<[Subscription]>('/buzz/subscription/list')
        .then(subscriptions => {
          return subscriptions
        })
    } catch (error) {
      throw new Error('listSubscriptions failed: ' + error)
    }
  }
}

const subscribe = new SubscribeAPI()
export default subscribe
