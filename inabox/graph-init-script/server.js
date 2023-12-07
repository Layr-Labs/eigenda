const http = require('http');

const server = http.createServer((req, res) => {
    if (req.url === '/health') {
        res.writeHead(200);
        res.end('OK');
    } else {
        res.writeHead(404);
        res.end('Not Found');
    }
});

server.listen(8080, () => {
    console.log('Server running on port 8080');
});