require('dotenv').config();
const express = require('express');
const setupSwagger = require('./swagger');
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

// 在应用中设置 Swagger
setupSwagger(app);

// 示例路由
/**
 * @swagger
 * /:
 *   get:
 *     description: Returns the current time from the database
 *     responses:
 *       200:
 *         description: Current time
 */
app.get('/', async (req, res) => {
  const result = await pool.query('SELECT NOW()');
  res.send(`Database time: ${result.rows[0].now}`);
});

app.get('/api/v1/example', (req, res) => {
  res.send({ message: 'Hello World' });
});

app.listen(port, () => {
  console.log(`Server is running on http://localhost:${port}`);
  console.log(`Swagger docs available at http://localhost:${port}/api-docs`);
});