const db = require('../db');

class Vote {
  // 获取所有激活的投票
  static async findAllActive() {
    return await db('votes')
      .where('is_active', true)
      .orderBy('created_at', 'desc');
  }

  // 根据ID获取投票详情（包含选项）
  static async findByIdWithOptions(id) {
    const vote = await db('votes').where('id', id).first();
    if (!vote) return null;

    const options = await db('vote_options')
      .where('vote_id', id)
      .orderBy('sort_order');

    return {
      ...vote,
      options
    };
  }

  // 创建新投票
  static async create(voteData, options) {
    const trx = await db.transaction();
    try {
      // 插入投票
      const [voteId] = await trx('votes').insert(voteData).returning('id');
      
      // 插入投票选项
      const optionsWithVoteId = options.map((option, index) => ({
        ...option,
        vote_id: voteId,
        sort_order: index + 1
      }));
      
      await trx('vote_options').insert(optionsWithVoteId);
      
      await trx.commit();
      return this.findByIdWithOptions(voteId);
    } catch (error) {
      await trx.rollback();
      throw error;
    }
  }

  // 用户投票
  static async submitVote(voteId, userId, optionIds) {
    const trx = await db.transaction();
    try {
      // 检查用户是否已经投票
      const existingVotes = await trx('user_votes')
        .where({ vote_id: voteId, user_id: userId });

      if (existingVotes.length > 0) {
        throw new Error('User has already voted');
      }

      // 检查投票是否仍然激活
      const vote = await trx('votes').where('id', voteId).first();
      if (!vote || !vote.is_active) {
        throw new Error('Vote is not active');
      }

      // 检查是否在投票时间范围内
      const now = new Date();
      if (vote.end_time && new Date(vote.end_time) < now) {
        throw new Error('Vote has ended');
      }

      // 检查选择数量是否超限
      if (optionIds.length > vote.max_choices) {
        throw new Error('Too many choices selected');
      }

      // 插入投票记录
      const voteRecords = optionIds.map(optionId => ({
        vote_id: voteId,
        user_id: userId,
        vote_option_id: optionId
      }));

      await trx('user_votes').insert(voteRecords);
      await trx.commit();
      
      return true;
    } catch (error) {
      await trx.rollback();
      throw error;
    }
  }

  // 获取投票结果统计
  static async getResults(voteId) {
    const results = await db('vote_options as vo')
      .leftJoin('user_votes as uv', 'vo.id', 'uv.vote_option_id')
      .where('vo.vote_id', voteId)
      .groupBy('vo.id', 'vo.title')
      .select(
        'vo.id',
        'vo.title',
        'vo.description',
        db.raw('COUNT(uv.id) as vote_count')
      )
      .orderBy('vo.sort_order');

    // 计算总票数
    const totalVotes = results.reduce((sum, result) => sum + parseInt(result.vote_count), 0);

    return {
      results: results.map(result => ({
        ...result,
        percentage: totalVotes > 0 ? ((result.vote_count / totalVotes) * 100).toFixed(2) : 0
      })),
      totalVotes
    };
  }

  // 检查用户是否已投票
  static async hasUserVoted(voteId, userId) {
    const vote = await db('user_votes')
      .where({ vote_id: voteId, user_id: userId })
      .first();
    return !!vote;
  }
}

module.exports = Vote;
