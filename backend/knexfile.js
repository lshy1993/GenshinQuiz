/**
 * @type { Object.<string, import("knex").Knex.Config> }
 */
module.exports = {
  // 本地开发环境使用 MySQL
  local: {
    client: 'mysql2',
    connection: {
      host: 'localhost',
      port: 3306,
      database: 'genshinquiz',
      user: 'root',
      password: 'your_mysql_password' // 请替换为你的 MySQL 密码
    },
    migrations: {
      directory: './db/migrations'
    },
    seeds: {
      directory: './db/seeds'
    }
  },
  
  // Docker 开发环境使用 PostgreSQL
  development: {
    client: 'postgresql',
    connection: {
      host: process.env.DB_HOST || 'localhost',
      port: process.env.DB_PORT || 5432,
      database: process.env.DB_NAME || 'genshinquiz',
      user: process.env.DB_USER || 'postgres',
      password: process.env.DB_PASSWORD || 'password'
    },
    migrations: {
      directory: './db/migrations'
    },
    seeds: {
      directory: './db/seeds'
    }
  },

  production: {
    client: 'postgresql',
    connection: {
      host: process.env.DB_HOST,
      port: process.env.DB_PORT,
      database: process.env.DB_NAME,
      user: process.env.DB_USER,
      password: process.env.DB_PASSWORD,
      ssl: process.env.DB_SSL ? { rejectUnauthorized: false } : false
    },
    migrations: {
      directory: './db/migrations'
    },
    seeds: {
      directory: './db/seeds'
    },
    pool: {
      min: 2,
      max: 10
    }
  }
};
