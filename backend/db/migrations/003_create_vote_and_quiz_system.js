/**
 * @param { import("knex").Knex } knex
 * @returns { Promise<void> }
 */
exports.up = function(knex) {
  return knex.schema
    .alterTable('quizzes', function (table) {
      // 增强题目表
      table.integer('created_by').unsigned().references('id').inTable('users').after('explanation'); // 题目创建者
      table.boolean('is_active').defaultTo(true).after('created_by'); // 题目状态
      table.json('options').after('answer'); // 选择题选项 JSON 格式
      table.enum('type', ['single_choice', 'multiple_choice', 'true_false', 'text']).defaultTo('single_choice').after('type');
    })
    
    // 创建投票系统表
    .createTable('votes', function (table) {
      table.increments('id').primary();
      table.string('title', 255).notNullable(); // 投票标题
      table.text('description'); // 投票描述
      table.enum('type', ['single_choice', 'multiple_choice']).defaultTo('single_choice'); // 投票类型
      table.integer('created_by').unsigned().references('id').inTable('users'); // 创建者
      table.boolean('is_active').defaultTo(true); // 投票状态
      table.boolean('is_anonymous').defaultTo(false); // 是否匿名投票
      table.timestamp('start_time').notNullable(); // 开始时间
      table.timestamp('end_time'); // 结束时间（可为空表示无限期）
      table.integer('max_choices').defaultTo(1); // 最大选择数量
      table.timestamps(true, true);
    })
    
    // 创建投票选项表
    .createTable('vote_options', function (table) {
      table.increments('id').primary();
      table.integer('vote_id').unsigned().references('id').inTable('votes').onDelete('CASCADE');
      table.string('title', 255).notNullable(); // 选项标题
      table.text('description'); // 选项描述
      table.string('image_url', 500); // 选项图片
      table.integer('sort_order').defaultTo(0); // 排序
      table.timestamps(true, true);
    })
    
    // 创建用户投票记录表
    .createTable('user_votes', function (table) {
      table.increments('id').primary();
      table.integer('vote_id').unsigned().references('id').inTable('votes').onDelete('CASCADE');
      table.integer('user_id').unsigned().references('id').inTable('users').onDelete('CASCADE');
      table.integer('vote_option_id').unsigned().references('id').inTable('vote_options').onDelete('CASCADE');
      table.timestamps(true, true);
      
      // 联合唯一索引：确保用户对每个投票的每个选项只能投一次
      table.unique(['vote_id', 'user_id', 'vote_option_id']);
    })
    
    // 创建问答记录表
    .createTable('quiz_attempts', function (table) {
      table.increments('id').primary();
      table.integer('user_id').unsigned().references('id').inTable('users').onDelete('CASCADE');
      table.integer('quiz_id').unsigned().references('id').inTable('quizzes').onDelete('CASCADE');
      table.json('user_answer'); // 用户答案（JSON格式支持多种题型）
      table.boolean('is_correct'); // 是否正确
      table.integer('score').defaultTo(0); // 得分
      table.integer('time_spent'); // 答题用时（秒）
      table.timestamps(true, true);
    })
    
    // 创建题目分类表
    .createTable('quiz_categories', function (table) {
      table.increments('id').primary();
      table.string('name', 100).notNullable(); // 分类名称
      table.string('description', 255); // 分类描述
      table.string('icon_url', 500); // 分类图标
      table.boolean('is_active').defaultTo(true);
      table.timestamps(true, true);
    });
};

/**
 * @param { import("knex").Knex } knex
 * @returns { Promise<void> }
 */
exports.down = function(knex) {
  return knex.schema
    .dropTable('quiz_categories')
    .dropTable('quiz_attempts')
    .dropTable('user_votes')
    .dropTable('vote_options')
    .dropTable('votes')
    .alterTable('quizzes', function (table) {
      table.dropColumn('created_by');
      table.dropColumn('is_active');
      table.dropColumn('options');
      table.dropColumn('type');
    });
};
