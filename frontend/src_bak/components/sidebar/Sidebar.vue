<template>
  <div>
    <div class="subscription-title">Subscription</div>
    <ul v-if="subscriptions" class="p-1">
      <li v-for="subscription in subscriptions" :key="subscription.channelName" class="subscription-item">
        <img :src="subscription.channelThumbnail" alt="thumbnail" class="thumbnail">
        <span>{{ subscription.channelName }}</span>
      </li>
    </ul>
    <div v-else>
      Loading subscriptions...
    </div>
  </div>
</template>

<script>
import Subscribe from '@/components/utils/Subscribe'

export default {
  name: 'Sidebar',
  components: {
    Subscribe
  },
  data () {
    return {
      subscriptions: null
    }
  },
  methods: {
    updateSubscription () {
      Subscribe.listSubscriptions()
        .then(subscriptions => {
          this.subscriptions = subscriptions
        })
        .catch(error => {
          alert('Internal error: ' + error)
        })
    }
  },
  mounted () {
    this.updateSubscription()
  }
}
</script>

<style scoped>
.subscription-title {
  text-align: left;
  font-size: 1.3rem;
  margin-bottom: 20px;
}

.subscription-item {
  display: flex;
  align-items: center;
  margin-bottom: 10px;
}

.thumbnail {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  margin-right: 10px;
}
</style>
