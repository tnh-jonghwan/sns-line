// eventHub 연결
let eventSource = null;
const messages = document.getElementById('messages');
const messageInput = document.getElementById('messageInput');
const sendButton = document.getElementById('sendButton');
const statusDot = document.getElementById('statusDot');
const statusText = document.getElementById('statusText');

// eventHub 연결
function connect() {
    eventSource = new EventSource('/events');
    
    eventSource.onopen = () => {
        console.log('eventHub 연결됨');
        statusDot.classList.add('connected');
        statusText.textContent = '연결됨';
        sendButton.disabled = false;
    };
    
    eventSource.onmessage = (event) => {
        const data = JSON.parse(event.data);
        addMessage(data.text, data.userId, 'received');
    };
    
    eventSource.onerror = (error) => {
        console.error('eventHub 에러:', error);
        statusDot.classList.remove('connected');
        statusText.textContent = '연결 끊김';
        sendButton.disabled = true;
        
        // 연결 재시도
        eventSource.close();
        setTimeout(connect, 5000);
    };
    
    // connected 이벤트 처리
    eventSource.addEventListener('connected', (event) => {
        console.log('서버 연결 확인:', event.data);
    });
}

// 메시지 추가
function addMessage(text, userId, type = 'received') {
    const messageDiv = document.createElement('div');
    messageDiv.className = `message ${type}`;
    
    const bubble = document.createElement('div');
    bubble.className = 'message-bubble';
    bubble.textContent = text;
    
    const time = document.createElement('div');
    time.className = 'message-time';
    time.textContent = new Date().toLocaleTimeString('ko-KR', { 
        hour: '2-digit', 
        minute: '2-digit' 
    });
    
    messageDiv.appendChild(bubble);
    messageDiv.appendChild(time);
    
    messages.appendChild(messageDiv);
    messages.scrollTop = messages.scrollHeight;
}

// 메시지 전송
function sendMessage() {
    const text = messageInput.value.trim();
    if (!text) return;
    
    // 서버로 전송
    fetch('/api/send', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ text }),
    })
    .then(response => response.json())
    .then(data => {
        if (data.success) {
            addMessage(text, 'me', 'sent');
            messageInput.value = '';
        } else {
            alert('메시지 전송 실패: ' + data.error);
        }
    })
    .catch(error => {
        console.error('전송 에러:', error);
        alert('메시지 전송 중 오류가 발생했습니다.');
    });
}

// 이벤트 리스너
sendButton.addEventListener('click', sendMessage);
messageInput.addEventListener('keypress', (e) => {
    if (e.key === 'Enter') {
        sendMessage();
    }
});

// 초기 연결
connect();
