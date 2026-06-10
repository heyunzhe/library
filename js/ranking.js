document.addEventListener('DOMContentLoaded', function() {
    const list = document.getElementById('rankingList');

    const token = localStorage.getItem('access_token');
    const headers = token ? { 'Authorization': 'Bearer ' + token } : {};
    const options = { method: 'POST', headers: headers };
    fetch('/ranking', options)
    .then(r => r.json())
    .then(data => {
        data.sort((a, b) => b.count - a.count);

        if (data.length === 0) {
            list.innerHTML = '<div style="text-align:center;padding:40px;color:#95a5a6">暂无借阅数据</div>';
            return;
        }

        data.forEach((book, index) => {
            const rankClass = index < 3 ? `rank-${index + 1}` : '';
            const coverUrl = book.cover ? '/' + book.cover : '/images/default-avatar.svg';

            const item = document.createElement('div');
            item.className = 'ranking-item';
            item.innerHTML = `
                <div class="rank ${rankClass}">${index + 1}</div>
                <img src="${coverUrl}" alt="${book.title}" class="book-cover" onerror="this.src='/images/default-avatar.svg'">
                <div class="book-info">
                    <div class="book-title">${book.title}</div>
                    <div class="book-meta">ISBN: ${book.isbn}</div>
                </div>
                <div class="borrow-count">${book.count} <small>次借阅</small></div>
            `;
            list.appendChild(item);
        });
    })
    .catch(() => {
        list.innerHTML = '<div style="text-align:center;padding:40px;color:#95a5a6">加载失败，请稍后再试</div>';
    });

    // 认证按钮
    const btn = document.getElementById('authButton');
    if (btn) {
        if (localStorage.getItem('isLoggedIn') === 'true') {
            btn.textContent = '个人中心';
            btn.href = '/user/library';
        } else {
            btn.href = '/login?redirect=' + encodeURIComponent(window.location.pathname);
        }
    }
});