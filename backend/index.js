require('dotenv').config();
const express = require('express');
const { Pool } = require('pg');

const app = express();
const port = 3082;

const pool = new Pool({
  user: process.env.DB_USER,
  password: process.env.DB_PASSWORD,
  host: process.env.DB_HOST,
  port: process.env.DB_PORT,
  database: process.env.DB_NAME,
});

app.get('/', async (req, res) => {
  const result = await pool.query('SELECT NOW()');
  res.send(`Database time: ${result.rows[0].now}`);
});

app.listen(port, () => {
  console.log(`Backend running at http://localhost:${port}`);
});