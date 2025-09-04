/**
 * @param { import("knex").Knex } knex
 * @returns { Promise<void> } 
 */
exports.seed = async function(knex) {
  // 清空现有数据
  await knex('quizzes').del();
  await knex('users').del();

  // 插入用户数据
  await knex('users').insert([
    { name: 'John Doe', email: 'john.doe@example.com' },
    { name: '旅行者', email: 'traveler@teyvat.com' },
    { name: 'Genshin Player', email: 'player@genshin.com' }
  ]);

  // 插入原神问答数据
  await knex('quizzes').insert([
    {
      question: 'What is the name of the main currency in Genshin Impact?',
      answer: 'Mora',
      category: 'Currency',
      difficulty: 'easy',
      explanation: 'Mora is the primary currency used throughout Teyvat for various transactions.'
    },
    {
      question: '原神中的七神分别对应哪七种元素？',
      answer: '风、岩、雷、草、水、火、冰',
      category: '设定',
      difficulty: 'medium',
      explanation: '提瓦特大陆的七神分别掌管风、岩、雷、草、水、火、冰七种元素力量。'
    },
    {
      question: 'Which character is known as the "Darknight Hero" of Mondstadt?',
      answer: 'Diluc',
      category: 'Characters',
      difficulty: 'easy',
      explanation: 'Diluc Ragnvindr secretly protects Mondstadt as the Darknight Hero during the night.'
    },
    {
      question: '璃月港的守护神是谁？',
      answer: '钟离（摩拉克斯）',
      category: '角色',
      difficulty: 'easy',
      explanation: '钟离是璃月的岩王帝君摩拉克斯的人间化身，是璃月港的守护神。'
    },
    {
      question: 'What is the maximum Adventure Rank in Genshin Impact?',
      answer: '60',
      category: 'Game Mechanics',
      difficulty: 'medium',
      explanation: 'The current maximum Adventure Rank is 60, though players can continue gaining AR EXP.'
    }
  ]);
};
