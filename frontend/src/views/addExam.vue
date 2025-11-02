<template>
  <div class="bg-gray-50 w-full">
    <div class="min-h-screen max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      <!-- Header -->
      <div class="mb-8">
        <div class="flex items-center gap-3 mb-4">
          <button
            @click="cancel"
            class="text-gray-600 hover:text-gray-900 p-2 rounded-lg hover:bg-gray-100"
          >
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
            </svg>
          </button>
          <h1 class="text-3xl font-bold text-gray-900">เพิ่มข้อสอบใหม่</h1>
        </div>
        <p class="text-gray-600">กรอกข้อมูลข้อสอบให้ครบถ้วน</p>
      </div>

      <!-- Form -->
      <div class="bg-white rounded-xl shadow-sm border border-gray-200 p-6">
       
        <div class="mb-6">
          <label class="block text-sm font-medium text-gray-700 mb-2">คำถาม</label>
          <textarea
            v-model="form.question"
            rows="3"
            class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent resize-none"
            placeholder="พิมพ์คำถาม..."
            required
          ></textarea>
        </div>

        
        <div class="mb-6">
          <label class="block text-sm font-medium text-gray-700 mb-4">ตัวเลือกคำตอบ</label>
          
          <div v-for="(answer, index) in form.answers" :key="index" class="mb-4 last:mb-0">
            <div class="flex items-start gap-3">
              <div class="flex items-center mt-1 w-8">
                <span class="text-sm font-medium text-gray-700">{{ index + 1 }}.</span>
              </div>
              <div class="flex-1">
                <input
                  v-model="answer.text"
                  type="text"
                  :placeholder="`${index + 1}. พิมพ์คำตอบ...`"
                  class="block w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  required
                >
              </div>
            </div>
          </div>
        </div>

        
        <div class="flex gap-3 justify-end pt-4 border-t border-gray-200">
          <button
            @click="cancel"
            class="px-4 py-1 text-gray-700 bg-gray-200 hover:bg-gray-300 rounded-lg font-medium transition-colors"
          >
            ยกเลิก
          </button>
          <button
            @click="save"
            :disabled="!isFormValid"
            :class="[
              'px-4 py-1 font-medium rounded-lg transition-colors',
              isFormValid
                ? 'bg-blue-600 hover:bg-blue-700 text-white'
                : 'bg-gray-300 text-gray-500 cursor-not-allowed'
            ]"
          >
            บันทึกข้อสอบ
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios'

export default {
  name: 'AddExam',
  data() {
    return {
      form: {
        question: '',
        answers: [
          { text: '' },
          { text: '' },
          { text: '' },
          { text: '' }
        ]
      }
    }
  },
  computed: {
    isFormValid() {
      return (
        this.form.question.trim() &&
        this.form.answers.every(answer => answer.text.trim())
      )
    }
  },
  methods: {
    async save() {
      if (!this.isFormValid) return

      try {
        await axios.post('http://localhost:8080/questions', {
          content: this.form.question,
          answers: this.form.answers.map(a => ({ text: a.text })) // ไม่ส่ง correct
        })
        this.$router.push('/')
      } catch (error) {
        console.error('Error saving question:', error)
        alert('เกิดข้อผิดพลาดในการบันทึก')
      }
    },
    cancel() {
      this.$router.push('/')
    }
  }
}
</script>