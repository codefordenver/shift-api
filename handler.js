const Sequelize = require('sequelize');
const db = require('./config/db.js');

//need a function to test connection and return error message if fails
pgconn = db.conn;
dataType = db.Sequelize;

const Model = pgconn.define('model', {
    geoid: {
        type: dataType.STRING,
        field: 'geoid10'
    },
    ttlpop: {
        type:dataType.INTEGER,
        field: 'b01001_001e'
    }
});

pgconn.query("SELECT geoid10, b01001_001e FROM acs5.county_state_b01001_2016 WHERE geoid10='08001'", {
    type: pgconn.QueryTypes.SELECT,
    model: Model,
    mapToModel: true
}).then(rst => {
    console.log(rst)
});

module.exports.listDatasets = async (event, context, callback) => {
    const textResponseHeaders = {
        'Content-Type': 'text/plain'
    };

    const jsonResponseHeaders = {
        'Content-Type': 'application/json'
    };

    db.conn.query("SELECT geoid10, b01001_001e FROM acs5.county_state_b01001_2016 WHERE geoid10='08001'", {
        type: sequelize.QueryTypes.SELECT,
        model: datasets(db.conn, db.Sequelize),
        mapToModel: true,

    }).then(mySchemas => {
        console.log(mySchemas);
        const response = {
            statusCode: 200,
            headers: jsonResponseHeaders,
            body: JSON.stringify(mySchemas)
        };
        callback(null, response);
    })
    .catch(error => {
        callback(null, {
            statusCode: 501,
            headers: jsonResponseHeaders,
            body: error
        });
    });
}