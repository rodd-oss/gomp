{
  const sockets = new Set();
  Array(100)
    .fill(1)
    .forEach(() => {
      sockets.add(new WebSocket("https://gomptest.d.roddtech.ru/ws"));
    });
}
