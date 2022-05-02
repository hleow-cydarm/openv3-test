const path = require('path');
const { generate } = require('openapi-typescript-validator');

generate({
  schemaFile: path.join(__dirname, '../spec/tictactoe.yaml'),
  schemaType: 'yaml',
  directory: path.join(__dirname, '/generated')
})