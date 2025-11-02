<template>
    <div class="min-h-screen bg-gray-50 w-full py-12">
        <div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">

            <div>
                <h1 class="text-3xl font-bold text-gray-900 mb-14">IT08</h1>
                <p class="text-gray-600">รายการข้อสอบทั้งหมด</p>
            </div>

            <div class="mb-6 py-4">
                <button @click="goToAdd"
                    class="bg-blue-600 hover:bg-blue-700 text-white font-medium py-1 px-6 rounded-lg shadow-md transition duration-200 flex items-center gap-2">
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path>
                    </svg>
                    เพิ่มข้อสอบ
                </button>
            </div>

            <div v-if="questions.length === 0" class="text-center py-16">
                <h3 class="text-lg font-medium text-gray-900 mb-2">ยังไม่มีข้อสอบ</h3>
            </div>

            <!-- รายการข้อสอบ -->
            <div v-else class="space-y-4">
                <div v-for="question in questions" :key="question.id"
                    class="mb-6 bg-white rounded-xl shadow-sm border border-gray-200 p-6 hover:shadow-md transition-shadow">
                    <div class="flex justify-between items-start">
                        <div class="flex-1">

                            <div class="flex items-center gap-3 mb-3">
                                <div class="w-8 h-8 bg-blue-100 rounded-full flex items-center justify-center">
                                    <span class="text-blue-600 font-semibold text-sm">{{ question.number }}</span>
                                </div>
                                <h3 class="text-lg font-semibold text-gray-900">{{ question.content }}</h3>
                            </div>

                            <div class="ml-11 space-y-3">
                                <label v-for="(answer, index) in question.answers" :key="index"
                                    class="flex items-center gap-3 cursor-default">
                                    <input type="radio" :name="'q' + question.id" :value="index"
                                        class="w-4 h-4 text-blue-600" />
                                    <span class="text-sm text-gray-700 flex-1">{{ answer.text }}</span>
                                </label>
                            </div>
                        </div>

                        <button @click="deleteQuestion(question.id)"
                            class="text-red-600 hover:text-red-800 hover:bg-red-50 p-2 rounded-lg transition-colors"
                            title="ลบข้อสอบ">
                            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                    d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                            </svg>
                        </button>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import axios from 'axios'

export default {
    name: 'ListExam',
    data() {
        return {
            questions: []
        }
    },
    mounted() {
        this.fetchQuestions()
    },
    methods: {
        async fetchQuestions() {
            try {
                const response = await axios.get('http://localhost:8080/questions')
                this.questions = response.data
            } catch (error) {
                console.error('Error fetching questions:', error)
            }
        },
        goToAdd() {
            this.$router.push('/add')
        },
        async deleteQuestion(id) {
            if (confirm('คุณต้องการลบข้อสอบนี้จริงๆ ใช่ไหม?')) {
                try {
                    await axios.delete(`http://localhost:8080/questions/${id}`)
                    this.questions = this.questions.filter(q => q.id !== id)
                } catch (error) {
                    console.error('Error deleting question:', error)
                    alert('เกิดข้อผิดพลาดในการลบ')
                }
            }
        }
    }
}
</script>