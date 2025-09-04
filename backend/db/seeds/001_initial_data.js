/**
 * @param { import("knex").Knex } knex
 * @returns { Promise<void> } 
 */
exports.seed = async function(knex) {
  // 清空现有数据（注意顺序，先删除有外键的表）
  await knex('user_votes').del();
  await knex('quiz_attempts').del();
  await knex('vote_options').del();
  await knex('votes').del();
  await knex('quizzes').del();
  await knex('quiz_categories').del();
  await knex('users').del();

  // 插入用户数据
  const users = await knex('users').insert([
    { 
      name: 'Admin User', 
      email: 'admin@genshin.com',
      password: '$2b$10$hash_password_here', // 实际使用时需要用 bcrypt 加密
      role: 'admin'
    },
    { 
      name: '旅行者', 
      email: 'traveler@teyvat.com',
      password: '$2b$10$hash_password_here',
      role: 'user'
    },
    { 
      name: 'Genshin Player', 
      email: 'player@genshin.com',
      password: '$2b$10$hash_password_here',
      role: 'user'
    }
  ]).returning('id');

  // 插入题目分类
  const categories = await knex('quiz_categories').insert([
    { name: '角色', description: '关于原神角色的问题', icon_url: '/icons/character.svg' },
    { name: '武器', description: '关于武器系统的问题', icon_url: '/icons/weapon.svg' },
    { name: '设定', description: '关于游戏世界观的问题', icon_url: '/icons/lore.svg' },
    { name: '元素', description: '关于元素系统的问题', icon_url: '/icons/element.svg' }
  ]).returning('id');

  // 插入原神问答数据
  await knex('quizzes').insert([
    {
      question: '原神中的主要货币叫什么？',
      answer: 'Mora',
      category: '基础',
      difficulty: 'easy',
      type: 'single_choice',
      options: JSON.stringify(['Mora', '原石', '创世结晶', '星辉']),
      explanation: 'Mora是提瓦特大陆通用的货币，用于各种交易和升级。',
      created_by: 1, // admin用户
      is_active: true
    },
    {
      question: '以下哪个角色属于蒙德地区？',
      answer: '温迪',
      category: '角色',
      difficulty: 'easy',
      type: 'single_choice',
      options: JSON.stringify(['温迪', '钟离', '雷电将军', '纳西妲']),
      explanation: '温迪是蒙德的风神巴巴托斯，也是蒙德地区的代表角色。',
      created_by: 1,
      is_active: true
    },
    {
      question: '原神中有多少种基础元素？',
      answer: '7',
      category: '元素',
      difficulty: 'medium',
      type: 'text',
      explanation: '原神有7种基础元素：风、岩、雷、草、水、火、冰。',
      created_by: 1,
      is_active: true
    }
  ]);

  // 插入投票示例
  const votes = await knex('votes').insert([
    {
      title: '你最喜欢的原神角色是？',
      description: '选择你最喜欢的原神角色',
      type: 'single_choice',
      created_by: 1,
      is_active: true,
      is_anonymous: false,
      start_time: knex.fn.now(),
      end_time: knex.raw("datetime('now', '+30 days')"),
      max_choices: 1
    },
    {
      title: '你希望哪些功能优先更新？',
      description: '多选你最期待的功能更新',
      type: 'multiple_choice',
      created_by: 1,
      is_active: true,
      is_anonymous: true,
      start_time: knex.fn.now(),
      max_choices: 3
    }
  ]).returning('id');

  // 插入投票选项
  await knex('vote_options').insert([
    // 第一个投票的选项
    { vote_id: votes[0].id || 1, title: '温迪', description: '自由的风神', sort_order: 1 },
    { vote_id: votes[0].id || 1, title: '钟离', description: '契约的岩神', sort_order: 2 },
    { vote_id: votes[0].id || 1, title: '雷电将军', description: '永恒的雷神', sort_order: 3 },
    { vote_id: votes[0].id || 1, title: '纳西妲', description: '智慧的草神', sort_order: 4 },
    
    // 第二个投票的选项
    { vote_id: votes[1].id || 2, title: '新角色', description: '更多有趣的角色', sort_order: 1 },
    { vote_id: votes[1].id || 2, title: '新地区', description: '探索新的地区', sort_order: 2 },
    { vote_id: votes[1].id || 2, title: '玩法优化', description: '改善现有玩法', sort_order: 3 },
    { vote_id: votes[1].id || 2, title: '剧情推进', description: '主线剧情发展', sort_order: 4 }
  ]);
};
