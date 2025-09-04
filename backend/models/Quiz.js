const db = require('../db');

class Quiz {
  static async findAll(filters = {}) {
    let query = db('quizzes').select('*');
    
    if (filters.category) {
      query = query.where('category', filters.category);
    }
    
    if (filters.difficulty) {
      query = query.where('difficulty', filters.difficulty);
    }
    
    return query.orderBy('created_at', 'desc');
  }

  static async findById(id) {
    return db('quizzes').where({ id }).first();
  }

  static async getRandomQuiz(filters = {}) {
    let query = db('quizzes');
    
    if (filters.category) {
      query = query.where('category', filters.category);
    }
    
    if (filters.difficulty) {
      query = query.where('difficulty', filters.difficulty);
    }
    
    return query.orderByRaw('RANDOM()').first();
  }

  static async create(quizData) {
    const [quiz] = await db('quizzes').insert(quizData).returning('*');
    return quiz;
  }

  static async update(id, quizData) {
    const [quiz] = await db('quizzes')
      .where({ id })
      .update({ ...quizData, updated_at: new Date() })
      .returning('*');
    return quiz;
  }

  static async delete(id) {
    return db('quizzes').where({ id }).del();
  }

  static async getCategories() {
    const categories = await db('quizzes')
      .distinct('category')
      .whereNotNull('category')
      .orderBy('category');
    return categories.map(row => row.category);
  }

  static async getStats() {
    const [stats] = await db('quizzes')
      .select([
        db.raw('COUNT(*) as total'),
        db.raw('COUNT(CASE WHEN difficulty = ? THEN 1 END) as easy', ['easy']),
        db.raw('COUNT(CASE WHEN difficulty = ? THEN 1 END) as medium', ['medium']),
        db.raw('COUNT(CASE WHEN difficulty = ? THEN 1 END) as hard', ['hard'])
      ]);
    return stats;
  }
}

module.exports = Quiz;
