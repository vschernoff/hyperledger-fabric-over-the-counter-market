import io from 'socket.io-client';

const TRANSACTION_TYPE = 'ENDORSER_TRANSACTION';

export const socketService = {
  subscribe
};

function subscribe(callback) {
  let host = window.location.hostname;
  const PORT = process.env.PORT || 4000;
  let socket = io.connect(host + ':' + PORT);

  socket.on('connect', function () {
    console.log('connect socket');
    socket.emit('listen_channel', {
      channelId: 'common',
      fullBlock: true
    });
  });

  socket.on('chainblock', function (data) {
    if (checkTransactionType(data.data.data[0].payload.header.channel_header.typeString)) {
      callback();
    }
  });

}

function checkTransactionType(type = '') {
  return type === TRANSACTION_TYPE;
}
