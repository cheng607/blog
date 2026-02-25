<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router';
import List from '../components/List.vue';
import { getTagByid } from '../api/tag';

interface typeObj {
    id: number,
    name: string,
    description: string
}
const route = useRoute()
const typeData = ref<typeObj>()
onMounted(async () => {
    const typeId = route.query.typeId
    const res = await getTagByid(Number(typeId))
    typeData.value = res.data
})

</script>

<template>
    <div class="detail-page">
        <div class="cat-header">
            <h1>{{ typeData?.name }}</h1>
            <p>{{ typeData?.description }}</p>
        </div>
        <ElDivider />
        <div class="cat-content">
            <List :typeId="typeData?.id"></List>
        </div>
    </div>
</template>

<style scoped>
.detail-page {
    width: 100%;
    padding: 40px 20px;
    box-sizing: border-box;
    max-width: 1300px;
    margin: 0 auto;
}

.cat-header {
    text-align: center;
    margin-bottom: 40px;
}

.cat-header h1 {
    font-size: 2.5rem;
    margin-bottom: 15px;
    color: var(--text-color);
}

.cat-header p {
    font-size: 1.1rem;
    color: #666;
}

.cat-content {
    margin-top: 30px;
    padding: 0 20px;
}
</style>
