import * as apiService from './api.service';
import {configService} from './config.service';

export const dealService = {
  getAll,
  edit,
  history
};

function getAll() {
  return apiService.query(
    apiService.channels.common,
    apiService.contracts.reference,
    'queryDeals',
    `[]`);
}

function edit(product) {
}

function history(product) {
}
