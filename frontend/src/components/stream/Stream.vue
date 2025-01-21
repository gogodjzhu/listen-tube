<template>
  <div id="stream" class="scroll-container" @scroll="handleScroll">
    <div v-if="contents.length" class="col-sm-12 col-md-10 col-lg-8 mx-auto">
      <div v-for="content in contents" :key="content.name" class="card d-flex flex-row" @click="playContent(content)">
        <div class="card-img-left-wrapper">
          <img :src="content.thumbNail" class="card-img-left" alt="..." />
        </div>
        <div class="card-body">
          <h5 class="text-truncate">{{ content.name }}</h5>
          <p class="text-truncate">{{ content.channelName }}</p>
          <p class="text-truncate">{{ content.length }} · {{ content.publishedTime }}</p>
        </div>
      </div>
      <div v-if="noMoreContents" class="no-more-contents">
        No more contents to load.
      </div>
    </div>
  </div>
</template>

<script>
import ContentAPI from '@/components/utils/Content'
import 'aplayer/dist/APlayer.min.css'
import APlayer from 'aplayer'

export default {
  name: 'Sidebar',
  components: {
    ContentAPI
  },
  data () {
    return {
      contents: [],
      pageIndex: 1,
      pageSize: 10,
      loading: false,
      noMoreContents: false
    }
  },
  methods: {
    updateContents () {
      if (this.loading || this.noMoreContents) return
      this.loading = true
      ContentAPI.listContents(this.pageIndex, this.pageSize)
        .then(newContents => {
          if (newContents.length === 0) {
            this.noMoreContents = true
          } else {
            this.contents = [...this.contents, ...newContents]
            this.pageIndex++
          }
          this.loading = false
        })
        .catch(error => {
          alert('Internal error: ' + error)
          this.loading = false
        })
    },
    handleScroll (event) {
      const bottom = event.target.scrollHeight - event.target.scrollTop === event.target.clientHeight
      if (bottom) {
        this.updateContents()
      }
    },
    playContent (content) {
      const ap = new APlayer({
        container: document.getElementById('aplayer'),
        audio: {
          name: content.name,
          artist: content.channelName,
          url: `http://172.17.0.2:8080/openapi/content/stream/` + content.credit,
          cover: content.thumbNail
        }
      })
      ap.play()
    }
  },
  mounted () {
    this.updateContents()
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
  height: calc(100vh - 60px); /* Adjust based on header and bottom height */
  overflow-y: auto;
}

.no-more-contents {
  text-align: center;
  padding: 1rem;
  color: #888;
}
</style>
