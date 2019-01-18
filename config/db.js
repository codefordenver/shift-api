const Sequelize = require('sequelize');
const env = require('./env.json');
const conn = new Sequelize( env.database, env.username, env.password, {
    dialect: env.dialect,
    host: env.host,
    port: env.port,
    operatorsAliases: false
});

module.exports = {
    'Sequelize': Sequelize,
    'conn': conn
}