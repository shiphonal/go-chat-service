document.addEventListener('DOMContentLoaded', function() {
    // Check if we're on the login page
    if (document.getElementById('login-form')) {
        const loginForm = document.getElementById('login-form');

        loginForm.addEventListener('submit', function(e) {
            e.preventDefault();

            const username = document.getElementById('username').value;
            const password = document.getElementById('password').value;

            // Here you would typically make an API call to authenticate
            // For now, we'll just redirect to the chat page
            document.cookie = `token=simulated_token_for_${username}; path=/`;
            window.location.href = '../templates/index.html';
        });
    }

    // Check if we're on the chat page
    if (document.getElementById('logout-btn')) {
        const logoutBtn = document.getElementById('logout-btn');
        const messageForm = document.getElementById('message-form');

        logoutBtn.addEventListener('click', function() {
            document.cookie = 'token=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;';
            window.location.href = '/static/login.html';
        });

        messageForm.addEventListener('submit', async function(e) {
            e.preventDefault();

            const messageType = document.getElementById('message-type').value;
            const content = document.getElementById('message-content').value;

            try {
                const response = await fetch('/messages', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/x-www-form-urlencoded',
                    },
                    body: `text=${encodeURIComponent(messageType)}&content=${encodeURIComponent(content)}`
                });

                if (response.ok) {
                    const result = await response.text();
                    console.log('Message sent:', result);
                    document.getElementById('message-content').value = '';

                    // Add message to the UI
                    const messagesDiv = document.getElementById('messages');
                    const messageDiv = document.createElement('div');
                    messageDiv.className = 'message sent';
                    messageDiv.innerHTML = `
                        <div class="message-info">You (${new Date().toLocaleTimeString()})</div>
                        <div class="message-content">${content}</div>
                    `;
                    messagesDiv.appendChild(messageDiv);
                    messagesDiv.scrollTop = messagesDiv.scrollHeight;
                } else {
                    console.error('Failed to send message');
                }
            } catch (error) {
                console.error('Error:', error);
            }
        });

        // Simulate loading some messages (in a real app, you'd fetch these from the server)
        setTimeout(() => {
            const messagesDiv = document.getElementById('messages');

            const sampleMessages = [
                { type: 'text', content: 'Hello there!', sent: false },
                { type: 'text', content: 'Hi! How are you?', sent: true },
                { type: 'text', content: 'I\'m doing well, thanks for asking!', sent: false }
            ];

            sampleMessages.forEach(msg => {
                const messageDiv = document.createElement('div');
                messageDiv.className = `message ${msg.sent ? 'sent' : 'received'}`;
                messageDiv.innerHTML = `
                    <div class="message-info">${msg.sent ? 'You' : 'Other User'} (${new Date().toLocaleTimeString()})</div>
                    <div class="message-content">${msg.content}</div>
                `;
                messagesDiv.appendChild(messageDiv);
            });

            messagesDiv.scrollTop = messagesDiv.scrollHeight;
        }, 500);
    }
});