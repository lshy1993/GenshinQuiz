const db = require('../db');

class User {
  static async findAll() {
    return db('users').select('*');
  }

  static async findById(id) {
    return db('users').where({ id }).first();
  }

  static async create(userData) {
    const [user] = await db('users').insert(userData).returning('*');
    return user;
  }

  static async update(id, userData) {
    const [user] = await db('users')
      .where({ id })
      .update({ ...userData, updated_at: new Date() })
      .returning('*');
    return user;
  }

  static async delete(id) {
    return db('users').where({ id }).del();
  }
}

module.exports = User;
