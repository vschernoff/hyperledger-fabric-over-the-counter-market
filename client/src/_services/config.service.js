import {config, getChannels, getChaincodes} from './api.service';

export const configService = {
  load,
  get,
  update,
  getPeers,
  getChannels: loadChannels,
  getChaincodes: loadChaincodes
};

let localConfig = {};
let isFetching = false;
let reqPromise;

function get() {
  return localConfig;
}

function load() {
  if (isFetching) {
    return reqPromise;
  }
  isFetching = true;

  reqPromise = config()
    .then(res => {
      localConfig = res;
      return localConfig;
    });

  return reqPromise;
}

// load secured information, like chaincodes, channels
function update() {
  return Promise.all([
    loadChannels(),
    loadChaincodes()
  ]);
}

function _extend(asyncFn, cfgKey, responseKey) {
  if (!localConfig[cfgKey]) {
    if (asyncFn.promise) {
      return asyncFn.promise;
    }
    asyncFn.promise = asyncFn()
      .then(response => {
        localConfig[cfgKey] = {};
        response[responseKey].forEach(obj => {
          localConfig[cfgKey][obj.name] = obj.name;
        });

        delete asyncFn.promise;
        return localConfig[cfgKey];
      });

    return asyncFn.promise;
  }

  return Promise.resolve(localConfig[cfgKey]);
}

function loadChannels() {
  return _extend(getChannels, 'channels', 'channels');
}

function loadChaincodes() {
  return _extend(getChaincodes, 'chaincodes', 'chaincodes');
}

function getPeers() {
  if (localConfig.peers) {
    return localConfig.peers;
  }
  const orgConfig = localConfig['network-config'][localConfig.org];
  localConfig.peers = Object.keys(orgConfig).filter(k => k.startsWith('peer'));
  return localConfig.peers;
}
