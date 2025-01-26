<template>
  <div>
    <div class="subscription-title">Subscription</div>
    <ul v-if="buzzInfo.subscriptions" class="p-1">
      <li v-for="subscription in buzzInfo.subscriptions" :key="subscription.channel_name" class="subscription-item">
        <img :src="subscription.channel_thumbnail" alt="thumbnail" class="thumbnail">
        <span>{{ subscription.channel_name }}</span>
      </li>
    </ul>
    <div v-else>
      Loading subscriptions...
    </div>
  </div>
</template>

<script lang="ts">
import subscribeAPI from '../utils/Subscribe'
import useUserInfoStore from '../../stores/UserInfo'
import useBuzzInfoStore from '../../stores/BuzzInfo'

export default {
  name: 'Sidebar',
  data () {
    return {
    }
  },
  setup () {
    const buzzInfo = useBuzzInfoStore()
    const userInfo = useUserInfoStore()

    return {
      buzzInfo: buzzInfo,
      userInfo: userInfo
    }
  },
  methods: {
    handleUpdateSubscription () {
      subscribeAPI.listSubscriptions()
        .then(subscriptions => {
          this.buzzInfo.updateSubscriptions(subscriptions)
        })
        .catch(error => {
          console.log('Internal error: ' + error)
        })
    }
  },
  mounted () {
    this.userInfo.$subscribe((mutation, state) => {
      this.buzzInfo.clearSubscriptions()
      this.handleUpdateSubscription()
    })
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
