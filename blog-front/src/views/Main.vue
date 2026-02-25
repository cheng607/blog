<script setup lang="ts">
import AnswerBook from '../components/AnswerBook.vue';
import List from '../components/List.vue';
import Mine from '../components/Mine.vue';
import Type from '../components/Type.vue';

</script>

<template>
  <div class="main-page">
    <div class="list-container-wrapper">
      <List></List>
    </div>
    <div class="right">
      <Mine></Mine>
      <Type></Type>
      <AnswerBook></AnswerBook>
    </div>
  </div>
</template>

<style scoped>
.main-page {
  height: auto;
  display: flex;
  justify-content: space-around;
  margin: 0 150px;
  padding: 20px;
}

.list-container-wrapper {
  width: 75%;
}

.right {
  width: 25%;
  margin-left: 20px;
  display: flex;
  flex-direction: column;
  gap: 20px;
}
</style>