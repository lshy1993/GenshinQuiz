const db = require('../db');

class User {
  static async findAll() {
    return db('users').select('id', 'name', 'email', 'role', 'avatar_url', 'is_active', 'created_at');
  }

  static async findById(id) {
    return db('users').where({ id }).select('id', 'name', 'email', 'role', 'avatar_url', 'is_active', 'created_at').first();
  }

  static async findByEmail(email) {
    return db('users').where({ email }).first();
  }

  static async create(userData) {
    const [user] = await db('users').insert({
      ...userData,
      created_at: new Date(),
      updated_at: new Date()
    }).returning(['id', 'name', 'email', 'role', 'avatar_url', 'is_active', 'created_at']);
    return user;
  }

  static async update(id, userData) {
    const [user] = await db('users')
      .where({ id })
      .update({ ...userData, updated_at: new Date() })
      .returning(['id', 'name', 'email', 'role', 'avatar_url', 'is_active', 'created_at']);
    return user;
  }

  static async updateLastLogin(id) {
    return db('users')
      .where({ id })
      .update({ last_login_at: new Date() });
  }

  static async delete(id) {
    return db('users').where({ id }).update({ is_active: false }); // 软删除
  }

  // 获取用户的投票历史
  static async getVoteHistory(userId) {
    return db('user_votes as uv')
      .join('votes as v', 'uv.vote_id', 'v.id')
      .join('vote_options as vo', 'uv.vote_option_id', 'vo.id')
      .where('uv.user_id', userId)
      .select(
        'v.id as vote_id',
        'v.title as vote_title',
        'vo.title as option_title',
        'uv.created_at as voted_at'
      )
      .orderBy('uv.created_at', 'desc');
  }

  // 获取用户的答题记录
  static async getQuizHistory(userId) {
    return db('quiz_attempts as qa')
      .join('quizzes as q', 'qa.quiz_id', 'q.id')
      .where('qa.user_id', userId)
      .select(
        'q.id as quiz_id',
        'q.question',
        'q.category',
        'qa.is_correct',
        'qa.score',
        'qa.time_spent',
        'qa.created_at as attempted_at'
      )
      .orderBy('qa.created_at', 'desc');
  }

  // 检查用户权限
  static async isAdmin(userId) {
    const user = await db('users').where({ id: userId }).first();
    return user && user.role === 'admin';
  }
}

module.exports = User;
