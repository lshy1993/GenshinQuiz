const db = require('../db');

class Quiz {
  static async findAll(filters = {}) {
    let query = db('quizzes')
      .where('is_active', true)
      .select('*');
    
    if (filters.category) {
      query = query.where('category', filters.category);
    }
    
    if (filters.difficulty) {
      query = query.where('difficulty', filters.difficulty);
    }
    
    if (filters.type) {
      query = query.where('type', filters.type);
    }
    
    return query.orderBy('created_at', 'desc');
  }

  static async findById(id) {
    const quiz = await db('quizzes').where({ id, is_active: true }).first();
    if (quiz && quiz.options) {
      quiz.options = JSON.parse(quiz.options);
    }
    return quiz;
  }

  static async getRandomQuiz(filters = {}) {
    let query = db('quizzes').where('is_active', true);
    
    if (filters.category) {
      query = query.where('category', filters.category);
    }
    
    if (filters.difficulty) {
      query = query.where('difficulty', filters.difficulty);
    }

    if (filters.type) {
      query = query.where('type', filters.type);
    }
    
    const quiz = await query.orderByRaw('RANDOM()').first();
    if (quiz && quiz.options) {
      quiz.options = JSON.parse(quiz.options);
    }
    return quiz;
  }

  static async create(quizData) {
    // 如果有选项，转换为 JSON 字符串
    if (quizData.options && Array.isArray(quizData.options)) {
      quizData.options = JSON.stringify(quizData.options);
    }
    
    const [quiz] = await db('quizzes').insert({
      ...quizData,
      created_at: new Date(),
      updated_at: new Date()
    }).returning('*');
    
    if (quiz.options) {
      quiz.options = JSON.parse(quiz.options);
    }
    return quiz;
  }

  static async update(id, quizData) {
    // 如果有选项，转换为 JSON 字符串
    if (quizData.options && Array.isArray(quizData.options)) {
      quizData.options = JSON.stringify(quizData.options);
    }
    
    const [quiz] = await db('quizzes')
      .where({ id })
      .update({ ...quizData, updated_at: new Date() })
      .returning('*');
    
    if (quiz && quiz.options) {
      quiz.options = JSON.parse(quiz.options);
    }
    return quiz;
  }

  static async delete(id) {
    // 软删除
    return db('quizzes').where({ id }).update({ is_active: false });
  }

  // 记录用户答题
  static async recordAttempt(userId, quizId, userAnswer, isCorrect, timeSpent) {
    const attemptData = {
      user_id: userId,
      quiz_id: quizId,
      user_answer: typeof userAnswer === 'object' ? JSON.stringify(userAnswer) : userAnswer,
      is_correct: isCorrect,
      score: isCorrect ? 1 : 0,
      time_spent: timeSpent,
      created_at: new Date(),
      updated_at: new Date()
    };

    const [attempt] = await db('quiz_attempts').insert(attemptData).returning('*');
    return attempt;
  }

  // 获取题目统计
  static async getStats(id) {
    const stats = await db('quiz_attempts')
      .where('quiz_id', id)
      .select(
        db.raw('COUNT(*) as total_attempts'),
        db.raw('SUM(CASE WHEN is_correct = true THEN 1 ELSE 0 END) as correct_attempts'),
        db.raw('AVG(time_spent) as avg_time_spent')
      )
      .first();

    return {
      totalAttempts: parseInt(stats.total_attempts) || 0,
      correctAttempts: parseInt(stats.correct_attempts) || 0,
      accuracy: stats.total_attempts > 0 
        ? ((stats.correct_attempts / stats.total_attempts) * 100).toFixed(2) 
        : 0,
      avgTimeSpent: Math.round(stats.avg_time_spent || 0)
    };
  }

  // 获取分类列表
  static async getCategories() {
    return db('quiz_categories').where('is_active', true).orderBy('name');
  }

  // 按分类获取题目数量
  static async getCategoryStats() {
    return db('quizzes')
      .where('is_active', true)
      .groupBy('category')
      .select('category')
      .count('* as count')
      .orderBy('category');
  }
}

module.exports = Quiz;
