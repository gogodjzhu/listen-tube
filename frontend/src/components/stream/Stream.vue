<template>
  <div id="stream" class="col-sm-12 col-md-10 col-lg-8 mx-auto">
    <div v-if="contents">
      <div v-for="content in contents" :key="content.name" class="card d-flex flex-row">
        <div class="card-img-left-wrapper">
          <img :src="content.thumbNail" class="card-img-left" alt="..." />
        </div>
        <div class="card-body">
          <h5 class="text-truncate">{{ content.name }}</h5>
          <p class="text-truncate">{{ content.channelName }}</p>
          <p class="text-truncate">{{ content.length }} · {{ content.publishedTime }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import ContentAPI from '@/components/utils/Content'

export default {
  name: 'Sidebar',
  components: {
    ContentAPI
  },
  data () {
    return {
      contents: null
    }
  },
  methods: {
    updateContents () {
      ContentAPI.listContents(1, 10)
        .then(contents => {
          this.contents = contents
        })
        .catch(error => {
          alert('Internal error: ' + error)
        })
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
</style>
