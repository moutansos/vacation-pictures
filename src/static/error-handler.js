window.addEventListener('error', function(e) {
    // Log the error to the server

    const error = `${e.message} at 
    ${e.filename}:${e.lineno}:${e.colno}
    ${e.error.stack}
    `;
    fetch('/api/log', {
        method: 'POST',
        body: JSON.stringify([{
            message: error,
            level: 'error',
        }]),
        headers: {
            'Content-Type': 'application/json',
        },
    });
});


