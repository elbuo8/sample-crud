#!/usr/bin/env node
const program = require('commander');
const request = require('request');
const querystring = require('querystring');

const API_HOST = process.env.SAMPLE_CRUD_API_HOST || 'http://localhost:3000'; // Use actual domain

/*
cli [action] [resource]
*/

const availableResources = new Set(['models']);
function isResourceAvailable(r) {
  return availableResources.has(r);
}

// Filters and Sort and Limit

program
  .command('get <resource>')
  .option('-i, --id [id]', 'ID of the resource')
  .option('-n, --name [name]', 'Model name')
  .option('--sortBy [sortBy]', 'Field to sort by')
  .option('--limit [limit]', 'Max number of results', parseInt)
  .action((resource, cmd) => {
    if (!isResourceAvailable(resource)) {
      return console.log('Unsupported Resource');
    }
    let resourceURL = `${API_HOST}/${resource}`;
    if (cmd.id) {
      resourceURL += `/${cmd.id}`;
    }
    resourceURL += '?';
    const qs = {};
    if (cmd.name) {
      qs.name = cmd.name;
    }
    if (cmd.sortBy) {
      qs.sortBy = cmd.sortBy;
    }
    if (cmd.limit) {
      qs.total = cmd.limit;
    }
    resourceURL += querystring.stringify(qs);
    request.get(resourceURL, (err, res, body) => {
      console.log(err || body);
    });
  });

program
  .command('create <resource> <name>')
  .option('-a, --accuracy [accuracy]', 'Accuracy of the model', parseFloat)
  .action((resource, name, cmd) => {
    if (!isResourceAvailable(resource)) {
      return console.log('Unsupported Resource');
    }
    let resourceURL = `${API_HOST}/${resource}`;
    const payload = { name };
    if (cmd.accuracy) {
      payload.accuracy = cmd.accuracy;
    }
    request.put(resourceURL, { json: payload }, (err, res, body) => {
      console.log(err || body);
    });
  });

program
  .command('update <resource> <id>')
  .option('-a, --accuracy [accuracy]', 'Accuracy of the model', parseFloat)
  .option('-n, --name [name]', 'Name of the model')
  .action((resource, id, cmd) => {
    if (!isResourceAvailable(resource)) {
      return console.log('Unsupported Resource');
    }
    let resourceURL = `${API_HOST}/${resource}/${id}`;
    const payload = {};
    if (cmd.accuracy) {
      payload.accuracy = cmd.accuracy;
    }
    if (cmd.name) {
      payload.name = cmd.name;
    }
    request.post(resourceURL, { json: payload }, (err, res, body) => {
      console.log(err || body);
    });
  });

program
  .command('delete <resource> <id>')
  .action((resource, id, cmd) => {
    if (!isResourceAvailable(resource)) {
      return console.log('Unsupported Resource');
    }
    let resourceURL = `${API_HOST}/${resource}/${id}`;
    request.delete(resourceURL, (err, res, body) => {
      console.log(err || body);
    });
  });

program.parse(process.argv);
