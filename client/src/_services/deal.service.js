import * as apiService from './api.service';
import {configService} from './config.service';

export const dealService = {
  getAll
};

const ACTIONS = {
  getAll: 'queryDeals'
};

const CHAINCODES = {
  reference: 'reference'
};

const CHANNELS = {
  common: 'common'
};

async function getAll() {
  const channels = await configService.getChannels();
  const chaincodes = await configService.getChaincodes();
  return await apiService.query(
    channels[CHANNELS.common],
    chaincodes[CHAINCODES.reference],
    ACTIONS.getAll,
    `[]`);
}

