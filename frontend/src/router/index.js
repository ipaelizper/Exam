import { createRouter, createWebHistory } from 'vue-router';
import ListExam from '../views/listExam.vue';
import AddExam from '../views/addExam.vue';

const routes = [
  { path: '/', name: 'IT08-1', component: ListExam },
  { path: '/add', name: 'IT08-2', component: AddExam },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;