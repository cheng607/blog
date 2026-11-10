<script setup lang="ts">

import AnswerBook from '../components/AnswerBook.vue';
import List from '../components/List.vue';
import Mine from '../components/Mine.vue';
import Type from '../components/Type.vue';

</script>

<template>
  <div class="main-page">
    <Type class="page-type" />
    <div class="list-container-wrapper">
      <List></List>
    </div>
    <div class="right">
      <Mine class="page-mine" />
      <AnswerBook class="page-answer" />
    </div>
  </div>
</template>

<style scoped>
.main-page {
  height: auto;
  display: grid;
  grid-template-columns: 1fr 280px;
  grid-template-areas:
    "list mine"
    "list type"
    "list answer";
  justify-content: center;
  max-width: 1400px;
  margin: 0 auto;
  padding: 20px 24px;
  gap: 20px;
  box-sizing: border-box;
}

.list-container-wrapper {
  grid-area: list;
  min-width: 0;
}

.page-type {
  grid-area: type;
  align-self: start;
}

.right {
  display: contents;
}

.page-mine {
  grid-area: mine;
  align-self: start;
}

.page-answer {
  grid-area: answer;
  align-self: start;
}

@media (max-width: 992px) {
  .main-page {
    grid-template-columns: 1fr;
    grid-template-areas:
      "type"
      "list"
      "mine"
      "answer";
    padding: 16px;
    margin: 0;
    gap: 12px;
  }

  .page-type {
    margin: 0 -4px;
  }
}
</style>