{
  const sockets = new Set();
  Array(10)
    .fill(1)
    .forEach(() => {
      sockets.add(new WebSocket("https://gomptest.d.roddtech.ru/ws"));
    });
}

{
  const sockets = new Set();
  Array(32)
    .fill(1)
    .forEach(() => {
      sockets.add(new WebSocket("http://localhost:3000/ws"));
    });
}
