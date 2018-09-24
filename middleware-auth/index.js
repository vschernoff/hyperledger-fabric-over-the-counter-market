module.exports = function auth() {
  return Promise.resolve();
};

/* example
const credentialStore = [{
  username: 'test',
  password: undefined,
  org: 'a'
}];

module.exports = function auth(username, password, orgname) {
  const user = credentialStore.find(obj => obj.username === username && obj.org === orgname && obj.password === password);
  return user ? Promise.resolve(user) : Promise.reject({status: 401, message: 'invalid credentials'});
};
*/
