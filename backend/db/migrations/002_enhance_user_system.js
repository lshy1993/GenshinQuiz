/**
 * @param { import("knex").Knex } knex
 * @returns { Promise<void> }
 */
exports.up = function(knex) {
  return knex.schema
    .alterTable('users', function (table) {
      table.string('password', 255).notNullable().after('email'); // 密码（应该加密）
      table.enum('role', ['user', 'admin']).defaultTo('user').after('password'); // 用户角色
      table.string('avatar_url', 500).after('role'); // 头像URL
      table.boolean('is_active').defaultTo(true).after('avatar_url'); // 账户状态
      table.timestamp('last_login_at').after('is_active'); // 最后登录时间
    });
};

/**
 * @param { import("knex").Knex } knex
 * @returns { Promise<void> }
 */
exports.down = function(knex) {
  return knex.schema
    .alterTable('users', function (table) {
      table.dropColumn('password');
      table.dropColumn('role');
      table.dropColumn('avatar_url');
      table.dropColumn('is_active');
      table.dropColumn('last_login_at');
    });
};
