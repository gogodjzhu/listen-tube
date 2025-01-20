<template>
  <div id="stream" class="py-2 px-3 col-md-8 mx-auto">
    <div v-if="contents">
      <div v-for="content in contents" :key="content.name">
        <div>
          <div class="card d-flex flex-row">
            <div class="card-img-left-wrapper">
              <img :src="content.thumbNail" class="card-img-left" alt="..."/>
              <div class="media-length">{{content.length}}</div>
            </div>
            <div class="card-body">
              <h5 class="text-truncate">{{content.name}}</h5>
              <p class="text-truncate">{{content.channelName}}</p>
              <p class="text-truncate">{{content.publishedTime}}</p>
            </div>
          </div>
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
    width: calc(100% - 60px); /* Default width when sidebar is collapsed */
    margin-bottom: 1rem;
    border: none;
    border-radius: 10px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.card-img-left-wrapper {
    position: relative;
    width: 150px; /* Fixed width for the image container */
    height: auto; /* Adjust height automatically */
    overflow: hidden;
}

.card-img-left {
    width: 100%; /* Ensure the image takes the full width of the container */
    height: auto; /* Maintain the aspect ratio */
    object-fit: cover;
}

.media-length {
    font-size: 8px;
    color: white;
    background-color: rgba(69, 69, 69, 0.7); /* Use rgba for semi-transparent background */
    position: absolute;
    bottom: 20px; /* Add some space from the bottom */
    right: 7px; /* Add some space from the right */
    padding: 2px 5px; /* Add some padding for better readability */
    z-index: 1; /* Ensure it appears on top of the image */
    width: auto; /* Ensure the width adjusts to the content */
    height: auto; /* Ensure the height adjusts to the content */
    max-width: calc(100% - 10px); /* Ensure it doesn't overflow the image width */
    max-height: calc(100% - 10px); /* Ensure it doesn't overflow the image height */
    border-radius: 5px; /* Add rounded corners */
}

.card-body {
    flex: 1;
    min-width: 0; /* Ensure the text truncation works within the flex container */
}

.text-truncate {
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}
</style>
