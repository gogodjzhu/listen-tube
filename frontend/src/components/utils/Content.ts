import API from './ApiService'

export interface Content {
  platform: string
  name: string
  credit: string
  channel_name: string
  channel_credit: string
  thumbnail: string
  published_time: string
  length: string
  state: number
  create_at: number
  update_at: number
}

class ContentAPI {
  public async listContents (pageIndex: number, pageSize: number): Promise<Content[]> {
    try {
      const params = { page_index: pageIndex, page_size: pageSize }
      return await API.GET<[Content]>('/buzz/content/list', params)
        .then(contents => {
          return contents
        })
    } catch (error) {
      throw new Error('listContents failed: ' + error)
    }
  }
}

const content = new ContentAPI()
export default content
