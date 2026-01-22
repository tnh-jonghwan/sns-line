// 현재 활성 플랫폼
let currentPlatform = 'line';

// eventHub 연결
let eventSource = null;
const messagesLine = document.getElementById('messagesLine');
const messagesInstagram = document.getElementById('messagesInstagram');
const messageInput = document.getElementById('messageInput');
const sendButton = document.getElementById('sendButton');
const statusDot = document.getElementById('statusDot');
const statusText = document.getElementById('statusText');
const tabs = document.querySelectorAll('.tab');
const messageContainers = document.querySelectorAll('.messages-container');

// 탭 전환
tabs.forEach(tab => {
    tab.addEventListener('click', () => {
        const platform = tab.dataset.platform;
        switchPlatform(platform);
    });
});

function switchPlatform(platform) {
    currentPlatform = platform;
    
    // 탭 활성화
    tabs.forEach(tab => {
        if (tab.dataset.platform === platform) {
            tab.classList.add('active');
        } else {
            tab.classList.remove('active');
        }
    });
    
    // 메시지 컨테이너 전환
    messageContainers.forEach(container => {
        if (container.dataset.platform === platform) {
            container.classList.add('active');
        } else {
            container.classList.remove('active');
        }
    });
}

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
        // userId로 플랫폼 구분 (간단한 방법: Instagram ID는 보통 더 길음)
        // 더 정확한 방법은 서버에서 platform 정보를 함께 보내는 것
        const platform = detectPlatform(data.userId);
        addMessage(data.text, data.userId, 'received', platform);
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

// 플랫폼 감지 (임시 - 서버에서 보내주면 더 정확함)
function detectPlatform(userId) {
    // Instagram IGSID는 보통 매우 긴 숫자
    // LINE userId는 보통 U로 시작하거나 다른 형식
    if (userId && userId.length > 15 && /^\d+$/.test(userId)) {
        return 'instagram';
    }
    return 'line';
}

// 메시지 추가
function addMessage(text, userId, type = 'received', platform = 'line') {
    const messageDiv = document.createElement('div');
    messageDiv.className = `message ${type}`;
    
    const bubble = document.createElement('div');
    bubble.className = 'message-bubble';
    
    // 수신 메시지에 sender ID 표시
    if (type === 'received' && userId) {
        const senderInfo = document.createElement('div');
        senderInfo.className = 'sender-id';
        senderInfo.textContent = userId;
        bubble.appendChild(senderInfo);
    }
    
    const textContent = document.createElement('div');
    textContent.textContent = text;
    bubble.appendChild(textContent);
    
    const time = document.createElement('div');
    time.className = 'message-time';
    time.textContent = new Date().toLocaleTimeString('ko-KR', { 
        hour: '2-digit', 
        minute: '2-digit' 
    });
    
    messageDiv.appendChild(bubble);
    messageDiv.appendChild(time);
    
    // 플랫폼에 맞는 컨테이너에 추가
    const targetContainer = platform === 'instagram' ? messagesInstagram : messagesLine;
    targetContainer.appendChild(messageDiv);
    targetContainer.scrollTop = targetContainer.scrollHeight;
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
            addMessage(text, 'me', 'sent', currentPlatform);
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
