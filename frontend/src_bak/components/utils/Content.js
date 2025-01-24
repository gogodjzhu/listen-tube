import Request from './Request'

class Content {
  constructor (platform, name, credit, channelName, channelCredit, thumbNail, publishedTime, length, state, createAt, updateAt) {
    this.platform = platform
    this.name = name
    this.credit = credit
    this.channelName = channelName
    this.channelCredit = channelCredit
    this.thumbNail = thumbNail
    this.publishedTime = publishedTime
    this.length = length
    this.state = state
    this.createAt = createAt
    this.updateAt = updateAt
  }
}

class ContentAPI {
  static async listContents (pageIndex, pageSize) {
    try {
      return await Request.get('/buzz/content/list', { params: { page_index: pageIndex, page_size: pageSize } })
        .then(response => {
          if (response.code !== 0) {
            return []
          }
          // map the response.contents to Content objects
          return response.contents
            .map(content => new Content(
              content.platform,
              content.name,
              content.credit,
              content.channel_name,
              content.channel_credit,
              content.thumbnail,
              content.published_time,
              content.length,
              content.state,
              content.create_at,
              content.update_at))
        })
    } catch (error) {
      if (error.response && error.response.status === 401) {
        return []
      }
      throw new Error('listContents failed: ' + error.message)
    }
  }
}

export default ContentAPI
