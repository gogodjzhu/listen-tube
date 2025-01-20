import Request from './Request'

class Subscription {
  constructor (platform, channelName, channelThumbnail, createAt, updateAt) {
    this.platform = platform
    this.channelName = channelName
    this.channelThumbnail = channelThumbnail
    this.createAt = createAt
    this.updateAt = updateAt
  }
}

class Subscribe {
  static async subscribe (channelId) {
    try {
      return await Request.post('/buzz/subscription/add', { channel_id: channelId })
        .then(response => {
          if (response.code !== 0) {
            throw new Error('subscribe failed: ' + response.message)
          }
          return true
        })
    } catch (error) {
      throw new Error('subscribe failed: ' + error.message)
    }
  }

  static async unsubscribe (channelId) {
    try {
      return await Request.post('/buzz/subscription/delete', { channel_id: channelId })
        .then(response => {
          if (response.code !== 0) {
            throw new Error('unsubscribe failed: ' + response.message)
          }
          return true
        })
    } catch (error) {
      throw new Error('unsubscribe failed: ' + error.message)
    }
  }

  static async listSubscriptions () {
    try {
      return await Request.get('/buzz/subscription/list')
        .then(response => {
          if (response.code !== 0) {
            return []
          }
          // map the response.subscriptions to Subscription objects
          return response.subscriptions
            .map(sub => new Subscription(
              sub.platform,
              sub.channel_name,
              sub.channel_thumbnail,
              sub.create_at,
              sub.update_at))
        })
    } catch (error) {
      if (error.response && error.response.status === 401) {
        return []
      }
      throw new Error('listSubscriptions failed: ' + error.message)
    }
  }
}

export default Subscribe
