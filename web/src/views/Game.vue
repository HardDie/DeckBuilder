<template>
  <div class="game-wrapper" v-if="game">
    <header class="game-wrapper__header game">
      <img class="game__image" :src="game.image">
      <div class="game__info">
        <div class="info__about">
          <h2>{{ game.name }}</h2>
          <span>{{ game.description }}</span>
        </div>
        <div class="info__control">
          <v-btn color="primary">Изменить</v-btn>
          <v-btn color="error">Удалить</v-btn>
        </div>
      </div>
    </header>
    <section class="decks">
      <h3 class="decks__title">Колоды</h3>
      <v-card width="300" height="500" class="d-flex align-center justify-center">
        <v-icon x-large>mdi-plus-thick</v-icon>
      </v-card>
    </section>
  </div>
</template>

<script>
import { mapGetters, mapActions } from 'vuex'
export default {
  name: 'Game',
  data () {
    return {
      game: null
    }
  },
  computed: {
    ...mapGetters('games', ['getGames'])
  },
  methods: {
    ...mapActions('games', ['fetchGames'])
  },
  async mounted () {
    if (!this.getGames.length) {
      await this.fetchGames()
    }
    this.game = this.getGames.find(game => game.name === this.$route.params?.name)
  }
}
</script>

<style lang="scss" scoped>
.game {
  display: flex;
  gap: 30px;
  &__image {
    width: 400px;
    height: 450px;
  }
  &__info {
    max-width: 600px;
  }
}
.info__about {
  display: flex;
  flex-direction: column;
  gap: 10px;
  word-break: break-all;
}
.info__control {
  margin-top: 30px;
  display: flex;
  gap: 20px;
}

.decks {
  margin-top: 50px;
  &__title {
    margin-bottom: 30px;
  }
}
</style>
