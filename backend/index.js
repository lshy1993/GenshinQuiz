const express = require('express');
const { Pool } = require('pg');

const app = express();
const port = 5000;

// 配置数据库连接
const pool = new Pool({
  connectionString: process.env.DATABASE_URL,
});

app.get('/', async (req, res) => {
  const result = await pool.query('SELECT NOW()');
  res.send(`Database time: ${result.rows[0].now}`);
});

app.listen(port, () => {
  console.log(`Backend running at http://localhost:${port}`);
});