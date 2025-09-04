/**
 * @param { import("knex").Knex } knex
 * @returns { Promise<void> }
 */
exports.up = function(knex) {
  return knex.schema
    .createTable('users', function (table) {
      table.increments('id').primary();
      table.string('name', 100);
      table.string('email', 100).unique();
      table.timestamps(true, true); // created_at, updated_at
    })
    .createTable('quizzes', function (table) {
      table.increments('id').primary();
      table.string('question', 255).notNullable();
      table.string('answer', 255).notNullable();
      table.string('category', 100); // 题目分类，如：角色、武器、地区等
      table.enum('difficulty', ['easy', 'medium', 'hard']).defaultTo('easy');
      table.text('explanation'); // 答案解释
      table.timestamps(true, true);
    });
};

/**
 * @param { import("knex").Knex } knex
 * @returns { Promise<void> }
 */
exports.down = function(knex) {
  return knex.schema
    .dropTable('quizzes')
    .dropTable('users');
};
