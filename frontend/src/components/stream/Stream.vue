<template>
  <div id="stream" class="scroll-container" @scroll="handleScroll">
    <div v-if="buzzInfo.contents.length" class="col-sm-12 col-md-10 col-lg-8 mx-auto">
      <div v-for="content in buzzInfo.contents" :key="content.name" class="card d-flex flex-row"
        @click="handlePlayContent(content)">
        <div class="card-img-left-wrapper">
          <img :src="content.thumbnail" class="card-img-left" alt="..." />
        </div>
        <div class="card-body">
          <h5 class="text-truncate">{{ content.name }}</h5>
          <p class="text-truncate">{{ content.channel_name }}</p>
          <p class="text-truncate">{{ content.length }} · {{ content.published_time }}</p>
        </div>
      </div>
      <div v-if="noMoreContents" class="no-more-contents">
        No more contents to load.
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import contentAPI from '../utils/Content'
import type { Content } from '../utils/Content'
import useUserInfoStore from '../../stores/UserInfo'
import useBuzzInfoStore from '../../stores/BuzzInfo'

import 'aplayer/dist/APlayer.min.css'
// @ts-ignore: aplayer
import APlayer from 'aplayer'

export default {
  name: 'Stream',
  data () {
    return {
      pageIndex: 1,
      pageSize: 10,
      loading: false,
      noMoreContents: false
    }
  },
  setup () {
    const userInfo = useUserInfoStore()
    const buzzInfo = useBuzzInfoStore()
    return {
      userInfo: userInfo,
      buzzInfo: buzzInfo
    }
  },
  methods: {
    handleUpdateContents () {
      if (this.loading || this.noMoreContents) return
      this.loading = true
      contentAPI.listContents(this.pageIndex, this.pageSize)
        .then(newContents => {
          if (newContents.length === 0) {
            this.noMoreContents = true
          } else {
            this.buzzInfo.appendContents(newContents)
            this.pageIndex++
          }
          this.loading = false
        })
        .catch(error => {
          console.log('Internal error: ' + error)
          this.loading = false
        })
    },
    // @ts-ignore: event type
    handleScroll (event) {
      const bottom = event.target.scrollHeight - event.target.scrollTop === event.target.clientHeight
      if (bottom) {
        this.handleUpdateContents()
      }
    },
    handlePlayContent (content: Content) {
      const ap = new APlayer({
        container: document.getElementById('aplayer'),
        audio: {
          name: content.name,
          artist: content.channel_name,
          // @ts-ignore: import.meta
          url: `${import.meta.env.VITE_API_BASE_URL}/openapi/content/stream/` + content.credit,
          cover: content.thumbnail
        }
      })
      ap.play()
    }
  },
  mounted () {
    this.userInfo.$subscribe((mutation, state) => {
      this.pageIndex = 1
      this.loading = false
      this.noMoreContents = false
      
      this.buzzInfo.clearContents()
      this.handleUpdateContents()
    })
  }
}
</script>

<style scoped>
.card {
  display: flex;
  flex-direction: row;
  width: auto;
  margin: 0.5rem;
  padding: 0.5rem;
  border: none;
  border-radius: 0.5rem;
  box-shadow: 0rem 0.2rem 0.4rem rgba(0, 0, 0, 0.1);
}

.card-img-left-wrapper {
  width: 10rem;
  overflow: hidden;
}

.card-img-left {
  width: 100%;
  height: auto;
  object-fit: cover;
}

.card-body {
  flex: 1;
  min-width: 0;
  text-align: left;
  padding: 0rem 1rem;
}

.text-truncate {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.scroll-container {
  height: calc(100vh - 60px);
  /* Adjust based on header and bottom height */
  overflow-y: auto;
}

.no-more-contents {
  text-align: center;
  padding: 1rem;
  color: #888;
}
</style>
